package cmd

import (
	"encoding/json"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
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

func getAppProcesses(app AppResource, cli plugin.CliConnection, displayAppChan chan<- DisplayApp) {
	var displayApp DisplayApp

	var appProcesses AppProcesses

	processUrl, _ := url.Parse(app.AppLinks.Processes.Href)

	cmd := []string{"curl", processUrl.Path}

	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &appProcesses)

	for _, appProcess := range appProcesses.Processes {
		if appProcess.AppProcessRelationship.AppRelationShip.Data.GUID == app.GUID {
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

	displayAppChan <- displayApp
}
