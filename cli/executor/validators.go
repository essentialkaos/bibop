package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/essentialkaos/ek/v12/env"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/sliceutil"
	"github.com/essentialkaos/ek/v12/strutil"
	"github.com/essentialkaos/ek/v12/system"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// checkRecipeWorkingDir checks recipe working dir
func checkRecipeWorkingDir(r *recipe.Recipe) error {
	switch {
	case !fsutil.IsExist(r.Dir):
		return fmt.Errorf("Directory %s doesn't exist", r.Dir)
	case !fsutil.IsDir(r.Dir):
		return fmt.Errorf("%s is not a directory", r.Dir)
	case !fsutil.IsReadable(r.Dir):
		return fmt.Errorf("Directory %s is not readable", r.Dir)
	}

	return nil
}

// checkRecipePrivileges checks user privileges
func checkRecipePrivileges(r *recipe.Recipe) error {
	if !r.RequireRoot {
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

// checkRecipeTags checks tags
func checkRecipeTags(r *recipe.Recipe, tags []string) []error {
	if len(tags) == 0 || sliceutil.Contains(tags, "*") {
		return nil
	}

	var knownTags []string

	for _, c := range r.Commands {
		if c.Tag != "" {
			knownTags = append(knownTags, c.Tag)
		}
	}

	var errs []error

	for _, tag := range tags {
		if !sliceutil.Contains(knownTags, tag) {
			errs = append(errs, fmt.Errorf("This recipe doesn't contain command with tag '%s'", tag))
		}
	}

	return errs
}

// checkRecipeVariables checks recipe for unknown variables
func checkRecipeVariables(r *recipe.Recipe) []error {
	var errs []error

	knownVars := append(recipe.DynamicVariables[:0:0], recipe.DynamicVariables...)
	varRegex := regexp.MustCompile(`\{([a-zA-Z0-9_-]+)\}`)

	for _, c := range r.Commands {
		submatch := varRegex.FindAllStringSubmatch(c.GetCmdline(), -1)

		if len(submatch) != 0 {
			errs = append(errs, convertSubmatchToErrors(nil, submatch, c.Line)...)
		}

		submatch = varRegex.FindAllStringSubmatch(c.User, -1)

		if len(submatch) != 0 {
			errs = append(errs, convertSubmatchToErrors(nil, submatch, c.Line)...)
		}

		for _, a := range c.Actions {
			knownVars = append(knownVars, getDynamicVars(a)...)

			for argIndex := range a.Arguments {
				arg, _ := a.GetS(argIndex)

				submatch = varRegex.FindAllStringSubmatch(arg, -1)

				if len(submatch) != 0 {
					errs = append(errs, convertSubmatchToErrors(knownVars, submatch, a.Line)...)
				}
			}
		}
	}

	return errs
}

// checkPackages checks if required packages are installed on the system
func checkPackages(r *recipe.Recipe) []error {
	if len(r.Packages) == 0 {
		return nil
	}

	switch {
	case env.Which("rpm") != "":
		return checkRPMPackages(r.Packages)
	case env.Which("dpkg") != "":
		return checkDEBPackages(r.Packages)
	}

	return []error{errors.New("Can't check required packages availability: Unsupported OS")}
}

// checkDependencies checks if all required binaries are present on the system
func checkDependencies(r *recipe.Recipe) []error {
	var errs []error

	binCache := make(map[string]bool)

	for _, c := range r.Commands {
		for _, a := range c.Actions {
			var binary string

			switch a.Name {
			case recipe.ACTION_SERVICE_PRESENT:
				binary = "systemctl"
			case recipe.ACTION_SERVICE_ENABLED:
				binary = "systemctl"
			case recipe.ACTION_SERVICE_WORKS:
				binary = "systemctl"
			case recipe.ACTION_WAIT_SERVICE:
				binary = "systemctl"
			case recipe.ACTION_LIB_LOADED:
				binary = "ldconfig"
			case recipe.ACTION_LIB_CONFIG:
				binary = "pkg-config"
			case recipe.ACTION_LIB_LINKED:
				binary = "readelf"
			case recipe.ACTION_LIB_RPATH:
				binary = "readelf"
			case recipe.ACTION_LIB_SONAME:
				binary = "readelf"
			case recipe.ACTION_LIB_EXPORTED:
				binary = "nm"
			case recipe.ACTION_PYTHON2_PACKAGE:
				binary = "python"
			case recipe.ACTION_PYTHON3_PACKAGE:
				binary = "python3"
			}

			if !hasBinary(binCache, binary) {
				errs = append(errs, fmt.Errorf(
					"Line %d: Action %q requires %q binary", a.Line, a.Name, binary,
				))
			}
		}
	}

	return errs
}

// hasBinary checks if binary is present on the system
func hasBinary(binCache map[string]bool, binary string) bool {
	isExist, ok := binCache[binary]

	if ok {
		return isExist
	}

	binCache[binary] = env.Which(binary) != ""

	return binCache[binary]
}

// getDynamicVars returns slice with dynamic vars
func getDynamicVars(a *recipe.Action) []string {
	switch a.Name {
	case recipe.ACTION_CHECKSUM_READ:
		v, _ := a.GetS(1)
		return []string{v}
	default:
		return nil
	}
}

// convertSubmatchToErrors convert slice with submatch data to error slice
func convertSubmatchToErrors(knownVars []string, data [][]string, line uint16) []error {
	var errs []error

	for _, match := range data {
		if sliceutil.Contains(knownVars, match[1]) {
			continue
		}

		errs = append(errs, fmt.Errorf("Line %d: can't find variable with name %s", line, match[1]))
	}

	return errs
}

// checkRPMPackages checks if rpm packages are installed
func checkRPMPackages(pkgs []string) []error {
	cmd := exec.Command("rpm", "-q", "--queryformat", "%{name}\n")
	cmd.Env = []string{"LC_ALL=C"}
	cmd.Args = append(cmd.Args, pkgs...)

	output, _ := cmd.Output()

	var errs []error

	for _, pkgInfo := range strings.Split(string(output), "\n") {
		if strings.Contains(pkgInfo, " ") {
			pkgName := strutil.ReadField(pkgInfo, 1, true)
			errs = append(errs, fmt.Errorf("Package %s is not installed", pkgName))
		}
	}

	return errs
}

// checkDEBPackages checks if deb packages are installed
func checkDEBPackages(pkgs []string) []error {
	cmd := exec.Command("dpkg-query", "-l")
	cmd.Env = []string{"LC_ALL=C"}
	cmd.Args = append(cmd.Args, pkgs...)

	output, _ := cmd.Output()

	var errs []error

	for _, pkgInfo := range strings.Split(string(output), "\n") {
		if strings.Contains(pkgInfo, "no packages found") {
			pkgName := strutil.Exclude(pkgInfo, "dpkg-query: no packages found matching ")
			errs = append(errs, fmt.Errorf("Package %s is not installed", pkgName))
		}
	}

	return errs
}
