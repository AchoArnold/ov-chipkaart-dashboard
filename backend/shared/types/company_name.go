package types

// CompanyName is the company to which a transaction belongs
type CompanyName string

// String returns the company name as a string
func (companyName CompanyName) String() string {
	return string(companyName)
}

const (
	// CompanyNameNS is the company name for the NS train
	CompanyNameNS = CompanyName("NS")
	// CompanyNameRET is the company name for RET
	CompanyNameRET = CompanyName("RET")
)
