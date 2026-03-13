package internal

import (
	"context"

	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// userCreateStep implements step.datadog_user_create
type userCreateStep struct {
	name       string
	moduleName string
}

func newUserCreateStep(name string, config map[string]any) (*userCreateStep, error) {
	return &userCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	email := resolveValue("email", current, config)
	if email == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "email is required"}}, nil
	}
	body := datadogV2.UserCreateRequest{
		Data: datadogV2.UserCreateData{
			Type: datadogV2.USERSTYPE_USERS,
			Attributes: datadogV2.UserCreateAttributes{
				Email: email,
			},
		},
	}
	if name := resolveValue("name", current, config); name != "" {
		body.Data.Attributes.SetName(name)
	}
	api := datadogV2.NewUsersApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateUser(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	userID := ""
	if resp.Data != nil {
		userID = derefStr(resp.Data.Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": userID, "email": email}}, nil
}

// userGetStep implements step.datadog_user_get
type userGetStep struct {
	name       string
	moduleName string
}

func newUserGetStep(name string, config map[string]any) (*userGetStep, error) {
	return &userGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	userID := resolveValue("user_id", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "user_id is required"}}, nil
	}
	api := datadogV2.NewUsersApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetUser(ddCtx.ctx, userID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	email := ""
	userName := ""
	if resp.Data != nil && resp.Data.Attributes != nil {
		email = derefStr(resp.Data.Attributes.Email)
		userName = derefStr(resp.Data.Attributes.Name.Get())
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":    userID,
		"email": email,
		"name":  userName,
	}}, nil
}

// userUpdateStep implements step.datadog_user_update
type userUpdateStep struct {
	name       string
	moduleName string
}

func newUserUpdateStep(name string, config map[string]any) (*userUpdateStep, error) {
	return &userUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	userID := resolveValue("user_id", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "user_id is required"}}, nil
	}
	body := datadogV2.UserUpdateRequest{
		Data: datadogV2.UserUpdateData{
			Type:       datadogV2.USERSTYPE_USERS,
			Id:         userID,
			Attributes: datadogV2.UserUpdateAttributes{},
		},
	}
	if name := resolveValue("name", current, config); name != "" {
		body.Data.Attributes.SetName(name)
	}
	api := datadogV2.NewUsersApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.UpdateUser(ddCtx.ctx, userID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	updatedID := ""
	if resp.Data != nil {
		updatedID = derefStr(resp.Data.Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": updatedID, "updated": true}}, nil
}

// userDisableStep implements step.datadog_user_disable
type userDisableStep struct {
	name       string
	moduleName string
}

func newUserDisableStep(name string, config map[string]any) (*userDisableStep, error) {
	return &userDisableStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userDisableStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	userID := resolveValue("user_id", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "user_id is required"}}, nil
	}
	api := datadogV2.NewUsersApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, err := api.DisableUser(ddCtx.ctx, userID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"disabled": true, "id": userID}}, nil
}

// userListStep implements step.datadog_user_list
type userListStep struct {
	name       string
	moduleName string
}

func newUserListStep(name string, config map[string]any) (*userListStep, error) {
	return &userListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewUsersApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListUsers(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	users := make([]any, 0)
	for _, u := range resp.Data {
			email := ""
			if u.Attributes != nil {
				email = derefStr(u.Attributes.Email)
			}
			users = append(users, map[string]any{
				"id":    derefStr(u.Id),
				"email": email,
			})
		}
	return &sdk.StepResult{Output: map[string]any{"users": users, "count": len(users)}}, nil
}

// userInviteStep implements step.datadog_user_invite
type userInviteStep struct {
	name       string
	moduleName string
}

func newUserInviteStep(name string, config map[string]any) (*userInviteStep, error) {
	return &userInviteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userInviteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	userID := resolveValue("user_id", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "user_id is required"}}, nil
	}
	body := datadogV2.UserInvitationsRequest{
		Data: []datadogV2.UserInvitationData{
			{
				Type: datadogV2.USERINVITATIONSTYPE_USER_INVITATIONS,
				Relationships: datadogV2.UserInvitationRelationships{
					User: datadogV2.RelationshipToUser{
						Data: datadogV2.RelationshipToUserData{
							Type: datadogV2.USERSTYPE_USERS,
							Id:   userID,
						},
					},
				},
			},
		},
	}
	api := datadogV2.NewUsersApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.SendInvitations(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	invitations := make([]any, 0)
	for _, inv := range resp.Data {
			invitations = append(invitations, map[string]any{
				"id": derefStr(inv.Id),
			})
		}
	return &sdk.StepResult{Output: map[string]any{"invitations": invitations, "count": len(invitations)}}, nil
}
