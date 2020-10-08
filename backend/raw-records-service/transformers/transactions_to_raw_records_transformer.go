package transformers

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/entities"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/transactions-service"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/types"
	"time"
)

func (t Transformers) TransactionsToRawRecords(
	transactions []*transactions_service.Transaction,
	analyzeRequestId id.ID,
	source types.RawRecordSource,
) (records []entities.RawRecord) {

	rawRecords := make([]entities.RawRecord, len(transactions))
	for index, record := range transactions {
		var fare *float64
		if record.Fare != nil {
			fare = &record.Fare.Value
		}

		var ePurseMut *float64
		if record.EPurseMut != nil {
			ePurseMut = &record.EPurseMut.Value
		}

		rawRecords[index] = entities.RawRecord{
			ID:                     id.New(),
			AnalyzeRequestID:      	analyzeRequestId,
			CheckInInfo:            record.CheckInInfo,
			CheckInText:            record.CheckInText,
			Fare:                   fare,
			FareCalculation:        record.FareCalculation,
			FareText:               record.FareText,
			ModalType:              record.ModalType,
			ProductInfo:            record.ProductInfo,
			ProductText:            record.ProductText,
			Pto:                    record.Pto,
			TransactionDateTime:    record.TransactionDateTime.AsTime(),
			TransactionInfo:        record.TransactionInfo,
			TransactionName:        types.TransactionName(record.TransactionName),
			EPurseMut:              ePurseMut,
			EPurseMutInfo:          record.GetEPurseMutInfo(),
			TransactionExplanation: record.GetTransactionExplanation(),
			TransactionPriority:    record.GetTransactionPriority(),
			Source:                 source,
			CreatedAt:              time.Now().UTC(),
			UpdatedAt:              time.Now().UTC(),
		}
	}

	return rawRecords
}
