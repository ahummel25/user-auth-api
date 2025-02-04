package graphql

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/vektah/gqlparser/v2/ast"

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

func DefaultTranslation() {
	directives.ValidateAddTranslation("email", " must be a valid email address")
}

// NewServer creates a GraphQL server with common configurations
func NewServer(es graphql.ExecutableSchema) *handler.Server {
	srv := handler.New(es)

	// Add transports in order of preference
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	// Set up query cache
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	// Add extensions
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

func init() {
	r := mux.NewRouter()
	DefaultTranslation()

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
	server := NewServer(schema)

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
