package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

// SpaceSearchResults represents top level attributes of JSON response from Cloud Foundry API
type SpaceSearchResults struct {
	TotalResults int                   `json:"total_results"`
	TotalPages   int                   `json:"total_pages"`
	Resources    []SpaceSearchResource `json:"resources"`
}

// SpaceSearchResource represents resources attribute of JSON response from Cloud Foundry API
type SpaceSearchResource struct {
	Name          string            `json:"name"`
	SpaceGUID     string            `json:"guid"`
	Relationships SpaceRelationship `json:"relationships"`
}

type SpaceRelationship struct {
	RelationshipsOrg RelationshipsOrg `json:"organization"`
}

type RelationshipsOrg struct {
	OrgData OrgData `json:"data"`
}

type OrgData struct {
	OrgGUID string `json:"guid"`
}

// GetSpaceData requests all of the Application data from Cloud Foundry
func getSpaces(cli plugin.CliConnection) map[string]SpaceSearchResource {
	var data map[string]SpaceSearchResource = make(map[string]SpaceSearchResource)
	spaces := getSpaceData(cli)

	for _, val := range spaces.Resources {
		data[val.SpaceGUID] = val
	}

	return data
}

// GetSpaceData requests all of the Application data from Cloud Foundry
func getSpaceData(cli plugin.CliConnection) SpaceSearchResults {
	var res SpaceSearchResults = unmarshallSpaceSearchResults("/v3/spaces", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v3/spaces?page=%d&per_page=50", i)
			tRes := unmarshallSpaceSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallSpaceSearchResults(apiUrl string, cli plugin.CliConnection) SpaceSearchResults {
	var tRes SpaceSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}
