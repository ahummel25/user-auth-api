package directives

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/src/user-auth-api/graphql/model"
)

func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role, action model.Action,
) (res interface{}, err error) {
	fc := graphql.GetFieldContext(ctx).Args

	switch action.String() {
	case model.ActionCreateUser.String():
		log.Printf("%+v\n", fc["user"].(model.NewUserInput).Email)
	case model.ActionDeleteUser.String():
		log.Printf("%+v\n", fc)
	}
	return next(ctx)
}
