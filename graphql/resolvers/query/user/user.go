package user

import (
	"context"

	"github.com/ahummel25/user-auth-api/graphql/model"
	"github.com/ahummel25/user-auth-api/service/user"
)

// Resolver contains the services that user query resolver calls into
type Resolver struct {
	UserService user.API
}

func (r *Resolver) AuthenticateUser(ctx context.Context, params model.AuthParams) (*model.UserObject, error) {
	return r.UserService.AuthenticateUser(ctx, params.UsernameOrEmail, params.Password)
}
