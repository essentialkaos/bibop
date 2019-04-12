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
	cmd := NewCommand([]string{"echo 123"}, 0)

	c.Assert(cmd.Cmdline, Equals, "echo 123")
	c.Assert(cmd.Description, Equals, "")
	c.Assert(cmd.Actions, HasLen, 0)

	cmd = NewCommand([]string{"echo 123", "Echo command"}, 0)

	c.Assert(cmd.Cmdline, Equals, "echo 123")
	c.Assert(cmd.Description, Equals, "Echo command")
	c.Assert(cmd.Actions, HasLen, 0)
}

func (s *RecipeSuite) TestBasicRecipe(c *C) {
	r := NewRecipe("/home/user/test.recipe")

	r.Dir = "/tmp"

	r.AddVariable("service", "nginx")
	r.AddVariable("user", "nginx")

	c1 := NewCommand([]string{"{user}:echo {service}"}, 0)
	c2 := NewCommand([]string{"echo ABCD 1.53 4000", "Echo command"}, 0)

	r.AddCommand(c1, "")
	r.AddCommand(c2, "special")

	c.Assert(r.RequireRoot, Equals, true)
	c.Assert(c1.User, Equals, "nginx")
	c.Assert(c2.Tag, Equals, "special")

	a1 := &Action{"copy", []string{"file1", "file2"}, true, 0, nil}
	a2 := &Action{"touch", []string{"{service}"}, false, 0, nil}
	a3 := &Action{"print", []string{"1.53", "4000", "ABCD"}, false, 0, nil}

	c1.AddAction(a1)
	c2.AddAction(a2)

	c.Assert(c1.GetCmdlineArgs(), DeepEquals, []string{"echo", "nginx"})
	c.Assert(c2.GetCmdlineArgs(), DeepEquals, []string{"echo", "ABCD", "1.53", "4000"})

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

	c.Assert(r.GetVariable("LIBDIR"), Not(Equals), "")
	c.Assert(r.GetVariable("PYTHON_SITELIB"), Not(Equals), "")
	c.Assert(r.GetVariable("PYTHON_SITEARCH"), Not(Equals), "")

	r.GetVariable("PYTHON3_SITELIB")
	r.GetVariable("PYTHON3_SITEARCH")

	c.Assert(getPythonSitePackages("999", false), Equals, "")
}

func (s *RecipeSuite) TestIndex(c *C) {
	r := NewRecipe("/home/user/test.recipe")
	c1 := NewCommand([]string{"test"}, 0)

	c.Assert(c1.Index(), Equals, -1)

	r.AddCommand(c1, "")

	c.Assert(c1.Index(), Equals, 0)

	a1 := &Action{"abcd", []string{}, true, 0, nil}

	c.Assert(a1.Index(), Equals, -1)

	c1.AddAction(a1)

	c.Assert(a1.Index(), Equals, 0)

	c1.Actions = make([]*Action, 0)

	c.Assert(a1.Index(), Equals, -1)

	r.Commands = make([]*Command, 0)

	c.Assert(c1.Index(), Equals, -1)
}

func (s *RecipeSuite) TestCommandsParser(c *C) {
	cmd := NewCommand(nil, 0)

	c.Assert(cmd.Cmdline, Equals, "")
	c.Assert(cmd.User, Equals, "")
	c.Assert(cmd.Description, Equals, "")

	cmd = NewCommand([]string{"echo 'abcd'"}, 0)

	c.Assert(cmd.Cmdline, Equals, "echo 'abcd'")
	c.Assert(cmd.User, Equals, "")
	c.Assert(cmd.Description, Equals, "")

	cmd = NewCommand([]string{"echo 'abcd'", "My command"}, 0)

	c.Assert(cmd.Cmdline, Equals, "echo 'abcd'")
	c.Assert(cmd.User, Equals, "")
	c.Assert(cmd.Description, Equals, "My command")

	cmd = NewCommand([]string{"nobody:echo 'abcd'", "My command"}, 0)

	c.Assert(cmd.Cmdline, Equals, "echo 'abcd'")
	c.Assert(cmd.User, Equals, "nobody")
	c.Assert(cmd.Description, Equals, "My command")
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

	c.Assert(renderVars(nil, "{abcd}"), Equals, "{abcd}")
	c.Assert(renderVars(r, "{abcd}"), Equals, "{abcd}")
	c.Assert(renderVars(r, "{test}.{test}"), Equals, "ABC.ABC")
}

// ////////////////////////////////////////////////////////////////////////////////// //
