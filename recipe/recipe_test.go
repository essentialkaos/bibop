package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/essentialkaos/ek/v13/timeutil"

	. "github.com/essentialkaos/check"
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
	c.Assert(r.variables.index, HasLen, 0)
}

func (s *RecipeSuite) TestCommandConstructor(c *C) {
	cmd := NewCommand([]string{"-", "Hollow command"}, 0)

	c.Assert(cmd.Cmdline, Equals, "")
	c.Assert(cmd.Description, Equals, "Hollow command")
	c.Assert(cmd.Actions, HasLen, 0)
	c.Assert(cmd.IsHollow(), Equals, true)
	c.Assert(cmd.String(), Equals, "Command{-1: Hollow command → <HOLLOW> | Actions: 0}")

	cmd = NewCommand([]string{"echo 123"}, 0)

	c.Assert(cmd.Cmdline, Equals, "echo 123")
	c.Assert(cmd.Description, Equals, "")
	c.Assert(cmd.Actions, HasLen, 0)
	c.Assert(cmd.IsHollow(), Equals, false)
	c.Assert(cmd.String(), Equals, "Command{-1: echo 123 | Actions: 0}")

	cmd = NewCommand([]string{"echo 123"}, 0)
	cmd.Tag = "special"

	c.Assert(cmd.Cmdline, Equals, "echo 123")
	c.Assert(cmd.Description, Equals, "")
	c.Assert(cmd.Actions, HasLen, 0)
	c.Assert(cmd.IsHollow(), Equals, false)
	c.Assert(cmd.String(), Equals, "Command{-1:special echo 123 | Actions: 0}")

	cmd = NewCommand([]string{"echo 123", "Echo command"}, 0)

	c.Assert(cmd.Cmdline, Equals, "echo 123")
	c.Assert(cmd.Description, Equals, "Echo command")
	c.Assert(cmd.Actions, HasLen, 0)
	c.Assert(cmd.String(), Equals, "Command{-1: Echo command → echo 123 | Actions: 0}")

	cmd = NewCommand([]string{"myapp: USER=john ID=251 echo 123", "Echo command"}, 0)

	c.Assert(cmd.Cmdline, Equals, "echo 123")
	c.Assert(cmd.User, Equals, "myapp")
	c.Assert(cmd.Env, DeepEquals, []string{"USER=john", "ID=251"})
	c.Assert(cmd.Description, Equals, "Echo command")
	c.Assert(cmd.Actions, HasLen, 0)
	c.Assert(cmd.String(), Equals, "Command{-1: Echo command → (myapp) [USER=john ID=251] echo 123 | Actions: 0}")
}

func (s *RecipeSuite) TestBasicRecipe(c *C) {
	err := os.Setenv("MY_TEST_VARIABLE", "TEST1234")

	if err != nil {
		c.Fatal(err.Error())
	}

	r := NewRecipe("/home/user/test.recipe")

	r.Dir = "/tmp"
	r.Packages = []string{"pkg1", "pkg2"}

	c.Assert(r.GetPackages(), Equals, "pkg1 pkg2")

	r.AddVariable("service", "nginx")
	r.AddVariable("user", "nginx")

	c.Assert(r.Commands.Last(), IsNil)

	err = r.AddVariable("group", "{group}1")

	c.Assert(err, DeepEquals, errors.New("Can't define variable \"group\": variable contains itself as a part of value"))

	c1 := NewCommand([]string{"{user}:USER=bob echo {service}", "Basic command"}, 0)
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
	c.Assert(c1.String(), Equals, "Command{0: Basic command → (nginx) [USER=bob] echo {service} | Actions: 0}")
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

	a4 := &Action{
		Name:      "print",
		Arguments: []string{"env_{ENV:MY_TEST_VARIABLE}"},
		Negative:  false, Line: 0, Command: nil,
	}

	a5 := &Action{
		Name:      "print",
		Arguments: []string{"env_{ENV:MY_TEST_VARIABLE_1}"},
		Negative:  false, Line: 0, Command: nil,
	}

	c.Assert(c1.Actions.Last(), IsNil)

	c1.AddAction(a1)
	c2.AddAction(a2)
	c2.AddAction(a3)
	c2.AddAction(a4)
	c2.AddAction(a5)

	c.Assert(a1.String(), Equals, "Action{0: !copy file1 file2}")

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

	vs, err = a4.GetS(0)
	c.Assert(vs, Equals, "env_TEST1234")
	c.Assert(err, IsNil)

	vs, err = a5.GetS(0)
	c.Assert(vs, Equals, "env_{ENV:MY_TEST_VARIABLE_1}")
	c.Assert(err, IsNil)
}

func (s *RecipeSuite) TestDynamicVariables(c *C) {
	r := NewRecipe("/home/user/test.recipe")

	r.Dir = "/tmp"

	c.Assert(r.GetVariable("WORKDIR", false), Equals, "/tmp")
	c.Assert(r.GetVariable("TIMESTAMP", false), HasLen, 10)
	c.Assert(r.GetVariable("HOSTNAME", false), Not(Equals), "")
	c.Assert(r.GetVariable("IP", false), Not(Equals), "")
	c.Assert(r.GetVariable("ARCH", false), Not(Equals), "")
	c.Assert(r.GetVariable("ARCH_NAME", false), Not(Equals), "")
	c.Assert(r.GetVariable("ARCH_BITS", false), Not(Equals), "")
	c.Assert(r.GetVariable("OS", false), Not(Equals), "")

	c.Assert(r.GetVariable("LIBDIR", false), Not(Equals), "")
	r.GetVariable("LIBDIR_LOCAL", false)

	erlangBaseDir = "/unknown"

	c.Assert(r.GetVariable("ERLANG_BIN_DIR", false), Equals, "/unknown/erts/bin")

	delete(dynVarCache, "ERLANG_BIN_DIR")
	erlangBaseDir = c.MkDir()
	os.Mkdir(erlangBaseDir+"/erts-0.0.0", 0755)

	c.Assert(r.GetVariable("ERLANG_BIN_DIR", false), Equals, erlangBaseDir+"/erts-0.0.0/bin")

	// Check cache
	c.Assert(r.GetVariable("ERLANG_BIN_DIR", false), Equals, erlangBaseDir+"/erts-0.0.0/bin")

	err := os.Setenv("MY_TEST_VARIABLE", "TEST1234")

	if err != nil {
		c.Fatal(err.Error())
	}

	c.Assert(r.GetVariable("ENV:MY_TEST_VARIABLE", false), Equals, "TEST1234")
	c.Assert(r.GetVariable("ENV:MY_TEST_VARIABLE_1", false), Equals, "")
	c.Assert(r.GetVariable(
		"DATE:%Y%m%d", false), Equals,
		timeutil.Format(time.Now(), "%Y%m%d"),
	)
}

func (s *RecipeSuite) TestPythonVariables(c *C) {
	r := NewRecipe("/home/user/test.recipe")

	r.GetVariable("PYTHON2_VERSION", false)
	r.GetVariable("PYTHON3_VERSION", false)
	r.GetVariable("PYTHON2_SITELIB", false)
	r.GetVariable("PYTHON3_SITELIB", false)
	r.GetVariable("PYTHON2_SITELIB_LOCAL", false)
	r.GetVariable("PYTHON3_SITELIB_LOCAL", false)
	r.GetVariable("PYTHON2_SITEARCH", false)
	r.GetVariable("PYTHON3_SITEARCH", false)
	r.GetVariable("PYTHON2_SITEARCH_LOCAL", false)
	r.GetVariable("PYTHON3_SITEARCH_LOCAL", false)
	r.GetVariable("PYTHON3_BINDING_SUFFIX", false)

	python3Bin = "_unknown_"
	c.Assert(evalPythonCode(3, "test"), Equals, "")
	c.Assert(getPythonBindingSuffix(3), Equals, "")

	python3Bin = "python3"
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

	c.Assert(r.GetVariable("unknown", false), Equals, "")

	err := r.AddVariable("test1$", "abc1")
	c.Assert(err, NotNil)

	err = r.AddVariable("test1", "abc1")
	c.Assert(err, IsNil)
	c.Assert(r.GetVariable("test1", false), Equals, "abc1")

	r.variables = &Variables{index: map[string]*Variable{}}

	c.Assert(r.SetVariable("test2", "abc2"), IsNil)
	c.Assert(r.GetVariable("test2", false), Equals, "abc2")

	r.AddVariable("test1", "abc1")
	c.Assert(r.GetVariable("test1", false), Equals, "abc1")

	c.Assert(r.SetVariable("test2", "abc3"), IsNil)
	c.Assert(r.GetVariable("test2", false), Equals, "abc3")

	c.Assert(r.SetVariable("test1", "abc"), NotNil)

	c.Assert(r.SetVariable("test3", "{test1}-{test2}"), IsNil)
	c.Assert(r.GetVariable("test3", true), Equals, "abc1-abc3")

	c.Assert(r.GetVariable("unknown", false), Equals, "")

	c.Assert(r.GetVariables(), DeepEquals, []string{"test2", "test1", "test3"})

	r.variables.index["longvar"] = &Variable{strings.Repeat("A", 300), false}
	r.variables.index["longvar:test"] = &Variable{strings.Repeat("A", 300), false}

	c.Assert(renderVars(r, "{longvar}{longvar}"), Equals, "{longvar}{longvar}")
	c.Assert(renderVars(r, "{longvar:test}{longvar:test}"), Equals, "{longvar:test}{longvar:test}")
}

func (s *RecipeSuite) TestAux(c *C) {
	r := NewRecipe("test.recipe")
	r.AddVariable("test", "ABC")

	k := NewCommand([]string{}, 0)

	r.AddCommand(k, "", false)

	c.Assert(renderVars(nil, "{abcd}"), Equals, "{abcd}")
	c.Assert(renderVars(r, "{abcd}"), Equals, "{abcd}")
	c.Assert(renderVars(r, "{test}.{test}"), Equals, "ABC.ABC")

	c.Assert(k.Data.Get("TEST"), Equals, "")
	c.Assert(k.Data.Has("TEST"), Equals, false)

	k.Data.Set("TEST", "ABCD")

	c.Assert(k.Data.Get("TEST"), Equals, "ABCD")
	c.Assert(k.Data.Has("TEST"), Equals, true)
}

func (s *RecipeSuite) TestTags(c *C) {
	r, k := &Recipe{}, &Command{}
	r.AddCommand(k, "teardown", false)

	c.Assert(r.HasTeardown(), Equals, true)

	r, k = &Recipe{}, &Command{}
	r.AddCommand(k, "", false)

	c.Assert(r.HasTeardown(), Equals, false)
}

func (s *RecipeSuite) TestNesting(c *C) {
	r := NewRecipe("/home/user/test.recipe")

	r.AddVariable("a", "{d}{d}{d}{d}{d}{d}{d}{d}{d}{d}{d}{d}")
	r.AddVariable("d", "{a}{a}{a}{a}{a}{a}{a}{a}{a}{a}{a}{a}")

	c1 := NewCommand([]string{"echo 1", "My command {d}"}, 0)

	r.AddCommand(c1, "", false)

	c.Assert(len(c1.Description) < 8192, Equals, true)
}

// ////////////////////////////////////////////////////////////////////////////////// //
