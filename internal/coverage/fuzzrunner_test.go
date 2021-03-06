// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package coverage_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	mocks4 "github.com/c4-project/c4t/internal/model/service/backend/mocks"

	"github.com/c4-project/c4t/internal/helper/srvrun"

	backend2 "github.com/c4-project/c4t/internal/model/service/backend"

	mocks3 "github.com/c4-project/c4t/internal/model/litmus/mocks"

	"github.com/c4-project/c4t/internal/model/recipe"

	"github.com/c4-project/c4t/internal/id"
	"github.com/stretchr/testify/mock"

	"github.com/c4-project/c4t/internal/model/litmus"
	"github.com/c4-project/c4t/internal/subject"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/coverage"
	"github.com/c4-project/c4t/internal/model/service/fuzzer"

	"github.com/c4-project/c4t/internal/stage/fuzzer/mocks"
)

// TestFuzzRunner_Run tests FuzzRunner.Run's happy path.
func TestFuzzRunner_Run(t *testing.T) {
	td := t.TempDir()

	var (
		f  mocks.SingleFuzzer
		l  mocks4.SingleLifter
		s  mocks3.StatDumper
		dr srvrun.DryRunner
	)
	f.Test(t)
	l.Test(t)
	s.Test(t)

	conf := fuzzer.Config{Params: map[string]string{"fus": "ro dah"}}
	fr := coverage.FuzzRunner{
		Fuzzer:     &f,
		Lifter:     &l,
		StatDumper: &s,
		Config:     &conf,
		Arch:       id.ArchX86,
		Runner:     dr,
	}
	sub := subject.NewOrPanic(litmus.NewOrPanic("foo.litmus", litmus.WithArch(id.ArchC)))
	rc := coverage.RunContext{
		Seed:        4321,
		BucketDir:   td,
		NumInBucket: 1,
		Input:       sub,
	}

	// TODO(@MattWindsor91): it'd be good if we didn't have to do this, but it's needed for the litmus arch scraper.
	err := os.WriteFile(rc.OutLitmus(), []byte("C foo1\n"), 0644)
	require.NoError(t, err, "couldn't write stub litmus output")

	f.On("Fuzz", mock.Anything, mock.MatchedBy(func(f fuzzer.Job) bool {
		return f.Seed == rc.Seed &&
			f.OutLitmus == rc.OutLitmus() &&
			f.Config != nil &&
			reflect.DeepEqual(conf, *(f.Config))
	})).Return(nil).Once()
	l.On("Lift", mock.Anything, mock.MatchedBy(func(l backend2.LiftJob) bool {
		return l.Arch.Equal(fr.Arch) &&
			l.In.Source == backend2.LiftLitmus &&
			l.In.Litmus.Filepath() == rc.OutLitmus() &&
			l.Out.Target == backend2.ToDefault &&
			l.Out.Dir == rc.LiftOutDir()
	}), dr).Return(recipe.Recipe{}, nil).Once()
	s.On("DumpStats", mock.Anything, mock.AnythingOfType("*litmus.Statset"), rc.OutLitmus()).Return(nil).Once()

	require.NoError(t, fr.Run(context.Background(), rc), "mock fuzz run shouldn't error")

	f.AssertExpectations(t)
	l.AssertExpectations(t)
	s.AssertExpectations(t)
}
