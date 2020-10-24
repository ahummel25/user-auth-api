package types

import (
	"github.com/graphql-go/graphql"
)

// UserType represents the resolved type from the user query.
type UserType struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// UserObject represents a user object.
var UserObject = graphql.NewObject(graphql.ObjectConfig{
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
