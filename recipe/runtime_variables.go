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

	"pkg.re/essentialkaos/ek.v10/fsutil"
	"pkg.re/essentialkaos/ek.v10/netutil"
	"pkg.re/essentialkaos/ek.v10/system"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// DynamicVariables contains list of dynamic variables
var DynamicVariables = []string{
	"WORKDIR",
	"TIMESTAMP",
	"DATE",
	"HOSTNAME",
	"IP",
	"PYTHON_SITELIB",
	"PYTHON2_SITELIB",
	"PYTHON_SITEARCH",
	"PYTHON2_SITEARCH",
	"PYTHON3_SITELIB",
	"PYTHON3_SITEARCH",
}

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

	case "PYTHON_SITELIB", "PYTHON2_SITELIB":
		return getPythonSitePackages("2", false)

	case "PYTHON_SITEARCH", "PYTHON2_SITEARCH":
		return getPythonSitePackages("2", true)

	case "PYTHON3_SITELIB":
		return getPythonSitePackages("3", false)

	case "PYTHON3_SITEARCH":
		return getPythonSitePackages("3", true)

	case "LIBDIR":
		return getLibDir()
	}

	return ""
}

// getPythonSitePackages return path Python site packages directory
func getPythonSitePackages(version string, arch bool) string {
	dir := "/usr/lib"

	if arch && fsutil.IsExist("/usr/lib64") {
		dir = "/usr/lib64"
	}

	dirList := fsutil.List(dir, true,
		fsutil.ListingFilter{
			MatchPatterns: []string{"python" + version + ".*"},
		},
	)

	if len(dirList) == 0 {
		return ""
	}

	return dir + "/" + dirList[0] + "/site-packages"
}

// getLibDir returns path to directory with libs
func getLibDir() string {
	if fsutil.IsExist("/usr/lib64") {
		return "/usr/lib64"
	}

	return "/usr/lib"
}
