package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rahulkj/app-info/cmd"
	"gopkg.in/yaml.v3"
)

func main() {
	cmd.Logger()

	startTime := time.Now()
	option := flag.String("option", "csv", "csv, json, yaml, packages")
	configFileLocation := flag.String("config", "", "Absolute path to config file that has the cloud foundry target and bearer token")
	includeEnvironmentVars := flag.Bool("include-env", false, "Optional flag to include environment variables in json / manifest output. (default false)")
	flag.Parse()

	if *option == "" {
		cmd.Red("Error: -option cannot be empty.\n")
		flag.Usage()
		os.Exit(1)
	}

	if *configFileLocation == "" {
		cmd.Red("Error: -config should be specified\n")
		flag.Usage()
		os.Exit(1)
	}

	var config cmd.Config

	if *configFileLocation != "" {
		c, err := checkConfigExists(*configFileLocation)

		if c == nil || err != nil {
			cmd.Red("Error: Specfied config does not have the required keys\n")
			flag.Usage()
			os.Exit(1)
		} else {
			config = *c
		}
	}

	info := cmd.GetInfo(config)
	fmt.Printf("Connecting to TAS version: %s\n", info.Build)

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
		cmd.Red("Error: -option is invalid.\n")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println()
	cmd.Yellow("***** Finished in %s *****\n", time.Since(startTime))
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
func printInCSVFormat(config cmd.Config, includeEnvVariables bool) {
	orgs, spaces, apps, unboundServices := cmd.GatherData(config, includeEnvVariables)

	cmd.Green("**** Following is the csv output for apps ****\n")
	fmt.Println()

	fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n", "ORG", "SPACE", "APPLICATION", "STATE", "INSTANCES", "MEMORY", "DISK", "HEALTH_CHECK", "STACK", "BUILDPACK", "DETECTED_BUILDPACK_FILENAME", "SERVICES")
	for _, val := range apps {

		space := spaces[val.SpaceGUID]
		spaceName := space.Name
		orgName := orgs[space.Relationships.RelationshipsOrg.Data.GUID]

		var service_names []string
		for _, service := range val.Services {
			service_names = append(service_names, service.Name)
		}

		fmt.Printf("%s,%s,%s,%s,%v,%v MB,%v MB,%s,%s,%s,%s,%s\n", orgName, spaceName, val.Name, val.State, val.Instances, val.Memory, val.Disk, val.HealthCheck, val.Stack, val.Buildpacks, val.DetectedBuildPackFileNames, service_names)
	}

	if len(unboundServices) != 0 {
		cmd.Green("**** Following is the csv output for unused services ****\n")
		fmt.Printf("%s,%s,%s,%s,%s,%s,%s\n", "SERVICENAME", "CREATED_AT", "UPDATED_AT", "TYPE", "VERSION", "DESCRIPTION", "SPACE")

		for _, service := range unboundServices {
			fmt.Printf("%s,%s,%s,%s,%s,%s,%s\n", service.Name, service.CreatedAt, service.UpdatedAt, service.Type, service.MaintenanceInfo.Version, service.MaintenanceInfo.Description, spaces[service.Relationships.Space.Data.GUID].Name)
		}
	}
}

// PrintVerboseOutputInJsonFormat prints the app state, instances, memroy and disk data to console
func printVerboseOutputInJsonFormat(config cmd.Config, includeEnvVariables bool) {
	_, _, apps, _ := cmd.GatherData(config, includeEnvVariables)

	b, err := json.Marshal(apps)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd.Green("**** Following is the json output ****\n")
	fmt.Println(string(b))
}

func downloadApplicationManifests(config cmd.Config, includeEnvVariables bool) {
	currentDir, err := os.Getwd()
	if err != nil {
		cmd.Red("Failed to access current directory: %s\n", err)
		return
	}

	currentDir = currentDir + "/output"

	cmd.Yellow("Output will be generated in: %s\n", currentDir)

	os.MkdirAll(currentDir, os.ModePerm)

	cmd.GenerateAppManifests(currentDir, config, includeEnvVariables)

	cmd.Green("Generate application manifests are located in: %s\n", currentDir)
}

func downloadApplicationPackages(config cmd.Config) {
	currentDir, err := os.Getwd()
	if err != nil {
		cmd.Red("Failed to access current directory: %s\n", err)
		return
	}

	currentDir = currentDir + "/output"

	cmd.Yellow("Packages will be downloaded into: %s\n", currentDir)

	os.MkdirAll(currentDir, os.ModePerm)

	cmd.DownloadApplicationPackages(currentDir, config)

	cmd.Green("Application packages are located in: %s\n", currentDir)
}
