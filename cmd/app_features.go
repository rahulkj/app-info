package cmd

import (
	"encoding/json"
	"net/url"
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

func getAppFeatures(app AppResource, cli plugin.CliConnection, displayAppChan chan<- DisplayApp) {
	var displayApp DisplayApp

	var appFeatures AppFeatures

	featuresUrl, _ := url.Parse(app.AppLinks.Features.Href)

	cmd := []string{"curl", featuresUrl.Path}

	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &appFeatures)

	displayApp.Features = appFeatures.Features

	displayAppChan <- displayApp
}
