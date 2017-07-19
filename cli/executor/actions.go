package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"pkg.re/essentialkaos/ek.v9/mathutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionWait is action processor for "exit"
func actionWait(action *recipe.Action) bool {
	durSec, err := action.GetF(0)

	if err != nil {
		return false
	}

	durSec = mathutil.BetweenF64(durSec, 0.01, 3600.0)

	time.Sleep(secondsToDuration(durSec))

	return true
}

// actionExpect is action processor for "expect"
func actionExpect(action *recipe.Action, output *outputStore) bool {
	var (
		err     error
		start   time.Time
		substr  string
		maxWait float64
	)

	substr, err = action.GetS(0)

	if err != nil {
		return false
	}

	if len(action.Arguments) > 1 {
		maxWait, err = action.GetF(1)

		if err != nil {
			return false
		}
	} else {
		maxWait = 5.0
	}

	maxWait = mathutil.BetweenF64(maxWait, 0.01, 3600.0)
	start = time.Now()

	for {
		if bytes.Contains(output.data.Bytes(), []byte(substr)) {
			return true
		}

		if time.Since(start) >= secondsToDuration(maxWait) {
			return false
		}

		time.Sleep(15 * time.Millisecond)
	}

	return false
}

// actionInput is action processor for "input"
func actionInput(action *recipe.Action, input io.Writer) bool {
	text, err := action.GetS(0)

	if err != nil {
		return false
	}

	if !strings.HasSuffix(text, "\n") {
		text = text + "\n"
	}

	_, err = input.Write([]byte(text))

	return err == nil
}

// actionExit is action processor for "exit"
func actionExit(action *recipe.Action, cmd *exec.Cmd) bool {
	var (
		err      error
		start    time.Time
		exitCode int
		maxWait  float64
	)

	go cmd.Wait()

	exitCode, err = action.GetI(0)

	if err != nil {
		return false
	}

	if len(action.Arguments) > 1 {
		maxWait, err = action.GetF(1)

		if err != nil {
			return false
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
			return false
		}
	}

	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)

	if !ok {
		return false
	}

	return status.ExitStatus() == exitCode
}

// actionOutputEqual is action processor for "output-equal"
func actionOutputEqual(action *recipe.Action, output *outputStore) bool {
	data, err := action.GetS(0)

	if err != nil {
		return false
	}

	return output.String() == data
}

// actionOutputContain is action processor for "output-contain"
func actionOutputContain(action *recipe.Action, output *outputStore) bool {
	data, err := action.GetS(0)

	if err != nil {
		return false
	}

	return strings.Contains(output.String(), data)
}

// actionOutputPrefix is action processor for "output-prefix"
func actionOutputPrefix(action *recipe.Action, output *outputStore) bool {
	data, err := action.GetS(0)

	if err != nil {
		return false
	}

	return strings.HasPrefix(output.String(), data)
}

// actionOutputSuffix is action processor for "output-suffix"
func actionOutputSuffix(action *recipe.Action, output *outputStore) bool {
	data, err := action.GetS(0)

	if err != nil {
		return false
	}

	return strings.HasSuffix(output.String(), data)
}

// actionOutputLength is action processor for "output-length"
func actionOutputLength(action *recipe.Action, output *outputStore) bool {
	size, err := action.GetI(0)

	if err != nil {
		return false
	}

	return len(output.String()) == size
}
