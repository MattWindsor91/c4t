// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter_test

import (
	"testing"

	"github.com/MattWindsor91/act-tester/internal/controller/lifter/mocks"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/controller/lifter"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// TestNew_errors tests the error result of New in various situations.
func TestNew_errors(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		// cdelta modifies the configuration from a known-working value.
		cdelta func(*lifter.Config) *lifter.Config
		// pdelta modifies the plan from a known-working value.
		pdelta func(*plan.Plan) *plan.Plan
		// err is any error expected to occur on constructing with the modified plan and configuraiton.
		err error
	}{
		"ok": {
			err: nil,
		},
		"nil-config": {
			cdelta: func(c *lifter.Config) *lifter.Config {
				return nil
			},
			err: lifter.ErrConfigNil,
		},
		"nil-maker": {
			cdelta: func(c *lifter.Config) *lifter.Config {
				c.Driver = nil
				return c
			},
			err: lifter.ErrDriverNil,
		},
		"nil-paths": {
			cdelta: func(c *lifter.Config) *lifter.Config {
				c.Paths = nil
				return c
			},
			err: iohelp.ErrPathsetNil,
		},
		"nil-plan": {
			pdelta: func(p *plan.Plan) *plan.Plan {
				return nil
			},
			err: plan.ErrNil,
		},
		"nil-backend": {
			pdelta: func(p *plan.Plan) *plan.Plan {
				p.Backend = nil
				return p
			},
			err: lifter.ErrNoBackend,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				msl mocks.SingleLifter
				mpl mocks.Pather
			)
			cfg := &lifter.Config{
				Driver: &msl,
				Paths:  &mpl,
				Stderr: nil,
			}
			if f := c.cdelta; f != nil {
				cfg = f(cfg)
			}

			p := plan.Mock()
			if f := c.pdelta; f != nil {
				p = f(p)
			}

			_, err := lifter.New(cfg, p)
			testhelp.ExpectErrorIs(t, err, c.err, "in New()")
			msl.AssertExpectations(t)
		})
	}
}
