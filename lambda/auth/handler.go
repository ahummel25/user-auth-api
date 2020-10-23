package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/src/user-auth-api/graphql/generated"
	"github.com/src/user-auth-api/graphql/model"
	"github.com/src/user-auth-api/graphql/resolvers"
	"github.com/src/user-auth-api/service/user"
	"github.com/src/user-auth-api/utils"
)

var muxAdapter *gorillamux.GorillaMuxAdapter

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/404.html")
}

func init() {
	r := mux.NewRouter()
	r.Use(injectDBCollection)

	userService := user.New()
	initResolvers := resolvers.Services{
		UserService: userService,
	}
	cfg := generated.Config{
		Resolvers: &initResolvers,
		Directives: generated.DirectiveRoot{
			HasRole: func(
				ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role, action model.Action,
			) (res interface{}, err error) {
				fc := graphql.GetFieldContext(ctx).Args

				switch action.String() {
				case model.ActionCreateUser.String():
					log.Printf("%+v\n", fc["user"].(model.CreateUserInput).Email)
				case model.ActionDeleteUser.String():
					log.Printf("%+v\n", fc["user"].(model.DeleteUserInput).Email)
				}

				return next(ctx)
			},
		},
	}
	schema := generated.NewExecutableSchema(cfg)
	server := handler.NewDefaultServer(schema)

	r.Handle("/graphiql", playground.Handler("GraphQL playground", "/graphql"))
	r.Handle("/graphql", server)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	muxAdapter = gorillamux.New(r)
}

// LambdaHandler is our lambda handler invoked by the `lambda.Start` function call
func LambdaHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	var (
		apiGWResponse                events.APIGatewayProxyResponse
		err                          error
		switchableAPIGatewayResponse *core.SwitchableAPIGatewayResponse
	)
	if switchableAPIGatewayResponse, err = muxAdapter.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request)); err != nil {
		apiGWResponse = utils.BuildErrorResponse(apiGWResponse, err.Error())
		return apiGWResponse, nil
	}

	apiGWResponse = *switchableAPIGatewayResponse.Version1()
	apiGWResponse.Headers = map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
		"Content-Type":                     "application/json",
	}
	apiGWResponse.IsBase64Encoded = false
	apiGWResponse.StatusCode = http.StatusOK

	return apiGWResponse, nil
}
