package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cloudfoundry/cli/plugin"

	"github.com/schollz/progressbar/v3"
)

// AppSearchResults represents top level attributes of JSON response from Cloud Foundry API
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

func GatherData(cli plugin.CliConnection, include_env_variables bool) (map[string]string, map[string]SpaceSearchResource, []DisplayApp) {
	orgs := getOrgs(cli)
	spaces := getSpaces(cli)
	apps := getAppData(cli)
	buildpacks := getBuildpacks(cli)
	routes := getAllRoutes(cli)
	services := getAllServices(cli)

	var displayApps []DisplayApp

	var wg sync.WaitGroup
	results := make(chan DisplayApp)

	bar := progressbar.Default(int64(len(apps.Resources)))
	for _, appResource := range apps.Resources {
		wg.Add(1)
		go getAppResourceData(appResource, routes, services, buildpacks, include_env_variables, cli, results, &wg)
	}

	for i := 0; i < len(apps.Resources); i++ {
		bar.Add(1)
		newapp := <-results
		displayApps = append(displayApps, newapp)
	}

	wg.Wait()

	return orgs, spaces, displayApps
}

func getAppResourceData(app AppResource, routes Routes, services []Service, buildpacks map[string]BuildpackResources, include_env_variables bool, cli plugin.CliConnection, results chan DisplayApp, wg *sync.WaitGroup) {
	defer wg.Done()

	var displayApp DisplayApp

	displayApp.Name = app.Name
	displayApp.AppGUID = app.GUID
	displayApp.Stack = app.Lifecycle.Data.Stack
	displayApp.Buildpacks = app.Lifecycle.Data.Buildpacks
	displayApp.State = app.State
	displayApp.SpaceGUID = app.RelationShips.Space.Data.SpaceGUID

	displayApp = getAppProcesses(displayApp, cli)

	displayApp = getAppFeatures(displayApp, cli)

	displayApp = getBuildpackDetails(displayApp, buildpacks)

	if include_env_variables {
		displayApp = getAppEnvironmentVariables(displayApp, cli)
	}

	displayApp = getAppRoutes(displayApp, routes, cli)
	displayApp = getAppServices(displayApp, services)

	time.Sleep(time.Second)

	results <- displayApp

}
