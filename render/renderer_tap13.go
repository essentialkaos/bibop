package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// TAPRenderer is Test Anything Protocol v13 compatible renderer
type TAPRenderer struct {
	index int
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Start prints info about started test
func (rr *TAPRenderer) Start(r *recipe.Recipe) {
	fmt.Println("TAP version 13")
	fmt.Printf("1..%d\n", rr.getTestCount(r))

	recipeFile, _ := filepath.Abs(r.File)
	workingDir, _ := filepath.Abs(r.Dir)

	fmt.Println("#")
	fmt.Println("# RECIPE INFO")
	fmt.Printf("# Recipe file: %s\n", recipeFile)
	fmt.Printf("# Working dir: %s\n", workingDir)
	fmt.Printf("# Unsafe actions: %t\n", r.UnsafeActions)
	fmt.Printf("# Require root: %t\n", r.RequireRoot)
	fmt.Printf("# Fast finish: %t\n", r.FastFinish)
	fmt.Printf("# Lock workdir: %t\n", r.LockWorkdir)
	fmt.Printf("# Unbuffered IO: %t\n", r.Unbuffer)

	rr.index = 1
}

// CommandStarted prints info about started command
func (rr *TAPRenderer) CommandStarted(c *recipe.Command) {
	fmt.Println("#")

	switch {
	case c.Cmdline == "-" && c.Description == "":
		fmt.Println("# - Empty command -")
	case c.Cmdline == "-" && c.Description != "":
		fmt.Printf("# %s\n", c.Description)
	case c.Cmdline != "-" && c.Description == "":
		fmt.Printf("# %s\n", c.Cmdline)
	case c.Cmdline != "-" && c.Description == "" && c.User != "":
		fmt.Printf("# [%s] %s\n", c.User, c.Cmdline)
	case c.Cmdline != "-" && c.Description != "" && c.User != "":
		fmt.Printf("# %s -> [%s] %s\n", c.Description, c.User, c.GetCmdline())
	default:
		fmt.Printf("# %s -> %s\n", c.Description, c.GetCmdline())
	}
}

// CommandFailed prints info about failed command
func (rr *TAPRenderer) CommandFailed(c *recipe.Command, err error) {
	fmt.Printf("Bail out! %v\n", err)
}

// CommandFailed prints info about executed command
func (rr *TAPRenderer) CommandDone(c *recipe.Command, isLast bool) {
	return
}

// ActionInProgress prints info about action in progress
func (rr *TAPRenderer) ActionStarted(a *recipe.Action) {
	return
}

// ActionFailed prints info about failed action
func (rr *TAPRenderer) ActionFailed(a *recipe.Action, err error) {
	fmt.Printf(
		"not ok %d - %s %s\n",
		rr.index,
		rr.formatActionName(a),
		rr.formatActionArgs(a),
	)

	fmt.Printf("  %v\n", err)

	rr.index++
}

// ActionDone prints info about successfully finished action
func (rr *TAPRenderer) ActionDone(a *recipe.Action, isLast bool) {
	fmt.Printf(
		"ok %d - %s %s\n",
		rr.index,
		rr.formatActionName(a),
		rr.formatActionArgs(a),
	)

	rr.index++
}

// Result prints info about test results
func (rr *TAPRenderer) Result(passes, fails int) {
	return
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getTestCount returns number of all tests in recipe
func (rr *TAPRenderer) getTestCount(r *recipe.Recipe) int {
	var num int

	for _, cmd := range r.Commands {
		num += len(cmd.Actions)
	}

	return num
}

// formatActionName format action name
func (rr *TAPRenderer) formatActionName(a *recipe.Action) string {
	if a.Negative {
		return "!" + a.Name
	}

	return a.Name
}

// formatActionArgs format command arguments and return it as string
func (rr *TAPRenderer) formatActionArgs(a *recipe.Action) string {
	var result string

	for index := range a.Arguments {
		arg, _ := a.GetS(index)

		if strings.Contains(arg, " ") {
			result += "\"" + arg + "\""
		} else {
			result += arg
		}

		if index+1 != len(a.Arguments) {
			result += " "
		}
	}

	return result
}
