package internal

import (
	"context"
	"fmt"

	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// apmRetentionFilterCreateStep implements step.datadog_apm_retention_filter_create
type apmRetentionFilterCreateStep struct {
	name       string
	moduleName string
}

func newAPMRetentionFilterCreateStep(name string, config map[string]any) (*apmRetentionFilterCreateStep, error) {
	return &apmRetentionFilterCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apmRetentionFilterCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	filterName := resolveValue("name", current, config)
	if filterName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	rate := resolveFloat64("rate", current, config)
	if rate == 0 {
		rate = 1.0
	}
	q := resolveValue("query", current, config)
	if q == "" {
		q = "*"
	}
	body := datadogV2.RetentionFilterCreateRequest{
		Data: datadogV2.RetentionFilterCreateData{
			Attributes: datadogV2.RetentionFilterCreateAttributes{
				Name:       filterName,
				Rate:       rate,
				Enabled:    true,
				FilterType: datadogV2.RETENTIONFILTERTYPE_SPANS_SAMPLING_PROCESSOR,
				Filter:     datadogV2.SpansFilterCreate{Query: q},
			},
			Type: datadogV2.APMRETENTIONFILTERTYPE_apm_retention_filter,
		},
	}
	api := datadogV2.NewAPMRetentionFiltersApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.CreateApmRetentionFilter(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	filterID := ""
	if resp.Data != nil {
		filterID = resp.Data.Id
	}
	return &sdk.StepResult{Output: map[string]any{"id": filterID, "name": filterName}}, nil
}

// apmRetentionFilterUpdateStep implements step.datadog_apm_retention_filter_update
type apmRetentionFilterUpdateStep struct {
	name       string
	moduleName string
}

func newAPMRetentionFilterUpdateStep(name string, config map[string]any) (*apmRetentionFilterUpdateStep, error) {
	return &apmRetentionFilterUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apmRetentionFilterUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	filterID := resolveValue("filter_id", current, config)
	if filterID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "filter_id is required"}}, nil
	}
	filterName := resolveValue("name", current, config)
	rate := resolveFloat64("rate", current, config)
	if rate == 0 {
		rate = 1.0
	}
	q := resolveValue("query", current, config)
	if q == "" {
		q = "*"
	}
	body := datadogV2.RetentionFilterUpdateRequest{
		Data: datadogV2.RetentionFilterUpdateData{
			Id:   filterID,
			Type: datadogV2.APMRETENTIONFILTERTYPE_apm_retention_filter,
			Attributes: datadogV2.RetentionFilterUpdateAttributes{
				Name:       filterName,
				Rate:       rate,
				Enabled:    true,
				FilterType: datadogV2.RETENTIONFILTERALLTYPE_SPANS_SAMPLING_PROCESSOR,
				Filter:     datadogV2.SpansFilterCreate{Query: q},
			},
		},
	}
	api := datadogV2.NewAPMRetentionFiltersApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.UpdateApmRetentionFilter(ddCtx.ctx, filterID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	updatedID := ""
	if resp.Data != nil {
		updatedID = resp.Data.Id
	}
	return &sdk.StepResult{Output: map[string]any{"id": updatedID, "updated": true}}, nil
}

// apmRetentionFilterDeleteStep implements step.datadog_apm_retention_filter_delete
type apmRetentionFilterDeleteStep struct {
	name       string
	moduleName string
}

func newAPMRetentionFilterDeleteStep(name string, config map[string]any) (*apmRetentionFilterDeleteStep, error) {
	return &apmRetentionFilterDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apmRetentionFilterDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	filterID := resolveValue("filter_id", current, config)
	if filterID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "filter_id is required"}}, nil
	}
	api := datadogV2.NewAPMRetentionFiltersApi(datadog.NewAPIClient(ddCtx.newConfig()))
	_, err := api.DeleteApmRetentionFilter(ddCtx.ctx, filterID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": filterID}}, nil
}

// apmRetentionFilterListStep implements step.datadog_apm_retention_filter_list
type apmRetentionFilterListStep struct {
	name       string
	moduleName string
}

func newAPMRetentionFilterListStep(name string, config map[string]any) (*apmRetentionFilterListStep, error) {
	return &apmRetentionFilterListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *apmRetentionFilterListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewAPMRetentionFiltersApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.ListApmRetentionFilters(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	filters := make([]any, 0)
	for _, f := range resp.Data {
		filters = append(filters, map[string]any{
			"id":   f.Id,
			"name": f.Attributes.GetName(),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"filters": filters, "count": len(filters)}}, nil
}

// spanSearchStep implements step.datadog_span_search
type spanSearchStep struct {
	name       string
	moduleName string
}

func newSpanSearchStep(name string, config map[string]any) (*spanSearchStep, error) {
	return &spanSearchStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *spanSearchStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	reqType := datadogV2.SPANSLISTREQUESTTYPE_SEARCH_REQUEST
	body := datadogV2.SpansListRequest{
		Data: &datadogV2.SpansListRequestData{
			Attributes: &datadogV2.SpansListRequestAttributes{
				Filter: &datadogV2.SpansQueryFilter{
					Query: &query,
				},
			},
			Type: &reqType,
		},
	}
	api := datadogV2.NewSpansApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.ListSpans(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	spans := make([]any, 0)
	for _, sp := range resp.Data {
		spans = append(spans, map[string]any{
			"id":   derefStr(sp.Id),
			"type": fmt.Sprintf("%v", sp.GetType()),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"spans": spans, "count": len(spans)}}, nil
}

// spanAggregateStep implements step.datadog_span_aggregate
type spanAggregateStep struct {
	name       string
	moduleName string
}

func newSpanAggregateStep(name string, config map[string]any) (*spanAggregateStep, error) {
	return &spanAggregateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *spanAggregateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	reqType := datadogV2.SPANSAGGREGATEREQUESTTYPE_AGGREGATE_REQUEST
	body := datadogV2.SpansAggregateRequest{
		Data: &datadogV2.SpansAggregateData{
			Attributes: &datadogV2.SpansAggregateRequestAttributes{
				Filter: &datadogV2.SpansQueryFilter{
					Query: &query,
				},
			},
			Type: &reqType,
		},
	}
	api := datadogV2.NewSpansApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.AggregateSpans(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	buckets := make([]any, 0)
	for _, b := range resp.Data {
		buckets = append(buckets, map[string]any{
			"id": derefStr(b.Id),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"buckets": buckets, "count": len(buckets)}}, nil
}
