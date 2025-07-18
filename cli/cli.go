package cli

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

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fmtutil"
	"github.com/essentialkaos/ek/v13/fmtutil/panel"
	"github.com/essentialkaos/ek/v13/fmtutil/table"
	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/req"
	"github.com/essentialkaos/ek/v13/strutil"
	"github.com/essentialkaos/ek/v13/support"
	"github.com/essentialkaos/ek/v13/support/deps"
	"github.com/essentialkaos/ek/v13/support/pkgs"
	"github.com/essentialkaos/ek/v13/support/resources"
	"github.com/essentialkaos/ek/v13/terminal"
	"github.com/essentialkaos/ek/v13/terminal/tty"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
	"github.com/essentialkaos/ek/v13/usage/update"

	"github.com/essentialkaos/bibop/cli/executor"
	"github.com/essentialkaos/bibop/parser"
	"github.com/essentialkaos/bibop/recipe"
	"github.com/essentialkaos/bibop/render"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Application info
const (
	APP  = "bibop"
	VER  = "8.2.0"
	DESC = "Utility for testing command-line tools"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_DRY_RUN            = "D:dry-run"
	OPT_LIST_PACKAGES      = "L:list-packages"
	OPT_LIST_PACKAGES_FLAT = "L1:list-packages-flat"
	OPT_VARIABLES          = "V:variables"
	OPT_BARCODE            = "B:barcode"
	OPT_EXTRA              = "X:extra"
	OPT_TIME               = "T:time"
	OPT_PAUSE              = "P:pause"
	OPT_FORMAT             = "f:format"
	OPT_DIR                = "d:dir"
	OPT_PATH               = "p:path"
	OPT_ERROR_DIR          = "e:error-dir"
	OPT_TAG                = "t:tag"
	OPT_QUIET              = "q:quiet"
	OPT_IGNORE_PACKAGES    = "ip:ignore-packages"
	OPT_NO_CLEANUP         = "nl:no-cleanup"
	OPT_NO_COLOR           = "nc:no-color"
	OPT_HELP               = "h:help"
	OPT_VER                = "v:version"

	OPT_UPDATE       = "U:update"
	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_DRY_RUN:            {Type: options.BOOL},
	OPT_LIST_PACKAGES:      {Type: options.BOOL},
	OPT_LIST_PACKAGES_FLAT: {Type: options.BOOL},
	OPT_VARIABLES:          {Type: options.BOOL},
	OPT_BARCODE:            {Type: options.BOOL},
	OPT_EXTRA:              {Type: options.INT, Value: 10, Min: 1, Max: 256},
	OPT_TIME:               {Type: options.BOOL},
	OPT_PAUSE:              {Type: options.FLOAT, Max: 60},
	OPT_FORMAT:             {},
	OPT_DIR:                {},
	OPT_PATH:               {},
	OPT_ERROR_DIR:          {},
	OPT_TAG:                {Mergeble: true},
	OPT_QUIET:              {Type: options.BOOL},
	OPT_IGNORE_PACKAGES:    {Type: options.BOOL},
	OPT_NO_CLEANUP:         {Type: options.BOOL},
	OPT_NO_COLOR:           {Type: options.BOOL},
	OPT_HELP:               {Type: options.BOOL},
	OPT_VER:                {Type: options.MIXED},

	OPT_UPDATE:       {Type: options.MIXED},
	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

var colorTagApp, colorTagVer string
var rawOutput bool

// ////////////////////////////////////////////////////////////////////////////////// //

func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.Error(" - "))
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			WithPackages(pkgs.Collect(
				"ca-certificates", "systemd", "systemd-sysv",
				"initscripts", "libc-bin", "dpkg",
				"gcc", "python2", "python3", "binutils",
			)).
			WithResources(resources.Collect()).
			Print()
		os.Exit(0)
	case withSelfUpdate && options.GetB(OPT_UPDATE):
		os.Exit(updateBinary())
	case options.GetB(OPT_HELP) || len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	configureSubsystems()

	validateOptions()
	process(args.Get(0).Clean().String())
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
		rawOutput = true
	}

	if os.Getenv("CI") == "" {
		fmtutil.SeparatorFullscreen = true
	} else {
		fmtc.DisableColors = false
	}

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{*}{#66CC99}", "{#66CC99}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{#85}", "{#85}"
	default:
		colorTagApp, colorTagVer = "{*}{c}", "{c}"
	}

	fmtutil.SeparatorTitleColorTag = colorTagApp
}

// configureUI configure user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	fmtutil.SeparatorSymbol = "–"
	panel.Indent = 5
}

// configureSubsystems configures bibop subsystems
func configureSubsystems() {
	req.Global.SetUserAgent(APP, VER)

	if !options.Has(OPT_PATH) {
		return
	}

	pathOpt, err := filepath.Abs(options.GetS(OPT_PATH))

	if err != nil {
		printErrorAndExit(err.Error())
	}

	newPath := os.Getenv("PATH") + ":" + pathOpt

	err = os.Setenv("PATH", newPath)

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// validateOptions validates options
func validateOptions() {
	errsDir := options.GetS(OPT_ERROR_DIR)

	if errsDir != "" {
		err := fsutil.ValidatePerms("DW", errsDir)

		if err != nil {
			printErrorAndExit(err.Error())
		}
	}

	wrkDir := options.GetS(OPT_DIR)

	if wrkDir != "" {
		err := fsutil.ValidatePerms("DR", wrkDir)

		if err != nil {
			printErrorAndExit(err.Error())
		}
	}
}

// process start recipe processing
func process(file string) {
	var errDir string

	r, err := parser.Parse(file)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	if options.Has(OPT_DIR) {
		r.Dir, _ = filepath.Abs(options.GetS(OPT_DIR))
	} else {
		r.Dir, _ = filepath.Abs(filepath.Dir(file))
	}

	switch {
	case options.GetB(OPT_LIST_PACKAGES),
		options.GetB(OPT_LIST_PACKAGES_FLAT):
		listPackages(r)
		os.Exit(0)
	case options.GetB(OPT_VARIABLES):
		listVariables(r)
		os.Exit(0)
	case options.GetB(OPT_BARCODE):
		printBarcode(r)
		os.Exit(0)
	}

	if options.Has(OPT_ERROR_DIR) {
		errDir, _ = filepath.Abs(options.GetS(OPT_ERROR_DIR))
	}

	cfg := &executor.Config{
		Quiet:          options.GetB(OPT_QUIET),
		DisableCleanup: options.GetB(OPT_NO_CLEANUP),
		DebugLines:     options.GetI(OPT_EXTRA),
		Pause:          options.GetF(OPT_PAUSE),
		ErrsDir:        errDir,
	}

	e := executor.NewExecutor(cfg)
	tags := strutil.Fields(options.GetS(OPT_TAG))

	validate(e, r, tags)

	rr := getRenderer()

	if !e.Run(rr, r, tags) {
		os.Exit(1)
	}
}

// validate validates recipe and print validation errors
func validate(e *executor.Executor, r *recipe.Recipe, tags []string) {
	errs := e.Validate(r, getValidationConfig(tags))

	if len(errs) == 0 {
		if options.GetB(OPT_DRY_RUN) {
			fmtc.Println("{g}This recipe has no issues{!}")
			os.Exit(0)
		}

		return
	}

	terminal.Error("Recipe validation errors:")

	for _, err := range errs {
		terminal.Error("  • %v", err)
	}

	os.Exit(1)
}

// listPackages shows list packages required by recipe
func listPackages(r *recipe.Recipe) {
	if len(r.Packages) == 0 {
		return
	}

	if options.GetB(OPT_LIST_PACKAGES_FLAT) {
		fmt.Println(strings.Join(r.Packages, " "))
	} else {
		fmtc.If(!rawOutput).NewLine()
		for _, pkg := range r.Packages {
			fmtc.If(!rawOutput).Printf("{s-}•{!} %s\n", pkg)
			fmtc.If(rawOutput).Printf("%s\n", pkg)
		}
		fmtc.If(!rawOutput).NewLine()
	}
}

// listVariables shows list of variables
func listVariables(r *recipe.Recipe) {
	t := table.NewTable("Name", "Value")

	for _, v := range r.GetVariables() {
		t.Add(v, strutil.Q(r.GetVariable(v, true), "{s-}—{!}"))
	}

	t.Separator()

	for _, v := range recipe.DynamicVariables {
		t.Add("{s}"+v+"{!}", strutil.Q(r.GetVariable(v, false), "{s-}—{!}"))
	}

	fmtc.NewLine()
	t.Render()
	fmtc.NewLine()
}

// getValidationConfig generates validation config
func getValidationConfig(tags []string) *executor.ValidationConfig {
	vc := &executor.ValidationConfig{Tags: tags}

	if options.GetB(OPT_DRY_RUN) {
		vc.IgnoreDependencies = true
		vc.IgnorePackages = true
		vc.IgnorePrivileges = true
	}

	if options.GetB(OPT_IGNORE_PACKAGES) {
		vc.IgnorePackages = true
	}

	return vc
}

// getRenderer returns renderer for executor
func getRenderer() render.Renderer {
	if options.GetB(OPT_QUIET) {
		return &render.QuietRenderer{}
	}

	if !options.Has(OPT_FORMAT) {
		return &render.TerminalRenderer{PrintExecTime: options.GetB(OPT_TIME)}
	}

	switch strings.ToLower(options.GetS(OPT_FORMAT)) {
	case "json":
		return &render.JSONRenderer{}
	case "xml":
		return &render.XMLRenderer{Version: VER}
	case "tap13":
		return &render.TAP13Renderer{Version: VER}
	case "tap14":
		return &render.TAP14Renderer{Version: VER}
	}

	printErrorAndExit("Unknown output format %s", options.GetS(OPT_FORMAT))

	return nil
}

// printErrorAndExit print error message and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	terminal.Error(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, APP, "recipe"))
	case "fish":
		fmt.Print(fish.Generate(info, APP))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, APP, "*.recipe"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "recipe")

	info.AppNameColorTag = colorTagApp

	info.AddOption(OPT_DRY_RUN, "Parse and validate recipe")
	info.AddOption(OPT_EXTRA, "Number of output lines for failed action {s-}(default: 10){!}", "lines")
	info.AddOption(OPT_PAUSE, "Pause between commands in seconds", "duration")
	info.AddOption(OPT_LIST_PACKAGES, "List required packages")
	info.AddOption(OPT_LIST_PACKAGES_FLAT, "List required packages in one line {s-}(useful for scripts){!}")
	info.AddOption(OPT_VARIABLES, "List recipe variables")
	info.AddOption(OPT_BARCODE, "Show unique barcode for test {s-}(based on recipe and required packages){!}")
	info.AddOption(OPT_TIME, "Print execution time for every action")
	info.AddOption(OPT_FORMAT, "Output format {s-}(tap13|tap14|json|xml){!}", "format")
	info.AddOption(OPT_DIR, "Path to working directory", "dir")
	info.AddOption(OPT_PATH, "Path to directory with binaries", "path")
	info.AddOption(OPT_ERROR_DIR, "Path to directory for errors data", "dir")
	info.AddOption(OPT_TAG, "One or more command tags to run", "tag")
	info.AddOption(OPT_QUIET, "Quiet mode")
	info.AddOption(OPT_IGNORE_PACKAGES, "Do not check system for installed packages")
	info.AddOption(OPT_NO_CLEANUP, "Disable deleting files created during tests")

	if withSelfUpdate {
		info.AddOption(OPT_UPDATE, "Update application to the latest version")
	}

	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"app.recipe",
		"Run tests from app.recipe",
	)

	info.AddExample(
		"app.recipe --quiet --error-dir bibop-errors",
		"Run tests from app.recipe in quiet mode and save errors data to bibop-errors directory",
	)

	info.AddExample(
		"app.recipe --tag init,service",
		"Run tests from app.recipe and execute commands with tags init and service",
	)

	info.AddExample(
		"app.recipe --extra",
		"Run tests from app.recipe and print the last 10 lines from command output if action was failed",
	)

	info.AddExample(
		"app.recipe --extra=50",
		"Run tests from app.recipe and print the last 50 lines from command output if action was failed",
	)

	info.AddExample(
		"app.recipe --format json 1> ~/results/app.json",
		"Run tests from app.recipe and save result in JSON format",
	)

	info.AddRawExample(
		"sudo dnf install $(bibop app.recipe -L1)",
		"Install all packages required for tests",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",

		License:    "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		BugTracker: "https://github.com/essentialkaos/bibop/issues",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
		about.UpdateChecker = usage.UpdateChecker{"essentialkaos/bibop", update.GitHubChecker}
	}

	return about
}
