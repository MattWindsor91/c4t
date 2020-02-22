package compiler

import (
	"context"
	"io"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// SingleRunner is the interface of things that can run compilers.
type SingleRunner interface {
	// RunCompiler runs the compiler pointed to by c on the input files infiles.
	// On success, it outputs a binary to outfile.
	// If applicable, errw will be connected to the compiler's standard error.
	RunCompiler(ctx context.Context, c *model.NamedCompiler, infiles []string, outfile string, errw io.Writer) error
}

// SubjectPather is the interface of types that can produce path sets for compilations.
type SubjectPather interface {
	// Prepare sets up the directories ready to serve through SubjectPaths.
	// It takes the list of compiler IDs that are to be represented in the pathset.
	Prepare(compilers []model.ID) error

	// SubjectPaths gets the binary and log file paths for the subject/compiler pair sc.
	SubjectPaths(sc SubjectCompile) subject.CompileFileset
}

// Config represents the configuration that goes into a batch compiler run.
type Config struct {
	// Driver is what the compiler should use to run single compiler jobs.
	Driver SingleRunner

	// MachineID is the machine ID to use when loading from a multi-machine plan.
	// This can be empty if the plan only contains one machine.
	MachineID model.ID

	// Paths is the pathset for this compiler run.
	Paths SubjectPather
}

// Run runs a compiler configured by this config.
func (c *Config) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	cm, err := New(c, p)
	if err != nil {
		return nil, err
	}
	return cm.Run(ctx)
}
