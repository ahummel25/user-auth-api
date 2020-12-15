package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"github.com/src/user-auth-api/graphql/generated"
)

type mutationResolver struct{ *Resolvers }
type queryResolver struct{ *Resolvers }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolvers) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolvers) Query() generated.QueryResolver { return &queryResolver{r} }
