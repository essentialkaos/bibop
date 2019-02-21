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

// OutputStore it is storage for stdout and stderr data
type OutputStore struct {
	Stdout *OutputContainer
	Stderr *OutputContainer
	Clear  bool
}

type OutputContainer struct {
	buf *bytes.Buffer
}

// ////////////////////////////////////////////////////////////////////////////////// //

// escapeCharRegex is regexp for searching escape characters
var escapeCharRegex = regexp.MustCompile(`\x1b\[[0-9\;]+m`)

// ////////////////////////////////////////////////////////////////////////////////// //

func NewOutputStore() *OutputStore {
	return &OutputStore{
		&OutputContainer{},
		&OutputContainer{},
		false,
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Write writes data into buffer
func (o *OutputContainer) Write(data []byte, _ error) {
	if len(data) == 0 {
		return
	}

	if o.buf == nil {
		o.buf = bytes.NewBuffer(nil)
	}

	o.buf.Write(sanitizeData(data))
}

// Bytes returns data as a byte slice
func (o *OutputContainer) Bytes() []byte {
	if o.buf == nil {
		return []byte{}
	}

	return o.buf.Bytes()
}

// String return data as a string
func (o *OutputContainer) String() string {
	if o.buf == nil {
		return ""
	}

	return o.buf.String()
}

// IsEmpty returns true if container is empty
func (o *OutputContainer) IsEmpty() bool {
	return o.buf == nil || o.buf.Len() == 0
}

// Purge clears data
func (o *OutputContainer) Purge() {
	if o.buf != nil {
		o.buf.Reset()
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// HasData returns true if store contains any amount of data
func (o *OutputStore) HasData() bool {
	return !o.Stdout.IsEmpty() || !o.Stderr.IsEmpty()
}

// Purge clears all data
func (o *OutputStore) Purge() {
	o.Stdout.Purge()
	o.Stderr.Purge()
}

// ////////////////////////////////////////////////////////////////////////////////// //

// sanitizeData removes escape characters
func sanitizeData(data []byte) []byte {
	return escapeCharRegex.ReplaceAll(data, nil)
}
