// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"errors"
	"log"

	"github.com/MattWindsor91/act-tester/internal/plan/analysis"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/stage/perturber"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/remote"

	"github.com/MattWindsor91/act-tester/internal/director/observer"

	"github.com/MattWindsor91/act-tester/internal/director/pathset"

	"github.com/MattWindsor91/act-tester/internal/stage/lifter"

	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer"

	"github.com/mitchellh/go-homedir"

	"github.com/MattWindsor91/act-tester/internal/stage/planner"

	"github.com/MattWindsor91/act-tester/internal/config"
)

var (
	// ErrObserverNil occurs when we try to pass a nil observer as an option.
	ErrObserverNil = errors.New("observer nil")

	// ErrNoMachines occurs when we try to build a director without defining any machines defined.
	ErrNoMachines = errors.New("no machines defined")

	// ErrNoOutDir occurs when we try to build a director with no output directory specified in the config.
	ErrNoOutDir = errors.New("no output directory specified in config")
)

// Option is the type of options for the director.
type Option func(*Director) error

// Options bundles the separate options ops into a single option.
func Options(ops ...Option) Option {
	return func(d *Director) error {
		for _, op := range ops {
			if err := op(d); err != nil {
				return err
			}
		}
		return nil
	}
}

// ObserveWith adds obs to the director's observer pool.
func ObserveWith(obs ...observer.Observer) Option {
	return func(d *Director) error {
		for _, o := range obs {
			if o == nil {
				return ErrObserverNil
			}
			d.observers = append(d.observers, o)
		}
		return nil
	}
}

// LogWith sets the director's logger to l.
func LogWith(l *log.Logger) Option {
	// TODO(@MattWindsor91): replace logger with observer
	return func(d *Director) error {
		d.l = l
		return nil
	}
}

// FilterMachines filters the director's machine set with glob.
func FilterMachines(glob id.ID) Option {
	return func(d *Director) error {
		if glob.IsEmpty() {
			return nil
		}
		var err error
		d.machines, err = d.machines.Filter(glob)
		return err
	}
}

// OverrideQuantities overrides the director's quantities with qs.
func OverrideQuantities(qs quantity.RootSet) Option {
	return func(d *Director) error {
		d.quantities.Override(qs)
		return nil
	}
}

// SSH sets the director's SSH config to s.
func SSH(s *remote.Config) Option {
	return func(d *Director) error {
		d.ssh = s
		return nil
	}
}

// OutDir sets the director's paths relative to dir.
// It performs home directory expansion in dir.
func OutDir(dir string) Option {
	return func(d *Director) error {
		if ystring.IsBlank(dir) {
			return ErrNoOutDir
		}
		edir, err := homedir.Expand(dir)
		if err != nil {
			return err
		}
		d.paths = pathset.New(edir)
		return nil
	}
}

// Filters adds fs to the set of filters to use for any analyses this director runs.
func Filters(fs analysis.FilterSet) Option {
	return func(d *Director) error {
		d.filters = append(d.filters, fs...)
		return nil
	}
}

// FiltersFromFile loads a filter set from path, if it is non-blank.
func FiltersFromFile(path string) Option {
	return func(d *Director) error {
		if ystring.IsBlank(path) {
			return nil
		}
		fs, err := analysis.LoadFilterSet(path)
		if err != nil {
			return err
		}
		return Filters(fs)(d)
	}
}

// Env groups together the bits of configuration that pertain to dealing with the environment.
type Env struct {
	// Fuzzer is a single-shot fuzzing driver.
	// TODO(@MattWindsor91): this overlaps nontrivially with Planner; both should use the same dumper!
	Fuzzer fuzzer.Driver

	// Lifter is a single-shot recipe lifter.
	Lifter lifter.SingleLifter

	// CInspector is the compiler inspector used for perturbing compiler optimisation levels.
	CInspector compiler.Inspector

	// Planner instructs any planners built for this director as to how to acquire information about compilers, etc.
	Planner planner.Source
}

// Check makes sure the environment is sensible.
func (e Env) Check() error {
	if e.Fuzzer == nil {
		return fuzzer.ErrDriverNil
	}
	if e.Lifter == nil {
		return lifter.ErrDriverNil
	}
	if e.CInspector == nil {
		return perturber.ErrCInspectorNil
	}
	return e.Planner.Check()
}

// ConfigFromGlobal extracts the parts of a global config file relevant to a director, and builds a config from them.
func ConfigFromGlobal(g *config.Config) Option {
	return Options(
		FiltersFromFile(g.Paths.FilterFile),
		OutDir(g.Paths.OutDir),
		OverrideQuantities(g.Quantities),
		SSH(g.SSH),
	)
}
