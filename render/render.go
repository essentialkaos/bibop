package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Renderer is interface for renderers
type Renderer interface {
	Start(r *recipe.Recipe)
	CommandStarted(c *recipe.Command)
	CommandFailed(c *recipe.Command, err error)
	CommandDone(c *recipe.Command, isLast bool)
	ActionStarted(a *recipe.Action)
	ActionFailed(a *recipe.Action, err error)
	ActionDone(a *recipe.Action, isLast bool)
	Result(passes, fails int)
}

// ////////////////////////////////////////////////////////////////////////////////// //
