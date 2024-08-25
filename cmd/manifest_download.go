package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

func GenerateAppManifests(currentDir string, config Config, include_env_variables bool) {
	orgs, spaces, apps := GatherData(config, include_env_variables)

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
	logger := log.New(os.Stdout, "Log: ", log.Ldate|log.Ltime)

	space := spaces[app.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.OrgData.OrgGUID]

	yamlData, err := yaml.Marshal(app)
	if err != nil {
		logger.Printf("Failed to marshal YAML: %s\n", err)
		return
	}

	fileName := app.Name + ".yml"

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	filePath := filepath.Join(orgDir, fileName)

	if err := os.WriteFile(filePath, []byte(yamlData), 0644); err != nil {
		logger.Printf("Failed to write file '%s': %s\n", fileName, err)
		return
	}

	fmt.Printf("File '%s' created successfully.\n", fileName)
}
