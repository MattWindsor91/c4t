// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"io"

	"github.com/c4-project/c4t/internal/stage/planner"

	"github.com/c4-project/c4t/internal/helper/iohelp"
	"github.com/c4-project/c4t/internal/observing"
)

// Option is the type of options to supply to the coverage testbed maker's constructor.
type Option func(*Maker) error

// Options applies each option in opts successively.
func Options(opts ...Option) Option {
	return func(maker *Maker) error {
		for _, o := range opts {
			if err := o(maker); err != nil {
				return err
			}
		}
		return nil
	}
}

// ObserveWith adds each observer o into the observer list.
func ObserveWith(o ...Observer) Option {
	return func(maker *Maker) error {
		if err := observing.CheckObservers(o); err != nil {
			return err
		}
		maker.observers = append(maker.observers, o...)
		return nil
	}
}

// OverrideQuantities overrides the maker's quantity set with qs.
func OverrideQuantities(qs QuantitySet) Option {
	return func(maker *Maker) error {
		maker.qs.Override(qs)
		return nil
	}
}

// AddInputs adds paths to the input list.
func AddInputs(paths ...string) Option {
	return func(maker *Maker) error {
		ps, err := iohelp.ExpandMany(paths)
		if err != nil {
			return err
		}
		// TODO(@MattWindsor91): handle the other uses of this expansion at the option level?
		fs, err := planner.ExpandLitmusInputs(ps)
		if err != nil {
			return err
		}
		maker.inputs = append(maker.inputs, fs...)

		return nil
	}
}

// SendStderrTo redirects stderr from commands to w.
func SendStderrTo(w io.Writer) Option {
	return func(maker *Maker) error {
		maker.errw = w
		return nil
	}
}
