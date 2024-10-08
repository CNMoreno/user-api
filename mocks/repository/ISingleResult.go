// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ISingleResult is an autogenerated mock type for the ISingleResult type
type ISingleResult struct {
	mock.Mock
}

// Decode provides a mock function with given fields: v
func (_m *ISingleResult) Decode(v interface{}) error {
	ret := _m.Called(v)

	if len(ret) == 0 {
		panic("no return value specified for Decode")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Err provides a mock function with given fields:
func (_m *ISingleResult) Err() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Err")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewISingleResult creates a new instance of ISingleResult. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewISingleResult(t interface {
	mock.TestingT
	Cleanup(func())
}) *ISingleResult {
	mock := &ISingleResult{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
