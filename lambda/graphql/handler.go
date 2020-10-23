package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/src/user-auth-api/graphql/directives"
	"github.com/src/user-auth-api/graphql/generated"
	"github.com/src/user-auth-api/graphql/resolvers"
	"github.com/src/user-auth-api/graphql/resolvers/mutations"
	userMutation "github.com/src/user-auth-api/graphql/resolvers/mutations/user"
	"github.com/src/user-auth-api/graphql/resolvers/query"
	userQuery "github.com/src/user-auth-api/graphql/resolvers/query/user"
	"github.com/src/user-auth-api/service"
	"github.com/src/user-auth-api/service/user"
)

var muxAdapter *gorillamux.GorillaMuxAdapter

type graphqlResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

func defaultTranslation() {
	directives.ValidateAddTranslation("email", " must be a valid email address")
}

func init() {
	r := mux.NewRouter()
	r.Use(injectDBCollection)
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
	// Ugly way of determining an appropriate response code
	os.Unsetenv("CLIENT_ERROR")
	if switchableAPIGatewayResponse, err = muxAdapter.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request)); err != nil {
		return apiGWResponse, err
	}

	apiGWResponse = *switchableAPIGatewayResponse.Version1()
	apiGWResponse.IsBase64Encoded = false
	apiGWResponse.StatusCode = http.StatusOK
	// Check for errors if this was a request to the /graphql endpoint
	if request.Path == "/graphql" {
		var gqlResponse graphqlResponse
		// Unmarshal GraphQL response and check for errors
		if err := json.Unmarshal([]byte(apiGWResponse.Body), &gqlResponse); err != nil {
			return apiGWResponse, fmt.Errorf("error parsing response: %w", err)
		}
		// Assign a response code based on the existence of errors.
		// TODO: Find better way to do this by not determining it based on an env variable.
		// With the lanbda proxying we're doing, there doesn't appear to be a clean way to communicate
		// service errors back out to the handler and then to API Gateway. The context doesn't appear to remain
		// populated during the before and after points of proxying the request.
		if len(gqlResponse.Errors) > 0 {
			clientError := os.Getenv("CLIENT_ERROR")
			switch clientError {
			case service.BAD_REQUEST:
				apiGWResponse.StatusCode = http.StatusBadRequest
			case service.NOT_FOUND:
				apiGWResponse.StatusCode = http.StatusNotFound
			default:
				apiGWResponse.StatusCode = http.StatusInternalServerError
			}
		}
	}
	return apiGWResponse, nil
}
