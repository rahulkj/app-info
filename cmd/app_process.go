package cmd

import (
	"encoding/json"
)

func getAppProcesses(app AppResource, config Config) DisplayApp {
	var displayApp DisplayApp

	apiUrl := app.AppLinks.Processes.Href

	var appProcesses AppProcesses
	output, _ := getResponse(config, apiUrl)

	json.Unmarshal([]byte(output), &appProcesses)

	for _, appProcess := range appProcesses.Processes {
		if appProcess.AppProcessRelationship.AppRelationShip.Data.GUID == app.GUID {
			displayApp.AppGUID = app.GUID
			displayApp.Instances = appProcess.Instances
			displayApp.Memory = appProcess.Memory
			displayApp.Disk = appProcess.Disk
			displayApp.LogRate = appProcess.LogRate
			displayApp.HealthCheck = appProcess.HealthCheck.Type
			displayApp.ReadinessHealthCheck = appProcess.ReadinessHealthCheck.Type
			displayApp.Type = appProcess.Type
			break
		}
	}

	return displayApp
}
