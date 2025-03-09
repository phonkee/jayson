// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// Extension is an autogenerated mock type for the Extension type
type Extension struct {
	mock.Mock
}

// ExtendResponseObject provides a mock function with given fields: _a0, _a1
func (_m *Extension) ExtendResponseObject(_a0 context.Context, _a1 map[string]interface{}) bool {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ExtendResponseObject")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ExtendResponseWriter provides a mock function with given fields: _a0, _a1
func (_m *Extension) ExtendResponseWriter(_a0 context.Context, _a1 http.ResponseWriter) bool {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ExtendResponseWriter")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, http.ResponseWriter) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewExtension creates a new instance of Extension. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExtension(t interface {
	mock.TestingT
	Cleanup(func())
}) *Extension {
	mock := &Extension{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
