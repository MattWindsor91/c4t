// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"context"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/plan/analysis"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/subject/status"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/MattWindsor91/act-tester/internal/subject/corpus"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

// TestAnalyse_errors tests various errors while analysing plans.
func TestAnalyse_errors(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		p   *plan.Plan
		ctx func() context.Context
		err error
	}{
		"no-plan": {
			p:   nil,
			ctx: context.Background,
			err: plan.ErrNil,
		},
		"no-corpus": {
			p:   &plan.Plan{Metadata: plan.Metadata{Version: plan.CurrentVer}},
			ctx: context.Background,
			err: corpus.ErrNone,
		},
		"done-context": {
			p: plan.Mock(),
			ctx: func() context.Context {
				wc, cf := context.WithCancel(context.Background())
				cf()
				return wc
			},
			err: context.Canceled,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := analysis.Analyse(c.ctx(), c.p, 10)
			testhelp.ExpectErrorIs(t, err, c.err, "analysing broken plan")
		})
	}
}

// TestAnalyse_mock tests that analysing an example plan with Analyse gives the expected collation.
func TestAnalyse_mock(t *testing.T) {
	t.Parallel()

	m := plan.Mock()
	crp, err := analysis.Analyse(context.Background(), m, 10)
	require.NoError(t, err, "unexpected error analysing")

	cases := map[string]struct {
		subc         status.Status
		wantSubjects []string
	}{
		"flagged":          {subc: status.Flagged, wantSubjects: []string{"baz"}},
		"run-failures":     {subc: status.RunFail, wantSubjects: []string{}},
		"run-timeouts":     {subc: status.RunTimeout, wantSubjects: []string{"barbaz"}},
		"compile-failures": {subc: status.CompileFail, wantSubjects: []string{"bar"}},
		"compile-timeouts": {subc: status.CompileTimeout, wantSubjects: []string{}},
		"successes":        {subc: status.Ok, wantSubjects: []string{"foo"}},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := crp.ByStatus[c.subc].Names()
			assert.Equal(t, c.wantSubjects, got, "wrong subjects")
		})
	}
}