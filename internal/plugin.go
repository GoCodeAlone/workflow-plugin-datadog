// Package internal implements the workflow-plugin-datadog plugin.
package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// Version is set at build time via -ldflags
// "-X github.com/GoCodeAlone/workflow-plugin-datadog/internal.Version=X.Y.Z".
// Default is a bare semver so plugin loaders that validate semver accept
// unreleased dev builds; goreleaser overrides with the real release tag.
var Version = "0.0.0"

// datadogPlugin implements sdk.PluginProvider, sdk.ModuleProvider, sdk.StepProvider,
// and sdk.SchemaProvider.
type datadogPlugin struct{}

// NewDatadogPlugin returns a new datadogPlugin instance.
func NewDatadogPlugin() sdk.PluginProvider {
	return &datadogPlugin{}
}

// Manifest returns plugin metadata.
func (p *datadogPlugin) Manifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Name:        "workflow-plugin-datadog",
		Version:     Version,
		Author:      "GoCodeAlone",
		Description: "Datadog observability platform plugin (~120 step types across all Datadog APIs)",
	}
}

// ModuleTypes returns the module type names this plugin provides.
func (p *datadogPlugin) ModuleTypes() []string {
	return []string{"datadog.provider"}
}

// CreateModule creates a module instance of the given type.
func (p *datadogPlugin) CreateModule(typeName, name string, config map[string]any) (sdk.ModuleInstance, error) {
	switch typeName {
	case "datadog.provider":
		m, err := newDatadogModule(name, config)
		if err != nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, fmt.Errorf("datadog plugin: unknown module type %q", typeName)
	}
}

// StepTypes returns the step type names this plugin provides.
func (p *datadogPlugin) StepTypes() []string {
	return allStepTypes()
}

// CreateStep creates a step instance of the given type.
func (p *datadogPlugin) CreateStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	return createStep(typeName, name, config)
}

// ModuleSchemas implements sdk.SchemaProvider. It returns contract descriptors
// for each module type this plugin provides, enabling strict-contract audits.
func (p *datadogPlugin) ModuleSchemas() []sdk.ModuleSchemaData {
	return []sdk.ModuleSchemaData{
		{
			Type:        "datadog.provider",
			Label:       "Datadog Provider",
			Category:    "observability",
			Description: "Authenticates with the Datadog API and makes a scoped client available to all Datadog step types.",
			ConfigFields: []sdk.ConfigField{
				{
					Name:        "apiKey",
					Type:        "string",
					Description: "Datadog API key (required).",
					Required:    true,
				},
				{
					Name:        "appKey",
					Type:        "string",
					Description: "Datadog application key (required).",
					Required:    true,
				},
				{
					Name:         "site",
					Type:         "string",
					Description:  "Datadog site (e.g. datadoghq.com, datadoghq.eu). Defaults to datadoghq.com.",
					DefaultValue: "datadoghq.com",
				},
				{
					Name:        "apiUrl",
					Type:        "string",
					Description: "Optional custom API URL override (useful for mocking/testing).",
				},
			},
		},
	}
}
