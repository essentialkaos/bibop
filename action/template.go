package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strconv"
	"text/template"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// varWrapper is a recipe wrapper for accessing variables
type varWrapper struct {
	r *recipe.Recipe
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Template is action processor for "template"
func Template(action *recipe.Action) error {
	mode := uint64(0644)
	source, err := action.GetS(0)

	if err != nil {
		return err
	}

	dest, err := action.GetS(1)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, source)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Action uses unsafe path (%s)", source)
	}

	isSafePath, err = checkPathSafety(action.Command.Recipe, dest)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Action uses unsafe path (%s)", dest)
	}

	if action.Has(2) {
		modeStr, _ := action.GetS(1)
		mode, err = strconv.ParseUint(modeStr, 8, 32)

		if err != nil {
			return err
		}
	}

	tmplData, err := os.ReadFile(source)

	if err != nil {
		fmt.Errorf("Can't read template %q: %v", source, err)
	}

	tmpl, err := template.New("").Parse(string(tmplData))

	if err != nil {
		return fmt.Errorf("Can't parse template %q: %v", source, err)
	}

	fd, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(mode))

	if err != nil {
		return fmt.Errorf("Can't save template data into %q: %v", dest, err)
	}

	defer fd.Close()

	vw := &varWrapper{action.Command.Recipe}
	err = tmpl.Execute(fd, vw)

	if err != nil {
		return fmt.Errorf("Can't render template %q: %v", source, err)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Var returns variable value
func (vw *varWrapper) Var(name string) string {
	return vw.r.GetVariable(name, true)
}

// Is compares variable value
func (vw *varWrapper) Is(name, value string) bool {
	return vw.r.GetVariable(name, true) == value
}

// ////////////////////////////////////////////////////////////////////////////////// //
