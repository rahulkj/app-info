package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rahulkj/app-info/cmd"
	"gopkg.in/yaml.v3"
)

// AppInfo represents Buildpack Usage CLI interface
type AppInfo struct{}

func main() {
	option := flag.String("option", "csv", "csv, json, yaml, packages")
	configFileLocation := flag.String("config", "", "Absolute path to config file that has the cloud foundry target and bearer token")
	includeEnvironmentVars := flag.Bool("include-env", false, "Optional flag to include environment variables in json / manifest output. (default false)")
	flag.Parse()

	if *option == "" {
		fmt.Println("Error: -option cannot be empty.")
		flag.Usage()
		os.Exit(1)
	}

	if *configFileLocation == "" {
		fmt.Println("Error: -config should be specified")
		flag.Usage()
		os.Exit(1)
	}

	var config cmd.Config

	if *configFileLocation != "" {
		c, err := checkConfigExists(*configFileLocation)

		if c == nil || err != nil {
			fmt.Println("Error: Specfied config does not have the required keys")
			flag.Usage()
			os.Exit(1)
		} else {
			config = *c
		}
	}

	switch *option {
	case "csv":
		printInCSVFormat(config, *includeEnvironmentVars)
	case "json":
		printVerboseOutputInJsonFormat(config, *includeEnvironmentVars)
	case "yaml":
		downloadApplicationManifests(config, *includeEnvironmentVars)
	case "packages":
		downloadApplicationPackages(config)
	default:
		fmt.Println("Error: -option cannot be empty.")
		flag.Usage()
		os.Exit(1)
	}

}

func checkConfigExists(filePath string) (*cmd.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config cmd.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.ApiEndpoint == "" || config.OauthToken == "" {
		return nil, nil
	}

	return &config, nil
}

// PrintInCSVFormat prints the app and buildpack used info on the console
func printInCSVFormat(config cmd.Config, include_env_variables bool) {
	orgs, spaces, apps := cmd.GatherData(config, include_env_variables)

	fmt.Println("**** Following is the csv output ****")
	fmt.Println()

	fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n", "ORG", "SPACE", "APPLICATION", "STATE", "INSTANCES", "MEMORY", "DISK", "HEALTH_CHECK", "STACK", "BUILDPACK", "DETECTED_BUILDPACK", "DETECTED_BUILDPACK_FILENAME")
	for _, val := range apps {

		space := spaces[val.SpaceGUID]
		spaceName := space.Name
		orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

		fmt.Printf("%s,%s,%s,%s,%v,%v MB,%v MB,%s,%s,%s,%s,%s\n", orgName, spaceName, val.Name, val.State, val.Instances, val.Memory, val.Disk, val.HealthCheck, val.Stack, val.Buildpacks, val.DetectedBuildPack, val.DetectedBuildPackFileNames)
	}
}

// PrintVerboseOutputInJsonFormat prints the app state, instances, memroy and disk data to console
func printVerboseOutputInJsonFormat(config cmd.Config, include_env_variables bool) {
	_, _, apps := cmd.GatherData(config, include_env_variables)

	b, err := json.Marshal(apps)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("**** Following is the json output ****")
	fmt.Println(string(b))
}

func downloadApplicationManifests(config cmd.Config, include_env_variables bool) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to access current directory: %s\n", err)
		return
	}

	currentDir = currentDir + "/output"

	fmt.Println("Output will be generated in: ", currentDir)

	os.MkdirAll(currentDir, os.ModePerm)

	cmd.GenerateAppManifests(currentDir, config, include_env_variables)

	fmt.Println("Generate application manifests are located in: ", currentDir)
}

func downloadApplicationPackages(config cmd.Config) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to access current directory: %s\n", err)
		return
	}

	currentDir = currentDir + "/output"

	fmt.Println("Packages will be downloaded into: ", currentDir)

	os.MkdirAll(currentDir, os.ModePerm)

	cmd.DownloadApplicationPackages(currentDir, config)

	fmt.Println("Application packages are located in: ", currentDir)
}
