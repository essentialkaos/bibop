package support

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os/exec"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/system"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// showOSInfo shows verbose information about system
func showOSInfo() {
	osInfo, err := system.GetOSInfo()

	if err == nil {
		fmtutil.Separator(false, "OS INFO")
		fmtc.Printf("  {*}%-15s{!} %s\n", "Name:", formatValue(osInfo.Name))
		fmtc.Printf("  {*}%-15s{!} %s\n", "Pretty Name:", formatValue(osInfo.PrettyName))
		fmtc.Printf("  {*}%-15s{!} %s\n", "Version:", formatValue(osInfo.VersionID))
		fmtc.Printf("  {*}%-15s{!} %s\n", "ID:", formatValue(osInfo.ID))
		fmtc.Printf("  {*}%-15s{!} %s\n", "ID Like:", formatValue(osInfo.IDLike))
		fmtc.Printf("  {*}%-15s{!} %s\n", "Version ID:", formatValue(osInfo.VersionID))
		fmtc.Printf("  {*}%-15s{!} %s\n", "Version Code:", formatValue(osInfo.VersionCodename))
		fmtc.Printf("  {*}%-15s{!} %s\n", "CPE:", formatValue(osInfo.CPEName))
	}

	systemInfo, err := system.GetSystemInfo()

	if err != nil {
		return
	} else {
		if osInfo == nil {
			fmtutil.Separator(false, "SYSTEM INFO")
			fmtc.Printf("  {*}%-15s{!} %s\n", "Name:", formatValue(systemInfo.OS))
			fmtc.Printf("  {*}%-15s{!} %s\n", "Version:", formatValue(systemInfo.Version))
		}
	}

	fmtc.Printf("  {*}%-15s{!} %s\n", "Arch:", formatValue(systemInfo.Arch))
	fmtc.Printf("  {*}%-15s{!} %s\n", "Kernel:", formatValue(systemInfo.Kernel))

	containerEngine := "No"

	switch {
	case fsutil.IsExist("/.dockerenv"):
		containerEngine = "Yes (Docker)"
	case fsutil.IsExist("/run/.containerenv"):
		containerEngine = "Yes (Podman)"
	}

	fmtc.NewLine()
	fmtc.Printf("  {*}%-15s{!} %s\n", "Container:", containerEngine)
}

// showEnvInfo shows info about environment
func showEnvInfo(pkgs []Pkg) {
	fmtutil.Separator(false, "ENVIRONMENT")

	for _, pkg := range pkgs {
		fmtc.Printf("  {*}%-18s{!} %s\n", pkg.Name+":", formatValue(pkg.Version))
	}
}

// collectEnvInfo collects info about packages
func collectEnvInfo() []Pkg {
	if isDEBBased() {
		return []Pkg{
			getPackageInfo("ca-certificates"),
			getPackageInfo("systemd"),
			getPackageInfo("systemd-sysv"),
			getPackageInfo("initscripts"),
			getPackageInfo("libc-bin"),
			getPackageInfo("dpkg"),
			getPackageInfo("python"),
			getPackageInfo("python3"),
			getPackageInfo("binutils"),
		}
	}

	return []Pkg{
		getPackageInfo("ca-certificates"),
		getPackageInfo("systemd"),
		getPackageInfo("systemd-sysv"),
		getPackageInfo("initscripts"),
		getPackageInfo("glibc"),
		getPackageInfo("rpm"),
		getPackageInfo("python"),
		getPackageInfo("python3"),
		getPackageInfo("binutils"),
	}
}

// getPackageVersion returns package name from rpm database
func getPackageInfo(name string) Pkg {
	switch {
	case isDEBBased():
		return getDEBPackageInfo(name)
	case isRPMBased():
		return getRPMPackageInfo(name)
	}

	return Pkg{name, ""}
}

// isDEBBased returns true if is DEB-based distro
func isRPMBased() bool {
	return fsutil.IsExist("/usr/bin/rpm")
}

// isDEBBased returns true if is DEB-based distro
func isDEBBased() bool {
	return fsutil.IsExist("/usr/bin/dpkg-query")
}

// getRPMPackageInfo returns info about RPM package
func getRPMPackageInfo(name string) Pkg {
	cmd := exec.Command("rpm", "-q", name)
	out, err := cmd.Output()

	if err != nil || len(out) == 0 {
		return Pkg{name, ""}
	}

	return Pkg{name, strings.TrimRight(string(out), "\n\r")}
}

// getDEBPackageInfo returns info about DEB package
func getDEBPackageInfo(name string) Pkg {
	cmd := exec.Command("dpkg-query", "--showformat=${Version}", "--show", name)
	out, err := cmd.Output()

	if err != nil || len(out) == 0 {
		return Pkg{name, ""}
	}

	return Pkg{name, string(out)}
}
