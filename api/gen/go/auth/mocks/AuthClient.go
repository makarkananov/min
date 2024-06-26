// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	authv1 "min/api/gen/go/auth"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// AuthClient is an autogenerated mock type for the AuthClient type
type AuthClient struct {
	mock.Mock
}

// ChangeLinksRemaining provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) ChangeLinksRemaining(ctx context.Context, in *authv1.ChangeLinksRemainingRequest, opts ...grpc.CallOption) (*authv1.ChangeLinksRemainingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ChangeLinksRemaining")
	}

	var r0 *authv1.ChangeLinksRemainingResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *authv1.ChangeLinksRemainingRequest, ...grpc.CallOption) (*authv1.ChangeLinksRemainingResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *authv1.ChangeLinksRemainingRequest, ...grpc.CallOption) *authv1.ChangeLinksRemainingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*authv1.ChangeLinksRemainingResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *authv1.ChangeLinksRemainingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) Login(ctx context.Context, in *authv1.LoginRequest, opts ...grpc.CallOption) (*authv1.LoginResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 *authv1.LoginResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *authv1.LoginRequest, ...grpc.CallOption) (*authv1.LoginResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *authv1.LoginRequest, ...grpc.CallOption) *authv1.LoginResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*authv1.LoginResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *authv1.LoginRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) Register(ctx context.Context, in *authv1.RegisterRequest, opts ...grpc.CallOption) (*authv1.RegisterResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 *authv1.RegisterResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *authv1.RegisterRequest, ...grpc.CallOption) (*authv1.RegisterResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *authv1.RegisterRequest, ...grpc.CallOption) *authv1.RegisterResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*authv1.RegisterResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *authv1.RegisterRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateToken provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) ValidateToken(ctx context.Context, in *authv1.ValidateTokenRequest, opts ...grpc.CallOption) (*authv1.ValidateTokenResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 *authv1.ValidateTokenResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *authv1.ValidateTokenRequest, ...grpc.CallOption) (*authv1.ValidateTokenResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *authv1.ValidateTokenRequest, ...grpc.CallOption) *authv1.ValidateTokenResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*authv1.ValidateTokenResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *authv1.ValidateTokenRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAuthClient creates a new instance of AuthClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthClient {
	mock := &AuthClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
