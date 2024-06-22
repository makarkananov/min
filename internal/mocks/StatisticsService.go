// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "min/internal/core/domain"

	mock "github.com/stretchr/testify/mock"
)

// StatisticsService is an autogenerated mock type for the StatisticsService type
type StatisticsService struct {
	mock.Mock
}

// AddEvent provides a mock function with given fields: ctx, event
func (_m *StatisticsService) AddEvent(ctx context.Context, event domain.Event) error {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for AddEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Event) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStatisticsService creates a new instance of StatisticsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStatisticsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *StatisticsService {
	mock := &StatisticsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
