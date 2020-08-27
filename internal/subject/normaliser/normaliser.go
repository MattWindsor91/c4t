// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package normaliser provides utilities for archiving and transferring plans, corpora, and subjects.
package normaliser

import (
	"errors"
	"fmt"
	"path"

	"github.com/MattWindsor91/act-tester/internal/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/subject"
)

// ErrCollision occurs if the normaliser tries to map two files to the same normalised path.
// Usually, this is an internal error.
var ErrCollision = errors.New("path already mapped by normaliser")

// Normaliser contains state necessary to normalise a single subject's paths.
// This is useful for archiving the subject inside a tarball, or copying it to another host.
type Normaliser struct {
	// root is the prefix to add to every normalised name.
	root string

	// err is the first error this normaliser encountered.
	err error

	// Mappings contains maps from normalised names to original names.
	// (The mappings are this way around to help us notice collisions.)
	Mappings Map
}

// New constructs a new Normaliser relative to root.
func New(root string) *Normaliser {
	return &Normaliser{
		root:     root,
		Mappings: make(map[string]Entry),
	}
}

// Normalise normalises mappings from subject component files to 'normalised' names.
func (n *Normaliser) Normalise(s subject.Subject) (*subject.Subject, error) {
	n.err = nil

	s.Source.Path = n.replaceAndAdd(s.Source.Path, filekind.Litmus, filekind.InOrig, FileOrigLitmus)
	s.Fuzz = n.fuzz(s.Fuzz)
	s.Compiles = n.compiles(s.Compiles)
	s.Recipes = n.recipes(s.Recipes)
	// No need to normalise runs
	return &s, n.err
}

func (n *Normaliser) fuzz(of *subject.Fuzz) *subject.Fuzz {
	if of == nil {
		return nil
	}
	f := *of
	f.Litmus.Path = n.replaceAndAdd(f.Litmus.Path, filekind.Litmus, filekind.InFuzz, FileFuzzLitmus)
	f.Trace = n.replaceAndAdd(f.Trace, filekind.Trace, filekind.InFuzz, FileFuzzTrace)
	return &f
}

func (n *Normaliser) recipes(rs map[string]recipe.Recipe) map[string]recipe.Recipe {
	if rs == nil {
		return nil
	}

	nrs := make(map[string]recipe.Recipe, len(rs))
	for archstr, r := range rs {
		nrs[archstr] = n.recipe(archstr, r)
	}
	return nrs
}

func (n *Normaliser) recipe(archstr string, h recipe.Recipe) recipe.Recipe {
	oldPaths := h.Paths()
	h.Dir = RecipeDir(n.root, archstr)
	for i, np := range h.Paths() {
		n.add(oldPaths[i], np, filekind.GuessFromFile(np), filekind.InRecipe)
	}
	return h
}

func (n *Normaliser) compiles(cs map[string]compilation.CompileResult) map[string]compilation.CompileResult {
	if cs == nil {
		return nil
	}
	ncs := make(map[string]compilation.CompileResult, len(cs))
	for cidstr, c := range cs {
		ncs[cidstr] = n.compile(cidstr, c)
	}
	return ncs
}

func (n *Normaliser) compile(cidstr string, c compilation.CompileResult) compilation.CompileResult {
	c.Files.Bin = n.replaceAndAdd(c.Files.Bin, filekind.Bin, filekind.InCompile, DirCompiles, cidstr, FileBin)
	c.Files.Log = n.replaceAndAdd(c.Files.Log, filekind.Log, filekind.InCompile, DirCompiles, cidstr, FileCompileLog)
	return c
}

// replaceAndAdd adds the path assembled by joining segs together as a mapping from opath.
// If opath is empty, this just returns "" and does no addition.
func (n *Normaliser) replaceAndAdd(opath string, k filekind.Kind, l filekind.Loc, segs ...string) string {
	if n.err != nil || opath == "" {
		return ""
	}
	return n.add(opath, path.Join(n.root, path.Join(segs...)), k, l)
}

// add tries to add the mapping between opath and npath to the normaliser's mappings, returning npath.
// It fails if there is a collision.
func (n *Normaliser) add(opath, npath string, k filekind.Kind, l filekind.Loc) string {
	if _, ok := n.Mappings[npath]; ok {
		n.err = fmt.Errorf("%w: %q", ErrCollision, npath)
		return npath
	}
	n.Mappings[npath] = Entry{
		Original: opath,
		Kind:     k,
		Loc:      l,
	}
	return npath
}