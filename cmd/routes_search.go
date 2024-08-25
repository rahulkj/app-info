package cmd

import (
	"encoding/json"
	"fmt"
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

type AppRoutes struct {
	AppGUID string
	Routes  []string
}

// GetOrgData requests all of the Application data from Cloud Foundry
func getAllRoutes(config Config) Routes {
	apiUrl := fmt.Sprintf("%s/v3/routes", config.ApiEndpoint)
	var res Routes = unmarshallRoutesSearchResults(apiUrl, config)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("%s?page=%d&per_page=100", apiUrl, i)
			tRes := unmarshallRoutesSearchResults(apiUrl, config)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallRoutesSearchResults(apiUrl string, config Config) Routes {
	var tRes Routes
	output, _ := getResponse(config, apiUrl)
	json.Unmarshal([]byte(output), &tRes)

	return tRes
}

func getAppRoutes(app AppResource, routes Routes) AppRoutes {
	// defer wg.Done()
	var appRoutes AppRoutes

	appRoutes.AppGUID = app.GUID

	for _, resource := range routes.Resources {
		for _, destination := range resource.Destinations {
			if destination.DestinationApp.GUID == app.GUID {
				appRoutes.Routes = append(appRoutes.Routes, resource.URL)
			}
		}
	}

	return appRoutes
}
