package subject

import "fmt"

// ExampleHarness_Paths is a testable example for Paths.
func ExampleHarness_Paths() {
	h := Harness{Dir: "foo/bar", Files: []string{"baz", "barbaz", "foobar"}}
	for _, f := range h.Paths() {
		fmt.Println(f)
	}

	// Output:
	// foo/bar/baz
	// foo/bar/barbaz
	// foo/bar/foobar
}
