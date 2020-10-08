package resolver

import (
	"context"
	"errors"
	raw_records_service "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/raw-records-service"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/types"
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
	grpcCtx, cancel := context.WithTimeout(grpcCtx, time.Second*5)
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

	var recordsResponse *transactions_service.TransactionsResponse
	var source types.RawRecordSource
	if analyzeRequest.InputType == entities.AnalyzeRequestInputTypeCredentials {
		source = types.RawRecordSourceAPI
		recordsResponse, err = r.transactionsServiceClient.FetchByCredentials(grpcCtx, &transactions_service.FetchByCredentialsRequest{
			Username:   *input.OvChipkaartUsername,
			Password:   *input.OvChipkaartPassword,
			CardNumber: input.OvChipkaartNumber,
			StartDate:  protoStartDate,
			EndDate:    protoEndDate,
		})
	} else {
		source = types.RawRecordSourceCSV
		data, err := ioutil.ReadAll(input.TravelHistoryFile.File)
		if err != nil {
			r.addError(ctx, "travelHistoryFile", "error while processing csv file", CodeValidationError)
			return false, err
		}

		recordsResponse, err = r.transactionsServiceClient.FetchFromBytes(grpcCtx, &transactions_service.FetchFromBytesRequest{
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

	if len(recordsResponse.Transactions) == 0 {
		r.addError(ctx, "startDate", "There are no transactions within this date range", CodeValidationError)
		r.addError(ctx, "endDate", "There are no transactions within this date range", CodeValidationError)
		return false, errors.New("error while processing ov chipkaart transactions")
	}

	err = r.db.AnalyzeRequestRepository().Store(analyzeRequest)
	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrap(err, "cannot save analyze request in the database"))
		return false, internalErrors.ErrInternalServerError
	}

	grpcCtx = context.WithValue(context.Background(), internalContext.KeyAnalyzeRequestID, analyzeRequest.ID.String())
	grpcCtx, cancel = context.WithTimeout(grpcCtx, time.Second*5)
	defer cancel()
	_, err = r.rawRecordsServiceClient.StoreTransactions(grpcCtx, &raw_records_service.StoreTransactionsRequest{
		Transactions:     recordsResponse.Transactions,
		Source:           source.String(),
		AnalyzeRequestId: analyzeRequest.ID.String(),
	})

	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrap(err, "could not save raw records"))
		return false, internalErrors.ErrInternalServerError
	}

	return true, nil
}
