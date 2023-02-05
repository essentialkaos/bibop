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

// checkRecipeVariables checks recipe for unnown variables
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

// checkPackages checks packages
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
		if strings.Contains(pkgInfo, "is not installed") {
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
