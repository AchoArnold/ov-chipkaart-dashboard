package database

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/entities"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
)

// AnalyzeRequestRepository is an instance of the user repository
type AnalyzeRequestRepository interface {
	Store(analyzeRequest entities.AnalyzeRequest) error
	FindByID(analyzeRequestID id.ID) (entities.AnalyzeRequest, error)
	IndexForUser(userID id.ID, skip *int, limit *int, sortBy *string, sortDirection *string) (analyzeRequests []entities.AnalyzeRequest, err error)
}
