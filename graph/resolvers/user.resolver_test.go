package resolvers_test

import (
	"errors"
	"log"
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
	mockLastName       = "Smith"
	mockResponse       struct {
		User struct {
			UserID    string
			FirstName string
			LastName  string
			Email     string
			UserName  string
		}
	}
	mockEmail       = "mock_email@gmail.com"
	mockPassword    = "mockPassword123"
	mockUserName    = "mock_username"
	testAuthService *mocks.MockedUserService
)

func setup() {
	testAuthService = new(mocks.MockedUserService)
}

func TestQueryResolver_AuthenticateUser(t *testing.T) {
	q := `
	query DoLogin ($email: String!, $password: String!) { 
	  userLogin (params: {email: $email, password: $password}) {
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

		resolvers := resolvers.Resolvers{UserService: testAuthService}

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
		u := model.User{
			UserID:    mockID,
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

		_ = c.Post(
			q, &mockResponse,
			client.Var("email", mockEmail),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", mockEmail, mockPassword)

		// require.Equal(t, mockID, mockResponse.Data.ID)
		// require.Equal(t, mockName, mockResponse.Data.Name)
	})

	t.Run("should respond with an error when an invalid password is provided", func(t *testing.T) {
		setup()

		testAuthService.ErrorInvalidPassword = true
		resolvers := resolvers.Resolvers{UserService: testAuthService}

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))

		testAuthService.On(
			"AuthenticateUser",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(errInvalidPassword, nil)

		err := c.Post(
			q, &mockResponse,
			client.Var("email", mockEmail),
			client.Var("password", mockPassword),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "AuthenticateUser", 1)
		testAuthService.AssertCalled(t, "AuthenticateUser", mockEmail, mockPassword)

		require.EqualError(t, err, `[{"message":"invalid password","path":["userLogin"]}]`)
	})
}

func TestMutationResolver_CreateUser(t *testing.T) {
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
		Email:     "mock_email@gmail.com",
		FirstName: "John",
		LastName:  "Smith",
		UserName:  "mock_username",
		Password:  "mockPassword123",
	}

	t.Run("should create user", func(t *testing.T) {
		setup()

		resolvers := resolvers.Resolvers{UserService: testAuthService}

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
		u := model.User{
			UserID:    mockID,
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

		_ = c.Post(
			q, &mockResponse,
			client.Var("createUserInput", createUserInput),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "CreateUser", 1)
		testAuthService.AssertCalled(t, "CreateUser", createUserInput)

		log.Printf("Hwereeee\n")
		log.Printf("%+v\n", mockResponse)

		// require.Equal(t, mockID, mockResponse.Data.ID)
		// require.Equal(t, mockName, mockResponse.Data.Name)
	})

	t.Run("should respond with an error when a username is already taken", func(t *testing.T) {
		setup()

		testAuthService.ErrorUserAlreadyExists = true
		resolvers := resolvers.Resolvers{UserService: testAuthService}

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
		u := model.User{
			UserID:    mockID,
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

		err := c.Post(
			q, &mockResponse,
			client.Var("createUserInput", createUserInput),
		)

		testAuthService.AssertExpectations(t)
		testAuthService.AssertNumberOfCalls(t, "CreateUser", 1)
		testAuthService.AssertCalled(t, "CreateUser", createUserInput)

		require.EqualError(t, err, `[{"message":"user name already exists","path":["createUser"]}]`)
	})
}
