// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Code generated by "stringer -type Target"; DO NOT EDIT.

package compiler

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Exe-0]
	_ = x[Obj-1]
}

const _Target_name = "ExeObj"

var _Target_index = [...]uint8{0, 3, 6}

func (i Target) String() string {
	if i >= Target(len(_Target_index)-1) {
		return "Target(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Target_name[_Target_index[i]:_Target_index[i+1]]
}
