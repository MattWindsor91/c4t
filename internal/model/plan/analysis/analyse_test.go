// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/plan"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"
)

// TestAnalyse_empty tests that analysing an empty corpus gives an error.
func TestAnalyse_empty(t *testing.T) {
	t.Parallel()

	_, err := analysis.Analyse(context.Background(), &plan.Plan{Header: plan.Header{Version: plan.CurrentVer}}, 10)
	testhelp.ExpectErrorIs(t, err, corpus.ErrNone, "analysing empty plan")
}

// TestAnalyse_mock tests that analysing an example corpus gives the expected collation.
func TestAnalyse_mock(t *testing.T) {
	t.Parallel()

	m := plan.Mock()
	crp, err := analysis.Analyse(context.Background(), m, 10)
	if err != nil {
		t.Fatal("unexpected error collating mock corpus:", err)
	}

	cases := map[string]struct {
		subc         subject.Status
		wantSubjects []string
	}{
		"flagged":          {subc: subject.StatusFlagged, wantSubjects: []string{"baz"}},
		"run-failures":     {subc: subject.StatusRunFail, wantSubjects: []string{}},
		"run-timeouts":     {subc: subject.StatusRunTimeout, wantSubjects: []string{"barbaz"}},
		"compile-failures": {subc: subject.StatusCompileFail, wantSubjects: []string{"bar"}},
		"compile-timeouts": {subc: subject.StatusCompileTimeout, wantSubjects: []string{}},
		"successes":        {subc: subject.StatusOk, wantSubjects: []string{"foo"}},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := crp.ByStatus[c.subc].Names()
			if !reflect.DeepEqual(got, c.wantSubjects) {
				t.Errorf("wrong subjects: got=%v; want=%v", got, c.wantSubjects)
			}
		})
	}
}