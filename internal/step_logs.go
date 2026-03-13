package internal

import (
	"context"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// logSubmitStep implements step.datadog_log_submit
type logSubmitStep struct {
	name       string
	moduleName string
}

func newLogSubmitStep(name string, config map[string]any) (*logSubmitStep, error) {
	return &logSubmitStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *logSubmitStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	message := resolveValue("message", current, config)
	if message == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "message is required"}}, nil
	}
	entry := datadogV2.HTTPLogItem{
		Message: message,
	}
	if service := resolveValue("service", current, config); service != "" {
		entry.SetService(service)
	}
	if source := resolveValue("source", current, config); source != "" {
		entry.SetDdsource(source)
	}
	if host := resolveValue("host", current, config); host != "" {
		entry.SetHostname(host)
	}
	if tags := resolveValue("tags", current, config); tags != "" {
		entry.SetDdtags(tags)
	}
	api := datadogV2.NewLogsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, _, err := api.SubmitLog(ddCtx.ctx, []datadogV2.HTTPLogItem{entry})
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"submitted": true}}, nil
}

// logSearchStep implements step.datadog_log_search
type logSearchStep struct {
	name       string
	moduleName string
}

func newLogSearchStep(name string, config map[string]any) (*logSearchStep, error) {
	return &logSearchStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *logSearchStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	body := datadogV1.LogsListRequest{
		Query: &query,
		Limit: datadog.PtrInt32(100),
		Time:  datadogV1.LogsListRequestTime{},
	}
	api := datadogV1.NewLogsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListLogs(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	logs := make([]any, 0)
	for _, l := range resp.Logs {
		msg := ""
		if l.Content != nil {
			msg = derefStr(l.Content.Message)
		}
		logs = append(logs, map[string]any{
			"id":      derefStr(l.Id),
			"message": msg,
		})
	}
	return &sdk.StepResult{Output: map[string]any{"logs": logs, "count": len(logs)}}, nil
}

// logAggregateStep implements step.datadog_log_aggregate
type logAggregateStep struct {
	name       string
	moduleName string
}

func newLogAggregateStep(name string, config map[string]any) (*logAggregateStep, error) {
	return &logAggregateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *logAggregateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	body := datadogV2.LogsAggregateRequest{
		Filter: &datadogV2.LogsQueryFilter{
			Query: &query,
		},
	}
	api := datadogV2.NewLogsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.AggregateLogs(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	buckets := make([]any, 0)
	if resp.Data != nil {
		for _, b := range resp.Data.Buckets {
			buckets = append(buckets, map[string]any{
				"by": b.By,
			})
		}
	}
	return &sdk.StepResult{Output: map[string]any{"buckets": buckets, "count": len(buckets)}}, nil
}

// logArchiveCreateStep implements step.datadog_log_archive_create
type logArchiveCreateStep struct {
	name       string
	moduleName string
}

func newLogArchiveCreateStep(name string, config map[string]any) (*logArchiveCreateStep, error) {
	return &logArchiveCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *logArchiveCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	archiveName := resolveValue("name", current, config)
	if archiveName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	query := resolveValue("query", current, config)
	// Default: use S3 destination type with a placeholder
	bucket := resolveValue("bucket", current, config)
	if bucket == "" {
		bucket = "my-bucket"
	}
	path := resolveValue("path", current, config)
	roleArn := resolveValue("role_arn", current, config)
	dest := datadogV2.LogsArchiveDestinationS3{
		Bucket:      bucket,
		Path:        &path,
		Integration: datadogV2.LogsArchiveIntegrationS3{AccountId: "123456789012", RoleName: roleArn},
		Type:        datadogV2.LOGSARCHIVEDESTINATIONS3TYPE_S3,
	}
	body := datadogV2.LogsArchiveCreateRequest{
		Data: &datadogV2.LogsArchiveCreateRequestDefinition{
			Attributes: &datadogV2.LogsArchiveCreateRequestAttributes{
				Name:        archiveName,
				Query:       query,
				Destination: datadogV2.LogsArchiveCreateRequestDestination{LogsArchiveDestinationS3: &dest},
			},
			Type: "archives",
		},
	}
	api := datadogV2.NewLogsArchivesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateLogsArchive(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	archiveID := ""
	if resp.Data != nil {
		archiveID = derefStr(resp.Data.Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": archiveID, "name": archiveName}}, nil
}

// logArchiveListStep implements step.datadog_log_archive_list
type logArchiveListStep struct {
	name       string
	moduleName string
}

func newLogArchiveListStep(name string, config map[string]any) (*logArchiveListStep, error) {
	return &logArchiveListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *logArchiveListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewLogsArchivesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListLogsArchives(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	archives := make([]any, 0)
	for _, a := range resp.Data {
		item := map[string]any{"id": derefStr(a.Id), "type": string(a.Type)}
		if a.Attributes != nil {
			item["name"] = a.Attributes.GetName()
		}
		archives = append(archives, item)
	}
	return &sdk.StepResult{Output: map[string]any{"archives": archives, "count": len(archives)}}, nil
}

// logArchiveDeleteStep implements step.datadog_log_archive_delete
type logArchiveDeleteStep struct {
	name       string
	moduleName string
}

func newLogArchiveDeleteStep(name string, config map[string]any) (*logArchiveDeleteStep, error) {
	return &logArchiveDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *logArchiveDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	archiveID := resolveValue("archive_id", current, config)
	if archiveID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "archive_id is required"}}, nil
	}
	api := datadogV2.NewLogsArchivesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, err := api.DeleteLogsArchive(ddCtx.ctx, archiveID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": archiveID}}, nil
}

// logPipelineCreateStep implements step.datadog_log_pipeline_create
type logPipelineCreateStep struct {
	name       string
	moduleName string
}

func newLogPipelineCreateStep(name string, config map[string]any) (*logPipelineCreateStep, error) {
	return &logPipelineCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *logPipelineCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	pipelineName := resolveValue("name", current, config)
	if pipelineName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	body := datadogV1.LogsPipeline{
		Name:       pipelineName,
		Processors: []datadogV1.LogsProcessor{},
	}
	api := datadogV1.NewLogsPipelinesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateLogsPipeline(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":   derefStr(resp.Id),
		"name": resp.Name,
	}}, nil
}

// logPipelineListStep implements step.datadog_log_pipeline_list
type logPipelineListStep struct {
	name       string
	moduleName string
}

func newLogPipelineListStep(name string, config map[string]any) (*logPipelineListStep, error) {
	return &logPipelineListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *logPipelineListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV1.NewLogsPipelinesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListLogsPipelines(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	pipelines := make([]any, 0, len(resp))
	for _, p := range resp {
		pipelines = append(pipelines, map[string]any{
			"id":   derefStr(p.Id),
			"name": p.Name,
		})
	}
	return &sdk.StepResult{Output: map[string]any{"pipelines": pipelines, "count": len(pipelines)}}, nil
}

// logPipelineDeleteStep implements step.datadog_log_pipeline_delete
type logPipelineDeleteStep struct {
	name       string
	moduleName string
}

func newLogPipelineDeleteStep(name string, config map[string]any) (*logPipelineDeleteStep, error) {
	return &logPipelineDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *logPipelineDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	pipelineID := resolveValue("pipeline_id", current, config)
	if pipelineID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "pipeline_id is required"}}, nil
	}
	api := datadogV1.NewLogsPipelinesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, err := api.DeleteLogsPipeline(ddCtx.ctx, pipelineID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": pipelineID}}, nil
}
