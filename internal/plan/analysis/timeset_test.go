// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/plan/analysis"
)

// ExampleNewTimeSet is a runnable example for NewTimeSet.
func ExampleNewTimeSet() {
	ts := analysis.NewTimeSet(1*time.Second, 1*time.Second, 2*time.Second, 4*time.Second)
	fmt.Println("min", ts.Min)
	fmt.Println("avg", ts.Mean())
	fmt.Println("max", ts.Max)
	fmt.Println("sum", ts.Sum)
	fmt.Println("count", ts.Count)

	// Output:
	// min 1s
	// avg 2s
	// max 4s
	// sum 8s
	// count 4
}
