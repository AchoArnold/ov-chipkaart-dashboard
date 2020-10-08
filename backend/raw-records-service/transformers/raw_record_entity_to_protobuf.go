package transformers

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/entities"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/raw-records-service"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (t Transformers) RawRecordsToRawRecordsResponse(rawRecords []entities.RawRecord) (*raw_records_service.FetchByRequestIdResponse, error) {
	records := make([]*raw_records_service.RawRecord, 0, len(rawRecords))
	for _, rawRecord := range rawRecords {
		var fare *wrappers.DoubleValue
		if rawRecord.Fare != nil {
			fare = wrapperspb.Double(*rawRecord.Fare)
		}

		var ePurseMut *wrappers.DoubleValue
		if rawRecord.EPurseMut != nil {
			ePurseMut = wrapperspb.Double(*rawRecord.EPurseMut)
		}

		transactionDateTime, err := ptypes.TimestampProto(rawRecord.TransactionDateTime)
		if err != nil {
			return nil, err
		}

		createdAt, err := ptypes.TimestampProto(rawRecord.TransactionDateTime)
		if err != nil {
			return nil, err
		}

		updatedAt, err := ptypes.TimestampProto(rawRecord.TransactionDateTime)
		if err != nil {
			return nil, err
		}


		records = append(records, &raw_records_service.RawRecord{
			CheckInInfo:            rawRecord.CheckInInfo,
			CheckInText:            rawRecord.CheckInText,
			Fare:                   fare,
			FareCalculation:        rawRecord.FareCalculation,
			FareText:               rawRecord.FareText,
			ModalType:              rawRecord.ModalType,
			ProductInfo:            rawRecord.ProductInfo,
			ProductText:           	rawRecord.ProductText,
			Pto:                    rawRecord.Pto,
			TransactionDateTime:    transactionDateTime,
			TransactionInfo:        rawRecord.TransactionInfo,
			TransactionName:       	rawRecord.TransactionName.String(),
			EPurseMut:              ePurseMut,
			EPurseMutInfo:          rawRecord.EPurseMutInfo,
			TransactionExplanation: rawRecord.TransactionExplanation,
			TransactionPriority:    rawRecord.TransactionPriority,
			CreatedAt:              createdAt,
			UpdatedAt:              updatedAt,
			Source:                 rawRecord.Source.String(),
		})
	}

	return &raw_records_service.FetchByRequestIdResponse{RawRecords: records}, nil
}
