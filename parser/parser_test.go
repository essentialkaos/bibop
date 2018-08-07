package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"io/ioutil"
	"testing"

	. "pkg.re/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _DATA = `
# This is comment
dir "/tmp"
unsafe-paths true

command "echo" "Simple echo command"
  exit 1
`

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type ParseSuite struct {
	TmpDir string
}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&ParseSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *ParseSuite) SetUpSuite(c *C) {
	s.TmpDir = c.MkDir()

	err := ioutil.WriteFile(s.TmpDir+"/test.recipe", []byte(_DATA), 0644)

	if err != nil {
		c.Fatal(err.Error())
	}
}

func (s *ParseSuite) TestBasicParsing(c *C) {
	recipe, err := Parse(s.TmpDir + "/test.recipe")

	c.Assert(err, IsNil)
	c.Assert(recipe, NotNil)

	recipe, err = Parse(s.TmpDir + "/test1.recipe")

	c.Assert(err, NotNil)
	c.Assert(recipe, IsNil)
}

func (s *ParseSuite) TestTokenParsingErrors(c *C) {
	_, _, err := parseToken("abcd test")
	c.Assert(err, NotNil)

	_, _, err = parseToken("  abcd test")
	c.Assert(err, NotNil)

	_, _, err = parseToken("  perms 1 2 3")
	c.Assert(err, NotNil)

	_, _, err = parseToken("  perms 1")
	c.Assert(err, NotNil)
}
