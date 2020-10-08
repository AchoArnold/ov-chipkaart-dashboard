package mongodb

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/database"
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

// UserRepository returns the user repository
func (db *MongoDB) UserRepository() database.UserRepository {
	return NewUserRepository(db.client, "users")
}

// AnalyzeRequestRepository returns the Analyze Request Repository
func (db *MongoDB) AnalyzeRequestRepository() database.AnalyzeRequestRepository {
	return NewAnalyzeRequestRepository(db.client, "analyze_requests")
}
