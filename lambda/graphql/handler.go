package main

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/ahummel25/user-auth-api/db"
	"github.com/ahummel25/user-auth-api/graphql/directives"
	"github.com/ahummel25/user-auth-api/graphql/generated"
	"github.com/ahummel25/user-auth-api/graphql/resolvers"
	"github.com/ahummel25/user-auth-api/graphql/resolvers/mutations"
	userMutation "github.com/ahummel25/user-auth-api/graphql/resolvers/mutations/user"
	"github.com/ahummel25/user-auth-api/graphql/resolvers/query"
	userQuery "github.com/ahummel25/user-auth-api/graphql/resolvers/query/user"
	"github.com/ahummel25/user-auth-api/service/user"
)

var muxAdapter *gorillamux.GorillaMuxAdapter

func defaultTranslation() {
	directives.ValidateAddTranslation("email", " must be a valid email address")
}

func init() {
	r := mux.NewRouter()
	defaultTranslation()

	userService := user.New()
	mutationResolvers := mutations.MutationResolvers{Resolver: userMutation.Resolver{
		UserService: userService,
	}}
	queryResolvers := query.QueryResolvers{Resolver: userQuery.Resolver{
		UserService: userService,
	}}
	resolvers := resolvers.Services{
		MutationResolvers: mutationResolvers,
		QueryResolvers:    queryResolvers,
	}
	cfg := generated.Config{
		Resolvers: &resolvers,
	}
	cfg.Directives.Binding = directives.Binding
	cfg.Directives.HasRole = directives.HasRole
	schema := generated.NewExecutableSchema(cfg)
	server := handler.NewDefaultServer(schema)

	r.Handle("/graphiql", playground.Handler("GraphQL playground", "/graphql"))
	r.Handle("/apollo", playground.ApolloSandboxHandler("GraphQL Apollo playground", "/graphql"))
	r.Handle("/graphql", server)
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

	// Setup DB context
	ctx, err = db.SetupDBContext(ctx)
	if err != nil {
		apiGWResponse.StatusCode = http.StatusInternalServerError
		apiGWResponse.Body = err.Error()
		return apiGWResponse, nil
	}

	if switchableAPIGatewayResponse, err = muxAdapter.ProxyWithContext(
		ctx,
		*core.NewSwitchableAPIGatewayRequestV1(&request),
	); err != nil {
		return apiGWResponse, err
	}

	apiGWResponse = *switchableAPIGatewayResponse.Version1()
	apiGWResponse.IsBase64Encoded = false
	apiGWResponse.StatusCode = http.StatusOK
	return apiGWResponse, nil
}
