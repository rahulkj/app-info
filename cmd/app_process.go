package cmd

import (
	"encoding/json"
)

type AppProcesses struct {
	Processes []AppProcessResource `json:"resources"`
}

type AppProcessResource struct {
	GUID                   string                 `json:"guid"`
	Type                   string                 `json:"type"`
	Instances              int                    `json:"instances"`
	Memory                 int                    `json:"memory_in_mb"`
	Disk                   int                    `json:"disk_in_mb"`
	LogRate                int                    `json:"log_rate_limit_in_bytes_per_second"`
	HealthCheck            AppHealthCheck         `json:"health_check"`
	ReadinessHealthCheck   AppHealthCheck         `json:"readiness_health_check"`
	AppProcessRelationship AppProcessRelationship `json:"relationships"`
}

type AppHealthCheck struct {
	Type string `json:"type"`
}

type AppProcessRelationship struct {
	AppRelationShip AppRelationShip `json:"app"`
}

type AppRelationShip struct {
	Data Data `json:"data"`
}

type Data struct {
	GUID string `json:"guid"`
}

func getAppProcesses(app AppResource, config Config) DisplayApp {
	// defer wg.Done()

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
