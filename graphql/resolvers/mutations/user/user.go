package user

import (
	"context"

	"github.com/src/user-auth-api/graphql/model"
	"github.com/src/user-auth-api/service/user"
)

// Resolver contains the services that user mutation resolver calls into
type Resolver struct {
	UserService user.API
}

func (r *Resolver) CreateUser(ctx context.Context, params model.CreateUserInput) (*model.UserObject, error) {
	return r.UserService.CreateUser(ctx, params)
}

func (r *Resolver) DeleteUser(ctx context.Context, userID string) (bool, error) {
	return r.UserService.DeleteUser(ctx, userID)
}
