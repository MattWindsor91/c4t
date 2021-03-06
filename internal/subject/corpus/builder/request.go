// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package builder

import (
	"context"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/model/recipe"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/subject"
)

// Request is the type of requests to a Builder.
type Request struct {
	// Name is the name of the subject to add or modify
	Name string `json:"name"`

	// Add is populated if this request is an Add.
	Add *Add `json:"add,omitempty"`

	// Compile is populated if this request is a Compile.
	Compile *Compile `json:"compile,omitempty"`

	// Recipe is populated if this request is a Recipe.
	Recipe *Recipe `json:"recipe,omitempty"`

	// Run is populated if this request is a Run.
	Run *Run `json:"run,omitempty"`
}

// SendTo tries to send this request down ch while checking to see if ctx has been cancelled.
func (b Request) SendTo(ctx context.Context, ch chan<- Request) error {
	select {
	case ch <- b:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Add is a request to add the given subject to the corpus.
type Add subject.Subject

// AddRequest constructs an add-subject request for subject s.
func AddRequest(s *subject.Named) Request {
	a := Add(s.Subject)
	return Request{Name: s.Name, Add: &a}
}

// Compile is a request to add the given compiler result to the named subject.
type Compile struct {
	// CompilerID is the ID of the compiler that produced this result.
	CompilerID id.ID `json:"compiler_id,omitempty"`

	// Result is the compile result.
	Result compilation.CompileResult `json:"result,omitempty"`
}

// CompileRequest constructs an add-compile request for the compilation with name name and result r.
func CompileRequest(name compilation.Name, r compilation.CompileResult) Request {
	return Request{Name: name.SubjectName, Compile: &Compile{CompilerID: name.CompilerID, Result: r}}
}

// Recipe is a request to add the given recipe to the named subject, under the named architecture.
type Recipe struct {
	// Arch is the ID of the architecture for which this lifting is occurring.
	Arch id.ID `json:"arch,omitempty"`

	// Recipe is the produced recipe.
	Recipe recipe.Recipe `json:"recipe,omitempty"`
}

// RecipeRequest constructs an add-recipe request for the subject with name sname, arch ID arch, and recipe r.
func RecipeRequest(sname string, arch id.ID, r recipe.Recipe) Request {
	return Request{Name: sname, Recipe: &Recipe{Arch: arch, Recipe: r}}
}

// Run is a request to add the given run result to the named subject.
type Run struct {
	// CompilerID is the ID of the compiler that produced this result.
	CompilerID id.ID `json:"compiler_id,omitempty"`

	// Run is the run result.
	Result compilation.RunResult `json:"result,omitempty"`
}

// RunRequest constructs an add-run request for the compilation with name name and result r.
func RunRequest(name compilation.Name, r compilation.RunResult) Request {
	return Request{Name: name.SubjectName, Run: &Run{CompilerID: name.CompilerID, Result: r}}
}
