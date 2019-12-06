package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

// OrgSearchResults represents top level attributes of JSON response from Cloud Foundry API
type OrgSearchResults struct {
	TotalResults int                  `json:"total_results"`
	TotalPages   int                  `json:"total_pages"`
	Resources    []OrgSearchResources `json:"resources"`
}

// OrgSearchResources represents resources attribute of JSON response from Cloud Foundry API
type OrgSearchResources struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

func (c AppInfo) GetOrgs(cli plugin.CliConnection) map[string]string {
	var data map[string]string
	data = make(map[string]string)
	orgs := c.GetOrgData(cli)

	for _, val := range orgs.Resources {
		data[val.GUID] = val.Name
	}

	return data
}

// GetOrgData requests all of the Application data from Cloud Foundry
func (c AppInfo) GetOrgData(cli plugin.CliConnection) OrgSearchResults {
	var res OrgSearchResults
	res = c.UnmarshallOrgSearchResults("/v3/organizations", cli)

	if res.TotalPages > 0 {
		for i := 0; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v3/organizations", strconv.Itoa(i))
			tRes := c.UnmarshallOrgSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func (c AppInfo) UnmarshallOrgSearchResults(apiUrl string, cli plugin.CliConnection) OrgSearchResults {
	var tRes OrgSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}
