package database

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/entities"
)

// RawRecordRepository is an instance of the user repository
type RawRecordRepository interface {
	StoreMany(records []entities.RawRecord) error
}
