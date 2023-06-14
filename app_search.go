package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/cloudfoundry/cli/plugin"
)

// AppSearchResults represents top level attributes of JSON response from Cloud Foundry API
type AppSearchResults struct {
	TotalResults int                 `json:"total_results"`
	TotalPages   int                 `json:"total_pages"`
	NextUrl      string              `json:"next_url"`
	Resources    []AppSearchResource `json:"resources"`
}

// AppSearchResource represents resources attribute of JSON response from Cloud Foundry API
type AppSearchResource struct {
	Metadata AppMetadata `json:"metadata"`
	Entity   AppEntity   `json:"entity"`
}

type AppMetadata struct {
	AppGUID string `json:"guid"`
}

type AppEntity struct {
	Name                string                  `json:"name"`
	Instances           int                     `json:"instances"`
	State               string                  `json:"state"`
	Memory              int                     `json:"memory"`
	DiskQuota           int                     `json:"disk_quota"`
	Buildpack           string                  `json:"buildpack"`
	DetectedBuildPack   string                  `json:"detected_buildpack"`
	SpaceGUID           string                  `json:"space_guid"`
	StartCommand        string                  `json:"detected_start_command"`
	Environment         map[string]string       `json:"environment_json"`
	Command             string                  `json:"command"`
	HealthCheck         string                  `json:"health_check_type"`
	HealthCheckEndpoint string                  `json:"health_check_http_endpoint"`
	Routes              []string                `json:"routes"`
	RoutesUrl           string                  `json:"routes_url"`
	Stack               string                  `json:"stack"`
	StackUrl            string                  `json:"stack_url"`
	ServiceInstances    []ServiceInstanceEntity `json:"service_instances"`
	ServiceUrl          string                  `json:"service_bindings_url"`
}

type Services struct {
	Resources []ServiceResource `json:"resources"`
}

type ServiceResource struct {
	Entity ServiceEntity `json:"entity"`
}

type ServiceEntity struct {
	ServiceInstanceUrl string `json:"service_instance_url"`
}

type Routes struct {
	Resources []RouteResource `json:"resources"`
}

type RouteResource struct {
	Entity RouteEntity `json:"entity"`
}

type RouteEntity struct {
	Host      string `json:"host"`
	DomainUrl string `json:"domain_url"`
}

type Entity struct {
	Entity EntityEntity `json:"entity"`
}

type EntityEntity struct {
	Name string `json:"name"`
}

type ServiceInstance struct {
	Entity ServiceInstanceEntity `json:"entity"`
}

type ServiceInstanceEntity struct {
	Name                   string            `json:"name"`
	Type                   string            `json:"type"`
	MaintenanceInfo        MaintenanceInfo   `json:"maintenance_info"`
	ServicePlanUrl         string            `json:"service_plan_url"`
	ServiceInstanceKeysUrl string            `json:"service_keys_url"`
	ServicePlanEntity      ServicePlanEntity `json:"service_plan_entity"`
}

type MaintenanceInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

type AppPackages struct {
	Resources []AppPackageResource `json:"resources"`
}

type AppPackageResource struct {
	GUID string `json:"guid"`
}

type ServicePlanEntity struct {
	ServicePlanEntityData ServicePlanEntityData `json:"entity"`
}

type ServicePlanEntityData struct {
	Name        string `json:"name"`
	Free        bool   `json:"free"`
	Description string `json:"description"`
	Public      string `json:"public"`
	Active      bool   `json:"active"`
}

// GetAppData requests all of the Application data from Cloud Foundry
func getAppData(cli plugin.CliConnection) AppSearchResults {
	fmt.Println("**** Gathering application metadata from all orgs and spaces ****")

	res := unmarshallAppSearchResults("/v2/apps", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v2/apps?order-direction=asc&page=%d&results-per-page=50", i)

			tRes := unmarshallAppSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallAppSearchResults(apiUrl string, cli plugin.CliConnection) AppSearchResults {
	var tRes AppSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)

	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}

func getRoutes(app *AppSearchResource, cli plugin.CliConnection) {
	var routeURLs []string
	var routes Routes
	cmd := []string{"curl", app.Entity.RoutesUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &routes)

	for _, route := range routes.Resources {
		var domain Entity
		cmd := []string{"curl", route.Entity.DomainUrl}
		output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
		json.Unmarshal([]byte(strings.Join(output, "")), &domain)

		var routeURL = route.Entity.Host + "." + domain.Entity.Name

		routeURLs = append(routeURLs, routeURL)
	}

	app.Entity.Routes = routeURLs
}

func getStacks(app *AppSearchResource, cli plugin.CliConnection) {
	var stack Entity
	cmd := []string{"curl", app.Entity.StackUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &stack)

	app.Entity.Stack = stack.Entity.Name
}

func getServices(app *AppSearchResource, cli plugin.CliConnection) {
	var services Services
	var serviceInstances []ServiceInstanceEntity

	var serviceInstance ServiceInstance
	cmd := []string{"curl", app.Entity.ServiceUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &services)

	for _, service := range services.Resources {
		cmd := []string{"curl", service.Entity.ServiceInstanceUrl}
		output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
		json.Unmarshal([]byte(strings.Join(output, "")), &serviceInstance)
		serviceInstance.Entity.ServicePlanEntity = getServicePlanDetails(serviceInstance.Entity, cli)

		serviceInstances = append(serviceInstances, serviceInstance.Entity)
	}

	app.Entity.ServiceInstances = serviceInstances
}

func getServicePlanDetails(serviceInstanceEntity ServiceInstanceEntity, cli plugin.CliConnection) ServicePlanEntity {
	var servicePlanEntity ServicePlanEntity
	cmd := []string{"curl", serviceInstanceEntity.ServicePlanUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &servicePlanEntity)
	return servicePlanEntity
}

func gatherData(cli plugin.CliConnection) (map[string]string, map[string]SpaceSearchResource, AppSearchResults) {
	orgs := getOrgs(cli)
	spaces := getSpaces(cli)
	apps := getAppData(cli)

	for i, app := range apps.Resources {
		getRoutes(&app, cli)
		getStacks(&app, cli)
		getServices(&app, cli)
		apps.Resources[i] = app
	}

	return orgs, spaces, apps
}

func generateAppManifests(currentDir string, cli plugin.CliConnection) {
	orgs, spaces, apps := gatherData(cli)

	var wg sync.WaitGroup
	for _, app := range apps.Resources {
		wg.Add(1)
		go func(orgs map[string]string, spaces map[string]SpaceSearchResource, app AppSearchResource, currentDir string) {
			defer wg.Done()
			createAppManifest(orgs, spaces, app, currentDir)
		}(orgs, spaces, app, currentDir)
	}

	wg.Wait()
}

func createAppManifest(orgs map[string]string, spaces map[string]SpaceSearchResource, app AppSearchResource, currentDir string) {
	space := spaces[app.Entity.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

	yamlData, err := yaml.Marshal(app)
	if err != nil {
		fmt.Printf("Failed to marshal YAML: %s\n", err)
		return
	}

	fileName := app.Entity.Name + ".yml"

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	filePath := filepath.Join(orgDir, fileName)

	if err := ioutil.WriteFile(filePath, []byte(yamlData), 0644); err != nil {
		fmt.Printf("Failed to write file '%s': %s\n", fileName, err)
		return
	}
	fmt.Printf("File '%s' created successfully.\n", fileName)
}

func downloadApplicationPackages(currentDir string, cli plugin.CliConnection) {
	orgs, spaces, apps := gatherData(cli)

	for _, app := range apps.Resources {
		downloadAppPackages(orgs, spaces, app, currentDir, cli)
	}
}

func downloadAppPackages(orgs map[string]string, spaces map[string]SpaceSearchResource, app AppSearchResource, currentDir string, cli plugin.CliConnection) {
	space := spaces[app.Entity.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

	cmd := []string{"curl", "/v3/apps/" + app.Metadata.AppGUID + "/packages"}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)

	var appPackages AppPackages
	json.Unmarshal([]byte(strings.Join(output, "")), &appPackages)

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	for _, appPackageResource := range appPackages.Resources {
		downloadAppPackage(app, appPackageResource, orgDir, cli)
	}
}

func downloadAppPackage(app AppSearchResource, appPackageResource AppPackageResource, orgDir string, cli plugin.CliConnection) {
	fileName := app.Entity.Name + "-" + appPackageResource.GUID + ".zip"

	filePath := filepath.Join(orgDir, fileName)

	cmd := []string{"curl", "/v3/packages/" + appPackageResource.GUID + "/download", "--output", filePath}
	cli.CliCommandWithoutTerminalOutput(cmd...)

	fmt.Printf("Package download successfully in '%s'\n", fileName)
}
