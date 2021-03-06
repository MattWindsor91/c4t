// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	compiler "github.com/c4-project/c4t/internal/model/service/compiler"
	mock "github.com/stretchr/testify/mock"

	optlevel "github.com/c4-project/c4t/internal/model/service/compiler/optlevel"

	stringhelp "github.com/c4-project/c4t/internal/helper/stringhelp"
)

// Inspector is an autogenerated mock type for the Inspector type
type Inspector struct {
	mock.Mock
}

// DefaultMOpts provides a mock function with given fields: c
func (_m *Inspector) DefaultMOpts(c *compiler.Compiler) (stringhelp.Set, error) {
	ret := _m.Called(c)

	var r0 stringhelp.Set
	if rf, ok := ret.Get(0).(func(*compiler.Compiler) stringhelp.Set); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(stringhelp.Set)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*compiler.Compiler) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DefaultOptLevels provides a mock function with given fields: c
func (_m *Inspector) DefaultOptLevels(c *compiler.Compiler) (stringhelp.Set, error) {
	ret := _m.Called(c)

	var r0 stringhelp.Set
	if rf, ok := ret.Get(0).(func(*compiler.Compiler) stringhelp.Set); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(stringhelp.Set)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*compiler.Compiler) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OptLevels provides a mock function with given fields: c
func (_m *Inspector) OptLevels(c *compiler.Compiler) (map[string]optlevel.Level, error) {
	ret := _m.Called(c)

	var r0 map[string]optlevel.Level
	if rf, ok := ret.Get(0).(func(*compiler.Compiler) map[string]optlevel.Level); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]optlevel.Level)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*compiler.Compiler) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
