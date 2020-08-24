// Code generated by mockery v2.2.0. DO NOT EDIT.

package mocks

import (
	builder "github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	mock "github.com/stretchr/testify/mock"

	observer "github.com/MattWindsor91/act-tester/internal/stage/mach/observer"
)

// Observer is an autogenerated mock type for the Observer type
type Observer struct {
	mock.Mock
}

// OnBuild provides a mock function with given fields: _a0
func (_m *Observer) OnBuild(_a0 builder.Message) {
	_m.Called(_a0)
}

// OnMachineNodeAction provides a mock function with given fields: _a0
func (_m *Observer) OnMachineNodeAction(_a0 observer.Message) {
	_m.Called(_a0)
}
