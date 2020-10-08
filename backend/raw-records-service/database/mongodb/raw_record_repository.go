package mongodb

import (
	"context"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/entities"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/errors"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/mongodb"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/types"
	"github.com/palantir/stacktrace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	fieldAnalyzeRequestId = "analyze_request_id"
	fieldTransactionDatetime = "transaction_datetime";
)

// RawRecordRepository creates a new instance of the raw record repository
type RawRecordRepository struct {
	mongodb.Repository
}

// NewRawRecordRepository creates a new instance of the raw record repository
func NewRawRecordRepository(db *mongo.Database, collection string) database.RawRecordRepository {
	return &RawRecordRepository{mongodb.NewRepository(db, collection)}
}

// StoreMany stores multiple raw records
func (repository *RawRecordRepository) StoreMany(records []entities.RawRecord) error {
	var documents []interface{}
	for _, record := range records {
		documents = append(documents, bson.M{
			"id":                      record.ID.String(),
			"user_id":                 record.UserID.String(),
			fieldAnalyzeRequestId:     record.AnalyzeRequestID.String(),
			"check_in_info":           record.CheckInInfo,
			"fare":                    record.Fare,
			"fare_calculation":        record.FareCalculation,
			"fare_text":               record.FareText,
			"modal_type":              record.ModalType,
			"product_info":            record.ProductInfo,
			"product_text":            record.ProductText,
			"pto":                     record.Pto,
			fieldTransactionDatetime:  primitive.NewDateTimeFromTime(record.TransactionDateTime),
			"transaction_info":        record.TransactionInfo,
			"transaction_name":        record.TransactionName.String(),
			"e_purse_mut":             record.EPurseMut,
			"e_purse_mut_info":        record.EPurseMutInfo,
			"transaction_explanation": record.TransactionExplanation,
			"transaction_priority":    record.TransactionPriority,
			"source":                  record.Source,
			"created_at":              primitive.NewDateTimeFromTime(record.CreatedAt),
			"updated_at":              primitive.NewDateTimeFromTime(record.UpdatedAt),
		})
	}
	_, err := repository.Collection().InsertMany(context.Background(), documents)
	if err != nil {
		return stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseError, "cannot insert raw records into the database")
	}

	return nil
}

// FetchByRequestId returns all raw records with a given request id sorted in ascending order.
func(repository *RawRecordRepository) FetchByRequestId(requestID id.ID) (records []entities.RawRecord, err error) {
	cursor, err := repository.Collection().Find(
		repository.DefaultTimeoutContext(),
		bson.M{fieldAnalyzeRequestId: requestID.String()},
		&options.FindOptions{Sort: bson.D{{fieldTransactionDatetime, mongodb.SortOrderAscending}}},
	)

	if err == mongo.ErrNoDocuments {
		return records, errors.ErrEntityNotFound
	}

	if err != nil {
		return records, stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseError, "could not fetch raw records by request id")
	}

	var rawResults []map[string]interface{}
	err = cursor.All(repository.DefaultTimeoutContext(), &rawResults)
	if err != nil {
		return records, stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseError, "could not fetch raw records by request id")
	}

	records = make([]entities.RawRecord, len(rawResults))
	for index, dbRecord := range rawResults {
		records[index], err = repository.hydrateRawRecordFromDBRecord(dbRecord)
		if err != nil {
			return records, stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseError, "error hydrating raw record into model")
		}
	}

	return records, nil
}

func (repository *RawRecordRepository) hydrateRawRecordFromDBRecord(dbRecord map[string]interface{}) (rawRecord entities.RawRecord, err error) {
	userID, err := id.FromString(dbRecord["user_id"].(string))
	if err != nil {
		return rawRecord, stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseHydrationError, "could not decode user id form string")
	}

	recordID, err := id.FromString(dbRecord["id"].(string))
	if err != nil {
		return rawRecord, stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseHydrationError, "could not decode raw record id form string")
	}

	analyzeRequestID, err := id.FromString(dbRecord["analyze_request_id"].(string))
	if err != nil {
		return rawRecord, stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseHydrationError, "could not decode raw analyze request id form string")
	}

	source, err :=  types.RawRecordSourceFromString(dbRecord["source"].(string))
	if err != nil {
		return rawRecord, stacktrace.PropagateWithCode(err, errors.ErrCodeDatabaseHydrationError, "could not create raw record source")
	}

	return entities.RawRecord{
		ID:                     recordID,
		UserID:                 userID,
		AnalyzeRequestID:       analyzeRequestID,
		CheckInInfo:            dbRecord["check_in_info"].(string),
		CheckInText:            dbRecord["check_in_info"].(string),
		Fare:                   dbRecord["fare"].(*float64),
		FareCalculation:        dbRecord["fare_calculation"].(string),
		FareText:               dbRecord["fare_text"].(string),
		ModalType:              dbRecord["modal_type"].(string),
		ProductInfo:            dbRecord["product_info"].(string),
		ProductText:            dbRecord["product_text"].(string),
		Pto:                    dbRecord["pto"].(string),
		TransactionDateTime:    dbRecord["transaction_datetime"].(primitive.DateTime).Time(),
		TransactionInfo:        dbRecord["transaction_info"].(string),
		TransactionName:        types.TransactionName(dbRecord["transaction_name"].(string)),
		EPurseMut:              dbRecord["e_purse_mut"].(*float64),
		EPurseMutInfo:          dbRecord["e_purse_mut_info"].(string),
		TransactionExplanation: dbRecord["transaction_explanation"].(string),
		TransactionPriority:    dbRecord["transaction_priority"].(string),
		Source:                	source,
		CreatedAt:              dbRecord["created_at"].(primitive.DateTime).Time(),
		UpdatedAt:              dbRecord["updated_at"].(primitive.DateTime).Time(),
	}, nil
}
