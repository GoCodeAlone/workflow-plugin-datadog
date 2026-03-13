package internal

import (
	"context"
	"sync"
)

// datadogContext holds the authenticated context for Datadog API calls.
type datadogContext struct {
	ctx  context.Context
	site string
}

var (
	clientMu       sync.RWMutex
	clientRegistry = make(map[string]*datadogContext)
)

// RegisterClient adds a Datadog context to the global registry under the given name.
func RegisterClient(name string, c *datadogContext) {
	clientMu.Lock()
	defer clientMu.Unlock()
	clientRegistry[name] = c
}

// GetClient looks up a Datadog context by name.
func GetClient(name string) (*datadogContext, bool) {
	clientMu.RLock()
	defer clientMu.RUnlock()
	c, ok := clientRegistry[name]
	return c, ok
}

// UnregisterClient removes a client from the registry.
func UnregisterClient(name string) {
	clientMu.Lock()
	defer clientMu.Unlock()
	delete(clientRegistry, name)
}
