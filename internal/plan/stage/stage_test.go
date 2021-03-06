// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stage_test

import (
	"fmt"
	"testing"

	"github.com/c4-project/c4t/internal/helper/testhelp"
	"github.com/c4-project/c4t/internal/plan/stage"
)

// ExampleFromString is a testable example for FromString.
func ExampleFromString() {
	s, err := stage.FromString("Plan")
	fmt.Println(s, err)

	_, err = stage.FromString("Nonsuch")
	fmt.Println(err)

	// Output:
	// Plan <nil>
	// unknown Stage: "Nonsuch"
}

// ExampleStage_String is a testable example for Stage.String.
func ExampleStage_String() {
	for i := stage.Unknown; i <= stage.Last+1; i++ {
		fmt.Println(i)
	}

	// Output:
	// Unknown
	// Plan
	// Perturb
	// Fuzz
	// Lift
	// Invoke
	// Mach
	// Compile
	// Run
	// Analyse
	// SetCompiler
	// Stage(11)
}

// ExampleStage_MarshalJSON is a runnable example for MarshalJSON.
func ExampleStage_MarshalJSON() {
	for i := stage.Unknown + 1; i <= stage.Last; i++ {
		bs, _ := i.MarshalJSON()
		fmt.Println(string(bs))
	}

	// Output:
	// "Plan"
	// "Perturb"
	// "Fuzz"
	// "Lift"
	// "Invoke"
	// "Mach"
	// "Compile"
	// "Run"
	// "Analyse"
	// "SetCompiler"
}

// TestStage_MarshalJSON_roundTrip tests Op's marshalling and unmarshalling by round-trip.
func TestStage_MarshalJSON_roundTrip(t *testing.T) {
	t.Parallel()
	for i := stage.Unknown; i <= stage.Last; i++ {
		i := i
		t.Run(i.String(), func(t *testing.T) {
			t.Parallel()
			testhelp.TestJSONRoundTrip(t, i, "round-trip Stage")
		})
	}
}
