package cmd

import (
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type Services struct {
	Resources []ServiceResource `json:"resources"`
}

type ServiceResource struct {
	Entity ServiceEntity `json:"entity"`
}

type ServiceEntity struct {
	ServiceInstanceUrl string `json:"service_instance_url"`
}

type ServiceInstance struct {
	Entity ServiceInstanceEntity `json:"entity"`
}

type ServiceInstanceEntity struct {
	Name                   string            `json:"name"`
	Type                   string            `json:"type"`
	MaintenanceInfo        MaintenanceInfo   `json:"maintenance_info"`
	ServicePlanUrl         string            `json:"service_plan_url"`
	ServiceInstanceKeysUrl string            `json:"service_keys_url"`
	ServicePlanEntity      ServicePlanEntity `json:"service_plan_entity"`
}

type MaintenanceInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

type ServicePlanEntity struct {
	ServicePlanEntityData ServicePlanEntityData `json:"entity"`
}

type ServicePlanEntityData struct {
	Name        string `json:"name"`
	Free        bool   `json:"free"`
	Description string `json:"description"`
	Public      string `json:"public"`
	Active      bool   `json:"active"`
}

func getServices(app *AppSearchResource, cli plugin.CliConnection) {
	var services Services
	var serviceInstances []ServiceInstanceEntity

	var serviceInstance ServiceInstance
	cmd := []string{"curl", app.Entity.ServiceUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &services)

	for _, service := range services.Resources {
		cmd := []string{"curl", service.Entity.ServiceInstanceUrl}
		output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
		json.Unmarshal([]byte(strings.Join(output, "")), &serviceInstance)
		serviceInstance.Entity.ServicePlanEntity = getServicePlanDetails(serviceInstance.Entity, cli)

		serviceInstances = append(serviceInstances, serviceInstance.Entity)
	}

	app.Entity.ServiceInstances = serviceInstances
}

func getServicePlanDetails(serviceInstanceEntity ServiceInstanceEntity, cli plugin.CliConnection) ServicePlanEntity {
	var servicePlanEntity ServicePlanEntity
	cmd := []string{"curl", serviceInstanceEntity.ServicePlanUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &servicePlanEntity)
	return servicePlanEntity
}
