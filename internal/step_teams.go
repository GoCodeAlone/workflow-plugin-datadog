package internal

import (
	"context"

	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// teamCreateStep implements step.datadog_team_create
type teamCreateStep struct {
	name       string
	moduleName string
}

func newTeamCreateStep(name string, config map[string]any) (*teamCreateStep, error) {
	return &teamCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *teamCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	teamName := resolveValue("name", current, config)
	if teamName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	handle := resolveValue("handle", current, config)
	if handle == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "handle is required"}}, nil
	}
	body := datadogV2.TeamCreateRequest{
		Data: datadogV2.TeamCreate{
			Type: datadogV2.TEAMTYPE_TEAM,
			Attributes: datadogV2.TeamCreateAttributes{
				Name:   teamName,
				Handle: handle,
			},
		},
	}
	if desc := resolveValue("description", current, config); desc != "" {
		body.Data.Attributes.SetDescription(desc)
	}
	api := datadogV2.NewTeamsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.CreateTeam(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	teamID := ""
	if resp.Data != nil {
		teamID = resp.Data.Id
	}
	return &sdk.StepResult{Output: map[string]any{"id": teamID, "name": teamName, "handle": handle}}, nil
}

// teamGetStep implements step.datadog_team_get
type teamGetStep struct {
	name       string
	moduleName string
}

func newTeamGetStep(name string, config map[string]any) (*teamGetStep, error) {
	return &teamGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *teamGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	teamID := resolveValue("team_id", current, config)
	if teamID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "team_id is required"}}, nil
	}
	api := datadogV2.NewTeamsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.GetTeam(ddCtx.ctx, teamID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	teamName := ""
	if resp.Data != nil {
		teamName = resp.Data.Attributes.GetName()
	}
	return &sdk.StepResult{Output: map[string]any{"id": teamID, "name": teamName}}, nil
}

// teamUpdateStep implements step.datadog_team_update
type teamUpdateStep struct {
	name       string
	moduleName string
}

func newTeamUpdateStep(name string, config map[string]any) (*teamUpdateStep, error) {
	return &teamUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *teamUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	teamID := resolveValue("team_id", current, config)
	if teamID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "team_id is required"}}, nil
	}
	teamName := resolveValue("name", current, config)
	handle := resolveValue("handle", current, config)
	if handle == "" {
		handle = teamID
	}
	body := datadogV2.TeamUpdateRequest{
		Data: datadogV2.TeamUpdate{
			Type: datadogV2.TEAMTYPE_TEAM,
			Attributes: datadogV2.TeamUpdateAttributes{
				Name:   teamName,
				Handle: handle,
			},
		},
	}
	api := datadogV2.NewTeamsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.UpdateTeam(ddCtx.ctx, teamID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	updatedID := ""
	if resp.Data != nil {
		updatedID = resp.Data.Id
	}
	return &sdk.StepResult{Output: map[string]any{"id": updatedID, "updated": true}}, nil
}

// teamDeleteStep implements step.datadog_team_delete
type teamDeleteStep struct {
	name       string
	moduleName string
}

func newTeamDeleteStep(name string, config map[string]any) (*teamDeleteStep, error) {
	return &teamDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *teamDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	teamID := resolveValue("team_id", current, config)
	if teamID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "team_id is required"}}, nil
	}
	api := datadogV2.NewTeamsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	_, err := api.DeleteTeam(ddCtx.ctx, teamID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": teamID}}, nil
}

// teamListStep implements step.datadog_team_list
type teamListStep struct {
	name       string
	moduleName string
}

func newTeamListStep(name string, config map[string]any) (*teamListStep, error) {
	return &teamListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *teamListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewTeamsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.ListTeams(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	teams := make([]any, 0)
	for _, t := range resp.Data {
		teams = append(teams, map[string]any{
			"id":   t.Id,
			"name": t.Attributes.GetName(),
		})
	}
	return &sdk.StepResult{Output: map[string]any{"teams": teams, "count": len(teams)}}, nil
}

// teamMemberAddStep implements step.datadog_team_member_add
type teamMemberAddStep struct {
	name       string
	moduleName string
}

func newTeamMemberAddStep(name string, config map[string]any) (*teamMemberAddStep, error) {
	return &teamMemberAddStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *teamMemberAddStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	teamID := resolveValue("team_id", current, config)
	userID := resolveValue("user_id", current, config)
	if teamID == "" || userID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "team_id and user_id are required"}}, nil
	}
	role := datadogV2.USERTEAMROLE_ADMIN
	attrs := datadogV2.UserTeamAttributes{}
	attrs.SetRole(role)
	body := datadogV2.UserTeamRequest{
		Data: datadogV2.UserTeamCreate{
			Type:       datadogV2.USERTEAMTYPE_TEAM_MEMBERSHIPS,
			Attributes: &attrs,
			Relationships: &datadogV2.UserTeamRelationships{
				User: &datadogV2.RelationshipToUserTeamUser{
					Data: datadogV2.RelationshipToUserTeamUserData{
						Type: datadogV2.USERTEAMUSERTYPE_USERS,
						Id:   userID,
					},
				},
			},
		},
	}
	api := datadogV2.NewTeamsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	resp, _, err := api.CreateTeamMembership(ddCtx.ctx, teamID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	memberID := ""
	if resp.Data != nil {
		memberID = resp.Data.Id
	}
	return &sdk.StepResult{Output: map[string]any{"id": memberID, "added": true}}, nil
}

// teamMemberRemoveStep implements step.datadog_team_member_remove
type teamMemberRemoveStep struct {
	name       string
	moduleName string
}

func newTeamMemberRemoveStep(name string, config map[string]any) (*teamMemberRemoveStep, error) {
	return &teamMemberRemoveStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *teamMemberRemoveStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	teamID := resolveValue("team_id", current, config)
	userID := resolveValue("user_id", current, config)
	if teamID == "" || userID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "team_id and user_id are required"}}, nil
	}
	api := datadogV2.NewTeamsApi(datadog.NewAPIClient(ddCtx.newConfig()))
	_, err := api.DeleteTeamMembership(ddCtx.ctx, teamID, userID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"removed": true, "team_id": teamID, "user_id": userID}}, nil
}
