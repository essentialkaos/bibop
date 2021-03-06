package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"pkg.re/essentialkaos/ek.v12/fsutil"
	"pkg.re/essentialkaos/ek.v12/strutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var headersDirs = []string{
	"/usr/include",
	"/usr/local/include",
}

var libDirs = []string{
	"/lib",
	"/lib64",
	"/usr/lib",
	"/usr/lib64",
	"/usr/local/lib",
	"/usr/local/lib64",
}

// ////////////////////////////////////////////////////////////////////////////////// //

// LibLoaded is action processor for "lib-loaded"
func LibLoaded(action *recipe.Action) error {
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

// LibHeader is action processor for "lib-header"
func LibHeader(action *recipe.Action) error {
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

// LibConfig is action processor for "lib-config"
func LibConfig(action *recipe.Action) error {
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

// LibExist is action processor for "lib-exist"
func LibExist(action *recipe.Action) error {
	lib, err := action.GetS(0)

	if err != nil {
		return err
	}

	var hasLib bool

	for _, libDir := range libDirs {
		if fsutil.IsExist(libDir + "/" + lib) {
			hasLib = true
			break
		}
	}

	switch {
	case !action.Negative && !hasLib:
		return fmt.Errorf("Library file %s not found on the system", lib)
	case action.Negative && hasLib:
		return fmt.Errorf("Library file %s found on the system", lib)
	}

	return nil
}

// LibLinked is action processor for "lib-linked"
func LibLinked(action *recipe.Action) error {
	binary, err := action.GetS(0)

	if err != nil {
		return err
	}

	lib, err := action.GetS(1)

	if err != nil {
		return err
	}

	isLinked, err := isELFHasTag(binary, "Shared library", lib)

	if err != nil {
		return fmt.Errorf("Can't get info from binary: %v", err)
	}

	switch {
	case !action.Negative && !isLinked:
		return fmt.Errorf("Binary %s is not linked with shared library %s", binary, lib)
	case action.Negative && isLinked:
		return fmt.Errorf("Binary %s is linked with shared library %s", binary, lib)
	}

	return nil
}

// LibRPath is action processor for "lib-rpath"
func LibRPath(action *recipe.Action) error {
	binary, err := action.GetS(0)

	if err != nil {
		return err
	}

	rpath, err := action.GetS(1)

	if err != nil {
		return err
	}

	hasRPath, err := isELFHasTag(binary, "Library rpath", rpath)

	if err != nil {
		return fmt.Errorf("Can't get info from binary: %v", err)
	}

	switch {
	case !action.Negative && !hasRPath:
		return fmt.Errorf("Binary %s does not use %s as rpath (run-time search path)", binary, rpath)
	case action.Negative && hasRPath:
		return fmt.Errorf("Binary %s uses %s as rpath (run-time search path)", binary, rpath)
	}

	return nil
}

// LibSOName is action processor for "lib-soname"
func LibSOName(action *recipe.Action) error {
	binary, err := action.GetS(0)

	if err != nil {
		return err
	}

	soname, err := action.GetS(1)

	if err != nil {
		return err
	}

	hasName, err := isELFHasTag(binary, "Library soname", soname)

	if err != nil {
		return fmt.Errorf("Can't get info from binary: %v", err)
	}

	switch {
	case !action.Negative && !hasName:
		return fmt.Errorf("Binary %s does not contain %s in soname field", binary, soname)
	case action.Negative && hasName:
		return fmt.Errorf("Binary %s contains %s in soname field", binary, soname)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func isLibLoaded(glob string) (bool, error) {
	cmd := exec.Command("ldconfig", "-p")
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

func isELFHasTag(file, tag, glob string) (bool, error) {
	tags, err := extractELFTags(file, tag)

	if err != nil {
		return false, err
	}

	for _, tag := range tags {
		match, _ := filepath.Match(glob, tag)

		if match {
			return true, nil
		}
	}

	return false, nil
}

func extractELFTags(file, tag string) ([]string, error) {
	var result []string

	cmd := exec.Command("readelf", "-d", file)
	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "(INIT)") {
			break
		}

		if !strings.Contains(line, tag+":") {
			continue
		}

		valueIndex := strings.Index(line, "[")

		if valueIndex == -1 {
			continue
		}

		result = append(result, strings.Trim(line[valueIndex:], "[]"))
	}

	return result, nil
}
