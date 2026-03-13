package internal

import (
	"context"
	"time"

	datadogV1 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// eventCreateStep implements step.datadog_event_create
type eventCreateStep struct {
	name       string
	moduleName string
}

func newEventCreateStep(name string, config map[string]any) (*eventCreateStep, error) {
	return &eventCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	title := resolveValue("title", current, config)
	if title == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "title is required"}}, nil
	}
	text := resolveValue("text", current, config)
	body := datadogV1.EventCreateRequest{
		Title: title,
		Text:  text,
	}
	if tags := resolveStringSlice("tags", current, config); len(tags) > 0 {
		body.SetTags(tags)
	}
	if host := resolveValue("host", current, config); host != "" {
		body.SetHost(host)
	}
	api := datadogV1.NewEventsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateEvent(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	eventID := int64(0)
	if resp.Event != nil && resp.Event.Id != nil {
		eventID = *resp.Event.Id
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":     eventID,
		"title":  title,
		"status": derefStr(resp.Status),
	}}, nil
}

// eventGetStep implements step.datadog_event_get
type eventGetStep struct {
	name       string
	moduleName string
}

func newEventGetStep(name string, config map[string]any) (*eventGetStep, error) {
	return &eventGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	eventID := resolveInt64("event_id", current, config)
	if eventID == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "event_id is required"}}, nil
	}
	api := datadogV1.NewEventsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetEvent(ddCtx.ctx, eventID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	out := map[string]any{"found": resp.Event != nil}
	if resp.Event != nil {
		e := resp.Event
		out["id"] = derefInt64(e.Id)
		out["title"] = derefStr(e.Title)
		out["text"] = derefStr(e.Text)
		out["host"] = derefStr(e.Host)
	}
	return &sdk.StepResult{Output: out}, nil
}

// eventListStep implements step.datadog_event_list
type eventListStep struct {
	name       string
	moduleName string
}

func newEventListStep(name string, config map[string]any) (*eventListStep, error) {
	return &eventListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	start := resolveInt64("start", current, config)
	end := resolveInt64("end", current, config)
	if start == 0 {
		start = time.Now().Add(-1 * time.Hour).Unix()
	}
	if end == 0 {
		end = time.Now().Unix()
	}
	params := datadogV1.NewListEventsOptionalParameters()
	if tags := resolveValue("tags", current, config); tags != "" {
		params.WithTags(tags)
	}
	api := datadogV1.NewEventsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListEvents(ddCtx.ctx, start, end, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	events := make([]any, 0)
	for _, e := range resp.Events {
		events = append(events, map[string]any{
			"id":    derefInt64(e.Id),
			"title": derefStr(e.Title),
			"text":  derefStr(e.Text),
			"host":  derefStr(e.Host),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"events": events, "count": len(events)}}, nil
}

// eventSearchStep implements step.datadog_event_search (alias for list with query)
type eventSearchStep struct {
	name       string
	moduleName string
}

func newEventSearchStep(name string, config map[string]any) (*eventSearchStep, error) {
	return &eventSearchStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventSearchStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	start := resolveInt64("start", current, config)
	end := resolveInt64("end", current, config)
	if start == 0 {
		start = time.Now().Add(-1 * time.Hour).Unix()
	}
	if end == 0 {
		end = time.Now().Unix()
	}
	params := datadogV1.NewListEventsOptionalParameters()
	if tags := resolveValue("tags", current, config); tags != "" {
		params.WithTags(tags)
	}
	if sources := resolveValue("sources", current, config); sources != "" {
		params.WithSources(sources)
	}
	api := datadogV1.NewEventsApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListEvents(ddCtx.ctx, start, end, *params)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	events := make([]any, 0)
	for _, e := range resp.Events {
		events = append(events, map[string]any{
			"id":    derefInt64(e.Id),
			"title": derefStr(e.Title),
			"text":  derefStr(e.Text),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"events": events, "count": len(events)}}, nil
}
