package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/src/user-auth-api/graph/generated"
	"github.com/src/user-auth-api/graph/model"
)

type mutationResolver struct{ *Resolvers }
type queryResolver struct{ *Resolvers }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolvers) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolvers) Query() generated.QueryResolver { return &queryResolver{r} }

func (r *mutationResolver) CreateUser(ctx context.Context, user model.CreateUserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) UserLogin(ctx context.Context, params model.AuthParams) (*model.User, error) {
	return r.AuthService.AuthenticateUser(params.Email, params.Password)
}
