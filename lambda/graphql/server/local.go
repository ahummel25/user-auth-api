package server

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"

	"github.com/ahummel25/user-auth-api/db"
	"github.com/ahummel25/user-auth-api/graphql/directives"
	"github.com/ahummel25/user-auth-api/graphql/generated"
	"github.com/ahummel25/user-auth-api/graphql/resolvers"
	"github.com/ahummel25/user-auth-api/graphql/resolvers/mutations"
	userMutation "github.com/ahummel25/user-auth-api/graphql/resolvers/mutations/user"
	"github.com/ahummel25/user-auth-api/graphql/resolvers/query"
	userQuery "github.com/ahummel25/user-auth-api/graphql/resolvers/query/user"
	mainHandler "github.com/ahummel25/user-auth-api/lambda/graphql"
	"github.com/ahummel25/user-auth-api/service/user"
)

func StartLocalServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	mainHandler.DefaultTranslation()

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
	server := mainHandler.NewServer(schema)

	r.Handle("/graphiql", playground.Handler("GraphQL playground", "/graphql"))
	r.Handle("/apollo", playground.ApolloSandboxHandler("GraphQL Apollo playground", "/graphql"))
	r.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Setup DB context
		ctx, err := db.SetupDBContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r = r.WithContext(ctx)
		server.ServeHTTP(w, r)
	}))

	log.Printf("Server is running on http://localhost:%s/", port)
	log.Printf("GraphQL playground available at http://localhost:%s/graphiql", port)
	log.Printf("Apollo playground available at http://localhost:%s/apollo", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
