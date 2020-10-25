package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/src/user-auth-api/graph/generated"
	"github.com/src/user-auth-api/graph/model"
)

type queryResolver struct{ *Resolver }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

func (r *queryResolver) UserLogin(ctx context.Context, params model.AuthParams) (*model.User, error) {
	return r.AuthService.AuthenticateUser(params.Username, params.Password)
}
