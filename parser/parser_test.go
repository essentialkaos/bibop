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

	recipe, err = Parse("../testdata/test7.recipe")

	c.Assert(err, NotNil)
	c.Assert(recipe, IsNil)

	recipe, err = Parse("../testdata/test8.recipe")

	c.Assert(err, NotNil)
	c.Assert(recipe, IsNil)
}

func (s *ParseSuite) TestBasicParsing(c *C) {
	recipe, err := Parse("../testdata/test1.recipe")

	c.Assert(err, IsNil)
	c.Assert(recipe, NotNil)

	c.Assert(recipe.File, Not(Equals), "")
	c.Assert(recipe.Dir, Not(Equals), "")
	c.Assert(recipe.UnsafeActions, Equals, true)
	c.Assert(recipe.RequireRoot, Equals, true)
	c.Assert(recipe.FastFinish, Equals, true)
	c.Assert(recipe.LockWorkdir, Equals, false)
	c.Assert(recipe.Commands, HasLen, 2)
	c.Assert(recipe.Packages, DeepEquals, []string{"package1", "package2"})

	c.Assert(recipe.Commands[0].User, Equals, "nobody")
	c.Assert(recipe.Commands[0].Tag, Equals, "")
	c.Assert(recipe.Commands[0].Cmdline, Equals, "echo")
	c.Assert(recipe.Commands[0].Description, Equals, "Simple echo command")
	c.Assert(recipe.Commands[0].Actions, HasLen, 2)

	c.Assert(recipe.Commands[1].User, Equals, "")
	c.Assert(recipe.Commands[1].Tag, Equals, "special")
	c.Assert(recipe.Commands[1].Cmdline, Equals, "echo")
	c.Assert(recipe.Commands[1].Description, Equals, "Simple echo command")
	c.Assert(recipe.Commands[1].Actions, HasLen, 1)
}

func (s *ParseSuite) TestTokenParsingErrors(c *C) {
	_, err := parseLine("abcd test")
	c.Assert(err, NotNil)

	_, err = parseLine("  abcd test")
	c.Assert(err, NotNil)

	_, err = parseLine("  perms 1 2 3")
	c.Assert(err, NotNil)

	_, err = parseLine("  perms 1")
	c.Assert(err, NotNil)

	_, err = parseLine("  ,")
	c.Assert(err, NotNil)

	_, err = parseLine("  !print 'asd'")
	c.Assert(err, NotNil)

	_, err = parseLine("  user-home abcd")
	c.Assert(err, NotNil)

	_, err = parseLine("  user-home abcd abcd abcd")
	c.Assert(err, NotNil)
}

func (s *ParseSuite) TestAux(c *C) {
	_, err := parseRecipeFile("../testdata/test0.recipe")
	c.Assert(err, NotNil)

	c.Assert(extractTag("command: \"asd\""), Equals, "")
}
