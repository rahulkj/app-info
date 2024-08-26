package cmd

type Config struct {
	ApiEndpoint string `yaml:"cf_endpoint"`
	OauthToken  string `yaml:"token"`
}

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
	AppRelationShip DataHolder `json:"app"`
}

type DataHolder struct {
	Data Data `json:"data"`
}

type Data struct {
	GUID string `json:"guid"`
}

type AppEnvironment struct {
	AppGUID     string
	Environment map[string]interface{} `json:"var"`
}

type AppFeatures struct {
	AppGUID  string
	Features []AppFeatureResource `json:"resources"`
}

type AppFeatureResource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// AppSearchResults represents top level attributes of JSON response from Cloud Foundry API
type Apps struct {
	Resources  []AppResource `json:"resources"`
	Pagination Pagination    `json:"pagination"`
}

type AppResource struct {
	GUID          string           `json:"guid"`
	Name          string           `json:"name"`
	State         string           `json:"state"`
	Lifecycle     Lifecycle        `json:"lifecycle"`
	RelationShips AppRelationShips `json:"relationships"`
	AppLinks      Links            `json:"links"`
}

type Lifecycle struct {
	Type string        `json:"type"`
	Data LifecycleData `json:"data"`
}

type LifecycleData struct {
	Buildpacks []string `json:"buildpacks"`
	Stack      string   `json:"stack"`
}

type AppRelationShips struct {
	Space DataHolder `json:"space"`
}

type Links struct {
	Self                      Link `json:"self"`
	Details                   Link `json:"details"`
	EnvironmentVars           Link `json:"environment_variables"`
	Space                     Link `json:"space"`
	Processes                 Link `json:"processes"`
	Packages                  Link `json:"packages"`
	CurrentDroplet            Link `json:"current_droplet"`
	Droplets                  Link `json:"droplets"`
	Tasks                     Link `json:"tasks"`
	Revisions                 Link `json:"revisions"`
	DeployedRevisions         Link `json:"deployed_revisions"`
	Features                  Link `json:"features"`
	ServiceCredentialBindings Link `json:"service_credential_bindings"`
	ServiceRouteBindings      Link `json:"service_route_bindings"`
	ServicePlan               Link `json:"service_plan"`
}

type Link struct {
	Href string `json:"href"`
}
type DisplayApp struct {
	Name                       string                 `json:"name"`
	AppGUID                    string                 `json:"guid"`
	Instances                  int                    `json:"instances"`
	State                      string                 `json:"state"`
	Memory                     int                    `json:"memory_in_mb"`
	Disk                       int                    `json:"disk_in_mb"`
	LogRate                    int                    `json:"log_rate_limit_in_bytes_per_second"`
	Buildpacks                 []string               `json:"buildpacks"`
	DetectedBuildPack          string                 `json:"detected_buildpack"`
	DetectedBuildPackFileNames []string               `json:"detected_buildpack_filenames"`
	SpaceGUID                  string                 `json:"space_guid"`
	Environment                map[string]interface{} `json:"environment_json"`
	HealthCheck                string                 `json:"health_check_type"`
	ReadinessHealthCheck       string                 `json:"readiness_health_check_type"`
	Type                       string                 `json:"type"`
	Routes                     []string               `json:"routes"`
	Stack                      string                 `json:"stack"`
	Services                   []Service              `json:"services"`
	Features                   []AppFeatureResource   `json:"resources"`
	StackGUID                  string                 `json:"stackguid"`
}

type Buildpacks struct {
	Resources  []BuildpackResources `json:"resources"`
	Pagination Pagination           `json:"pagination"`
}

type BuildpackResources struct {
	GUID     string `json:"guid"`
	Name     string `json:"name"`
	Stack    string `json:"stack"`
	State    string `json:"state"`
	Position int    `json:"position"`
	Filename string `json:"filename"`
	Enabled  bool   `json:"enabled"`
	Locked   bool   `json:"locked"`
}

type AppDetectedBuildpacks struct {
	AppGUID                    string
	DetectedBuildPackFileNames []string
}

type AppPackages struct {
	Resources []Data `json:"resources"`
}

type CurrentDroplet struct {
	GUID  string             `json:"guid"`
	Links CurrentDropletLink `json:"links"`
}

type CurrentDropletLink struct {
	Package Link `json:"package"`
}

// OrgSearchResults represents top level attributes of JSON response from Cloud Foundry API
type OrgSearchResults struct {
	Pagination Pagination          `json:"pagination"`
	Resources  []OrgSearchResource `json:"resources"`
}

// OrgSearchResource represents resources attribute of JSON response from Cloud Foundry API
type OrgSearchResource struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

type Routes struct {
	Resources  []RouteResources `json:"resources"`
	Pagination Pagination       `json:"pagination"`
}

type RouteResources struct {
	GUID         string        `json:"guid"`
	Protocol     string        `json:"protocol"`
	Host         string        `json:"host"`
	Path         string        `json:"path"`
	URL          string        `json:"url"`
	Destinations []Destination `json:"destinations"`
}

type Destination struct {
	GUID           string         `json:"guid"`
	Port           int            `json:"port"`
	DestinationApp DestinationApp `json:"app"`
}

type DestinationApp struct {
	GUID    string                `json:"guid"`
	Process DestinationAppProcess `json:"process"`
}

type DestinationAppProcess struct {
	Type string `json:"type"`
}

type AppRoutes struct {
	AppGUID string
	Routes  []string
}

// SpaceSearchResults represents top level attributes of JSON response from Cloud Foundry API
type SpaceSearchResults struct {
	Pagination Pagination            `json:"pagination"`
	Resources  []SpaceSearchResource `json:"resources"`
}

// SpaceSearchResource represents resources attribute of JSON response from Cloud Foundry API
type SpaceSearchResource struct {
	Name          string            `json:"name"`
	SpaceGUID     string            `json:"guid"`
	Relationships SpaceRelationship `json:"relationships"`
}

type SpaceRelationship struct {
	RelationshipsOrg DataHolder `json:"organization"`
}

type Stacks struct {
	Resources  []StackResource `json:"resources"`
	Pagination Pagination      `json:"pagination"`
}

type StackResource struct {
	GUID        string `json:"guid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

type Pagination struct {
	TotalResults int `json:"total_results"`
	TotalPages   int `json:"total_pages"`
}

type ServiceInstances struct {
	Pagination Pagination                 `json:"pagination"`
	Resources  []ServiceInstancesResource `json:"resources"`
}

type ServiceInstancesResource struct {
	GUID            string                               `json:"guid"`
	CreatedAt       string                               `json:"created_at"`
	UpdatedAt       string                               `json:"updated_at"`
	Name            string                               `json:"name"`
	Type            string                               `json:"type"`
	MaintenanceInfo MaintenanceInfo                      `json:"maintenance_info"`
	Relationships   ServiceInstancesResourceRelationship `json:"relationships"`
	Links           Links                                `json:"links"`
}

type ServiceInstancesResourceRelationship struct {
	Space       DataHolder `json:"space"`
	ServicePlan DataHolder `json:"service_plan"`
}

type MaintenanceInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

type ServiceInstanceBindings struct {
	Pagination Pagination                       `json:"pagination"`
	Resources  []ServiceInstanceBindingResource `json:"resources"`
}

type ServiceInstanceBindingResource struct {
	ServiceInstanceGUID string
	GUID                string                                     `json:"guid"`
	Type                string                                     `json:"type"`
	Relationships       ServiceInstanceBindingResourceRelationship `json:"relationships"`
	Links               Links                                      `json:"links"`
}

type ServiceInstanceBindingResourceRelationship struct {
	App             DataHolder `json:"app"`
	ServiceInstance DataHolder `json:"service_instance"`
}

type Service struct {
	GUID        string `json:"guid"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type Info struct {
	Name        string `json:"name"`
	Build       string `json:"build"`
	Description string `json:"description"`
}
