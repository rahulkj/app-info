package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

// OrgSearchResults represents top level attributes of JSON response from Cloud Foundry API
type OrgSearchResults struct {
	TotalResults int                 `json:"total_results"`
	TotalPages   int                 `json:"total_pages"`
	Resources    []OrgSearchResource `json:"resources"`
}

// OrgSearchResource represents resources attribute of JSON response from Cloud Foundry API
type OrgSearchResource struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

func getOrgs(cli plugin.CliConnection) map[string]string {
	var data map[string]string
	data = make(map[string]string)
	orgs := getOrgData(cli)

	for _, val := range orgs.Resources {
		data[val.GUID] = val.Name
	}

	return data
}

// GetOrgData requests all of the Application data from Cloud Foundry
func getOrgData(cli plugin.CliConnection) OrgSearchResults {
	var res OrgSearchResults = unmarshallOrgSearchResults("/v3/organizations", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v3/organizations?page=%d&per_page=50", i)
			tRes := unmarshallOrgSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallOrgSearchResults(apiUrl string, cli plugin.CliConnection) OrgSearchResults {
	var tRes OrgSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}
