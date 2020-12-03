package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"pkg.re/essentialkaos/ek.v12/env"
	"pkg.re/essentialkaos/ek.v12/fsutil"
	"pkg.re/essentialkaos/ek.v12/mathutil"
	"pkg.re/essentialkaos/ek.v12/pid"
	"pkg.re/essentialkaos/ek.v12/signal"
	"pkg.re/essentialkaos/ek.v12/timeutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// ErrCantReadPIDFile is returned if PID can't be read
var ErrCantReadPIDFile = fmt.Errorf("Can't read PID from PID file")

// ////////////////////////////////////////////////////////////////////////////////// //

// ProcessWorks is action processor for "process-works"
func ProcessWorks(action *recipe.Action) error {
	pidFile, err := action.GetS(0)

	if err != nil {
		return err
	}

	ppid := pid.Read(pidFile)

	if ppid == -1 {
		return ErrCantReadPIDFile
	}

	isWorks := fsutil.IsExist("/proc/" + strconv.Itoa(ppid))

	switch {
	case !action.Negative && !isWorks:
		return fmt.Errorf("Process with PID %d from PID file %s doesn't exist", ppid, pidFile)
	case action.Negative && isWorks:
		return fmt.Errorf("Process with PID %d from PID file %s exists", ppid, pidFile)
	}

	return err
}

// WaitPID is action processor for "wait-pid"
func WaitPID(action *recipe.Action) error {
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
	timeoutDur := timeutil.SecondsToDuration(timeout)

	for range time.NewTicker(25 * time.Millisecond).C {
		if time.Since(start) >= timeoutDur {
			break
		}

		switch {
		case !action.Negative && fsutil.IsExist(pidFile):
			ppid := pid.Read(pidFile)

			if ppid == -1 {
				continue
			}

			if fsutil.IsExist("/proc/" + strconv.Itoa(ppid)) {
				return nil
			}
		case action.Negative && !fsutil.IsExist(pidFile):
			time.Sleep(time.Second)
			return nil
		}
	}

	switch action.Negative {
	case false:
		return fmt.Errorf("Timeout (%g sec) reached, and PID file %s didn't appear", timeout, pidFile)
	default:
		return fmt.Errorf("Timeout (%g sec) reached, and PID file %s still exists", timeout, pidFile)
	}
}

// WaitFS is action processor for "wait-fs"
func WaitFS(action *recipe.Action) error {
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
	timeoutDur := timeutil.SecondsToDuration(timeout)

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

// Connect is action processor for "connect"
func Connect(action *recipe.Action) error {
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

// App is action processor for "app"
func App(action *recipe.Action) error {
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

// Signal is action processor for "signal"
func Signal(action *recipe.Action, cmd *exec.Cmd) error {
	var sig syscall.Signal
	var sigVal, pidFile string
	var err error

	sigVal, err = action.GetS(0)

	if err != nil {
		return err
	}

	if action.Has(1) {
		pidFile, err = action.GetS(1)

		if err != nil {
			return err
		}
	}

	sig, err = parseSignal(sigVal)

	if err != nil {
		return err
	}

	if pidFile != "" {
		return sendSignalToPID(sig, pidFile)
	}

	return sendSignalToCmd(sig, cmd)
}

// Env is action processor for "env"
func Env(action *recipe.Action) error {
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
		return fmt.Errorf("Environment variable %s has different value (%s ≠ %s)", name, fmtValue(envValue), value)
	case action.Negative && envValue == value:
		return fmt.Errorf("Environment variable %s has invalid value (%s)", name, value)
	}

	return nil
}

// EnvSet is action processor for "env-set"
func EnvSet(action *recipe.Action) error {
	name, err := action.GetS(0)

	if err != nil {
		return err
	}

	value, err := action.GetS(1)

	if err != nil {
		return err
	}

	err = os.Setenv(name, value)

	if err != nil {
		return fmt.Errorf("Can't set environment variable: %v", err)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// isNumber returns true if given value is number
func isNumber(v string) bool {
	if len(v) != 0 && strings.Trim(v, "0123456789") == "" {
		return true
	}

	return false
}

// parseSignal parses signal definition
func parseSignal(v string) (syscall.Signal, error) {
	if isNumber(v) {
		sigCode, err := strconv.Atoi(v)

		if err != nil {
			return syscall.Signal(-1), err
		}

		return signal.GetByCode(sigCode)
	}

	return signal.GetByName(v)
}

// sendSignalToCmd sends signal to PID from command
func sendSignalToCmd(sig syscall.Signal, cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return fmt.Errorf("Can't find process PID (process already dead?)")
	}

	if cmd.ProcessState.Exited() {
		return fmt.Errorf("Can't send signal - process already dead")
	}

	return signal.Send(cmd.ProcessState.Pid(), sig)
}

// sendSignalToPID sends signal to PID from PID file
func sendSignalToPID(sig syscall.Signal, pidFile string) error {
	ppid := pid.Read(pidFile)

	if ppid == -1 {
		return ErrCantReadPIDFile
	}

	return signal.Send(ppid, sig)
}
