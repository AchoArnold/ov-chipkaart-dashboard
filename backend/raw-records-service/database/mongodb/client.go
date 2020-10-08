package mongodb

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/database"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB is the struct for mongodb
type MongoDB struct {
	client *mongo.Database
}

// NewMongoDB creates a new instance of the mongodb client
func NewMongoDB(client *mongo.Database) database.DB {
	return &MongoDB{
		client: client,
	}
}

// RawRecordRepository represents a raw record
func (db *MongoDB) RawRecordRepository() database.RawRecordRepository {
	return NewRawRecordRepository(db.client, "raw_records")
}
