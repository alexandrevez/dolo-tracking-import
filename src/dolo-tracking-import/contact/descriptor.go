package contact

// Descriptor is what we retrieve from the CSV file
type Descriptor struct {
	CompanyName string
	DomainName  string
	Email       string
}

// NewDescriptor builds a descriptor yo
func NewDescriptor(companyName string, domainName string, email string) Descriptor {
	return Descriptor{
		CompanyName: companyName,
		DomainName:  domainName,
		Email:       email,
	}
}
