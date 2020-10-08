package handlers

import (
	"context"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/errors"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	raw_records_service "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/raw-records-service"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

func (s *Server) FetchByRequestId(_ context.Context, request *raw_records_service.FetchByRequestIdRequest) (response *raw_records_service.FetchByRequestIdResponse, err error) {
	requestID, err := id.FromString(request.GetRequestID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	records, err := s.DB.RawRecordRepository().FetchByRequestId(requestID)
	if err != nil && err == errors.ErrEntityNotFound {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response, err = s.Transformers.RawRecordsToRawRecordsResponse(records)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, err
}
