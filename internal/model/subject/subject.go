// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package subject contains types and functions for dealing with test subject records.
// Such subjects generally live in a plan; the separate package exists to accommodate the large amount of subject
// specific types and functions in relation to the other parts of a test plan.

package subject

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/model"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Normalise represents a single test subject in a corpus.
type Subject struct {
	// Stats is the statistics set for this subject.
	Stats model.Statset `toml:"stats,omitempty" json:"stats,omitempty"`

	// Fuzz is the fuzzing pathset for this subject, if it has been fuzzed.
	Fuzz *Fuzz `toml:"fuzz,omitempty" json:"fuzz,omitempty"`

	// Litmus is the (slashed) path to this subject's original Litmus file.
	OrigLitmus string `toml:"orig_litmus,omitempty" json:"orig_litmus,omitempty"`

	// Compiles contains information about this subject's compilation attempts.
	// It maps from the string form of each compiler's ID.
	// If nil, this subject hasn't had any compilations.
	Compiles map[string]CompileResult `toml:"compiles,omitempty" json:"compiles,omitempty"`

	// Recipes contains information about this subject's lifted test recipes.
	// It maps the string form of each recipe's target architecture's ID.
	// If nil, this subject hasn't had a recipe generated.
	Recipes map[string]recipe.Recipe `toml:"recipes,omitempty" json:"recipes,omitempty"`

	// Runs contains information about this subject's runs so far.
	// It maps from the string form of each compiler's ID.
	// If nil, this subject hasn't had any runs.
	Runs map[string]RunResult `toml:"runs,omitempty" json:"runs,omitempty"`
}

// New is a convenience constructor for subjects.
func New(origLitmus string, opt ...Option) (*Subject, error) {
	s := Subject{OrigLitmus: origLitmus}
	for _, c := range opt {
		if err := c(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// NewOrPanic is like New, but panics if there is an error.
// Use in tests only.
func NewOrPanic(origLitmus string, opt ...Option) *Subject {
	n, err := New(origLitmus, opt...)
	if err != nil {
		panic(err)
	}
	return n
}

// Option is the type of (functional) options to the New constructor.
type Option func(*Subject) error

// WithThreads is an option that sets the subject's threads to threads.
func WithThreads(threads int) Option {
	return func(s *Subject) error {
		s.Stats.Threads = threads
		return nil
	}
}

// BestLitmus tries to get the 'best' litmus test path for further development.
//
// When there is a fuzzing record for this subject, the fuzz output is the best path.
// Otherwise, if there is a non-empty Litmus file for this subject, that file is the best path.
// Else, BestLitmus returns an error.
func (s *Subject) BestLitmus() (string, error) {
	switch {
	case s.HasFuzzFile():
		return s.Fuzz.Files.Litmus, nil
	case s.OrigLitmus != "":
		return s.OrigLitmus, nil
	default:
		return "", ErrNoBestLitmus
	}
}

// HasFuzzFile gets whether this subject has a fuzzed testcase file.
func (s *Subject) HasFuzzFile() bool {
	return s.Fuzz != nil && ystring.IsNotBlank(s.Fuzz.Files.Litmus)
}

// Note that all of these maps work in basically the same way; their being separate and duplicated is just a
// consequence of Go not (yet) having generics.

// CompileResult gets the compilation result for the compiler ID cid.
func (s *Subject) CompileResult(cid id.ID) (CompileResult, error) {
	key := cid.String()
	c, ok := s.Compiles[key]
	if !ok {
		return CompileResult{}, fmt.Errorf("%w: compiler=%q", ErrMissingCompile, key)
	}
	return c, nil
}

// AddCompileResult sets the compilation information for compiler ID cid to c in this subject.
// It fails if there already _is_ a compilation.
func (s *Subject) AddCompileResult(cid id.ID, c CompileResult) error {
	s.ensureCompileMap()
	key := cid.String()
	if _, ok := s.Compiles[key]; ok {
		return fmt.Errorf("%w: compiler=%q", ErrDuplicateCompile, key)
	}
	s.Compiles[key] = c
	return nil
}

// ensureCompileMap makes sure this subject has a compile result map.
func (s *Subject) ensureCompileMap() {
	if s.Compiles == nil {
		s.Compiles = make(map[string]CompileResult)
	}
}

// Recipe gets the harness for the architecture with id arch.
func (s *Subject) Recipe(arch id.ID) (recipe.Recipe, error) {
	key := arch.String()
	h, ok := s.Recipes[key]
	if !ok {
		return recipe.Recipe{}, fmt.Errorf("%w: arch=%q", ErrMissingHarness, key)
	}
	return h, nil
}

// AddRecipe sets the harness information for arch to h in this subject.
// It fails if there already _is_ a harness for arch.
func (s *Subject) AddRecipe(arch id.ID, h recipe.Recipe) error {
	s.ensureHarnessMap()
	key := arch.String()
	if _, ok := s.Recipes[key]; ok {
		return fmt.Errorf("%w: arch=%q", ErrDuplicateHarness, key)
	}
	s.Recipes[key] = h
	return nil
}

// ensureHarnessMap makes sure this subject has a harness map.
func (s *Subject) ensureHarnessMap() {
	if s.Recipes == nil {
		s.Recipes = make(map[string]recipe.Recipe)
	}
}

// RunOf gets the run for the compiler with id cid.
func (s *Subject) RunOf(cid id.ID) (RunResult, error) {
	key := cid.String()
	h, ok := s.Runs[key]
	if !ok {
		return RunResult{}, fmt.Errorf("%w: compiler=%q", ErrMissingRun, key)
	}
	return h, nil
}

// AddRun sets the run information for cid to r in this subject.
// It fails if there already _is_ a run for cid.
func (s *Subject) AddRun(cid id.ID, r RunResult) error {
	s.ensureRunMap()
	key := cid.String()
	if _, ok := s.Runs[key]; ok {
		return fmt.Errorf("%w: compiler=%q", ErrDuplicateRun, key)
	}
	s.Runs[key] = r
	return nil
}

// ensureHarnessMap makes sure this subject has a harness map.
func (s *Subject) ensureRunMap() {
	if s.Runs == nil {
		s.Runs = make(map[string]RunResult)
	}
}

// AddName copies this subject into a new Named with the given name.
func (s *Subject) AddName(name string) *Named {
	return &Named{Name: name, Subject: *s}
}
