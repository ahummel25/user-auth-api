package resolvers

import (
	"github.com/src/user-auth-api/graphql/generated"
	"github.com/src/user-auth-api/service/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Services struct {
	UserService user.API
}

type mutationResolver struct{ *Services }
type queryResolver struct{ *Services }

// Mutation returns generated.MutationResolver implementation.
func (r *Services) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Services) Query() generated.QueryResolver { return &queryResolver{r} }
