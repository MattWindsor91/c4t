// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import "errors"

var (
	// ErrDirEmpty occurs when the local directory filepath is empty.
	ErrDirEmpty = errors.New("local dir is empty string")

	// ErrInvokerNil occurs when the invocation getter is nil.
	ErrInvokerNil = errors.New("invocation getter nil")
)