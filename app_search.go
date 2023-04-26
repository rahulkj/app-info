package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
	Name                string            `json:"name"`
	Instances           int               `json:"instances"`
	State               string            `json:"state"`
	Memory              int               `json:"memory"`
	DiskQuota           int               `json:"disk_quota"`
	Buildpack           string            `json:"buildpack"`
	DetectedBuildPack   string            `json:"detected_buildpack"`
	SpaceGUID           string            `json:"space_guid"`
	StartCommand        string            `json:"detected_start_command"`
	Environment         map[string]string `json:"environment_json"`
	Command             string            `json:"command"`
	HealthCheck         string            `json:"health_check_type"`
	HealthCheckEndpoint string            `json:"health_check_http_endpoint"`
	Routes              []string
	RoutesUrl           string `json:"routes_url"`
	Stack               string
	StackUrl            string `json:"stack_url"`
	ServiceInstances    []ServiceInstanceEntity
	ServiceUrl          string `json:"service_bindings_url"`
}

type AppServices struct {
	ServiceName    string
	ServiceVersion string
	ServiceType    string
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

// GetAppData requests all of the Application data from Cloud Foundry
func (c AppInfo) GetAppData(cli plugin.CliConnection) AppSearchResults {
	var res AppSearchResults
	res = c.UnmarshallAppSearchResults("/v2/apps", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v2/apps", strconv.Itoa(i))
			tRes := c.UnmarshallAppSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	for _, app := range res.Resources {
		routes := c.getRoutes(app, cli)
		app.Entity.Routes = routes

		stack := c.getStacks(app, cli)
		app.Entity.Stack = stack

		services := c.getServices(app, cli)
		app.Entity.ServiceInstances = services
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

func (c AppInfo) getRoutes(app AppSearchResources, cli plugin.CliConnection) []string {
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

	return routeURLs
}

func (c AppInfo) getStacks(app AppSearchResources, cli plugin.CliConnection) string {
	var stack Entity
	cmd := []string{"curl", app.Entity.StackUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &stack)
	return stack.Entity.Name
}

func (c AppInfo) getServices(app AppSearchResources, cli plugin.CliConnection) []ServiceInstanceEntity {
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

	return serviceInstances
}
