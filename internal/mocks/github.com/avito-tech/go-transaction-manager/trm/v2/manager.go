// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	trm "github.com/avito-tech/go-transaction-manager/trm/v2"
	mock "github.com/stretchr/testify/mock"
)

// MockManager is an autogenerated mock type for the Manager type
type MockManager struct {
	mock.Mock
}

// Do provides a mock function with given fields: _a0, _a1
func (_m *MockManager) Do(_a0 context.Context, _a1 func(context.Context) error) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Do")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DoWithSettings provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockManager) DoWithSettings(_a0 context.Context, _a1 trm.Settings, _a2 func(context.Context) error) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for DoWithSettings")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, trm.Settings, func(context.Context) error) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockManager creates a new instance of MockManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockManager {
	mock := &MockManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
