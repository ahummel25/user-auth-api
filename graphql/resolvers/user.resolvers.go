package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/src/user-auth-api/graphql/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, params model.CreateUserInput) (*model.UserObject, error) {
	return r.UserService.CreateUser(params)
}

func (r *queryResolver) AuthenticateUser(ctx context.Context, params model.AuthParams) (*model.UserObject, error) {
	return r.UserService.AuthenticateUser(params.Email, params.Password)
}
