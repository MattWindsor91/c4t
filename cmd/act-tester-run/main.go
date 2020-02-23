package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/runner"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const (
	defaultOutDir    = "run_results"
	usageFlagMachine = "specifies a machine in a multi-machine plan to run"
)

func main() {
	if err := run(os.Args, os.Stdout, os.Stderr); err != nil {
		ux.LogTopError(err)
	}
}

func run(args []string, outw, errw io.Writer) error {
	var (
		act   interop.ActRunner
		dir   string
		pfile string
		pmach string
	)

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	ux.ActRunnerFlags(fs, &act)
	ux.OutDirFlag(fs, &dir, defaultOutDir)
	ux.PlanFileFlag(fs, &pfile)
	fs.StringVar(&pmach, ux.FlagMachine, "", usageFlagMachine)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	cfg := runner.Config{
		Logger:    log.New(errw, "", 0),
		Parser:    &act,
		MachineID: model.IDFromString(pmach),
		Paths:     runner.NewPathset(dir),
	}
	return makeAndRunRunner(&cfg, pfile, outw)
}

func makeAndRunRunner(c *runner.Config, pfile string, outw io.Writer) error {
	p, perr := ux.LoadPlan(pfile)
	if perr != nil {
		return perr
	}
	run, rerr := runner.New(c, p)
	if rerr != nil {
		return rerr
	}
	out, oerr := run.Run(context.Background())
	if oerr != nil {
		return oerr
	}
	je := json.NewEncoder(outw)
	return je.Encode(out)
}
