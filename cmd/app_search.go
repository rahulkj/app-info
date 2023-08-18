package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/cloudfoundry/cli/plugin"
)

// AppSearchResults represents top level attributes of JSON response from Cloud Foundry API
type AppSearchResults struct {
	TotalResults int                 `json:"total_results"`
	TotalPages   int                 `json:"total_pages"`
	NextUrl      string              `json:"next_url"`
	Resources    []AppSearchResource `json:"resources"`
}

// AppSearchResource represents resources attribute of JSON response from Cloud Foundry API
type AppSearchResource struct {
	Metadata AppMetadata `json:"metadata"`
	Entity   AppEntity   `json:"entity"`
}

type AppMetadata struct {
	AppGUID string `json:"guid"`
}

type AppEntity struct {
	Name                      string                  `json:"name"`
	Instances                 int                     `json:"instances"`
	State                     string                  `json:"state"`
	Memory                    int                     `json:"memory"`
	DiskQuota                 int                     `json:"disk_quota"`
	Buildpack                 string                  `json:"buildpack"`
	DetectedBuildPack         string                  `json:"detected_buildpack"`
	DetectedBuildPackGUID     string                  `json:"detected_buildpack_guid"`
	DetectedBuildPackFileName string                  `json:"detected_buildpack_filename"`
	SpaceGUID                 string                  `json:"space_guid"`
	StartCommand              string                  `json:"detected_start_command"`
	Environment               map[string]string       `json:"environment_json"`
	Command                   string                  `json:"command"`
	HealthCheck               string                  `json:"health_check_type"`
	HealthCheckEndpoint       string                  `json:"health_check_http_endpoint"`
	Routes                    []string                `json:"routes"`
	RoutesUrl                 string                  `json:"routes_url"`
	Stack                     string                  `json:"stack"`
	StackUrl                  string                  `json:"stack_url"`
	ServiceInstances          []ServiceInstanceEntity `json:"service_instances"`
	ServiceUrl                string                  `json:"service_bindings_url"`
}

type AppPackages struct {
	Resources []AppPackageResource `json:"resources"`
}

type AppPackageResource struct {
	GUID string `json:"guid"`
}

// GetAppData requests all of the Application data from Cloud Foundry
func getAppData(cli plugin.CliConnection) AppSearchResults {
	fmt.Println("**** Gathering application metadata from all orgs and spaces ****")

	res := unmarshallAppSearchResults("/v2/apps", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v2/apps?order-direction=asc&page=%d&results-per-page=50", i)

			tRes := unmarshallAppSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func unmarshallAppSearchResults(apiUrl string, cli plugin.CliConnection) AppSearchResults {
	var tRes AppSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)

	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}

func getBuildpackDetails(app *AppSearchResource, buildpacks map[string]BuildpackResources) {
	buildpack := buildpacks[app.Entity.DetectedBuildPackGUID]

	app.Entity.DetectedBuildPackFileName = buildpack.Filename
}

func gatherData(cli plugin.CliConnection) (map[string]string, map[string]SpaceSearchResource, AppSearchResults) {
	orgs := getOrgs(cli)
	spaces := getSpaces(cli)
	apps := getAppData(cli)
	buildpacks := getBuildpacks(cli)

	for i, app := range apps.Resources {
		getRoutes(&app, cli)
		getStacks(&app, cli)
		getServices(&app, cli)
		getBuildpackDetails(&app, buildpacks)
		apps.Resources[i] = app
	}

	return orgs, spaces, apps
}

func generateAppManifests(currentDir string, cli plugin.CliConnection) {
	orgs, spaces, apps := gatherData(cli)

	var wg sync.WaitGroup
	for _, app := range apps.Resources {
		wg.Add(1)
		go func(orgs map[string]string, spaces map[string]SpaceSearchResource, app AppSearchResource, currentDir string) {
			defer wg.Done()
			createAppManifest(orgs, spaces, app, currentDir)
		}(orgs, spaces, app, currentDir)
	}

	wg.Wait()
}

func createAppManifest(orgs map[string]string, spaces map[string]SpaceSearchResource, app AppSearchResource, currentDir string) {
	space := spaces[app.Entity.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

	yamlData, err := yaml.Marshal(app)
	if err != nil {
		fmt.Printf("Failed to marshal YAML: %s\n", err)
		return
	}

	fileName := app.Entity.Name + ".yml"

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	filePath := filepath.Join(orgDir, fileName)

	if err := os.WriteFile(filePath, []byte(yamlData), 0644); err != nil {
		fmt.Printf("Failed to write file '%s': %s\n", fileName, err)
		return
	}
	fmt.Printf("File '%s' created successfully.\n", fileName)
}

func downloadApplicationPackages(currentDir string, cli plugin.CliConnection) {
	orgs, spaces, apps := gatherData(cli)

	for _, app := range apps.Resources {
		downloadAppPackages(orgs, spaces, app, currentDir, cli)
	}
}

func downloadAppPackages(orgs map[string]string, spaces map[string]SpaceSearchResource, app AppSearchResource, currentDir string, cli plugin.CliConnection) {
	space := spaces[app.Entity.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

	cmd := []string{"curl", "/v3/apps/" + app.Metadata.AppGUID + "/packages"}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)

	var appPackages AppPackages
	json.Unmarshal([]byte(strings.Join(output, "")), &appPackages)

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	for _, appPackageResource := range appPackages.Resources {
		downloadAppPackage(app, appPackageResource, orgDir, cli)
	}
}

func downloadAppPackage(app AppSearchResource, appPackageResource AppPackageResource, orgDir string, cli plugin.CliConnection) {
	fileName := app.Entity.Name + "-" + appPackageResource.GUID + ".zip"

	filePath := filepath.Join(orgDir, fileName)

	cmd := []string{"curl", "/v3/packages/" + appPackageResource.GUID + "/download", "--output", filePath}
	cli.CliCommandWithoutTerminalOutput(cmd...)

	fmt.Printf("Package download successfully in '%s'\n", fileName)
}
