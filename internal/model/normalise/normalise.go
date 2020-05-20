// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package normalise provides utilities for archiving and transferring plans, corpora, and subjects.
package normalise

import (
	"errors"
	"fmt"
	"path"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// ErrCollision occurs if the normaliser tries to map two files to the same normalised path.
// Usually, this is an internal error.
var ErrCollision = errors.New("path already mapped by normaliser")

const (
	FileBin        = "a.out"
	FileCompileLog = "compile.log"
	FileOrigLitmus = "orig.litmus"
	FileFuzzLitmus = "fuzz.litmus"
	FileFuzzTrace  = "fuzz.trace"
	// DirCompiles is the normalised directory for compile results.
	DirCompiles = "compiles"
	// DirHarnesses is the normalised directory for harness results.
	DirHarnesses = "harnesses"
)

// Normaliser contains state necessary to normalise a single subject's paths.
// This is useful for archiving the subject inside a tarball, or copying it to another host.
type Normaliser struct {
	// root is the prefix to add to every normalised name.
	root string

	// Mappings contains maps from normalised names to original names.
	// (The mappings are this way around to help us notice collisions.)
	Mappings map[string]Normalisation
}

// Normalisation is a record in the normaliser's mappings.
type Normalisation struct {
	// Original is the original path.
	Original string
	// Kind is the kind of path to which this mapping belongs.
	// This exists mainly to make it possible to use a Normaliser to work out how to copy a plan to another host,
	// but only copy selective subsets of files.
	Kind NormalisationKind
}

type NormalisationKind int

const (
	// NKFuzz marks that a mapping pertains to the original litmus file.
	NKOrigLitmus NormalisationKind = iota
	// NKFuzz marks that a mapping is part of a fuzz.
	NKFuzz
	// NKCompile marks that a mapping is part of a compile.
	NKCompile
	// NKHarness marks that a mapping is part of a harness.
	NKHarness
)

// NewNormaliser constructs a new Normaliser relative to root.
func NewNormaliser(root string) *Normaliser {
	return &Normaliser{
		root:     root,
		Mappings: make(map[string]Normalisation),
	}
}

// MappingsOfKind filters this normaliser's map to only the files matching kind nk.
func (n *Normaliser) MappingsOfKind(nk NormalisationKind) map[string]string {
	fs := make(map[string]string)
	for k, v := range n.Mappings {
		if v.Kind == nk {
			fs[k] = v.Original
		}
	}
	return fs
}

// Corpus normalises mappings for each subject in c.
func (n *Normaliser) Corpus(c corpus.Corpus) (corpus.Corpus, error) {
	c2 := make(corpus.Corpus, len(c))
	for name, s := range c {
		// The aliasing of Mappings here is deliberate.
		snorm := Normaliser{root: path.Join(n.root, name), Mappings: n.Mappings}
		ns, err := snorm.Subject(s)
		if err != nil {
			return nil, fmt.Errorf("normalising %s: %w", name, err)
		}
		c2[name] = *ns
	}
	return c2, nil
}

// Subject normalises mappings from subject component files to 'normalised' names.
func (n *Normaliser) Subject(s subject.Subject) (*subject.Subject, error) {
	var err error
	s.OrigLitmus, err = n.replaceAndAdd(s.OrigLitmus, NKOrigLitmus, FileOrigLitmus)
	if s.Fuzz != nil && err == nil {
		s.Fuzz, err = n.fuzz(*s.Fuzz)
	}
	if s.Compiles != nil && err == nil {
		s.Compiles, err = n.compiles(s.Compiles)
	}
	if s.Harnesses != nil && err == nil {
		s.Harnesses, err = n.harnesses(s.Harnesses)
	}
	// No need to normalise runs
	return &s, err
}

func (n *Normaliser) fuzz(f subject.Fuzz) (*subject.Fuzz, error) {
	var err error
	if f.Files.Litmus, err = n.replaceAndAdd(f.Files.Litmus, NKFuzz, FileFuzzLitmus); err != nil {
		return nil, err
	}
	f.Files.Trace, err = n.replaceAndAdd(f.Files.Trace, NKFuzz, FileFuzzTrace)
	return &f, err
}

func (n *Normaliser) harnesses(hs map[string]subject.Harness) (map[string]subject.Harness, error) {
	nhs := make(map[string]subject.Harness, len(hs))
	for archstr, h := range hs {
		var err error
		nhs[archstr], err = n.harness(archstr, h)
		if err != nil {
			return nil, err
		}
	}
	return nhs, nil
}

func (n *Normaliser) harness(archstr string, h subject.Harness) (subject.Harness, error) {
	oldPaths := h.Paths()
	h.Dir = path.Join(n.root, DirHarnesses, archstr)
	for i, np := range h.Paths() {
		if err := n.add(oldPaths[i], np, NKHarness); err != nil {
			return h, err
		}
	}
	return h, nil
}

func (n *Normaliser) compiles(cs map[string]subject.CompileResult) (map[string]subject.CompileResult, error) {
	ncs := make(map[string]subject.CompileResult, len(cs))
	for cidstr, c := range cs {
		var err error
		ncs[cidstr], err = n.compile(cidstr, c)
		if err != nil {
			return nil, err
		}
	}

	return ncs, nil
}

func (n *Normaliser) compile(cidstr string, c subject.CompileResult) (subject.CompileResult, error) {
	var err error
	if c.Files.Bin, err = n.replaceAndAdd(c.Files.Bin, NKCompile, DirCompiles, cidstr, FileBin); err != nil {
		return c, err
	}
	c.Files.Log, err = n.replaceAndAdd(c.Files.Log, NKCompile, DirCompiles, cidstr, FileCompileLog)
	return c, err
}

// replaceAndAdd adds the path assembled by joining segs together as a mapping from opath.
// If opath is empty, this just returns ("", nil) and does no addition.
func (n *Normaliser) replaceAndAdd(opath string, nk NormalisationKind, segs ...string) (string, error) {
	if opath == "" {
		return "", nil
	}

	npath := path.Join(n.root, path.Join(segs...))
	err := n.add(opath, npath, nk)
	return npath, err
}

// add tries to add the mapping between opath and npath to the normaliser's mappings.
// It fails if there is a collision.
func (n *Normaliser) add(opath, npath string, nk NormalisationKind) error {
	if _, ok := n.Mappings[npath]; ok {
		return fmt.Errorf("%w: %q", ErrCollision, npath)
	}
	n.Mappings[npath] = Normalisation{
		Original: opath,
		Kind:     nk,
	}
	return nil
}
