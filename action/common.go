package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"pkg.re/essentialkaos/ek.v12/mathutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Handler is action handler function
type Handler func(action *recipe.Action) error

// ////////////////////////////////////////////////////////////////////////////////// //

// Wait is action processor for "exit"
func Wait(action *recipe.Action) error {
	durSec, err := action.GetF(0)

	if err != nil {
		return err
	}

	durSec = mathutil.BetweenF64(durSec, 0.01, 3600.0)

	time.Sleep(secondsToDuration(durSec))

	return nil
}

// Exit is action processor for "exit"
func Exit(action *recipe.Action, cmd *exec.Cmd) error {
	if cmd == nil {
		return nil
	}

	var err error
	var start time.Time
	var exitCode int
	var timeout float64

	exitCode, err = action.GetI(0)

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

	start = time.Now()

	for range time.NewTicker(25 * time.Millisecond).C {
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			break
		}

		if time.Since(start) > secondsToDuration(timeout) {
			return fmt.Errorf("Reached timeout (%g sec)", timeout)
		}
	}

	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)

	if !ok {
		return fmt.Errorf("Can't get exit code from process state")
	}

	switch {
	case !action.Negative && status.ExitStatus() != exitCode:
		return fmt.Errorf("The process has exited with invalid exit code (%d ≠ %d)", status.ExitStatus(), exitCode)
	case action.Negative && status.ExitStatus() == exitCode:
		return fmt.Errorf("The process has exited with invalid exit code (%d)", status.ExitStatus())
	}

	return nil
}
