package output

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"regexp"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Store it is storage for stdout and stderr data
type Store struct {
	Stdout *Container
	Stderr *Container
	Clear  bool
}

type Container struct {
	buf  *bytes.Buffer
	size int
}

// ////////////////////////////////////////////////////////////////////////////////// //

// escapeCharRegex is regexp for searching escape characters
var escapeCharRegex = regexp.MustCompile(`\x1b\[[0-9\;]+m`)

// ////////////////////////////////////////////////////////////////////////////////// //

func NewStore(size int) *Store {
	return &Store{
		&Container{size: size},
		&Container{size: size},
		false,
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Write writes data into buffer
func (c *Container) Write(data []byte, _ error) {
	if len(data) == 0 {
		return
	}

	if c.buf == nil {
		c.buf = bytes.NewBuffer(nil)
	}

	dataLen := len(data)

	if dataLen >= c.size {
		c.buf.Reset()
		data = data[len(data)-c.size:]
	}

	if c.buf.Len()+dataLen > c.size {
		c.buf.Next(c.buf.Len() + dataLen - c.size)
	}

	c.buf.Write(sanitizeData(data))
}

// Bytes returns data as a byte slice
func (c *Container) Bytes() []byte {
	if c.buf == nil {
		return []byte{}
	}

	return c.buf.Bytes()
}

// String return data as a string
func (c *Container) String() string {
	if c.buf == nil {
		return ""
	}

	return c.buf.String()
}

// IsEmpty returns true if container is empty
func (c *Container) IsEmpty() bool {
	return c.buf == nil || c.buf.Len() == 0
}

// Purge clears data
func (c *Container) Purge() {
	if c.buf != nil {
		c.buf.Reset()
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// HasData returns true if store contains any amount of data
func (s *Store) HasData() bool {
	return !s.Stdout.IsEmpty() || !s.Stderr.IsEmpty()
}

// Purge clears all data
func (s *Store) Purge() {
	s.Stdout.Purge()
	s.Stderr.Purge()
}

// ////////////////////////////////////////////////////////////////////////////////// //

// sanitizeData removes escape characters
func sanitizeData(data []byte) []byte {
	return escapeCharRegex.ReplaceAll(data, nil)
}
