package internal

import (
	"context"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// tagsGetStep implements step.datadog_tags_get
type tagsGetStep struct {
	name       string
	moduleName string
}

func newTagsGetStep(name string, config map[string]any) (*tagsGetStep, error) {
	return &tagsGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *tagsGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	hostName := resolveValue("host_name", current, config)
	if hostName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "host_name is required"}}, nil
	}
	api := datadogV1.NewTagsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.GetHostTags(ddCtx.ctx, hostName)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"host": derefStr(resp.Host),
		"tags": resp.Tags,
	}}, nil
}

// tagsUpdateStep implements step.datadog_tags_update
type tagsUpdateStep struct {
	name       string
	moduleName string
}

func newTagsUpdateStep(name string, config map[string]any) (*tagsUpdateStep, error) {
	return &tagsUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *tagsUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	hostName := resolveValue("host_name", current, config)
	if hostName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "host_name is required"}}, nil
	}
	tags := resolveStringSlice("tags", current, config)
	body := datadogV1.HostTags{
		Host: &hostName,
		Tags: tags,
	}
	api := datadogV1.NewTagsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.UpdateHostTags(ddCtx.ctx, hostName, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"host": derefStr(resp.Host),
		"tags": resp.Tags,
	}}, nil
}

// tagsDeleteStep implements step.datadog_tags_delete
type tagsDeleteStep struct {
	name       string
	moduleName string
}

func newTagsDeleteStep(name string, config map[string]any) (*tagsDeleteStep, error) {
	return &tagsDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *tagsDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	hostName := resolveValue("host_name", current, config)
	if hostName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "host_name is required"}}, nil
	}
	api := datadogV1.NewTagsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	_, err := api.DeleteHostTags(ddCtx.ctx, hostName)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "host": hostName}}, nil
}

// tagsListStep implements step.datadog_tags_list
type tagsListStep struct {
	name       string
	moduleName string
}

func newTagsListStep(name string, config map[string]any) (*tagsListStep, error) {
	return &tagsListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *tagsListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV1.NewTagsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.ListHostTags(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"tags": resp.Tags}}, nil
}
