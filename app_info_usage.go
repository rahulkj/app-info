package main

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
)

// BuildpackUsage represents Buildpack Usage CLI interface
type BuildpackUsage struct{}

// Metadata is the data retrived from the response json
type Metadata struct {
	GUID string `json:"guid"`
}

// GetMetadata provides the Cloud Foundry CLI with metadata to provide user about how to use buildpack-usage command
func (c *BuildpackUsage) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "app-info",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 1,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "app-info",
				HelpText: "Command to view all apps running across all orgs/spaces in the cf deployment with specific details",
				UsageDetails: plugin.Usage{
					Usage: "cf app-info\n   cf app-info --verbose",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(BuildpackUsage))
}

// Run is what is executed by the Cloud Foundry CLI when the buildpack-usage command is specified
func (c BuildpackUsage) Run(cli plugin.CliConnection, args []string) {
	if args[0] == "app-info" {
		orgs := c.GetOrgs(cli)
		spaces := c.GetSpaces(cli)
		apps := c.GetAppData(cli)
		if len(args) == 2 {
			if args[1] == "--verbose" {
				c.PrintVerboseOutputInCSVFormat(orgs, spaces, apps)
			}
		} else {
			c.PrintInCSVFormat(orgs, spaces, apps)
		}

	}
}

// PrintInCSVFormat prints the app and buildpack used info on the console
func (c BuildpackUsage) PrintInCSVFormat(orgs map[string]string, spaces map[string]SpaceSearchEntity, apps AppSearchResults) {
	fmt.Println("")

	fmt.Printf("Following is the csv output \n\n")

	fmt.Printf("%s,%s,%s,%s,%s\n", "ORG", "SPACE", "APPLICATION", "STATE", "BUILDPACK")

	for _, val := range apps.Resources {
		bp := val.Entity.Buildpack
		if bp == "" {
			bp = val.Entity.DetectedBuildpack
		}

		space := spaces[val.Entity.SpaceGUID]
		spaceName := space.Name
		orgName := orgs[space.OrgGUID]

		fmt.Printf("%s,%s,%s,%s,%s\n", orgName, spaceName, val.Entity.Name, val.Entity.State, bp)
	}
}

// PrintVerboseOutputInCSVFormat prints the app state, instances, memroy and disk data to console
func (c BuildpackUsage) PrintVerboseOutputInCSVFormat(orgs map[string]string, spaces map[string]SpaceSearchEntity, apps AppSearchResults) {
	fmt.Println("")

	fmt.Printf("Following is the csv output \n\n")

	fmt.Printf("%s,%s,%s,%s,%s,%s,%s\n", "ORG", "SPACE", "APPLICATION", "STATE", "INSTANCES", "MEMORY", "DISK")

	for _, val := range apps.Resources {

		space := spaces[val.Entity.SpaceGUID]
		spaceName := space.Name
		orgName := orgs[space.OrgGUID]

		fmt.Printf("%s,%s,%s,%s,%v,%v MB,%v MB\n", orgName, spaceName, val.Entity.Name, val.Entity.State, val.Entity.Instances, val.Entity.Memory, val.Entity.DiskQuota)
	}
}
