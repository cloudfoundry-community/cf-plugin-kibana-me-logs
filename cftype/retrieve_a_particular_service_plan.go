package cftype

// TODO: generate this code from http://apidocs.cloudfoundry.org/204/apps/retrieve_a_particular_service_plan.html

// RetrieveAParticularServicePlan Retrieve a Particular Service Plan
// GET /v2/service_plans/:guid
// http://apidocs.cloudfoundry.org/204/apps/retrieve_a_particular_service_plan.html
type RetrieveAParticularServicePlan struct {
	Metadata RetrieveAParticularServicePlanMetadata
	Entity   RetrieveAParticularServicePlanEntity
}

// RetrieveAParticularServicePlanMetadata ...
type RetrieveAParticularServicePlanMetadata struct {
	GUID string
	URL  string
}

// RetrieveAParticularServicePlanEntity ...
type RetrieveAParticularServicePlanEntity struct {
	Name        string `json:"name"`
	ServiceGUID string `json:"service_guid"`
	ServiceURL  string `json:"service_url"`
}
