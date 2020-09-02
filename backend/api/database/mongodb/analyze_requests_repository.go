package mongodb

import (
	"context"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/entities"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/time"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AnalyzeRequestRepository creates a new instance of the user repository
type AnalyzeRequestRepository struct {
	repository
}

// NewAnalyzeRequestRepository creates a new instance of the user repository
func NewAnalyzeRequestRepository(db *mongo.Database, collection string) database.AnalyzeRequestRepository {
	return &AnalyzeRequestRepository{repository{db, collection}}
}

// Store stores a user on the mongodb repository
func (repository *AnalyzeRequestRepository) Store(analyzeRequest entities.AnalyzeRequest) error {
	_, err := repository.Collection().InsertOne(context.Background(), bson.M{
		"id":                  analyzeRequest.ID.String(),
		"user_id":             analyzeRequest.UserID.String(),
		"input_type":          analyzeRequest.InputType,
		"ov_chipkaart_number": analyzeRequest.OvChipkaartNumber,
		"start_date":          analyzeRequest.StartDate.Format(time.DateFormat),
		"end_date":            analyzeRequest.EndDate.Format(time.DateFormat),
		"status":              string(analyzeRequest.Status),
		"created_at":          analyzeRequest.CreatedAt,
		"updated_at":          analyzeRequest.UpdatedAt,
	})

	return err
}

// FindByID finds a user in the database using it's ID
func (repository *AnalyzeRequestRepository) FindByID(ID id.ID) (analyzeRequest entities.AnalyzeRequest, err error) {
	dbRecord := map[string]interface{}{}
	err = repository.Collection().FindOne(repository.DefaultTimeoutContext(), bson.M{"id": ID.String()}).Decode(&dbRecord)

	if err == mongo.ErrNoDocuments {
		return analyzeRequest, database.ErrEntityNotFound
	}
	if err != nil {
		return analyzeRequest, errors.Wrap(err, "error fetching single analzye request from the database by id")
	}

	return repository.hydrateAnalyzeFromDBRecord(dbRecord)
}

func (repository *AnalyzeRequestRepository) hydrateAnalyzeFromDBRecord(dbRecord map[string]interface{}) (analyzeRequest entities.AnalyzeRequest, err error) {
	requestID, err := id.FromString(dbRecord["id"].(string))
	if err != nil {
		return analyzeRequest, errors.Wrap(err, "could not decode analyze request id form string")
	}

	userID, err := id.FromString(dbRecord["user_id"].(string))
	if err != nil {
		return analyzeRequest, errors.Wrap(err, "could not decode user id form string")
	}

	startDate, err := time.FromDate(dbRecord["start_date"].(string))
	if err != nil {
		return analyzeRequest, errors.Wrap(err, "cannot decode start date")
	}

	endDate, err := time.FromDate(dbRecord["end_date"].(string))
	if err != nil {
		return analyzeRequest, errors.Wrap(err, "cannot decode end date")
	}

	return entities.AnalyzeRequest{
		ID:                requestID,
		UserID:            userID,
		StartDate:         startDate,
		EndDate:           endDate,
		InputType:         dbRecord["input_type"].(string),
		OvChipkaartNumber: dbRecord["ov_chipkaart_number"].(string),
		CreatedAt:         dbRecord["created_at"].(primitive.DateTime).Time(),
		UpdatedAt:         dbRecord["updated_at"].(primitive.DateTime).Time(),
	}, err
}
