package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v10/mathutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionExpect is action processor for "expect"
func actionExpect(action *recipe.Action, output *OutputStore) error {
	var timeout float64

	substr, err := action.GetS(0)

	if err != nil {
		return err
	}

	if action.Has(1) {
		timeout, err = action.GetF(1)

		if err != nil {
			return err
		}
	} else {
		timeout = 5.0
	}

	start := time.Now()
	timeout = mathutil.BetweenF64(timeout, 0.01, 3600.0)
	timeoutDur := secondsToDuration(timeout)

	stdout := output.Stdout
	stderr := output.Stdout

	for range time.NewTicker(25 * time.Millisecond).C {
		if bytes.Contains(stdout.Bytes(), []byte(substr)) || bytes.Contains(stderr.Bytes(), []byte(substr)) {
			output.Clear = true
			return nil
		}

		if time.Since(start) >= timeoutDur {
			break
		}
	}

	output.Clear = true

	return fmt.Errorf("Timeout (%g sec) reached", timeout)
}

// actionWaitOutput is action processor for "wait-output"
func actionWaitOutput(action *recipe.Action, output *OutputStore) error {
	timeout, err := action.GetF(0)

	if err != nil {
		return err
	}

	start := time.Now()
	timeoutDur := secondsToDuration(timeout)

	for range time.NewTicker(25 * time.Millisecond).C {
		if output.HasData() {
			return nil
		}

		if time.Since(start) >= timeoutDur {
			break
		}
	}

	return fmt.Errorf("Timeout (%g sec) reached, but output still empty", timeout)
}

// actionInput is action processor for "input"
func actionInput(action *recipe.Action, input io.Writer, output *OutputStore) error {
	text, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.HasSuffix(text, "\n") {
		text = text + "\n"
	}

	_, err = input.Write([]byte(text))

	output.Clear = true

	return err
}

// actionOutputMatch is action processor for "output-match"
func actionOutputMatch(action *recipe.Action, output *OutputStore) error {
	pattern, err := action.GetS(0)

	if err != nil {
		return err
	}

	rg := regexp.MustCompile(pattern)
	isMatch := rg.Match(output.Stdout.Bytes()) || rg.Match(output.Stderr.Bytes())

	switch {
	case !action.Negative && !isMatch:
		return fmt.Errorf("Output doesn't contains data with pattern %s", pattern)
	case action.Negative && isMatch:
		return fmt.Errorf("Output contains data with pattern %s", pattern)
	}

	return nil
}

// actionOutputContains is action processor for "output-contains"
func actionOutputContains(action *recipe.Action, output *OutputStore) error {
	substr, err := action.GetS(0)

	if err != nil {
		return err
	}

	stdout := output.Stdout.String()
	stderr := output.Stderr.String()

	isMatch := strings.Contains(stdout, substr) || strings.Contains(stderr, substr)

	switch {
	case !action.Negative && !isMatch:
		return fmt.Errorf("Output doesn't contains substring \"%s\"", substr)
	case action.Negative && isMatch:
		return fmt.Errorf("Output  contains substring \"%s\"", substr)
	}

	return nil
}

// actionOutputTrim is action processor for "output-trim"
func actionOutputTrim(action *recipe.Action, output *OutputStore) error {
	output.Clear = true
	return nil
}
