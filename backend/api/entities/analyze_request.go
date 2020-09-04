package entities

import (
	"time"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
)

const (
	// AnalyzeRequestInputTypeCSV for when the user uploaded a CSV file
	AnalyzeRequestInputTypeCSV = "csv"
	// AnalyzeRequestInputTypeCredentials type when user submits username and password
	AnalyzeRequestInputTypeCredentials = "username/password"
)

// AnalyzeRequestStatus is the status of the request
type AnalyzeRequestStatus string

// String converts form status to string
func (status AnalyzeRequestStatus) String() string {
	return string(status)
}

const (
	// AnalyzeRequestStatusInProgress indicates that the request is in progress
	AnalyzeRequestStatusInProgress = AnalyzeRequestStatus("in-progress")
)

// AnalyzeRequest entity
type AnalyzeRequest struct {
	ID                id.ID
	UserID            id.ID
	InputType         string
	OvChipkaartNumber string
	Status            AnalyzeRequestStatus
	StartDate         time.Time
	EndDate           time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
