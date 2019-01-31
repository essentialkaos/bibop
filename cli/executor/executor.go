package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v10/fmtc"
	"pkg.re/essentialkaos/ek.v10/fmtutil"
	"pkg.re/essentialkaos/ek.v10/fsutil"
	"pkg.re/essentialkaos/ek.v10/log"
	"pkg.re/essentialkaos/ek.v10/pluralize"
	"pkg.re/essentialkaos/ek.v10/system"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Executor is executor struct
type Executor struct {
	quiet  bool
	start  time.Time
	passes int
	fails  int
	logger *log.Logger
}

// ////////////////////////////////////////////////////////////////////////////////// //

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

// SetupLogger setup logger
func (e *Executor) SetupLogger(file string) error {
	logger, err := log.New(file, 0644)

	if err != nil {
		return err
	}

	e.logger = logger

	return nil
}

// Validate validate recipe
func (e *Executor) Validate(r *recipe.Recipe) error {
	err := checkRecipeWorkingDir(r.Dir)

	if err != nil {
		return err
	}

	return checkRecipePriveleges(r.RequireRoot)
}

// Run run recipe on given executor
func (e *Executor) Run(r *recipe.Recipe) bool {
	printBasicRecipeInfo(e, r)
	logBasicRecipeInfo(e, r)

	e.start = time.Now()

	fsutil.Push(r.Dir)

	for index, command := range r.Commands {
		printCommandHeader(e, command)

		err := runCommand(e, command)

		if err != nil {
			// We don't use logger.Error because we log only errors
			e.logger.Info("(command %-2d) %v", index+1, err)
			e.fails++
		} else {
			e.passes++
		}
	}

	fsutil.Pop()

	printResultInfo(e)
	logResultInfo(e)

	return e.fails == 0
}

// newOutputStore create new output store
func newOutputStore() *outputStore {
	return &outputStore{data: bytes.NewBuffer(nil)}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Shrink clear output data
func (os *outputStore) Shrink() {
	os.data.Reset()
	os.clear = false
}

// String return output data as string
func (os *outputStore) String() string {
	return os.data.String()
}

// ////////////////////////////////////////////////////////////////////////////////// //

// runCommand run command
func runCommand(e *Executor, c *recipe.Command) error {
	var (
		err         error
		cmd         *exec.Cmd
		stdinWriter io.WriteCloser
		output      *outputStore
	)

	totalActions := len(c.Actions)

	if c.Cmdline != "-" {
		fullCmd := c.GetCommand()
		cmd = exec.Command(fullCmd[0], fullCmd[1:]...)
		stdinWriter, _ = cmd.StdinPipe()
		output = createOutputStore(cmd)

		err = cmd.Start()

		if err != nil {
			return err
		}

		go cmd.Wait()
	}

	for index, action := range c.Actions {
		if !e.quiet {
			fmtc.TPrintf(
				"  {s-}└{!} {s~-}● {!}%s {s}%s{!} {s-}[%s]{!}",
				action.Name, formatActionArgs(action),
				formatDuration(time.Since(e.start)),
			)
		}

		if action.Name == "exit" {
			err = actionExit(action, cmd)
		} else {
			err = runAction(action, output, stdinWriter)
		}

		if !e.quiet {
			if err != nil {
				fmtc.TPrintf("  {s-}└{!} {r}✖ {!}%s {r}%s{!}\n\n", action.Name, formatActionArgs(action))
			} else {
				if index+1 == totalActions {
					fmtc.TPrintf("  {s-}└{!} {g}✔ {!}%s {s}%s{!}\n\n", action.Name, formatActionArgs(action))
				} else {
					fmtc.TPrintf("  {s-}├{!} {g}✔ {!}%s {s}%s{!}\n", action.Name, formatActionArgs(action))
				}
			}
		}

		if err != nil {
			return fmt.Errorf("(action %-2d) %v", index+1, err)
		}
	}

	return nil
}

// logBasicRecipeInfo print path to recipe and working dir to log file
func logBasicRecipeInfo(e *Executor, r *recipe.Recipe) {
	e.logger.Aux(strings.Repeat("-", 80))
	e.logger.Info("Recipe: %s | Working Dir: %s", r.File, r.Dir)
}

// printResultInfo print info about finished test
func logResultInfo(e *Executor) {
	e.logger.Info(
		"Pass: %d | Fail: %d | Duration: %s",
		e.passes, e.fails, formatDuration(time.Since(e.start)),
	)
}

// printBasicRecipeInfo print path to recipe and working dir
func printBasicRecipeInfo(e *Executor, r *recipe.Recipe) {
	if e.quiet {
		return
	}

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
	if e.quiet {
		return
	}

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
func printCommandHeader(e *Executor, command *recipe.Command) {
	if e.quiet {
		return
	}

	fmtc.Printf("  ")

	if command.Description != "" {
		fmtc.Printf("{*}%s{!} → ", command.Description)
	}

	if command.Cmdline == "-" {
		fmtc.Printf("{y}<empty command>{!}")
	} else {
		fmtc.Printf("{c}%s{!}", command.Cmdline)
	}

	fmtc.NewLine()
}

// runAction run action on command
func runAction(action *recipe.Action, output *outputStore, input io.Writer) error {
	var err error

	if output != nil && input != nil {
		switch action.Name {
		case "expect":
			err = actionExpect(action, output)
			output.clear = true
		case "print", "input":
			err = actionInput(action, input)
			output.clear = true
		case "output-equal":
			err = actionOutputEqual(action, output)
		case "output-contains":
			err = actionOutputContains(action, output)
		case "output-prefix":
			err = actionOutputPrefix(action, output)
		case "output-suffix":
			err = actionOutputSuffix(action, output)
		case "output-length":
			err = actionOutputLength(action, output)
		case "output-trim":
			err = actionOutputTrim(action, output)
		}
	}

	switch action.Name {
	case "wait", "sleep":
		err = actionWait(action)
	case "perms":
		err = actionPerms(action)
	case "owner":
		err = actionOwner(action)
	case "exist":
		err = actionExist(action)
	case "readable":
		err = actionReadable(action)
	case "writable":
		err = actionWritable(action)
	case "directory":
		err = actionDirectory(action)
	case "empty":
		err = actionEmpty(action)
	case "empty-directory":
		err = actionEmptyDirectory(action)
	case "checksum":
		err = actionChecksum(action)
	case "file-contains":
		err = actionFileContains(action)
	case "copy":
		err = actionCopy(action)
	case "move":
		err = actionMove(action)
	case "touch":
		err = actionTouch(action)
	case "mkdir":
		err = actionMkdir(action)
	case "remove":
		err = actionRemove(action)
	case "chmod":
		err = actionChmod(action)
	case "process-works":
		err = actionProcessWorks(action)
	}

	return err
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

// formatActionArgs format command arguments and return it as string
func formatActionArgs(action *recipe.Action) string {
	var result string

	for index := range action.Arguments {
		arg, _ := action.GetS(index)
		result += arg

		if index+1 != len(action.Arguments) {
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

// isSafePath return true if path is save
func isSafePath(r *recipe.Recipe, path string) bool {
	if r.UnsafePaths {
		return true
	}

	var err error

	path, err = filepath.Abs(path)

	if err != nil {
		return false
	}

	if !strings.HasPrefix(path, r.Dir) {
		return false
	}

	return true
}

// checkRecipeWorkingDir checks working dir
func checkRecipeWorkingDir(dir string) error {
	switch {
	case !fsutil.IsExist(dir):
		return fmt.Errorf("Directory %s doesn't exist", dir)
	case !fsutil.IsDir(dir):
		return fmt.Errorf("%s is not a directory", dir)
	case !fsutil.IsWritable(dir):
		return fmt.Errorf("Directory %s is not writable", dir)
	case !fsutil.IsReadable(dir):
		return fmt.Errorf("Directory %s is not readable", dir)
	}

	return nil
}

// checkRecipePriveleges checks if bibop has superuser privileges
func checkRecipePriveleges(requireRoot bool) error {
	if !requireRoot {
		return nil
	}

	curUser, err := system.CurrentUser(true)

	if err != nil {
		return fmt.Errorf("Can't check user privileges: %v", err)
	}

	if !curUser.IsRoot() {
		return fmt.Errorf("This recipe require root privileges")
	}

	return nil
}
