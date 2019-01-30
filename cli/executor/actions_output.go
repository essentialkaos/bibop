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
	var (
		err     error
		start   time.Time
		substr  string
		maxWait float64
	)

	substr, err = action.GetS(0)

	if err != nil {
		return err
	}

	if len(action.Arguments) > 1 {
		maxWait, err = action.GetF(1)

		if err != nil {
			return err
		}
	} else {
		maxWait = 5.0
	}

	maxWait = mathutil.BetweenF64(maxWait, 0.01, 3600.0)
	start = time.Now()

	for {
		if bytes.Contains(output.data.Bytes(), []byte(substr)) {
			return nil
		}

		if time.Since(start) >= secondsToDuration(maxWait) {
			return fmt.Errorf("Reached max wait time (%g sec)", maxWait)
		}

		time.Sleep(15 * time.Millisecond)
	}
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
