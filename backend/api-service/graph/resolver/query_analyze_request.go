package resolver

import (
	"context"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/errors"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/model"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/time"
	"github.com/palantir/stacktrace"
)

func (r *queryResolver) analyzeRequests(ctx context.Context, skip *int, limit *int, orderBy *string, orderDirection *string) ([]*model.AnalyzeRequest, error) {
	// check that the user is authorized
	userID, err := r.userIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorizedRequest
	}

	validationResult := r.validator.ValidateAnalzyeRequestsInput(skip, limit, orderBy, orderDirection, r.languageTagFromContext(ctx))
	if validationResult.HasError {
		r.addValidationErrors(ctx, validationResult)
		return nil, errors.ErrValidationError
	}

	dbResults, err := r.db.AnalyzeRequestRepository().IndexForUser(userID, skip, limit, orderBy, orderDirection)
	if err != nil {
		r.errorHandler.CaptureError(ctx, stacktrace.Propagate(err, "error handling the analzye requests query"))
		return nil, errors.ErrInternalServerError
	}

	results := make([]*model.AnalyzeRequest, len(dbResults))
	for index, dbResult := range dbResults {
		results[index] = &model.AnalyzeRequest{
			StartDate:         dbResult.StartDate.Format(time.DateFormat),
			EndDate:           dbResult.EndDate.Format(time.DateFormat),
			OvChipkaartNumber: dbResult.OvChipkaartNumber,
			ID:                dbResult.ID.String(),
			Status:            dbResult.Status.String(),
			CreatedAt:         dbResult.CreatedAt.Format(time.DefaultFormat),
			UpdatedAt:         dbResult.UpdatedAt.Format(time.DefaultFormat),
		}
	}

	return results, nil
}
