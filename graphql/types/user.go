package types

import (
	"github.com/graphql-go/graphql"
)

var UserMap = map[string]string{
	"ahummel25": "Welcome123",
}

// UserType represents the resolved type from the user query.
type UserType struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// User represents a user account.
var User = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})
