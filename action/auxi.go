package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/essentialkaos/bibop/recipe"
	"github.com/essentialkaos/ek/v12/strutil"
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
var escapeCharRegex = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

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

// Tail return the last lines of output
func (c *OutputContainer) Tail(lines int) string {
	if c.IsEmpty() {
		return ""
	}

	if c.buf.Len() < lines+2 {
		return c.String()
	}

	data := c.buf.Bytes()
	line := 0

	for i := len(data) - 2; i >= 0; i-- {
		if data[i] == '\n' {
			line++
		}

		if line == lines {
			return strings.Trim(string(data[i+1:]), " \n\r")
		}
	}

	return strings.Trim(string(data), " \n\r")
}

// IsEmpty returns true if container is empty
func (c *OutputContainer) IsEmpty() bool {
	return c == nil || c.buf == nil || c.buf.Len() == 0
}

// Purge clears data
func (c *OutputContainer) Purge() {
	if c == nil || c.buf == nil {
		return
	}

	c.buf.Reset()
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

// fmtHash format SHA-256 hash
func fmtHash(v string) string {
	if len(v) < 64 {
		return v
	}

	return strutil.Head(v, 7) + "…" + strutil.Tail(v, 7)
}

// sanitizeData removes escape characters
func sanitizeData(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte("\r"), nil)
	return escapeCharRegex.ReplaceAll(data, nil)
}
