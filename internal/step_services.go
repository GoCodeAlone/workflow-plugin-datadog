package internal

import (
	"context"

	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// serviceDefinitionUpsertStep implements step.datadog_service_definition_upsert
type serviceDefinitionUpsertStep struct {
	name       string
	moduleName string
}

func newServiceDefinitionUpsertStep(name string, config map[string]any) (*serviceDefinitionUpsertStep, error) {
	return &serviceDefinitionUpsertStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *serviceDefinitionUpsertStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	serviceName := resolveValue("service_name", current, config)
	if serviceName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "service_name is required"}}, nil
	}
	// Build a v2 service definition payload
	payload := map[string]any{
		"schema-version": "v2",
		"dd-service":     serviceName,
	}
	if team := resolveValue("team", current, config); team != "" {
		payload["team"] = team
	}
	if desc := resolveValue("description", current, config); desc != "" {
		payload["description"] = desc
	}
	body := datadogV2.ServiceDefinitionsCreateRequest{
		ServiceDefinitionV2Dot2: &datadogV2.ServiceDefinitionV2Dot2{
			SchemaVersion: datadogV2.SERVICEDEFINITIONV2DOT2VERSION_V2_2,
			DdService:     serviceName,
		},
	}
	if team := resolveValue("team", current, config); team != "" {
		body.ServiceDefinitionV2Dot2.SetTeam(team)
	}
	api := datadogV2.NewServiceDefinitionApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateOrUpdateServiceDefinitions(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	_ = resp
	return &sdk.StepResult{Output: map[string]any{"upserted": true, "service_name": serviceName}}, nil
}

// serviceDefinitionGetStep implements step.datadog_service_definition_get
type serviceDefinitionGetStep struct {
	name       string
	moduleName string
}

func newServiceDefinitionGetStep(name string, config map[string]any) (*serviceDefinitionGetStep, error) {
	return &serviceDefinitionGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *serviceDefinitionGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	serviceName := resolveValue("service_name", current, config)
	if serviceName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "service_name is required"}}, nil
	}
	api := datadogV2.NewServiceDefinitionApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetServiceDefinition(ddCtx.ctx, serviceName)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	_ = resp
	return &sdk.StepResult{Output: map[string]any{"service_name": serviceName, "found": true}}, nil
}

// serviceDefinitionDeleteStep implements step.datadog_service_definition_delete
type serviceDefinitionDeleteStep struct {
	name       string
	moduleName string
}

func newServiceDefinitionDeleteStep(name string, config map[string]any) (*serviceDefinitionDeleteStep, error) {
	return &serviceDefinitionDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *serviceDefinitionDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	serviceName := resolveValue("service_name", current, config)
	if serviceName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "service_name is required"}}, nil
	}
	api := datadogV2.NewServiceDefinitionApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, err := api.DeleteServiceDefinition(ddCtx.ctx, serviceName)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "service_name": serviceName}}, nil
}

// serviceDefinitionListStep implements step.datadog_service_definition_list
type serviceDefinitionListStep struct {
	name       string
	moduleName string
}

func newServiceDefinitionListStep(name string, config map[string]any) (*serviceDefinitionListStep, error) {
	return &serviceDefinitionListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *serviceDefinitionListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewServiceDefinitionApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListServiceDefinitions(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	services := make([]any, 0)
	for _, svc := range resp.Data {
			svcName := ""
			if svc.Attributes != nil && svc.Attributes.Schema != nil {
				if v2 := svc.Attributes.Schema.ServiceDefinitionV2Dot2; v2 != nil {
					svcName = v2.DdService
				}
			}
			services = append(services, map[string]any{
				"service_name": svcName,
			})
		}
	return &sdk.StepResult{Output: map[string]any{"services": services, "count": len(services)}}, nil
}
