package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v10/env"
	"pkg.re/essentialkaos/ek.v10/fsutil"
	"pkg.re/essentialkaos/ek.v10/mathutil"
	"pkg.re/essentialkaos/ek.v10/pluralize"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionProcessWorks is action processor for "process-works"
func actionProcessWorks(action *recipe.Action) error {
	pidFile, err := action.GetS(0)

	if err != nil {
		return err
	}

	pid, err := readPID(pidFile)

	if err != nil {
		return err
	}

	isWorks := fsutil.IsExist("/proc/" + pid)

	switch {
	case !action.Negative && !isWorks:
		return fmt.Errorf("Process with PID %s from PID file %s doesn't exist", pid, pidFile)
	case action.Negative && isWorks:
		return fmt.Errorf("Process with PID %s from PID file %s exists", pid, pidFile)
	}

	return err
}

// actionWaitPID is action processor for "wait-pid"
func actionWaitPID(action *recipe.Action) error {
	pidFile, err := action.GetS(0)

	if err != nil {
		return err
	}

	var timeout int
	var counter int

	if action.Has(1) {
		timeout, err = action.GetI(1)

		if err != nil {
			return err
		}
	} else {
		timeout = 60
	}

	timeout = mathutil.Between(timeout, 1, 3600)

	for range time.NewTicker(time.Second).C {
		if fsutil.IsExist(pidFile) {
			return nil
		}

		switch {
		case !action.Negative && fsutil.IsExist(pidFile):
			return nil
		case action.Negative && !fsutil.IsExist(pidFile):
			time.Sleep(time.Second)
			return nil
		}

		counter++

		if counter > timeout {
			break
		}
	}

	switch action.Negative {
	case false:
		return fmt.Errorf(
			"Timeout (%s) reached, and PID file didn't appear",
			pluralize.Pluralize(timeout, "second", "seconds"),
		)
	default:
		return fmt.Errorf(
			"Timeout (%s) reached, and PID file still exists",
			pluralize.Pluralize(timeout, "second", "seconds"),
		)
	}
}

// actionConnect is action processor for "connect"
func actionConnect(action *recipe.Action) error {
	network, err := action.GetS(0)

	if err != nil {
		return err
	}

	address, err := action.GetS(1)

	if err != nil {
		return err
	}

	conn, err := net.DialTimeout(network, address, time.Second)

	if conn != nil {
		conn.Close()
	}

	switch {
	case !action.Negative && err != nil:
		return fmt.Errorf("Can't connect to %s (%s)", address, network)
	case action.Negative && err == nil:
		return fmt.Errorf("Successfully connected to %s (%s)", address, network)
	}

	return nil
}

// actionApp is action processor for "app"
func actionApp(action *recipe.Action) error {
	appName, err := action.GetS(0)

	if err != nil {
		return err
	}

	switch {
	case !action.Negative && env.Which(appName) == "":
		return fmt.Errorf("Application %s not found in PATH", appName)
	case action.Negative && env.Which(appName) != "":
		return fmt.Errorf("Application %s found in PATH", appName)
	}

	return nil
}

// actionEnv is action processor for "env"
func actionEnv(action *recipe.Action) error {
	name, err := action.GetS(0)

	if err != nil {
		return err
	}

	value, err := action.GetS(1)

	if err != nil {
		return err
	}

	envValue := env.Get().GetS(name)

	switch {
	case !action.Negative && envValue != value:
		return fmt.Errorf("Environment variable %s has different value (%s ≠ %s)", name, envValue, value)
	case action.Negative && envValue == value:
		return fmt.Errorf("Environment variable %s has invalid value (%s)", name, value)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// readPID reads PID from PID file
func readPID(file string) (string, error) {
	pidFileData, err := ioutil.ReadFile(file)

	if err != nil {
		return "", fmt.Errorf("Can't read PID file %s: %v", file, err)
	}

	return strings.TrimRight(string(pidFileData), " \n\r"), nil
}
