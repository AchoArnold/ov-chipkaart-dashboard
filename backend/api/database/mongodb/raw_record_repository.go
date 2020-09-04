package mongodb

import (
	"context"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/entities"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RawRecordRepository creates a new instance of the raw record repository
type RawRecordRepository struct {
	repository
}

// NewRawRecordRepository creates a new instance of the raw record repository
func NewRawRecordRepository(db *mongo.Database, collection string) database.RawRecordRepository {
	return &RawRecordRepository{repository{db, collection}}
}

// StoreMany stores multiple raw records
func (repository *RawRecordRepository) StoreMany(records []entities.RawRecord) error {
	var documents []interface{}
	for _, record := range records {
		documents = append(documents, bson.M{
			"id":                      record.ID.String(),
			"user_id":                 record.UserID.String(),
			"analyze_request_id":      record.AnalyzeRequestID.String(),
			"check_in_info":           record.CheckInInfo,
			"fare":                    record.Fare,
			"fare_calculation":        record.FareCalculation,
			"fare_text":               record.FareText,
			"modal_type":              record.ModalType,
			"product_info":            record.ProductInfo,
			"product_text":            record.ProductText,
			"pto":                     record.Pto,
			"transaction_datetime":    primitive.NewDateTimeFromTime(record.TransactionDateTime),
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
	_, err := repository.db.Collection(repository.collection).InsertMany(context.Background(), documents)
	if err != nil {
		return errors.Wrapf(err, "cannot insert raw records into the database")
	}

	return nil
}

func (repository *RawRecordRepository) hydrateRawRecordFromDBRecord(dbRecord map[string]interface{}) (rawRecord entities.RawRecord, err error) {
	userID, err := id.FromString(dbRecord["user_id"].(string))
	if err != nil {
		return rawRecord, errors.Wrap(err, "could not decode user id form string")
	}

	recordID, err := id.FromString(dbRecord["id"].(string))
	if err != nil {
		return rawRecord, errors.Wrap(err, "could not decode raw record id form string")
	}

	analyzeRequestID, err := id.FromString(dbRecord["analyze_request_id"].(string))
	if err != nil {
		return rawRecord, errors.Wrap(err, "could not decode raw analyze request id form string")
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
		TransactionName:        entities.TransactionName(dbRecord["transaction_name"].(string)),
		EPurseMut:              dbRecord["e_purse_mut"].(*float64),
		EPurseMutInfo:          dbRecord["e_purse_mut_info"].(string),
		TransactionExplanation: dbRecord["transaction_explanation"].(string),
		TransactionPriority:    dbRecord["transaction_priority"].(string),
		Source:                 entities.RawRecordSource(dbRecord["source"].(string)),
		CreatedAt:              dbRecord["created_at"].(primitive.DateTime).Time(),
		UpdatedAt:              dbRecord["updated_at"].(primitive.DateTime).Time(),
	}, nil
}
