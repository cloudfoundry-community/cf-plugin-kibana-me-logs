package cftype

// TODO: generate this code from http://apidocs.cloudfoundry.org/204/apps/retrieve_a_particular_service_instance.html

// RetrieveAParticularServiceInstance Retrieve a Particular Service Instance
// GET /v2/service_instances/:guid
// http://apidocs.cloudfoundry.org/204/apps/retrieve_a_particular_service_instance.html
type RetrieveAParticularServiceInstance struct {
	Metadata RetrieveAParticularServiceInstanceMetadata
	Entity   RetrieveAParticularServiceInstanceEntity
}

// RetrieveAParticularServiceInstanceMetadata ...
type RetrieveAParticularServiceInstanceMetadata struct {
	GUID string
	URL  string
}

// RetrieveAParticularServiceInstanceEntity ...
type RetrieveAParticularServiceInstanceEntity struct {
	Name              string                 `json:"name"`
	Credentials       map[string]interface{} `json:"credentials"`
	ServicePlanGUID   string                 `json:"service_plan_guid"`
	SpaceGUID         string                 `json:"space_guid"`
	DashboardURL      string                 `json:"dashboard_url"`
	EntityType        string                 `json:"type"`
	SpaceURL          string                 `json:"space_url"`
	ServicePlanURL    string                 `json:"service_plan_url"`
	ServiceBindingURL string                 `json:"service_binding_url"`
}
