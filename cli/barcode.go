package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os/exec"

	"github.com/essentialkaos/ek/v12/env"
	"github.com/essentialkaos/ek/v12/fmtutil/barcode"
	"github.com/essentialkaos/ek/v12/hash"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// printBarcode prints barcode based on list of installed packages and recipe content
func printBarcode(r *recipe.Recipe) {
	pkgsInfo, err := getPackagesInfo(r.Packages)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	recipeCrc := hash.FileHash(r.File)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	fmt.Println(barcode.Dots(
		append(pkgsInfo, []byte(recipeCrc)...),
	))
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getPackagesInfo returns unique info about installed packages
func getPackagesInfo(pkgs []string) ([]byte, error) {
	switch {
	case env.Which("rpm") != "":
		return getRPMPackagesInfo(pkgs)
	case env.Which("dpkg") != "":
		return getDEBPackagesInfo(pkgs)
	}

	return []byte{}, fmt.Errorf("Can't generate barcode for current system")
}

// getRPMPackagesInfo returns info about installed packages from rpm
func getRPMPackagesInfo(pkgs []string) ([]byte, error) {
	cmd := exec.Command("rpm", "-q", "--qf", "%{pkgid}\n")
	cmd.Env = []string{"LC_ALL=C"}
	cmd.Args = append(cmd.Args, pkgs...)

	output, _ := cmd.Output()

	return output, nil
}

// getDEBPackagesChecksums returns info about installed packages from apt
func getDEBPackagesInfo(pkgs []string) ([]byte, error) {
	cmd := exec.Command("apt-cache", "show")
	cmd.Env = []string{"LC_ALL=C"}
	cmd.Args = append(cmd.Args, pkgs...)

	output, _ := cmd.Output()

	return output, nil
}
