package resolvers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ahummel25/user-auth-api/graphql/directives"
	"github.com/ahummel25/user-auth-api/graphql/generated"
	"github.com/ahummel25/user-auth-api/graphql/model"
	"github.com/ahummel25/user-auth-api/graphql/resolvers"
	"github.com/ahummel25/user-auth-api/graphql/resolvers/mutations"
	userMutation "github.com/ahummel25/user-auth-api/graphql/resolvers/mutations/user"
	"github.com/ahummel25/user-auth-api/graphql/resolvers/query"
	userQuery "github.com/ahummel25/user-auth-api/graphql/resolvers/query/user"
	userMocks "github.com/ahummel25/user-auth-api/service/user/mocks"
)

var (
	errInvalidPassword    = errors.New("invalid password")
	errNoUserFound        = errors.New("no user found")
	errUserAlreadyExists  = errors.New("user name or email already exists")
	mockUserID            = "dfb8fe7f-56e4-47dc-b5bc-f6f0f524402b"
	mockFirstName         = "Test"
	mockLastName          = "User"
	mockEmail             = "mock_email@gmail.com"
	mockInvalidEmail      = "mock_email@gmail"
	mockPassword          = "mockPassword123"
	mockUserName          = "mock_username"
	mockRole              = model.RoleUser
	mockUserLoginResponse struct {
		AuthenticateUser struct {
			User struct {
				ID        string
				FirstName string
				LastName  string
				Email     string
				UserName  string
				Role      model.Role
			}
		}
	}

	mockCreateUserResponse struct {
		CreateUser struct {
			User struct {
				ID        string
				FirstName string
				LastName  string
				Email     string
				UserName  string
				Role      model.Role
			}
		}
	}

	mockDeleteUserResponse struct {
		DeleteUser bool
	}
)

var (
	loginQuery = `query Login($usernameOrEmail: String!, $password: String!) { 
	  authenticateUser(params: {usernameOrEmail: $usernameOrEmail, password: $password}) {
		user {
		  id
		  firstName
		  lastName
		  email
		  userName
		  role
		}
	  }
	}`

	createUser = `mutation CreateUser($newUserInput: NewUserInput!) {
	  createUser(user: $newUserInput) {
		user {
		  id
		  email
		  firstName
		  lastName
		  userName
		  role
		}
	  }
	}`

	deleteUser = `mutation DeleteUser($userID: ID!) {
		deleteUser(userID: $userID)
	}`
)

func setup() (*client.Client, *userMocks.MockAPI) {
	mockUserService := new(userMocks.MockAPI)
	mutationResolvers := mutations.MutationResolvers{
		Resolver: userMutation.Resolver{UserService: mockUserService},
	}
	queryResolvers := query.QueryResolvers{
		Resolver: userQuery.Resolver{UserService: mockUserService},
	}
	mockResolvers := resolvers.Services{
		MutationResolvers: mutationResolvers,
		QueryResolvers:    queryResolvers,
	}
	cfg := generated.Config{Resolvers: &mockResolvers}
	cfg.Directives.Binding = directives.Binding
	cfg.Directives.HasRole = directives.HasRole
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(cfg)))
	return c, mockUserService
}

func Test_AuthenticateUser_Success(t *testing.T) {
	c, mockUserService := setup()
	ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
		return assert.NotEmpty(t, actual)
	})
	u := model.User{
		ID:        mockUserID,
		FirstName: mockFirstName,
		LastName:  mockLastName,
		Email:     mockEmail,
		UserName:  mockUserName,
		Role:      mockRole,
	}

	uu := model.UserObject{User: &u}

	mockUserService.On(
		"AuthenticateUser",
		ctxArgMatcher,
		mockUserName,
		mockPassword,
	).Return(&uu, nil)

	c.MustPost(
		loginQuery, &mockUserLoginResponse,
		client.Var("usernameOrEmail", mockUserName),
		client.Var("password", mockPassword),
	)

	mockUserService.AssertExpectations(t)
	mockUserService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
	mockUserService.AssertCalled(t, "AuthenticateUser", ctxArgMatcher, mockUserName, mockPassword)

	require.Equal(t, mockUserID, mockUserLoginResponse.AuthenticateUser.User.ID)
	require.Equal(t, mockFirstName, mockUserLoginResponse.AuthenticateUser.User.FirstName)
	require.Equal(t, mockLastName, mockUserLoginResponse.AuthenticateUser.User.LastName)
	require.Equal(t, mockEmail, mockUserLoginResponse.AuthenticateUser.User.Email)
	require.Equal(t, mockUserName, mockUserLoginResponse.AuthenticateUser.User.UserName)
	require.Equal(t, mockRole, mockUserLoginResponse.AuthenticateUser.User.Role)
}

func Test_AuthenticateUser_Error(t *testing.T) {
	t.Run("should respond with an error when a valid user is not found", func(t *testing.T) {
		c, mockUserService := setup()
		ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
			return assert.NotEmpty(t, actual)
		})

		mockUserService.On(
			"AuthenticateUser",
			ctxArgMatcher,
			mockUserName,
			mockPassword,
		).Return(nil, errNoUserFound)

		err := c.Post(
			loginQuery, &mockUserLoginResponse,
			client.Var("usernameOrEmail", mockUserName),
			client.Var("password", mockPassword),
		)

		mockUserService.AssertExpectations(t)
		mockUserService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		mockUserService.AssertCalled(t, "AuthenticateUser", ctxArgMatcher, mockUserName, mockPassword)

		require.Empty(t, mockUserLoginResponse)
		require.EqualError(t, err, `[{"message":"`+errNoUserFound.Error()+`","path":["authenticateUser"]}]`)
	})

	t.Run("should respond with an error when an invalid password is provided", func(t *testing.T) {
		c, mockUserService := setup()
		ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
			return assert.NotEmpty(t, actual)
		})

		mockUserService.On(
			"AuthenticateUser",
			ctxArgMatcher,
			mockUserName,
			mockPassword,
		).Return(nil, errInvalidPassword)

		err := c.Post(
			loginQuery, &mockUserLoginResponse,
			client.Var("usernameOrEmail", mockUserName),
			client.Var("password", mockPassword),
		)

		mockUserService.AssertExpectations(t)
		mockUserService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		mockUserService.AssertCalled(t, "AuthenticateUser", ctxArgMatcher, mockUserName, mockPassword)

		require.Empty(t, mockUserLoginResponse)
		require.EqualError(t, err, `[{"message":"`+errInvalidPassword.Error()+`","path":["authenticateUser"]}]`)
	})
}

func Test_CreateUser_Success(t *testing.T) {
	c, mockUserService := setup()
	var newUserInput = model.NewUserInput{
		Email:     mockEmail,
		FirstName: mockFirstName,
		LastName:  mockLastName,
		UserName:  mockUserName,
		Password:  mockPassword,
	}
	ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
		return assert.NotEmpty(t, actual)
	})
	u := model.User{
		ID:        mockUserID,
		FirstName: mockFirstName,
		LastName:  mockLastName,
		Email:     mockEmail,
		UserName:  mockUserName,
		Role:      mockRole,
	}

	uu := model.UserObject{User: &u}

	mockUserService.On(
		"CreateUser",
		ctxArgMatcher,
		newUserInput,
	).Return(&uu, nil)

	c.MustPost(
		createUser, &mockCreateUserResponse,
		client.Var("newUserInput", newUserInput),
	)

	mockUserService.AssertExpectations(t)
	mockUserService.AssertNumberOfCalls(t, "CreateUser", 1)
	mockUserService.AssertCalled(t, "CreateUser", ctxArgMatcher, newUserInput)

	require.Equal(t, mockUserID, mockCreateUserResponse.CreateUser.User.ID)
	require.Equal(t, mockFirstName, mockCreateUserResponse.CreateUser.User.FirstName)
	require.Equal(t, mockLastName, mockCreateUserResponse.CreateUser.User.LastName)
	require.Equal(t, mockEmail, mockCreateUserResponse.CreateUser.User.Email)
	require.Equal(t, mockUserName, mockCreateUserResponse.CreateUser.User.UserName)
	require.Equal(t, mockRole, mockCreateUserResponse.CreateUser.User.Role)
}

func Test_CreateUser_Error(t *testing.T) {
	t.Run("should respond with an error when a username is already taken", func(t *testing.T) {
		c, mockUserService := setup()
		var newUserInput = model.NewUserInput{
			Email:     mockEmail,
			FirstName: mockFirstName,
			LastName:  mockLastName,
			UserName:  mockUserName,
			Password:  mockPassword,
		}
		ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
			return assert.NotEmpty(t, actual)
		})

		mockUserService.On(
			"CreateUser",
			ctxArgMatcher,
			newUserInput,
		).Return(nil, errUserAlreadyExists)

		err := c.Post(
			createUser, &mockCreateUserResponse,
			client.Var("newUserInput", newUserInput),
		)

		mockUserService.AssertExpectations(t)
		mockUserService.AssertNumberOfCalls(t, "CreateUser", 1)
		mockUserService.AssertCalled(t, "CreateUser", ctxArgMatcher, newUserInput)

		require.Empty(t, mockCreateUserResponse)
		require.EqualError(t, err, `[{"message":"`+errUserAlreadyExists.Error()+`","path":["createUser"]}]`)
	})
	t.Run("should respond with an error when the provided password is too short", func(t *testing.T) {
		c, mockUserService := setup()
		var newUserInput = model.NewUserInput{
			Email:     mockEmail,
			FirstName: mockFirstName,
			LastName:  mockLastName,
			UserName:  mockUserName,
			Password:  "123ABC",
		}
		err := c.Post(
			createUser, &mockCreateUserResponse,
			client.Var("newUserInput", newUserInput),
		)

		mockUserService.AssertExpectations(t)
		// CreateUser is not called since the binding directive catches the invalid password
		mockUserService.AssertNumberOfCalls(t, "CreateUser", 0)
		mockUserService.AssertNotCalled(t, "CreateUser")

		require.Empty(t, mockCreateUserResponse)
		require.EqualError(t, err, `[{"message":"password must be at least 8 characters in length","path":["createUser","user","password"]}]`)
	})
	t.Run("should respond with an error when an invalid email provided", func(t *testing.T) {
		c, mockUserService := setup()
		var newUserInput = model.NewUserInput{
			Email:     mockInvalidEmail,
			FirstName: mockFirstName,
			LastName:  mockLastName,
			UserName:  mockUserName,
			Password:  mockPassword,
		}
		err := c.Post(
			createUser, &mockCreateUserResponse,
			client.Var("newUserInput", newUserInput),
		)

		mockUserService.AssertExpectations(t)
		// CreateUser is not called since the binding directive catches the invalid email format
		mockUserService.AssertNumberOfCalls(t, "CreateUser", 0)
		mockUserService.AssertNotCalled(t, "CreateUser")

		require.Empty(t, mockCreateUserResponse)
		require.EqualError(t, err, `[{"message":"email must be a valid email address","path":["createUser","user","email"]}]`)
	})
}

func Test_DeleteUser_Success(t *testing.T) {
	c, mockUserService := setup()
	ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
		return assert.NotEmpty(t, actual)
	})

	mockUserService.On(
		"DeleteUser",
		ctxArgMatcher,
		mockUserID,
	).Return(true, nil)

	c.MustPost(
		deleteUser, &mockDeleteUserResponse,
		client.Var("userID", mockUserID),
	)

	mockUserService.AssertExpectations(t)
	mockUserService.AssertNumberOfCalls(t, "DeleteUser", 1)
	mockUserService.AssertCalled(t, "DeleteUser", ctxArgMatcher, mockUserID)
	require.Equal(t, true, mockDeleteUserResponse.DeleteUser)
}

func Test_DeleteUser_Error(t *testing.T) {
	t.Run("should respond with an error when a user trying to be deleted does not exist", func(t *testing.T) {
		c, mockUserService := setup()
		ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
			return assert.NotEmpty(t, actual)
		})

		mockUserService.On(
			"DeleteUser",
			ctxArgMatcher,
			mockUserID,
		).Return(false, errNoUserFound)

		err := c.Post(
			deleteUser, &mockDeleteUserResponse,
			client.Var("userID", mockUserID),
		)

		mockUserService.AssertExpectations(t)
		mockUserService.AssertNumberOfCalls(t, "DeleteUser", 1)
		mockUserService.AssertCalled(t, "DeleteUser", ctxArgMatcher, mockUserID)
		require.Empty(t, mockDeleteUserResponse)
		require.Equal(t, false, mockDeleteUserResponse.DeleteUser)
		require.EqualError(t, err, `[{"message":"`+errNoUserFound.Error()+`","path":["deleteUser"]}]`)
	})
}
