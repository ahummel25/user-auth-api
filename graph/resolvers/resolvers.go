package resolvers

import "github.com/src/user-auth-api/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolvers struct {
	AuthService services.AuthService
}
