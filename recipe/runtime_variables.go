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
	"PYTHON_SITELIB_LOCAL",
	"PYTHON2_SITELIB_LOCAL",
	"PYTHON3_SITELIB_LOCAL",
}

// ////////////////////////////////////////////////////////////////////////////////// //

// dynVarCache is dynamic variables cache
var dynVarCache map[string]string

// ////////////////////////////////////////////////////////////////////////////////// //

// getRuntimeVariable return run-time variable
func getRuntimeVariable(name string, r *Recipe) string {
	if dynVarCache == nil {
		dynVarCache = make(map[string]string)
	}

	value := dynVarCache[name]

	if value != "" {
		return value
	}

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
			dynVarCache[name] = systemInfo.Hostname
		}

	case "IP":
		dynVarCache[name] = netutil.GetIP()

	case "PYTHON_SITELIB", "PYTHON2_SITELIB":
		dynVarCache[name] = getPythonSitePackages("2", false, false)

	case "PYTHON_SITELIB_LOCAL", "PYTHON2_SITELIB_LOCAL":
		dynVarCache[name] = getPythonSitePackages("2", false, true)

	case "PYTHON_SITEARCH", "PYTHON2_SITEARCH":
		dynVarCache[name] = getPythonSitePackages("2", true, false)

	case "PYTHON3_SITELIB":
		dynVarCache[name] = getPythonSitePackages("3", false, false)

	case "PYTHON3_SITELIB_LOCAL":
		dynVarCache[name] = getPythonSitePackages("3", false, true)

	case "PYTHON3_SITEARCH":
		dynVarCache[name] = getPythonSitePackages("3", true, false)

	case "LIBDIR":
		dynVarCache[name] = getLibDir(false)

	case "LIBDIR_LOCAL":
		dynVarCache[name] = getLibDir(true)
	}

	return dynVarCache[name]
}

// getPythonSitePackages return path Python site packages directory
func getPythonSitePackages(version string, arch, local bool) string {
	prefix := "/usr"

	if local {
		prefix = "/usr/local"
	}

	dir := prefix + "/lib"

	if arch && fsutil.IsExist(prefix+"/lib64") {
		dir = prefix + "/lib64"
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
func getLibDir(local bool) string {
	prefix := "/usr"

	if local {
		prefix = "/usr/local"
	}

	if fsutil.IsExist(prefix + "/lib64") {
		return prefix + "/lib64"
	}

	return prefix + "/lib"
}
