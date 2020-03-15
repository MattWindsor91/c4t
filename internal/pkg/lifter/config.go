// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"context"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
)

// Config contains configuration used to run a lifter for a particular machine, perhaps across multiple plans.
type Config struct {
	// Maker is a harness maker.
	Maker HarnessMaker

	// Logger is the logger to use for this lifter.
	// This may be nil, in which case the lifter will log silently.
	Logger *log.Logger

	// Observer tracks the lifter's progress across a corpus.
	Observer builder.Observer

	// Paths does path resolution and preparation for the incoming lifter.
	Paths Pather
}

// Run is shorthand for constructing a Lifter using c, then running it with p.
func (c *Config) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	l, err := New(c, p)
	if err != nil {
		return nil, err
	}
	return l.Run(ctx)
}
