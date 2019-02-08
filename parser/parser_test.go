package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"

	. "pkg.re/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type ParseSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&ParseSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *ParseSuite) TestGlobalErrors(c *C) {
	recipe, err := Parse("../testdata/test0.recipe")

	c.Assert(err, NotNil)
	c.Assert(recipe, IsNil)

	recipe, err = Parse("../testdata/test2.recipe")

	c.Assert(err, NotNil)
	c.Assert(recipe, IsNil)

	recipe, err = Parse("../testdata/test3.recipe")

	c.Assert(err, NotNil)
	c.Assert(recipe, IsNil)

	recipe, err = Parse("../testdata/test4.recipe")

	c.Assert(err, NotNil)
	c.Assert(recipe, IsNil)

	recipe, err = Parse("../testdata/test5.recipe")

	c.Assert(err, NotNil)
	c.Assert(recipe, IsNil)

	recipe, err = Parse("../testdata/test6.recipe")

	c.Assert(err, NotNil)
	c.Assert(recipe, IsNil)
}

func (s *ParseSuite) TestBasicParsing(c *C) {
	recipe, err := Parse("../testdata/test1.recipe")

	c.Assert(err, IsNil)
	c.Assert(recipe, NotNil)
}

func (s *ParseSuite) TestTokenParsingErrors(c *C) {
	_, _, _, err := parseToken("abcd test")
	c.Assert(err, NotNil)

	_, _, _, err = parseToken("  abcd test")
	c.Assert(err, NotNil)

	_, _, _, err = parseToken("  perms 1 2 3")
	c.Assert(err, NotNil)

	_, _, _, err = parseToken("  perms 1")
	c.Assert(err, NotNil)

	_, _, _, err = parseToken("  ,")
	c.Assert(err, NotNil)

	_, _, _, err = parseToken("  !print 'asd'")
	c.Assert(err, NotNil)
}

func (s *ParseSuite) TestAux(c *C) {
	_, err := parseRecipeFile("../testdata/test0.recipe")
	c.Assert(err, NotNil)
}
