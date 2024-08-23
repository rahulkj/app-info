package cmd

import (
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type AppFeatures struct {
	Features []AppFeatureResource `json:"resources"`
}

type AppFeatureResource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func getAppFeatures(app DisplayApp, cli plugin.CliConnection) (displayApp DisplayApp) {
	var appFeatures AppFeatures
	cmd := []string{"curl", "/v3/apps/" + app.AppGUID + "/features"}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &appFeatures)

	app.Features = appFeatures.Features
	return app
}
