package cmd

import (
	"encoding/json"
)

type AppEnvironment struct {
	AppGUID     string
	Environment map[string]interface{} `json:"var"`
}

func getAppEnvironmentVariables(app AppResource, include_env_variables bool, config Config) AppEnvironment {
	// defer wg.Done()

	var appEnvironment AppEnvironment

	appEnvironment.AppGUID = app.GUID

	if include_env_variables {

		apiUrl := app.AppLinks.EnvironmentVars.Href

		output, _ := getResponse(config, apiUrl)
		json.Unmarshal([]byte(output), &appEnvironment)
	}

	return appEnvironment
}
