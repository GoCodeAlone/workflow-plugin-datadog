package internal

import (
	"context"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// dashboardCreateStep implements step.datadog_dashboard_create
type dashboardCreateStep struct {
	name       string
	moduleName string
}

func newDashboardCreateStep(name string, config map[string]any) (*dashboardCreateStep, error) {
	return &dashboardCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *dashboardCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	title := resolveValue("title", current, config)
	if title == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "title is required"}}, nil
	}
	layoutType := datadogV1.DASHBOARDLAYOUTTYPE_ORDERED
	if lt := resolveValue("layout_type", current, config); lt == "free" {
		layoutType = datadogV1.DASHBOARDLAYOUTTYPE_FREE
	}
	body := datadogV1.Dashboard{
		Title:      title,
		LayoutType: layoutType,
		Widgets:    []datadogV1.Widget{},
	}
	if desc := resolveValue("description", current, config); desc != "" {
		body.SetDescription(desc)
	}
	api := datadogV1.NewDashboardsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateDashboard(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":    derefStr(resp.Id),
		"title": resp.Title,
		"url":   derefStr(resp.Url),
	}}, nil
}

// dashboardGetStep implements step.datadog_dashboard_get
type dashboardGetStep struct {
	name       string
	moduleName string
}

func newDashboardGetStep(name string, config map[string]any) (*dashboardGetStep, error) {
	return &dashboardGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *dashboardGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	dashboardID := resolveValue("dashboard_id", current, config)
	if dashboardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "dashboard_id is required"}}, nil
	}
	api := datadogV1.NewDashboardsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetDashboard(ddCtx.ctx, dashboardID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":          derefStr(resp.Id),
		"title":       resp.Title,
		"url":         derefStr(resp.Url),
		"description": resp.GetDescription(),
	}}, nil
}

// dashboardUpdateStep implements step.datadog_dashboard_update
type dashboardUpdateStep struct {
	name       string
	moduleName string
}

func newDashboardUpdateStep(name string, config map[string]any) (*dashboardUpdateStep, error) {
	return &dashboardUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *dashboardUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	dashboardID := resolveValue("dashboard_id", current, config)
	if dashboardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "dashboard_id is required"}}, nil
	}
	title := resolveValue("title", current, config)
	if title == "" {
		title = "Updated Dashboard"
	}
	layoutType := datadogV1.DASHBOARDLAYOUTTYPE_ORDERED
	body := datadogV1.Dashboard{
		Title:      title,
		LayoutType: layoutType,
		Widgets:    []datadogV1.Widget{},
	}
	api := datadogV1.NewDashboardsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.UpdateDashboard(ddCtx.ctx, dashboardID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":    derefStr(resp.Id),
		"title": resp.Title,
	}}, nil
}

// dashboardDeleteStep implements step.datadog_dashboard_delete
type dashboardDeleteStep struct {
	name       string
	moduleName string
}

func newDashboardDeleteStep(name string, config map[string]any) (*dashboardDeleteStep, error) {
	return &dashboardDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *dashboardDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	dashboardID := resolveValue("dashboard_id", current, config)
	if dashboardID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "dashboard_id is required"}}, nil
	}
	api := datadogV1.NewDashboardsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.DeleteDashboard(ddCtx.ctx, dashboardID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": derefStr(resp.DeletedDashboardId)}}, nil
}

// dashboardListStep implements step.datadog_dashboard_list
type dashboardListStep struct {
	name       string
	moduleName string
}

func newDashboardListStep(name string, config map[string]any) (*dashboardListStep, error) {
	return &dashboardListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *dashboardListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV1.NewDashboardsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListDashboards(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	dashboards := make([]any, 0)
	for _, d := range resp.Dashboards {
		dashboards = append(dashboards, map[string]any{
			"id":    derefStr(d.Id),
			"title": derefStr(d.Title),
			"url":   derefStr(d.Url),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"dashboards": dashboards, "count": len(dashboards)}}, nil
}
