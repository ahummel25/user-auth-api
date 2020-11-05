package resolvers_test

import (
	"errors"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/src/user-auth-api/graph/generated"
	"github.com/src/user-auth-api/graph/model"
	"github.com/src/user-auth-api/graph/resolvers"
	"github.com/src/user-auth-api/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errInvalidPassword = errors.New("invalid password")
	mockID             = "1"
	mockFirstName      = "John"
	mockResponse       struct {
		Data struct {
			UserID string
			Name   string
		}
	}
	mockUsername    = "ahummel25"
	mockPassword    = "Welcome123"
	testAuthService *mocks.MockedAuthService
)

func setup() {
	testAuthService = new(mocks.MockedAuthService)
}

func TestQueryResolver_AuthenticateUser(t *testing.T) {
	q := `
	query DoLogin ($username: String!, $password: String!) { 
	  userLogin (params: {username: $username, password: $password}) {
		  userID
		  name
	  }
	}
  `

	t.Run("should authenticate user correctly", func(t *testing.T) {
		setup()
		resolvers := resolvers.Resolvers{AuthService: testAuthService}

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
		u := model.User{UserID: mockID, FirstName: mockFirstName}

		testAuthService.On(
			"AuthenticateUser",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&u, nil)

		_ = c.Post(
			q, &mockResponse,
			client.Var("username", mockUsername),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", mockUsername, mockPassword)

		// require.Equal(t, mockID, mockResponse.Data.ID)
		// require.Equal(t, mockName, mockResponse.Data.Name)
	})

	t.Run("should respond with an error when an invalid password is provided", func(t *testing.T) {
		setup()
		testAuthService.ErrorInvalidPassword = true
		resolvers := resolvers.Resolvers{AuthService: testAuthService}

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))

		testAuthService.On(
			"AuthenticateUser",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(errInvalidPassword, nil)

		err := c.Post(
			q, &mockResponse,
			client.Var("username", mockUsername),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", mockUsername, mockPassword)

		require.EqualError(t, err, `[{"message":"invalid password","path":["userLogin"]}]`)
	})
}
