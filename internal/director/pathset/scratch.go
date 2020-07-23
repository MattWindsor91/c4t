// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

// Scratch contains the pre-computed paths for a machine run.
type Scratch struct {
	// DirFuzz is the directory to which fuzzed subjects will be output.
	DirFuzz string
	// DirLift is the directory to which lifter outputs will be written.
	DirLift string
	// DirPlan is the directory to which plans will be written.
	DirPlan string
	// DirRun is the directory into which act-tester-mach output will go.
	DirRun string
}

// NewScratch creates a machine pathset rooted at root.
func NewScratch(root string) *Scratch {
	return &Scratch{
		DirFuzz: filepath.Join(root, segFuzz),
		DirLift: filepath.Join(root, segLift),
		DirPlan: filepath.Join(root, segPlan),
		DirRun:  filepath.Join(root, segRun),
	}
}

// PlanForStage gets the path to the plan file for stage s.
// Note that neither Prepare nor this method create or otherwise access the plan file.
func (p *Scratch) PlanForStage(s stage.Stage) string {
	file := fmt.Sprintf("plan.%s%s", strings.ToLower(s.String()), plan.Ext)
	return filepath.Join(p.DirPlan, file)
}
