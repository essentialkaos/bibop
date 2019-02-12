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
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v10/mathutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionExpect is action processor for "expect"
func actionExpect(action *recipe.Action, output *outputStore) error {
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

	for range time.NewTicker(25 * time.Millisecond).C {
		if bytes.Contains(output.data.Bytes(), []byte(substr)) {
			return nil
		}

		if time.Since(start) >= timeoutDur {
			break
		}
	}

	return fmt.Errorf("Timeout (%g sec) reached", timeout)
}

// actionWaitOutput is action processor for "wait-output"
func actionWaitOutput(action *recipe.Action, output *outputStore) error {
	timeout, err := action.GetF(0)

	if err != nil {
		return err
	}

	start := time.Now()
	timeoutDur := secondsToDuration(timeout)

	for range time.NewTicker(25 * time.Millisecond).C {
		if output.data.Len() != 0 {
			return nil
		}

		if time.Since(start) >= timeoutDur {
			break
		}
	}

	return fmt.Errorf("Timeout (%g sec) reached, but output still empty", timeout)
}

// actionInput is action processor for "input"
func actionInput(action *recipe.Action, input io.Writer) error {
	text, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.HasSuffix(text, "\n") {
		text = text + "\n"
	}

	_, err = input.Write([]byte(text))

	return err
}

// actionOutputEqual is action processor for "output-equal"
func actionOutputEqual(action *recipe.Action, output *outputStore) error {
	data, err := action.GetS(0)

	if err != nil {
		return err
	}

	if output.String() != data {
		return fmt.Errorf("Output doesn't equals substring \"%s\"", data)
	}

	return nil
}

// actionOutputContains is action processor for "output-contains"
func actionOutputContains(action *recipe.Action, output *outputStore) error {
	data, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.Contains(output.String(), data) {
		return fmt.Errorf("Output doesn't contains substring \"%s\"", data)
	}

	return nil
}

// actionOutputPrefix is action processor for "output-prefix"
func actionOutputPrefix(action *recipe.Action, output *outputStore) error {
	data, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.HasPrefix(output.String(), data) {
		return fmt.Errorf("Output doesn't have prefix \"%s\"", data)
	}

	return nil
}

// actionOutputSuffix is action processor for "output-suffix"
func actionOutputSuffix(action *recipe.Action, output *outputStore) error {
	data, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.HasSuffix(output.String(), data) {
		return fmt.Errorf("Output doesn't have suffix \"%s\"", data)
	}

	return nil
}

// actionOutputLength is action processor for "output-length"
func actionOutputLength(action *recipe.Action, output *outputStore) error {
	size, err := action.GetI(0)

	if err != nil {
		return err
	}

	outputSize := len(output.String())

	if outputSize == size {
		return fmt.Errorf("Output has different length (%d ≠ %d)", outputSize, size)
	}

	return nil
}

// actionOutputTrim is action processor for "output-trim"
func actionOutputTrim(action *recipe.Action, output *outputStore) error {
	output.clear = true
	return nil
}
