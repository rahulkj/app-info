package cmd

import (
	"encoding/json"
	"fmt"
)

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

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
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
