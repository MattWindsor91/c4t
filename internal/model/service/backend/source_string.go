// Code generated by "stringer -type Source"; DO NOT EDIT.

package backend

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LiftUnknown-0]
	_ = x[LiftLitmus-1]
}

const _Source_name = "LiftUnknownLiftLitmus"

var _Source_index = [...]uint8{0, 11, 21}

func (i Source) String() string {
	if i >= Source(len(_Source_index)-1) {
		return "Source(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Source_name[_Source_index[i]:_Source_index[i+1]]
}
