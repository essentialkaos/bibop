package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type OutputStore struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
	Clear  bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Shrink clears all output data
func (o *OutputStore) Shrink() {
	o.Stdout.Reset()
	o.Stderr.Reset()
	o.Clear = false
}

// HasData returns true if sttore contains any data
func (o *OutputStore) HasData() bool {
	return o.Stdout.Len() != 0 || o.Stderr.Len() != 0
}
