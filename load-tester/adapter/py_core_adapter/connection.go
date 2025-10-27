package py_core_adapter

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

// handleDisconnect cleans up connection state and pending requests on disconnect
func (adapter *Adapter) handleDisconnect() {
	adapter.connMutex.Lock()
	defer adapter.connMutex.Unlock()

	if adapter.conn != nil {
		adapter.conn.Close()
	}
	adapter.isConnected = false

	// Notify all waiting requests about the disconnect
	adapter.requestsMutex.Lock()
	for messageID, ch := range adapter.pendingRequests {
		select {
		case ch <- map[string]any{
			"err_info":   "connection lost",
			"status":     "996",
			"id_message": messageID,
		}:
		default:
		}
		delete(adapter.pendingRequests, messageID)
	}
	adapter.requestsMutex.Unlock()
}

// Connect establishes a connection to the server if not already connected
func (adapter *Adapter) Connect() error {
	const op = "py_core_adapter/Connect"

	adapter.connMutex.Lock()
	defer adapter.connMutex.Unlock()

	if adapter.isConnected {
		return nil
	}

	adapter.logger.WithFields(logrus.Fields{
		"op":      op,
		"address": adapter.address,
	}).Info("Establishing connection to server")

	// Create a TCP connection
	conn, err := net.Dial("tcp", adapter.address)
	if err != nil {
		adapter.logger.WithFields(logrus.Fields{
			"op":      op,
			"address": adapter.address,
			"err":     err.Error(),
		}).Error("Failed to connect to server")

		return fmt.Errorf("error connecting to server: %w", err)
	}

	adapter.conn = conn
	adapter.isConnected = true

	adapter.logger.WithFields(logrus.Fields{
		"op":      op,
		"address": adapter.address,
	}).Info("TCP connection established")

	return nil
}

// Close closes the connection to the server
func (adapter *Adapter) Close() {
	const op = "py_core_adapter/Close"

	adapter.connMutex.Lock()
	defer adapter.connMutex.Unlock()

	if !adapter.isConnected {
		return
	}

	adapter.logger.WithFields(logrus.Fields{
		"op":      op,
		"address": adapter.address,
	}).Info("Closing connection to server")

	if adapter.conn != nil {
		adapter.conn.Close()
	}

	adapter.isConnected = false
}
