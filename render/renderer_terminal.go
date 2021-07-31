package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/fmtutil"
	"pkg.re/essentialkaos/ek.v12/strutil"
	"pkg.re/essentialkaos/ek.v12/terminal/window"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// TerminalRenderer renders info to terminal
type TerminalRenderer struct {
	isStarted  bool
	isFinished bool
	start      time.Time
}

// ////////////////////////////////////////////////////////////////////////////////// //

var isCI = os.Getenv("CI") != ""

// ////////////////////////////////////////////////////////////////////////////////// //

// Start prints info about started test
func (rr *TerminalRenderer) Start(r *recipe.Recipe) {
	if rr.isStarted {
		return
	}

	rr.start = time.Now()

	rr.printRecipeInfo(r)
	rr.printSeparator("ACTIONS")

	rr.isStarted = true
}

// CommandStarted prints info about started command
func (rr *TerminalRenderer) CommandStarted(c *recipe.Command) {
	prefix := "  "

	if c.Tag != "" {
		prefix += fmt.Sprintf("{s}(%s){!} ", c.Tag)
	}

	switch {
	case c.Cmdline == "-" && c.Description == "":
		rr.renderMessage(prefix + "{*}- Empty command -{!}")
	case c.Cmdline == "-" && c.Description != "":
		rr.renderMessage(prefix+"{*}%s{!}", c.Description)
	case c.Cmdline != "-" && c.Description == "":
		rr.renderMessage(prefix+"{c-}%s{!}", c.Cmdline)
	case c.Cmdline != "-" && c.Description == "" && c.User != "":
		rr.renderMessage(prefix+"{c*}[%s]{!} {c-}%s{!}", c.User, c.Cmdline)
	case c.Cmdline != "-" && c.Description != "" && c.User != "":
		rr.renderMessage(
			prefix+"{*}%s{!} {s}→{!} {c*}[%s]{!} {c-}%s{!}",
			c.Description, c.User, c.GetCmdline(),
		)
	default:
		rr.renderMessage(
			prefix+"{*}%s{!} {s}→{!} {c-}%s{!}",
			c.Description, c.GetCmdline(),
		)
	}

	fmtc.NewLine()
}

// CommandFailed prints info about failed command
func (rr *TerminalRenderer) CommandFailed(c *recipe.Command, err error) {
	fmtc.Printf("  {r}%v{!}\n", err)
}

// CommandFailed prints info about executed command
func (rr *TerminalRenderer) CommandDone(c *recipe.Command, isLast bool) {
	if !isLast {
		fmtc.NewLine()
	}
}

// ActionInProgress prints info about action in progress
func (rr *TerminalRenderer) ActionStarted(a *recipe.Action) {
	if isCI {
		return
	}

	rr.renderTmpMessage(
		"  {s-}└─{!} {s~-}●  {!}"+rr.formatActionName(a)+" {s}%s{!} {s-}[%s]{!}",
		rr.formatActionArgs(a),
		rr.formatDuration(time.Since(rr.start), false),
	)
}

// ActionFailed prints info about failed action
func (rr *TerminalRenderer) ActionFailed(a *recipe.Action, err error) {
	rr.renderTmpMessage(
		"  {s-}└─{!} {r}✖  {!}"+rr.formatActionName(a)+" {s}%s{!}",
		rr.formatActionArgs(a),
	)

	if !isCI {
		fmtc.NewLine()
	}

	fmtc.Printf("     {r}%v{!}\n", err)
}

// ActionDone prints info about successfully finished action
func (rr *TerminalRenderer) ActionDone(a *recipe.Action, isLast bool) {
	if isLast {
		rr.renderTmpMessage(
			"  {s-}└─{!} {g}✔  {!}"+rr.formatActionName(a)+" {s}%s{!}",
			rr.formatActionArgs(a),
		)
	} else {
		rr.renderTmpMessage(
			"  {s-}├─{!} {g}✔  {!}"+rr.formatActionName(a)+" {s}%s{!}",
			rr.formatActionArgs(a),
		)
	}

	if !isCI {
		fmtc.NewLine()
	}
}

// Result prints info about test results
func (rr *TerminalRenderer) Result(passes, fails int) {
	if rr.isFinished {
		return
	}

	rr.printSeparator("RESULTS")

	if passes == 0 {
		fmtc.Printf("  {*}Pass:{!} {r}%d{!}\n", passes)
	} else {
		fmtc.Printf("  {*}Pass:{!} {g}%d{!}\n", passes)
	}

	if fails == 0 {
		fmtc.Printf("  {*}Fail:{!} {g}%d{!}\n", fails)
	} else {
		fmtc.Printf("  {*}Fail:{!} {r}%d{!}\n", fails)
	}

	d := rr.formatDuration(time.Since(rr.start), true)
	d = strings.Replace(d, ".", "{s-}.", -1) + "{!}"

	fmtc.NewLine()
	fmtc.Println("  {*}Duration:{!} " + d)
	fmtc.NewLine()

	rr.isFinished = true
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printSeparator prints separator
func (rr *TerminalRenderer) printSeparator(name string) {
	fmtutil.Separator(false, name)
}

// printRecipeInfo prints path to recipe and working dir
func (rr *TerminalRenderer) printRecipeInfo(r *recipe.Recipe) {
	recipeFile, _ := filepath.Abs(r.File)
	workingDir, _ := filepath.Abs(r.Dir)

	rr.printSeparator("RECIPE")

	fmtc.Printf("  {*}%-15s{!} %s\n", "Recipe file:", recipeFile)
	fmtc.Printf("  {*}%-15s{!} %s\n", "Working dir:", workingDir)

	rr.printOptionFlag("Unsafe actions", r.UnsafeActions)
	rr.printOptionFlag("Require root", r.RequireRoot)
	rr.printOptionFlag("Fast finish", r.FastFinish)
	rr.printOptionFlag("Lock workdir", r.LockWorkdir)
	rr.printOptionFlag("Unbuffered IO", r.Unbuffer)
}

// printOptionFlag formats and prints option value
func (rr *TerminalRenderer) printOptionFlag(name string, flag bool) {
	fmtc.Printf("  {*}%-15s{!} ", name+":")

	switch flag {
	case true:
		fmtc.Println("Yes")
	case false:
		fmtc.Println("No")
	}
}

// formatDuration formats duration
func (rr *TerminalRenderer) formatDuration(d time.Duration, withMS bool) string {
	var m, s, ms time.Duration

	m = d / time.Minute
	d -= (m * time.Minute)
	s = d / time.Second
	d -= (s * time.Second)
	ms = d / time.Millisecond

	switch withMS {
	case true:
		return fmtc.Sprintf("%d:%02d.%03d", m, s, ms)
	default:
		return fmtc.Sprintf("%d:%02d", m, s)
	}
}

// renderTmpMessage prints temporary message limited by window size
func (rr *TerminalRenderer) renderTmpMessage(f string, a ...interface{}) {
	if isCI {
		fmtc.Printf(f+"\n", a...)
		return
	}

	ww := window.GetWidth()

	if ww <= 0 {
		fmtc.TPrintf(f, a...)
		return
	}

	textSize := strutil.Len(fmtc.Clean(fmt.Sprintf(f, a...)))

	if textSize < ww {
		fmtc.TPrintf(f, a...)
		return
	}

	ww--

	fmtc.TLPrintf(ww, f, a...)
	fmtc.Printf("{s}…{!}")
}

// renderMessage prints message limited by window size
func (rr *TerminalRenderer) renderMessage(f string, a ...interface{}) {
	ww := window.GetWidth()

	if ww <= 0 {
		fmtc.Printf(f, a...)
		return
	}

	textSize := strutil.Len(fmtc.Clean(fmt.Sprintf(f, a...)))

	if textSize < ww {
		fmtc.Printf(f, a...)
		return
	}

	ww--

	fmtc.LPrintf(ww, f, a...)
	fmtc.Printf("{s}…{!}")
}

// formatActionName format action name
func (rr *TerminalRenderer) formatActionName(a *recipe.Action) string {
	if a.Negative {
		return "{s}!{!}" + a.Name
	}

	return a.Name
}

// formatActionArgs format command arguments and return it as string
func (rr *TerminalRenderer) formatActionArgs(a *recipe.Action) string {
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
