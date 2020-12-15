package main

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/src/user-auth-api/graphql/generated"
	"github.com/src/user-auth-api/graphql/model"
	"github.com/src/user-auth-api/graphql/resolvers"
	"github.com/src/user-auth-api/services"
	"github.com/src/user-auth-api/utils"
)

var muxAdapter *gorillamux.GorillaMuxAdapter

func init() {
	r := mux.NewRouter()

	userService := services.NewUserService()

	initResolvers := resolvers.Resolvers{
		UserService: userService,
	}

	c := generated.Config{
		Resolvers: &initResolvers,
		Directives: generated.DirectiveRoot{
			HasRole: func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
				fc := graphql.GetFieldContext(ctx).Args

				log.Printf("%+v\n", fc["user"].(model.CreateUserInput).Email)

				return next(ctx)
			},
		},
	}

	// From server.go
	schema := generated.NewExecutableSchema(c)
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
