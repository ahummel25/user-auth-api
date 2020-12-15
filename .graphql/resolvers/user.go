package resolvers

import (
	"errors"

	"github.com/graphql-go/graphql"

	"github.com/src/user-auth-api/graphql/types"
)

var errInvalidPassword = errors.New("invalid password")

var userMap = map[string]string{
	"ahummel25": "Welcome123",
}

// UserResolver resolves the user.
var UserResolver = func(params graphql.ResolveParams) (interface{}, error) {
	userPassword := userMap[params.Args["username"].(string)]

	if ok := userPassword == params.Args["password"]; !ok {
		return nil, errInvalidPassword
	}

	user := types.UserType{
		ID:   1,
		Name: "Andrew",
	}

	return user, nil
}
