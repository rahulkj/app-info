package cmd

import (
	"encoding/json"
)

type AppFeatures struct {
	AppGUID  string
	Features []AppFeatureResource `json:"resources"`
}

type AppFeatureResource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func getAppFeatures(app AppResource, config Config) AppFeatures {
	// defer wg.Done()

	var appFeatures AppFeatures

	appFeatures.AppGUID = app.GUID

	apiUrl := app.AppLinks.Features.Href

	output, _ := getResponse(config, apiUrl)
	json.Unmarshal([]byte(output), &appFeatures)

	return appFeatures
}
