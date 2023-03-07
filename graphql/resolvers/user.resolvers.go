package resolvers

import (
	"context"

	"github.com/src/user-auth-api/graphql/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, params model.CreateUserInput) (*model.UserObject, error) {
	return r.UserService.CreateUser(params)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, params model.DeleteUserInput) (string, error) {
	return r.UserService.DeleteUser(params)
}

func (r *queryResolver) AuthenticateUser(ctx context.Context, params model.AuthParams) (*model.UserObject, error) {
	return r.UserService.AuthenticateUser(ctx, params.Username, params.Password)
}
