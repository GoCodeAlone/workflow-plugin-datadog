package internal

import (
	"context"
	"time"

	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// auditLogSearchStep implements step.datadog_audit_log_search
type auditLogSearchStep struct {
	name       string
	moduleName string
}

func newAuditLogSearchStep(name string, config map[string]any) (*auditLogSearchStep, error) {
	return &auditLogSearchStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *auditLogSearchStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	from := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	to := time.Now().Format(time.RFC3339)
	body := datadogV2.AuditLogsSearchEventsRequest{
		Filter: &datadogV2.AuditLogsQueryFilter{
			Query: &query,
			From:  &from,
			To:    &to,
		},
	}
	params := datadogV2.NewSearchAuditLogsOptionalParameters().WithBody(body)
	api := datadogV2.NewAuditApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.SearchAuditLogs(ddCtx.ctx, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	events := make([]any, 0)
	for _, e := range resp.Data {
		events = append(events, map[string]any{
			"id": derefStr(e.Id),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"events": events, "count": len(events)}}, nil
}

// auditLogListStep implements step.datadog_audit_log_list
type auditLogListStep struct {
	name       string
	moduleName string
}

func newAuditLogListStep(name string, config map[string]any) (*auditLogListStep, error) {
	return &auditLogListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *auditLogListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	params := datadogV2.NewListAuditLogsOptionalParameters()
	if q := resolveValue("query", current, config); q != "" {
		params.WithFilterQuery(q)
	}
	api := datadogV2.NewAuditApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListAuditLogs(ddCtx.ctx, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	events := make([]any, 0)
	for _, e := range resp.Data {
		events = append(events, map[string]any{
			"id": derefStr(e.Id),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"events": events, "count": len(events)}}, nil
}
