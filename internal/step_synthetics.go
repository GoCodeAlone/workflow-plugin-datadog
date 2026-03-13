package internal

import (
	"context"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// syntheticsTestCreateStep implements step.datadog_synthetics_test_create
type syntheticsTestCreateStep struct {
	name       string
	moduleName string
}

func newSyntheticsTestCreateStep(name string, config map[string]any) (*syntheticsTestCreateStep, error) {
	return &syntheticsTestCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *syntheticsTestCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	testName := resolveValue("name", current, config)
	if testName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	url := resolveValue("url", current, config)
	if url == "" {
		url = "https://example.com"
	}
	locations := resolveStringSlice("locations", current, config)
	if len(locations) == 0 {
		locations = []string{"aws:us-east-1"}
	}
	body := datadogV1.SyntheticsAPITest{
		Name:      testName,
		Type:      datadogV1.SYNTHETICSAPITESTTYPE_API,
		Locations: locations,
		Config: datadogV1.SyntheticsAPITestConfig{
			Request: &datadogV1.SyntheticsTestRequest{
				Url:    &url,
				Method: datadog.PtrString("GET"),
			},
		},
		Options: datadogV1.SyntheticsTestOptions{},
	}
	if tags := resolveStringSlice("tags", current, config); len(tags) > 0 {
		body.SetTags(tags)
	}
	api := datadogV1.NewSyntheticsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateSyntheticsAPITest(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"public_id": derefStr(resp.PublicId),
		"name":      resp.Name,
	}}, nil
}

// syntheticsTestGetStep implements step.datadog_synthetics_test_get
type syntheticsTestGetStep struct {
	name       string
	moduleName string
}

func newSyntheticsTestGetStep(name string, config map[string]any) (*syntheticsTestGetStep, error) {
	return &syntheticsTestGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *syntheticsTestGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	publicID := resolveValue("public_id", current, config)
	if publicID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "public_id is required"}}, nil
	}
	api := datadogV1.NewSyntheticsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetTest(ddCtx.ctx, publicID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"public_id": derefStr(resp.PublicId),
		"name":      derefStr(resp.Name),
	}}, nil
}

// syntheticsTestUpdateStep implements step.datadog_synthetics_test_update
type syntheticsTestUpdateStep struct {
	name       string
	moduleName string
}

func newSyntheticsTestUpdateStep(name string, config map[string]any) (*syntheticsTestUpdateStep, error) {
	return &syntheticsTestUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *syntheticsTestUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	publicID := resolveValue("public_id", current, config)
	if publicID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "public_id is required"}}, nil
	}
	testName := resolveValue("name", current, config)
	url := resolveValue("url", current, config)
	if url == "" {
		url = "https://example.com"
	}
	locations := resolveStringSlice("locations", current, config)
	if len(locations) == 0 {
		locations = []string{"aws:us-east-1"}
	}
	body := datadogV1.SyntheticsAPITest{
		Name:      testName,
		Type:      datadogV1.SYNTHETICSAPITESTTYPE_API,
		Locations: locations,
		Config: datadogV1.SyntheticsAPITestConfig{
			Request: &datadogV1.SyntheticsTestRequest{
				Url:    &url,
				Method: datadog.PtrString("GET"),
			},
		},
		Options: datadogV1.SyntheticsTestOptions{},
	}
	api := datadogV1.NewSyntheticsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.UpdateAPITest(ddCtx.ctx, publicID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"public_id": derefStr(resp.PublicId),
		"name":      resp.Name,
	}}, nil
}

// syntheticsTestDeleteStep implements step.datadog_synthetics_test_delete
type syntheticsTestDeleteStep struct {
	name       string
	moduleName string
}

func newSyntheticsTestDeleteStep(name string, config map[string]any) (*syntheticsTestDeleteStep, error) {
	return &syntheticsTestDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *syntheticsTestDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	publicIDs := resolveStringSlice("public_ids", current, config)
	if len(publicIDs) == 0 {
		if id := resolveValue("public_id", current, config); id != "" {
			publicIDs = []string{id}
		}
	}
	if len(publicIDs) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "public_id or public_ids is required"}}, nil
	}
	body := datadogV1.SyntheticsDeleteTestsPayload{PublicIds: publicIDs}
	api := datadogV1.NewSyntheticsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.DeleteTests(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	deleted := make([]any, 0)
	for _, d := range resp.DeletedTests {
		deleted = append(deleted, derefStr(d.PublicId))
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": deleted, "count": len(deleted)}}, nil
}

// syntheticsTestListStep implements step.datadog_synthetics_test_list
type syntheticsTestListStep struct {
	name       string
	moduleName string
}

func newSyntheticsTestListStep(name string, config map[string]any) (*syntheticsTestListStep, error) {
	return &syntheticsTestListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *syntheticsTestListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV1.NewSyntheticsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListTests(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	tests := make([]any, 0)
	for _, t := range resp.Tests {
		tests = append(tests, map[string]any{
			"public_id": derefStr(t.PublicId),
			"name":      derefStr(t.Name),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"tests": tests, "count": len(tests)}}, nil
}

// syntheticsTestTriggerStep implements step.datadog_synthetics_test_trigger
type syntheticsTestTriggerStep struct {
	name       string
	moduleName string
}

func newSyntheticsTestTriggerStep(name string, config map[string]any) (*syntheticsTestTriggerStep, error) {
	return &syntheticsTestTriggerStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *syntheticsTestTriggerStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	publicIDs := resolveStringSlice("public_ids", current, config)
	if len(publicIDs) == 0 {
		if id := resolveValue("public_id", current, config); id != "" {
			publicIDs = []string{id}
		}
	}
	if len(publicIDs) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "public_id or public_ids is required"}}, nil
	}
	tests := make([]datadogV1.SyntheticsTriggerTest, 0, len(publicIDs))
	for _, id := range publicIDs {
		tests = append(tests, datadogV1.SyntheticsTriggerTest{PublicId: id})
	}
	body := datadogV1.SyntheticsTriggerBody{Tests: tests}
	api := datadogV1.NewSyntheticsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.TriggerTests(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	triggered := make([]any, 0)
	for _, r := range resp.Results {
		triggered = append(triggered, map[string]any{
			"public_id": derefStr(r.PublicId),
			"result_id": derefStr(r.ResultId),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"results": triggered, "count": len(triggered)}}, nil
}

// syntheticsResultsGetStep implements step.datadog_synthetics_results_get
type syntheticsResultsGetStep struct {
	name       string
	moduleName string
}

func newSyntheticsResultsGetStep(name string, config map[string]any) (*syntheticsResultsGetStep, error) {
	return &syntheticsResultsGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *syntheticsResultsGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	publicID := resolveValue("public_id", current, config)
	resultID := resolveValue("result_id", current, config)
	if publicID == "" || resultID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "public_id and result_id are required"}}, nil
	}
	api := datadogV1.NewSyntheticsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetAPITestResult(ddCtx.ctx, publicID, resultID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"result_id": derefStr(resp.ResultId),
		"status":    int64(resp.GetStatus()),
	}}, nil
}

// syntheticsGlobalVarCreateStep implements step.datadog_synthetics_global_var_create
type syntheticsGlobalVarCreateStep struct {
	name       string
	moduleName string
}

func newSyntheticsGlobalVarCreateStep(name string, config map[string]any) (*syntheticsGlobalVarCreateStep, error) {
	return &syntheticsGlobalVarCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *syntheticsGlobalVarCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	varName := resolveValue("name", current, config)
	if varName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	value := resolveValue("value", current, config)
	body := datadogV1.SyntheticsGlobalVariableRequest{
		Name: varName,
		Value: &datadogV1.SyntheticsGlobalVariableValue{
			Value: &value,
		},
	}
	if desc := resolveValue("description", current, config); desc != "" {
		body.SetDescription(desc)
	}
	api := datadogV1.NewSyntheticsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateGlobalVariable(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":   derefStr(resp.Id),
		"name": resp.Name,
	}}, nil
}

// syntheticsGlobalVarListStep implements step.datadog_synthetics_global_var_list
type syntheticsGlobalVarListStep struct {
	name       string
	moduleName string
}

func newSyntheticsGlobalVarListStep(name string, config map[string]any) (*syntheticsGlobalVarListStep, error) {
	return &syntheticsGlobalVarListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *syntheticsGlobalVarListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV1.NewSyntheticsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListGlobalVariables(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	vars := make([]any, 0)
	for _, v := range resp.Variables {
		vars = append(vars, map[string]any{
			"id":   derefStr(v.Id),
			"name": v.Name,
		})
	}
	return &sdk.StepResult{Output: map[string]any{"variables": vars, "count": len(vars)}}, nil
}

// syntheticsGlobalVarDeleteStep implements step.datadog_synthetics_global_var_delete
type syntheticsGlobalVarDeleteStep struct {
	name       string
	moduleName string
}

func newSyntheticsGlobalVarDeleteStep(name string, config map[string]any) (*syntheticsGlobalVarDeleteStep, error) {
	return &syntheticsGlobalVarDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *syntheticsGlobalVarDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	varID := resolveValue("variable_id", current, config)
	if varID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "variable_id is required"}}, nil
	}
	api := datadogV1.NewSyntheticsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, err := api.DeleteGlobalVariable(ddCtx.ctx, varID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": varID}}, nil
}
