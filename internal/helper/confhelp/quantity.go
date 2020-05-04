// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package confhelp

import (
	"log"
	"reflect"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// GenericOverride substitutes any quantities in new that are non-zero for those in *old (which must be a pointer).
func GenericOverride(old, new interface{}) {
	qv := reflect.ValueOf(old).Elem()
	nv := reflect.ValueOf(new)

	nf := nv.NumField()
	for i := 0; i < nf; i++ {
		k := nv.Field(i)
		if !k.IsZero() {
			qv.Field(i).Set(k)
		}
	}
}

// LogWorkers dumps the number of workers configured by nworkers to the logger l.
func LogWorkers(l *log.Logger, nworkers int) {
	l.Println("running across", iohelp.PluralQuantity(nworkers, "worker", "", "s"))
}
