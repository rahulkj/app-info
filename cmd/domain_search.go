package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
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

func getDomains(cli plugin.CliConnection) map[string]string {
	var data map[string]string
	data = make(map[string]string)
	domains := getDomainsData(cli)

	for _, val := range domains.Resources {
		data[val.GUID] = val.Name
	}

	return data
}

// GetOrgData requests all of the Application data from Cloud Foundry
func getDomainsData(cli plugin.CliConnection) Domains {
	var res Domains = unmarshallDomainsSearchResults("/v3/domains", cli)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v3/domains?page=%d&per_page=50", i)
			tRes := unmarshallDomainsSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallDomainsSearchResults(apiUrl string, cli plugin.CliConnection) Domains {
	var tRes Domains
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}
