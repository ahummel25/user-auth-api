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
	errNoUserFound           = errors.New("no user found!")
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

	mockDeleteUserResponse struct {
		DeleteUser string
	}
	mockResolvers   resolvers.Services
	testAuthService *mocks.MockedUserService
)

func setup() {
	testAuthService = new(mocks.MockedUserService)
	mockResolvers = resolvers.Services{UserService: testAuthService}
}

func TestQueryResolver_AuthenticateUser(t *testing.T) {
	q := `
	query DoLogin ($username: String!, $password: String!) { 
	  authenticateUser(params: {username: $username, password: $password}) {
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
			client.Var("username", mockUserName),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", mockUserName, mockPassword)

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
			client.Var("username", mockUserName),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", mockUserName, mockPassword)

		require.Empty(t, mockUserLoginResponse)
		require.EqualError(t, err, `[{"message":"`+errNoUserFound.Error()+`","path":["authenticateUser"]}]`)
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
			client.Var("username", mockUserName),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", mockUserName, mockPassword)

		require.Empty(t, mockUserLoginResponse)
		require.EqualError(t, err, `[{"message":"`+errInvalidPassword.Error()+`","path":["authenticateUser"]}]`)
	})
}

func TestMutationResolver_CreateUser(t *testing.T) {
	cfg := generated.Config{
		Resolvers: &mockResolvers,
		Directives: generated.DirectiveRoot{
			HasRole: func(
				ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role, action model.Action,
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
		require.EqualError(t, err, `[{"message":"`+errUserNameAlreadyExists.Error()+`","path":["createUser"]}]`)
	})
}

func TestMutationResolver_DeleteUser(t *testing.T) {
	cfg := generated.Config{
		Resolvers: &mockResolvers,
		Directives: generated.DirectiveRoot{
			HasRole: func(
				ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role, action model.Action,
			) (res interface{}, err error) {
				return next(ctx)
			},
		},
	}

	q := `
	mutation DeleteUser($deleteUserInput: DeleteUserInput!) {
		deleteUser(user: $deleteUserInput)
	}
  `

	var deleteUserInput = model.DeleteUserInput{
		Email:     mockEmail,
		FirstName: mockFirstName,
		LastName:  mockLastName,
		UserID:    mockUserID,
		UserName:  mockUserName,
	}

	t.Run("should delete the user", func(t *testing.T) {
		setup()

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(cfg)))
		r := deleteUserInput.UserName + " successfully deleted"

		testAuthService.On(
			"DeleteUser",
			mock.AnythingOfType("DeleteUserInput"),
		).Return(r, nil)

		c.MustPost(
			q, &mockDeleteUserResponse,
			client.Var("deleteUserInput", deleteUserInput),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "DeleteUser", 1)
		testAuthService.AssertCalled(t, "DeleteUser", deleteUserInput)

		require.Equal(t, r, mockDeleteUserResponse.DeleteUser)
	})

	t.Run("should respond with an error when a user trying to be deleted does not exist", func(t *testing.T) {
		setup()

		testAuthService.ErrorNoUserFound = true

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(cfg)))

		testAuthService.On(
			"DeleteUser",
			mock.AnythingOfType("DeleteUserInput"),
		).Return(nil, errNoUserFound)

		err := c.Post(
			q, &mockDeleteUserResponse,
			client.Var("deleteUserInput", deleteUserInput),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "DeleteUser", 1)
		testAuthService.AssertCalled(t, "DeleteUser", deleteUserInput)

		require.Empty(t, mockDeleteUserResponse)
		require.EqualError(t, err, `[{"message":"`+errNoUserFound.Error()+`","path":["deleteUser"]}]`)
	})
}
