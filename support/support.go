package support

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"runtime"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/hash"
	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/depsy"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Pkg contains basic package info
type Pkg struct {
	Name    string
	Version string
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ShowSupportInfo prints verbose info about application, system, dependencies and
// important environment
func ShowSupportInfo(app, ver, gitRev string, gomod []byte) {
	pkgs := collectEnvInfo()

	fmtutil.SeparatorTitleColorTag = "{*}"
	fmtutil.SeparatorFullscreen = false

	showApplicationInfo(app, ver, gitRev)
	showOSInfo()
	showEnvInfo(pkgs)
	showDepsInfo(gomod)
	fmtutil.Separator(false)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showApplicationInfo shows verbose information about application
func showApplicationInfo(app, ver, gitRev string) {
	fmtutil.Separator(false, "APPLICATION INFO")

	fmtc.Printf("  {*}%-10s{!} %s\n", "Name:", app)
	fmtc.Printf("  {*}%-10s{!} %s\n", "Version:", ver)

	fmtc.Printf(
		"  {*}%-10s{!} %s {s}(%s/%s){!}\n", "Go:",
		strings.TrimLeft(runtime.Version(), "go"),
		runtime.GOOS, runtime.GOARCH,
	)

	if gitRev != "" {
		if !fmtc.DisableColors && fmtc.IsTrueColorSupported() {
			fmtc.Printf("  {*}%-10s{!} %s {#"+strutil.Head(gitRev, 6)+"}●{!}\n", "Git SHA:", gitRev)
		} else {
			fmtc.Printf("  {*}%-10s{!} %s\n", "Git SHA:", gitRev)
		}
	}

	bin, _ := os.Executable()
	binSHA := hash.FileHash(bin)

	if binSHA != "" {
		binSHA = strutil.Head(binSHA, 7)
		if !fmtc.DisableColors && fmtc.IsTrueColorSupported() {
			fmtc.Printf("  {*}%-10s{!} %s {#"+strutil.Head(binSHA, 6)+"}●{!}\n", "Bin SHA:", binSHA)
		} else {
			fmtc.Printf("  {*}%-10s{!} %s\n", "Bin SHA:", binSHA)
		}
	}
}

// showDepsInfo shows information about all dependencies
func showDepsInfo(gomod []byte) {
	deps := depsy.Extract(gomod, false)

	if len(deps) == 0 {
		return
	}

	fmtutil.Separator(false, "DEPENDENCIES")

	for _, dep := range deps {
		if dep.Extra == "" {
			fmtc.Printf(" {s}%8s{!}  %s\n", dep.Version, dep.Path)
		} else {
			fmtc.Printf(" {s}%8s{!}  %s {s-}(%s){!}\n", dep.Version, dep.Path, dep.Extra)
		}
	}
}

// formatValue formats value for output
func formatValue(v string) string {
	if v == "" {
		return fmtc.Sprintf("{s-}—{!}")
	}

	return v
}
