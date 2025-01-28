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

// Renderer is interface for renderers
type Renderer interface {
	// Start prints info about started test
	Start(r *recipe.Recipe)

	// CommandStarted prints info about started command
	CommandStarted(c *recipe.Command)

	// CommandSkipped prints info about skipped command
	CommandSkipped(c *recipe.Command, isLast bool)

	// CommandFailed prints info about failed command
	CommandFailed(c *recipe.Command, err error)

	// CommandFailed prints info about executed command
	CommandDone(c *recipe.Command, isLast bool)

	// ActionStarted prints info about action in progress
	ActionStarted(a *recipe.Action)

	// ActionFailed prints info about failed action
	ActionFailed(a *recipe.Action, err error)

	// ActionDone prints info about successfully finished action
	ActionDone(a *recipe.Action, isLast bool)

	// Result prints info about test results
	Result(passes, fails, skips int)
}

// ////////////////////////////////////////////////////////////////////////////////// //
