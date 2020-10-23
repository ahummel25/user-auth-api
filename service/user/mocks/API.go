// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/src/user-auth-api/graphql/model"
	mock "github.com/stretchr/testify/mock"
)

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

// AuthenticateUser provides a mock function with given fields: ctx, username, password
func (_m *API) AuthenticateUser(ctx context.Context, username string, password string) (*model.UserObject, error) {
	ret := _m.Called(ctx, username, password)

	var r0 *model.UserObject
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*model.UserObject, error)); ok {
		return rf(ctx, username, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.UserObject); ok {
		r0 = rf(ctx, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserObject)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, params
func (_m *API) CreateUser(ctx context.Context, params model.NewUserInput) (*model.UserObject, error) {
	ret := _m.Called(ctx, params)

	var r0 *model.UserObject
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.NewUserInput) (*model.UserObject, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.NewUserInput) *model.UserObject); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserObject)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.NewUserInput) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: ctx, userID
func (_m *API) DeleteUser(ctx context.Context, userID string) (bool, error) {
	ret := _m.Called(ctx, userID)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAPI interface {
	mock.TestingT
	Cleanup(func())
}

// NewAPI creates a new instance of API. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAPI(t mockConstructorTestingTNewAPI) *API {
	mock := &API{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
