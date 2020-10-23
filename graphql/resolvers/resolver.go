package resolvers

import (
	"github.com/src/user-auth-api/graphql/generated"
	"github.com/src/user-auth-api/graphql/resolvers/mutations"
	"github.com/src/user-auth-api/graphql/resolvers/query"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Services struct {
	mutations.MutationResolvers
	query.QueryResolvers
}

type MutationResolver struct{ *Services }
type QueryResolver struct{ *Services }

// Mutation returns generated.MutationResolver implementation.
func (r *Services) Mutation() generated.MutationResolver { return &MutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Services) Query() generated.QueryResolver { return &QueryResolver{r} }
