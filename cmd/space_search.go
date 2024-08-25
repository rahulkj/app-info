package cmd

import (
	"encoding/json"
	"fmt"
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
func getSpaces(config Config) map[string]SpaceSearchResource {
	var data map[string]SpaceSearchResource = make(map[string]SpaceSearchResource)
	spaces := getSpaceData(config)

	for _, val := range spaces.Resources {
		data[val.SpaceGUID] = val
	}

	return data
}

// GetSpaceData requests all of the Application data from Cloud Foundry
func getSpaceData(config Config) SpaceSearchResults {
	apiUrl := fmt.Sprintf("%s/v3/spaces", config.ApiEndpoint)
	var res SpaceSearchResults = unmarshallSpaceSearchResults(apiUrl, config)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("%s?page=%d&per_page=100", apiUrl, i)
			tRes := unmarshallSpaceSearchResults(apiUrl, config)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallSpaceSearchResults(apiUrl string, config Config) SpaceSearchResults {
	var tRes SpaceSearchResults
	output, _ := getResponse(config, apiUrl)
	json.Unmarshal([]byte(output), &tRes)

	return tRes
}
