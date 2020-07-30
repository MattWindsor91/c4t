// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package singleobs

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/observing"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
)

// Logger lifts a Logger to an observer.
type Logger log.Logger

// OnBuild logs build messages.
func (l *Logger) OnBuild(b builder.Message) {
	switch b.Kind {
	case observing.BatchStep:
		l.onBuildRequest(b.Request)
	}
}

// OnBuildRequest logs failed compile and run results.
func (l *Logger) onBuildRequest(r *builder.Request) {
	switch {
	case r.Compile != nil && r.Compile.Result.Status != status.Ok:
		(*log.Logger)(l).Printf("subject %q on compiler %q: %s", r.Name, r.Compile.CompilerID.String(), r.Compile.Result.Status)
	case r.Run != nil && r.Run.Result.Status != status.Ok:
		(*log.Logger)(l).Printf("subject %q on compiler %q: %s", r.Name, r.Run.CompilerID.String(), r.Run.Result.Status)
	}
}

// OnBuildFinish does nothing.
func (l *Logger) OnBuildFinish() {}

// OnCompilerConfig logs compiler config messages.
func (l *Logger) OnCompilerConfig(m compiler.Message) {
	switch m.Kind {
	case observing.BatchStart:
		l.onCompilerPlanStart(m.Num)
	case observing.BatchStep:
		l.onCompilerPlan(*m.Configuration)
	}
}

// onCompilerPlanStart briefly logs a compiler start.
func (l *Logger) onCompilerPlanStart(ncompilers int) {
	(*log.Logger)(l).Printf("planning %d compiler(s)...\n", ncompilers)
}

// onCompilerPlan logs a compiler plan.
func (l *Logger) onCompilerPlan(nc compiler.Named) {
	(*log.Logger)(l).Printf("compiler %s:\n", nc.ID)
	if nc.SelectedOpt != nil {
		(*log.Logger)(l).Printf(" - opt: %q:\n", nc.SelectedOpt.Name)
	}
	if !ystring.IsBlank(nc.SelectedMOpt) {
		(*log.Logger)(l).Printf(" - m/opt: %q:\n", nc.SelectedMOpt)
	}
}

// OnCopyStart briefly logs a file-copy start.
func (l *Logger) OnCopyStart(nfiles int) {
	(*log.Logger)(l).Printf("copying %d files...\n", nfiles)
}

// OnCopy does nothing.
func (l *Logger) OnCopy(_, _ string) {}

// OnCopyFinish does nothing.
func (l *Logger) OnCopyFinish() {}
