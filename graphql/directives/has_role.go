package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/ahummel25/user-auth-api/graphql/model"
)

func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role, action model.Action,
) (res interface{}, err error) {
	fc := graphql.GetFieldContext(ctx).Args

	switch action.String() {
	case model.ActionCreateUser.String():
		_, ok := fc["user"].(model.NewUserInput)
		if !ok {
			return nil, fmt.Errorf("invalid user")
		}
	case model.ActionDeleteUser.String():
		// The userID is directly available in the args
		_, ok := fc["userID"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid userID")
		}
		// Here you should implement the logic to check if the current user has the required role
		// This might involve getting the current user from the context and checking their role
		// For example:
		// currentUser := auth.GetUserFromContext(ctx)
		// if currentUser.Role != role {
		//     return nil, fmt.Errorf("user does not have the required role")
		// }
	}
	return next(ctx)
}
