package cftype

// TODO: generate this code from http://apidocs.cloudfoundry.org/204/apps/retrieve_a_particular_domain.html

// RetrieveAParticularDomain Retrieve a Particular Domain
// GET /v2/domains/:guid
// http://apidocs.cloudfoundry.org/204/apps/retrieve_a_particular_domain.html
type RetrieveAParticularDomain struct {
	Metadata RetrieveAParticularDomainMetadata
	Entity   RetrieveAParticularDomainEntity
}

// RetrieveAParticularDomainMetadata ...
type RetrieveAParticularDomainMetadata struct {
	GUID string
	URL  string
}

// RetrieveAParticularDomainEntity ...
type RetrieveAParticularDomainEntity struct {
	Name string `json:"name"`
}
