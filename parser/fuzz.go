//go:build gofuzz
// +build gofuzz

package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Fuzz(data []byte) int {
	_, err := parseRecipeData("", bytes.NewReader(data))

	if err != nil {
		return 0
	}

	return 1
}
