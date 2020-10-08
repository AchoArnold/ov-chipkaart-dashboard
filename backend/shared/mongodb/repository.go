package mongodb

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbOperationTimeout = 5 * time.Second

	// SortOrderDescending MongoDB's sort descending flag
	SortOrderDescending = -1

	// SortOrderAscending MongoDB sort ascending flag
	SortOrderAscending  = 1

	fieldCreatedAt = "created_at"

	defaultLimit = 50
	defaultSkip  = 0
)

type Repository struct {
	db         *mongo.Database
	collection string
}

func NewRepository(db *mongo.Database, collection string) Repository {
	return Repository{
		db:         db,
		collection: collection,
	}
}

func (repository Repository) DB() *mongo.Database {
	return repository.db
}

func (repository Repository) Collection() *mongo.Collection {
	return repository.DB().Collection(repository.collection)
}

func (repository Repository) DefaultTimeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)
	return ctx
}

func (repository Repository) GetFindOptions(skip *int, take *int, sortBy *string, sortDirection *string) *options.FindOptions {
	findOptions := options.Find()

	sortKey := fieldCreatedAt
	if sortBy != nil {
		sortKey = *sortBy
	}

	sortValue := SortOrderDescending
	if sortDirection != nil {
		if strings.ToLower(*sortDirection) == "asc" {
			sortValue = SortOrderAscending
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

