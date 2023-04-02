package resolvers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/src/user-auth-api/graphql/generated"
	"github.com/src/user-auth-api/graphql/model"
	"github.com/src/user-auth-api/graphql/resolvers"
	userMocks "github.com/src/user-auth-api/service/user/mocks"
)

var (
	errInvalidPassword       = errors.New("invalid password")
	errNoUserFound           = errors.New("no user found")
	errUserNameAlreadyExists = errors.New("user name already exists")
	mockUserID               = "1"
	mockFirstName            = "John"
	mockLastName             = "Smith"
	mockEmail                = "mock_email@gmail.com"
	mockPassword             = "mockPassword123"
	mockUserName             = "mock_username"
	mockRole                 = model.RoleUser
	mockUserLoginResponse    struct {
		AuthenticateUser struct {
			User struct {
				UserID    string
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
				UserID    string
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
	mockResolvers   resolvers.Services
	testAuthService *userMocks.API
)

var (
	loginQuery = `query Login($username: String!, $password: String!) { 
	  authenticateUser(params: {username: $username, password: $password}) {
		user {
		  firstName
		  lastName
		  email
		  userName
		  userID
		  role
		}
	  }
	}`

	createUser = `mutation CreateUser($createUserInput: CreateUserInput!) {
	  createUser(user: $createUserInput) {
		user {
		  userID
		  email
		  firstName
		  lastName
		  userName
		  role
		}
	  }
	}`

	deleteUser = `
	mutation DeleteUser($deleteUserInput: DeleteUserInput!) {
		deleteUser(user: $deleteUserInput)
	}`
)

func setup() {
	testAuthService = new(userMocks.API)
	mockResolvers = resolvers.Services{UserService: testAuthService}
}

func Test_AuthenticateUser_Success(t *testing.T) {
	setup()

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &mockResolvers})))
	ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
		return assert.NotEmpty(t, actual)
	})
	u := model.User{
		UserID:    mockUserID,
		FirstName: mockFirstName,
		LastName:  mockLastName,
		Email:     mockEmail,
		UserName:  mockUserName,
		Role:      mockRole,
	}

	uu := model.UserObject{User: &u}

	testAuthService.On(
		"AuthenticateUser",
		ctxArgMatcher,
		mockUserName,
		mockPassword,
	).Return(&uu, nil)

	c.MustPost(
		loginQuery, &mockUserLoginResponse,
		client.Var("username", mockUserName),
		client.Var("password", mockPassword),
	)

	testAuthService.AssertExpectations(t)
	testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
	testAuthService.AssertCalled(t, "AuthenticateUser", ctxArgMatcher, mockUserName, mockPassword)

	require.Equal(t, mockUserID, mockUserLoginResponse.AuthenticateUser.User.UserID)
	require.Equal(t, mockFirstName, mockUserLoginResponse.AuthenticateUser.User.FirstName)
	require.Equal(t, mockLastName, mockUserLoginResponse.AuthenticateUser.User.LastName)
	require.Equal(t, mockEmail, mockUserLoginResponse.AuthenticateUser.User.Email)
	require.Equal(t, mockUserName, mockUserLoginResponse.AuthenticateUser.User.UserName)
	require.Equal(t, mockRole, mockUserLoginResponse.AuthenticateUser.User.Role)
}

func Test_AuthenticateUser_Error(t *testing.T) {
	t.Run("should respond with an error when a valid user is not found", func(t *testing.T) {
		setup()

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &mockResolvers})))
		ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
			return assert.NotEmpty(t, actual)
		})

		testAuthService.On(
			"AuthenticateUser",
			ctxArgMatcher,
			mockUserName,
			mockPassword,
		).Return(nil, errNoUserFound)

		err := c.Post(
			loginQuery, &mockUserLoginResponse,
			client.Var("username", mockUserName),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", ctxArgMatcher, mockUserName, mockPassword)

		require.Empty(t, mockUserLoginResponse)
		require.EqualError(t, err, `[{"message":"`+errNoUserFound.Error()+`","path":["authenticateUser"]}]`)
	})

	t.Run("should respond with an error when an invalid password is provided", func(t *testing.T) {
		setup()

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &mockResolvers})))
		ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
			return assert.NotEmpty(t, actual)
		})

		testAuthService.On(
			"AuthenticateUser",
			ctxArgMatcher,
			mockUserName,
			mockPassword,
		).Return(nil, errInvalidPassword)

		err := c.Post(
			loginQuery, &mockUserLoginResponse,
			client.Var("username", mockUserName),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", ctxArgMatcher, mockUserName, mockPassword)

		require.Empty(t, mockUserLoginResponse)
		require.EqualError(t, err, `[{"message":"`+errInvalidPassword.Error()+`","path":["authenticateUser"]}]`)
	})
}

func Test_CreateUser_Success(t *testing.T) {
	setup()
	var createUserInput = model.CreateUserInput{
		Email:     mockEmail,
		FirstName: mockFirstName,
		LastName:  mockLastName,
		UserName:  mockUserName,
		Password:  mockPassword,
	}
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

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(cfg)))
	ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
		return assert.NotEmpty(t, actual)
	})
	u := model.User{
		UserID:    mockUserID,
		FirstName: mockFirstName,
		LastName:  mockLastName,
		Email:     mockEmail,
		UserName:  mockUserName,
		Role:      mockRole,
	}

	uu := model.UserObject{User: &u}

	testAuthService.On(
		"CreateUser",
		ctxArgMatcher,
		createUserInput,
	).Return(&uu, nil)

	c.MustPost(
		createUser, &mockCreateUserResponse,
		client.Var("createUserInput", createUserInput),
	)

	testAuthService.AssertExpectations(t)
	testAuthService.AssertNumberOfCalls(t, "CreateUser", 1)
	testAuthService.AssertCalled(t, "CreateUser", ctxArgMatcher, createUserInput)

	require.Equal(t, mockUserID, mockCreateUserResponse.CreateUser.User.UserID)
	require.Equal(t, mockFirstName, mockCreateUserResponse.CreateUser.User.FirstName)
	require.Equal(t, mockLastName, mockCreateUserResponse.CreateUser.User.LastName)
	require.Equal(t, mockEmail, mockCreateUserResponse.CreateUser.User.Email)
	require.Equal(t, mockUserName, mockCreateUserResponse.CreateUser.User.UserName)
	require.Equal(t, mockRole, mockCreateUserResponse.CreateUser.User.Role)
}

func Test_CreateUser_Error(t *testing.T) {
	t.Run("should respond with an error when a username is already taken", func(t *testing.T) {
		setup()
		var createUserInput = model.CreateUserInput{
			Email:     mockEmail,
			FirstName: mockFirstName,
			LastName:  mockLastName,
			UserName:  mockUserName,
			Password:  mockPassword,
		}
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

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(cfg)))
		ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
			return assert.NotEmpty(t, actual)
		})

		testAuthService.On(
			"CreateUser",
			ctxArgMatcher,
			createUserInput,
		).Return(nil, errUserNameAlreadyExists)

		err := c.Post(
			createUser, &mockCreateUserResponse,
			client.Var("createUserInput", createUserInput),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "CreateUser", 1)
		testAuthService.AssertCalled(t, "CreateUser", ctxArgMatcher, createUserInput)

		require.Empty(t, mockCreateUserResponse)
		require.EqualError(t, err, `[{"message":"`+errUserNameAlreadyExists.Error()+`","path":["createUser"]}]`)
	})
}

func Test_DeleteUser_Success(t *testing.T) {
	setup()
	var deleteUserInput = model.DeleteUserInput{
		Email:     mockEmail,
		FirstName: mockFirstName,
		LastName:  mockLastName,
		UserID:    mockUserID,
		UserName:  mockUserName,
	}
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

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(cfg)))
	ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
		return assert.NotEmpty(t, actual)
	})

	testAuthService.On(
		"DeleteUser",
		ctxArgMatcher,
		deleteUserInput,
	).Return(true, nil)

	c.MustPost(
		deleteUser, &mockDeleteUserResponse,
		client.Var("deleteUserInput", deleteUserInput),
	)

	testAuthService.AssertExpectations(t)
	testAuthService.AssertNumberOfCalls(t, "DeleteUser", 1)
	testAuthService.AssertCalled(t, "DeleteUser", ctxArgMatcher, deleteUserInput)
	require.Equal(t, true, mockDeleteUserResponse.DeleteUser)
}

func Test_DeleteUser_Error(t *testing.T) {
	t.Run("should respond with an error when a user trying to be deleted does not exist", func(t *testing.T) {
		setup()
		var deleteUserInput = model.DeleteUserInput{
			Email:     mockEmail,
			FirstName: mockFirstName,
			LastName:  mockLastName,
			UserID:    mockUserID,
			UserName:  mockUserName,
		}
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

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(cfg)))
		ctxArgMatcher := mock.MatchedBy(func(actual context.Context) bool {
			return assert.NotEmpty(t, actual)
		})

		testAuthService.On(
			"DeleteUser",
			ctxArgMatcher,
			deleteUserInput,
		).Return(false, errNoUserFound)

		err := c.Post(
			deleteUser, &mockDeleteUserResponse,
			client.Var("deleteUserInput", deleteUserInput),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "DeleteUser", 1)
		testAuthService.AssertCalled(t, "DeleteUser", ctxArgMatcher, deleteUserInput)
		require.Empty(t, mockDeleteUserResponse)
		require.Equal(t, false, mockDeleteUserResponse.DeleteUser)
		require.EqualError(t, err, `[{"message":"`+errNoUserFound.Error()+`","path":["deleteUser"]}]`)
	})
}
