package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
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

// TAP14Renderer is Test Anything Protocol v14 compatible renderer
type TAP14Renderer struct {
	Version string

	commandFailed bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Start prints info about started test
func (rr *TAP14Renderer) Start(r *recipe.Recipe) {
	fmt.Println("TAP version 14")
	fmt.Printf("1..%d\n", len(r.Commands))

	recipeFile, _ := filepath.Abs(r.File)
	workingDir, _ := filepath.Abs(r.Dir)

	fmt.Println("")
	fmt.Printf("# RECIPE INFO | bibop %s\n", rr.Version)
	fmt.Printf("# Recipe file: %s\n", recipeFile)
	fmt.Printf("# Working dir: %s\n", workingDir)
	fmt.Printf("# Unsafe actions: %t\n", r.UnsafeActions)
	fmt.Printf("# Require root: %t\n", r.RequireRoot)
	fmt.Printf("# Fast finish: %t\n", r.FastFinish)
	fmt.Printf("# Lock workdir: %t\n", r.LockWorkdir)
	fmt.Printf("# Unbuffered IO: %t\n", r.Unbuffer)
}

// CommandStarted prints info about started command
func (rr *TAP14Renderer) CommandStarted(c *recipe.Command) {
	fmt.Println("")
	fmt.Printf("# Subtest: %s\n", rr.getCommandInfo(c))
	fmt.Printf("    1..%d\n", len(c.Actions))

	rr.commandFailed = false
}

// CommandSkipped prints info about skipped command
func (rr *TAP14Renderer) CommandSkipped(c *recipe.Command, isLast bool) {
	fmt.Println("")
	fmt.Printf("ok %d - %s # SKIP\n", c.Index()+1, rr.getCommandInfo(c))
}

// CommandFailed prints info about failed command
func (rr *TAP14Renderer) CommandFailed(c *recipe.Command, err error) {
	fmt.Printf("Bail out! %v\n", err)
}

// CommandFailed prints info about executed command
func (rr *TAP14Renderer) CommandDone(c *recipe.Command, isLast bool) {
	if rr.commandFailed {
		fmt.Printf("not ok %d - %s\n", c.Index()+1, rr.getCommandInfo(c))
	} else {
		fmt.Printf("ok %d - %s\n", c.Index()+1, rr.getCommandInfo(c))
	}
}

// ActionInProgress prints info about action in progress
func (rr *TAP14Renderer) ActionStarted(a *recipe.Action) {
	return
}

// ActionFailed prints info about failed action
func (rr *TAP14Renderer) ActionFailed(a *recipe.Action, err error) {
	fmt.Printf(
		"    not ok %d - %s %s\n",
		a.Index()+1,
		rr.formatActionName(a),
		rr.formatActionArgs(a),
	)
	fmt.Print("      ---\n")
	fmt.Printf("      message: '%v'\n", err)

	rr.commandFailed = true
}

// ActionDone prints info about successfully finished action
func (rr *TAP14Renderer) ActionDone(a *recipe.Action, isLast bool) {
	fmt.Printf(
		"    ok %d - %s %s\n",
		a.Index()+1,
		rr.formatActionName(a),
		rr.formatActionArgs(a),
	)
}

// Result prints info about test results
func (rr *TAP14Renderer) Result(passes, fails, skips int) {
	fmt.Println("")
	fmt.Printf(
		"# Passed: %d | Failed: %d | Skipped: %d\n\n",
		passes, fails, skips,
	)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getCommandInfo returns command info
func (rr *TAP14Renderer) getCommandInfo(c *recipe.Command) string {
	var info string

	if c.IsHollow() {
		if c.Description == "" {
			info = "- Empty command -"
		} else {
			info = fmt.Sprintf("%s", c.Description)
		}
	} else {
		if c.Description != "" {
			info += fmt.Sprintf("%s -> ", c.Description)
		}

		if c.User != "" {
			info += fmt.Sprintf("[%s] ", c.User)
		}

		if len(c.Env) != 0 {
			info += strings.Join(c.Env, " ") + " "
		}

		info += c.GetCmdline()
	}

	return info
}

// formatActionName format action name
func (rr *TAP14Renderer) formatActionName(a *recipe.Action) string {
	if a.Negative {
		return "!" + a.Name
	}

	return a.Name
}

// formatActionArgs format command arguments and return it as string
func (rr *TAP14Renderer) formatActionArgs(a *recipe.Action) string {
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
