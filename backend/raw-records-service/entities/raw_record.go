package entities

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/types"
	"time"
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
	TransactionName        types.TransactionName
	EPurseMut              *float64
	EPurseMutInfo          string
	TransactionExplanation string
	TransactionPriority    string
	Source                 types.RawRecordSource
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// IsCheckIn determines if a record is a check in record
func (record RawRecord) IsCheckIn() bool {
	return record.TransactionName.IsTheSameAs(types.TransactionNameCheckIn)
}

// IsNSSupplement determines if a records is a surcharge
func (record RawRecord) IsNSSupplement() bool {
	return record.TransactionName.IsTheSameAs(types.TransactionNameIntercityDirectSurcharge)
}

// IsCheckOut determines if a record is checkout transaction.
func (record RawRecord) IsCheckOut() bool {
	return record.TransactionName.IsTheSameAs(types.TransactionNameCheckOut)
}

// IsRET is used to determine if a raw record is from the RET company
func (record RawRecord) IsRET() bool {
	return record.Pto == types.CompanyNameRET.String()
}

// IsNS is used to determine if a raw record is from the NS company
func (record RawRecord) IsNS() bool {
	return record.Pto == types.CompanyNameNS.String() && record.ModalType == "Trein"
}
