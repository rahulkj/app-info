package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type Routes struct {
	Resources  []RouteResources `json:"resources"`
	Pagination RoutesPagination `json:"pagination"`
}

type RoutesPagination struct {
	TotalPages int `json:"total_pages"`
}

type RouteResources struct {
	GUID         string        `json:"guid"`
	Protocol     string        `json:"protocol"`
	Host         string        `json:"host"`
	Path         string        `json:"path"`
	URL          string        `json:"url"`
	Destinations []Destination `json:"destinations"`
}

type Destination struct {
	GUID           string         `json:"guid"`
	Port           int            `json:"port"`
	DestinationApp DestinationApp `json:"app"`
}

type DestinationApp struct {
	GUID    string                `json:"guid"`
	Process DestinationAppProcess `json:"process"`
}

type DestinationAppProcess struct {
	Type string `json:"type"`
}

// GetOrgData requests all of the Application data from Cloud Foundry
func getAllRoutes(cli plugin.CliConnection) Routes {
	var res Routes = unmarshallRoutesSearchResults("/v3/routes", cli)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v3/buildpacks?page=%d&per_page=50", i)
			tRes := unmarshallRoutesSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallRoutesSearchResults(apiUrl string, cli plugin.CliConnection) Routes {
	var tRes Routes
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}

func getAppRoutes(app AppResource, routes Routes, displayAppChan chan<- DisplayApp) {
	var displayApp DisplayApp

	var routeURLs []string

	for _, resource := range routes.Resources {
		for _, destination := range resource.Destinations {
			if destination.DestinationApp.GUID == app.GUID {
				routeURLs = append(routeURLs, resource.URL)
			}
		}
	}

	displayApp.Routes = routeURLs

	displayAppChan <- displayApp
}
