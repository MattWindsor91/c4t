// Code generated by mockery v2.1.0. DO NOT EDIT.

package mocks

import (
	saver "github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"
	mock "github.com/stretchr/testify/mock"
)

// Observer is an autogenerated mock type for the Observer type
type Observer struct {
	mock.Mock
}

// OnArchive provides a mock function with given fields: s
func (_m *Observer) OnArchive(s saver.ArchiveMessage) {
	_m.Called(s)
}