// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	machine "github.com/c4-project/c4t/internal/machine"
	mock "github.com/stretchr/testify/mock"
)

// Observer is an autogenerated mock type for the Observer type
type Observer struct {
	mock.Mock
}

// OnMachines provides a mock function with given fields: m
func (_m *Observer) OnMachines(m machine.Message) {
	_m.Called(m)
}
