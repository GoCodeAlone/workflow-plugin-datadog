package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// stepConstructor is a function that creates a StepInstance.
type stepConstructor func(name string, config map[string]any) (sdk.StepInstance, error)

// stepRegistry maps step type strings to constructor functions.
var stepRegistry = map[string]stepConstructor{
	// Metrics
	"step.datadog_metric_submit":            func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetricSubmitStep(n, c) },
	"step.datadog_metric_query":             func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetricQueryStep(n, c) },
	"step.datadog_metric_query_scalar":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetricQueryScalarStep(n, c) },
	"step.datadog_metric_metadata_get":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetricMetadataGetStep(n, c) },
	"step.datadog_metric_metadata_update":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetricMetadataUpdateStep(n, c) },
	"step.datadog_metric_list_active":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetricListActiveStep(n, c) },
	"step.datadog_metric_tag_config_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetricTagConfigCreateStep(n, c) },
	"step.datadog_metric_tag_config_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetricTagConfigUpdateStep(n, c) },
	"step.datadog_metric_tag_config_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetricTagConfigDeleteStep(n, c) },
	"step.datadog_metric_tag_config_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetricTagConfigListStep(n, c) },

	// Events
	"step.datadog_event_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventCreateStep(n, c) },
	"step.datadog_event_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventGetStep(n, c) },
	"step.datadog_event_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventListStep(n, c) },
	"step.datadog_event_search": func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventSearchStep(n, c) },

	// Monitors
	"step.datadog_monitor_create":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newMonitorCreateStep(n, c) },
	"step.datadog_monitor_get":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newMonitorGetStep(n, c) },
	"step.datadog_monitor_update":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newMonitorUpdateStep(n, c) },
	"step.datadog_monitor_delete":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newMonitorDeleteStep(n, c) },
	"step.datadog_monitor_list":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newMonitorListStep(n, c) },
	"step.datadog_monitor_search":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newMonitorSearchStep(n, c) },
	"step.datadog_monitor_validate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newMonitorValidateStep(n, c) },

	// Dashboards
	"step.datadog_dashboard_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDashboardCreateStep(n, c) },
	"step.datadog_dashboard_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newDashboardGetStep(n, c) },
	"step.datadog_dashboard_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDashboardUpdateStep(n, c) },
	"step.datadog_dashboard_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDashboardDeleteStep(n, c) },
	"step.datadog_dashboard_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newDashboardListStep(n, c) },

	// Logs
	"step.datadog_log_submit":          func(n string, c map[string]any) (sdk.StepInstance, error) { return newLogSubmitStep(n, c) },
	"step.datadog_log_search":          func(n string, c map[string]any) (sdk.StepInstance, error) { return newLogSearchStep(n, c) },
	"step.datadog_log_aggregate":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newLogAggregateStep(n, c) },
	"step.datadog_log_archive_create":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newLogArchiveCreateStep(n, c) },
	"step.datadog_log_archive_list":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newLogArchiveListStep(n, c) },
	"step.datadog_log_archive_delete":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newLogArchiveDeleteStep(n, c) },
	"step.datadog_log_pipeline_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newLogPipelineCreateStep(n, c) },
	"step.datadog_log_pipeline_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newLogPipelineListStep(n, c) },
	"step.datadog_log_pipeline_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newLogPipelineDeleteStep(n, c) },

	// Synthetics
	"step.datadog_synthetics_test_create":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newSyntheticsTestCreateStep(n, c) },
	"step.datadog_synthetics_test_get":           func(n string, c map[string]any) (sdk.StepInstance, error) { return newSyntheticsTestGetStep(n, c) },
	"step.datadog_synthetics_test_update":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newSyntheticsTestUpdateStep(n, c) },
	"step.datadog_synthetics_test_delete":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newSyntheticsTestDeleteStep(n, c) },
	"step.datadog_synthetics_test_list":          func(n string, c map[string]any) (sdk.StepInstance, error) { return newSyntheticsTestListStep(n, c) },
	"step.datadog_synthetics_test_trigger":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newSyntheticsTestTriggerStep(n, c) },
	"step.datadog_synthetics_results_get":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newSyntheticsResultsGetStep(n, c) },
	"step.datadog_synthetics_global_var_create":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newSyntheticsGlobalVarCreateStep(n, c) },
	"step.datadog_synthetics_global_var_list":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newSyntheticsGlobalVarListStep(n, c) },
	"step.datadog_synthetics_global_var_delete":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newSyntheticsGlobalVarDeleteStep(n, c) },

	// SLOs
	"step.datadog_slo_create":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newSLOCreateStep(n, c) },
	"step.datadog_slo_get":         func(n string, c map[string]any) (sdk.StepInstance, error) { return newSLOGetStep(n, c) },
	"step.datadog_slo_update":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newSLOUpdateStep(n, c) },
	"step.datadog_slo_delete":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newSLODeleteStep(n, c) },
	"step.datadog_slo_list":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newSLOListStep(n, c) },
	"step.datadog_slo_search":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newSLOSearchStep(n, c) },
	"step.datadog_slo_history_get": func(n string, c map[string]any) (sdk.StepInstance, error) { return newSLOHistoryGetStep(n, c) },

	// Downtimes
	"step.datadog_downtime_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDowntimeCreateStep(n, c) },
	"step.datadog_downtime_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newDowntimeGetStep(n, c) },
	"step.datadog_downtime_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDowntimeUpdateStep(n, c) },
	"step.datadog_downtime_cancel": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDowntimeCancelStep(n, c) },
	"step.datadog_downtime_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newDowntimeListStep(n, c) },

	// Incidents
	"step.datadog_incident_create":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newIncidentCreateStep(n, c) },
	"step.datadog_incident_get":         func(n string, c map[string]any) (sdk.StepInstance, error) { return newIncidentGetStep(n, c) },
	"step.datadog_incident_update":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newIncidentUpdateStep(n, c) },
	"step.datadog_incident_delete":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newIncidentDeleteStep(n, c) },
	"step.datadog_incident_list":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newIncidentListStep(n, c) },
	"step.datadog_incident_todo_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newIncidentTodoCreateStep(n, c) },
	"step.datadog_incident_todo_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newIncidentTodoUpdateStep(n, c) },
	"step.datadog_incident_todo_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newIncidentTodoDeleteStep(n, c) },

	// Security
	"step.datadog_security_rule_create":         func(n string, c map[string]any) (sdk.StepInstance, error) { return newSecurityRuleCreateStep(n, c) },
	"step.datadog_security_rule_get":            func(n string, c map[string]any) (sdk.StepInstance, error) { return newSecurityRuleGetStep(n, c) },
	"step.datadog_security_rule_update":         func(n string, c map[string]any) (sdk.StepInstance, error) { return newSecurityRuleUpdateStep(n, c) },
	"step.datadog_security_rule_delete":         func(n string, c map[string]any) (sdk.StepInstance, error) { return newSecurityRuleDeleteStep(n, c) },
	"step.datadog_security_rule_list":           func(n string, c map[string]any) (sdk.StepInstance, error) { return newSecurityRuleListStep(n, c) },
	"step.datadog_security_signal_list":         func(n string, c map[string]any) (sdk.StepInstance, error) { return newSecuritySignalListStep(n, c) },
	"step.datadog_security_signal_state_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newSecuritySignalStateUpdateStep(n, c) },

	// Users
	"step.datadog_user_create":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserCreateStep(n, c) },
	"step.datadog_user_get":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserGetStep(n, c) },
	"step.datadog_user_update":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserUpdateStep(n, c) },
	"step.datadog_user_disable": func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserDisableStep(n, c) },
	"step.datadog_user_list":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserListStep(n, c) },
	"step.datadog_user_invite":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserInviteStep(n, c) },

	// Roles
	"step.datadog_role_create":           func(n string, c map[string]any) (sdk.StepInstance, error) { return newRoleCreateStep(n, c) },
	"step.datadog_role_get":              func(n string, c map[string]any) (sdk.StepInstance, error) { return newRoleGetStep(n, c) },
	"step.datadog_role_update":           func(n string, c map[string]any) (sdk.StepInstance, error) { return newRoleUpdateStep(n, c) },
	"step.datadog_role_delete":           func(n string, c map[string]any) (sdk.StepInstance, error) { return newRoleDeleteStep(n, c) },
	"step.datadog_role_list":             func(n string, c map[string]any) (sdk.StepInstance, error) { return newRoleListStep(n, c) },
	"step.datadog_role_permission_add":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newRolePermissionAddStep(n, c) },
	"step.datadog_role_permission_remove": func(n string, c map[string]any) (sdk.StepInstance, error) { return newRolePermissionRemoveStep(n, c) },

	// Teams
	"step.datadog_team_create":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newTeamCreateStep(n, c) },
	"step.datadog_team_get":           func(n string, c map[string]any) (sdk.StepInstance, error) { return newTeamGetStep(n, c) },
	"step.datadog_team_update":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newTeamUpdateStep(n, c) },
	"step.datadog_team_delete":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newTeamDeleteStep(n, c) },
	"step.datadog_team_list":          func(n string, c map[string]any) (sdk.StepInstance, error) { return newTeamListStep(n, c) },
	"step.datadog_team_member_add":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newTeamMemberAddStep(n, c) },
	"step.datadog_team_member_remove": func(n string, c map[string]any) (sdk.StepInstance, error) { return newTeamMemberRemoveStep(n, c) },

	// Key Management
	"step.datadog_api_key_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAPIKeyCreateStep(n, c) },
	"step.datadog_api_key_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newAPIKeyGetStep(n, c) },
	"step.datadog_api_key_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAPIKeyUpdateStep(n, c) },
	"step.datadog_api_key_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAPIKeyDeleteStep(n, c) },
	"step.datadog_api_key_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAPIKeyListStep(n, c) },
	"step.datadog_app_key_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppKeyCreateStep(n, c) },
	"step.datadog_app_key_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppKeyListStep(n, c) },
	"step.datadog_app_key_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppKeyDeleteStep(n, c) },

	// Notebooks
	"step.datadog_notebook_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newNotebookCreateStep(n, c) },
	"step.datadog_notebook_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newNotebookGetStep(n, c) },
	"step.datadog_notebook_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newNotebookUpdateStep(n, c) },
	"step.datadog_notebook_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newNotebookDeleteStep(n, c) },
	"step.datadog_notebook_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newNotebookListStep(n, c) },

	// Hosts
	"step.datadog_host_list":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newHostListStep(n, c) },
	"step.datadog_host_mute":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newHostMuteStep(n, c) },
	"step.datadog_host_unmute":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newHostUnmuteStep(n, c) },
	"step.datadog_host_totals_get": func(n string, c map[string]any) (sdk.StepInstance, error) { return newHostTotalsGetStep(n, c) },

	// Tags
	"step.datadog_tags_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newTagsGetStep(n, c) },
	"step.datadog_tags_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newTagsUpdateStep(n, c) },
	"step.datadog_tags_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newTagsDeleteStep(n, c) },
	"step.datadog_tags_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newTagsListStep(n, c) },

	// Service Catalog
	"step.datadog_service_definition_upsert": func(n string, c map[string]any) (sdk.StepInstance, error) { return newServiceDefinitionUpsertStep(n, c) },
	"step.datadog_service_definition_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newServiceDefinitionGetStep(n, c) },
	"step.datadog_service_definition_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newServiceDefinitionDeleteStep(n, c) },
	"step.datadog_service_definition_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newServiceDefinitionListStep(n, c) },

	// APM
	"step.datadog_apm_retention_filter_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAPMRetentionFilterCreateStep(n, c) },
	"step.datadog_apm_retention_filter_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAPMRetentionFilterUpdateStep(n, c) },
	"step.datadog_apm_retention_filter_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAPMRetentionFilterDeleteStep(n, c) },
	"step.datadog_apm_retention_filter_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAPMRetentionFilterListStep(n, c) },
	"step.datadog_span_search":                 func(n string, c map[string]any) (sdk.StepInstance, error) { return newSpanSearchStep(n, c) },
	"step.datadog_span_aggregate":              func(n string, c map[string]any) (sdk.StepInstance, error) { return newSpanAggregateStep(n, c) },

	// Audit
	"step.datadog_audit_log_search": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuditLogSearchStep(n, c) },
	"step.datadog_audit_log_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuditLogListStep(n, c) },
}

// createStep dispatches to the appropriate step constructor.
func createStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	constructor, ok := stepRegistry[typeName]
	if !ok {
		return nil, fmt.Errorf("datadog plugin: unknown step type %q", typeName)
	}
	return constructor(name, config)
}

// allStepTypes returns all registered step type strings.
func allStepTypes() []string {
	types := make([]string, 0, len(stepRegistry))
	for k := range stepRegistry {
		types = append(types, k)
	}
	return types
}
