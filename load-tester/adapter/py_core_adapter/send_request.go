package py_core_adapter

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// SendRequest sends a JSON request to the TCP server and returns the response
func (adapter *Adapter) SendRequest(payload []byte) (map[string]any, error) {
	const op = "py_core_adapter/SendRequest"

	adapter.logger.WithFields(logrus.Fields{
		"op":      op,
		"address": adapter.address,
		"payload": string(payload),
	}).Info("Processing TCP request")

	// Extract the message_id from the payload to track the request
	var requestMap map[string]any
	if err := json.Unmarshal(payload, &requestMap); err != nil {
		return nil, fmt.Errorf("error parsing payload: %w", err)
	}

	messageID, ok := adapter.extractMessageID(requestMap)
	if !ok {
		return nil, fmt.Errorf("request payload missing message_id field")
	}

	// Create a channel to receive the response for this specific message
	responseChan := make(chan map[string]any, 1)

	// Register this request in our tracking map
	adapter.requestsMutex.Lock()
	adapter.pendingRequests[messageID] = responseChan
	adapter.requestsMutex.Unlock()

	// Ensure we have a connection, or create one
	if err := adapter.Connect(); err != nil {
		// Clean up the pending request
		adapter.requestsMutex.Lock()
		delete(adapter.pendingRequests, messageID)
		adapter.requestsMutex.Unlock()

		return nil, err
	}

	// Send the request to the server under mutex
	adapter.connMutex.Lock()
	conn := adapter.conn // Keep a reference to the current connection

	// Send the length of the JSON data (2 bytes)
	length := len(payload)
	if length > 65535 {
		adapter.connMutex.Unlock()

		// Clean up the pending request
		adapter.requestsMutex.Lock()
		delete(adapter.pendingRequests, messageID)
		adapter.requestsMutex.Unlock()

		adapter.logger.WithFields(logrus.Fields{
			"op":     op,
			"length": length,
		}).Error("Payload too large")

		return nil, fmt.Errorf("payload too large: %d bytes", length)
	}

	lengthBytes := []byte{byte(length >> 8), byte(length & 0xff)}
	_, err := conn.Write(lengthBytes)
	if err != nil {
		adapter.connMutex.Unlock()
		adapter.handleDisconnect() // Handle connection failure

		// Clean up the pending request
		adapter.requestsMutex.Lock()
		delete(adapter.pendingRequests, messageID)
		adapter.requestsMutex.Unlock()

		adapter.logger.WithFields(logrus.Fields{
			"op":  op,
			"err": err.Error(),
		}).Error("Failed to send payload length")

		return nil, fmt.Errorf("error sending length: %w", err)
	}

	// Send the JSON data
	_, err = conn.Write(payload)
	adapter.connMutex.Unlock() // Release mutex after sending

	if err != nil {
		adapter.handleDisconnect() // Handle connection failure

		// Clean up the pending request
		adapter.requestsMutex.Lock()
		delete(adapter.pendingRequests, messageID)
		adapter.requestsMutex.Unlock()

		adapter.logger.WithFields(logrus.Fields{
			"op":  op,
			"err": err.Error(),
		}).Error("Failed to send payload data")

		return nil, fmt.Errorf("error sending JSON data: %w", err)
	}

	adapter.logger.WithFields(logrus.Fields{
		"op":         op,
		"message_id": messageID,
		"length":     length,
	}).Info("Request sent successfully")

	// Wait for the response - no client-side timeout, rely on server for timeout indication
	response, ok := <-responseChan
	if !ok {
		// Channel was closed, likely due to disconnection
		return nil, fmt.Errorf("response channel closed, connection lost")
	}

	adapter.logger.WithFields(logrus.Fields{
		"op":         op,
		"id_message": messageID,
	}).Debug("Received response from server via channel")

	return response, nil
}
