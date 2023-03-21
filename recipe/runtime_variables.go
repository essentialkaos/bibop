package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/netutil"
	"github.com/essentialkaos/ek/v12/strutil"
	"github.com/essentialkaos/ek/v12/system"
	"github.com/essentialkaos/ek/v12/timeutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// DynamicVariables contains list of dynamic variables
var DynamicVariables = []string{
	"WORKDIR",
	"TIMESTAMP",
	"HOSTNAME",
	"IP",
	"ARCH",
	"ARCH_BITS",
	"ARCH_NAME",
	"OS",
	"LIBDIR",
	"LIBDIR_LOCAL",
	"PYTHON2_VERSION",
	"PYTHON2_SITELIB",
	"PYTHON2_SITEARCH",
	"PYTHON3_VERSION",
	"PYTHON3_SITELIB",
	"PYTHON3_SITEARCH",
	"PYTHON3_BINDING_SUFFIX",
	"ERLANG_BIN_DIR",
}

// ////////////////////////////////////////////////////////////////////////////////// //

// dynVarCache is dynamic variables cache
var dynVarCache map[string]string

// systemInfoCache is cached system info
var systemInfoCache *system.SystemInfo

// prefixDir is path to base prefix dir
var prefixDir = "/usr"

// localPrefixDir is path to base local prefix dir
var localPrefixDir = "/usr/local"

// erlangBaseDir is path to directory with Erlang data
var erlangBaseDir = "/usr/lib64/erlang"

// ////////////////////////////////////////////////////////////////////////////////// //

var python2Bin = "python2"
var python3Bin = "python3"

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

	switch {
	case strings.HasPrefix(name, "ENV:"):
		return getEnvVariable(name)
	case strings.HasPrefix(name, "DATE:"):
		return getDateVariable(name)
	}

	switch name {
	case "WORKDIR":
		return r.Dir

	case "TIMESTAMP":
		return strconv.FormatInt(time.Now().Unix(), 10)

	case "HOSTNAME":
		systemInfo := getSystemInfo()

		if systemInfo != nil {
			dynVarCache[name] = systemInfo.Hostname
		}

	case "ARCH":
		systemInfo := getSystemInfo()

		if systemInfo != nil {
			dynVarCache[name] = systemInfo.Arch
		}

	case "ARCH_BITS":
		systemInfo := getSystemInfo()

		if systemInfo != nil {
			dynVarCache[name] = strconv.Itoa(systemInfo.ArchBits)
		}

	case "ARCH_NAME":
		systemInfo := getSystemInfo()

		if systemInfo != nil {
			dynVarCache[name] = systemInfo.ArchName
		}

	case "OS":
		systemInfo := getSystemInfo()

		if systemInfo != nil {
			dynVarCache[name] = strings.ToLower(systemInfo.OS)
		}

	case "IP":
		dynVarCache[name] = netutil.GetIP()

	case "PYTHON2_VERSION":
		dynVarCache[name] = getPythonVersion(2)

	case "PYTHON2_SITELIB":
		dynVarCache[name] = getPythonSiteLib(2)

	case "PYTHON2_SITEARCH":
		dynVarCache[name] = getPythonSiteArch(2)

	case "PYTHON3_VERSION":
		dynVarCache[name] = getPythonVersion(3)

	case "PYTHON3_SITELIB":
		dynVarCache[name] = getPythonSiteLib(3)

	case "PYTHON3_SITEARCH":
		dynVarCache[name] = getPythonSiteArch(3)

	case "PYTHON3_BINDING_SUFFIX":
		dynVarCache[name] = getPythonBindingSuffix()

	case "LIBDIR":
		dynVarCache[name] = getLibDir(false)

	case "LIBDIR_LOCAL":
		dynVarCache[name] = getLibDir(true)

	case "ERLANG_BIN_DIR":
		dynVarCache[name] = getErlangBinDir()
	}

	return dynVarCache[name]
}

// getSystemInfo returns struct with system info
func getSystemInfo() *system.SystemInfo {
	if systemInfoCache != nil {
		return systemInfoCache
	}

	systemInfoCache, _ = system.GetSystemInfo()

	return systemInfoCache
}

// getPythonVersion returns Python version
func getPythonVersion(majorVersion int) string {
	return evalPythonCode(majorVersion, "import sys; print('{0}.{1}'.format(sys.version_info.major,sys.version_info.minor))")
}

// getPythonSiteLib returns Python site lib
func getPythonSiteLib(majorVersion int) string {
	return evalPythonCode(majorVersion, "from distutils.sysconfig import get_python_lib; print(get_python_lib())")
}

// getPythonSiteArch returns Python site arch
func getPythonSiteArch(majorVersion int) string {
	return evalPythonCode(majorVersion, "from distutils.sysconfig import get_python_lib; print(get_python_lib(plat_specific=True))")
}

// getPythonBindingSuffix returns suffix for Python bindings
func getPythonBindingSuffix() string {
	version := getPythonVersion(3)

	if version == "" {
		return ""
	}

	version = strutil.Exclude(version, ".")
	systemInfo := getSystemInfo()

	if systemInfo == nil {
		return ""
	}

	return fmt.Sprintf(".cpython-%sm-%s-linux-gnu.so", version, systemInfo.Arch)
}

// evalPythonCode evaluates Python code
func evalPythonCode(majorVersion int, code string) string {
	var bin string

	switch majorVersion {
	case 2:
		bin = python2Bin
	case 3:
		bin = python3Bin
	}

	cmd := exec.Command(bin, "-c", code)
	out, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.Trim(string(out), "\r\n")
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

// getDateVariable returns date variable
func getDateVariable(name string) string {
	name = strutil.Exclude(name, "DATE:")
	return timeutil.Format(time.Now(), name)
}
