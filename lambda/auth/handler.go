package main

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/src/user-auth-api/graph/generated"
	"github.com/src/user-auth-api/graph/resolver"
	"github.com/src/user-auth-api/services"
	"github.com/src/user-auth-api/utils"
)

var muxAdapter *gorillamux.GorillaMuxAdapter

func init() {
	r := mux.NewRouter()

	authService := services.NewAuthService()

	resolvers := resolver.Resolver{
		AuthService: authService,
	}

	// From server.go
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})
	server := handler.NewDefaultServer(schema)

	r.Handle("/graphiql", playground.Handler("GraphQL playground", "/graphql"))
	r.Handle("/graphql", server)

	muxAdapter = gorillamux.New(r)
}

// LambdaHandler is our lambda handler invoked by the `lambda.Start` function call
func LambdaHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	var (
		err      error
		response = events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			IsBase64Encoded: false,
		}
	)

	if response, err = muxAdapter.Proxy(request); err != nil {
		response = utils.BuildErrorResponse(response, err.Error())

		return response, nil
	}

	return response, err
}
