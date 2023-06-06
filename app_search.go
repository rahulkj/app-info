package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/cloudfoundry/cli/plugin"
)

// AppSearchResults represents top level attributes of JSON response from Cloud Foundry API
type AppSearchResults struct {
	TotalResults int                  `json:"total_results"`
	TotalPages   int                  `json:"total_pages"`
	Resources    []AppSearchResources `json:"resources"`
}

// AppSearchResources represents resources attribute of JSON response from Cloud Foundry API
type AppSearchResources struct {
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
	ServiceInstances    []ServiceInstanceEntity `json:"service_instances`
	ServiceUrl          string                  `json:"service_bindings_url"`
}

type Services struct {
	Resources []ServiceResources `json:"resources"`
}

type ServiceResources struct {
	Entity ServiceEntity `json:"entity"`
}

type ServiceEntity struct {
	ServiceInstanceUrl string `json:"service_instance_url"`
}

type Routes struct {
	Resources []RouteResources `json:"resources"`
}

type RouteResources struct {
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
	Name            string          `json:"name"`
	Type            string          `json:"type"`
	MaintenanceInfo MaintenanceInfo `json:"maintenance_info"`
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

// GetAppData requests all of the Application data from Cloud Foundry
func (c AppInfo) GetAppData(cli plugin.CliConnection) AppSearchResults {
	fmt.Println("**** Gathering application metadata from all orgs and spaces ****")

	res := c.UnmarshallAppSearchResults("/v2/apps", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v2/apps", strconv.Itoa(i))
			tRes := c.UnmarshallAppSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	for i, app := range res.Resources {
		c.getRoutes(&app, cli)

		c.getStacks(&app, cli)

		c.getServices(&app, cli)

		res.Resources[i] = app
	}

	return res
}

func (c AppInfo) UnmarshallAppSearchResults(apiUrl string, cli plugin.CliConnection) AppSearchResults {
	var tRes AppSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}

func (c AppInfo) getRoutes(app *AppSearchResources, cli plugin.CliConnection) {
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

func (c AppInfo) getStacks(app *AppSearchResources, cli plugin.CliConnection) {
	var stack Entity
	cmd := []string{"curl", app.Entity.StackUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &stack)

	app.Entity.Stack = stack.Entity.Name
}

func (c AppInfo) getServices(app *AppSearchResources, cli plugin.CliConnection) {
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
		serviceInstances = append(serviceInstances, serviceInstance.Entity)
	}

	app.Entity.ServiceInstances = serviceInstances
}

func (c AppInfo) GatherData(cli plugin.CliConnection) (map[string]string, map[string]SpaceSearchResources, AppSearchResults) {
	orgs := c.GetOrgs(cli)
	spaces := c.GetSpaces(cli)
	apps := c.GetAppData(cli)

	return orgs, spaces, apps
}

func (c AppInfo) GenerateAppManifests(currentDir string, cli plugin.CliConnection) {
	orgs, spaces, apps := c.GatherData(cli)

	for _, app := range apps.Resources {

		space := spaces[app.Entity.SpaceGUID]
		spaceName := space.Name
		orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

		yamlData, err := yaml.Marshal(app)
		if err != nil {
			fmt.Println("Failed to marshal YAML: %s", err)
			return
		}

		fileName := app.Entity.Name + ".yml"

		orgDir := currentDir + "/" + orgName + "/" + spaceName

		os.MkdirAll(orgDir, os.ModePerm)

		filePath := filepath.Join(orgDir, fileName)

		if err := ioutil.WriteFile(filePath, []byte(yamlData), 0644); err != nil {
			fmt.Println("Failed to write file '%s': %s", fileName, err)
			return
		}
		fmt.Println("File '%s' created successfully.", fileName)
	}
}

func (c AppInfo) DownloadApplicationPackages(currentDir string, cli plugin.CliConnection) {
	orgs, spaces, apps := c.GatherData(cli)

	for _, app := range apps.Resources {

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
			fileName := app.Entity.Name + "-" + appPackageResource.GUID + ".zip"

			filePath := filepath.Join(orgDir, fileName)

			cmd := []string{"curl", "/v3/packages/" + appPackageResource.GUID + "/download", "--output", filePath}
			cli.CliCommandWithoutTerminalOutput(cmd...)

			fmt.Println("Package download successfully in '%s'", fileName)
		}
	}
}
