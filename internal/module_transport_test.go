package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestDefaultTransportMutation_Vulnerability proves the pre-fix bug:
// Init() with apiUrl set used to replace http.DefaultTransport globally.
// This test is kept as documentation — after the fix it should FAIL to detect
// the mutation (i.e. DefaultTransport must remain unchanged).
func TestDefaultTransportUnchangedAfterInit(t *testing.T) {
	original := http.DefaultTransport

	m, err := newDatadogModule("test-transport", map[string]any{
		"apiKey": "test-api-key",
		"appKey": "test-app-key",
		"apiUrl": "http://proxy.example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	defer UnregisterClient("test-transport")

	if http.DefaultTransport != original {
		t.Error("security: Init() must not mutate http.DefaultTransport (affects all HTTP clients in the process)")
	}
}

// TestModuleScopedHTTPClient verifies that when apiUrl is set, the module stores
// a custom HTTP client scoped to the datadogContext, and that http.DefaultClient
// is NOT affected by the URL rewrite.
func TestModuleScopedHTTPClient(t *testing.T) {
	// Track which server received each request.
	realHit := false
	proxyHit := false

	realServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		realHit = true
		w.WriteHeader(http.StatusOK)
	}))
	defer realServer.Close()

	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyHit = true
		w.WriteHeader(http.StatusOK)
	}))
	defer proxyServer.Close()

	m, err := newDatadogModule("test-scoped", map[string]any{
		"apiKey": "test-api-key",
		"appKey": "test-app-key",
		"apiUrl": proxyServer.URL,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	defer UnregisterClient("test-scoped")

	// The module's internal client should use the rewrite transport.
	ddCtx, ok := GetClient("test-scoped")
	if !ok || ddCtx == nil {
		t.Fatal("expected client to be registered")
	}
	if ddCtx.httpClient == nil {
		t.Fatal("expected module to store a scoped httpClient")
	}

	// A request via the module's httpClient goes to proxyServer (URL rewritten).
	resp, err := ddCtx.httpClient.Get(realServer.URL + "/test")
	if err != nil {
		t.Fatalf("module httpClient.Get: %v", err)
	}
	resp.Body.Close()
	if !proxyHit {
		t.Error("expected module httpClient to route through the proxy server")
	}

	// Reset and verify http.DefaultClient is unaffected.
	realHit = false
	proxyHit = false
	resp2, err := http.DefaultClient.Get(realServer.URL + "/test")
	if err != nil {
		t.Fatalf("http.DefaultClient.Get: %v", err)
	}
	resp2.Body.Close()
	if !realHit {
		t.Error("expected http.DefaultClient to reach the real server")
	}
	if proxyHit {
		t.Error("security: http.DefaultClient must not be routed through the proxy")
	}
}

// TestNewConfig_UsesModuleHTTPClient verifies that newConfig() returns a
// datadog.Configuration whose HTTPClient is the module-scoped client.
func TestNewConfig_UsesModuleHTTPClient(t *testing.T) {
	m, err := newDatadogModule("test-cfg", map[string]any{
		"apiKey": "test-api-key",
		"appKey": "test-app-key",
		"apiUrl": "http://proxy.example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	defer UnregisterClient("test-cfg")

	ddCtx, ok := GetClient("test-cfg")
	if !ok {
		t.Fatal("expected client")
	}

	cfg := ddCtx.newConfig()
	if cfg.HTTPClient != ddCtx.httpClient {
		t.Error("newConfig() must set HTTPClient to the module-scoped client")
	}
}

// TestNewConfig_NoApiUrl returns a config with no custom HTTPClient when
// apiUrl is not set (normal production use).
func TestNewConfig_NoApiUrl(t *testing.T) {
	m, err := newDatadogModule("test-nourl", map[string]any{
		"apiKey": "test-api-key",
		"appKey": "test-app-key",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	defer UnregisterClient("test-nourl")

	ddCtx, ok := GetClient("test-nourl")
	if !ok {
		t.Fatal("expected client")
	}
	if ddCtx.httpClient != nil {
		t.Error("httpClient should be nil when apiUrl is not set")
	}
	cfg := ddCtx.newConfig()
	if cfg.HTTPClient != nil {
		t.Error("newConfig() should not set HTTPClient when apiUrl is not set")
	}
}
