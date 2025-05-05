package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
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
func (rr *QuietRenderer) Start(r *recipe.Recipe) {}

// CommandStarted prints info about started command
func (rr *QuietRenderer) CommandStarted(c *recipe.Command) {}

// CommandSkipped prints info about skipped command
func (rr *QuietRenderer) CommandSkipped(c *recipe.Command, isLast bool) {}

// CommandFailed prints info about failed command
func (rr *QuietRenderer) CommandFailed(c *recipe.Command, err error) {}

// CommandFailed prints info about executed command
func (rr *QuietRenderer) CommandDone(c *recipe.Command, isLast bool) {}

// ActionStarted prints info about action in progress
func (rr *QuietRenderer) ActionStarted(a *recipe.Action) {}

// ActionFailed prints info about failed action
func (rr *QuietRenderer) ActionFailed(a *recipe.Action, err error) {}

// ActionDone prints info about successfully finished action
func (rr *QuietRenderer) ActionDone(a *recipe.Action, isLast bool) {}

// Result prints info about test results
func (rr *QuietRenderer) Result(passes, fails, skips int) {}
