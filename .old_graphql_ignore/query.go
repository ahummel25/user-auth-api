package graphql

import (
	"log"

	"github.com/graphql-go/graphql"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/src/user-auth-api/graphql/fields"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"user": fields.User,
	},
})

// ExecuteQuery will execute a graphql query.
func ExecuteQuery(request RequestInput) (*graphql.Result, gqlerrors.FormattedError) {
	var params = graphql.Params{
		Schema:         BuildGraphQLSchema(),
		RequestString:  request.Query,
		OperationName:  request.OperationName,
		VariableValues: request.Variables,
	}

	result := graphql.Do(params)

	if len(result.Errors) > 0 {
		log.Printf("wrong result, unexpected errors: %v", result.Errors)
		return nil, result.Errors[0]
	}
	return result, gqlerrors.FormattedError{}
}
