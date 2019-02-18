package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"pkg.re/essentialkaos/ek.v10/fsutil"
	"pkg.re/essentialkaos/ek.v10/strutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var headersDirs = []string{
	"/usr/include",
	"/usr/local/include",
}

var libDirs = []string{
	"/usr/lib",
	"/usr/lib64",
}

// ////////////////////////////////////////////////////////////////////////////////// //

// actionLibLoaded is action processor for "lib-loaded"
func actionLibLoaded(action *recipe.Action) error {
	lib, err := action.GetS(0)

	if err != nil {
		return err
	}

	isLoaded, err := isLibLoaded(lib)

	if err != nil {
		return fmt.Errorf("Can't get info from ldconfig: %v", err)
	}

	switch {
	case !action.Negative && !isLoaded:
		return fmt.Errorf("Shared library %s is not loaded to the dynamic linker cache", lib)
	case action.Negative && isLoaded:
		return fmt.Errorf("Shared library %s is present in the dynamic linker cache", lib)
	}

	return nil
}

// actionLibHeader is action processor for "lib-header"
func actionLibHeader(action *recipe.Action) error {
	header, err := action.GetS(0)

	if err != nil {
		return err
	}

	var isHeaderExist bool

	for _, dir := range headersDirs {
		switch {
		case fsutil.IsExist(dir + "/" + header),
			fsutil.IsExist(dir + "/" + header + ".h"):
			isHeaderExist = true
			break
		}
	}

	switch {
	case !action.Negative && !isHeaderExist:
		return fmt.Errorf("Header %s is not found on the system", header)
	case action.Negative && isHeaderExist:
		return fmt.Errorf("Header %s found on the system", header)
	}

	return nil
}

// actionLibConfig is action processor for "lib-config"
func actionLibConfig(action *recipe.Action) error {
	lib, err := action.GetS(0)

	if err != nil {
		return err
	}

	var hasConfig bool

	for _, libDir := range libDirs {
		if fsutil.IsExist(libDir + "/pkgconfig/" + lib + ".pc") {
			hasConfig = true
			break
		}
	}

	switch {
	case !action.Negative && !hasConfig:
		return fmt.Errorf("Configuration file for %s library not found on the system", lib)
	case action.Negative && hasConfig:
		return fmt.Errorf("Configuration file for %s library found on the system", lib)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func isLibLoaded(glob string) (bool, error) {
	cmd := exec.Command("/usr/sbin/ldconfig", "-p")
	output, err := cmd.Output()

	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(string(output), "\n") {
		if !strings.Contains(line, "=>") {
			continue
		}

		line = strings.TrimSpace(line)
		line = strutil.ReadField(line, 0, false, " ")

		match, _ := filepath.Match(glob, line)

		if match {
			return true, nil
		}
	}

	return false, nil
}
