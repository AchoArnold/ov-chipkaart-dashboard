package database

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/entities"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
)

// AnalyzeRequestRepository is an instance of the user repository
type AnalyzeRequestRepository interface {
	Store(analyzeRequest entities.AnalyzeRequest) error
	FindByID(analyzeRequestID id.ID) (entities.AnalyzeRequest, error)
}
