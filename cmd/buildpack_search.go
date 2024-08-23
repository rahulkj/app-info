package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
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

func getBuildpacks(cli plugin.CliConnection) map[string]BuildpackResources {
	data := make(map[string]BuildpackResources)
	buildpacks := getBuildPacksData(cli)

	for _, val := range buildpacks.Resources {
		data[val.Name] = val
	}

	return data
}

// GetOrgData requests all of the Application data from Cloud Foundry
func getBuildPacksData(cli plugin.CliConnection) Buildpacks {
	var res Buildpacks = unmarshallBuildpackSearchResults("/v3/buildpacks", cli)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v3/buildpacks?page=%d&per_page=50", i)
			tRes := unmarshallBuildpackSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallBuildpackSearchResults(apiUrl string, cli plugin.CliConnection) Buildpacks {
	var tRes Buildpacks
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}

func getBuildpackDetails(app DisplayApp, buildpacks map[string]BuildpackResources) (displayApp DisplayApp) {
	for _, buildpack := range app.Buildpacks {
		app.DetectedBuildPackFileNames = append(app.DetectedBuildPackFileNames, buildpacks[buildpack].Filename)
	}

	return app
}
