package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"testing"

	. "pkg.re/essentialkaos/check.v1"
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
	r.Packages = []string{"pkg1", "pkg2"}

	c.Assert(r.GetPackages(), Equals, "pkg1 pkg2")

	r.AddVariable("service", "nginx")
	r.AddVariable("user", "nginx")

	c.Assert(r.Commands.Last(), IsNil)

	c1 := NewCommand([]string{"{user}:echo {service}"}, 0)
	c2 := NewCommand([]string{"echo ABCD 1.53 4000", "Echo command for service {service}"}, 0)
	c3 := NewCommand([]string{"echo 1234"}, 0)
	c4 := NewCommand([]string{"echo test"}, 0)

	r.AddCommand(c1, "", false)
	r.AddCommand(c2, "special", false)
	r.AddCommand(c3, "", false)
	r.AddCommand(c4, "", true)

	c.Assert(r.Commands.Last(), Equals, c4)
	c.Assert(r.Commands.Has(3), Equals, true)
	c.Assert(r.Commands.Has(99), Equals, false)

	c.Assert(r.RequireRoot, Equals, true)
	c.Assert(c1.User, Equals, "nginx")
	c.Assert(c2.Tag, Equals, "special")
	c.Assert(c2.Description, Equals, "Echo command for service nginx")

	c.Assert(c3.GroupID, Equals, c4.GroupID)

	a1 := &Action{Name: "copy",
		Arguments: []string{"file1", "file2"},
		Negative:  true, Line: 0, Command: nil,
	}

	a2 := &Action{
		Name:      "touch",
		Arguments: []string{"{service}"},
		Negative:  false, Line: 0, Command: nil,
	}

	a3 := &Action{
		Name:      "print",
		Arguments: []string{"1.53", "4000", "ABCD"},
		Negative:  false, Line: 0, Command: nil,
	}

	c.Assert(c1.Actions.Last(), IsNil)

	c1.AddAction(a1)
	c2.AddAction(a2)

	c.Assert(c1.GetCmdlineArgs(), DeepEquals, []string{"echo", "nginx"})
	c.Assert(c2.GetCmdlineArgs(), DeepEquals, []string{"echo", "ABCD", "1.53", "4000"})

	c.Assert(c1.Actions.Last(), Equals, a1)
	c.Assert(c1.Actions.Has(0), Equals, true)
	c.Assert(c1.Actions.Has(99), Equals, false)

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
}

func (s *RecipeSuite) TestDynamicVariables(c *C) {
	r := NewRecipe("/home/user/test.recipe")

	r.Dir = "/tmp"

	c.Assert(r.GetVariable("WORKDIR"), Equals, "/tmp")
	c.Assert(r.GetVariable("TIMESTAMP"), HasLen, 10)
	c.Assert(r.GetVariable("DATE"), Not(Equals), "")
	c.Assert(r.GetVariable("HOSTNAME"), Not(Equals), "")
	c.Assert(r.GetVariable("IP"), Not(Equals), "")

	c.Assert(r.GetVariable("LIBDIR"), Not(Equals), "")

	r.GetVariable("PYTHON_SITELIB")
	r.GetVariable("PYTHON_SITEARCH")
	r.GetVariable("PYTHON_SITELIB_LOCAL")
	r.GetVariable("PYTHON3_SITELIB")
	r.GetVariable("PYTHON3_SITELIB_LOCAL")
	r.GetVariable("PYTHON3_SITEARCH")
	r.GetVariable("PYTHON_SITEARCH_LOCAL")
	r.GetVariable("PYTHON2_SITEARCH_LOCAL")
	r.GetVariable("PYTHON3_SITEARCH_LOCAL")
	r.GetVariable("LIBDIR_LOCAL")

	c.Assert(getPythonSitePackages("999", false, false), Equals, "")

	erlangBaseDir = "/unknown"

	c.Assert(r.GetVariable("ERLANG_BIN_DIR"), Equals, "/unknown/erts/bin")

	delete(dynVarCache, "ERLANG_BIN_DIR")
	erlangBaseDir = c.MkDir()
	os.Mkdir(erlangBaseDir+"/erts-0.0.0", 0755)

	c.Assert(r.GetVariable("ERLANG_BIN_DIR"), Equals, erlangBaseDir+"/erts-0.0.0/bin")

	// Check cache
	c.Assert(r.GetVariable("ERLANG_BIN_DIR"), Equals, erlangBaseDir+"/erts-0.0.0/bin")

	err := os.Setenv("MY_TEST_VARIABLE", "TEST1234")

	if err != nil {
		c.Fatal(err.Error())
	}

	c.Assert(r.GetVariable("ENV:MY_TEST_VARIABLE"), Equals, "TEST1234")
	c.Assert(r.GetVariable("ENV:MY_TEST_VARIABLE_1"), Equals, "")
}

func (s *RecipeSuite) TestGgetPythonSitePackages(c *C) {
	prefixDir = c.MkDir()

	os.Mkdir(prefixDir+"/lib", 0755)
	os.Mkdir(prefixDir+"/lib/python3.6", 0755)
	os.Mkdir(prefixDir+"/lib64", 0755)
	os.Mkdir(prefixDir+"/lib64/python3.6", 0755)

	c.Assert(getPythonSitePackages("3", false, false), Equals, prefixDir+"/lib/python3.6/site-packages")
	c.Assert(getPythonSitePackages("3", true, false), Equals, prefixDir+"/lib64/python3.6/site-packages")

	prefixDir = "/usr"
}

func (s *RecipeSuite) TestGetLibDir(c *C) {
	prefixDir = c.MkDir()

	os.Mkdir(prefixDir+"/lib", 0755)
	c.Assert(getLibDir(false), Equals, prefixDir+"/lib")

	os.Mkdir(prefixDir+"/lib64", 0755)
	c.Assert(getLibDir(false), Equals, prefixDir+"/lib64")

	prefixDir = "/usr"
}

func (s *RecipeSuite) TestIndex(c *C) {
	r := NewRecipe("/home/user/test.recipe")
	c1 := NewCommand([]string{"test"}, 0)

	c.Assert(c1.Index(), Equals, -1)

	r.AddCommand(c1, "", false)

	c.Assert(c1.Index(), Equals, 0)

	a1 := &Action{
		Name: "abcd", Arguments: []string{},
		Negative: true, Line: 0, Command: nil,
	}

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

	k := &Command{}

	r.AddCommand(k, "", false)

	c.Assert(renderVars(nil, "{abcd}"), Equals, "{abcd}")
	c.Assert(renderVars(r, "{abcd}"), Equals, "{abcd}")
	c.Assert(renderVars(r, "{test}.{test}"), Equals, "ABC.ABC")

	c.Assert(k.GetProp("TEST"), Equals, "")
	c.Assert(k.HasProp("TEST"), Equals, false)

	k.SetProp("TEST", "ABCD")

	c.Assert(k.GetProp("TEST"), Equals, "ABCD")
	c.Assert(k.HasProp("TEST"), Equals, true)
}

func (s *RecipeSuite) TestTags(c *C) {
	r, k := &Recipe{}, &Command{}
	r.AddCommand(k, "teardown", false)

	c.Assert(r.HasTeardown(), Equals, true)

	r, k = &Recipe{}, &Command{}
	r.AddCommand(k, "", false)

	c.Assert(r.HasTeardown(), Equals, false)
}

// ////////////////////////////////////////////////////////////////////////////////// //
