// +build linux, !darwin, !windows

package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	CLI "github.com/essentialkaos/bibop/cli"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	CLI.Init()
}
