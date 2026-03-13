package internal

import (
	"context"

	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// incidentCreateStep implements step.datadog_incident_create
type incidentCreateStep struct {
	name       string
	moduleName string
}

func newIncidentCreateStep(name string, config map[string]any) (*incidentCreateStep, error) {
	return &incidentCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *incidentCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	title := resolveValue("title", current, config)
	if title == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "title is required"}}, nil
	}
	customerImpacted := resolveBool("customer_impacted", current, config)
	body := datadogV2.IncidentCreateRequest{
		Data: datadogV2.IncidentCreateData{
			Type: datadogV2.INCIDENTTYPE_INCIDENTS,
			Attributes: datadogV2.IncidentCreateAttributes{
				Title:            title,
				CustomerImpacted: customerImpacted,
			},
		},
	}
	api := datadogV2.NewIncidentsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateIncident(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":    resp.Data.Id,
		"title": title,
	}}, nil
}

// incidentGetStep implements step.datadog_incident_get
type incidentGetStep struct {
	name       string
	moduleName string
}

func newIncidentGetStep(name string, config map[string]any) (*incidentGetStep, error) {
	return &incidentGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *incidentGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	incidentID := resolveValue("incident_id", current, config)
	if incidentID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "incident_id is required"}}, nil
	}
	api := datadogV2.NewIncidentsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetIncident(ddCtx.ctx, incidentID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":    resp.Data.Id,
		"title": resp.Data.Attributes.Title,
	}}, nil
}

// incidentUpdateStep implements step.datadog_incident_update
type incidentUpdateStep struct {
	name       string
	moduleName string
}

func newIncidentUpdateStep(name string, config map[string]any) (*incidentUpdateStep, error) {
	return &incidentUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *incidentUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	incidentID := resolveValue("incident_id", current, config)
	if incidentID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "incident_id is required"}}, nil
	}
	body := datadogV2.IncidentUpdateRequest{
		Data: datadogV2.IncidentUpdateData{
			Id:   incidentID,
			Type: datadogV2.INCIDENTTYPE_INCIDENTS,
			Attributes: &datadogV2.IncidentUpdateAttributes{},
		},
	}
	if title := resolveValue("title", current, config); title != "" {
		body.Data.Attributes.SetTitle(title)
	}
	api := datadogV2.NewIncidentsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.UpdateIncident(ddCtx.ctx, incidentID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":      resp.Data.Id,
		"updated": true,
	}}, nil
}

// incidentDeleteStep implements step.datadog_incident_delete
type incidentDeleteStep struct {
	name       string
	moduleName string
}

func newIncidentDeleteStep(name string, config map[string]any) (*incidentDeleteStep, error) {
	return &incidentDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *incidentDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	incidentID := resolveValue("incident_id", current, config)
	if incidentID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "incident_id is required"}}, nil
	}
	api := datadogV2.NewIncidentsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, err := api.DeleteIncident(ddCtx.ctx, incidentID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": incidentID}}, nil
}

// incidentListStep implements step.datadog_incident_list
type incidentListStep struct {
	name       string
	moduleName string
}

func newIncidentListStep(name string, config map[string]any) (*incidentListStep, error) {
	return &incidentListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *incidentListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewIncidentsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListIncidents(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	incidents := make([]any, 0)
	for _, inc := range resp.Data {
			incidents = append(incidents, map[string]any{
				"id":    inc.Id,
				"title": inc.Attributes.Title,
			})
		}
	return &sdk.StepResult{Output: map[string]any{"incidents": incidents, "count": len(incidents)}}, nil
}

// incidentTodoCreateStep implements step.datadog_incident_todo_create
type incidentTodoCreateStep struct {
	name       string
	moduleName string
}

func newIncidentTodoCreateStep(name string, config map[string]any) (*incidentTodoCreateStep, error) {
	return &incidentTodoCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *incidentTodoCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	incidentID := resolveValue("incident_id", current, config)
	if incidentID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "incident_id is required"}}, nil
	}
	content := resolveValue("content", current, config)
	if content == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "content is required"}}, nil
	}
	body := datadogV2.IncidentTodoCreateRequest{
		Data: datadogV2.IncidentTodoCreateData{
			Type: datadogV2.INCIDENTTODOTYPE_INCIDENT_TODOS,
			Attributes: datadogV2.IncidentTodoAttributes{
				Content: content,
			},
		},
	}
	api := datadogV2.NewIncidentsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateIncidentTodo(ddCtx.ctx, incidentID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":          resp.Data.Id,
		"incident_id": incidentID,
	}}, nil
}

// incidentTodoUpdateStep implements step.datadog_incident_todo_update
type incidentTodoUpdateStep struct {
	name       string
	moduleName string
}

func newIncidentTodoUpdateStep(name string, config map[string]any) (*incidentTodoUpdateStep, error) {
	return &incidentTodoUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *incidentTodoUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	incidentID := resolveValue("incident_id", current, config)
	todoID := resolveValue("todo_id", current, config)
	if incidentID == "" || todoID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "incident_id and todo_id are required"}}, nil
	}
	content := resolveValue("content", current, config)
	body := datadogV2.IncidentTodoPatchRequest{
		Data: datadogV2.IncidentTodoPatchData{
			Type: datadogV2.INCIDENTTODOTYPE_INCIDENT_TODOS,
			Attributes: datadogV2.IncidentTodoAttributes{
				Content: content,
			},
		},
	}
	completed := resolveBool("completed", current, config)
	if completed {
		body.Data.Attributes.SetCompleted("true")
	}
	api := datadogV2.NewIncidentsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.UpdateIncidentTodo(ddCtx.ctx, incidentID, todoID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":      resp.Data.Id,
		"updated": true,
	}}, nil
}

// incidentTodoDeleteStep implements step.datadog_incident_todo_delete
type incidentTodoDeleteStep struct {
	name       string
	moduleName string
}

func newIncidentTodoDeleteStep(name string, config map[string]any) (*incidentTodoDeleteStep, error) {
	return &incidentTodoDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *incidentTodoDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	incidentID := resolveValue("incident_id", current, config)
	todoID := resolveValue("todo_id", current, config)
	if incidentID == "" || todoID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "incident_id and todo_id are required"}}, nil
	}
	api := datadogV2.NewIncidentsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, err := api.DeleteIncidentTodo(ddCtx.ctx, incidentID, todoID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": todoID}}, nil
}
