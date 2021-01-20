package main

import (
	"context"
	"log"
	"os"

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

/*func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/404.html")
}*/

func init() {
	r := mux.NewRouter()
	env := os.Getenv("ENV")

	userService := services.NewUserService()

	initResolvers := resolvers.Resolvers{
		UserService: userService,
	}

	c := generated.Config{
		Resolvers: &initResolvers,
		Directives: generated.DirectiveRoot{
			HasRole: func(
				ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role, action model.Action,
			) (res interface{}, err error) {
				fc := graphql.GetFieldContext(ctx).Args

				switch action.String() {
				case model.ActionCreateUser.String():
					log.Printf("%+v\n", fc["user"].(model.CreateUserInput).Email)
					break
				case model.ActionDeleteUser.String():
					log.Printf("%+v\n", fc["user"].(model.DeleteUserInput).Email)
					break
				}

				return next(ctx)
			},
		},
	}

	// From server.go
	schema := generated.NewExecutableSchema(c)
	server := handler.NewDefaultServer(schema)

	if env != "dev" {
		r.Handle("/graphiql", playground.Handler("GraphQL playground", "/graphql"))
	}

	r.Handle("/graphql", server)
	// r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

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
