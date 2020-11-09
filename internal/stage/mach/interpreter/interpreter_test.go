// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package interpreter_test

import (
	"context"
	"errors"
	"io/ioutil"
	"path"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/interpreter"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/stretchr/testify/mock"

	"github.com/MattWindsor91/act-tester/internal/model/job/compile"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	mdl "github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler/mocks"
)

// TestInterpreter_Interpret tests Interpret on an example recipe.
func TestInterpreter_Interpret(t *testing.T) {
	t.Parallel()

	var mc mocks.Compiler

	r := recipe.New(
		"in",
		recipe.OutExe,
		recipe.AddFiles("body.c", "harness.c", "body.h"),
		recipe.CompileFileToObj(path.Join("in", "body.c")),
		recipe.CompileAllCToExe(),
	)
	c := mdl.Configuration{}
	it, err := interpreter.NewInterpreter(&mc, &c, "a.out", r)
	require.NoError(t, err, "error while making interpreter")

	mc.On("RunCompiler",
		mock.Anything,
		*compile.New(compile.Obj, &c, path.Join("in", "obj_0.o"), path.Join("in", "body.c")),
		ioutil.Discard,
	).Return(nil).Once().On("RunCompiler",
		mock.Anything,
		*compile.New(compile.Exe, &c, "a.out", path.Join("in", "obj_0.o"), path.Join("in", "harness.c")),
		ioutil.Discard,
	).Return(nil).Once()

	err = it.Interpret(context.Background())
	require.NoError(t, err, "error while running interpreter")

	mc.AssertExpectations(t)
}

// TestInterpreter_Interpret_compileError tests Interpret's response to a compiler error.
func TestInterpreter_Interpret_compileError(t *testing.T) {
	t.Parallel()

	var mc mocks.Compiler

	werr := errors.New("no me gusta")

	r := recipe.New(
		"in",
		recipe.OutExe,
		recipe.AddFiles("body.c", "harness.c", "body.h"),
		recipe.AddInstructions(recipe.Instruction{Op: recipe.Nop}),
		recipe.CompileFileToObj(path.Join("in", "body.c")),
		recipe.CompileAllCToExe(),
	)
	c := mdl.Configuration{}
	it, err := interpreter.NewInterpreter(&mc, &c, "a.out", r)
	require.NoError(t, err, "error while making interpreter")

	mc.On("RunCompiler",
		mock.Anything,
		*compile.New(compile.Obj, &c, path.Join("in", "obj_0.o"), path.Join("in", "body.c")),
		ioutil.Discard,
	).Return(werr).Once()
	// The second compile job should not be run.

	err = it.Interpret(context.Background())
	testhelp.ExpectErrorIs(t, err, werr, "wrong error while running interpreter")

	mc.AssertExpectations(t)
}

// TestInterpreter_Interpret_badInstruction tests whether a bad interpreter instruction is caught correctly.
func TestInterpreter_Interpret_badInstruction(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  []recipe.Instruction
		err error
	}{
		"bad-op":   {in: []recipe.Instruction{{Op: 42}}, err: interpreter.ErrBadOp},
		"bad-file": {in: []recipe.Instruction{recipe.PushInputInst("nonsuch.c")}, err: interpreter.ErrFileUnavailable},
		"reused-file": {in: []recipe.Instruction{
			recipe.PushInputInst("body.c"),
			recipe.PushInputInst("body.c"),
		}, err: interpreter.ErrFileUnavailable,
		},
		"reused-file-inputs": {in: []recipe.Instruction{
			recipe.PushInputsInst(filekind.CSrc),
			recipe.PushInputInst("body.c"),
		}, err: interpreter.ErrFileUnavailable,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var mc mocks.Compiler
			mc.Test(t)

			r := recipe.New(
				"in",
				recipe.OutExe,
				recipe.AddFiles("body.c", "harness.c", "body.h"),
				recipe.AddInstructions(c.in...),
			)
			cmp := mdl.Configuration{}
			it, err := interpreter.NewInterpreter(&mc, &cmp, "a.out", r)
			require.NoError(t, err, "error while making interpreter")

			err = it.Interpret(context.Background())
			testhelp.ExpectErrorIs(t, err, c.err, "running interpreter on bad instruction")

			mc.AssertExpectations(t)
		})
	}
}

// TestInterpreter_Interpret_tooManyObjs tests the interpreter's object overflow by setting its cap to a comically low
// amount, then overflowing it.
func TestInterpreter_Interpret_tooManyObjs(t *testing.T) {
	t.Parallel()

	var mc mocks.Compiler
	mc.Test(t)

	r := recipe.New(
		"in",
		recipe.OutExe,
		recipe.AddFiles("body.c", "harness.c", "body.h"),
		recipe.CompileFileToObj(path.Join("in", "body.c")),
		recipe.CompileFileToObj(path.Join("in", "harness.c")),
	)
	c := mdl.Configuration{}
	mc.On("RunCompiler",
		mock.Anything,
		*compile.New(compile.Obj, &c, path.Join("in", "obj_0.o"), path.Join("in", "body.c")),
		ioutil.Discard).Return(nil).Once()

	it, err := interpreter.NewInterpreter(&mc, &c, "a.out", r, interpreter.SetMaxObjs(1))
	require.NoError(t, err, "error while making interpreter")

	err = it.Interpret(context.Background())
	testhelp.ExpectErrorIs(t, err, interpreter.ErrObjOverflow, "running interpreter with overflowing objs")

	mc.AssertExpectations(t)
}

// TestNewInterpreter_errors tests NewInterpreter in various error conditions.
func TestNewInterpreter_errors(t *testing.T) {
	t.Parallel()

	r := recipe.New(
		"in",
		recipe.OutExe,
		recipe.AddFiles("body.c", "harness.c", "body.h"),
		recipe.CompileFileToObj(path.Join("in", "body.c")),
		recipe.CompileFileToObj(path.Join("in", "harness.c")),
	)
	cmp := mdl.Configuration{}

	cases := map[string]struct {
		useDriver bool
		compiler  *compiler.Configuration
		err       error
	}{
		"no-driver": {
			useDriver: false,
			compiler:  &cmp,
			err:       interpreter.ErrDriverNil,
		},
		"no-compiler-cfg": {
			useDriver: true,
			compiler:  nil,
			err:       interpreter.ErrCompilerConfigNil,
		},
		"ok": {
			useDriver: true,
			compiler:  &cmp,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var mc mocks.Compiler
			mc.Test(t)
			var d interpreter.Driver
			if c.useDriver {
				d = &mc
			}

			_, err := interpreter.NewInterpreter(d, c.compiler, "a.out", r)
			testhelp.ExpectErrorIs(t, err, c.err, "constructing new interpreter")
		})
	}
}
