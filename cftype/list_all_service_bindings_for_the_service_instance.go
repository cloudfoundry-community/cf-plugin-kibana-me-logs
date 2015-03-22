package cftype

// TODO: generate this code from http://apidocs.cloudfoundry.org/204/apps/list_all_service_bindings_for_the_service_instance.html

// ListAllServiceBindingsForTheServiceInstance List all Service Bindings for the Service Instance
// GET /v2/service_instances/:guid/service_bindings
// http://apidocs.cloudfoundry.org/204/apps/list_all_service_bindings_for_the_service_instance.html
type ListAllServiceBindingsForTheServiceInstance struct {
	TotalResults int    `json:"total_results"`
	TotalPages   int    `json:"total_pages"`
	PrevURL      string `json:"prev_url"`
	NextURL      string `json:"next_url"`
	Resources    []ListAllServiceBindingsForTheServiceInstanceResource
}

// ListAllServiceBindingsForTheServiceInstanceResource ...
type ListAllServiceBindingsForTheServiceInstanceResource struct {
	Entity ListAllServiceBindingsForTheServiceInstanceResourceEntity
}

// ListAllServiceBindingsForTheServiceInstanceResourceEntity ...
type ListAllServiceBindingsForTheServiceInstanceResourceEntity struct {
	AppGUID             string                 `json:"app_guid"`
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	Credentials         map[string]interface{} `json:"credentials"`
	BindingOptions      map[string]interface{} `json:"binding_options"`
	SyslogDrainURL      string                 `json:"syslog_drain_url"`
	AppURL              string                 `json:"app_url"`
	ServiceInstanceURL  string                 `json:"service_instance_url"`
}
