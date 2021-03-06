// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"context"
	"io"

	"github.com/c4-project/c4t/internal/plan/stage"
)

// Runner is the interface of parts of the tester that transform plans.
type Runner interface {
	// Stage returns the stage that this runner implements.
	Stage() stage.Stage

	// Run runs this type's processing stage on the plan pointed to by p.
	// It also takes a context, which can be used to cancel the process.
	// It returns an updated plan (which may or may not be p edited in-place), or an error.
	Run(ctx context.Context, p *Plan) (*Plan, error)

	// Closer captures that some stages have resources that need to be freed after use.
	io.Closer

	// TODO(@MattWindsor91): implement Close() calls for parents of Runners other than the director.
}
