package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"pkg.re/essentialkaos/ek.v12/errutil"
	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/fsutil"
	"pkg.re/essentialkaos/ek.v12/log"
	"pkg.re/essentialkaos/ek.v12/passwd"
	"pkg.re/essentialkaos/ek.v12/sliceutil"
	"pkg.re/essentialkaos/ek.v12/strutil"
	"pkg.re/essentialkaos/ek.v12/system"
	"pkg.re/essentialkaos/ek.v12/timeutil"
	"pkg.re/essentialkaos/ek.v12/tmp"

	"github.com/essentialkaos/bibop/action"
	"github.com/essentialkaos/bibop/output"
	"github.com/essentialkaos/bibop/recipe"
	"github.com/essentialkaos/bibop/render"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const MAX_STORAGE_SIZE = 8 * 1024 * 1024 // 2 MB

// ////////////////////////////////////////////////////////////////////////////////// //

// Executor is executor struct
type Executor struct {
	config  *Config     // Config
	start   time.Time   // Time when recipe execution started
	passes  int         // Number of passed commands
	fails   int         // Number of failed commands
	skipped int         // Number of skipped commands
	logger  *log.Logger // Pointer to logger
}

// ExecutorConfig contains executor configuration
type Config struct {
	Quiet          bool
	DisableCleanup bool
	ErrsDir        string
}

// ValidationConfig is config for validation
type ValidationConfig struct {
	Tags               []string
	IgnoreDependencies bool
	IgnorePrivileges   bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

var handlers = map[string]action.Handler{
	recipe.ACTION_WAIT:            action.Wait,
	recipe.ACTION_CHDIR:           action.Chdir,
	recipe.ACTION_MODE:            action.Mode,
	recipe.ACTION_OWNER:           action.Owner,
	recipe.ACTION_EXIST:           action.Exist,
	recipe.ACTION_READABLE:        action.Readable,
	recipe.ACTION_WRITABLE:        action.Writable,
	recipe.ACTION_EXECUTABLE:      action.Executable,
	recipe.ACTION_DIR:             action.Dir,
	recipe.ACTION_EMPTY:           action.Empty,
	recipe.ACTION_EMPTY_DIR:       action.EmptyDir,
	recipe.ACTION_CHECKSUM:        action.Checksum,
	recipe.ACTION_CHECKSUM_READ:   action.ChecksumRead,
	recipe.ACTION_FILE_CONTAINS:   action.FileContains,
	recipe.ACTION_COPY:            action.Copy,
	recipe.ACTION_MOVE:            action.Move,
	recipe.ACTION_TOUCH:           action.Touch,
	recipe.ACTION_MKDIR:           action.Mkdir,
	recipe.ACTION_REMOVE:          action.Remove,
	recipe.ACTION_CHMOD:           action.Chmod,
	recipe.ACTION_PROCESS_WORKS:   action.ProcessWorks,
	recipe.ACTION_WAIT_PID:        action.WaitPID,
	recipe.ACTION_WAIT_FS:         action.WaitFS,
	recipe.ACTION_CONNECT:         action.Connect,
	recipe.ACTION_APP:             action.App,
	recipe.ACTION_ENV:             action.Env,
	recipe.ACTION_ENV_SET:         action.EnvSet,
	recipe.ACTION_USER_EXIST:      action.UserExist,
	recipe.ACTION_USER_ID:         action.UserID,
	recipe.ACTION_USER_GID:        action.UserGID,
	recipe.ACTION_USER_GROUP:      action.UserGroup,
	recipe.ACTION_USER_SHELL:      action.UserShell,
	recipe.ACTION_USER_HOME:       action.UserHome,
	recipe.ACTION_GROUP_EXIST:     action.GroupExist,
	recipe.ACTION_GROUP_ID:        action.GroupID,
	recipe.ACTION_SERVICE_PRESENT: action.ServicePresent,
	recipe.ACTION_SERVICE_ENABLED: action.ServiceEnabled,
	recipe.ACTION_SERVICE_WORKS:   action.ServiceWorks,
	recipe.ACTION_HTTP_STATUS:     action.HTTPStatus,
	recipe.ACTION_HTTP_HEADER:     action.HTTPHeader,
	recipe.ACTION_HTTP_CONTAINS:   action.HTTPContains,
	recipe.ACTION_HTTP_JSON:       action.HTTPJSON,
	recipe.ACTION_HTTP_SET_AUTH:   action.HTTPSetAuth,
	recipe.ACTION_HTTP_SET_HEADER: action.HTTPSetHeader,
	recipe.ACTION_LIB_LOADED:      action.LibLoaded,
	recipe.ACTION_LIB_HEADER:      action.LibHeader,
	recipe.ACTION_LIB_CONFIG:      action.LibConfig,
	recipe.ACTION_LIB_EXIST:       action.LibExist,
	recipe.ACTION_LIB_LINKED:      action.LibLinked,
	recipe.ACTION_PYTHON_MODULE:   action.PythonModule,
	recipe.ACTION_PYTHON3_MODULE:  action.Python3Module,
}

var temp *tmp.Temp
var tempDir string

// ////////////////////////////////////////////////////////////////////////////////// //

// NewExecutor create new executor struct
func NewExecutor(cfg *Config) *Executor {
	return &Executor{config: cfg}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Validate validates recipe
func (e *Executor) Validate(r *recipe.Recipe, cfg *ValidationConfig) []error {
	errs := errutil.NewErrors()

	errs.Add(checkRecipeWorkingDir(r))
	errs.Add(checkRecipeTags(r, cfg.Tags))
	errs.Add(checkRecipeVariables(r))

	if !cfg.IgnorePrivileges {
		errs.Add(checkRecipePrivileges(r))
	}

	if !cfg.IgnoreDependencies {
		errs.Add(checkPackages(r))
	}

	if !errs.HasErrors() {
		return nil
	}

	return errs.All()
}

// Run run recipe on given executor
func (e *Executor) Run(rr render.Renderer, r *recipe.Recipe, tags []string) bool {
	rr.Start(r)

	cwd, _ := os.Getwd()

	if r.Dir != "" {
		os.Chdir(r.Dir)
	}

	processRecipe(e, rr, r, tags)

	os.Chdir(cwd)

	rr.Result(e.passes, e.fails)

	cleanTempData()
	cleanupWorkingDir(e, r.Dir)

	return e.fails == 0
}

// ////////////////////////////////////////////////////////////////////////////////// //

// processRecipe execute commands in recipe
func processRecipe(e *Executor, rr render.Renderer, r *recipe.Recipe, tags []string) {
	e.start = time.Now()
	e.skipped = len(r.Commands)

	for index, command := range r.Commands {
		if r.LockWorkdir && r.Dir != "" {
			os.Chdir(r.Dir) // Set current dir to working dir for every command
		}

		if skipCommand(command, tags) {
			e.skipped--
			continue
		}

		rr.CommandStarted(command)

		ok := runCommand(e, rr, command)

		e.skipped--

		if !ok {
			e.fails++

			if r.FastFinish {
				break
			}
		} else {
			e.passes++
		}

		rr.CommandDone(command, index+1 == len(r.Commands))

		if r.Delay > 0 {
			time.Sleep(timeutil.SecondsToDuration(r.Delay))
		}
	}
}

// runCommand executes command and all actions
func runCommand(e *Executor, rr render.Renderer, c *recipe.Command) bool {
	var err error
	var cmd *exec.Cmd
	var input io.Writer

	outputStore := output.NewStore(MAX_STORAGE_SIZE)

	if c.Cmdline != "-" {
		cmd, input, err = execCommand(c, outputStore)

		if err != nil {
			rr.CommandFailed(c, err)

			logError(e, c, nil, outputStore, err)

			return false
		}
	}

	for index, action := range c.Actions {
		rr.ActionStarted(action)

		err = runAction(action, cmd, input, outputStore)

		if err != nil {
			rr.ActionFailed(action, err)
		} else {
			rr.ActionDone(action, index+1 == len(c.Actions))
		}

		if err != nil {
			logError(e, c, action, outputStore, err)
			return false
		}
	}

	return true
}

// execCommand executes command
func execCommand(c *recipe.Command, outputStore *output.Store) (*exec.Cmd, io.Writer, error) {
	cmd, err := createCommand(c)

	if err != nil {
		return nil, nil, err
	}

	input, err := cmd.StdinPipe()

	if err != nil {
		return nil, nil, err
	}

	connectOutputStore(cmd, outputStore)

	err = cmd.Start()

	if err != nil {
		return nil, nil, err
	}

	go cmd.Wait()

	return cmd, input, nil
}

// createCommand creates command
func createCommand(c *recipe.Command) (*exec.Cmd, error) {
	var cmdSlice []string

	if c.User != "" {
		if !system.IsUserExist(c.User) {
			return nil, fmt.Errorf("Can't execute the command: user %s doesn't exist on the system", c.User)
		}

		cmdSlice = append(cmdSlice, "/sbin/runuser", "-s", "/bin/bash", c.User, "-c")

		if c.Recipe.Unbuffer {
			cmdSlice = append(cmdSlice, "stdbuf -o0 -e0 -i0 "+c.GetCmdline())
		} else {
			cmdSlice = append(cmdSlice, c.GetCmdline())
		}
	} else {
		if c.Recipe.Unbuffer {
			cmdSlice = append(cmdSlice, "stdbuf", "-o0", "-e0", "-i0")
		}

		cmdSlice = append(cmdSlice, c.GetCmdlineArgs()...)
	}

	return exec.Command(cmdSlice[0], cmdSlice[1:]...), nil
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
	case recipe.ACTION_OUTPUT_CONTAINS, recipe.ACTION_OUTPUT_MATCH, recipe.ACTION_OUTPUT_TRIM:
		time.Sleep(25 * time.Millisecond)
	}

	switch a.Name {
	case recipe.ACTION_EXIT:
		return action.Exit(a, cmd)
	case recipe.ACTION_EXPECT:
		return action.Expect(a, outputStore)
	case recipe.ACTION_EXPECT_STDOUT:
		return action.ExpectStdout(a, outputStore)
	case recipe.ACTION_EXPECT_STDERR:
		return action.ExpectStderr(a, outputStore)
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
	case recipe.ACTION_SIGNAL:
		return action.Signal(a, cmd)
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

	go outputIOLoop(cmd, stdoutReader, outputStore.Stdout)
	go outputIOLoop(cmd, stderrReader, outputStore.Stderr)
}

// outputIOLoop reads data from reader and writes it to output store container
func outputIOLoop(cmd *exec.Cmd, r io.Reader, c *output.Container) {
	buf := make([]byte, 16384)

	for {
		n, _ := r.Read(buf[:cap(buf)])

		if n > 0 {
			c.Write(buf[:n])
		}

		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			return
		}
	}
}

// skipCommand return true if command should be skipped
func skipCommand(c *recipe.Command, tags []string) bool {
	if c.Tag == "" {
		return false
	}

	return !sliceutil.Contains(tags, c.Tag) && !sliceutil.Contains(tags, "*")
}

// logError log error data
func logError(e *Executor, c *recipe.Command, a *recipe.Action, o *output.Store, err error) {
	if e.config.ErrsDir == "" {
		return
	}

	recipeName := strutil.Exclude(filepath.Base(c.Recipe.File), ".recipe")

	if e.logger == nil {
		err := setupLogger(e, fmt.Sprintf("%s/%s.log", e.config.ErrsDir, recipeName))

		if err != nil {
			fmtc.Printf("{r}Can't create error log: %v{!}\n", err)
			return
		}
	}

	id := passwd.GenPassword(8, passwd.STRENGTH_MEDIUM)
	origin := getErrorOrigin(c, a, id)

	e.logger.Info("(%s) %v", origin, err)

	if !o.Stdout.IsEmpty() {
		output := fmt.Sprintf("%s-stdout-%s.log", recipeName, id)
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s", e.config.ErrsDir, output), o.Stdout.Bytes(), 0644)

		if err != nil {
			e.logger.Info("(%s) Can't save stdout data: %v", origin, err)
		}
	}

	if !o.Stderr.IsEmpty() {
		output := fmt.Sprintf("%s-stderr-%s.log", recipeName, id)
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s", e.config.ErrsDir, output), o.Stderr.Bytes(), 0644)

		if err != nil {
			e.logger.Info("(%s) Can't save stderr data: %v", origin, err)
		}
	}
}

// getErrorOrigin returns info about error orign
func getErrorOrigin(c *recipe.Command, a *recipe.Action, id string) string {
	switch a {
	case nil:
		return fmt.Sprintf(
			"id: %s | command: %d | line: %d",
			id, c.Index()+1, c.Line,
		)
	default:
		return fmt.Sprintf(
			"id: %s | command: %d | action: %d:%s | line: %d",
			id, c.Index()+1, a.Index()+1, a.Name, a.Line,
		)
	}
}

// setupLogger configures logger
func setupLogger(e *Executor, file string) error {
	var err error

	e.logger, err = log.New(file, 0644)

	return err
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

// cleanupWorkingDir cleanup working dir
func cleanupWorkingDir(e *Executor, workingDir string) {
	if e.config.DisableCleanup {
		return
	}

	targets := fsutil.ListAll(workingDir, false, fsutil.ListingFilter{
		CTimeYounger: e.start.Unix(),
	})

	fsutil.ListToAbsolute(workingDir, targets)

	if len(targets) == 0 {
		return
	}

	for _, target := range targets {
		os.RemoveAll(target)
	}
}
