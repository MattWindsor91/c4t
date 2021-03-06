// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package gccnt contains the app definition for c4t-gccnt.
package gccnt

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/c4-project/c4t/internal/mutation"

	"github.com/c4-project/c4t/internal/ux/stdflag"

	"github.com/c4-project/c4t/internal/serviceimpl/compiler/gcc"

	"github.com/c4-project/c4t/internal/tool/gccnt"

	// This name is because every single time I try to use v2 named as 'cli', my IDE decides to replace it with v1.
	// Yes, I know, I shouldn't work around IDE issues by obfuscating my code, but I'm at my wit's end.
	c "github.com/urfave/cli/v2"
)

// App creates the c4t-gccnt app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:   "c4t-gccnt",
		Usage:  "wraps gcc with various optional failure modes",
		Flags:  flags(),
		Action: run,
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

const (
	flagOutput           = "o"
	flagBin              = "nt-bin"
	flagDryRun           = "nt-dryrun"
	flagDivergeOnOpt     = "nt-diverge-opt"
	flagErrorOnOpt       = "nt-error-opt"
	flagDivergeMutPeriod = "nt-diverge-mutant-period"
	flagErrorMutPeriod   = "nt-error-mutant-period"
	flagHitMutPeriod     = "nt-hit-mutant-period"
	flagMutant           = "nt-mutant"
	flagPthread          = "pthread"
	flagStd              = "std"
	flagMarch            = "march"
	flagMcpu             = "mcpu"
)

func flags() []c.Flag {
	fs := []c.Flag{
		&c.PathFlag{
			Name:      flagOutput,
			Usage:     "output file",
			Required:  true,
			TakesFile: true,
			Value:     "a.out",
		},
		&c.StringFlag{
			Name:  flagBin,
			Usage: "the 'real' compiler `command` to run",
			Value: "gcc",
		},
		&c.BoolFlag{
			Name:  flagDryRun,
			Usage: "print the outcome of running gccn't instead of doing it",
		},
		&c.StringSliceFlag{
			Name:  flagErrorOnOpt,
			Usage: "optimisation `levels` (minus the '-O') on which gccn't should exit with an error",
		},
		&c.StringSliceFlag{
			Name:  flagDivergeOnOpt,
			Usage: "optimisation `levels` (minus the '-O') on which gccn't should diverge",
		},
		&c.Uint64Flag{
			Name:        flagDivergeMutPeriod,
			Usage:       "diverge when the mutant number is a multiple of this `period`",
			DefaultText: "disabled",
		},
		&c.Uint64Flag{
			Name:        flagErrorMutPeriod,
			Usage:       "error when the mutant number is a multiple of this `period`",
			DefaultText: "disabled",
		},
		&c.Uint64Flag{
			Name:        flagHitMutPeriod,
			Usage:       "report a hit when the mutant number is a multiple of this `period`",
			DefaultText: "disabled",
		},
		&c.Uint64Flag{
			Name:        flagMutant,
			EnvVars:     []string{mutation.EnvVar},
			Usage:       "the mutant `number` to use if simulating mutation testing",
			DefaultText: "mutation disabled",
		},
		&c.StringFlag{
			Name:  flagStd,
			Usage: "`standard` to pass through to gcc",
		},
		&c.StringFlag{
			Name:  flagMarch,
			Usage: "architecture optimisation `spec` to pass through to gcc",
		},
		&c.StringFlag{
			Name:  flagMcpu,
			Usage: "cpu optimisation `spec` to pass through to gcc",
		},
		&c.BoolFlag{
			Name:  flagPthread,
			Usage: "passes through pthread to gcc",
		},
	}
	return append(fs, oflags()...)
}

func oflags() []c.Flag {
	flags := make([]c.Flag, len(gcc.OptLevelNames))
	for i, o := range gcc.OptLevelNames {
		flags[i] = &c.BoolFlag{
			Name:  "O" + o,
			Usage: fmt.Sprintf("optimisation level '%s'", o),
		}
	}
	return flags
}

func run(ctx *c.Context) error {
	olevel, err := geto(ctx)
	if err != nil {
		return err
	}

	g := gccnt.Gccnt{
		Mutant:   ctx.Uint64(flagMutant),
		Bin:      ctx.String(flagBin),
		In:       ctx.Args().Slice(),
		Out:      ctx.Path(flagOutput),
		OptLevel: olevel,
		Conds:    makeConditionSet(ctx),
		March:    ctx.String(flagMarch),
		Mcpu:     ctx.String(flagMcpu),
		Pthread:  ctx.Bool(flagPthread),
		Std:      ctx.String(flagStd),
	}

	if ctx.Bool(flagDryRun) {
		return g.DryRun(ctx.Context, os.Stderr)
	}
	return g.Run(ctx.Context, os.Stdout, os.Stderr)
}

func makeConditionSet(ctx *c.Context) gccnt.ConditionSet {
	return gccnt.ConditionSet{
		Diverge: gccnt.Condition{
			Opts:      ctx.StringSlice(flagDivergeOnOpt),
			MutPeriod: ctx.Uint64(flagDivergeMutPeriod),
		},
		Error: gccnt.Condition{
			Opts:      ctx.StringSlice(flagErrorOnOpt),
			MutPeriod: ctx.Uint64(flagErrorMutPeriod),
		},
		MutHitPeriod: ctx.Uint64(flagHitMutPeriod),
	}
}

func geto(ctx *c.Context) (string, error) {
	set := false
	o := "0"

	for _, possible := range gcc.OptLevelNames {
		if ctx.Bool("O" + possible) {
			o = possible
			if set {
				return "", errors.New("multiple optimisation levels defined")
			}
			set = true
		}
	}

	return o, nil
}

/*
if g.Mutant, err = getMutant(); err != nil || g.Mutant == 0 {
return err
}
func getMutant() (uint64, error) {
	return strconv.ParseUint(os.Getenv(mutation.EnvVar), 10, 64)
}
*/
