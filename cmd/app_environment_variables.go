package cmd

import (
	"encoding/json"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type AppEnvironment struct {
	Environment map[string]string `json:"var"`
}

func getAppEnvironmentVariables(app AppResource, include_env_variables bool, cli plugin.CliConnection, displayAppChan chan DisplayApp) {
	var displayApp DisplayApp

	if include_env_variables {

		var appEnvironment AppEnvironment

		envVarsUrl, _ := url.Parse(app.AppLinks.Features.Href)

		cmd := []string{"curl", envVarsUrl.Path}

		output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
		json.Unmarshal([]byte(strings.Join(output, "")), &appEnvironment)

		displayApp.Environment = appEnvironment.Environment
	}

	displayAppChan <- displayApp
}
