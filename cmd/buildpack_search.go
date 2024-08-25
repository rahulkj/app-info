package cmd

import (
	"encoding/json"
	"fmt"
)

type Buildpacks struct {
	Resources  []BuildpackResources `json:"resources"`
	Pagination BuildpacksPagination `json:"pagination"`
}

type BuildpacksPagination struct {
	TotalPages int `json:"total_pages"`
}

type BuildpackResources struct {
	GUID     string `json:"guid"`
	Name     string `json:"name"`
	Stack    string `json:"stack"`
	State    string `json:"state"`
	Position int    `json:"position"`
	Filename string `json:"filename"`
	Enabled  bool   `json:"enabled"`
	Locked   bool   `json:"locked"`
}

type AppDetectedBuildpacks struct {
	AppGUID                    string
	DetectedBuildPackFileNames []string
}

func getBuildpacks(config Config) map[string]BuildpackResources {
	data := make(map[string]BuildpackResources)
	buildpacks := getBuildPacksData(config)

	for _, val := range buildpacks.Resources {
		data[val.Name] = val
	}

	return data
}

// GetOrgData requests all of the Application data from Cloud Foundry
func getBuildPacksData(config Config) Buildpacks {
	apiUrl := fmt.Sprintf("%s/v3/buildpacks", config.ApiEndpoint)
	var res Buildpacks = unmarshallBuildpackSearchResults(apiUrl, config)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("%s?page=%d&per_page=100", apiUrl, i)
			tRes := unmarshallBuildpackSearchResults(apiUrl, config)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallBuildpackSearchResults(apiUrl string, config Config) Buildpacks {
	var tRes Buildpacks

	output, _ := getResponse(config, apiUrl)
	json.Unmarshal([]byte(output), &tRes)

	return tRes
}

func getBuildpackDetails(app AppResource, buildpacks map[string]BuildpackResources) AppDetectedBuildpacks {
	// defer wg.Done()
	var appDetectedBuildpacks AppDetectedBuildpacks

	for _, buildpack := range app.Lifecycle.Data.Buildpacks {
		appDetectedBuildpacks.AppGUID = app.GUID
		appDetectedBuildpacks.DetectedBuildPackFileNames = append(appDetectedBuildpacks.DetectedBuildPackFileNames, buildpacks[buildpack].Filename)
	}

	return appDetectedBuildpacks
}
