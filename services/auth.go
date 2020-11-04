package services

import (
	"database/sql"
	"errors"
	"log"

	dbHelper "github.com/src/user-auth-api/db"
	"github.com/src/user-auth-api/graph/model"
)

// AuthService contains signatures for any auth functions.
type AuthService interface {
	AuthenticateUser(username string, password string) (*model.User, error)
}

type authService struct{}

var (
	errInvalidPassword  = errors.New("invalid password")
	errUserDoesNotExist = errors.New("user does not exist")
)

var mockUserDB = map[string]string{
	"ahummel25": "Welcome123",
}

// NewAuthService returns a pointer to a new auth service.
func NewAuthService() *authService {
	return &authService{}
}

// AuthenticateUser authenticates the user.
func (a *authService) AuthenticateUser(username string, password string) (*model.User, error) {
	var (
		db  *sql.DB
		err error
	)
	if db, err = dbHelper.ConnectToDB(); err != nil {
		return nil, err
	}

	log.Printf("%v\n", db)

	defer db.Close()

	if mockUserDB[username] == "" {
		return nil, errUserDoesNotExist
	}

	dbUserPassword := mockUserDB[username]

	if dbUserPassword != password {
		return nil, errInvalidPassword
	}

	id := "1"
	name := "Andrew"

	user := &model.User{
		UserID: id,
		Name:   name,
	}

	return user, nil
}
