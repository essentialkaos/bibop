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
	"github.com/essentialkaos/ek/v12/system"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// showOSInfo shows verbose information about system
func showOSInfo() {
	systemInfo, err := system.GetSystemInfo()

	if err != nil {
		return
	}

	fmtutil.Separator(false, "SYSTEM INFO")

	fmtc.Printf("  {*}%-16s{!} %s\n", "Name:", formatValue(systemInfo.OS))
	fmtc.Printf("  {*}%-16s{!} %s\n", "Version:", formatValue(systemInfo.Version))
	fmtc.Printf("  {*}%-16s{!} %s\n", "Arch:", formatValue(systemInfo.Arch))
	fmtc.Printf("  {*}%-16s{!} %s\n", "Kernel:", formatValue(systemInfo.Kernel))
}
