package cmd

import (
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type Routes struct {
	Resources []RouteResource `json:"resources"`
}

type RouteResource struct {
	Entity RouteResourceEntity `json:"entity"`
}

type RouteResourceEntity struct {
	Host       string `json:"host"`
	DomainGuid string `json:"domain_guid"`
}

type Route struct {
	Entity RouteEntity `json:"entity"`
}

type RouteEntity struct {
	Name string `json:"name"`
}

func getRoutes(app *AppSearchResource, domains map[string]string, cli plugin.CliConnection) {
	var routeURLs []string
	var routes Routes
	cmd := []string{"curl", app.Entity.RoutesUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &routes)

	for _, route := range routes.Resources {
		domain := domains[route.Entity.DomainGuid]
		var routeURL = route.Entity.Host + "." + domain
		routeURLs = append(routeURLs, routeURL)
	}

	app.Entity.Routes = routeURLs
}
