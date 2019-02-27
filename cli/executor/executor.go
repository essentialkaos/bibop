package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v10/fmtc"
	"pkg.re/essentialkaos/ek.v10/fmtutil"
	"pkg.re/essentialkaos/ek.v10/fsutil"
	"pkg.re/essentialkaos/ek.v10/log"
	"pkg.re/essentialkaos/ek.v10/sliceutil"
	"pkg.re/essentialkaos/ek.v10/strutil"
	"pkg.re/essentialkaos/ek.v10/system"
	"pkg.re/essentialkaos/ek.v10/terminal/window"
	"pkg.re/essentialkaos/ek.v10/tmp"

	"github.com/essentialkaos/bibop/action"
	"github.com/essentialkaos/bibop/output"
	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const MAX_STORAGE_SIZE = 2 * 1024 * 1024 // 2 MB

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

// ////////////////////////////////////////////////////////////////////////////////// //

var handlers = map[string]action.Handler{
	"wait":            action.Wait,
	"sleep":           action.Wait,
	"chdir":           action.Chdir,
	"perms":           action.Perms,
	"owner":           action.Owner,
	"exist":           action.Exist,
	"readable":        action.Readable,
	"writable":        action.Writable,
	"executable":      action.Executable,
	"dir":             action.Dir,
	"empty":           action.Empty,
	"empty-dir":       action.EmptyDir,
	"checksum":        action.Checksum,
	"checksum-read":   action.ChecksumRead,
	"file-contains":   action.FileContains,
	"copy":            action.Copy,
	"move":            action.Move,
	"touch":           action.Touch,
	"mkdir":           action.Mkdir,
	"remove":          action.Remove,
	"chmod":           action.Chmod,
	"process-works":   action.ProcessWorks,
	"wait-pid":        action.WaitPID,
	"wait-fs":         action.WaitFS,
	"connect":         action.Connect,
	"app":             action.App,
	"env":             action.Env,
	"user-exist":      action.UserExist,
	"user-id":         action.UserID,
	"user-gid":        action.UserGID,
	"user-group":      action.UserGroup,
	"user-shell":      action.UserShell,
	"user-home":       action.UserHome,
	"group-exist":     action.GroupExist,
	"group-id":        action.GroupID,
	"service-present": action.ServicePresent,
	"service-enabled": action.ServiceEnabled,
	"service-works":   action.ServiceWorks,
	"http-status":     action.HTTPStatus,
	"http-header":     action.HTTPHeader,
	"http-contains":   action.HTTPContains,
	"lib-loaded":      action.LibLoaded,
	"lib-header":      action.LibHeader,
	"lib-config":      action.LibConfig,
}

var temp *tmp.Temp
var tempDir string

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
func (e *Executor) Validate(r *recipe.Recipe, tags []string) error {
	err := checkRecipeWorkingDir(r.Dir)

	if err != nil {
		return err
	}

	return checkRecipePriveleges(r.RequireRoot)
}

// Run run recipe on given executor
func (e *Executor) Run(r *recipe.Recipe, tags []string) bool {
	printBasicRecipeInfo(e, r)
	logBasicRecipeInfo(e, r)

	printSeparator("ACTIONS", e.quiet)

	cwd, _ := os.Getwd()

	processRecipe(e, r, tags)

	os.Chdir(cwd)

	printSeparator("RESULTS", e.quiet)

	printResultInfo(e)
	logResultInfo(e)
	cleanTempData()

	return e.fails == 0
}

// ////////////////////////////////////////////////////////////////////////////////// //

// processRecipe execute commands in recipe
func processRecipe(e *Executor, r *recipe.Recipe, tags []string) {
	e.start = time.Now()
	e.skipped = len(r.Commands)

	for index, command := range r.Commands {
		if r.LockWorkdir {
			os.Chdir(r.Dir) // Set current dir to working dir for every command
		}

		if command.Tag != "" && !sliceutil.Contains(tags, command.Tag) {
			e.skipped--
			continue
		}

		printCommandHeader(e, command)

		err := runCommand(e, command)

		if index+1 != len(r.Commands) {
			fmtc.NewLine()
		}

		if err != nil {
			e.logger.Info("(command %d) %v", index+1, err) // We don't use logger.Error because we log only errors
			e.fails++

			if r.FastFinish {
				break
			}
		} else {
			e.passes++
			e.skipped--
		}
	}
}

// runCommand executes command and all actions
func runCommand(e *Executor, c *recipe.Command) error {
	var err error
	var cmd *exec.Cmd
	var input io.Writer

	outputStore := output.NewStore(MAX_STORAGE_SIZE)

	if c.Cmdline != "-" {
		cmd, input, err = execCommand(c, outputStore)

		if err != nil {
			if !e.quiet {
				fmtc.Printf("  {r}%v{!}\n", err)
			}

			return err
		}
	}

	for index, action := range c.Actions {
		if !e.quiet {
			renderTmpMessage(
				"  {s-}┖╴{!} {s~-}● {!}"+formatActionName(action)+" {s}%s{!} {s-}[%s]{!}",
				formatActionArgs(action),
				formatDuration(time.Since(e.start), false),
			)
		}

		err = runAction(action, cmd, input, outputStore)

		if !e.quiet {
			if err != nil {
				renderTmpMessage("  {s-}┖╴{!} {r}✖ {!}"+formatActionName(action)+" {s}%s{!}", formatActionArgs(action))
				fmtc.NewLine()
				fmtc.Printf("     {r}%v{!}\n", err)
			} else {
				if index+1 == len(c.Actions) {
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

// execCommand executes command
func execCommand(c *recipe.Command, outputStore *output.Store) (*exec.Cmd, io.Writer, error) {
	var cmd *exec.Cmd

	if c.User == "" {
		cmdArgs := c.GetCmdlineArgs()
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	} else {
		if !system.IsUserExist(c.User) {
			return nil, nil, fmt.Errorf("Can't execute the command: user %s doesn't exist on the system", c.User)
		}

		cmd = exec.Command("/sbin/runuser", "-s", "/bin/bash", c.User, "-c", c.GetCmdline())
	}

	input, _ := cmd.StdinPipe()

	connectOutputStore(cmd, outputStore)

	err := cmd.Start()

	if err != nil {
		return nil, nil, err
	}

	go cmd.Wait()

	return cmd, input, nil
}

// logBasicRecipeInfo print path to recipe and working dir to log file
func logBasicRecipeInfo(e *Executor, r *recipe.Recipe) {
	recipeFile, _ := filepath.Abs(r.File)
	workingDir, _ := filepath.Abs(r.Dir)

	e.logger.Aux(strings.Repeat("-", 80))
	e.logger.Info(
		"Recipe: %s | Working dir: %s | Unsafe actions: %t | Require root: %t | Fast finish: %t | Lock Workdir: %t",
		recipeFile, workingDir, r.UnsafeActions, r.RequireRoot, r.FastFinish, r.LockWorkdir,
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

	fmtc.Printf("  {*}%-15s{!} %s\n", "Recipe file:", recipeFile)
	fmtc.Printf("  {*}%-15s{!} %s\n", "Working dir:", workingDir)

	printRecipeOptionFlag("Unsafe actions", r.UnsafeActions)
	printRecipeOptionFlag("Require root", r.RequireRoot)
	printRecipeOptionFlag("Fast finish", r.FastFinish)
	printRecipeOptionFlag("Lock workdir", r.LockWorkdir)
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
	case c.Cmdline != "-" && c.Description == "" && c.User != "":
		renderMessage("  {c*}[%s]{!} {c-}%s{!}", c.User, c.Cmdline)
	case c.Cmdline != "-" && c.Description != "" && c.User != "":
		renderMessage(
			"  {*}%s{!} {s}→{!} {c*}[%s]{!} {c-}%s{!}",
			c.Description, c.User, c.GetCmdline(),
		)
	default:
		renderMessage(
			"  {*}%s{!} {s}→{!} {c-}%s{!}",
			c.Description, c.GetCmdline(),
		)
	}

	fmtc.NewLine()
}

// runAction run action on command
func runAction(a *recipe.Action, cmd *exec.Cmd, input io.Writer, outputStore *output.Store) error {
	var err error
	var tmpDir string

	if a.Name == recipe.ACTION_BACKUP || a.Name == recipe.ACTION_BACKUP_RESTORE {
		tmpDir, err = getTempDir()

		if err != nil {
			return err
		}
	}

	switch a.Name {
	case recipe.ACTION_EXIT:
		return action.Exit(a, cmd)
	case recipe.ACTION_EXPECT:
		return action.Expect(a, outputStore)
	case recipe.ACTION_PRINT:
		return action.Input(a, input, outputStore)
	case recipe.ACTION_WAIT_OUTPUT:
		return action.WaitOutput(a, outputStore)
	case recipe.ACTION_OUTPUT_CONTAINS:
		return action.OutputContains(a, outputStore)
	case recipe.ACTION_OUTPUT_MATCH:
		return action.OutputMatch(a, outputStore)
	case recipe.ACTION_OUTPUT_TRIM:
		return action.OutputTrim(a, outputStore)
	case recipe.ACTION_BACKUP:
		return action.Backup(a, tmpDir)
	case recipe.ACTION_BACKUP_RESTORE:
		return action.BackupRestore(a, tmpDir)
	}

	handler, ok := handlers[a.Name]

	if !ok {
		return fmt.Errorf("Can't find handler for action %s", a.Name)
	}

	return handler(a)
}

// connectOutputStore create output store
func connectOutputStore(cmd *exec.Cmd, outputStore *output.Store) {
	stdoutReader, _ := cmd.StdoutPipe()
	stderrReader, _ := cmd.StderrPipe()

	go func(stdout, stderr io.Reader, outputStore *output.Store) {
		for {
			if outputStore.Clear {
				outputStore.Purge()
			}

			outputStore.Stdout.Write(ioutil.ReadAll(stdout))
			outputStore.Stderr.Write(ioutil.ReadAll(stderr))

			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				return
			}
		}
	}(stdoutReader, stderrReader, outputStore)
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

// printRecipeOptionFlag formats and prints option value
func printRecipeOptionFlag(name string, flag bool) {
	fmtc.Printf("  {*}%-15s{!} ", name+":")

	switch flag {
	case true:
		fmtc.Println("Yes")
	case false:
		fmtc.Println("No")
	}
}

// getTempDir return path to directory for temporary data
func getTempDir() (string, error) {
	if tempDir != "" {
		return tempDir, nil
	}

	var err error

	temp, err = tmp.NewTemp()

	if err != nil {
		return "", fmt.Errorf("Can't create directory for temporary data: %v", err)
	}

	tempDir, err = temp.MkDir("bibop")

	if err != nil {
		return "", fmt.Errorf("Can't create directory for temporary data: %v", err)
	}

	return tempDir, nil
}

// cleanTempData removes temporary data
func cleanTempData() {
	if temp == nil {
		return
	}

	temp.Clean()
}
