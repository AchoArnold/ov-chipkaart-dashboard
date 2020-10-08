package resolver

import (
	"context"
	"errors"
	"io/ioutil"
	"time"

	internalContext "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/context"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/entities"
	internalErrors "github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/errors"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/graph/model"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/transactions-service"
	internalTime "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/time"
	"github.com/golang/protobuf/ptypes"
	pkgErrors "github.com/pkg/errors"
)

func (r *mutationResolver) storeAnalyzeRequest(ctx context.Context, input model.StoreAnalyzeRequestInput) (bool, error) {
	// check that the user is authorized
	userID, err := r.userIDFromContext(ctx)
	if err != nil {
		return false, ErrUnauthorizedRequest
	}

	validationResult := r.validator.ValidateStoreAnalzyeRequest(input, r.languageTagFromContext(ctx))
	if validationResult.HasError {
		r.addValidationErrors(ctx, validationResult)
		return false, internalErrors.ErrValidationError
	}

	inputType := entities.AnalyzeRequestInputTypeCSV
	if input.OvChipkaartUsername != nil {
		inputType = entities.AnalyzeRequestInputTypeCredentials
	}
	startDate, _ := internalTime.FromDate(input.StartDate)
	endDate, _ := internalTime.FromDate(input.EndDate)

	analyzeRequest := entities.AnalyzeRequest{
		ID:                id.New(),
		UserID:            userID,
		InputType:         inputType,
		OvChipkaartNumber: input.OvChipkaartNumber,
		StartDate:         startDate,
		EndDate:           endDate,
		Status:            entities.AnalyzeRequestStatusInProgress,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	grpcCtx := context.WithValue(context.Background(), internalContext.KeyAnalyzeRequestID, analyzeRequest.ID.String())
	grpcCtx, cancel := context.WithTimeout(grpcCtx, time.Second*15)
	defer cancel()

	protoStartDate, err := ptypes.TimestampProto(analyzeRequest.StartDate)
	if err != nil {
		r.errorHandler.CaptureError(ctx, err)
		return false, internalErrors.ErrInternalServerError
	}

	protoEndDate, err := ptypes.TimestampProto(analyzeRequest.EndDate)
	if err != nil {
		r.errorHandler.CaptureError(ctx, err)
		return false, internalErrors.ErrInternalServerError
	}

	var recordsResponse *transactions_service.RawRecordsResponse
	var source entities.RawRecordSource
	if analyzeRequest.InputType == entities.AnalyzeRequestInputTypeCredentials {
		source = entities.RawRecordSourceAPI
		recordsResponse, err = r.transactionsServiceClient.RawRecordsWithCredentials(grpcCtx, &transactions_service.CredentialsRawRecordsRequest{
			Username:   *input.OvChipkaartUsername,
			Password:   *input.OvChipkaartPassword,
			CardNumber: input.OvChipkaartNumber,
			StartDate:  protoStartDate,
			EndDate:    protoEndDate,
		})
	} else {
		source = entities.RawRecordSourceCSV
		data, err := ioutil.ReadAll(input.TravelHistoryFile.File)
		if err != nil {
			r.addError(ctx, "travelHistoryFile", "error while processing csv file", CodeValidationError)
			return false, err
		}

		recordsResponse, err = r.transactionsServiceClient.RawRecordsFromBytes(grpcCtx, &transactions_service.BytesRawRecordsRequest{
			CardNumber: analyzeRequest.OvChipkaartNumber,
			StartDate:  protoStartDate,
			EndDate:    protoEndDate,
			Data:       data,
		})
	}

	if err != nil {
		r.errorHandler.CaptureError(grpcCtx, err)
		return false, errors.New("error while fetching ov chipkaart transactions")
	}

	if len(recordsResponse.RawRecords) == 0 {
		r.addError(ctx, "startDate", "There are no transactions within this date range", CodeValidationError)
		r.addError(ctx, "endDate", "There are no transactions within this date range", CodeValidationError)
		return false, errors.New("error while processing ov chipkaart transactions")
	}

	rawRecords := make([]entities.RawRecord, len(recordsResponse.RawRecords))
	for index, record := range recordsResponse.RawRecords {
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
			UserID:                 analyzeRequest.UserID,
			AnalyzeRequestID:       analyzeRequest.ID,
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
			TransactionName:        entities.TransactionName(record.TransactionName),
			EPurseMut:              ePurseMut,
			EPurseMutInfo:          record.GetEPurseMutInfo(),
			TransactionExplanation: record.GetTransactionExplanation(),
			TransactionPriority:    record.GetTransactionPriority(),
			Source:                 source,
			CreatedAt:              time.Now().UTC(),
			UpdatedAt:              time.Now().UTC(),
		}
	}

	err = r.db.AnalyzeRequestRepository().Store(analyzeRequest)
	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrap(err, "cannot save analyze request in the database"))
		return false, internalErrors.ErrInternalServerError
	}

	err = r.db.RawRecordRepository().StoreMany(rawRecords)
	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrap(err, "could not save raw records"))
		return false, internalErrors.ErrInternalServerError
	}

	return true, nil
}
