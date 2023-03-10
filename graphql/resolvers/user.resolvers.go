package resolvers

import (
	"context"

	"github.com/src/user-auth-api/graphql/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, params model.CreateUserInput) (*model.UserObject, error) {
	return r.UserService.CreateUser(ctx, params)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, params model.DeleteUserInput) (bool, error) {
	return r.UserService.DeleteUser(ctx, params)
}

func (r *queryResolver) AuthenticateUser(ctx context.Context, params model.AuthParams) (*model.UserObject, error) {
	return r.UserService.AuthenticateUser(ctx, params.Username, params.Password)
}
