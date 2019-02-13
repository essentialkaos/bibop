package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"regexp"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// OutputStore it is storage for data
type OutputStore struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
	Clear  bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// escapeCharRegex is regexp for searching escape characters
var escapeCharRegex = regexp.MustCompile(`\x1b\[[0-9\;]+m`)

// ////////////////////////////////////////////////////////////////////////////////// //

// WriteStdout writes stdout data
func (o *OutputStore) WriteStdout(data []byte, _ error) {
	if len(data) == 0 {
		return
	}

	o.Stdout.Write(sanitizeData(data))
}

// WriteStderr writes stderr data
func (o *OutputStore) WriteStderr(data []byte, _ error) {
	if len(data) == 0 {
		return
	}

	o.Stderr.Write(sanitizeData(data))
}

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

// ////////////////////////////////////////////////////////////////////////////////// //

// sanitizeData removes escape characters
func sanitizeData(data []byte) []byte {
	return escapeCharRegex.ReplaceAll(data, nil)
}
