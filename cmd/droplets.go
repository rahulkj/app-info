package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

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
	Package PackageLink `json:"package"`
}

type PackageLink struct {
	Href string `json:"href"`
}

func DownloadApplicationPackages(currentDir string, config Config) {
	orgs, spaces, apps := GatherData(config, false)

	var messages []string

	var wg sync.WaitGroup
	result := make(chan string)

	bar := progressbar.Default(int64(len(apps)), "Downloading App Packages")

	for _, app := range apps {
		wg.Add(1)
		go downloadAppPackages(orgs, spaces, app, currentDir, config, result, &wg)
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

func downloadAppPackages(orgs map[string]string, spaces map[string]SpaceSearchResource, app DisplayApp, currentDir string, config Config, result chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	space := spaces[app.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

	apiUrl := fmt.Sprintf("%s/v3/apps/%s/droplets/current", config.ApiEndpoint, app.AppGUID)

	output, _ := getResponse(config, apiUrl)

	var currentDroplet CurrentDroplet
	json.Unmarshal([]byte(output), &currentDroplet)

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	downloadAppPackage(app, currentDroplet, orgDir, config, result)
}

func downloadAppPackage(app DisplayApp, currentDroplet CurrentDroplet, orgDir string, config Config, result chan string) {
	var message string

	if currentDroplet.Links.Package.Href != "" {
		apiUrl := currentDroplet.Links.Package.Href + "/download"

		fileName := app.Name + "-" + app.AppGUID + ".zip"

		filePath := filepath.Join(orgDir, fileName)

		ok, err := downloadFile(config, apiUrl, filePath)

		if ok && err == nil {
			message = fmt.Sprintf("Package download successfully for %s in: %s", app.Name, filePath)
		} else {
			message = fmt.Sprintf("Package download unsuccessfully for: %s", app.Name)
		}
	} else {
		message = fmt.Sprintf("Package does not exist for: %s", app.Name)
	}

	result <- message
}
