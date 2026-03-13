package internal

import (
	"context"
	"time"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// sloCreateStep implements step.datadog_slo_create
type sloCreateStep struct {
	name       string
	moduleName string
}

func newSLOCreateStep(name string, config map[string]any) (*sloCreateStep, error) {
	return &sloCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sloCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	sloName := resolveValue("name", current, config)
	if sloName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	sloType := datadogV1.SLOTYPE_METRIC
	if t := resolveValue("type", current, config); t == "monitor" {
		sloType = datadogV1.SLOTYPE_MONITOR
	}
	thresholds := []datadogV1.SLOThreshold{
		{
			Target:    99.9,
			Timeframe: datadogV1.SLOTIMEFRAME_THIRTY_DAYS,
		},
	}
	body := datadogV1.ServiceLevelObjectiveRequest{
		Name:       sloName,
		Type:       sloType,
		Thresholds: thresholds,
	}
	if desc := resolveValue("description", current, config); desc != "" {
		body.SetDescription(desc)
	}
	api := datadogV1.NewServiceLevelObjectivesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateSLO(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	sloID := ""
	if len(resp.Data) > 0 {
		sloID = derefStr(resp.Data[0].Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": sloID, "name": sloName}}, nil
}

// sloGetStep implements step.datadog_slo_get
type sloGetStep struct {
	name       string
	moduleName string
}

func newSLOGetStep(name string, config map[string]any) (*sloGetStep, error) {
	return &sloGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sloGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	sloID := resolveValue("slo_id", current, config)
	if sloID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "slo_id is required"}}, nil
	}
	api := datadogV1.NewServiceLevelObjectivesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetSLO(ddCtx.ctx, sloID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	out := map[string]any{"found": resp.Data != nil}
	if resp.Data != nil {
		out["id"] = derefStr(resp.Data.Id)
		out["name"] = derefStr(resp.Data.Name)
	}
	return &sdk.StepResult{Output: out}, nil
}

// sloUpdateStep implements step.datadog_slo_update
type sloUpdateStep struct {
	name       string
	moduleName string
}

func newSLOUpdateStep(name string, config map[string]any) (*sloUpdateStep, error) {
	return &sloUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sloUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	sloID := resolveValue("slo_id", current, config)
	if sloID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "slo_id is required"}}, nil
	}
	sloName := resolveValue("name", current, config)
	body := datadogV1.ServiceLevelObjective{
		Id:   &sloID,
		Name: sloName,
		Type: datadogV1.SLOTYPE_METRIC,
		Thresholds: []datadogV1.SLOThreshold{
			{Target: 99.9, Timeframe: datadogV1.SLOTIMEFRAME_THIRTY_DAYS},
		},
	}
	api := datadogV1.NewServiceLevelObjectivesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.UpdateSLO(ddCtx.ctx, sloID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	updatedID := ""
	if len(resp.Data) > 0 {
		updatedID = derefStr(resp.Data[0].Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": updatedID, "updated": true}}, nil
}

// sloDeleteStep implements step.datadog_slo_delete
type sloDeleteStep struct {
	name       string
	moduleName string
}

func newSLODeleteStep(name string, config map[string]any) (*sloDeleteStep, error) {
	return &sloDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sloDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	sloID := resolveValue("slo_id", current, config)
	if sloID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "slo_id is required"}}, nil
	}
	api := datadogV1.NewServiceLevelObjectivesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, _, err := api.DeleteSLO(ddCtx.ctx, sloID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": sloID}}, nil
}

// sloListStep implements step.datadog_slo_list
type sloListStep struct {
	name       string
	moduleName string
}

func newSLOListStep(name string, config map[string]any) (*sloListStep, error) {
	return &sloListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sloListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	params := datadogV1.NewListSLOsOptionalParameters()
	if q := resolveValue("query", current, config); q != "" {
		params.WithQuery(q)
	}
	api := datadogV1.NewServiceLevelObjectivesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListSLOs(ddCtx.ctx, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	slos := make([]any, 0)
	for _, slo := range resp.Data {
			slos = append(slos, map[string]any{
				"id":   derefStr(slo.Id),
				"name": slo.Name,
			})
		}
	return &sdk.StepResult{Output: map[string]any{"slos": slos, "count": len(slos)}}, nil
}

// sloSearchStep implements step.datadog_slo_search (uses list with query filter)
type sloSearchStep struct {
	name       string
	moduleName string
}

func newSLOSearchStep(name string, config map[string]any) (*sloSearchStep, error) {
	return &sloSearchStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sloSearchStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	params := datadogV1.NewListSLOsOptionalParameters()
	if q := resolveValue("query", current, config); q != "" {
		params.WithQuery(q)
	}
	api := datadogV1.NewServiceLevelObjectivesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListSLOs(ddCtx.ctx, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	slos := make([]any, 0)
	for _, slo := range resp.Data {
			slos = append(slos, map[string]any{
				"id":   derefStr(slo.Id),
				"name": slo.Name,
			})
		}
	return &sdk.StepResult{Output: map[string]any{"slos": slos, "count": len(slos)}}, nil
}

// sloHistoryGetStep implements step.datadog_slo_history_get
type sloHistoryGetStep struct {
	name       string
	moduleName string
}

func newSLOHistoryGetStep(name string, config map[string]any) (*sloHistoryGetStep, error) {
	return &sloHistoryGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sloHistoryGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	sloID := resolveValue("slo_id", current, config)
	if sloID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "slo_id is required"}}, nil
	}
	fromTs := resolveInt64("from", current, config)
	toTs := resolveInt64("to", current, config)
	if fromTs == 0 {
		fromTs = time.Now().Add(-30 * 24 * time.Hour).Unix()
	}
	if toTs == 0 {
		toTs = time.Now().Unix()
	}
	api := datadogV1.NewServiceLevelObjectivesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetSLOHistory(ddCtx.ctx, sloID, fromTs, toTs)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	_ = resp
	return &sdk.StepResult{Output: map[string]any{"slo_id": sloID, "from": fromTs, "to": toTs}}, nil
}
