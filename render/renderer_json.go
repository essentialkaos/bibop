package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// JSONRenderer is JSON renderer
type JSONRenderer struct {
	report     *report
	curCommand *command
}

// ////////////////////////////////////////////////////////////////////////////////// //

type report struct {
	RecipeInfo *recipeInfo `json:"recipe"`
	Commands   []*command  `json:"commands"`
	Results    *results    `json:"result"`
}

type recipeInfo struct {
	RecipeFile    string `json:"recipe_file"`
	WorkingDir    string `json:"working_dir"`
	UnsafeActions bool   `json:"unsafe_actions"`
	RequireRoot   bool   `json:"require_root"`
	FastFinish    bool   `json:"fast_finish"`
	LockWorkdir   bool   `json:"lock_workdir"`
	Unbuffer      bool   `json:"unbuffer"`
}

type command struct {
	Actions      []*action `json:"actions"`
	User         string    `json:"user,omitempty"`
	Tag          string    `json:"tag,omitempty"`
	Cmdline      string    `json:"cmdline"`
	Description  string    `json:"description"`
	ErrorMessage string    `json:"error_message,omitempty"`
	IsFailed     bool      `json:"is_failed"`
}

type action struct {
	Arguments    []string `json:"arguments"`
	Name         string   `json:"name"`
	ErrorMessage string   `json:"error_message,omitempty"`
	IsFailed     bool     `json:"is_failed"`
}

type results struct {
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Start prints info about started test
func (rr *JSONRenderer) Start(r *recipe.Recipe) {
	rr.report = &report{}
	rr.report.RecipeInfo = &recipeInfo{
		UnsafeActions: r.UnsafeActions,
		RequireRoot:   r.RequireRoot,
		FastFinish:    r.FastFinish,
		LockWorkdir:   r.LockWorkdir,
		Unbuffer:      r.Unbuffer,
	}

	rr.report.RecipeInfo.RecipeFile, _ = filepath.Abs(r.File)
	rr.report.RecipeInfo.WorkingDir, _ = filepath.Abs(r.Dir)
}

// CommandStarted prints info about started command
func (rr *JSONRenderer) CommandStarted(c *recipe.Command) {
	rr.curCommand = rr.convertCommand(c)
}

// CommandSkipped prints info about skipped command
func (rr *JSONRenderer) CommandSkipped(c *recipe.Command) {
	return
}

// CommandFailed prints info about failed command
func (rr *JSONRenderer) CommandFailed(c *recipe.Command, err error) {
	rr.curCommand.IsFailed = true
	rr.curCommand.ErrorMessage = err.Error()

	rr.report.Commands = append(rr.report.Commands, rr.curCommand)
}

// CommandFailed prints info about executed command
func (rr *JSONRenderer) CommandDone(c *recipe.Command, isLast bool) {
	rr.report.Commands = append(rr.report.Commands, rr.curCommand)
}

// ActionStarted prints info about action in progress
func (rr *JSONRenderer) ActionStarted(a *recipe.Action) {
	return
}

// ActionFailed prints info about failed action
func (rr *JSONRenderer) ActionFailed(a *recipe.Action, err error) {
	action := rr.convertAction(a)

	action.IsFailed = true
	action.ErrorMessage = err.Error()

	rr.curCommand.Actions = append(rr.curCommand.Actions, action)
}

// ActionDone prints info about successfully finished action
func (rr *JSONRenderer) ActionDone(a *recipe.Action, isLast bool) {
	rr.curCommand.Actions = append(rr.curCommand.Actions, rr.convertAction(a))
}

// Result prints info about test results
func (rr *JSONRenderer) Result(passes, fails, skips int) {
	rr.report.Results = &results{passes, fails, skips}
	data, _ := json.MarshalIndent(rr.report, "", "  ")
	fmt.Println(string(data))
}

// ////////////////////////////////////////////////////////////////////////////////// //

// convertAction converts command to inner format
func (rr *JSONRenderer) convertCommand(c *recipe.Command) *command {
	return &command{
		User:        c.User,
		Tag:         c.Tag,
		Cmdline:     c.GetCmdline(),
		Description: c.Description,
	}
}

// convertAction converts action to inner format
func (rr *JSONRenderer) convertAction(a *recipe.Action) *action {
	action := &action{}

	if a.Negative {
		action.Name = "!" + a.Name
	} else {
		action.Name = a.Name
	}

	for index := range a.Arguments {
		arg, _ := a.GetS(index)
		action.Arguments = append(action.Arguments, arg)
	}

	return action
}
