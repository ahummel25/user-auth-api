package resolver

import "github.com/src/user-auth-api/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthService services.AuthService
}
