package output

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"regexp"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Store it is storage for stdout and stderr data
type Store struct {
	buf  *bytes.Buffer
	size int
}

// ////////////////////////////////////////////////////////////////////////////////// //

// escapeCharRegex is regexp for searching escape characters
var escapeCharRegex = regexp.MustCompile(`\x1b\[[0-9\;]+m`)

// ////////////////////////////////////////////////////////////////////////////////// //

func NewStore(size int) *Store {
	return &Store{size: size}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Write writes data into buffer
func (s *Store) Write(data []byte) {
	if s == nil || len(data) == 0 {
		return
	}

	if s.buf == nil {
		s.buf = bytes.NewBuffer(nil)
	}

	dataLen := len(data)

	if dataLen >= s.size {
		s.buf.Reset()
		data = data[len(data)-s.size:]
	}

	if s.buf.Len()+dataLen > s.size {
		s.buf.Next(s.buf.Len() + dataLen - s.size)
	}

	s.buf.Write(sanitizeData(data))
}

// Bytes returns data as a byte slice
func (s *Store) Bytes() []byte {
	if s == nil || s.buf == nil {
		return []byte{}
	}

	return s.buf.Bytes()
}

// String return data as a string
func (s *Store) String() string {
	if s == nil || s.buf == nil {
		return ""
	}

	return s.buf.String()
}

// IsEmpty returns true if container is empty
func (s *Store) IsEmpty() bool {
	return s == nil || s.buf == nil || s.buf.Len() == 0
}

// Purge clears data
func (s *Store) Purge() {
	if s == nil || s.buf != nil {
		s.buf.Reset()
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// sanitizeData removes escape characters
func sanitizeData(data []byte) []byte {
	return escapeCharRegex.ReplaceAll(data, nil)
}
