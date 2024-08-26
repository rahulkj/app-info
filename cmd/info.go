package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetInfo(config Config) Info {
	apiUrl := fmt.Sprintf("%s/v3/info", config.ApiEndpoint)
	var res Info = unmarshallInfo(apiUrl, config)
	return res
}

func unmarshallInfo(apiUrl string, config Config) Info {
	var tRes Info
	output, err := getResponse(config, apiUrl)

	if output != "" {
		json.Unmarshal([]byte(output), &tRes)
	} else {
		Red("Failed to connect to the provided endpoint: %s due to %s\n", config.ApiEndpoint, err)
		os.Exit(1)
	}

	return tRes
}
