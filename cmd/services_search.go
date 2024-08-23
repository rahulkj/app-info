package cmd

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
)

type Service struct {
	GUID             string        `json:"guid"`
	Name             string        `json:"name"`
	ServicePlan      ServicePlan   `json:"service_plan"`
	Service          ServiceFields `json:"service_fields"`
	ApplicationNames []string      `json:"application_names"`
	IsUserProvided   bool          `json:"isUserProvided"`
}

type ServicePlan struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

type ServiceFields struct {
	Name string `json:"name"`
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

func getAllServices(cli plugin.CliConnection) []Service {
	results, _ := cli.GetServices()

	var services []Service

	for _, result := range results {
		var service Service
		service.ApplicationNames = result.ApplicationNames
		service.GUID = result.Guid
		service.Name = result.Name
		service.ServicePlan.GUID = result.ServicePlan.Guid
		service.ServicePlan.Name = result.ServicePlan.Name
		service.IsUserProvided = result.IsUserProvided
		service.Service.Name = result.Service.Name

		services = append(services, service)
	}

	return services
}

func getAppServices(app DisplayApp, services []Service) (displayApp DisplayApp) {
	for _, service := range services {
		for _, serviceApp := range service.ApplicationNames {
			if serviceApp == app.Name {
				fmt.Println("MATCHED " + app.Name)
				app.Services = append(app.Services, service)
				break
			}
		}
	}

	return app
}
