package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/generated"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.AuthOutput, error) {
	return r.createUser(ctx, input)
}

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthOutput, error) {
	return r.login(ctx, input)
}

func (r *mutationResolver) CancelToken(ctx context.Context, input model.CancelTokenInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StoreAnalyzeRequest(ctx context.Context, input model.StoreAnalyzeRequestInput) (bool, error) {
	return r.storeAnalyzeRequest(ctx, input)
}

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	return &model.User{}, nil
}

func (r *queryResolver) AnalyzeRequests(ctx context.Context, skip *int, take *int, orderBy *string) ([]*model.AnalzyeRequestDetails, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) StoreRequest(ctx context.Context, input *model.StoreAnalyzeRequestInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}
