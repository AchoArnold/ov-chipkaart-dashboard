package database

import "github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/database"

// DB is a database
type DB interface {
	RawRecordRepository() database.RawRecordRepository
}
