package support

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
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
		fmtc.Printf("  {*}%-16s{!} %s\n", "Name:", formatValue(osInfo.Name))
		fmtc.Printf("  {*}%-16s{!} %s\n", "Pretty Name:", formatValue(osInfo.PrettyName))
		fmtc.Printf("  {*}%-16s{!} %s\n", "Version:", formatValue(osInfo.VersionID))
		fmtc.Printf("  {*}%-16s{!} %s\n", "ID:", formatValue(osInfo.ID))
		fmtc.Printf("  {*}%-16s{!} %s\n", "ID Like:", formatValue(osInfo.IDLike))
		fmtc.Printf("  {*}%-16s{!} %s\n", "Version ID:", formatValue(osInfo.VersionID))
		fmtc.Printf("  {*}%-16s{!} %s\n", "Version Code:", formatValue(osInfo.VersionCodename))
		fmtc.Printf("  {*}%-16s{!} %s\n", "CPE:", formatValue(osInfo.CPEName))
	}

	systemInfo, err := system.GetSystemInfo()

	if err != nil {
		return
	} else {
		if osInfo == nil {
			fmtutil.Separator(false, "SYSTEM INFO")
			fmtc.Printf("  {*}%-16s{!} %s\n", "Name:", formatValue(systemInfo.OS))
			fmtc.Printf("  {*}%-16s{!} %s\n", "Version:", formatValue(systemInfo.Version))
		}
	}

	fmtc.Printf("  {*}%-16s{!} %s\n", "Arch:", formatValue(systemInfo.Arch))
	fmtc.Printf("  {*}%-16s{!} %s\n", "Kernel:", formatValue(systemInfo.Kernel))

	containerEngine := "No"

	switch {
	case fsutil.IsExist("/.dockerenv"):
		containerEngine = "Yes (Docker)"
	case fsutil.IsExist("/run/.containerenv"):
		containerEngine = "Yes (Podman)"
	}

	fmtc.NewLine()
	fmtc.Printf("  {*}%-16s{!} %s\n", "Container:", containerEngine)
}
