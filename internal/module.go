package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
)

// datadogModule creates a Datadog API context and registers it.
type datadogModule struct {
	name   string
	config map[string]any
}

func newDatadogModule(name string, config map[string]any) (*datadogModule, error) {
	return &datadogModule{name: name, config: config}, nil
}

// Init creates the Datadog API context and registers it in the global registry.
func (m *datadogModule) Init() error {
	apiKey, _ := m.config["apiKey"].(string)
	appKey, _ := m.config["appKey"].(string)
	if apiKey == "" {
		return fmt.Errorf("datadog.provider %q: apiKey is required", m.name)
	}
	if appKey == "" {
		return fmt.Errorf("datadog.provider %q: appKey is required", m.name)
	}

	site, _ := m.config["site"].(string)
	if site == "" {
		site = "datadoghq.com"
	}

	keys := map[string]datadog.APIKey{
		"apiKeyAuth": {Key: apiKey},
		"appKeyAuth": {Key: appKey},
	}

	ctx := datadog.NewDefaultContext(context.Background())
	ctx = context.WithValue(ctx, datadog.ContextAPIKeys, keys)
	ctx = context.WithValue(ctx, datadog.ContextServerVariables, map[string]string{
		"site": site,
	})

	// If custom apiUrl provided, install a URL-rewriting HTTP transport so all
	// Datadog SDK calls (including operation-specific endpoints like logs intake)
	// are redirected to the target URL. This is necessary because the SDK creates
	// fresh Configuration/APIClient instances in each step and some API operations
	// use hardcoded server templates that can't be overridden via context variables.
	if apiUrl, ok := m.config["apiUrl"].(string); ok && apiUrl != "" {
		target, err := url.Parse(strings.TrimRight(apiUrl, "/"))
		if err != nil {
			return fmt.Errorf("datadog.provider %q: invalid apiUrl: %w", m.name, err)
		}
		base := http.DefaultTransport
		if base == nil {
			base = &http.Transport{}
		}
		http.DefaultTransport = &urlRewriteTransport{target: target, base: base}
	}

	RegisterClient(m.name, &datadogContext{ctx: ctx, site: site})
	return nil
}

// urlRewriteTransport rewrites all outgoing HTTP requests to point to a target URL,
// preserving the original request path and query. Used for mock/test redirection.
type urlRewriteTransport struct {
	target *url.URL
	base   http.RoundTripper
}

func (t *urlRewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.URL.Scheme = t.target.Scheme
	req2.URL.Host = t.target.Host
	req2.Host = t.target.Host
	return t.base.RoundTrip(req2)
}

// Start is a no-op for this module.
func (m *datadogModule) Start(_ context.Context) error { return nil }

// Stop unregisters the Datadog client.
func (m *datadogModule) Stop(_ context.Context) error {
	UnregisterClient(m.name)
	return nil
}
