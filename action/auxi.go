package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Handler is action handler function
type Handler func(action *recipe.Action) error

// Store it is storage for stdout and stderr data
type OutputContainer struct {
	buf  *bytes.Buffer
	size int
}

// ////////////////////////////////////////////////////////////////////////////////// //

// escapeCharRegex is regexp for searching escape characters
var escapeCharRegex = regexp.MustCompile(`\x1b\[[0-9\;]+m`)

// ////////////////////////////////////////////////////////////////////////////////// //

func NewOutputContainer(size int) *OutputContainer {
	return &OutputContainer{size: size}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Write writes data into buffer
func (c *OutputContainer) Write(data []byte) {
	if c == nil || len(data) == 0 {
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
func (c *OutputContainer) Bytes() []byte {
	if c == nil || c.buf == nil {
		return []byte{}
	}

	return c.buf.Bytes()
}

// String return data as a string
func (c *OutputContainer) String() string {
	if c == nil || c.buf == nil {
		return ""
	}

	return c.buf.String()
}

// IsEmpty returns true if container is empty
func (c *OutputContainer) IsEmpty() bool {
	return c == nil || c.buf == nil || c.buf.Len() == 0
}

// Purge clears data
func (c *OutputContainer) Purge() {
	if c == nil || c.buf != nil {
		c.buf.Reset()
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// checkPathSafety return true if path is save
func checkPathSafety(r *recipe.Recipe, path string) (bool, error) {
	if r.UnsafeActions {
		return true, nil
	}

	targetPath, err := filepath.Abs(path)

	if err != nil {
		return false, err
	}

	workingDir, err := filepath.Abs(r.Dir)

	if err != nil {
		return false, err
	}

	return strings.HasPrefix(targetPath, workingDir), nil
}

// fmtValue formats value
func fmtValue(v string) string {
	if v == "" {
		return `""`
	}

	return v
}

// sanitizeData removes escape characters
func sanitizeData(data []byte) []byte {
	return escapeCharRegex.ReplaceAll(data, nil)
}
