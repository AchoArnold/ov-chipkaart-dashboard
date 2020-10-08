package handlers

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/transformers"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/logger"
	raw_records_service "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/raw-records-service"
)

type Server struct {
	raw_records_service.UnimplementedRawRecordsServiceServer
	DB database.DB
	Logger logger.Logger
	Transformers transformers.Transformers
}