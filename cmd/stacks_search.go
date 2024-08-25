package cmd

import (
	"encoding/json"
	"fmt"
)

func getAppStack(app *DisplayApp, stacks map[string]StackResource) {
	var stackResource StackResource = stacks[app.Stack]
	app.StackGUID = stackResource.GUID
}

func getStacks(config Config) map[string]StackResource {
	var data map[string]StackResource
	data = make(map[string]StackResource)
	Stacks := getStacksData(config)

	for _, val := range Stacks.Resources {
		data[val.Name] = val
	}

	return data
}

// GetOrgData requests all of the Application data from Cloud Foundry
func getStacksData(config Config) Stacks {
	apiUrl := fmt.Sprintf("%s/v3/stacks", config.ApiEndpoint)
	var res Stacks = unmarshallStacksearchResults(apiUrl, config)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("%s?page=%d&per_page=100", apiUrl, i)
			tRes := unmarshallStacksearchResults(apiUrl, config)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallStacksearchResults(apiUrl string, config Config) Stacks {
	var tRes Stacks
	output, _ := getResponse(config, apiUrl)
	json.Unmarshal([]byte(output), &tRes)

	return tRes
}
