package cmd

import (
	"encoding/json"
)

func getAppEnvironmentVariables(app AppResource, includeEnvVariables bool, config Config) AppEnvironment {
	var appEnvironment AppEnvironment

	appEnvironment.AppGUID = app.GUID

	if includeEnvVariables {
		apiUrl := app.AppLinks.EnvironmentVars.Href

		output, _ := getResponse(config, apiUrl)
		json.Unmarshal([]byte(output), &appEnvironment)
	}

	return appEnvironment
}
