package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strings"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

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
		return fmt.Errorf("Output have different length (%d ≠ %d)", outputSize, size)
	}

	return nil
}

// actionOutputTrim is action processor for "output-trim"
func actionOutputTrim(action *recipe.Action, output *outputStore) error {
	output.clear = true
	return nil
}
