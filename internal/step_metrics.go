package internal

import (
	"context"
	"time"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// metricSubmitStep implements step.datadog_metric_submit
type metricSubmitStep struct {
	name       string
	moduleName string
}

func newMetricSubmitStep(name string, config map[string]any) (*metricSubmitStep, error) {
	return &metricSubmitStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metricSubmitStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	metricName := resolveValue("metric", current, config)
	if metricName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metric is required"}}, nil
	}
	value := resolveFloat64("value", current, config)
	metricType := datadogV2.METRICINTAKETYPE_UNSPECIFIED
	if t := resolveValue("type", current, config); t != "" {
		switch t {
		case "count":
			metricType = datadogV2.METRICINTAKETYPE_COUNT
		case "rate":
			metricType = datadogV2.METRICINTAKETYPE_RATE
		case "gauge":
			metricType = datadogV2.METRICINTAKETYPE_GAUGE
		}
	}
	ts := time.Now().Unix()
	tags := resolveStringSlice("tags", current, config)

	point := datadogV2.MetricPoint{
		Timestamp: datadog.PtrInt64(ts),
		Value:     datadog.PtrFloat64(value),
	}
	series := datadogV2.MetricSeries{
		Metric: metricName,
		Type:   &metricType,
		Points: []datadogV2.MetricPoint{point},
	}
	if len(tags) > 0 {
		series.Tags = tags
	}
	body := datadogV2.MetricPayload{Series: []datadogV2.MetricSeries{series}}
	api := datadogV2.NewMetricsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	_, _, err := api.SubmitMetrics(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"submitted": true, "metric": metricName}}, nil
}

// metricQueryStep implements step.datadog_metric_query
type metricQueryStep struct {
	name       string
	moduleName string
}

func newMetricQueryStep(name string, config map[string]any) (*metricQueryStep, error) {
	return &metricQueryStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metricQueryStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	if query == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "query is required"}}, nil
	}
	fromTs := resolveInt64("from", current, config)
	toTs := resolveInt64("to", current, config)
	if fromTs == 0 {
		fromTs = time.Now().Add(-1 * time.Hour).Unix()
	}
	if toTs == 0 {
		toTs = time.Now().Unix()
	}
	api := datadogV1.NewMetricsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.QueryMetrics(ddCtx.ctx, fromTs, toTs, query)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	series := make([]any, 0)
	for _, s := range resp.Series {
		series = append(series, map[string]any{
			"metric": derefStr(s.Metric),
			"scope":  derefStr(s.Scope),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"series": series, "count": len(series)}}, nil
}

// metricQueryScalarStep implements step.datadog_metric_query_scalar
type metricQueryScalarStep struct {
	name       string
	moduleName string
}

func newMetricQueryScalarStep(name string, config map[string]any) (*metricQueryScalarStep, error) {
	return &metricQueryScalarStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metricQueryScalarStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	if query == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "query is required"}}, nil
	}
	fromTs := resolveInt64("from", current, config)
	toTs := resolveInt64("to", current, config)
	if fromTs == 0 {
		fromTs = time.Now().Add(-1 * time.Hour).Unix()
	}
	if toTs == 0 {
		toTs = time.Now().Unix()
	}
	queryName := "a"
	scalarReq := datadogV2.ScalarFormulaRequest{
		Attributes: datadogV2.ScalarFormulaRequestAttributes{
			Formulas: []datadogV2.QueryFormula{{Formula: queryName}},
			From:     fromTs * 1000,
			To:       toTs * 1000,
			Queries: []datadogV2.ScalarQuery{
				datadogV2.MetricsScalarQueryAsScalarQuery(&datadogV2.MetricsScalarQuery{
					Aggregator: datadogV2.METRICSAGGREGATOR_AVG,
					DataSource: datadogV2.METRICSDATASOURCE_METRICS,
					Query:      query,
					Name:       &queryName,
				}),
			},
		},
		Type: datadogV2.SCALARFORMULAREQUESTTYPE_SCALAR_REQUEST,
	}
	api := datadogV2.NewMetricsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.QueryScalarData(ddCtx.ctx, datadogV2.ScalarFormulaQueryRequest{Data: scalarReq})
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	_ = resp
	return &sdk.StepResult{Output: map[string]any{"queried": true, "query": query}}, nil
}

// metricMetadataGetStep implements step.datadog_metric_metadata_get
type metricMetadataGetStep struct {
	name       string
	moduleName string
}

func newMetricMetadataGetStep(name string, config map[string]any) (*metricMetadataGetStep, error) {
	return &metricMetadataGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metricMetadataGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	metricName := resolveValue("metric", current, config)
	if metricName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metric is required"}}, nil
	}
	api := datadogV1.NewMetricsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	meta, _, err := api.GetMetricMetadata(ddCtx.ctx, metricName)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"type":        derefStr(meta.Type),
		"description": derefStr(meta.Description),
		"unit":        derefStr(meta.Unit),
		"short_name":  derefStr(meta.ShortName),
	}}, nil
}

// metricMetadataUpdateStep implements step.datadog_metric_metadata_update
type metricMetadataUpdateStep struct {
	name       string
	moduleName string
}

func newMetricMetadataUpdateStep(name string, config map[string]any) (*metricMetadataUpdateStep, error) {
	return &metricMetadataUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metricMetadataUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	metricName := resolveValue("metric", current, config)
	if metricName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metric is required"}}, nil
	}
	body := datadogV1.MetricMetadata{}
	if desc := resolveValue("description", current, config); desc != "" {
		body.SetDescription(desc)
	}
	if unit := resolveValue("unit", current, config); unit != "" {
		body.SetUnit(unit)
	}
	if t := resolveValue("type", current, config); t != "" {
		body.SetType(t)
	}
	api := datadogV1.NewMetricsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	meta, _, err := api.UpdateMetricMetadata(ddCtx.ctx, metricName, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"type":        derefStr(meta.Type),
		"description": derefStr(meta.Description),
		"unit":        derefStr(meta.Unit),
	}}, nil
}

// metricListActiveStep implements step.datadog_metric_list_active
type metricListActiveStep struct {
	name       string
	moduleName string
}

func newMetricListActiveStep(name string, config map[string]any) (*metricListActiveStep, error) {
	return &metricListActiveStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metricListActiveStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	from := resolveInt64("from", current, config)
	if from == 0 {
		from = time.Now().Add(-1 * time.Hour).Unix()
	}
	params := datadogV1.NewListActiveMetricsOptionalParameters()
	if host := resolveValue("host", current, config); host != "" {
		params.WithHost(host)
	}
	api := datadogV1.NewMetricsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.ListActiveMetrics(ddCtx.ctx, from, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	metrics := make([]any, 0)
	for _, m := range resp.Metrics {
		metrics = append(metrics, m)
	}
	return &sdk.StepResult{Output: map[string]any{"metrics": metrics, "count": len(metrics)}}, nil
}

// metricTagConfigCreateStep implements step.datadog_metric_tag_config_create
type metricTagConfigCreateStep struct {
	name       string
	moduleName string
}

func newMetricTagConfigCreateStep(name string, config map[string]any) (*metricTagConfigCreateStep, error) {
	return &metricTagConfigCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metricTagConfigCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	metricName := resolveValue("metric", current, config)
	if metricName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metric is required"}}, nil
	}
	tags := resolveStringSlice("tags", current, config)
	metricType := datadogV2.METRICTAGCONFIGURATIONMETRICTYPES_GAUGE
	body := datadogV2.MetricTagConfigurationCreateRequest{
		Data: datadogV2.MetricTagConfigurationCreateData{
			Type: datadogV2.METRICTAGCONFIGURATIONTYPE_MANAGE_TAGS,
			Id:   metricName,
			Attributes: &datadogV2.MetricTagConfigurationCreateAttributes{
				Tags:       tags,
				MetricType: metricType,
			},
		},
	}
	api := datadogV2.NewMetricsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.CreateTagConfiguration(ddCtx.ctx, metricName, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	_ = resp
	return &sdk.StepResult{Output: map[string]any{"created": true, "metric": metricName}}, nil
}

// metricTagConfigUpdateStep implements step.datadog_metric_tag_config_update
type metricTagConfigUpdateStep struct {
	name       string
	moduleName string
}

func newMetricTagConfigUpdateStep(name string, config map[string]any) (*metricTagConfigUpdateStep, error) {
	return &metricTagConfigUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metricTagConfigUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	metricName := resolveValue("metric", current, config)
	if metricName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metric is required"}}, nil
	}
	tags := resolveStringSlice("tags", current, config)
	body := datadogV2.MetricTagConfigurationUpdateRequest{
		Data: datadogV2.MetricTagConfigurationUpdateData{
			Type: datadogV2.METRICTAGCONFIGURATIONTYPE_MANAGE_TAGS,
			Id:   metricName,
			Attributes: &datadogV2.MetricTagConfigurationUpdateAttributes{
				Tags: tags,
			},
		},
	}
	api := datadogV2.NewMetricsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.UpdateTagConfiguration(ddCtx.ctx, metricName, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	_ = resp
	return &sdk.StepResult{Output: map[string]any{"updated": true, "metric": metricName}}, nil
}

// metricTagConfigDeleteStep implements step.datadog_metric_tag_config_delete
type metricTagConfigDeleteStep struct {
	name       string
	moduleName string
}

func newMetricTagConfigDeleteStep(name string, config map[string]any) (*metricTagConfigDeleteStep, error) {
	return &metricTagConfigDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metricTagConfigDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	metricName := resolveValue("metric", current, config)
	if metricName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "metric is required"}}, nil
	}
	api := datadogV2.NewMetricsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	_, err := api.DeleteTagConfiguration(ddCtx.ctx, metricName)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "metric": metricName}}, nil
}

// metricTagConfigListStep implements step.datadog_metric_tag_config_list
type metricTagConfigListStep struct {
	name       string
	moduleName string
}

func newMetricTagConfigListStep(name string, config map[string]any) (*metricTagConfigListStep, error) {
	return &metricTagConfigListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *metricTagConfigListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewMetricsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.ListTagConfigurations(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	items := make([]any, 0)
	for _, d := range resp.Data {
		id := ""
		if d.Metric != nil {
			id = d.Metric.GetId()
		} else if d.MetricTagConfiguration != nil {
			id = d.MetricTagConfiguration.GetId()
		}
		items = append(items, map[string]any{
			"id": id,
		})
	}
	return &sdk.StepResult{Output: map[string]any{"configurations": items, "count": len(items)}}, nil
}
