package cmd

import (
	"encoding/json"
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
	AppProcessRelationship AppProcessRelationship `json:"relationships"`
}

type AppHealthCheck struct {
	Type string `json:"type"`
}

type AppProcessRelationship struct {
	AppRelationShip AppRelationShip `json:"app"`
}

type AppRelationShip struct {
	AppRelationShipData AppRelationShipData `json:"data"`
}

type AppRelationShipData struct {
	AppRelationShipGUID string `json:"guid"`
}

func getAppProcesses(app DisplayApp, cli plugin.CliConnection) (displayApp DisplayApp) {
	var appProcesses AppProcesses
	cmd := []string{"curl", "/v3/apps/" + app.AppGUID + "/processes"}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &appProcesses)

	for _, appProcess := range appProcesses.Processes {
		if appProcess.GUID == app.AppGUID {
			app.Instances = appProcess.Instances
			app.Memory = appProcess.Memory
			app.Disk = appProcess.Disk
			app.LogRate = appProcess.LogRate
			app.HealthCheck = appProcess.HealthCheck.Type
			break
		}
	}

	return app
}
