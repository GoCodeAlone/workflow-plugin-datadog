package internal

import (
	"context"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// monitorCreateStep implements step.datadog_monitor_create
type monitorCreateStep struct {
	name       string
	moduleName string
}

func newMonitorCreateStep(name string, config map[string]any) (*monitorCreateStep, error) {
	return &monitorCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *monitorCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	monitorName := resolveValue("name", current, config)
	if monitorName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	query := resolveValue("query", current, config)
	if query == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "query is required"}}, nil
	}
	monitorType := datadogV1.MONITORTYPE_METRIC_ALERT
	if t := resolveValue("type", current, config); t != "" {
		switch t {
		case "service check":
			monitorType = datadogV1.MONITORTYPE_SERVICE_CHECK
		case "event alert":
			monitorType = datadogV1.MONITORTYPE_EVENT_ALERT
		case "log alert":
			monitorType = datadogV1.MONITORTYPE_LOG_ALERT
		case "query alert":
			monitorType = datadogV1.MONITORTYPE_QUERY_ALERT
		}
	}
	body := datadogV1.Monitor{
		Name:  &monitorName,
		Type:  monitorType,
		Query: query,
	}
	if msg := resolveValue("message", current, config); msg != "" {
		body.SetMessage(msg)
	}
	if tags := resolveStringSlice("tags", current, config); len(tags) > 0 {
		body.SetTags(tags)
	}
	api := datadogV1.NewMonitorsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateMonitor(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":   derefInt64(resp.Id),
		"name": derefStr(resp.Name),
	}}, nil
}

// monitorGetStep implements step.datadog_monitor_get
type monitorGetStep struct {
	name       string
	moduleName string
}

func newMonitorGetStep(name string, config map[string]any) (*monitorGetStep, error) {
	return &monitorGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *monitorGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	monitorID := resolveInt64("monitor_id", current, config)
	if monitorID == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "monitor_id is required"}}, nil
	}
	api := datadogV1.NewMonitorsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetMonitor(ddCtx.ctx, monitorID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":      derefInt64(resp.Id),
		"name":    derefStr(resp.Name),
		"query":   resp.Query,
		"message": derefStr(resp.Message),
	}}, nil
}

// monitorUpdateStep implements step.datadog_monitor_update
type monitorUpdateStep struct {
	name       string
	moduleName string
}

func newMonitorUpdateStep(name string, config map[string]any) (*monitorUpdateStep, error) {
	return &monitorUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *monitorUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	monitorID := resolveInt64("monitor_id", current, config)
	if monitorID == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "monitor_id is required"}}, nil
	}
	body := datadogV1.MonitorUpdateRequest{}
	if n := resolveValue("name", current, config); n != "" {
		body.SetName(n)
	}
	if q := resolveValue("query", current, config); q != "" {
		body.SetQuery(q)
	}
	if msg := resolveValue("message", current, config); msg != "" {
		body.SetMessage(msg)
	}
	api := datadogV1.NewMonitorsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.UpdateMonitor(ddCtx.ctx, monitorID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":   derefInt64(resp.Id),
		"name": derefStr(resp.Name),
	}}, nil
}

// monitorDeleteStep implements step.datadog_monitor_delete
type monitorDeleteStep struct {
	name       string
	moduleName string
}

func newMonitorDeleteStep(name string, config map[string]any) (*monitorDeleteStep, error) {
	return &monitorDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *monitorDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	monitorID := resolveInt64("monitor_id", current, config)
	if monitorID == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "monitor_id is required"}}, nil
	}
	api := datadogV1.NewMonitorsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.DeleteMonitor(ddCtx.ctx, monitorID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": derefInt64(resp.DeletedMonitorId)}}, nil
}

// monitorListStep implements step.datadog_monitor_list
type monitorListStep struct {
	name       string
	moduleName string
}

func newMonitorListStep(name string, config map[string]any) (*monitorListStep, error) {
	return &monitorListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *monitorListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	params := datadogV1.NewListMonitorsOptionalParameters()
	if tags := resolveValue("tags", current, config); tags != "" {
		params.WithTags(tags)
	}
	if name := resolveValue("name", current, config); name != "" {
		params.WithName(name)
	}
	api := datadogV1.NewMonitorsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListMonitors(ddCtx.ctx, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	monitors := make([]any, 0, len(resp))
	for _, m := range resp {
		monitors = append(monitors, map[string]any{
			"id":   derefInt64(m.Id),
			"name": derefStr(m.Name),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"monitors": monitors, "count": len(monitors)}}, nil
}

// monitorSearchStep implements step.datadog_monitor_search
type monitorSearchStep struct {
	name       string
	moduleName string
}

func newMonitorSearchStep(name string, config map[string]any) (*monitorSearchStep, error) {
	return &monitorSearchStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *monitorSearchStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	params := datadogV1.NewSearchMonitorsOptionalParameters()
	if q := resolveValue("query", current, config); q != "" {
		params.WithQuery(q)
	}
	api := datadogV1.NewMonitorsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.SearchMonitors(ddCtx.ctx, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	monitors := make([]any, 0)
	for _, m := range resp.Monitors {
		monitors = append(monitors, map[string]any{
			"id":   derefInt64(m.Id),
			"name": derefStr(m.Name),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"monitors": monitors, "count": len(monitors)}}, nil
}

// monitorValidateStep implements step.datadog_monitor_validate
type monitorValidateStep struct {
	name       string
	moduleName string
}

func newMonitorValidateStep(name string, config map[string]any) (*monitorValidateStep, error) {
	return &monitorValidateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *monitorValidateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	if query == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "query is required"}}, nil
	}
	body := datadogV1.Monitor{
		Type:  datadogV1.MONITORTYPE_METRIC_ALERT,
		Query: query,
	}
	api := datadogV1.NewMonitorsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, _, err := api.ValidateMonitor(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"valid": false, "error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"valid": true}}, nil
}
