package entities

import (
	"time"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/entities"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
)

// TransactionType represents the type of transaction
type TransactionType string

// String returns the transaction type as a string
func (transactionType TransactionType) String() string {
	return string(transactionType)
}

const (
	// TransactionTypeTravel represents a travel transaction
	TransactionTypeTravel = TransactionType("Travel")

	// TransactionTypeSupplement represents a top up transaction
	TransactionTypeSupplement = TransactionType("Supplement")
)

// EnrichedRecord represents an enriched record.
type EnrichedRecord struct {
	ID               id.ID
	RawRecordID      id.ID
	StartTime        time.Time
	EndTime          time.Time
	StartTimeIsExact bool
	FromStationCode  string
	ToStationCode    string
	CompanyName      entities.CompanyName
	TransactionType  TransactionType
	Duration         time.Duration
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NSJourney returns the NSJourney for a given enriched record
func (record EnrichedRecord) NSJourney() NSJourney {
	return NewNSJourney(record.StartTime, record.FromStationCode, record.ToStationCode)
}

// IsSupplement determines if the enriched record is a supplement
func (record EnrichedRecord) IsSupplement() bool {
	return record.TransactionType == TransactionTypeSupplement
}

// IsNSJourney determines if the enriched record is an NSJourney
func (record EnrichedRecord) IsNSJourney() bool {
	return record.TransactionType == TransactionTypeTravel
}
