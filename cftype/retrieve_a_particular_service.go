package cftype

// TODO: generate this code from http://apidocs.cloudfoundry.org/204/apps/retrieve_a_particular_service.html

// RetrieveAParticularService Retrieve a Particular Service
// GET /v2/services/:guid
// http://apidocs.cloudfoundry.org/204/apps/retrieve_a_particular_service.html
type RetrieveAParticularService struct {
	Metadata RetrieveAParticularServiceMetadata
	Entity   RetrieveAParticularServiceEntity
}

// RetrieveAParticularServiceMetadata ...
type RetrieveAParticularServiceMetadata struct {
	GUID string
	URL  string
}

// RetrieveAParticularServiceEntity ...
type RetrieveAParticularServiceEntity struct {
	Label string `json:"label"`
}
