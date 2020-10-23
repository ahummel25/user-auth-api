package fields

import (
	"github.com/graphql-go/graphql"

	"github.com/src/user-auth-api/graphql/resolvers"
	"github.com/src/user-auth-api/graphql/types"
)

// User represents a single user.
var User = &graphql.Field{
	Name: "User",
	Type: types.UserObject,
	Args: graphql.FieldConfigArgument{
		"username": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"password": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: resolvers.UserResolver,
}
