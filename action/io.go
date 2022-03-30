package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
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

	"github.com/essentialkaos/ek/v12/mathutil"
	"github.com/essentialkaos/ek/v12/timeutil"

	"github.com/essentialkaos/bibop/output"
	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _DATA_READ_PERIOD = 10 * time.Millisecond

// ////////////////////////////////////////////////////////////////////////////////// //

// Expect is action processor for "expect"
func Expect(action *recipe.Action, outputStore *output.Store) error {
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
	timeoutDur := timeutil.SecondsToDuration(timeout)

	for range time.NewTicker(_DATA_READ_PERIOD).C {
		if bytes.Contains(outputStore.Bytes(), []byte(substr)) {
			outputStore.Purge()
			return nil
		}

		if time.Since(start) >= timeoutDur {
			break
		}
	}

	outputStore.Purge()

	return fmt.Errorf("Timeout (%g sec) reached", timeout)
}

// WaitOutput is action processor for "wait-output"
func WaitOutput(action *recipe.Action, outputStore *output.Store) error {
	timeout, err := action.GetF(0)

	if err != nil {
		return err
	}

	start := time.Now()
	timeoutDur := timeutil.SecondsToDuration(timeout)

	for range time.NewTicker(_DATA_READ_PERIOD).C {
		if !outputStore.IsEmpty() {
			return nil
		}

		if time.Since(start) >= timeoutDur {
			break
		}
	}

	return fmt.Errorf("Timeout (%g sec) reached, but output still empty", timeout)
}

// Input is action processor for "input"
func Input(action *recipe.Action, input *os.File, outputStore *output.Store) error {
	text, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.HasSuffix(text, "\n") {
		text = text + "\n"
	}

	outputStore.Purge()

	_, err = input.Write([]byte(text))

	return err
}

// OutputMatch is action processor for "output-match"
func OutputMatch(action *recipe.Action, outputStore *output.Store) error {
	pattern, err := action.GetS(0)

	if err != nil {
		return err
	}

	rg := regexp.MustCompile(pattern)
	isMatch := rg.Match(outputStore.Bytes())

	switch {
	case !action.Negative && !isMatch:
		return fmt.Errorf("Output doesn't contains data with pattern %s", pattern)
	case action.Negative && isMatch:
		return fmt.Errorf("Output contains data with pattern %s", pattern)
	}

	return nil
}

// OutputContains is action processor for "output-contains"
func OutputContains(action *recipe.Action, outputStore *output.Store) error {
	substr, err := action.GetS(0)

	if err != nil {
		return err
	}

	isMatch := strings.Contains(outputStore.String(), substr)

	switch {
	case !action.Negative && !isMatch:
		return fmt.Errorf("Output doesn't contains substring \"%s\"", substr)
	case action.Negative && isMatch:
		return fmt.Errorf("Output  contains substring \"%s\"", substr)
	}

	return nil
}

// OutputTrim is action processor for "output-trim"
func OutputTrim(action *recipe.Action, outputStore *output.Store) error {
	outputStore.Purge()
	return nil
}
