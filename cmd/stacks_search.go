package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type Stacks struct {
	Resources  []StackResource  `json:"resources"`
	Pagination StacksPagination `json:"pagination"`
}

type StacksPagination struct {
	TotalPages int `json:"total_pages"`
}

type StackResource struct {
	GUID        string `json:"guid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

func getAppStack(app *AppSearchResource, stacks map[string]StackResource) {
	var stackResource StackResource = stacks[app.Entity.StackGUID]
	app.Entity.Stack = stackResource.Name
}

func getStacks(cli plugin.CliConnection) map[string]StackResource {
	var data map[string]StackResource
	data = make(map[string]StackResource)
	Stacks := getStacksData(cli)

	for _, val := range Stacks.Resources {
		data[val.GUID] = val
	}

	return data
}

// GetOrgData requests all of the Application data from Cloud Foundry
func getStacksData(cli plugin.CliConnection) Stacks {
	var res Stacks = unmarshallStacksearchResults("/v3/stacks", cli)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v3/stacks?page=%d&per_page=50", i)
			tRes := unmarshallStacksearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallStacksearchResults(apiUrl string, cli plugin.CliConnection) Stacks {
	var tRes Stacks
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}
