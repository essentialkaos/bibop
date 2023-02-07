package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/sliceutil"
	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const PROP_LIB_EXPORTED = "LIB_EXPORTED"

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

	if hasConfig && !action.Negative {
		if !isLibPkgConfigValid(lib) {
			return fmt.Errorf(
				"Configuration file for %s library is not valid (try 'pkg-config --exists --print-errors %s' for check)",
				lib, lib,
			)
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

	hasLib := getLibPath(lib) != ""

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

// LibExported is action processor for "lib-exported"
func LibExported(action *recipe.Action) error {
	var libFile string

	command := action.Command

	lib, err := action.GetS(0)

	if err != nil {
		return err
	}

	symbol, err := action.GetS(1)

	if err != nil {
		return err
	}

	if fsutil.IsExist(lib) {
		libFile = lib
	} else {
		libFile = getLibPath(lib)
	}

	if libFile == "" {
		return fmt.Errorf("Library file %s not found on the system", lib)
	}

	var symbols []string

	if command.Data.Has(PROP_LIB_EXPORTED) {
		symbols = command.Data.Get(PROP_LIB_EXPORTED).([]string)
	} else {
		symbols, err = extractSOExports(libFile)

		if err != nil {
			return err
		}

		command.Data.Set(PROP_LIB_EXPORTED, symbols)
	}

	hasSymbol := sliceutil.Contains(symbols, symbol)

	switch {
	case !action.Negative && !hasSymbol:
		return fmt.Errorf("Library %s doesn't export symbol %q", lib, symbol)
	case action.Negative && hasSymbol:
		return fmt.Errorf("Library %s exports symbol %q", lib, symbol)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getLibPath returns path to library file
func getLibPath(lib string) string {
	for _, libDir := range libDirs {
		if fsutil.IsExist(libDir + "/" + lib) {
			return libDir + "/" + lib
		}
	}

	return ""
}

// isLibLoaded returns true if library is loaded by linker
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

// isLibPkgConfigValid checks if library package config is loaded and valid
func isLibPkgConfigValid(lib string) bool {
	return exec.Command("pkg-config", "--exists", lib).Run() == nil
}

// isELFHasTag returns true if elf file contains given tag
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

// extractELFTags extracts tags from ELF file
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

// extractSOExports returns slice with exported symbols
func extractSOExports(file string) ([]string, error) {
	var result []string

	cmd := exec.Command("nm", "--dynamic", "--defined-only", file)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf(string(output))
	}

	for _, line := range strings.Split(string(output), "\n") {
		switch strutil.ReadField(line, 1, false, " ") {
		case "T", "R", "D":
			result = append(result, strutil.ReadField(line, 2, false, " "))
		}
	}

	return result, nil
}
