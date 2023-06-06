package main

import (
	"encoding/json"
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin"
)

// AppInfo represents Buildpack Usage CLI interface
type AppInfo struct{}

// GetMetadata provides the Cloud Foundry CLI with metadata to provide user about how to use buildpack-usage command
func (c *AppInfo) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "app-info",
		Version: plugin.VersionType{
			Major: 2,
			Minor: 1,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     "app-info",
				HelpText: "Command to view all apps running across all orgs/spaces in the cf deployment",
				UsageDetails: plugin.Usage{
					Usage: "cf app-info [flags]",
					Options: map[string]string{
						"--csv or -c":       "Minimal application details",
						"--json or -j":      "All application details in json format",
						"--manifests or -m": "Generate application mainfests in current working directory",
						"--packages or -p":  "Download applications packages in current working directory",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(AppInfo))
}

// Run is what is executed by the Cloud Foundry CLI when the buildpack-usage command is specified
func (c AppInfo) Run(cli plugin.CliConnection, args []string) {
	if args[0] == "app-info" {
		if len(args) < 2 {
			fmt.Printf("Missing flags, please run help to see the valid options")
			os.Exit(0)
		}

		if args[1] == "--json" || args[1] == "-j" {
			c.printVerboseOutputInJsonFormat(cli)
		} else if args[1] == "--manifests" || args[1] == "-m" {
			c.downloadApplicationManifests(cli)
		} else if args[1] == "--csv" || args[1] == "-c" {
			c.printInCSVFormat(cli)
		} else if args[1] == "--packages" || args[1] == "-p" {
			c.downloadApplicationPackages(cli)
		} else {
			fmt.Printf("Invalid flags, please run help to see the valid options")
		}
	}
}

// PrintInCSVFormat prints the app and buildpack used info on the console
func (c AppInfo) printInCSVFormat(cli plugin.CliConnection) {
	orgs, spaces, apps := c.GatherMinimalData(cli)

	fmt.Println("**** Following is the csv output ****")
	fmt.Println()

	fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n", "ORG", "SPACE", "APPLICATION", "STATE", "INSTANCES", "MEMORY", "DISK", "BUILDPACK", "DETECTED_BUILDPACK", "HEALTH_CHECK")
	for _, val := range apps.Resources {

		space := spaces[val.Entity.SpaceGUID]
		spaceName := space.Name
		orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

		fmt.Printf("%s,%s,%s,%s,%v,%v MB,%v MB,%s,%s,%s\n", orgName, spaceName, val.Entity.Name, val.Entity.State, val.Entity.Instances, val.Entity.Memory, val.Entity.DiskQuota, val.Entity.Buildpack, val.Entity.DetectedBuildPack, val.Entity.HealthCheck)
	}
}

// PrintVerboseOutputInJsonFormat prints the app state, instances, memroy and disk data to console
func (c AppInfo) printVerboseOutputInJsonFormat(cli plugin.CliConnection) {
	_, _, apps := c.GatherData(cli)

	b, err := json.Marshal(apps)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("**** Following is the json output ****")
	fmt.Println(string(b))
}

func (c AppInfo) downloadApplicationManifests(cli plugin.CliConnection) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to access current directory: %s\n", err)
		return
	}

	currentDir = currentDir + "/output"

	fmt.Println("Output will be generated in: ", currentDir)

	os.MkdirAll(currentDir, os.ModePerm)

	c.GenerateAppManifests(currentDir, cli)

	fmt.Println("Generate application manifests are located in: ", currentDir)
}

func (c AppInfo) downloadApplicationPackages(cli plugin.CliConnection) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to access current directory: %s\n", err)
		return
	}

	currentDir = currentDir + "/output"

	fmt.Println("Packages will be downloaded into: ", currentDir)

	os.MkdirAll(currentDir, os.ModePerm)

	c.DownloadApplicationPackages(currentDir, cli)

	fmt.Println("Application packages are located in: ", currentDir)
}
