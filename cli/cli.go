package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"path/filepath"

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
	"pkg.re/essentialkaos/ek.v12/usage/update"

	"github.com/essentialkaos/bibop/cli/executor"
	"github.com/essentialkaos/bibop/parser"
	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Application info
const (
	APP  = "bibop"
	VER  = "2.2.0"
	DESC = "Utility for testing command-line tools"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_DIR           = "d:dir"
	OPT_ERROR_DIR     = "e:error-dir"
	OPT_TAG           = "t:tag"
	OPT_QUIET         = "q:quiet"
	OPT_DRY_RUN       = "D:dry-run"
	OPT_LIST_PACKAGES = "L:list-packages"
	OPT_NO_CLEANUP    = "NC:no-cleanup"
	OPT_NO_COLOR      = "nc:no-color"
	OPT_HELP          = "h:help"
	OPT_VER           = "v:version"

	OPT_COMPLETION = "completion"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_DIR:           {},
	OPT_ERROR_DIR:     {},
	OPT_TAG:           {Mergeble: true},
	OPT_QUIET:         {Type: options.BOOL},
	OPT_DRY_RUN:       {Type: options.BOOL},
	OPT_LIST_PACKAGES: {Type: options.BOOL},
	OPT_NO_CLEANUP:    {Type: options.BOOL},
	OPT_NO_COLOR:      {Type: options.BOOL},
	OPT_HELP:          {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:           {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION: {},
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
		genCompletion()
	}

	configureUI()
	configureSubsystems()

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	validateOptions()
	process(args[0])
}

// configureUI configure user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	fmtutil.SeparatorFullscreen = true
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

	if !e.Run(r, tags) {
		os.Exit(1)
	}
}

// validate validates recipe and print validation errors
func validate(e *executor.Executor, r *recipe.Recipe, tags []string) {
	vc := &executor.ValidationConfig{Tags: tags}

	if options.GetB(OPT_DRY_RUN) {
		vc.IgnoreDependencies = true
		vc.IgnorePrivileges = true
	}

	errs := e.Validate(r, vc)

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

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "recipe")

	info.AddOption(OPT_DIR, "Path to working directory", "dir")
	info.AddOption(OPT_ERROR_DIR, "Path to directory for errors data", "dir")
	info.AddOption(OPT_TAG, "Command tag", "tag")
	info.AddOption(OPT_QUIET, "Quiet mode")
	info.AddOption(OPT_DRY_RUN, "Parse and validate recipe")
	info.AddOption(OPT_LIST_PACKAGES, "List required packages")
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

	return info
}

// genCompletion generates completion for different shells
func genCompletion() {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "bibop"))
	case "fish":
		fmt.Printf(fish.Generate(info, "bibop"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "bibop"))
	default:
		os.Exit(1)
	}

	os.Exit(0)
}

// showAbout prints info about version
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/bibop", update.GitHubChecker},
	}

	about.Render()
}
