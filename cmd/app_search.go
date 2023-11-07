package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

type Apps struct {
	Resources  []AppResource  `json:"resources"`
	Pagination AppsPagination `json:"pagination"`
}

type AppsPagination struct {
	TotalPages int `json:"total_pages"`
}

type AppResource struct {
	GUID          string           `json:"guid"`
	Name          string           `json:"name"`
	State         string           `json:"state"`
	Lifecycle     Lifecycle        `json:"lifecycle"`
	RelationShips AppRelationShips `json:"relationships"`
}

type Lifecycle struct {
	Type string        `json:"type"`
	Data LifecycleData `json:"data"`
}

type LifecycleData struct {
	Buildpacks []string `json:"buildpacks"`
	Stack      string   `json:"stack"`
}

type AppRelationShips struct {
	Space AppSpace `json:"space"`
}

type AppSpace struct {
	Data AppSpaceData `json:"data"`
}

type AppSpaceData struct {
	SpaceGUID string `json:"guid"`
}

type DisplayApp struct {
	Name                       string               `json:"name"`
	AppGUID                    string               `json:"guid"`
	Instances                  int                  `json:"instances"`
	State                      string               `json:"state"`
	Memory                     int                  `json:"memory_in_mb"`
	Disk                       int                  `json:"disk_in_mb"`
	LogRate                    int                  `json:"log_rate_limit_in_bytes_per_second"`
	Buildpacks                 []string             `json:"buildpacks"`
	DetectedBuildPack          string               `json:"detected_buildpack"`
	DetectedBuildPackFileNames []string             `json:"detected_buildpack_filenames"`
	SpaceGUID                  string               `json:"space_guid"`
	StartCommand               string               `json:"detected_start_command"`
	Environment                map[string]string    `json:"environment_json"`
	Command                    string               `json:"command"`
	HealthCheck                string               `json:"health_check_type"`
	HealthCheckEndpoint        string               `json:"health_check_http_endpoint"`
	Routes                     []string             `json:"routes"`
	Stack                      string               `json:"stack"`
	Services                   []Service            `json:"services"`
	Features                   []AppFeatureResource `json:"resources"`
	StackGUID                  string               `json:"stackguid"`
}

// GetAppData requests all of the Application data from Cloud Foundry
func getAppData(cli plugin.CliConnection) Apps {
	fmt.Println("**** Gathering application metadata from all orgs and spaces ****")

	res := unmarshallAppSearchResults("/v3/apps", cli)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v3/apps?page=%d&per_page=50", i)

			tRes := unmarshallAppSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallAppSearchResults(apiUrl string, cli plugin.CliConnection) Apps {
	var tRes Apps
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)

	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}

func GatherData(cli plugin.CliConnection) (map[string]string, map[string]SpaceSearchResource, []DisplayApp) {
	orgs := getOrgs(cli)
	spaces := getSpaces(cli)
	apps := getAppData(cli)
	buildpacks := getBuildpacks(cli)
	routes := getAllRoutes(cli)
	services := getAllServices(cli)

	var displayApps []DisplayApp

	for _, app := range apps.Resources {
		var displayApp DisplayApp

		displayApp.Name = app.Name
		displayApp.AppGUID = app.GUID
		displayApp.Stack = app.Lifecycle.Data.Stack
		displayApp.Buildpacks = app.Lifecycle.Data.Buildpacks
		displayApp.State = app.State
		displayApp.SpaceGUID = app.RelationShips.Space.Data.SpaceGUID

		getAppFeatures(&displayApp, cli)
		getBuildpackDetails(&displayApp, buildpacks)
		getAppEnvironmentVariables(&displayApp, cli)
		getAppProcesses(&displayApp, cli)
		getAppRoutes(&displayApp, routes, cli)
		getAppServices(&displayApp, services)

		displayApps = append(displayApps, displayApp)
	}

	return orgs, spaces, displayApps
}
