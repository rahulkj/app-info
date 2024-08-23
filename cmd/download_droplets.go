package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

type AppPackages struct {
	Resources []AppPackageResource `json:"resources"`
}

type AppPackageResource struct {
	GUID string `json:"guid"`
}

func DownloadApplicationPackages(currentDir string, cli plugin.CliConnection) {
	orgs, spaces, apps := GatherData(cli, false)

	for _, app := range apps {
		downloadAppPackages(orgs, spaces, app, currentDir, cli)
	}
}

func downloadAppPackages(orgs map[string]string, spaces map[string]SpaceSearchResource, app DisplayApp, currentDir string, cli plugin.CliConnection) {
	space := spaces[app.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

	cmd := []string{"curl", "/v3/apps/" + app.AppGUID + "/packages"}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)

	var appPackages AppPackages
	json.Unmarshal([]byte(strings.Join(output, "")), &appPackages)

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	for _, appPackageResource := range appPackages.Resources {
		downloadAppPackage(app, appPackageResource, orgDir, cli)
	}
}

func downloadAppPackage(app DisplayApp, appPackageResource AppPackageResource, orgDir string, cli plugin.CliConnection) {
	fileName := app.Name + "-" + appPackageResource.GUID + ".zip"

	filePath := filepath.Join(orgDir, fileName)

	cmd := []string{"curl", "/v3/packages/" + appPackageResource.GUID + "/download", "--output", filePath}
	cli.CliCommandWithoutTerminalOutput(cmd...)

	fmt.Printf("Package download successfully in '%s'\n", fileName)
}
