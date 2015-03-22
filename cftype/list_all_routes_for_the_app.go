package cftype

// TODO: generate this code from http://apidocs.cloudfoundry.org/204/apps/list_all_routes_for_the_app.html

// ListAllRoutesForTheApp List all Routes for the App
// GET /v2/apps/:guid/routes
// http://apidocs.cloudfoundry.org/204/apps/list_all_routes_for_the_app.html
type ListAllRoutesForTheApp struct {
	TotalResults int    `json:"total_results"`
	TotalPages   int    `json:"total_pages"`
	PrevURL      string `json:"prev_url"`
	NextURL      string `json:"next_url"`
	Resources    []ListAllRoutesForTheAppResource
}

// ListAllRoutesForTheAppResource ...
type ListAllRoutesForTheAppResource struct {
	Entity ListAllRoutesForTheAppResourceEntity
}

// ListAllRoutesForTheAppResourceEntity ...
type ListAllRoutesForTheAppResourceEntity struct {
	Host       string `json:"host"`
	DomainGUID string `json:"domain_guid"`
	SpaceGUID  string `json:"space_guid"`
	DomainURL  string `json:"domain_url"`
	SpaceURL   string `json:"space_url"`
	AppsURL    string `json:"apps_url"`
}
