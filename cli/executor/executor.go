package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strconv"
	"time"

	"pkg.re/essentialkaos/ek.v9/fmtc"
	"pkg.re/essentialkaos/ek.v9/fmtutil"
	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/pluralize"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Executor is executor struct
type Executor struct {
	quiet  bool
	start  time.Time
	passes int
	fails  int
}

type outputStore struct {
	data  *bytes.Buffer
	clear bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewExecutor create new executor struct
func NewExecutor(quiet bool) *Executor {
	return &Executor{quiet: quiet}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Run run recipe on given executor
func (e *Executor) Run(r *recipe.Recipe) bool {
	if !e.quiet {
		printBasicRecipeInfo(r)
	}

	e.start = time.Now()

	fsutil.Push(r.Dir)

	for _, c := range r.Commands {
		if !e.quiet {
			printCommandHeader(c)
		}

		ok := runCommand(e, r, c)

		if ok {
			e.passes++
		} else {
			e.fails++
		}
	}

	fsutil.Pop()

	if !e.quiet {
		printResultInfo(e)
	}

	return e.fails == 0
}

// newOutputStore create new ouput store
func newOutputStore() *outputStore {
	return &outputStore{data: bytes.NewBuffer(nil)}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (os *outputStore) Shrink() {
	os.data.Reset()
	os.clear = false
}

func (os *outputStore) String() string {
	return os.data.String()
}

// ////////////////////////////////////////////////////////////////////////////////// //

// runCommand run command
func runCommand(e *Executor, r *recipe.Recipe, c *recipe.Command) bool {
	fullCmd := c.GetFullCommand()
	cmd := exec.Command(fullCmd[0], fullCmd[1:]...)
	stdinWriter, _ := cmd.StdinPipe()
	output := createOutputStore(cmd)
	totalActions := len(c.Actions)

	err := cmd.Start()

	if err != nil {
		return false
	}

	var ok bool
	var t *fmtc.T

	go cmd.Wait()

	for index, action := range c.Actions {
		if !e.quiet {
			t = fmtc.NewT()
			t.Printf(
				"  {s-}└{!} {s~-}● {!}%s {s}%s{!} {s-}[%s]{!}",
				action.Name, formatArguments(action.Arguments),
				formatDuration(time.Since(e.start)),
			)
		}

		if action.Name == "exit" {
			ok = actionExit(action, cmd)
		} else {
			ok = runAction(action, output, stdinWriter)
		}

		if !e.quiet {
			if !ok {
				t.Printf("  {s-}└{!} {r}✖ {!}%s {r}%s{!}\n\n", action.Name, formatArguments(action.Arguments))
			} else {
				if index+1 == totalActions {
					t.Printf("  {s-}└{!} {g}✔ {!}%s {s}%s{!}\n\n", action.Name, formatArguments(action.Arguments))
				} else {
					t.Printf("  {s-}├{!} {g}✔ {!}%s {s}%s{!}\n", action.Name, formatArguments(action.Arguments))
				}
			}
		}

		if !ok {
			return false
		}
	}

	return true
}

// printBasicRecipeInfo print path to recipe and working dir
func printBasicRecipeInfo(r *recipe.Recipe) {
	fmtutil.Separator(false)

	fmtc.Printf(
		"  {*}Recipe:{!} %s {s-}(%s){!}\n", r.File,
		pluralize.Pluralize(len(r.Commands), "command", "commands"),
	)

	fmtc.Printf("  {*}Working Dir:{!} %s\n", r.Dir)

	fmtutil.Separator(false)
}

// printResultInfo print info about finished test
func printResultInfo(e *Executor) {
	fmtutil.Separator(true)
	fmtc.NewLine()

	if e.passes == 0 {
		fmtc.Printf("  {*}Pass:{!} {r}%d{!}\n", e.passes)
	} else {
		fmtc.Printf("  {*}Pass:{!} {g}%d{!}\n", e.passes)
	}

	if e.fails == 0 {
		fmtc.Printf("  {*}Fail:{!} {g}%d{!}\n", e.fails)
	} else {
		fmtc.Printf("  {*}Fail:{!} {r}%d{!}\n", e.fails)
	}

	fmtc.NewLine()
	fmtc.Printf("  {*}Duration:{!} %s\n", formatDuration(time.Since(e.start)))
	fmtc.NewLine()
}

// printCommandHeader print header for executed command
func printCommandHeader(command *recipe.Command) {
	if command.Description != "" {
		fmtc.Printf("  {*}%s{!} → {c}%s{!}\n", command.Description, command.Cmdline)
	} else {
		fmtc.Printf("  {c}%s{!}\n", command.Cmdline)
	}
}

// runAction run action on command
func runAction(action *recipe.Action, output *outputStore, input io.Writer) bool {
	var status bool

	switch action.Name {
	case "wait", "sleep":
		status = actionWait(action)
	case "expect":
		status = actionExpect(action, output)
		output.clear = true
	case "print", "input":
		status = actionInput(action, input)
		output.clear = true
	case "output-equal":
		status = actionOutputEqual(action, output)
	case "output-contains":
		status = actionOutputContains(action, output)
	case "output-prefix":
		status = actionOutputPrefix(action, output)
	case "output-suffix":
		status = actionOutputSuffix(action, output)
	case "output-length":
		status = actionOutputLength(action, output)
	case "output-trim":
		status = actionOutputTrim(action, output)
	case "perms":
		status = actionPerms(action)
	case "owner":
		status = actionOwner(action)
	case "exist":
		status = actionExist(action)
	case "readable":
		status = actionReadable(action)
	case "writable":
		status = actionWritable(action)
	case "directory":
		status = actionDirectory(action)
	case "empty":
		status = actionEmpty(action)
	case "empty-directory":
		status = actionEmptyDirectory(action)
	case "not-exist":
		status = actionNotExist(action)
	case "not-readable":
		status = actionNotReadable(action)
	case "not-writable":
		status = actionNotWritable(action)
	case "not-directory":
		status = actionNotDirectory(action)
	case "not-empty":
		status = actionNotEmpty(action)
	case "not-empty-directory":
		status = actionNotEmptyDirectory(action)
	case "checksum":
		status = actionChecksum(action)
	case "file-contains":
		status = actionFileContains(action)
	case "process-works":
		status = actionProcessWorks(action)
	}

	return status
}

// createOutputStore create output store
func createOutputStore(cmd *exec.Cmd) *outputStore {
	stdoutReader, _ := cmd.StdoutPipe()
	stderrReader, _ := cmd.StderrPipe()
	multiReader := io.MultiReader(stdoutReader, stderrReader)
	outputReader := bufio.NewReader(multiReader)

	output := newOutputStore()

	go func(r *bufio.Reader, s *outputStore) {
		for {
			if s.clear {
				s.Shrink()
			}

			text, err := r.ReadString('\n')

			if err != nil {
				break
			}

			s.data.WriteString(text + "\n")
		}
	}(outputReader, output)

	return output
}

// secondsToDuration convert float seconds num to time.Duration
func secondsToDuration(sec float64) time.Duration {
	return time.Duration(sec*float64(time.Millisecond)) * 1000
}

// formatArguments format command arguments and return it as string
func formatArguments(args []string) string {
	var result string

	for index, arg := range args {
		_, err := strconv.ParseFloat(arg, 64)

		if err == nil {
			result += arg
		} else {
			result += "\"" + arg + "\""
		}

		if index+1 != len(args) {
			result += " "
		}
	}

	return result
}

// formatDuration format duration
func formatDuration(d time.Duration) string {
	var m, s time.Duration

	m = d / time.Minute
	s = (d - (m * time.Minute)) / time.Second

	return fmtc.Sprintf("%d:%02d", m, s)
}
