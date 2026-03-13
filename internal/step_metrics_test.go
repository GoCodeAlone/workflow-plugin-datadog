package internal

import (
	"context"
	"testing"
)

func TestModuleInit_RegistersClient(t *testing.T) {
	m, err := newDatadogModule("test-init", map[string]any{
		"apiKey": "test-api-key",
		"appKey": "test-app-key",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	c, ok := GetClient("test-init")
	if !ok || c == nil {
		t.Error("expected client to be registered")
	}
	// cleanup
	UnregisterClient("test-init")
}

func TestModuleStop_UnregistersClient(t *testing.T) {
	m, _ := newDatadogModule("test-stop", map[string]any{
		"apiKey": "test-api-key",
		"appKey": "test-app-key",
	})
	_ = m.Init()
	_ = m.Stop(context.Background())
	_, ok := GetClient("test-stop")
	if ok {
		t.Error("expected client to be unregistered after stop")
	}
}

func TestModuleInit_MissingAPIKey(t *testing.T) {
	m, err := newDatadogModule("test-missing", map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err == nil {
		t.Error("expected error for missing apiKey")
		UnregisterClient("test-missing")
	}
}

func TestModuleInit_MissingAppKey(t *testing.T) {
	m, err := newDatadogModule("test-missing-app", map[string]any{
		"apiKey": "test-api-key",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err == nil {
		t.Error("expected error for missing appKey")
		UnregisterClient("test-missing-app")
	}
}

func TestModuleInit_WithCustomSite(t *testing.T) {
	m, err := newDatadogModule("test-site", map[string]any{
		"apiKey": "test-api-key",
		"appKey": "test-app-key",
		"site":   "datadoghq.eu",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	c, ok := GetClient("test-site")
	if !ok || c == nil {
		t.Error("expected client to be registered")
	}
	if c.site != "datadoghq.eu" {
		t.Errorf("expected site datadoghq.eu, got %s", c.site)
	}
	UnregisterClient("test-site")
}

func TestMetricSubmitStep_NoClient(t *testing.T) {
	step, err := newMetricSubmitStep("test", map[string]any{"module": "nonexistent"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"metric": "test.metric",
		"value":  1.0,
	}, nil, map[string]any{"module": "nonexistent"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing client")
	}
}

func TestMetricSubmitStep_MissingMetric(t *testing.T) {
	// Register a client first
	m, _ := newDatadogModule("test-submit", map[string]any{
		"apiKey": "test-api-key",
		"appKey": "test-app-key",
	})
	_ = m.Init()
	defer UnregisterClient("test-submit")

	step, err := newMetricSubmitStep("test", map[string]any{"module": "test-submit"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{"module": "test-submit"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing metric")
	}
}

func TestMonitorCreateStep_NoClient(t *testing.T) {
	step, err := newMonitorCreateStep("test", map[string]any{"module": "nonexistent"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"name":  "test monitor",
		"query": "avg(last_5m):avg:system.cpu.user{*} > 90",
	}, nil, map[string]any{"module": "nonexistent"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing client")
	}
}

func TestEventCreateStep_MissingTitle(t *testing.T) {
	m, _ := newDatadogModule("test-event", map[string]any{
		"apiKey": "test-api-key",
		"appKey": "test-app-key",
	})
	_ = m.Init()
	defer UnregisterClient("test-event")

	step, err := newEventCreateStep("test", map[string]any{"module": "test-event"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{"module": "test-event"})
	if err != nil {
		t.Fatal(err)
	}
	if errMsg, ok := result.Output["error"].(string); !ok || errMsg != "title is required" {
		t.Errorf("expected 'title is required' error, got: %v", result.Output["error"])
	}
}

func TestStepRegistry_AllTypesRegistered(t *testing.T) {
	types := allStepTypes()
	if len(types) < 120 {
		t.Errorf("expected at least 120 step types, got %d", len(types))
	}
}

func TestCreateStep_UnknownType(t *testing.T) {
	_, err := createStep("step.datadog_unknown_type", "test", map[string]any{})
	if err == nil {
		t.Error("expected error for unknown step type")
	}
}

func TestCreateStep_AllRegisteredTypes(t *testing.T) {
	for typeName := range stepRegistry {
		step, err := createStep(typeName, "test", map[string]any{})
		if err != nil {
			t.Errorf("failed to create step %q: %v", typeName, err)
		}
		if step == nil {
			t.Errorf("step %q returned nil", typeName)
		}
	}
}

func TestHelpers_ResolveValue(t *testing.T) {
	current := map[string]any{"key": "from_current"}
	config := map[string]any{"key": "from_config"}
	val := resolveValue("key", current, config)
	if val != "from_current" {
		t.Errorf("expected 'from_current', got %q", val)
	}

	// Empty current, falls back to config
	val = resolveValue("key", map[string]any{}, config)
	if val != "from_config" {
		t.Errorf("expected 'from_config', got %q", val)
	}
}

func TestHelpers_GetModuleName(t *testing.T) {
	// Default
	name := getModuleName(map[string]any{})
	if name != "datadog" {
		t.Errorf("expected 'datadog', got %q", name)
	}
	// Custom
	name = getModuleName(map[string]any{"module": "my-dd"})
	if name != "my-dd" {
		t.Errorf("expected 'my-dd', got %q", name)
	}
}
