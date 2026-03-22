package internal

import (
	"context"
	"time"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// downtimeCreateStep implements step.datadog_downtime_create
type downtimeCreateStep struct {
	name       string
	moduleName string
}

func newDowntimeCreateStep(name string, config map[string]any) (*downtimeCreateStep, error) {
	return &downtimeCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *downtimeCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	scope := resolveValue("scope", current, config)
	if scope == "" {
		scope = "*"
	}
	start := resolveInt64("start", current, config)
	if start == 0 {
		start = time.Now().Unix()
	}
	end := resolveInt64("end", current, config)
	if end == 0 {
		end = start + 3600
	}
	body := datadogV1.Downtime{
		Scope: []string{scope},
		Start: &start,
	}
	body.End = *datadog.NewNullableInt64(&end)
	if msg := resolveValue("message", current, config); msg != "" {
		body.SetMessage(msg)
	}
	api := datadogV1.NewDowntimesApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.CreateDowntime(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":      derefInt64(resp.Id),
		"scope":   resp.Scope,
		"message": resp.GetMessage(),
	}}, nil
}

// downtimeGetStep implements step.datadog_downtime_get
type downtimeGetStep struct {
	name       string
	moduleName string
}

func newDowntimeGetStep(name string, config map[string]any) (*downtimeGetStep, error) {
	return &downtimeGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *downtimeGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	downtimeID := resolveInt64("downtime_id", current, config)
	if downtimeID == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "downtime_id is required"}}, nil
	}
	api := datadogV1.NewDowntimesApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.GetDowntime(ddCtx.ctx, downtimeID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":      derefInt64(resp.Id),
		"scope":   resp.Scope,
		"message": resp.GetMessage(),
	}}, nil
}

// downtimeUpdateStep implements step.datadog_downtime_update
type downtimeUpdateStep struct {
	name       string
	moduleName string
}

func newDowntimeUpdateStep(name string, config map[string]any) (*downtimeUpdateStep, error) {
	return &downtimeUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *downtimeUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	downtimeID := resolveInt64("downtime_id", current, config)
	if downtimeID == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "downtime_id is required"}}, nil
	}
	body := datadogV1.Downtime{}
	if msg := resolveValue("message", current, config); msg != "" {
		body.SetMessage(msg)
	}
	if scope := resolveValue("scope", current, config); scope != "" {
		body.SetScope([]string{scope})
	}
	api := datadogV1.NewDowntimesApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.UpdateDowntime(ddCtx.ctx, downtimeID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":      derefInt64(resp.Id),
		"updated": true,
	}}, nil
}

// downtimeCancelStep implements step.datadog_downtime_cancel
type downtimeCancelStep struct {
	name       string
	moduleName string
}

func newDowntimeCancelStep(name string, config map[string]any) (*downtimeCancelStep, error) {
	return &downtimeCancelStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *downtimeCancelStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	downtimeID := resolveInt64("downtime_id", current, config)
	if downtimeID == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "downtime_id is required"}}, nil
	}
	api := datadogV1.NewDowntimesApi(datadog.NewAPIClient(ddCtx.newConfig()))
	_, err := api.CancelDowntime(ddCtx.ctx, downtimeID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"cancelled": true, "id": downtimeID}}, nil
}

// downtimeListStep implements step.datadog_downtime_list
type downtimeListStep struct {
	name       string
	moduleName string
}

func newDowntimeListStep(name string, config map[string]any) (*downtimeListStep, error) {
	return &downtimeListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *downtimeListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	params := datadogV1.NewListDowntimesOptionalParameters()
	if current_only := resolveBool("current_only", current, config); current_only {
		params.WithCurrentOnly(true)
	}
	api := datadogV1.NewDowntimesApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.ListDowntimes(ddCtx.ctx, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	downtimes := make([]any, 0, len(resp))
	for _, d := range resp {
		downtimes = append(downtimes, map[string]any{
			"id":      derefInt64(d.Id),
			"scope":   d.Scope,
			"message": d.GetMessage(),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"downtimes": downtimes, "count": len(downtimes)}}, nil
}
