package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"pkg.re/essentialkaos/ek.v10/mathutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionWait is action processor for "exit"
func actionWait(action *recipe.Action) error {
	durSec, err := action.GetF(0)

	if err != nil {
		return err
	}

	durSec = mathutil.BetweenF64(durSec, 0.01, 3600.0)

	time.Sleep(secondsToDuration(durSec))

	return nil
}

// actionExit is action processor for "exit"
func actionExit(action *recipe.Action, cmd *exec.Cmd) error {
	if cmd == nil {
		return nil
	}

	var (
		err      error
		start    time.Time
		exitCode int
		maxWait  float64
	)

	go cmd.Wait()

	exitCode, err = action.GetI(0)

	if err != nil {
		return err
	}

	if len(action.Arguments) > 1 {
		maxWait, err = action.GetF(1)

		if err != nil {
			return err
		}
	} else {
		maxWait = 60.0
	}

	start = time.Now()

	for {
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			break
		}

		if time.Since(start) > secondsToDuration(maxWait) {
			return fmt.Errorf("Reached max wait time (%g sec)", maxWait)
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
