// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter_test

import (
	"errors"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/stage/lifter/mocks"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/stage/lifter"
)

// TestNew_errors tests the error result of New in various situations.
func TestNew_errors(t *testing.T) {
	t.Parallel()

	opterr := errors.New("oopsie")

	cases := map[string]struct {
		// ddelta modifies the driver from a known-working value.
		ddelta func(*mocks.SingleLifter) lifter.SingleLifter
		// padelta modifies the pather from a known-working value.
		pdelta func(*mocks.Pather) lifter.Pather
		// os adds options to the constructor.
		os []lifter.Option
		// err is any error expected to occur on constructing with the modified plan and configuraiton.
		err error
	}{
		"ok": {
			err: nil,
		},
		"nil-driver": {
			ddelta: func(l *mocks.SingleLifter) lifter.SingleLifter {
				return nil
			},
			err: lifter.ErrDriverNil,
		},
		"nil-paths": {
			pdelta: func(p *mocks.Pather) lifter.Pather {
				return nil
			},
			err: iohelp.ErrPathsetNil,
		},
		"opt-err": {
			os: []lifter.Option{func(*lifter.Lifter) error {
				return opterr
			}},
			err: opterr,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				msl mocks.SingleLifter
				sl  lifter.SingleLifter
				mpl mocks.Pather
				pl  lifter.Pather
			)

			if f := c.ddelta; f != nil {
				sl = f(&msl)
			} else {
				sl = &msl
			}
			if f := c.pdelta; f != nil {
				pl = f(&mpl)
			} else {
				pl = &mpl
			}

			_, err := lifter.New(sl, pl, c.os...)
			testhelp.ExpectErrorIs(t, err, c.err, "in New()")
			msl.AssertExpectations(t)
		})
	}
}
