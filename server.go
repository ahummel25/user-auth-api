package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/src/user-auth-api/graph/generated"
	"github.com/src/user-auth-api/graph/resolvers"
	"github.com/src/user-auth-api/services"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	userService := services.NewUserService()

	resolvers := resolvers.Resolvers{
		UserService: userService,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers}))

	http.Handle("/graphiql", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/graphiql for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
