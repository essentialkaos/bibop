package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"path/filepath"
	"strings"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// checkPathSafety return true if path is save
func checkPathSafety(r *recipe.Recipe, path string) (bool, error) {
	if r.UnsafeActions {
		return true, nil
	}

	targetPath, err := filepath.Abs(path)

	if err != nil {
		return false, err
	}

	workingDir, err := filepath.Abs(r.Dir)

	if err != nil {
		return false, err
	}

	return strings.HasPrefix(targetPath, workingDir), nil
}

// fmtValue formats value
func fmtValue(v string) string {
	if v == "" {
		return `""`
	}

	return v
}
