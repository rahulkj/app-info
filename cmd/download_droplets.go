package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/schollz/progressbar/v3"
)

type AppPackages struct {
	Resources []AppPackageResource `json:"resources"`
}

type AppPackageResource struct {
	GUID string `json:"guid"`
}

type CurrentDroplet struct {
	GUID  string             `json:"guid"`
	Links CurrentDropletLink `json:"links"`
}

type CurrentDropletLink struct {
	App     AppLink     `json:"app"`
	Package PackageLink `json:"package"`
}

type AppLink struct {
	Href string `json:"href"`
}

type PackageLink struct {
	Href string `json:"href"`
}

func DownloadApplicationPackages(currentDir string, cli plugin.CliConnection) {
	orgs, spaces, apps := GatherData(cli, false)

	var messages []string

	var wg sync.WaitGroup
	result := make(chan string)

	bar := progressbar.Default(int64(len(apps)), "Downloading App Packages")

	for _, app := range apps {
		wg.Add(1)
		go downloadAppPackages(orgs, spaces, app, currentDir, cli, result, &wg)
	}

	for i := 0; i < len(apps); i++ {
		message := <-result
		bar.Add(1)
		messages = append(messages, message)
	}

	wg.Wait()

	for _, message := range messages {
		fmt.Println(message)
	}
}

func downloadAppPackages(orgs map[string]string, spaces map[string]SpaceSearchResource, app DisplayApp, currentDir string, cli plugin.CliConnection, result chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	space := spaces[app.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

	cmd := []string{"curl", "/v3/apps/" + app.AppGUID + "/droplets/current"}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)

	var currentDroplet CurrentDroplet
	json.Unmarshal([]byte(strings.Join(output, "")), &currentDroplet)

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	downloadAppPackage(app, currentDroplet, orgDir, cli, result)
}

func downloadAppPackage(app DisplayApp, currentDroplet CurrentDroplet, orgDir string, cli plugin.CliConnection, result chan string) {

	currentPackageUrl, _ := url.Parse(currentDroplet.Links.Package.Href)

	fileName := app.Name + "-" + app.AppGUID + ".zip"

	filePath := filepath.Join(orgDir, fileName)

	cmd := []string{"curl", currentPackageUrl.Path + "/download", "--output", filePath}
	cli.CliCommandWithoutTerminalOutput(cmd...)

	result <- "Package download successfully in: " + fileName
}
