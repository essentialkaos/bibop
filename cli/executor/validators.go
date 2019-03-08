package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"regexp"

	"pkg.re/essentialkaos/ek.v10/fsutil"
	"pkg.re/essentialkaos/ek.v10/sliceutil"
	"pkg.re/essentialkaos/ek.v10/system"

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

// checkRecipePriveleges checks user priveleges
func checkRecipePriveleges(r *recipe.Recipe) error {
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
	var knownVars []string

	varRegex := regexp.MustCompile(`\{([a-zA-Z0-9_-]+)\}`)

	for _, c := range r.Commands {
		submatch := varRegex.FindAllStringSubmatch(c.GetCmdline(), -1)

		if len(submatch) != 0 {
			errs = append(errs, convertSubmatchToErrors(nil, submatch)...)
		}

		submatch = varRegex.FindAllStringSubmatch(c.User, -1)

		if len(submatch) != 0 {
			errs = append(errs, convertSubmatchToErrors(nil, submatch)...)
		}

		for _, a := range c.Actions {
			knownVars = append(knownVars, getDynamicVars(a)...)

			for argIndex := range a.Arguments {
				arg, _ := a.GetS(argIndex)
				submatch = varRegex.FindAllStringSubmatch(arg, -1)

				if len(submatch) != 0 {
					errs = append(errs, convertSubmatchToErrors(knownVars, submatch)...)
				}
			}
		}
	}

	return errs
}

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
func convertSubmatchToErrors(knownVars []string, data [][]string) []error {
	var errs []error

	for _, match := range data {
		if sliceutil.Contains(knownVars, match[1]) {
			continue
		}

		errs = append(errs, fmt.Errorf("Can't find veriable with name %s", match[1]))
	}

	return errs
}
