// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analyser

import (
	"errors"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/observer"
)

// ErrObserverNil occurs if we pass a nil Observer to ObserveWith.
var ErrObserverNil = errors.New("observer nil")

// Option is the type of options to the analyser stage constructor.
type Option func(*Analyser) error

// Options applies each option in opts in turn.
func Options(opts ...Option) Option {
	return func(a *Analyser) error {
		for _, o := range opts {
			if err := o(a); err != nil {
				return err
			}
		}
		return nil
	}
}

// ParWorkers sets the number of parallel analyser workers to n.
func ParWorkers(n int) Option {
	return func(a *Analyser) error {
		a.nworkers = n
		return nil
	}
}

// ObserveWith adds each observer in obs to the observer set.
func ObserveWith(obs ...observer.Observer) Option {
	return func(a *Analyser) error {
		for _, o := range obs {
			if o == nil {
				return ErrObserverNil
			}
		}
		a.observers = append(a.observers, obs...)
		return nil
	}
}

// SaveToPathset makes this analyser stage save to the given pathset.
// This can be nil, in which case saving is disabled.
func SaveToPathset(ps *saver.Pathset) Option {
	return func(a *Analyser) error {
		// ps can be nil
		a.savePaths = ps
		return nil
	}
}