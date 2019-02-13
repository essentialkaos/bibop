package recipe

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

type RecipeSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&RecipeSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *RecipeSuite) TestRecipeConstructor(c *C) {
	r := NewRecipe("/home/user/test.recipe")

	c.Assert(r.File, Equals, "/home/user/test.recipe")
	c.Assert(r.Dir, Equals, "")
	c.Assert(r.UnsafeActions, Equals, false)
	c.Assert(r.RequireRoot, Equals, false)
	c.Assert(r.Commands, HasLen, 0)
	c.Assert(r.variables, HasLen, 0)
}

func (s *RecipeSuite) TestCommandConstructor(c *C) {
	cmd := NewCommand([]string{"echo 123"})

	c.Assert(cmd.Cmdline, Equals, "echo 123")
	c.Assert(cmd.Description, Equals, "")
	c.Assert(cmd.Actions, HasLen, 0)

	cmd = NewCommand([]string{"echo 123", "Echo command"})

	c.Assert(cmd.Cmdline, Equals, "echo 123")
	c.Assert(cmd.Description, Equals, "Echo command")
	c.Assert(cmd.Actions, HasLen, 0)
}

func (s *RecipeSuite) TestActionConstructor(c *C) {
	a := NewAction("copy", []string{"file1", "file2"}, true)

	c.Assert(a.Name, Equals, "copy")
	c.Assert(a.Negative, Equals, true)
	c.Assert(a.Arguments, HasLen, 2)
}

func (s *RecipeSuite) TestBasicRecipe(c *C) {
	r := NewRecipe("/home/user/test.recipe")

	r.Dir = "/tmp"

	r.AddVariable("service", "nginx")

	c1 := NewCommand([]string{"echo {service}"})
	c2 := NewCommand([]string{"echo ABCD 1.53 4000", "Echo command"})

	r.AddCommand(c1)
	r.AddCommand(c2)

	a1 := NewAction("copy", []string{"file1", "file2"}, true)
	a2 := NewAction("touch", []string{"{service}"}, false)
	a3 := NewAction("print", []string{"1.53", "4000", "ABCD"}, false)

	c1.AddAction(a1)
	c2.AddAction(a2)

	c.Assert(c1.Arguments(), DeepEquals, []string{"echo", "nginx"})
	c.Assert(c2.Arguments(), DeepEquals, []string{"echo", "ABCD", "1.53", "4000"})

	vs, err := a1.GetS(0)
	c.Assert(vs, Equals, "file1")
	c.Assert(err, IsNil)

	vs, err = a2.GetS(0)
	c.Assert(vs, Equals, "nginx")
	c.Assert(err, IsNil)

	vf, err := a3.GetF(0)
	c.Assert(vf, Equals, 1.53)
	c.Assert(err, IsNil)

	vi, err := a3.GetI(1)
	c.Assert(vi, Equals, 4000)
	c.Assert(err, IsNil)

	vs, err = a3.GetS(99)
	c.Assert(vs, Equals, "")
	c.Assert(err, NotNil)

	vi, err = a3.GetI(99)
	c.Assert(vi, Equals, 0)
	c.Assert(err, NotNil)

	vf, err = a3.GetF(99)
	c.Assert(vf, Equals, 0.0)
	c.Assert(err, NotNil)

	vi, err = a3.GetI(2)
	c.Assert(vi, Equals, 0)
	c.Assert(err, NotNil)

	vf, err = a3.GetF(2)
	c.Assert(vf, Equals, 0.0)
	c.Assert(err, NotNil)

	c.Assert(r.GetVariable("WORKDIR"), Equals, "/tmp")
	c.Assert(r.GetVariable("TIMESTAMP"), HasLen, 10)
	c.Assert(r.GetVariable("DATE"), Not(Equals), "")
	c.Assert(r.GetVariable("HOSTNAME"), Not(Equals), "")
	c.Assert(r.GetVariable("IP"), Not(Equals), "")
}

func (s *RecipeSuite) TestVariables(c *C) {
	r := NewRecipe("/home/user/test.recipe")

	c.Assert(r.GetVariable("unknown"), Equals, "")

	r.AddVariable("test1", "abc1")
	c.Assert(r.GetVariable("test1"), Equals, "abc1")

	r.variables = nil

	c.Assert(r.SetVariable("test2", "abc2"), IsNil)
	c.Assert(r.GetVariable("test2"), Equals, "abc2")

	r.AddVariable("test1", "abc1")
	c.Assert(r.GetVariable("test1"), Equals, "abc1")

	c.Assert(r.SetVariable("test2", "abc3"), IsNil)
	c.Assert(r.GetVariable("test2"), Equals, "abc3")

	c.Assert(r.SetVariable("test1", "abc"), NotNil)

	c.Assert(r.GetVariable("unknown"), Equals, "")
}

func (s *RecipeSuite) TestAux(c *C) {
	r := &Recipe{
		variables: map[string]*Variable{"test": &Variable{"ABC", true}},
	}

	c.Assert(renderVars("{abcd}", nil), Equals, "{abcd}")
	c.Assert(renderVars("{abcd}", r), Equals, "{abcd}")
	c.Assert(renderVars("{test}.{test}", r), Equals, "ABC.ABC")
}

// ////////////////////////////////////////////////////////////////////////////////// //
