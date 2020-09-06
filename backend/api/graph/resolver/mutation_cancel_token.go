package resolver

import (
	"context"

	internalErrors "github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/errors"
)

func (r *mutationResolver) cancelToken(ctx context.Context) (bool, error) {
	jwt, err := r.tokenFromContext(ctx)
	if err != nil {
		r.errorHandler.CaptureError(ctx, err)
		return false, ErrUnauthorizedRequest
	}

	if !r.jwtService.IsValid(jwt) {
		return false, ErrUnauthorizedRequest
	}

	err = r.jwtService.InvalidateToken(jwt)
	if err != nil {
		r.errorHandler.CaptureError(ctx, err)
		return false, internalErrors.ErrInternalServerError
	}

	return true, nil
}
