package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

// SpaceSearchResults represents top level attributes of JSON response from Cloud Foundry API
type SpaceSearchResults struct {
	TotalResults int                    `json:"total_results"`
	TotalPages   int                    `json:"total_pages"`
	Resources    []SpaceSearchResources `json:"resources"`
}

// SpaceSearchResources represents resources attribute of JSON response from Cloud Foundry API
type SpaceSearchResources struct {
	Name          string             `json:"name"`
	SpaceGUID     string             `json:"guid"`
	Relationships SpaceRelationships `json:"relationships"`
}

type SpaceRelationships struct {
	RelationshipsOrg RelationshipsOrg `json:"organization"`
}

type RelationshipsOrg struct {
	OrgData OrgData `json:"data"`
}

type OrgData struct {
	OrgGUID string `json:"guid"`
}

// GetSpaceData requests all of the Application data from Cloud Foundry
func (c AppInfo) GetSpaces(cli plugin.CliConnection) map[string]SpaceSearchResources {
	var data map[string]SpaceSearchResources = make(map[string]SpaceSearchResources)
	spaces := c.GetSpaceData(cli)

	for _, val := range spaces.Resources {
		data[val.SpaceGUID] = val
	}

	return data
}

// GetSpaceData requests all of the Application data from Cloud Foundry
func (c AppInfo) GetSpaceData(cli plugin.CliConnection) SpaceSearchResults {
	var res SpaceSearchResults = c.UnmarshallSpaceSearchResults("/v3/spaces", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v3/spaces?page=%d&per_page=50", i)
			tRes := c.UnmarshallSpaceSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func (c AppInfo) UnmarshallSpaceSearchResults(apiUrl string, cli plugin.CliConnection) SpaceSearchResults {
	var tRes SpaceSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}
