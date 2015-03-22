package cftype

// TODO: generate this code from http://apidocs.cloudfoundry.org/204/apps/retrieve_a_particular_app.html

// RetrieveAParticularApp Retrieve a Particular App
// GET /v2/apps/:guid
// http://apidocs.cloudfoundry.org/204/apps/retrieve_a_particular_app.html
type RetrieveAParticularApp struct {
	Metadata RetrieveAParticularAppMetadata
	Entity   RetrieveAParticularAppEntity
}

// RetrieveAParticularAppMetadata ...
type RetrieveAParticularAppMetadata struct {
	GUID string
	URL  string
}

// RetrieveAParticularAppEntity ...
type RetrieveAParticularAppEntity struct {
	Name                 string `json:"name"`
	SpaceGUID            string `json:"space_guid"`
	StackGUID            string `json:"stack_guid"`
	Buildpack            string `json:"buildpack"`
	DetectedBuildpack    string `json:"detected_buildpack"`
	EnvironmentJSON      string `json:"environment_json"`
	Memory               int    `json:"memory"`
	Instances            int    `json:"instances"`
	DiskQuota            int    `json:"disk_quota"`
	State                string `json:"state"`
	Version              string `json:"version"`
	Command              string `json:"command"`
	StagingTaskID        string `json:"staging_task_id"`
	PackageState         string `json:"package_state"`
	HealthCheckType      string `json:"health_check_type"`
	HealthCheckTimeout   string `json:"health_check_timeout"`
	StagingFailedReason  string `json:"staging_failed_reason"`
	Diego                bool   `json:"diego"`
	DockerImage          string `json:"docker_image"`
	PackageUpdatedAt     string `json:"package_updated_at"`
	DetectedStartCommand string `json:"detected_start_command"`
	SpaceURL             string `json:"space_url"`
	StackURL             string `json:"stack_url"`
	EventsURL            string `json:"events_url"`
	ServiceBindingsURL   string `json:"service_bindings_url"`
	RoutesURL            string `json:"routes_url"`
}
