package internal

import (
	"context"

	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// apiKeyCreateStep implements step.datadog_api_key_create
type apiKeyCreateStep struct {
	name       string
	moduleName string
}

func newAPIKeyCreateStep(name string, config map[string]any) (*apiKeyCreateStep, error) {
	return &apiKeyCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apiKeyCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	keyName := resolveValue("name", current, config)
	if keyName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	body := datadogV2.APIKeyCreateRequest{
		Data: datadogV2.APIKeyCreateData{
			Type: datadogV2.APIKEYSTYPE_API_KEYS,
			Attributes: datadogV2.APIKeyCreateAttributes{
				Name: keyName,
			},
		},
	}
	api := datadogV2.NewKeyManagementApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.CreateAPIKey(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	keyID := ""
	if resp.Data != nil {
		keyID = derefStr(resp.Data.Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": keyID, "name": keyName}}, nil
}

// apiKeyGetStep implements step.datadog_api_key_get
type apiKeyGetStep struct {
	name       string
	moduleName string
}

func newAPIKeyGetStep(name string, config map[string]any) (*apiKeyGetStep, error) {
	return &apiKeyGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apiKeyGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	keyID := resolveValue("key_id", current, config)
	if keyID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "key_id is required"}}, nil
	}
	api := datadogV2.NewKeyManagementApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.GetAPIKey(ddCtx.ctx, keyID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	keyName := ""
	if resp.Data != nil && resp.Data.Attributes != nil {
		keyName = resp.Data.Attributes.GetName()
	}
	return &sdk.StepResult{Output: map[string]any{"id": keyID, "name": keyName}}, nil
}

// apiKeyUpdateStep implements step.datadog_api_key_update
type apiKeyUpdateStep struct {
	name       string
	moduleName string
}

func newAPIKeyUpdateStep(name string, config map[string]any) (*apiKeyUpdateStep, error) {
	return &apiKeyUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apiKeyUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	keyID := resolveValue("key_id", current, config)
	if keyID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "key_id is required"}}, nil
	}
	keyName := resolveValue("name", current, config)
	body := datadogV2.APIKeyUpdateRequest{
		Data: datadogV2.APIKeyUpdateData{
			Type: datadogV2.APIKEYSTYPE_API_KEYS,
			Id:   keyID,
			Attributes: datadogV2.APIKeyUpdateAttributes{
				Name: keyName,
			},
		},
	}
	api := datadogV2.NewKeyManagementApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.UpdateAPIKey(ddCtx.ctx, keyID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	updatedID := ""
	if resp.Data != nil {
		updatedID = derefStr(resp.Data.Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": updatedID, "updated": true}}, nil
}

// apiKeyDeleteStep implements step.datadog_api_key_delete
type apiKeyDeleteStep struct {
	name       string
	moduleName string
}

func newAPIKeyDeleteStep(name string, config map[string]any) (*apiKeyDeleteStep, error) {
	return &apiKeyDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apiKeyDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	keyID := resolveValue("key_id", current, config)
	if keyID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "key_id is required"}}, nil
	}
	api := datadogV2.NewKeyManagementApi(datadog.NewAPIClient(ddCtx.newConfig()))
	_, err := api.DeleteAPIKey(ddCtx.ctx, keyID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": keyID}}, nil
}

// apiKeyListStep implements step.datadog_api_key_list
type apiKeyListStep struct {
	name       string
	moduleName string
}

func newAPIKeyListStep(name string, config map[string]any) (*apiKeyListStep, error) {
	return &apiKeyListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apiKeyListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewKeyManagementApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.ListAPIKeys(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	keys := make([]any, 0)
	for _, k := range resp.Data {
			keyName := ""
			if k.Attributes != nil {
				keyName = k.Attributes.GetName()
			}
			keys = append(keys, map[string]any{
				"id":   derefStr(k.Id),
				"name": keyName,
			})
		}
	return &sdk.StepResult{Output: map[string]any{"keys": keys, "count": len(keys)}}, nil
}

// appKeyCreateStep implements step.datadog_app_key_create
type appKeyCreateStep struct {
	name       string
	moduleName string
}

func newAppKeyCreateStep(name string, config map[string]any) (*appKeyCreateStep, error) {
	return &appKeyCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appKeyCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	keyName := resolveValue("name", current, config)
	if keyName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	body := datadogV2.ApplicationKeyCreateRequest{
		Data: datadogV2.ApplicationKeyCreateData{
			Type: datadogV2.APPLICATIONKEYSTYPE_APPLICATION_KEYS,
			Attributes: datadogV2.ApplicationKeyCreateAttributes{
				Name: keyName,
			},
		},
	}
	api := datadogV2.NewKeyManagementApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.CreateCurrentUserApplicationKey(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	keyID := ""
	if resp.Data != nil {
		keyID = derefStr(resp.Data.Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": keyID, "name": keyName}}, nil
}

// appKeyListStep implements step.datadog_app_key_list
type appKeyListStep struct {
	name       string
	moduleName string
}

func newAppKeyListStep(name string, config map[string]any) (*appKeyListStep, error) {
	return &appKeyListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appKeyListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewKeyManagementApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.ListApplicationKeys(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	keys := make([]any, 0)
	for _, k := range resp.Data {
			keyName := ""
			if k.Attributes != nil {
				keyName = k.Attributes.GetName()
			}
			keys = append(keys, map[string]any{
				"id":   derefStr(k.Id),
				"name": keyName,
			})
		}
	return &sdk.StepResult{Output: map[string]any{"keys": keys, "count": len(keys)}}, nil
}

// appKeyDeleteStep implements step.datadog_app_key_delete
type appKeyDeleteStep struct {
	name       string
	moduleName string
}

func newAppKeyDeleteStep(name string, config map[string]any) (*appKeyDeleteStep, error) {
	return &appKeyDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appKeyDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	keyID := resolveValue("key_id", current, config)
	if keyID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "key_id is required"}}, nil
	}
	api := datadogV2.NewKeyManagementApi(datadog.NewAPIClient(ddCtx.newConfig()))
	_, err := api.DeleteCurrentUserApplicationKey(ddCtx.ctx, keyID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": keyID}}, nil
}
