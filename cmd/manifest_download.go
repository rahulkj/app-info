package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cloudfoundry/cli/plugin"
	"gopkg.in/yaml.v3"
)

func GenerateAppManifests(currentDir string, cli plugin.CliConnection, include_env_variables bool) {
	orgs, spaces, apps := GatherData(cli, include_env_variables)

	var wg sync.WaitGroup
	for _, app := range apps {
		wg.Add(1)
		go func(orgs map[string]string, spaces map[string]SpaceSearchResource, app DisplayApp, currentDir string) {
			defer wg.Done()
			createAppManifest(orgs, spaces, app, currentDir)
		}(orgs, spaces, app, currentDir)
	}

	wg.Wait()
}

func createAppManifest(orgs map[string]string, spaces map[string]SpaceSearchResource, app DisplayApp, currentDir string) {
	space := spaces[app.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

	yamlData, err := yaml.Marshal(app)
	if err != nil {
		fmt.Printf("Failed to marshal YAML: %s\n", err)
		return
	}

	fileName := app.Name + ".yml"

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	filePath := filepath.Join(orgDir, fileName)

	if err := os.WriteFile(filePath, []byte(yamlData), 0644); err != nil {
		fmt.Printf("Failed to write file '%s': %s\n", fileName, err)
		return
	}
	fmt.Printf("File '%s' created successfully.\n", fileName)
}
