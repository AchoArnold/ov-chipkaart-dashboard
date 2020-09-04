package mongodb

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/database"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbOperationTimeout = 5 * time.Second

	sortOrderDescending = -1
	sortOrderAscending  = 1

	fieldCreatedAt = "created_at"

	defaultLimit = 50
	defaultSkip  = 0
)

type repository struct {
	db         *mongo.Database
	collection string
}

func (repository repository) Collection() *mongo.Collection {
	return repository.db.Collection(repository.collection)
}

func (repository repository) DefaultTimeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)
	return ctx
}

func (repository repository) GetFindOptions(take *int, skip *int, sortBy *string, sortDirection *string) *options.FindOptions {
	findOptions := options.Find()

	sortKey := fieldCreatedAt
	if sortBy != nil {
		sortKey = *sortBy
	}

	sortValue := sortOrderDescending
	if sortDirection != nil {
		if strings.ToLower(*sortDirection) == "asc" {
			sortValue = sortOrderAscending
		}
	}

	findOptions.SetSort(bson.D{{sortKey, sortValue}})

	findOptions.SetSkip(defaultSkip)
	if skip != nil {
		findOptions.SetSkip(int64(*skip))
	}

	findOptions.SetLimit(defaultLimit)
	if take != nil {
		findOptions.SetLimit(int64(*take))
	}

	return findOptions
}

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

// RawRecordRepository represents a raw record
func (db *MongoDB) RawRecordRepository() database.RawRecordRepository {
	return NewRawRecordRepository(db.client, "raw_records")
}
