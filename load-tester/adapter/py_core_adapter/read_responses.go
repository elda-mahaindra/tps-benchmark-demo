package py_core_adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

// readResponses continuously reads responses from the server and passes them to processResponse
func (adapter *Adapter) readResponses() {
	const op = "py_core_adapter/readResponses"

	for {
		// If not connected, wait and retry
		if !adapter.isConnected {
			time.Sleep(100 * time.Millisecond) // Consider making this configurable or use exponential backoff
			continue
		}

		// Get the connection under mutex (read lock might suffice if conn assignment is atomic, but full lock is safer)
		adapter.connMutex.Lock()
		conn := adapter.conn
		adapter.connMutex.Unlock()

		if conn == nil {
			// Small sleep to avoid busy-waiting if conn is temporarily nil during reconnect
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Read the length of the response JSON data (2 bytes)
		responseLengthBytes := make([]byte, 2)
		_, err := conn.Read(responseLengthBytes)
		if err != nil {
			adapter.logger.WithFields(logrus.Fields{
				"op":  op,
				"err": err.Error(),
			}).Error("Failed to read response length, handling disconnect")

			adapter.handleDisconnect() // This will set isConnected to false and clean up
			continue                   // Loop will pause at the top due to isConnected check
		}

		responseLength := int(responseLengthBytes[0])<<8 | int(responseLengthBytes[1])

		// Read the response JSON data
		responseData := make([]byte, responseLength)
		// Use ReadFull to ensure all bytes are read
		bytesRead, err := io.ReadFull(conn, responseData)
		if err != nil {
			adapter.logger.WithFields(logrus.Fields{
				"op":        op,
				"err":       err.Error(),
				"readBytes": bytesRead, // Log how many bytes were read before error
				"expected":  responseLength,
			}).Error("Failed to read response data, handling disconnect")

			adapter.handleDisconnect()
			continue
		}

		// Launch a goroutine to process the response without blocking the read loop
		go adapter.processResponse(responseData)
	}
}

// processResponse handles decoding a single response and routing it to the correct waiter.
// This runs in its own goroutine, launched by readResponses.
func (adapter *Adapter) processResponse(responseData []byte) {
	const op = "py_core_adapter/processResponse"

	// Decode the response JSON data
	var responsePayload map[string]any
	if err := json.Unmarshal(responseData, &responsePayload); err != nil {
		adapter.logger.WithFields(logrus.Fields{
			"op":  op,
			"err": err.Error(),
			// Avoid logging raw data in production if it might contain sensitive info
			// "rawData": string(responseData),
		}).Error("Failed to unmarshal response")
		return // Cannot proceed without unmarshalling
	}

	// Find the message ID in the response
	messageID, ok := adapter.extractMessageID(responsePayload)
	if !ok {
		adapter.logger.WithFields(logrus.Fields{
			"op":       op,
			"response": fmt.Sprintf("%+v", responsePayload),
		}).Error("Response missing id_message")
		return // Cannot route without message ID
	}

	// Find the waiting channel for this message ID and remove it atomically
	adapter.requestsMutex.Lock() // Use full lock as we modify the map
	responseChan, exists := adapter.pendingRequests[messageID]
	if exists {
		// IMPORTANT: Remove the channel *before* sending to prevent potential race
		// where SendRequest times out *after* this goroutine checks existence
		// but *before* it sends, leading to a send on a potentially closed channel
		// if the timeout handler also tried to clean up.
		delete(adapter.pendingRequests, messageID)
	}
	adapter.requestsMutex.Unlock() // Unlock promptly after map access

	if !exists {
		adapter.logger.WithFields(logrus.Fields{
			"op":         op,
			"id_message": messageID,
			// "response":   fmt.Sprintf("%+v", responsePayload), // Maybe too verbose for just a missing waiter
		}).Warn("No waiting request found for id_message (potentially timed out or already processed)")
		return
	}

	// Send response to waiting handler
	// Use a select with a short timeout to avoid blocking indefinitely if the receiver
	// (SendRequest) somehow vanished or isn't ready.
	select {
	case responseChan <- responsePayload:
		adapter.logger.WithFields(logrus.Fields{
			"op":         op,
			"id_message": messageID,
		}).Debug("Response delivered to waiting request") // Changed level to Debug for less noise
	case <-time.After(1 * time.Second): // Short timeout for safety
		adapter.logger.WithFields(logrus.Fields{
			"op":         op,
			"id_message": messageID,
		}).Error("Failed to send response to channel within timeout (receiver likely gone)")
	}
	// Note: We don't close the channel here; SendRequest is responsible for its lifecycle.
}

// extractMessageID gets the id_message from a response payload
func (adapter *Adapter) extractMessageID(payload map[string]any) (string, bool) {
	if msgID, ok := payload["id_message"]; ok {
		if id, ok := msgID.(string); ok {
			return id, true
		}
	}
	return "", false
}
