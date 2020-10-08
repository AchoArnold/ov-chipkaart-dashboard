package database

import "github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/database"

// DB is a database
type DB interface {
	RawRecordRepository() database.RawRecordRepository
}
