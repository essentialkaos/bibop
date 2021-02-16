package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/fmtutil"
	"pkg.re/essentialkaos/ek.v12/fsutil"
	"pkg.re/essentialkaos/ek.v12/options"
	"pkg.re/essentialkaos/ek.v12/req"
	"pkg.re/essentialkaos/ek.v12/strutil"
	"pkg.re/essentialkaos/ek.v12/usage"
	"pkg.re/essentialkaos/ek.v12/usage/completion/bash"
	"pkg.re/essentialkaos/ek.v12/usage/completion/fish"
	"pkg.re/essentialkaos/ek.v12/usage/completion/zsh"
	"pkg.re/essentialkaos/ek.v12/usage/man"
	"pkg.re/essentialkaos/ek.v12/usage/update"

	"github.com/essentialkaos/bibop/cli/executor"
	"github.com/essentialkaos/bibop/parser"
	"github.com/essentialkaos/bibop/recipe"
	"github.com/essentialkaos/bibop/render"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Application info
const (
	APP  = "bibop"
	VER  = "4.4.1"
	DESC = "Utility for testing command-line tools"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_DRY_RUN         = "D:dry-run"
	OPT_LIST_PACKAGES   = "L:list-packages"
	OPT_FORMAT          = "f:format"
	OPT_DIR             = "d:dir"
	OPT_PATH            = "p:path"
	OPT_ERROR_DIR       = "e:error-dir"
	OPT_TAG             = "t:tag"
	OPT_QUIET           = "q:quiet"
	OPT_INGORE_PACKAGES = "ip:ignore-packages"
	OPT_NO_CLEANUP      = "nl:no-cleanup"
	OPT_NO_COLOR        = "nc:no-color"
	OPT_HELP            = "h:help"
	OPT_VER             = "v:version"

	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_DRY_RUN:         {Type: options.BOOL},
	OPT_LIST_PACKAGES:   {Type: options.BOOL},
	OPT_FORMAT:          {},
	OPT_DIR:             {},
	OPT_PATH:            {},
	OPT_ERROR_DIR:       {},
	OPT_TAG:             {Mergeble: true},
	OPT_QUIET:           {Type: options.BOOL},
	OPT_INGORE_PACKAGES: {Type: options.BOOL},
	OPT_NO_CLEANUP:      {Type: options.BOOL},
	OPT_NO_COLOR:        {Type: options.BOOL},
	OPT_HELP:            {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:             {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func Init() {
	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	if options.Has(OPT_COMPLETION) {
		os.Exit(genCompletion())
	}

	if options.Has(OPT_GENERATE_MAN) {
		os.Exit(genMan())
	}

	configureUI()
	configureSubsystems()

	if options.GetB(OPT_VER) {
		showAbout()
		os.Exit(0)
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		os.Exit(0)
	}

	validateOptions()
	process(args[0])
}

// configureUI configure user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if os.Getenv("CI") == "" {
		fmtutil.SeparatorFullscreen = true
	}

	fmtutil.SeparatorSymbol = "–"

	if fmtc.Is256ColorsSupported() {
		fmtutil.SeparatorTitleColorTag = "{#85}"
	} else {
		fmtutil.SeparatorTitleColorTag = "{c*}"
	}
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
		switch {
		case !fsutil.IsExist(errsDir):
			printErrorAndExit("Directory %s doesn't exist", errsDir)

		case !fsutil.IsDir(errsDir):
			printErrorAndExit("Object %s is not a directory", errsDir)

		case !fsutil.IsWritable(errsDir):
			printErrorAndExit("Directory %s is not writable", errsDir)
		}
	}

	wrkDir := options.GetS(OPT_DIR)

	if wrkDir != "" {
		switch {
		case !fsutil.IsExist(wrkDir):
			printErrorAndExit("Directory %s doesn't exist", wrkDir)
		case !fsutil.IsDir(wrkDir):
			printErrorAndExit("Object %s is not a directory", wrkDir)
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

	if options.GetB(OPT_LIST_PACKAGES) {
		listPackages(r.Packages)
	}

	if options.Has(OPT_ERROR_DIR) {
		errDir, _ = filepath.Abs(options.GetS(OPT_ERROR_DIR))
	}

	cfg := &executor.Config{
		Quiet:          options.GetB(OPT_QUIET),
		DisableCleanup: options.GetB(OPT_NO_CLEANUP),
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

	printError("Recipe validation errors:")

	for _, err := range errs {
		printError("  • %v", err)
	}

	os.Exit(1)
}

// listPackages list packages required bye recipe
func listPackages(pkgs []string) {
	if len(pkgs) == 0 {
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		fmtc.Println(pkg)
	}

	os.Exit(0)
}

// getValidationConfig generates validation config
func getValidationConfig(tags []string) *executor.ValidationConfig {
	vc := &executor.ValidationConfig{Tags: tags}

	if options.GetB(OPT_DRY_RUN) {
		vc.IgnoreDependencies = true
		vc.IgnorePrivileges = true
	}

	if options.GetB(OPT_INGORE_PACKAGES) {
		vc.IgnoreDependencies = true
	}

	return vc
}

// getRenderer returns renderer for executor
func getRenderer() render.Renderer {
	if options.GetB(OPT_QUIET) {
		return &render.QuietRenderer{}
	}

	if !options.Has(OPT_FORMAT) {
		return &render.TerminalRenderer{}
	}

	switch strings.ToLower(options.GetS(OPT_FORMAT)) {
	case "json":
		return &render.JSONRenderer{}
	case "xml":
		return &render.XMLRenderer{}
	case "tap", "tap13":
		return &render.TAPRenderer{}
	}

	printErrorAndExit("Unknown output format %s", options.GetS(OPT_FORMAT))

	return nil
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage prints usage info
func showUsage() {
	genUsage().Render()
}

// showAbout prints info about version
func showAbout() {
	genAbout().Render()
}

// genCompletion generates completion for different shells
func genCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "bibop"))
	case "fish":
		fmt.Printf(fish.Generate(info, "bibop"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "bibop"))
	default:
		return 1
	}

	return 0
}

// genMan generates man page
func genMan() int {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(),
		),
	)

	return 0
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "recipe")

	info.AddOption(OPT_DRY_RUN, "Parse and validate recipe")
	info.AddOption(OPT_LIST_PACKAGES, "List required packages")
	info.AddOption(OPT_FORMAT, "Output format {s-}(tap|json|xml){!}", "format")
	info.AddOption(OPT_DIR, "Path to working directory", "dir")
	info.AddOption(OPT_PATH, "Path to directory with binaries", "path")
	info.AddOption(OPT_ERROR_DIR, "Path to directory for errors data", "dir")
	info.AddOption(OPT_TAG, "Command tag", "tag")
	info.AddOption(OPT_QUIET, "Quiet mode")
	info.AddOption(OPT_INGORE_PACKAGES, "Do not check system for installed packages")
	info.AddOption(OPT_NO_CLEANUP, "Disable deleting files created during tests")
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
		"app.recipe --format json 1> ~/results/app.json",
		"Run tests from app.recipe and save result in JSON format",
	)

	return info
}

// genAbout generates info about version
func genAbout() *usage.About {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		BugTracker:    "https://github.com/essentialkaos/bibop/issues",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/bibop", update.GitHubChecker},
	}

	return about
}
