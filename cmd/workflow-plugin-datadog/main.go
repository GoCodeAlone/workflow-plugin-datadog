package main

import (
	"github.com/GoCodeAlone/workflow-plugin-datadog/internal"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

var version = "dev"

func main() {
	sdk.Serve(internal.NewDatadogPlugin())
}
