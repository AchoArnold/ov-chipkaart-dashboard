package entities

import (
	"strings"
	"time"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
)

// RawRecordSource is the source for a raw record entry
type RawRecordSource string

// String converts a raw record source to a string
func (source RawRecordSource) String() string {
	return string(source)
}

const (
	// RawRecordSourceAPI is an API source
	RawRecordSourceAPI = RawRecordSource("API")

	// RawRecordSourceCSV is a CSV source
	RawRecordSourceCSV = RawRecordSource("CSV")
)

// TransactionName  represents the various transaction names
type TransactionName string

// String returns the transaction name as a string
func (name TransactionName) String() string {
	return string(name)
}

// IsTheSameAs is used to compare 2 transaction names
func (name TransactionName) IsTheSameAs(comp TransactionName) bool {
	return strings.ToLower(name.String()) == strings.ToLower(comp.String())
}

const (
	// TransactionNameCheckIn represents a check in transaction
	TransactionNameCheckIn = TransactionName("Check-in")

	// TransactionNameCheckOut represents a check out transaction
	TransactionNameCheckOut = TransactionName("Check-uit")

	// TransactionNameIntercityDirectSurcharge represents an intercity direct transaction
	TransactionNameIntercityDirectSurcharge = TransactionName("Toeslag Intercity Direct")
)

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

// RawRecord represents a transaction record
type RawRecord struct {
	ID                     id.ID
	UserID                 id.ID
	AnalyzeRequestID       id.ID
	CheckInInfo            string
	CheckInText            string
	Fare                   *float64
	FareCalculation        string
	FareText               string
	ModalType              string
	ProductInfo            string
	ProductText            string
	Pto                    string
	TransactionDateTime    time.Time
	TransactionInfo        string
	TransactionName        TransactionName
	EPurseMut              *float64
	EPurseMutInfo          string
	TransactionExplanation string
	TransactionPriority    string
	Source                 RawRecordSource
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// IsCheckIn determines if a record is a check in record
func (record RawRecord) IsCheckIn() bool {
	return record.TransactionName.IsTheSameAs(TransactionNameCheckIn)
}

// IsNSSupplement determines if a records is a surcharge
func (record RawRecord) IsNSSupplement() bool {
	return record.TransactionName.IsTheSameAs(TransactionNameIntercityDirectSurcharge)
}

// IsCheckOut determines if a record is checkout transaction.
func (record RawRecord) IsCheckOut() bool {
	return record.TransactionName.IsTheSameAs(TransactionNameCheckOut)
}

// IsRET is used to determine if a raw record is from the RET company
func (record RawRecord) IsRET() bool {
	return record.Pto == CompanyNameRET.String()
}

// IsNS is used to determine if a raw record is from the NS company
func (record RawRecord) IsNS() bool {
	return record.Pto == CompanyNameNS.String() && record.ModalType == "Trein"
}
