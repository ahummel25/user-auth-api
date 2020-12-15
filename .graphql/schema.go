package graphql

import (
	"github.com/graphql-go/graphql"
)

// BuildGraphQLSchema buildsthe graphql schema.
func BuildGraphQLSchema() graphql.Schema {
	var schema graphql.Schema

	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	return schema
}
