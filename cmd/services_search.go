package cmd

import (
	"encoding/json"
	"fmt"
	"sync"
)

// GetSpaceData requests all of the Application data from Cloud Foundry
func getAllServices(config Config) map[string]ServiceInstancesResource {
	apiUrl := fmt.Sprintf("%s/v3/service_instances", config.ApiEndpoint)
	var res ServiceInstances = unmarshallServiceInstancesResults(apiUrl, config)

	if res.Pagination.TotalPages > 1 {
		for i := 2; i <= res.Pagination.TotalPages; i++ {
			apiUrl := fmt.Sprintf("%s?page=%d&per_page=100", apiUrl, i)
			tRes := unmarshallServiceInstancesResults(apiUrl, config)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	serviceInstanceResources := make(map[string]ServiceInstancesResource)

	for _, serviceInstanceResource := range res.Resources {
		serviceInstanceResources[serviceInstanceResource.GUID] = serviceInstanceResource
	}

	return serviceInstanceResources
}

func unmarshallServiceInstancesResults(apiUrl string, config Config) ServiceInstances {
	var tRes ServiceInstances
	output, _ := getResponse(config, apiUrl)
	json.Unmarshal([]byte(output), &tRes)

	return tRes
}

func getAllAppServiceBindings(config Config, serviceInstanceResource ServiceInstancesResource) ServiceInstanceBindings {

	var res ServiceInstanceBindings

	var result ServiceInstanceBindings

	apiUrl := serviceInstanceResource.Links.ServiceCredentialBindings.Href

	if apiUrl != "" {
		result = unmarshallServiceInstanceBindingResults(apiUrl, config)

		if res.Pagination.TotalPages > 1 {
			for i := 2; i <= result.Pagination.TotalPages; i++ {
				apiUrl := fmt.Sprintf("%s?page=%d&per_page=100", apiUrl, i)
				tRes := unmarshallServiceInstanceBindingResults(apiUrl, config)
				result.Resources = append(result.Resources, tRes.Resources...)
			}
		}
	}

	if len(result.Resources) == 0 {
		var sib ServiceInstanceBindingResource
		result.Resources = append(result.Resources, sib)
	}

	for _, sib := range result.Resources {
		sib.ServiceInstanceGUID = serviceInstanceResource.GUID
		res.Resources = append(res.Resources, sib)
	}

	return res

}

func unmarshallServiceInstanceBindingResults(apiUrl string, config Config) ServiceInstanceBindings {
	var tRes ServiceInstanceBindings
	output, _ := getResponse(config, apiUrl)
	json.Unmarshal([]byte(output), &tRes)

	return tRes
}

func getServices(config Config) (map[string][]ServiceInstancesResource, []ServiceInstancesResource) {
	appServicesBinding := make(map[string][]ServiceInstancesResource)
	var unboundServices []ServiceInstancesResource

	services := getAllServices(config)

	var wg sync.WaitGroup
	appDataCh := make(chan ServiceInstanceBindings)

	for _, serviceInstanceResource := range services {
		wg.Add(1)

		go func() {
			defer wg.Done()
			res := getAllAppServiceBindings(config, serviceInstanceResource)

			appDataCh <- res
		}()
	}

	go func() {
		wg.Wait()
		close(appDataCh) // Close results channel after all workers are done
	}()

	for data := range appDataCh {
		for _, serviceCredentialBindingResource := range data.Resources {
			appGuid := serviceCredentialBindingResource.Relationships.App.Data.GUID
			serviceInstanceGuid := serviceCredentialBindingResource.Relationships.ServiceInstance.Data.GUID

			if appGuid != "" {
				serviceInstanceResources := appServicesBinding[appGuid]
				serviceInstanceResources = append(serviceInstanceResources, services[serviceInstanceGuid])
				appServicesBinding[appGuid] = serviceInstanceResources
			} else {
				unboundServices = append(unboundServices, services[serviceCredentialBindingResource.ServiceInstanceGUID])
			}
		}
	}

	return appServicesBinding, unboundServices
}

func getAppSevices(app AppResource, services map[string][]ServiceInstancesResource) DisplayApp {
	var displayApp DisplayApp

	appServices := services[app.GUID]

	for _, service := range appServices {
		var appService Service

		appService.Name = service.Name
		appService.GUID = service.GUID
		appService.Type = service.Type

		appService.Description = service.MaintenanceInfo.Description
		appService.Version = service.MaintenanceInfo.Version

		displayApp.Services = append(displayApp.Services, appService)
	}

	return displayApp
}
