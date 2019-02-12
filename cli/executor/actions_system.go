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
	var timeout float64

	pidFile, err := action.GetS(0)

	if err != nil {
		return err
	}

	if action.Has(1) {
		timeout, err = action.GetF(1)

		if err != nil {
			return err
		}
	} else {
		timeout = 60.0
	}

	start := time.Now()
	timeout = mathutil.BetweenF64(timeout, 0.01, 3600.0)
	timeoutDur := secondsToDuration(timeout)

	for range time.NewTicker(25 * time.Millisecond).C {
		switch {
		case !action.Negative && fsutil.IsExist(pidFile):
			pid, err := readPID(pidFile)

			if err != nil {
				return err
			}

			if fsutil.IsExist("/proc/" + pid) {
				return nil
			}
		case action.Negative && !fsutil.IsExist(pidFile):
			time.Sleep(time.Second)
			return nil
		}

		if time.Since(start) >= timeoutDur {
			break
		}
	}

	switch action.Negative {
	case false:
		return fmt.Errorf("Timeout (%g sec) reached, and PID file %s didn't appear", pidFile)
	default:
		return fmt.Errorf("Timeout (%g sec) reached, and PID file %s still exists", pidFile)
	}
}

// actionWaitFS is action processor for "wait-fs"
func actionWaitFS(action *recipe.Action) error {
	var timeout float64

	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if action.Has(1) {
		timeout, err = action.GetF(1)

		if err != nil {
			return err
		}
	} else {
		timeout = 60.0
	}

	start := time.Now()
	timeout = mathutil.BetweenF64(timeout, 0.01, 3600.0)
	timeoutDur := secondsToDuration(timeout)

	for range time.NewTicker(25 * time.Millisecond).C {
		switch {
		case !action.Negative && fsutil.IsExist(file),
			action.Negative && !fsutil.IsExist(file):
			return nil
		}

		if time.Since(start) >= timeoutDur {
			break
		}
	}

	switch action.Negative {
	case false:
		return fmt.Errorf("Timeout (%g sec) reached, and %s didn't appear", timeout, file)
	default:
		return fmt.Errorf("Timeout (%g sec) reached, and %s still exists", timeout, file)
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
