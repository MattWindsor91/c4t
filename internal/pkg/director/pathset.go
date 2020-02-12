package director

// Pathset contains the pre-computed paths used by the director.
type Pathset struct {
	// DirRoot is the root directory of this director's path set.
	DirRoot string

	CurRun RunPathset

	// DirLitmus is the directory to which litmus tests will be written.
	DirLitmus string

	// DirTrace is the directory to which traces will be written.
	DirTrace string
}

// RunPathset contains the pre-computed paths used by a run of the director.
type RunPathset struct {
	// DirRoot is the root directory of this run.
	DirRoot string

	// DirFuzz is the fuzzing directory for this run.
	DirFuzz string

	// DirLift is the lifting directory for this run.
	DirLift string
}
