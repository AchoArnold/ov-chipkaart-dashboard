package mongodb

import (
	"context"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/mongodb"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/entities"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/errors"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/time"
	"github.com/palantir/stacktrace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AnalyzeRequestRepository creates a new instance of the user repository
type AnalyzeRequestRepository struct {
	mongodb.Repository
}

// NewAnalyzeRequestRepository creates a new instance of the user repository
func NewAnalyzeRequestRepository(db *mongo.Database, collection string) database.AnalyzeRequestRepository {
	return &AnalyzeRequestRepository{mongodb.NewRepository(db, collection)}
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
		"created_at":          primitive.NewDateTimeFromTime(analyzeRequest.CreatedAt),
		"updated_at":          primitive.NewDateTimeFromTime(analyzeRequest.UpdatedAt),
	})

	return err
}

// IndexForUser fetches all the analyze requests for the given user
func (repository *AnalyzeRequestRepository) IndexForUser(userID id.ID, skip *int, limit *int, sortBy *string, sortDirection *string) (analyzeRequests []entities.AnalyzeRequest, err error) {
	cursor, err := repository.Collection().Find(repository.DefaultTimeoutContext(), bson.M{"user_id": userID.String()}, repository.GetFindOptions(skip, limit, sortBy, sortDirection))
	if err != nil {
		return analyzeRequests, stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseError, "error fetching analyze requests from the database")
	}

	var rawResults []map[string]interface{}
	err = cursor.All(repository.DefaultTimeoutContext(), &rawResults)
	if err != nil {
		return analyzeRequests, stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseError, "could not analyze requests from the response")
	}

	analyzeRequests = make([]entities.AnalyzeRequest, len(rawResults))
	for index, val := range rawResults {
		analyzeRequests[index], err = repository.hydrateAnalyzeRequestFromDBRecord(val)
		if err != nil {
			return analyzeRequests, stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseError, "error hydrating raw record into model")
		}
	}

	return analyzeRequests, nil
}

// FindByID finds a user in the database using it's ID
func (repository *AnalyzeRequestRepository) FindByID(ID id.ID) (analyzeRequest entities.AnalyzeRequest, err error) {
	dbRecord := map[string]interface{}{}
	err = repository.Collection().FindOne(repository.DefaultTimeoutContext(), bson.M{"id": ID.String()}).Decode(&dbRecord)

	if err == mongo.ErrNoDocuments {
		return analyzeRequest, errors.ErrEntityNotFound
	}
	if err != nil {
		return analyzeRequest, stacktrace.Propagate(err, "error fetching single analzye request from the database by id")
	}

	return repository.hydrateAnalyzeRequestFromDBRecord(dbRecord)
}

func (repository *AnalyzeRequestRepository) hydrateAnalyzeRequestFromDBRecord(dbRecord map[string]interface{}) (analyzeRequest entities.AnalyzeRequest, err error) {
	requestID, err := id.FromString(dbRecord["id"].(string))
	if err != nil {
		return analyzeRequest, stacktrace.Propagate(err, "could not decode analyze request id form string")
	}

	userID, err := id.FromString(dbRecord["user_id"].(string))
	if err != nil {
		return analyzeRequest, stacktrace.Propagate(err, "could not decode user id form string")
	}

	startDate, err := time.FromDate(dbRecord["start_date"].(string))
	if err != nil {
		return analyzeRequest, stacktrace.Propagate(err, "cannot decode start date")
	}

	endDate, err := time.FromDate(dbRecord["end_date"].(string))
	if err != nil {
		return analyzeRequest, stacktrace.Propagate(err, "cannot decode end date")
	}

	return entities.AnalyzeRequest{
		ID:                requestID,
		UserID:            userID,
		StartDate:         startDate,
		EndDate:           endDate,
		Status:            entities.AnalyzeRequestStatus(dbRecord["status"].(string)),
		InputType:         dbRecord["input_type"].(string),
		OvChipkaartNumber: dbRecord["ov_chipkaart_number"].(string),
		CreatedAt:         dbRecord["created_at"].(primitive.DateTime).Time(),
		UpdatedAt:         dbRecord["updated_at"].(primitive.DateTime).Time(),
	}, err
}
