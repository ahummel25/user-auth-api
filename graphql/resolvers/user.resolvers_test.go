package resolvers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/src/user-auth-api/graphql/generated"
	"github.com/src/user-auth-api/graphql/model"
	"github.com/src/user-auth-api/graphql/resolvers"
	"github.com/src/user-auth-api/mocks"
)

var (
	errInvalidPassword       = errors.New("invalid password")
	errNoUserFound           = errors.New("no user found by that email address")
	errUserNameAlreadyExists = errors.New("user name already exists")
	mockUserID               = "1"
	mockFirstName            = "John"
	mockLastName             = "Smith"
	mockEmail                = "mock_email@gmail.com"
	mockPassword             = "mockPassword123"
	mockUserName             = "mock_username"
	mockUserLoginResponse    struct {
		AuthenticateUser struct {
			User struct {
				UserID    string
				FirstName string
				LastName  string
				Email     string
				UserName  string
			}
		}
	}

	mockCreateUserResponse struct {
		CreateUser struct {
			User struct {
				UserID    string
				FirstName string
				LastName  string
				Email     string
				UserName  string
			}
		}
	}
	mockResolvers   resolvers.Resolvers
	testAuthService *mocks.MockedUserService
)

func setup() {
	testAuthService = new(mocks.MockedUserService)
	mockResolvers = resolvers.Resolvers{UserService: testAuthService}
}

func TestQueryResolver_AuthenticateUser(t *testing.T) {
	q := `
	query DoLogin ($email: String!, $password: String!) { 
	  authenticateUser(params: {email: $email, password: $password}) {
		user {
		  firstName
		  lastName
		  email
		  userName
		  userID
		}
	  }
	}
  `

	t.Run("should authenticate user correctly", func(t *testing.T) {
		setup()

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &mockResolvers})))
		u := model.User{
			UserID:    mockUserID,
			FirstName: mockFirstName,
			LastName:  mockLastName,
			Email:     mockEmail,
			UserName:  mockUserName,
		}

		uu := model.UserObject{User: &u}

		testAuthService.On(
			"AuthenticateUser",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&uu, nil)

		c.MustPost(
			q, &mockUserLoginResponse,
			client.Var("email", mockEmail),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", mockEmail, mockPassword)

		require.Equal(t, mockUserID, mockUserLoginResponse.AuthenticateUser.User.UserID)
		require.Equal(t, mockFirstName, mockUserLoginResponse.AuthenticateUser.User.FirstName)
		require.Equal(t, mockLastName, mockUserLoginResponse.AuthenticateUser.User.LastName)
		require.Equal(t, mockEmail, mockUserLoginResponse.AuthenticateUser.User.Email)
		require.Equal(t, mockUserName, mockUserLoginResponse.AuthenticateUser.User.UserName)
	})

	t.Run("should respond with an error when a valid user is not found", func(t *testing.T) {
		setup()

		testAuthService.ErrorNoUserFound = true

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &mockResolvers})))

		testAuthService.On(
			"AuthenticateUser",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, errNoUserFound)

		err := c.Post(
			q, &mockUserLoginResponse,
			client.Var("email", mockEmail),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", mockEmail, mockPassword)

		require.Empty(t, mockUserLoginResponse)
		require.EqualError(t, err, `[{"message":"no user found by that email address","path":["authenticateUser"]}]`)
	})

	t.Run("should respond with an error when an invalid password is provided", func(t *testing.T) {
		setup()

		testAuthService.ErrorInvalidPassword = true

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &mockResolvers})))

		testAuthService.On(
			"AuthenticateUser",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, errInvalidPassword)

		err := c.Post(
			q, &mockUserLoginResponse,
			client.Var("email", mockEmail),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", mockEmail, mockPassword)

		require.Empty(t, mockUserLoginResponse)
		require.EqualError(t, err, `[{"message":"invalid password","path":["authenticateUser"]}]`)
	})
}

func TestMutationResolver_CreateUser(t *testing.T) {
	cfg := generated.Config{
		Resolvers: &mockResolvers,
		Directives: generated.DirectiveRoot{
			HasRole: func(
				ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role,
			) (res interface{}, err error) {
				return next(ctx)
			},
		},
	}

	q := `
	mutation CreateUser($createUserInput: CreateUserInput!) {
		createUser(user: $createUserInput) {
		  user {
			userID
			email
			firstName
			lastName
			userName
		  }
		}
	  }
  `

	var createUserInput = model.CreateUserInput{
		Email:     mockEmail,
		FirstName: mockFirstName,
		LastName:  mockLastName,
		UserName:  mockUserName,
		Password:  mockPassword,
	}

	t.Run("should create user", func(t *testing.T) {
		setup()

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(cfg)))
		u := model.User{
			UserID:    mockUserID,
			FirstName: mockFirstName,
			LastName:  mockLastName,
			Email:     mockEmail,
			UserName:  mockUserName,
		}

		uu := model.UserObject{User: &u}

		testAuthService.On(
			"CreateUser",
			mock.AnythingOfType("CreateUserInput"),
		).Return(&uu, nil)

		c.MustPost(
			q, &mockCreateUserResponse,
			client.Var("createUserInput", createUserInput),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "CreateUser", 1)
		testAuthService.AssertCalled(t, "CreateUser", createUserInput)

		require.Equal(t, mockUserID, mockCreateUserResponse.CreateUser.User.UserID)
		require.Equal(t, mockFirstName, mockCreateUserResponse.CreateUser.User.FirstName)
		require.Equal(t, mockLastName, mockCreateUserResponse.CreateUser.User.LastName)
		require.Equal(t, mockEmail, mockCreateUserResponse.CreateUser.User.Email)
		require.Equal(t, mockUserName, mockCreateUserResponse.CreateUser.User.UserName)
	})

	t.Run("should respond with an error when a username is already taken", func(t *testing.T) {
		setup()

		testAuthService.ErrorUserAlreadyExists = true

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(cfg)))

		testAuthService.On(
			"CreateUser",
			mock.AnythingOfType("CreateUserInput"),
		).Return(nil, errUserNameAlreadyExists)

		err := c.Post(
			q, &mockCreateUserResponse,
			client.Var("createUserInput", createUserInput),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "CreateUser", 1)
		testAuthService.AssertCalled(t, "CreateUser", createUserInput)

		require.Empty(t, mockCreateUserResponse)
		require.EqualError(t, err, `[{"message":"user name already exists","path":["createUser"]}]`)
	})
}
