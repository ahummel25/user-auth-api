package graphql

import (
	"errors"
	"log"

	"github.com/graphql-go/graphql"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/src/user-auth-api/graphql/types"
)

var errInvalidPassword = errors.New("invalid password")

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Name: "User",
			Type: types.User,
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userPassword := types.UserMap[params.Args["username"].(string)]

				if ok := userPassword == params.Args["password"]; !ok {
					return nil, errInvalidPassword
				}

				user := types.UserType{
					ID:   1,
					Name: "Andrew",
				}

				return user, nil
			},
		},
	},
})

// ExecuteQuery will execute a graphql query.
func ExecuteQuery(request RequestInput, schema graphql.Schema) (*graphql.Result, gqlerrors.FormattedError) {
	var params = graphql.Params{
		Schema:         schema,
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
