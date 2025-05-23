package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fmtutil"
	"github.com/essentialkaos/ek/v13/strutil"
	"github.com/essentialkaos/ek/v13/terminal/tty"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	_ANIMATION_STARTED uint8 = 1
	_ANIMATION_STOP    uint8 = 2
)

// ////////////////////////////////////////////////////////////////////////////////// //

// TerminalRenderer renders info to terminal
type TerminalRenderer struct {
	start      time.Time
	curAction  *recipe.Action
	syncChan   chan uint8
	isStarted  bool
	isFinished bool

	PrintExecTime bool
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
	fmtc.NewLine()
	fmtutil.Separator(true, "ACTIONS")

	rr.isStarted = true
}

// CommandStarted prints info about started command
func (rr *TerminalRenderer) CommandStarted(c *recipe.Command) {
	fmtc.NewLine()
	rr.renderMessage("  " + rr.formatCommandInfo(c))
	fmtc.NewLine()
}

// CommandSkipped prints info about skipped command
func (rr *TerminalRenderer) CommandSkipped(c *recipe.Command, isLast bool) {
	info := fmtc.Clean(rr.formatCommandInfo(c))

	fmtc.NewLine()

	if fmtc.DisableColors {
		fmtc.Printfn("  [SKIPPED] %s", info)
	} else {
		fmtc.Printfn("  {s-}%s{!}", info)
	}
}

// CommandFailed prints info about failed command
func (rr *TerminalRenderer) CommandFailed(c *recipe.Command, err error) {
	fmtc.NewLine()
	fmtc.Printfn("  {r}%v{!}", err)
}

// CommandFailed prints info about executed command
func (rr *TerminalRenderer) CommandDone(c *recipe.Command, isLast bool) {}

// ActionInProgress prints info about action in progress
func (rr *TerminalRenderer) ActionStarted(a *recipe.Action) {
	if isCI {
		return
	}

	rr.curAction = a
	rr.syncChan = make(chan uint8)

	go rr.renderCurrentActionProgress()

	// Wait until animation started
	<-rr.syncChan
}

// ActionFailed prints info about failed action
func (rr *TerminalRenderer) ActionFailed(a *recipe.Action, err error) {
	if !isCI {
		rr.syncChan <- _ANIMATION_STOP
	}

	var execTime string

	if rr.PrintExecTime {
		execTime = fmt.Sprintf(" {s-}(%s){!}", rr.formatActionTime(a))
	}

	rr.renderTmpMessage(
		"  {s-}└─{!} {r}✖  {!}"+rr.formatActionName(a)+" {s}%s{!}"+execTime,
		rr.formatActionArgs(a),
	)

	if !isCI {
		fmtc.NewLine()
	}

	fmtc.Printfn("     {r}%v{!}", err)
}

// ActionDone prints info about successfully finished action
func (rr *TerminalRenderer) ActionDone(a *recipe.Action, isLast bool) {
	if !isCI {
		rr.syncChan <- _ANIMATION_STOP
	}

	var execTime string

	if rr.PrintExecTime {
		execTime = fmt.Sprintf(" {s-}(%s){!}", rr.formatActionTime(a))
	}

	if isLast {
		rr.renderTmpMessage(
			"  {s-}└─{!} {g}✔  {!}"+rr.formatActionName(a)+" {s}%s{!}"+execTime,
			rr.formatActionArgs(a),
		)
	} else {
		rr.renderTmpMessage(
			"  {s-}├─{!} {g}✔  {!}"+rr.formatActionName(a)+" {s}%s{!}"+execTime,
			rr.formatActionArgs(a),
		)
	}

	if !isCI {
		fmtc.NewLine()
	}
}

// Result prints info about test results
func (rr *TerminalRenderer) Result(passes, fails, skips int) {
	if rr.isFinished {
		return
	}

	fmtutil.Separator(false, "RESULTS")

	if passes == 0 {
		fmtc.Printfn("  {*}Passed:{!} {r}%d{!}", passes)
	} else {
		fmtc.Printfn("  {*}Passed:{!} {g}%d{!}", passes)
	}

	if fails == 0 {
		fmtc.Printfn("  {*}Failed:{!} {g}%d{!}", fails)
	} else {
		fmtc.Printfn("  {*}Failed:{!} {r}%d{!}", fails)
	}

	if skips != 0 {
		fmtc.Printfn("  {*}Skipped:{!} {s}%d{!}", skips)
	}

	d := rr.formatDuration(time.Since(rr.start), true)
	d = strings.ReplaceAll(d, ".", "{s-}.") + "{!}"

	fmtc.NewLine()
	fmtc.Println("  {*}Duration:{!} " + d)
	fmtc.NewLine()

	rr.isFinished = true
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printRecipeInfo prints path to recipe and working dir
func (rr *TerminalRenderer) printRecipeInfo(r *recipe.Recipe) {
	recipeFile, _ := filepath.Abs(r.File)
	workingDir, _ := filepath.Abs(r.Dir)

	fmtutil.Separator(false, "RECIPE")

	fmtc.Printfn("  {*}%-15s{!} %s", "Recipe file:", recipeFile)
	fmtc.Printfn("  {*}%-15s{!} %s", "Working dir:", workingDir)

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
		fmtc.Printfn(f, a...)
		return
	}

	ww := tty.GetWidth()

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

// renderCurrentActionProgress renders info about current action
func (rr *TerminalRenderer) renderCurrentActionProgress() {
	frame := 0

	rr.renderTmpMessage(
		"  {s-}└─{!} {s-}●{!}  {!}"+rr.formatActionName(rr.curAction)+" {s}%s{!} {s-}[%s]{!}",
		rr.formatActionArgs(rr.curAction),
		rr.formatDuration(time.Since(rr.start), false),
	)

	ticker := time.NewTicker(time.Second / 5)
	defer ticker.Stop()

	rr.syncChan <- _ANIMATION_STARTED

	for {
		select {
		case <-rr.syncChan:
			return
		case <-ticker.C:
			dot := " "
			frame++

			switch frame {
			case 2:
				dot = "{s-}●{!}"
			case 3:
				dot, frame = "{s}●{!}", 0
			}

			rr.renderTmpMessage(
				"  {s-}└─{!} "+dot+"  {!}"+rr.formatActionName(rr.curAction)+" {s}%s{!} {s-}[%s]{!}",
				rr.formatActionArgs(rr.curAction),
				rr.formatDuration(time.Since(rr.start), false),
			)
		}
	}
}

// renderMessage prints message limited by window size
func (rr *TerminalRenderer) renderMessage(f string, a ...interface{}) {
	ww := tty.GetWidth()

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

// formatCommandInfo formats command info
func (rr *TerminalRenderer) formatCommandInfo(c *recipe.Command) string {
	var info string

	if c.Tag != "" {
		info += fmt.Sprintf("{s}(%s){!} ", c.Tag)
	}

	if c.IsHollow() {
		if c.Description == "" {
			return info + "{*}- Empty command -{!}"
		}

		return info + fmt.Sprintf("{*}%s{!} ", c.Description)
	}

	if c.Description != "" {
		info += fmt.Sprintf("{*}%s{!} {s}→{!} ", c.Description)
	}

	if c.User != "" {
		info += fmt.Sprintf("{c*}[%s]{!} ", c.User)
	}

	if len(c.Env) != 0 {
		info += fmt.Sprintf("{s}%s{!} ", strings.Join(c.Env, " "))
	}

	info += fmt.Sprintf("{c-}%s{!}", c.GetCmdline())

	return info
}

// formatActionName formats action name
func (rr *TerminalRenderer) formatActionName(a *recipe.Action) string {
	if a.Negative {
		return "{^r}!{!}" + a.Name
	}

	return a.Name
}

// formatActionTime formats action execution time
func (rr *TerminalRenderer) formatActionTime(a *recipe.Action) string {
	if a.Started.IsZero() {
		return "—"
	}

	d := time.Since(a.Started)

	switch {
	case d >= time.Second:
		return fmt.Sprintf("%g s", fmtutil.Float(float64(d)/float64(time.Second)))

	case d >= time.Millisecond:
		return fmt.Sprintf("%g ms", fmtutil.Float(float64(d)/float64(time.Millisecond)))

	case d >= time.Microsecond:
		return fmt.Sprintf("%g μs", fmtutil.Float(float64(d)/float64(time.Microsecond)))

	default:
		return fmt.Sprintf("%d ns", d.Nanoseconds())
	}
}

// formatActionArgs formats command arguments and return it as string
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
