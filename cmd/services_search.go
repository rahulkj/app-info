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
	Name                       string                `json:"name"`
	Type                       string                `json:"type"`
	MaintenanceInfo            MaintenanceInfo       `json:"maintenance_info"`
	ServicePlanUrl             string                `json:"service_plan_url"`
	ServiceInstanceKeysUrl     string                `json:"service_keys_url"`
	ServiceInstancePlanDetails ServicePlanEntityData `json:"service_plan_details"`
}

type MaintenanceInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

type ServicePlanEntity struct {
	ServicePlanEntityData ServicePlanEntityData `json:"entity"`
}

type ServicePlanEntityData struct {
	Name              string `json:"name"`
	Free              bool   `json:"free"`
	Description       string `json:"description"`
	Active            bool   `json:"active"`
	Bindable          bool   `json:"bindable"`
	ServiceURL        string `json:"service_url"`
	Label             string `json:"label"`
	ServiceBrokerName string `json:"service_broker_name"`
}

type ServiceData struct {
	Service ServiceDataEntity `json:"entity"`
}

type ServiceDataEntity struct {
	Label             string `json:"label"`
	Description       string `json:"description"`
	ServiceBrokerName string `json:"service_broker_name"`
}

func getServices(app *AppSearchResource, cli plugin.CliConnection) {
	var services Services
	var serviceInstances []ServiceInstanceEntity

	cmd := []string{"curl", app.Entity.ServiceBindingsUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &services)

	for _, service := range services.Resources {
		var serviceInstance ServiceInstance
		cmd := []string{"curl", service.Entity.ServiceInstanceUrl}
		output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
		json.Unmarshal([]byte(strings.Join(output, "")), &serviceInstance)
		serviceInstance.Entity.ServiceInstancePlanDetails = getServicePlanDetails(serviceInstance.Entity, cli).ServicePlanEntityData

		serviceInstances = append(serviceInstances, serviceInstance.Entity)
	}

	app.Entity.ServiceInstances = serviceInstances
}

func getServicePlanDetails(serviceInstanceEntity ServiceInstanceEntity, cli plugin.CliConnection) ServicePlanEntity {
	var servicePlanEntity ServicePlanEntity
	cmd := []string{"curl", serviceInstanceEntity.ServicePlanUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &servicePlanEntity)

	serviceData := getServiceDetails(servicePlanEntity, cli)

	servicePlanEntity.ServicePlanEntityData.Label = serviceData.Service.Label
	servicePlanEntity.ServicePlanEntityData.ServiceBrokerName = serviceData.Service.ServiceBrokerName

	return servicePlanEntity
}

func getServiceDetails(servicePlanEntity ServicePlanEntity, cli plugin.CliConnection) ServiceData {
	var serviceData ServiceData
	cmd := []string{"curl", servicePlanEntity.ServicePlanEntityData.ServiceURL}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &serviceData)
	return serviceData
}
