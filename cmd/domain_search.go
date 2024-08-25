package cmd

import (
	"encoding/json"
	"fmt"
)

type Domains struct {
	Resources  []DomainResources `json:"resources"`
	Pagination DomainsPagination `json:"pagination"`
}

type DomainsPagination struct {
	TotalPages int `json:"total_pages"`
}

type DomainResources struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

func getDomains(config Config) map[string]string {
	var data map[string]string
	data = make(map[string]string)
	domains := getDomainsData(config)

	for _, val := range domains.Resources {
		data[val.GUID] = val.Name
	}

	return data
}

// GetOrgData requests all of the Application data from Cloud Foundry
func getDomainsData(config Config) Domains {
	apiUrl := fmt.Sprintf("%s/v3/domains", config.ApiEndpoint)
	var res Domains = unmarshallDomainsSearchResults(apiUrl, config)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("%s?page=%d&per_page=100", apiUrl, i)
			tRes := unmarshallDomainsSearchResults(apiUrl, config)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallDomainsSearchResults(apiUrl string, config Config) Domains {
	var tRes Domains
	output, _ := getResponse(config, apiUrl)
	json.Unmarshal([]byte(output), &tRes)

	return tRes
}
