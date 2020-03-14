// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package lifter contains the part of the tester framework that lifts litmus tests to compilable C.
// It does so by means of a backend HarnessMaker.
package lifter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

var (
	// ErrMakerNil occurs when a lifter runs without a HarnessMaker set.
	ErrMakerNil = errors.New("harness maker nil")

	// ErrNoBackend occurs when backend information is missing.
	ErrNoBackend = errors.New("no backend provided")
)

// HarnessMaker is an interface capturing the ability to make test harnesses.
type HarnessMaker interface {
	// MakeHarness asks the harness maker to make the test harness described by spec.
	// It returns a list outfiles of files created (C files, header files, etc.), and/or an error err.
	MakeHarness(ctx context.Context, spec model.HarnessSpec) (outFiles []string, err error)
}

// Lifter holds the main configuration for the lifter part of the tester framework.
type Lifter struct {
	// conf is the configuration used for this lifter.
	conf Config

	// plan is the plan on which this lifter is operating.
	plan plan.Plan

	// l is the logger to use for this lifter.
	l *log.Logger
}

// New constructs a new Lifter given config c and plan p.
func New(c *Config, p *plan.Plan) (*Lifter, error) {
	if c.Maker == nil {
		return nil, ErrMakerNil
	}
	if p == nil {
		return nil, plan.ErrNil
	}
	if p.Backend == nil || p.Backend.ID.IsEmpty() {
		return nil, ErrNoBackend
	}

	l := Lifter{
		conf: *c,
		plan: *p,
		l:    iohelp.EnsureLog(c.Logger),
	}
	return &l, nil
}

// Run runs a lifting job: taking every test subject in a plan and using a backend to lift each into a test harness.
func (l *Lifter) Run(ctx context.Context) (*plan.Plan, error) {
	l.l.Println("making output directory", l.conf.OutDir)
	if err := os.Mkdir(l.conf.OutDir, 0744); err != nil {
		return nil, err
	}

	err := l.lift(ctx)
	return &l.plan, err
}

func (l *Lifter) lift(ctx context.Context) error {
	l.l.Println("now lifting")

	b, err := corpus.NewBuilder(corpus.BuilderConfig{
		Init:  l.plan.Corpus,
		NReqs: l.count(),
		Obs:   l.conf.Observer,
	})
	if err != nil {
		return fmt.Errorf("when making builder: %w", err)
	}
	mrng := l.plan.Header.Rand()

	var lerr error
	l.plan.Corpus, lerr = l.liftInner(ctx, mrng, b)
	return lerr
}

func (l *Lifter) liftInner(ctx context.Context, mrng *rand.Rand, b *corpus.Builder) (corpus.Corpus, error) {
	eg, ectx := errgroup.WithContext(ctx)
	var lc corpus.Corpus
	// It's very likely this will be a single element array.
	for _, a := range l.plan.Arches() {
		dir, derr := buildAndMkDir(l.conf.OutDir, a.Tags()...)
		if derr != nil {
			return nil, derr
		}
		j := l.makeJob(a, dir, mrng, b.SendCh)
		eg.Go(func() error {
			return j.Lift(ectx)
		})
	}
	eg.Go(func() error {
		var err error
		lc, err = b.Run(ectx)
		return err
	})
	err := eg.Wait()
	return lc, err
}

func (l *Lifter) makeJob(a model.ID, dir string, mrng *rand.Rand, resCh chan<- corpus.BuilderReq) Job {
	return Job{
		Arch:    a,
		Backend: l.plan.Backend.FQID(),
		OutDir:  dir,
		Maker:   l.conf.Maker,
		Corpus:  l.plan.Corpus,
		Rng:     rand.New(rand.NewSource(mrng.Int63())),
		ResCh:   resCh,
	}
}

// count counts the number of liftings that need doing.
func (l *Lifter) count() int {
	return len(l.plan.Arches()) * len(l.plan.Corpus)
}
