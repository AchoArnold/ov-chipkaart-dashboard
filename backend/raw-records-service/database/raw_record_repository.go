package database

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/entities"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
)

// RawRecordRepository is an instance of the user repository
type RawRecordRepository interface {
	StoreMany(records []entities.RawRecord) error
	FetchByRequestId(requestID id.ID) (records []entities.RawRecord, err error)
}
