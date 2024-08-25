package cmd

import (
	"encoding/json"
)

func getAppFeatures(app AppResource, config Config) AppFeatures {
	// defer wg.Done()

	var appFeatures AppFeatures

	appFeatures.AppGUID = app.GUID

	apiUrl := app.AppLinks.Features.Href

	output, _ := getResponse(config, apiUrl)
	json.Unmarshal([]byte(output), &appFeatures)

	return appFeatures
}
