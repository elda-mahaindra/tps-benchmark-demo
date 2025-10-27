package py_switching_adapter

import (
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type Adapter struct {
	logger  *logrus.Logger
	address string

	// Connection management
	conn        net.Conn
	connMutex   sync.Mutex
	isConnected bool

	// Response tracking
	pendingRequests map[string]chan map[string]any
	requestsMutex   sync.RWMutex
}

// NewAdapter creates a new TCP adapter instance
func NewAdapter(
	logger *logrus.Logger,
	address string,
) *Adapter {
	adapter := &Adapter{
		logger:          logger,
		address:         address,
		pendingRequests: make(map[string]chan map[string]any),
	}

	// Start response reader goroutine when creating the adapter
	go adapter.readResponses()

	return adapter
}
