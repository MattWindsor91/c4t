// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package rmem implements rudimentary backend support for RMEM.
//
// Presently, rmem is implemented as a herdtools-style backend, despite not being a herdtools project.
// This will likely change later on.
package rmem

import (
	"context"
	"fmt"
	"io"

	"github.com/MattWindsor91/c4t/internal/model/id"

	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend"

	"github.com/MattWindsor91/c4t/internal/model/service"
	backend2 "github.com/MattWindsor91/c4t/internal/model/service/backend"
)

var armArgs = []string{
	"-model", "promising",
	"-model", "promise_first",
	"-model", "promising_parallel_thread_state_search",
	"-model", "promising_parallel_without_follow_trace",
	"-priority_reduction", "false",
	"-interactive", "false",
	"-hash_prune", "false",
	"-allow_partial", "true",
	"-loop_limit", "2",
}

// Rmem holds implementations of various backend responsiblities for Rmem.
type Rmem struct{}

func (Rmem) LiftStandalone(ctx context.Context, j backend2.LiftJob, r service.RunInfo, x service.Runner, w io.Writer) error {
	// TODO(@MattWindsor91): this is only useful when the source is an assembly Litmus test.
	// We should distinguish between the two.

	if j.Arch.IsEmpty() {
		j.Arch = id.ArchAArch64
	}
	// TODO(@MattWindsor91): eventually support other things here
	if !j.Arch.Equal(id.ArchAArch64) {
		return fmt.Errorf("%w: only AArch64 is supported for now", backend.ErrNotSupported)
	}

	// TODO(@MattWindsor91): sanitise here
	r.Override(service.RunInfo{Args: append(armArgs, j.In.Litmus.Path)})
	return x.WithStdout(w).Run(ctx, r)
}

// LiftExe doesn't work.
func (Rmem) LiftExe(context.Context, backend2.LiftJob, service.RunInfo, service.Runner) error {
	return fmt.Errorf("%w: harness making", backend.ErrNotSupported)
}