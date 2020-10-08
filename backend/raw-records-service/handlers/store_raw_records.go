package handlers

import (
	"context"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	raw_records_service "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/raw-records-service"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/types"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) StoreTransactions(_ context.Context, request *raw_records_service.StoreTransactionsRequest) (*empty.Empty, error) {
	analyzeRequestId, err := id.FromString(request.AnalyzeRequestId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	source, err := types.RawRecordSourceFromString(request.Source)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rawRecords := s.Transformers.TransactionsToRawRecords(request.Transactions, analyzeRequestId, source)

	err = s.DB.RawRecordRepository().StoreMany(rawRecords)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}