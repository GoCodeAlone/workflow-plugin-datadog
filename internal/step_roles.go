package internal

import (
	"context"

	datadogV2 "github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// roleCreateStep implements step.datadog_role_create
type roleCreateStep struct {
	name       string
	moduleName string
}

func newRoleCreateStep(name string, config map[string]any) (*roleCreateStep, error) {
	return &roleCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *roleCreateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	roleName := resolveValue("name", current, config)
	if roleName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "name is required"}}, nil
	}
	body := datadogV2.RoleCreateRequest{
		Data: datadogV2.RoleCreateData{
			Type: datadogV2.ROLESTYPE_ROLES.Ptr(),
			Attributes: datadogV2.RoleCreateAttributes{
				Name: roleName,
			},
		},
	}
	api := datadogV2.NewRolesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.CreateRole(ddCtx.ctx, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	roleID := ""
	if resp.Data != nil {
		roleID = derefStr(resp.Data.Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": roleID, "name": roleName}}, nil
}

// roleGetStep implements step.datadog_role_get
type roleGetStep struct {
	name       string
	moduleName string
}

func newRoleGetStep(name string, config map[string]any) (*roleGetStep, error) {
	return &roleGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *roleGetStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	roleID := resolveValue("role_id", current, config)
	if roleID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "role_id is required"}}, nil
	}
	api := datadogV2.NewRolesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.GetRole(ddCtx.ctx, roleID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	roleName := ""
	if resp.Data != nil && resp.Data.Attributes != nil {
		roleName = resp.Data.Attributes.GetName()
	}
	return &sdk.StepResult{Output: map[string]any{"id": roleID, "name": roleName}}, nil
}

// roleUpdateStep implements step.datadog_role_update
type roleUpdateStep struct {
	name       string
	moduleName string
}

func newRoleUpdateStep(name string, config map[string]any) (*roleUpdateStep, error) {
	return &roleUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *roleUpdateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	roleID := resolveValue("role_id", current, config)
	if roleID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "role_id is required"}}, nil
	}
	roleName := resolveValue("name", current, config)
	body := datadogV2.RoleUpdateRequest{
		Data: datadogV2.RoleUpdateData{
			Type: datadogV2.ROLESTYPE_ROLES,
			Id:   roleID,
			Attributes: datadogV2.RoleUpdateAttributes{
				Name: &roleName,
			},
		},
	}
	api := datadogV2.NewRolesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.UpdateRole(ddCtx.ctx, roleID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	updatedID := ""
	if resp.Data != nil {
		updatedID = derefStr(resp.Data.Id)
	}
	return &sdk.StepResult{Output: map[string]any{"id": updatedID, "updated": true}}, nil
}

// roleDeleteStep implements step.datadog_role_delete
type roleDeleteStep struct {
	name       string
	moduleName string
}

func newRoleDeleteStep(name string, config map[string]any) (*roleDeleteStep, error) {
	return &roleDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *roleDeleteStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	roleID := resolveValue("role_id", current, config)
	if roleID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "role_id is required"}}, nil
	}
	api := datadogV2.NewRolesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, err := api.DeleteRole(ddCtx.ctx, roleID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "id": roleID}}, nil
}

// roleListStep implements step.datadog_role_list
type roleListStep struct {
	name       string
	moduleName string
}

func newRoleListStep(name string, config map[string]any) (*roleListStep, error) {
	return &roleListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *roleListStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	api := datadogV2.NewRolesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.ListRoles(ddCtx.ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	roles := make([]any, 0)
	for _, r := range resp.Data {
			roleName := ""
			if r.Attributes != nil {
				roleName = r.Attributes.GetName()
			}
			roles = append(roles, map[string]any{
				"id":   derefStr(r.Id),
				"name": roleName,
			})
		}
	return &sdk.StepResult{Output: map[string]any{"roles": roles, "count": len(roles)}}, nil
}

// rolePermissionAddStep implements step.datadog_role_permission_add
type rolePermissionAddStep struct {
	name       string
	moduleName string
}

func newRolePermissionAddStep(name string, config map[string]any) (*rolePermissionAddStep, error) {
	return &rolePermissionAddStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *rolePermissionAddStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	roleID := resolveValue("role_id", current, config)
	permissionID := resolveValue("permission_id", current, config)
	if roleID == "" || permissionID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "role_id and permission_id are required"}}, nil
	}
	body := datadogV2.RelationshipToPermission{
		Data: &datadogV2.RelationshipToPermissionData{
			Type: datadogV2.PERMISSIONSTYPE_PERMISSIONS.Ptr(),
			Id:   &permissionID,
		},
	}
	api := datadogV2.NewRolesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	resp, _, err := api.AddPermissionToRole(ddCtx.ctx, roleID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	count := len(resp.Data)
	return &sdk.StepResult{Output: map[string]any{"added": true, "permission_count": count}}, nil
}

// rolePermissionRemoveStep implements step.datadog_role_permission_remove
type rolePermissionRemoveStep struct {
	name       string
	moduleName string
}

func newRolePermissionRemoveStep(name string, config map[string]any) (*rolePermissionRemoveStep, error) {
	return &rolePermissionRemoveStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *rolePermissionRemoveStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	ddCtx, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "datadog client not found: " + s.moduleName}}, nil
	}
	roleID := resolveValue("role_id", current, config)
	permissionID := resolveValue("permission_id", current, config)
	if roleID == "" || permissionID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "role_id and permission_id are required"}}, nil
	}
	body := datadogV2.RelationshipToPermission{
		Data: &datadogV2.RelationshipToPermissionData{
			Type: datadogV2.PERMISSIONSTYPE_PERMISSIONS.Ptr(),
			Id:   &permissionID,
		},
	}
	api := datadogV2.NewRolesApi(datadog.NewAPIClient(datadog.NewConfiguration()))
	_, _, err := api.RemovePermissionFromRole(ddCtx.ctx, roleID, body)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"removed": true, "role_id": roleID, "permission_id": permissionID}}, nil
}
