package cmd

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type AppData interface{}

type BasicAppData struct {
	Name       string   `json:"name"`
	AppGUID    string   `json:"guid"`
	Stack      string   `json:"stack"`
	State      string   `json:"state"`
	Buildpacks []string `json:"buildpacks"`
	SpaceGUID  string   `json:"space_guid"`
}

// GetAppData requests all of the Application data from Cloud Foundry
func getAppsData(config Config) Apps {
	Yellow("**** Gathering application metadata from all orgs and spaces ****\n")

	res := unmarshallAppSearchResults("/v3/apps", config)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v3/apps?page=%d&per_page=100", i)

			tRes := unmarshallAppSearchResults(apiUrl, config)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallAppSearchResults(path string, config Config) Apps {
	var tRes Apps
	url := config.ApiEndpoint + path
	output, _ := getResponse(config, url)

	json.Unmarshal([]byte(output), &tRes)

	return tRes
}

func GatherData(config Config, includeEnvVariables bool) (map[string]string, map[string]SpaceSearchResource, []DisplayApp, []ServiceInstancesResource) {

	orgs := getOrgs(config)
	spaces := getSpaces(config)
	apps := getAppsData(config)
	buildpacks := getBuildpacks(config)
	routes := getAllRoutes(config)
	services, unboundServices := getServices(config)

	var displayApps []DisplayApp

	var wg sync.WaitGroup
	appDataCh := make(chan DisplayApp)

	bar := progressbar.Default(int64(len(apps.Resources)), "Gathering App Data")
	for _, appResource := range apps.Resources {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var displayApp DisplayApp
			basicData := getAppBasicData(appResource)
			displayApp.Name = basicData.Name
			displayApp.AppGUID = basicData.AppGUID
			displayApp.Stack = basicData.Stack
			displayApp.Buildpacks = basicData.Buildpacks
			displayApp.State = basicData.State
			displayApp.SpaceGUID = basicData.SpaceGUID

			disApp := getAppProcesses(appResource, config)
			displayApp.Instances = disApp.Instances
			displayApp.Memory = disApp.Memory
			displayApp.Disk = disApp.Disk
			displayApp.LogRate = disApp.LogRate
			displayApp.HealthCheck = disApp.HealthCheck
			displayApp.ReadinessHealthCheck = disApp.ReadinessHealthCheck
			displayApp.Type = disApp.Type

			appDetectedBuildpacks := getBuildpackDetails(appResource, buildpacks)
			displayApp.DetectedBuildPackFileNames = appDetectedBuildpacks.DetectedBuildPackFileNames

			appRoutes := getAppRoutes(appResource, routes)
			displayApp.Routes = appRoutes.Routes

			appEnvironment := getAppEnvironmentVariables(appResource, includeEnvVariables, config)
			displayApp.Environment = appEnvironment.Environment
			appFeatures := getAppFeatures(appResource, config)
			displayApp.Features = appFeatures.Features

			appServices := getAppSevices(appResource, services)
			displayApp.Services = appServices.Services

			appDataCh <- displayApp
		}()
	}

	go func() {
		wg.Wait()
		close(appDataCh) // Close results channel after all workers are done
	}()

	for data := range appDataCh {
		bar.Add(1)
		displayApps = append(displayApps, data)
	}

	return orgs, spaces, displayApps, unboundServices
}

func getAppBasicData(app AppResource) BasicAppData {
	var basicData BasicAppData

	basicData.Name = app.Name
	basicData.AppGUID = app.GUID
	basicData.Stack = app.Lifecycle.Data.Stack
	basicData.Buildpacks = app.Lifecycle.Data.Buildpacks
	basicData.State = app.State
	basicData.SpaceGUID = app.RelationShips.Space.Data.GUID

	return basicData
}
