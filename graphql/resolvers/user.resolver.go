package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/src/user-auth-api/graphql/model"
)

// var c generated.Config

// func init() {
// 	userService := services.NewUserService()

// 	r := Resolvers{
// 		UserService: userService,
// 	}

// 	c = generated.Config{Resolvers: &r}
// }
func (r *mutationResolver) CreateUser(ctx context.Context, params model.CreateUserInput) (*model.UserObject, error) {

	// c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
	// 	fmt.Println("HasRole called!")
	// 	return next(ctx)
	// }

	log.Printf("CreateUser reolver called\n")

	return r.UserService.CreateUser(params)
}

func (r *queryResolver) AuthenticateUser(ctx context.Context, params model.AuthParams) (*model.UserObject, error) {
	return r.UserService.AuthenticateUser(params.Email, params.Password)
}
