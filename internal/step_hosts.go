package internal

import (
	"context"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// hostListStep implements step.datadog_host_list
type hostListStep struct {
	name       string
	moduleName string
}

func newHostListStep(name string, config map[string]any) (*hostListStep, error) {
	return &hostListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *hostListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	params := datadogV1.NewListHostsOptionalParameters()
	if filter := resolveValue("filter", current, config); filter != "" {
		params.WithFilter(filter)
	}
	api := datadogV1.NewHostsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.ListHosts(ddCtx.ctx, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	hosts := make([]any, 0)
	for _, h := range resp.HostList {
			hosts = append(hosts, map[string]any{
				"name": derefStr(h.Name),
				"id":   derefInt64(h.Id),
			})
		}
	return &sdk.StepResult{Output: map[string]any{"hosts": hosts, "count": len(hosts)}}, nil
}

// hostMuteStep implements step.datadog_host_mute
type hostMuteStep struct {
	name       string
	moduleName string
}

func newHostMuteStep(name string, config map[string]any) (*hostMuteStep, error) {
	return &hostMuteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *hostMuteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	hostName := resolveValue("host_name", current, config)
	if hostName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "host_name is required"}}, nil
	}
	body := datadogV1.HostMuteSettings{}
	if msg := resolveValue("message", current, config); msg != "" {
		body.SetMessage(msg)
	}
	api := datadogV1.NewHostsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.MuteHost(ddCtx.ctx, hostName, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"hostname": derefStr(resp.Hostname),
		"action":   derefStr(resp.Action),
	}}, nil
}

// hostUnmuteStep implements step.datadog_host_unmute
type hostUnmuteStep struct {
	name       string
	moduleName string
}

func newHostUnmuteStep(name string, config map[string]any) (*hostUnmuteStep, error) {
	return &hostUnmuteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *hostUnmuteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	hostName := resolveValue("host_name", current, config)
	if hostName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "host_name is required"}}, nil
	}
	api := datadogV1.NewHostsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.UnmuteHost(ddCtx.ctx, hostName)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"hostname": derefStr(resp.Hostname),
		"action":   derefStr(resp.Action),
	}}, nil
}

// hostTotalsGetStep implements step.datadog_host_totals_get
type hostTotalsGetStep struct {
	name       string
	moduleName string
}

func newHostTotalsGetStep(name string, config map[string]any) (*hostTotalsGetStep, error) {
	return &hostTotalsGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *hostTotalsGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV1.NewHostsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.GetHostTotals(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"total_active": derefInt64(resp.TotalActive),
		"total_up":     derefInt64(resp.TotalUp),
	}}, nil
}
