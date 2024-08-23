package cmd

import (
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type AppEnvironment struct {
	Environment map[string]string `json:"var"`
}

func getAppEnvironmentVariables(app DisplayApp, cli plugin.CliConnection) (displayApp DisplayApp) {
	var appEnvironment AppEnvironment
	cmd := []string{"curl", "/v3/apps/" + app.AppGUID + "/environment_variables"}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &appEnvironment)

	app.Environment = appEnvironment.Environment

	return app
}
