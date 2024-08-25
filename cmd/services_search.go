package cmd

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

func getAllServices(config Config) {

}

func getAppServices(app DisplayApp, services []Service) {

}
