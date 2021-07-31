package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// QuietRenderer doesn't print any info
type QuietRenderer struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

// Start prints info about started test
func (rr *QuietRenderer) Start(r *recipe.Recipe) {
	return
}

// CommandStarted prints info about started command
func (rr *QuietRenderer) CommandStarted(c *recipe.Command) {
	return
}

// CommandFailed prints info about failed command
func (rr *QuietRenderer) CommandFailed(c *recipe.Command, err error) {
	return
}

// CommandFailed prints info about executed command
func (rr *QuietRenderer) CommandDone(c *recipe.Command, isLast bool) {
	return
}

// ActionStarted prints info about action in progress
func (rr *QuietRenderer) ActionStarted(a *recipe.Action) {
	return
}

// ActionFailed prints info about failed action
func (rr *QuietRenderer) ActionFailed(a *recipe.Action, err error) {
	return
}

// ActionDone prints info about successfully finished action
func (rr *QuietRenderer) ActionDone(a *recipe.Action, isLast bool) {
	return
}

// Result prints info about test results
func (rr *QuietRenderer) Result(passes, fails int) {
	return
}
