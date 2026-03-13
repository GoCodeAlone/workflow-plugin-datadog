package internal

import (
	"context"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

func defaultNotebookTime() datadogV1.NotebookGlobalTime {
	return datadogV1.NotebookRelativeTimeAsNotebookGlobalTime(
		datadogV1.NewNotebookRelativeTime(datadogV1.WIDGETLIVESPAN_PAST_ONE_HOUR),
	)
}

// notebookCreateStep implements step.datadog_notebook_create
type notebookCreateStep struct {
	name       string
	moduleName string
}

func newNotebookCreateStep(name string, config map[string]any) (*notebookCreateStep, error) {
	return &notebookCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *notebookCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	notebookName := resolveValue("name", current, config)
	if notebookName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	body := datadogV1.NotebookCreateRequest{
		Data: datadogV1.NotebookCreateData{
			Type: datadogV1.NOTEBOOKRESOURCETYPE_NOTEBOOKS,
			Attributes: datadogV1.NotebookCreateDataAttributes{
				Name:  notebookName,
				Cells: []datadogV1.NotebookCellCreateRequest{},
				Time:  defaultNotebookTime(),
			},
		},
	}
	api := datadogV1.NewNotebooksApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateNotebook(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	notebookID := int64(0)
	if resp.Data != nil {
		notebookID = resp.Data.GetId()
	}
	return &sdk.StepResult{Output: map[string]any{"id": notebookID, "name": notebookName}}, nil
}

// notebookGetStep implements step.datadog_notebook_get
type notebookGetStep struct {
	name       string
	moduleName string
}

func newNotebookGetStep(name string, config map[string]any) (*notebookGetStep, error) {
	return &notebookGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *notebookGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	notebookID := resolveInt64("notebook_id", current, config)
	if notebookID == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "notebook_id is required"}}, nil
	}
	api := datadogV1.NewNotebooksApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetNotebook(ddCtx.ctx, notebookID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	notebookName := ""
	if resp.Data != nil {
		notebookName = resp.Data.Attributes.GetName()
	}
	return &sdk.StepResult{Output: map[string]any{"id": notebookID, "name": notebookName}}, nil
}

// notebookUpdateStep implements step.datadog_notebook_update
type notebookUpdateStep struct {
	name       string
	moduleName string
}

func newNotebookUpdateStep(name string, config map[string]any) (*notebookUpdateStep, error) {
	return &notebookUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *notebookUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	notebookID := resolveInt64("notebook_id", current, config)
	if notebookID == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "notebook_id is required"}}, nil
	}
	notebookName := resolveValue("name", current, config)
	body := datadogV1.NotebookUpdateRequest{
		Data: datadogV1.NotebookUpdateData{
			Type: datadogV1.NOTEBOOKRESOURCETYPE_NOTEBOOKS,
			Attributes: datadogV1.NotebookUpdateDataAttributes{
				Name:  notebookName,
				Cells: []datadogV1.NotebookUpdateCell{},
				Time:  defaultNotebookTime(),
			},
		},
	}
	api := datadogV1.NewNotebooksApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.UpdateNotebook(ddCtx.ctx, notebookID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	updatedID := int64(0)
	if resp.Data != nil {
		updatedID = resp.Data.GetId()
	}
	return &sdk.StepResult{Output: map[string]any{"id": updatedID, "updated": true}}, nil
}

// notebookDeleteStep implements step.datadog_notebook_delete
type notebookDeleteStep struct {
	name       string
	moduleName string
}

func newNotebookDeleteStep(name string, config map[string]any) (*notebookDeleteStep, error) {
	return &notebookDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *notebookDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	notebookID := resolveInt64("notebook_id", current, config)
	if notebookID == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "notebook_id is required"}}, nil
	}
	api := datadogV1.NewNotebooksApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, err := api.DeleteNotebook(ddCtx.ctx, notebookID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": notebookID}}, nil
}

// notebookListStep implements step.datadog_notebook_list
type notebookListStep struct {
	name       string
	moduleName string
}

func newNotebookListStep(name string, config map[string]any) (*notebookListStep, error) {
	return &notebookListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *notebookListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	params := datadogV1.NewListNotebooksOptionalParameters()
	if q := resolveValue("query", current, config); q != "" {
		params.WithQuery(q)
	}
	api := datadogV1.NewNotebooksApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListNotebooks(ddCtx.ctx, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	notebooks := make([]any, 0)
	for _, nb := range resp.Data {
		nbName := nb.Attributes.GetName()
		notebooks = append(notebooks, map[string]any{
			"id":   nb.GetId(),
			"name": nbName,
		})
	}
	return &sdk.StepResult{Output: map[string]any{"notebooks": notebooks, "count": len(notebooks)}}, nil
}
