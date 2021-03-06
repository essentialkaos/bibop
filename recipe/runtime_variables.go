package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"strconv"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v12/fsutil"
	"pkg.re/essentialkaos/ek.v12/netutil"
	"pkg.re/essentialkaos/ek.v12/strutil"
	"pkg.re/essentialkaos/ek.v12/system"
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
	"PYTHON_SITEARCH_LOCAL",
	"PYTHON2_SITEARCH_LOCAL",
	"PYTHON3_SITEARCH_LOCAL",
	"ERLANG_BIN_DIR",
}

// ////////////////////////////////////////////////////////////////////////////////// //

// dynVarCache is dynamic variables cache
var dynVarCache map[string]string

// prefixDir is path to base prefix dir
var prefixDir = "/usr"

// localPrefixDir is path to base local prefix dir
var localPrefixDir = "/usr/local"

// erlangBaseDir is path to directory with Erlang data
var erlangBaseDir = "/usr/lib64/erlang"

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

	if strings.HasPrefix(name, "ENV:") {
		return getEnvVariable(name)
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

	case "PYTHON_SITEARCH_LOCAL", "PYTHON2_SITEARCH_LOCAL":
		dynVarCache[name] = getPythonSitePackages("2", true, true)

	case "PYTHON3_SITELIB":
		dynVarCache[name] = getPythonSitePackages("3", false, false)

	case "PYTHON3_SITELIB_LOCAL":
		dynVarCache[name] = getPythonSitePackages("3", false, true)

	case "PYTHON3_SITEARCH":
		dynVarCache[name] = getPythonSitePackages("3", true, false)

	case "PYTHON3_SITEARCH_LOCAL":
		dynVarCache[name] = getPythonSitePackages("3", true, true)

	case "LIBDIR":
		dynVarCache[name] = getLibDir(false)

	case "LIBDIR_LOCAL":
		dynVarCache[name] = getLibDir(true)

	case "ERLANG_BIN_DIR":
		dynVarCache[name] = getErlangBinDir()
	}

	return dynVarCache[name]
}

// getPythonSitePackages return path Python site packages directory
func getPythonSitePackages(version string, arch, local bool) string {
	prefix := getPrefixDir(local)
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
	prefix := getPrefixDir(local)

	if fsutil.IsExist(prefix + "/lib64") {
		return prefix + "/lib64"
	}

	return prefix + "/lib"
}

// getPrefixDir returns path to prefix dir
func getPrefixDir(local bool) string {
	switch local {
	case true:
		return localPrefixDir
	default:
		return prefixDir
	}
}

// getErlangBinDir returns path to Erlang bin directory
func getErlangBinDir() string {
	ertsDir := fsutil.List(
		erlangBaseDir, true,
		fsutil.ListingFilter{MatchPatterns: []string{"erts-*"}, Perms: "DX"},
	)

	if len(ertsDir) == 0 {
		return erlangBaseDir + "/erts/bin"
	}

	return erlangBaseDir + "/" + ertsDir[0] + "/bin"
}

// getEnvVariable returns environment variable
func getEnvVariable(name string) string {
	name = strutil.Exclude(name, "ENV:")
	return os.Getenv(name)
}
