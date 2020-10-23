package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/ahummel25/user-auth-api/graphql/generated"
	"github.com/ahummel25/user-auth-api/graphql/model"
	"github.com/ahummel25/user-auth-api/graphql/resolvers"
	"github.com/ahummel25/user-auth-api/graphql/resolvers/mutations"
	userMutation "github.com/ahummel25/user-auth-api/graphql/resolvers/mutations/user"
	"github.com/ahummel25/user-auth-api/graphql/resolvers/query"
	userQuery "github.com/ahummel25/user-auth-api/graphql/resolvers/query/user"
	"github.com/ahummel25/user-auth-api/service/user"
)

const defaultPort = "8080"

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/404.html")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

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

	c := generated.Config{
		Resolvers: &resolvers,
		Directives: generated.DirectiveRoot{
			HasRole: func(
				ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role, action model.Action,
			) (res interface{}, err error) {
				fc := graphql.GetFieldContext(ctx).Args
				log.Printf("Action: %s\n", action.String())

				switch action.String() {
				case model.ActionCreateUser.String():
					log.Printf("%+v\n", fc["user"].(model.NewUserInput).Email)
				case model.ActionDeleteUser.String():
					log.Printf("%+v\n", fc)
				}
				return next(ctx)
			},
		},
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	http.Handle("/graphiql", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)
	http.HandleFunc("/", notFoundHandler)

	log.Printf("connect to http://localhost:%s/graphiql for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
