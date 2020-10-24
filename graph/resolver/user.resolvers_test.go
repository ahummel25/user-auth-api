package resolver_test

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/src/user-auth-api/graph/generated"
	"github.com/src/user-auth-api/graph/model"
	"github.com/src/user-auth-api/graph/resolver"
	"github.com/src/user-auth-api/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	mockID       = "1"
	mockName     = "John Smith"
	mockResponse struct {
		Data struct {
			ID   string
			Name string
		}
	}
	mockUsername = "mockUsername123"
	mockPassword = "mockPassword123"
)

func TestQueryResolver_AuthenticateUser(t *testing.T) {
	t.Run("should authenticate user correctly", func(t *testing.T) {
		testAuthService := new(mocks.MockedAuthService)
		resolvers := resolver.Resolver{AuthService: testAuthService}

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
		u := model.User{ID: mockID, Name: mockName}

		testAuthService.On(
			"AuthenticateUser",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&u)

		q := `
		  query DoLogin ($username: String!, $password: String!) { 
			userLogin (params: {username: $username, password: $password}) {
				id
				name
			}
		  }
		`

		c.MustPost(
			q, &mockResponse,
			client.Var("username", mockUsername),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertCalled(t, "AuthenticateUser", mockUsername, mockPassword)

		require.Equal(t, mockID, mockResponse.Data.ID)
		require.Equal(t, mockName, mockResponse.Data.Name)
	})
}
