package internal_test

import (
	"testing"

	"github.com/GoCodeAlone/workflow/wftest"
)

func TestIntegration_EventCreate(t *testing.T) {
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  send-event:
    steps:
      - name: create_event
        type: step.datadog_event_create
        config:
          title: "Deploy complete"
          text: "v1.2.3 deployed successfully"
          tags:
            - "env:production"
      - name: confirm
        type: step.set
        config:
          values:
            done: true
`),
		wftest.MockStep("step.datadog_event_create", wftest.Returns(map[string]any{
			"id":     int64(12345),
			"title":  "Deploy complete",
			"status": "ok",
		})),
	)

	result := h.ExecutePipeline("send-event", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["done"] != true {
		t.Error("expected done=true")
	}
	eventOut := result.StepOutput("create_event")
	if eventOut["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", eventOut["status"])
	}
}

func TestIntegration_MonitorCreate(t *testing.T) {
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  create-monitor:
    steps:
      - name: monitor
        type: step.datadog_monitor_create
        config:
          name: "High CPU Alert"
          query: "avg(last_5m):avg:system.cpu.user{*} > 90"
          message: "CPU usage is too high!"
      - name: result
        type: step.set
        config:
          values:
            monitor_created: true
`),
		wftest.MockStep("step.datadog_monitor_create", wftest.Returns(map[string]any{
			"id":   int64(67890),
			"name": "High CPU Alert",
		})),
	)

	result := h.ExecutePipeline("create-monitor", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["monitor_created"] != true {
		t.Error("expected monitor_created=true")
	}
	monitorOut := result.StepOutput("monitor")
	if monitorOut["name"] != "High CPU Alert" {
		t.Errorf("expected name='High CPU Alert', got %v", monitorOut["name"])
	}
}

func TestIntegration_MetricSubmit(t *testing.T) {
	recorder := wftest.RecordStep("step.datadog_metric_submit").WithOutput(map[string]any{
		"accepted": true,
	})

	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  submit-metric:
    steps:
      - name: metric
        type: step.datadog_metric_submit
        config:
          metric: "app.request.count"
          value: 42.0
          tags:
            - "env:staging"
      - name: status
        type: step.set
        config:
          values:
            submitted: true
`),
		recorder,
	)

	result := h.ExecutePipeline("submit-metric", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["submitted"] != true {
		t.Error("expected submitted=true")
	}
	calls := recorder.Calls()
	if len(calls) != 1 {
		t.Errorf("expected 1 call to metric submit, got %d", len(calls))
	}
}
