package cmd

import (
	"encoding/json"
	"fmt"
)

func getOrgs(config Config) map[string]string {
	data := make(map[string]string)
	orgs := getOrgData(config)

	for _, val := range orgs.Resources {
		data[val.GUID] = val.Name
	}

	return data
}

// GetOrgData requests all of the Application data from Cloud Foundry
func getOrgData(config Config) OrgSearchResults {
	apiUrl := fmt.Sprintf("%s/v3/organizations", config.ApiEndpoint)
	var res OrgSearchResults = unmarshallOrgSearchResults(apiUrl, config)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("%s?page=%d&per_page=100", apiUrl, i)
			tRes := unmarshallOrgSearchResults(apiUrl, config)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallOrgSearchResults(apiUrl string, config Config) OrgSearchResults {
	var tRes OrgSearchResults
	output, _ := getResponse(config, apiUrl)
	json.Unmarshal([]byte(output), &tRes)

	return tRes
}
