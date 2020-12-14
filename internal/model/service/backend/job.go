// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/MattWindsor91/c4t/internal/model/recipe"

	"github.com/MattWindsor91/c4t/internal/subject/compilation"
	"github.com/MattWindsor91/c4t/internal/subject/obs"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/c4t/internal/model/litmus"

	"github.com/MattWindsor91/c4t/internal/model/id"
)

var (
	// ErrBadSource occurs when the input of a LiftJob has a source set to an unknown value.
	ErrBadSource = errors.New("bad input source")
	// ErrBadTarget occurs when the output of a LiftJob has a target set to an unknown value.
	ErrBadTarget = errors.New("bad output target")
	// ErrInLitmusBlank occurs when the input file of a lifter job is checked and found to be blank.
	ErrInLitmusBlank = errors.New("input litmus file path blank")
	// ErrOutDirBlank occurs when the output directory of a lifter job is checked and found to be blank.
	ErrOutDirBlank = errors.New("output directory path blank")
)

// LiftJob is a specification of how to lift a test into a compilable recipe.
type LiftJob struct {
	// Backend specifies the backend to use to perform the lifting.
	Backend *Spec

	// Arch is the ID of the architecture for which a recipe should be prepared, if the recipe is architecture-specific.
	Arch id.ID

	// In is the input specification for this job.
	In LiftInput

	// Out is the output specification for this job.
	Out LiftOutput
}

// Check performs several in-flight checks on a lifter job.
func (l LiftJob) Check() error {
	if err := l.In.Check(); err != nil {
		return err
	}
	return l.Out.Check()
}

// LiftInput is a specification of the input of a lifting operation.
type LiftInput struct {
	// Source specifies the kind of thing that the lifter should consume.
	Source Source

	// Litmus gives information about an input Litmus test, if any.
	Litmus litmus.Litmus
}

// LiftLitmusInput is shorthand for creating a LiftLitmus-style LiftInput over l.
func LiftLitmusInput(l litmus.Litmus) LiftInput {
	return LiftInput{Source: LiftLitmus, Litmus: l}
}

// Check makes sure that this lift input has a valid source and the data required for it.
func (l LiftInput) Check() error {
	switch l.Source {
	case LiftLitmus:
		return l.checkLitmus()
	default:
		return fmt.Errorf("%w: %s", ErrBadSource, l.Source)
	}
}

func (l LiftInput) checkLitmus() error {
	if !l.Litmus.HasPath() {
		return ErrInLitmusBlank
	}
	return nil
}

// LiftOutput is a specification of the output of a lifting operation.
type LiftOutput struct {
	// Dir specifies the output directory into which the lifter should put outputs.
	Dir string

	// Target specifies the kind of thing that the lifter should create.
	Target Target
}

// Check makes sure that this lift output has a valid target and the data required for it.
func (l LiftOutput) Check() error {
	// TODO(@MattWindsor91): ToStandalone shouldn't need a directory
	if ystring.IsBlank(l.Dir) {
		return ErrOutDirBlank
	}
	switch l.Target {
	case ToDefault:
		return nil
	case ToExeRecipe:
		return nil
	case ToObjRecipe:
		return nil
	case ToStandalone:
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrBadTarget, l.Target)
	}
}

// Files reads s.OutDir as a directory and returns its contents as qualified paths.
// This is useful for using a recipe job to feed a compiler job.
func (l LiftOutput) Files() ([]string, error) {
	fs, err := ioutil.ReadDir(l.Dir)
	if err != nil {
		return nil, err
	}

	ps := make([]string, len(fs))
	i := 0
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		ps[i] = f.Name()
		i++
	}
	return ps[:i], nil
}

// RunJob is the type of jobs being sent to a backend for running.
type RunJob struct {
	// Backend specifies the backend to use to perform the running.
	Backend *Spec

	// Recipe is a pointer to the recipe that was fed into the compile stage for this compilation; this is useful for
	// backends that don't compile, and instead peruse files from the compiler recipe.
	Recipe *recipe.Recipe

	// CompileResult is a pointer to the result of any compilation that was done for the running.
	// It may be nil if there was no compilation.
	CompileResult *compilation.CompileResult

	// Obs points to the observation record that should be filled out by the runner.
	Obs *obs.Obs
}