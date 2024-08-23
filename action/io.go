package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v13/mathutil"
	"github.com/essentialkaos/ek/v13/timeutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _DATA_READ_PERIOD = 10 * time.Millisecond

// ////////////////////////////////////////////////////////////////////////////////// //

// Expect is action processor for "expect"
func Expect(action *recipe.Action, output *OutputContainer) error {
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
	timeout = mathutil.Between(timeout, 0.01, 3600.0)
	timeoutDur := timeutil.SecondsToDuration(timeout)

	for range time.NewTicker(_DATA_READ_PERIOD).C {
		if bytes.Contains(output.Bytes(), []byte(substr)) {
			output.Purge()
			return nil
		}

		if time.Since(start) >= timeoutDur {
			break
		}
	}

	return fmt.Errorf("Timeout (%g sec) reached", timeout)
}

// WaitOutput is action processor for "wait-output"
func WaitOutput(action *recipe.Action, output *OutputContainer) error {
	timeout, err := action.GetF(0)

	if err != nil {
		return err
	}

	start := time.Now()
	timeoutDur := timeutil.SecondsToDuration(timeout)

	for range time.NewTicker(_DATA_READ_PERIOD).C {
		if !output.IsEmpty() {
			return nil
		}

		if time.Since(start) >= timeoutDur {
			break
		}
	}

	return fmt.Errorf("Timeout (%g sec) reached, but output still empty", timeout)
}

// Input is action processor for "input"
func Input(action *recipe.Action, input *os.File, output *OutputContainer) error {
	text, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	output.Purge()

	_, err = input.Write([]byte(text))

	return err
}

// OutputMatch is action processor for "output-match"
func OutputMatch(action *recipe.Action, output *OutputContainer) error {
	pattern, err := action.GetS(0)

	if err != nil {
		return err
	}

	rg := regexp.MustCompile(pattern)
	isMatch := rg.Match(output.Bytes())

	switch {
	case !action.Negative && !isMatch:
		return fmt.Errorf("Output doesn't contains data with pattern %q", pattern)
	case action.Negative && isMatch:
		return fmt.Errorf("Output contains data with pattern %q", pattern)
	}

	return nil
}

// OutputContains is action processor for "output-contains"
func OutputContains(action *recipe.Action, output *OutputContainer) error {
	substr, err := action.GetS(0)

	if err != nil {
		return err
	}

	isMatch := strings.Contains(output.String(), substr)

	switch {
	case !action.Negative && !isMatch:
		return fmt.Errorf("Output doesn't contains substring %q", substr)
	case action.Negative && isMatch:
		return fmt.Errorf("Output contains substring %q", substr)
	}

	return nil
}

// OutputEmpty is action processor for "output-empty"
func OutputEmpty(action *recipe.Action, output *OutputContainer) error {
	switch {
	case !action.Negative && !output.IsEmpty():
		return fmt.Errorf("Output contains data")
	case action.Negative && output.IsEmpty():
		return fmt.Errorf("Output is empty")
	}

	return nil
}

// OutputTrim is action processor for "output-trim"
func OutputTrim(action *recipe.Action, output *OutputContainer) error {
	output.Purge()
	return nil
}
