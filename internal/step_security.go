package internal

import (
	"context"

	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// securityRuleCreateStep implements step.datadog_security_rule_create
type securityRuleCreateStep struct {
	name       string
	moduleName string
}

func newSecurityRuleCreateStep(name string, config map[string]any) (*securityRuleCreateStep, error) {
	return &securityRuleCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *securityRuleCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	ruleName := resolveValue("name", current, config)
	if ruleName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	message := resolveValue("message", current, config)
	enabled := true
	severity := datadogV2.SECURITYMONITORINGRULESEVERITY_INFO
	cases := []datadogV2.SecurityMonitoringRuleCaseCreate{
		{
			Status:  severity,
			Name:    &ruleName,
		},
	}
	options := datadogV2.SecurityMonitoringRuleOptions{
		EvaluationWindow:  datadogV2.SECURITYMONITORINGRULEEVALUATIONWINDOW_ZERO_MINUTES.Ptr(),
		KeepAlive:         datadogV2.SECURITYMONITORINGRULEKEEPALIVE_ZERO_MINUTES.Ptr(),
		MaxSignalDuration: datadogV2.SECURITYMONITORINGRULEMAXSIGNALDURATION_ZERO_MINUTES.Ptr(),
	}
	body := datadogV2.SecurityMonitoringRuleCreatePayload{
		SecurityMonitoringStandardRuleCreatePayload: &datadogV2.SecurityMonitoringStandardRuleCreatePayload{
			Name:    ruleName,
			Message: message,
			IsEnabled: enabled,
			Cases:   cases,
			Options: options,
			Queries: []datadogV2.SecurityMonitoringStandardRuleQuery{},
		},
	}
	api := datadogV2.NewSecurityMonitoringApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateSecurityMonitoringRule(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	ruleID := ""
	if resp.SecurityMonitoringStandardRuleResponse != nil {
		ruleID = derefStr(resp.SecurityMonitoringStandardRuleResponse.Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": ruleID, "name": ruleName}}, nil
}

// securityRuleGetStep implements step.datadog_security_rule_get
type securityRuleGetStep struct {
	name       string
	moduleName string
}

func newSecurityRuleGetStep(name string, config map[string]any) (*securityRuleGetStep, error) {
	return &securityRuleGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *securityRuleGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	ruleID := resolveValue("rule_id", current, config)
	if ruleID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "rule_id is required"}}, nil
	}
	api := datadogV2.NewSecurityMonitoringApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetSecurityMonitoringRule(ddCtx.ctx, ruleID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	ruleName := ""
	if resp.SecurityMonitoringStandardRuleResponse != nil {
		ruleName = derefStr(resp.SecurityMonitoringStandardRuleResponse.Name)
	}
	return &sdk.StepResult{Output: map[string]any{"id": ruleID, "name": ruleName}}, nil
}

// securityRuleUpdateStep implements step.datadog_security_rule_update
type securityRuleUpdateStep struct {
	name       string
	moduleName string
}

func newSecurityRuleUpdateStep(name string, config map[string]any) (*securityRuleUpdateStep, error) {
	return &securityRuleUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *securityRuleUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	ruleID := resolveValue("rule_id", current, config)
	if ruleID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "rule_id is required"}}, nil
	}
	ruleName := resolveValue("name", current, config)
	body := datadogV2.SecurityMonitoringRuleUpdatePayload{
		Name: &ruleName,
	}
	if msg := resolveValue("message", current, config); msg != "" {
		body.SetMessage(msg)
	}
	api := datadogV2.NewSecurityMonitoringApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.UpdateSecurityMonitoringRule(ddCtx.ctx, ruleID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	updatedID := ""
	if resp.SecurityMonitoringStandardRuleResponse != nil {
		updatedID = derefStr(resp.SecurityMonitoringStandardRuleResponse.Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": updatedID, "updated": true}}, nil
}

// securityRuleDeleteStep implements step.datadog_security_rule_delete
type securityRuleDeleteStep struct {
	name       string
	moduleName string
}

func newSecurityRuleDeleteStep(name string, config map[string]any) (*securityRuleDeleteStep, error) {
	return &securityRuleDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *securityRuleDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	ruleID := resolveValue("rule_id", current, config)
	if ruleID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "rule_id is required"}}, nil
	}
	api := datadogV2.NewSecurityMonitoringApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, err := api.DeleteSecurityMonitoringRule(ddCtx.ctx, ruleID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": ruleID}}, nil
}

// securityRuleListStep implements step.datadog_security_rule_list
type securityRuleListStep struct {
	name       string
	moduleName string
}

func newSecurityRuleListStep(name string, config map[string]any) (*securityRuleListStep, error) {
	return &securityRuleListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *securityRuleListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewSecurityMonitoringApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListSecurityMonitoringRules(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	rules := make([]any, 0)
	for _, r := range resp.Data {
			ruleID := ""
			ruleName := ""
			if r.SecurityMonitoringStandardRuleResponse != nil {
				ruleID = derefStr(r.SecurityMonitoringStandardRuleResponse.Id)
				ruleName = derefStr(r.SecurityMonitoringStandardRuleResponse.Name)
			}
			rules = append(rules, map[string]any{"id": ruleID, "name": ruleName})
		}
	return &sdk.StepResult{Output: map[string]any{"rules": rules, "count": len(rules)}}, nil
}

// securitySignalListStep implements step.datadog_security_signal_list
type securitySignalListStep struct {
	name       string
	moduleName string
}

func newSecuritySignalListStep(name string, config map[string]any) (*securitySignalListStep, error) {
	return &securitySignalListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *securitySignalListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	body := datadogV2.SecurityMonitoringSignalListRequest{
		Filter: &datadogV2.SecurityMonitoringSignalListRequestFilter{
			Query: &query,
		},
	}
	api := datadogV2.NewSecurityMonitoringApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.SearchSecurityMonitoringSignals(ddCtx.ctx, *datadogV2.NewSearchSecurityMonitoringSignalsOptionalParameters().WithBody(body))
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	signals := make([]any, 0)
	for _, sig := range resp.Data {
			signals = append(signals, map[string]any{
				"id": sig.Id,
			})
		}
	return &sdk.StepResult{Output: map[string]any{"signals": signals, "count": len(signals)}}, nil
}

// securitySignalStateUpdateStep implements step.datadog_security_signal_state_update
type securitySignalStateUpdateStep struct {
	name       string
	moduleName string
}

func newSecuritySignalStateUpdateStep(name string, config map[string]any) (*securitySignalStateUpdateStep, error) {
	return &securitySignalStateUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *securitySignalStateUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	signalID := resolveValue("signal_id", current, config)
	if signalID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "signal_id is required"}}, nil
	}
	state := datadogV2.SECURITYMONITORINGSIGNALSTATE_OPEN
	if s2 := resolveValue("state", current, config); s2 == "archived" {
		state = datadogV2.SECURITYMONITORINGSIGNALSTATE_ARCHIVED
	} else if s2 == "under_review" {
		state = datadogV2.SECURITYMONITORINGSIGNALSTATE_UNDER_REVIEW
	}
	body := datadogV2.SecurityMonitoringSignalStateUpdateRequest{
		Data: datadogV2.SecurityMonitoringSignalStateUpdateData{
			Attributes: datadogV2.SecurityMonitoringSignalStateUpdateAttributes{
				State: state,
			},
		},
	}
	api := datadogV2.NewSecurityMonitoringApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, _, err := api.EditSecurityMonitoringSignalState(ddCtx.ctx, signalID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"updated": true, "signal_id": signalID}}, nil
}
