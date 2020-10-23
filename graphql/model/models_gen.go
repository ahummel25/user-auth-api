// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

// The input needed to authenticate a user.
type AuthParams struct {
	// The user's username or email address
	UsernameOrEmail string `json:"usernameOrEmail"`
	// The user's password
	Password string `json:"password"`
}

type Mutation struct {
}

// The input required to create a new user.
type NewUserInput struct {
	// The user's e-mail address
	Email string `json:"email"`
	// The user's first name
	FirstName string `json:"firstName"`
	// The user's last name
	LastName string `json:"lastName"`
	// The user's username
	UserName string `json:"userName"`
	// The user's role
	Role *Role `json:"role,omitempty"`
	// The user's password
	Password string `json:"password"`
}

type Query struct {
}

// An object representing an individual user.
type User struct {
	// The user's unique user ID
	ID string `json:"id"`
	// The user's e-mail address
	Email string `json:"email"`
	// The user's first name
	FirstName string `json:"firstName"`
	// The user's last name
	LastName string `json:"lastName"`
	// The user's username
	UserName string `json:"userName"`
	// The user's role
	Role Role `json:"role"`
}

type UserObject struct {
	// The user object pertaining to the given user.
	User *User `json:"user"`
}

type Action string

const (
	// Create User Action
	ActionCreateUser Action = "CREATE_USER"
	// Delete User Action
	ActionDeleteUser Action = "DELETE_USER"
)

var AllAction = []Action{
	ActionCreateUser,
	ActionDeleteUser,
}

func (e Action) IsValid() bool {
	switch e {
	case ActionCreateUser, ActionDeleteUser:
		return true
	}
	return false
}

func (e Action) String() string {
	return string(e)
}

func (e *Action) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Action(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Action", str)
	}
	return nil
}

func (e Action) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Role string

const (
	// ADMIN Role
	RoleAdmin Role = "ADMIN"
	// USER Role
	RoleUser Role = "USER"
)

var AllRole = []Role{
	RoleAdmin,
	RoleUser,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
