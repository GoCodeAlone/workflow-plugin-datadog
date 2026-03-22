package internal

import (
	"context"
	"net/http"
	"sync"

	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
)

// datadogContext holds the authenticated context for Datadog API calls.
type datadogContext struct {
	ctx        context.Context
	site       string
	httpClient *http.Client // non-nil only when apiUrl is configured
}

// newConfig returns a datadog.Configuration with the module-scoped HTTP client
// injected (when apiUrl was set), so steps never need to touch http.DefaultTransport.
func (c *datadogContext) newConfig() *datadog.Configuration {
	cfg := datadog.NewConfiguration()
	if c.httpClient != nil {
		cfg.HTTPClient = c.httpClient
	}
	return cfg
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
