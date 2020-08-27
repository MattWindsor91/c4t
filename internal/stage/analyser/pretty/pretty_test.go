// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pretty_test

import (
	"context"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/plan/analysis"
	"github.com/MattWindsor91/act-tester/internal/stage/analyser/pretty"
)

// ExamplePrinter_OnAnalysis is a testable example for Printer.OnAnalysis.
func ExamplePrinter_OnAnalysis() {
	p := plan.Mock()
	a, err := analysis.Analyse(context.Background(), p, 10)
	if err != nil {
		fmt.Println("analysis error:", err)
		return
	}
	pw, err := pretty.NewPrinter(pretty.ShowCompilers(true))
	if err != nil {
		fmt.Println("printer init error:", err)
		return
	}
	pw.OnAnalysis(*a)

	// Output:
	// # Compilers
	//   ## clang
	//     - style: gcc
	//     - arch: x86
	//     - opt: none
	//     - mopt: none
	//     ### Times (sec)
	//       - compile: Min 200 Avg 200 Max 200
	//       - run: Min 0 Avg 0 Max 0
	//     ### Results
	//       - Ok: 1 subject(s)
	//
	//   ## gcc
	//     - style: gcc
	//     - arch: ppc.64le.power9
	//     - opt: none
	//     - mopt: none
	//     ### Times (sec)
	//       - compile: Min 200 Avg 200 Max 200
	//       - run: Min 0 Avg 0 Max 0
	//     ### Results
	//       - Flagged: 1 subject(s)
	//       - CompileFail: 1 subject(s)
}