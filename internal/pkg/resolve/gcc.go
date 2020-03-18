// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package resolve

import (
	"context"
	"io"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// GCC represents GCC-style compilers such as GCC and Clang.
type GCC struct {
	// DefaultRun is the default run information for the particular compiler.
	DefaultRun model.CompilerRunInfo
}

// Compile compiles j according to run using a GCC-friendly invocation.
func (g GCC) Compile(ctx context.Context, _ id.ID, run *model.CompilerRunInfo, j model.CompileJob, errw io.Writer) error {
	orun := g.DefaultRun.Override(run)
	args := GCCArgs(orun, j)
	cmd := exec.CommandContext(ctx, orun.Cmd, args...)
	cmd.Stderr = errw
	return cmd.Run()
}

// GCCArgs computes the arguments to pass to GCC for running job j with run info run.
func GCCArgs(run model.CompilerRunInfo, j model.CompileJob) []string {
	args := run.Args
	args = append(args, "-o", j.Out)
	args = append(args, j.In...)
	return args
}
