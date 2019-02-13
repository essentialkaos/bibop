package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"strconv"
	"time"

	"pkg.re/essentialkaos/ek.v10/netutil"
	"pkg.re/essentialkaos/ek.v10/system"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// getRuntimeVariable return run-time variable
func getRuntimeVariable(name string, r *Recipe) string {
	switch name {
	case "WORKDIR":
		return r.Dir

	case "TIMESTAMP":
		return strconv.FormatInt(time.Now().Unix(), 10)

	case "DATE":
		return time.Now().String()

	case "HOSTNAME":
		systemInfo, err := system.GetSystemInfo()

		if err == nil {
			return systemInfo.Hostname
		}

	case "IP":
		return netutil.GetIP()
	}

	return ""
}
