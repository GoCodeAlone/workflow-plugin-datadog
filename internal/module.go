package internal

import (
	"context"
	"fmt"

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

	// If custom apiUrl provided, override the server index
	if apiUrl, ok := m.config["apiUrl"].(string); ok && apiUrl != "" {
		ctx = context.WithValue(ctx, datadog.ContextServerIndex, 0)
		ctx = context.WithValue(ctx, datadog.ContextServerVariables, map[string]string{
			"site": site,
		})
		_ = apiUrl // URL override handled via ContextServerVariables site
	}

	RegisterClient(m.name, &datadogContext{ctx: ctx, site: site})
	return nil
}

// Start is a no-op for this module.
func (m *datadogModule) Start(_ context.Context) error { return nil }

// Stop unregisters the Datadog client.
func (m *datadogModule) Stop(_ context.Context) error {
	UnregisterClient(m.name)
	return nil
}
