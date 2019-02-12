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
	"pkg.re/essentialkaos/ek.v10/strutil"
	"pkg.re/essentialkaos/ek.v10/system"
	"pkg.re/essentialkaos/ek.v10/terminal/window"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Executor is executor struct
type Executor struct {
	quiet   bool
	start   time.Time
	passes  int
	fails   int
	skipped int
	logger  *log.Logger
}

// ActionHandler is action handler function
type ActionHandler func(action *recipe.Action) error

// ////////////////////////////////////////////////////////////////////////////////// //

type outputStore struct {
	data  *bytes.Buffer
	clear bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

var handlers = map[string]ActionHandler{
	"wait":            actionWait,
	"sleep":           actionWait,
	"perms":           actionPerms,
	"owner":           actionOwner,
	"exist":           actionExist,
	"readable":        actionReadable,
	"writable":        actionWritable,
	"executable":      actionExecutable,
	"directory":       actionDirectory,
	"empty":           actionEmpty,
	"empty-directory": actionEmptyDirectory,
	"checksum":        actionChecksum,
	"checksum-read":   actionChecksumRead,
	"file-contains":   actionFileContains,
	"copy":            actionCopy,
	"move":            actionMove,
	"touch":           actionTouch,
	"mkdir":           actionMkdir,
	"remove":          actionRemove,
	"chmod":           actionChmod,
	"process-works":   actionProcessWorks,
	"wait-pid":        actionWaitPID,
	"wait-fs":         actionWaitFS,
	"connect":         actionConnect,
	"app":             actionApp,
	"env":             actionEnv,
	"user-exist":      actionUserExist,
	"user-id":         actionUserID,
	"user-gid":        actionUserGID,
	"user-group":      actionUserGroup,
	"user-shell":      actionUserShell,
	"user-home":       actionUserHome,
	"group-exist":     actionGroupExist,
	"group-id":        actionGroupID,
	"service-present": actionServicePresent,
	"service-enabled": actionServiceEnabled,
	"service-works":   actionServiceWorks,
	"http-status":     actionHTTPStatus,
	"http-header":     actionHTTPHeader,
	"http-contains":   actionHTTPContains,
	"lib-loaded":      actionLibLoaded,
	"lib-header":      actionLibHeader,
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

	printSeparator("ACTIONS", e.quiet)

	e.start = time.Now()
	e.skipped = len(r.Commands)

	fsutil.Push(r.Dir)

	for index, command := range r.Commands {
		printCommandHeader(e, command)

		err := runCommand(e, command)

		e.skipped--

		if index+1 != len(r.Commands) {
			fmtc.NewLine()
		}

		if err != nil {
			// We don't use logger.Error because we log only errors
			e.logger.Info("(command %d) %v", index+1, err)
			e.fails++

			if r.FastFinish {
				break
			}
		} else {
			e.passes++
		}
	}

	fsutil.Pop()

	printSeparator("RESULTS", e.quiet)

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
		fullCmd := c.Arguments()
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
			renderTmpMessage(
				"  {s-}┖╴{!} {s~-}● {!}"+formatActionName(action)+" {s}%s{!} {s-}[%s]{!}",
				formatActionArgs(action),
				formatDuration(time.Since(e.start), false),
			)
		}

		err = runAction(action, cmd, output, stdinWriter)

		if !e.quiet {
			if err != nil {
				renderTmpMessage("  {s-}┖╴{!} {r}✖ {!}"+formatActionName(action)+" {r}%s{!}", formatActionArgs(action))
				fmtc.NewLine()
				fmtc.Printf("     {r}%v{!}\n", err)
			} else {
				if index+1 == totalActions {
					renderTmpMessage("  {s-}┖╴{!} {g}✔ {!}"+formatActionName(action)+" {s}%s{!}", formatActionArgs(action))
				} else {
					renderTmpMessage("  {s-}┠╴{!} {g}✔ {!}"+formatActionName(action)+" {s}%s{!}", formatActionArgs(action))
				}

				fmtc.NewLine()
			}
		}

		if err != nil {
			return fmt.Errorf("(action %d) %v", index+1, err)
		}
	}

	return nil
}

// logBasicRecipeInfo print path to recipe and working dir to log file
func logBasicRecipeInfo(e *Executor, r *recipe.Recipe) {
	recipeFile, _ := filepath.Abs(r.File)
	workingDir, _ := filepath.Abs(r.Dir)

	e.logger.Aux(strings.Repeat("-", 80))
	e.logger.Info(
		"Recipe: %s | Working dir: %s | Unsafe actions: %t | Require root: %t | Fast finish: %t",
		recipeFile, workingDir, r.UnsafeActions, r.RequireRoot, r.FastFinish,
	)
}

// printResultInfo print info about finished test
func logResultInfo(e *Executor) {
	e.logger.Info(
		"Pass: %d | Fail: %d | Skipped: %d | Duration: %s",
		e.passes, e.fails, e.skipped, formatDuration(time.Since(e.start), true),
	)
}

// printBasicRecipeInfo print path to recipe and working dir
func printBasicRecipeInfo(e *Executor, r *recipe.Recipe) {
	if e.quiet {
		return
	}

	recipeFile, _ := filepath.Abs(r.File)
	workingDir, _ := filepath.Abs(r.Dir)

	fmtutil.Separator(false, "RECIPE")

	fmtc.Printf("  {*}Recipe file:{!} %s\n", recipeFile)
	fmtc.Printf("  {*}Working dir:{!} %s\n", workingDir)

	fmtc.Printf("  {*}Unsafe actions:{!} ")

	if r.UnsafeActions {
		fmtc.Println("Allowed")
	} else {
		fmtc.Println("Not allowed")
	}

	fmtc.Printf("  {*}Require root:{!} ")

	if r.RequireRoot {
		fmtc.Println("Yes")
	} else {
		fmtc.Println("No")
	}

	fmtc.Printf("  {*}Fast finish:{!} ")

	if r.FastFinish {
		fmtc.Println("Yes")
	} else {
		fmtc.Println("No")
	}
}

// printResultInfo print info about finished test
func printResultInfo(e *Executor) {
	if e.quiet {
		return
	}

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

	fmtc.Printf("  {*}Skipped:{!} %d\n", e.skipped)

	duration := formatDuration(time.Since(e.start), true)
	duration = strings.Replace(duration, ".", "{s-}.", -1) + "{!}"

	fmtc.NewLine()
	fmtc.Println("  {*}Duration:{!} " + duration)
	fmtc.NewLine()
}

// printCommandHeader print header for executed command
func printCommandHeader(e *Executor, c *recipe.Command) {
	if e.quiet {
		return
	}

	switch {
	case c.Cmdline == "-" && c.Description == "":
		renderMessage("  {*}- Empty command -{!}")
	case c.Cmdline == "-" && c.Description != "":
		renderMessage("  {*}%s{!}", c.Description)
	case c.Cmdline != "-" && c.Description == "":
		renderMessage("  {c-}%s{!}", c.Cmdline)
	default:
		renderMessage(
			"  {*}%s{!} {s}→{!} {c-}%s{!}",
			c.Description,
			strings.Join(c.Arguments(), " "),
		)
	}

	fmtc.NewLine()
}

// runAction run action on command
func runAction(a *recipe.Action, cmd *exec.Cmd, output *outputStore, input io.Writer) error {
	var err error

	switch a.Name {
	case "exit":
		return actionExit(a, cmd)
	case "expect":
		err = actionExpect(a, output)
		output.clear = true
		return err
	case "print", "input":
		err = actionInput(a, input)
		output.clear = true
		return err
	case "output-equal":
		return actionOutputEqual(a, output)
	case "output-contains":
		return actionOutputContains(a, output)
	case "output-prefix":
		return actionOutputPrefix(a, output)
	case "output-suffix":
		return actionOutputSuffix(a, output)
	case "output-length":
		return actionOutputLength(a, output)
	case "output-trim":
		return actionOutputTrim(a, output)
	}

	handler, ok := handlers[a.Name]

	if !ok {
		return fmt.Errorf("Can't find handler for action %s", a.Name)
	}

	return handler(a)
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

// formatActionName format action name
func formatActionName(a *recipe.Action) string {
	if a.Negative {
		return "{s}!{!}" + a.Name
	}

	return a.Name
}

// formatActionArgs format command arguments and return it as string
func formatActionArgs(a *recipe.Action) string {
	var result string

	for index := range a.Arguments {
		arg, _ := a.GetS(index)
		result += arg

		if index+1 != len(a.Arguments) {
			result += " "
		}
	}

	return result
}

// formatDuration format duration
func formatDuration(d time.Duration, withMS bool) string {
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

// checkPathSafety return true if path is save
func checkPathSafety(r *recipe.Recipe, path string) (bool, error) {
	if r.UnsafeActions {
		return true, nil
	}

	targetPath, err := filepath.Abs(path)

	if err != nil {
		return false, err
	}

	workingDir, err := filepath.Abs(r.Dir)

	if err != nil {
		return false, err
	}

	return strings.HasPrefix(targetPath, workingDir), nil
}

// checkRecipeWorkingDir checks working dir
func checkRecipeWorkingDir(dir string) error {
	switch {
	case !fsutil.IsExist(dir):
		return fmt.Errorf("Directory %s doesn't exist", dir)
	case !fsutil.IsDir(dir):
		return fmt.Errorf("%s is not a directory", dir)
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

// renderMessage prints message limited by window size
func renderMessage(f string, a ...interface{}) {
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

// renderTmpMessage prints temporary message limited by window size
func renderTmpMessage(f string, a ...interface{}) {
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

// printSeparator prints separator if quiet mode not enabled
func printSeparator(name string, quiet bool) {
	if quiet {
		return
	}

	fmtutil.Separator(false, name)
}
