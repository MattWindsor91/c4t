// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"
	mocks2 "github.com/MattWindsor91/act-tester/internal/model/litmus/mocks"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/MattWindsor91/act-tester/internal/plan"
)

// makePlan makes a 'valid' plan.
func makePlan() *plan.Plan {
	return &plan.Plan{
		Metadata: plan.Metadata{Version: plan.CurrentVer},
		Corpus: map[string]subject.Subject{
			"foo": *subject.NewOrPanic(litmus.New("foo.litmus", litmus.WithThreads(1))),
			"bar": *subject.NewOrPanic(litmus.New("bar.litmus", litmus.WithThreads(2))),
			"baz": *subject.NewOrPanic(litmus.New("baz.litmus", litmus.WithThreads(3))),
		},
	}
}

// TestNew_error checks that New processes various error conditions correctly.
func TestNew_error(t *testing.T) {
	t.Parallel()

	md := new(mocks2.StatDumper)
	mp := new(mocks.SubjectPather)

	opterr := errors.New("sup")

	cases := map[string]struct {
		// driver sets the driver for the constructor call.
		driver fuzzer.Driver
		// paths sets the pathset for the constructor call.
		paths fuzzer.SubjectPather
		// opts sets the options for the constructor call.
		opts []fuzzer.Option
		// err is any error expected to occur on constructing with the modified plan and configuraiton.
		err error
	}{
		"ok": {
			driver: fuzzer.AggregateDriver{Single: fuzzer.NopFuzzer{}, Stat: md},
			paths:  mp,
			err:    nil,
		},
		"nil-driver": {
			driver: nil,
			paths:  mp,
			err:    fuzzer.ErrDriverNil,
		},
		"nil-paths": {
			driver: fuzzer.AggregateDriver{Single: fuzzer.NopFuzzer{}, Stat: md},
			paths:  nil,
			err:    iohelp.ErrPathsetNil,
		},
		"bad-cycles": {
			driver: fuzzer.AggregateDriver{Single: fuzzer.NopFuzzer{}, Stat: md},
			paths:  mp,
			opts: []fuzzer.Option{
				fuzzer.OverrideQuantities(fuzzer.QuantitySet{SubjectCycles: -1}),
			},
			err: corpus.ErrSmall,
		},
		"err-option": {
			driver: fuzzer.AggregateDriver{Single: fuzzer.NopFuzzer{}, Stat: md},
			paths:  mp,
			opts: []fuzzer.Option{
				func(*fuzzer.Fuzzer) error {
					return opterr
				},
			},
			err: opterr,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := fuzzer.New(c.driver, c.paths, c.opts...)
			testhelp.ExpectErrorIs(t, err, c.err, "unexpected error in New")
		})
	}
}

// TestFuzzer_Run_error tests various error cases on Run.
func TestFuzzer_Run_error(t *testing.T) {
	t.Parallel()

	md := new(mocks2.StatDumper)
	mp := new(mocks.SubjectPather)

	cases := map[string]struct {
		pdelta func(*plan.Plan) *plan.Plan
		opts   []fuzzer.Option
		err    error
	}{
		"nil-plan": {
			pdelta: func(p *plan.Plan) *plan.Plan {
				return nil
			},
			err: plan.ErrNil,
		},
		"no-corpus": {
			pdelta: func(p *plan.Plan) *plan.Plan {
				p.Corpus = nil
				return p
			},
			err: corpus.ErrNone,
		},
		"bad-version": {
			pdelta: func(p *plan.Plan) *plan.Plan {
				p.Metadata.Version = 0
				return p
			},
			err: plan.ErrVersionMismatch,
		},
		"no-stage": {
			pdelta: func(p *plan.Plan) *plan.Plan {
				p.Metadata.Stages = []stage.Record{}
				return p
			},
			err: plan.ErrMissingStage,
		},
		"small-corpus": {
			opts: []fuzzer.Option{
				fuzzer.OverrideQuantities(
					fuzzer.QuantitySet{CorpusSize: 255},
				),
			},
			err: corpus.ErrSmall,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f, err := fuzzer.New(fuzzer.AggregateDriver{Single: fuzzer.NopFuzzer{}, Stat: md}, mp, c.opts...)
			require.NoError(t, err, "there shouldn't be an error yet!")

			p := makePlan()
			p.Metadata.ConfirmStage(stage.Plan, time.Now(), time.Now())
			if f := c.pdelta; f != nil {
				p = f(p)
			}
			_, err = f.Run(context.Background(), p)
			testhelp.ExpectErrorIs(t, err, c.err, "running fuzzer")
		})
	}
}

// TestFuzzer_Run_nop tests the happy path of running the Fuzzer with a driver that doesn't do anything.
func TestFuzzer_Run_nop(t *testing.T) {
	t.Parallel()

	md := new(mocks2.StatDumper)
	mp := new(mocks.SubjectPather)
	f, err := fuzzer.New(
		fuzzer.AggregateDriver{Single: fuzzer.NopFuzzer{}, Stat: md},
		mp,
	)
	require.NoError(t, err, "unexpected error in New")

	mp.On("Prepare").Return(nil).Once()
	mp.On("SubjectLitmus", mock.Anything).Return("fuzz.litmus")
	mp.On("SubjectTrace", mock.Anything).Return("fuzz.trace.txt")
	md.On("DumpStats", mock.Anything, mock.Anything, "fuzz.litmus").Return(nil)

	p := makePlan()
	p.Metadata.ConfirmStage(stage.Plan, time.Now(), time.Now())

	p2, err := f.Run(context.Background(), p)
	require.NoError(t, err, "unexpected error in Run")

	for name, s := range p2.Corpus {
		sc, err := fuzzer.ParseSubjectCycle(name)
		require.NoError(t, err, "name of fuzzer output not a subject-cycle name:", name)

		sf, ok := p.Corpus[sc.Name]
		require.Truef(t, ok, "subject %s in fuzzer output has no corresponding input", name)

		// This isn't exhaustive, but should be enough to catch out issues.
		//assert.Equal(t, sf.Stats, s.Stats, "stats mismatch")
		assert.Equal(t, sf.Source, s.Source, "litmus file mismatch")
	}

	mp.AssertExpectations(t)
	md.AssertExpectations(t)
}
