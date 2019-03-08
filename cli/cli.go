package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"

	"pkg.re/essentialkaos/ek.v10/fmtc"
	"pkg.re/essentialkaos/ek.v10/fmtutil"
	"pkg.re/essentialkaos/ek.v10/options"
	"pkg.re/essentialkaos/ek.v10/strutil"
	"pkg.re/essentialkaos/ek.v10/usage"
	"pkg.re/essentialkaos/ek.v10/usage/update"

	"github.com/essentialkaos/bibop/cli/executor"
	"github.com/essentialkaos/bibop/parser"
	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Application info
const (
	APP     = "bibop"
	VER     = "0.0.1"
	RELEASE = "β6"
	DESC    = "Utility for testing command-line tools"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_DIR      = "d:dir"
	OPT_LOG      = "l:log"
	OPT_TAG      = "t:tag"
	OPT_QUIET    = "q:quiet"
	OPT_DRY_RUN  = "D:dry-run"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_DIR:      {},
	OPT_LOG:      {},
	OPT_TAG:      {Mergeble: true},
	OPT_QUIET:    {Type: options.BOOL},
	OPT_DRY_RUN:  {Type: options.BOOL},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},
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

	configureUI()

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		return
	}

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

// process start recipe processing
func process(file string) {
	r, err := parser.Parse(file)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	if options.Has(OPT_DIR) {
		r.Dir = options.GetS(OPT_DIR)
	}

	e := executor.NewExecutor(options.GetB(OPT_QUIET))

	if options.Has(OPT_LOG) {
		err = e.SetupLogger(options.GetS(OPT_LOG))

		if err != nil {
			printErrorAndExit(err.Error())
		}
	}

	tags := strutil.Fields(options.GetS(OPT_TAG))

	validate(e, r, tags)

	if !e.Run(r, tags) {
		os.Exit(1)
	}
}

// validate validates recipe and print validation errors
func validate(e *executor.Executor, r *recipe.Recipe, tags []string) {
	errs := e.Validate(r, tags)

	if len(errs) == 0 {
		if options.GetB(OPT_DRY_RUN) {
			fmtc.Println("{g}This recipe has no issues{!}")
			os.Exit(0)
		}

		return
	}

	for _, err := range errs {
		printError(err.Error())
	}

	os.Exit(1)
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

func showUsage() {
	info := usage.NewInfo("", "recipe")

	info.AddOption(OPT_DIR, "Path to working directory")
	info.AddOption(OPT_LOG, "Path to log file for verbose info about errors")
	info.AddOption(OPT_TAG, "Command tag", "tag")
	info.AddOption(OPT_QUIET, "Quiet mode")
	info.AddOption(OPT_DRY_RUN, "Parse and validate recipe")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"application.recipe",
		"Run tests from application.recipe",
	)

	info.AddExample(
		"application.recipe --quiet --log errors.log ",
		"Run tests from application.recipe in quiet mode and log errors to errors.log",
	)

	info.AddExample(
		"application.recipe --quiet --log errors.log --tag init,service",
		"Run tests from application.recipe and execute commands with tags init and service",
	)

	info.Render()
}

func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Release:       " " + RELEASE,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/bibop", update.GitHubChecker},
	}

	about.Render()
}
