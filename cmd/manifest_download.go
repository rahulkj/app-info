package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

func GenerateManifests(currentDir string, config Config, include_env_variables bool) {
	orgs, spaces, apps, unboundServices := GatherData(config, include_env_variables)

	var wg sync.WaitGroup
	for _, app := range apps {
		wg.Add(1)
		go func(orgs map[string]string, spaces map[string]SpaceSearchResource, app DisplayApp, currentDir string) {
			defer wg.Done()
			createAppManifest(orgs, spaces, app, currentDir)
		}(orgs, spaces, app, currentDir)
	}

	for _, service := range unboundServices {
		wg.Add(1)

		go func(orgs map[string]string, spaces map[string]SpaceSearchResource, service Service, currentDir string) {
			defer wg.Done()
			createUnboundServiceManifest(orgs, spaces, service, currentDir)
		}(orgs, spaces, service, currentDir)
	}

	wg.Wait()
}

func createAppManifest(orgs map[string]string, spaces map[string]SpaceSearchResource, app DisplayApp, currentDir string) {

	space := spaces[app.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.Data.GUID]

	yamlData, err := yaml.Marshal(app)
	if err != nil {
		log.Printf("Failed to marshal YAML: %s\n", err)
		return
	}

	fileName := app.Name + "-application.yml"

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	filePath := filepath.Join(orgDir, fileName)

	if err := os.WriteFile(filePath, []byte(yamlData), 0644); err != nil {
		log.Printf("Failed to write file '%s': %s\n", fileName, err)
		return
	}

	fmt.Printf("File '%s' created successfully.\n", fileName)
}

func createUnboundServiceManifest(orgs map[string]string, spaces map[string]SpaceSearchResource, service Service, currentDir string) {

	space := spaces[service.SpaceGUID]
	spaceName := space.Name
	orgName := orgs[space.Relationships.RelationshipsOrg.Data.GUID]

	yamlData, err := yaml.Marshal(service)
	if err != nil {
		log.Printf("Failed to marshal YAML: %s\n", err)
		return
	}

	fileName := service.Name + "-unbound-service.yml"

	orgDir := currentDir + "/" + orgName + "/" + spaceName

	os.MkdirAll(orgDir, os.ModePerm)

	filePath := filepath.Join(orgDir, fileName)

	if err := os.WriteFile(filePath, []byte(yamlData), 0644); err != nil {
		log.Printf("Failed to write file '%s': %s\n", fileName, err)
		return
	}

	fmt.Printf("File '%s' created successfully.\n", fileName)
}
