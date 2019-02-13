package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"pkg.re/essentialkaos/ek.v10/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Recipe contains recipe data
type Recipe struct {
	File          string     // Path to recipe
	Dir           string     // Working dir
	UnsafeActions bool       // Allow unsafe actions
	RequireRoot   bool       // Require root privileges
	FastFinish    bool       // Fast finish flag
	LockWorkdir   bool       // Locking workdir flag
	Commands      []*Command // Commands

	variables map[string]*Variable // Variables
}

// Command contains command with all actions
type Command struct {
	Cmdline     string
	Description string
	Actions     []*Action

	Recipe *Recipe
}

// Action contains action name and slice with arguments
type Action struct {
	Name      string
	Arguments []string
	Negative  bool

	Command *Command
}

type Variable struct {
	Value    string
	ReadOnly bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// varRegex is regexp for parsing variables
var varRegex = regexp.MustCompile(`\{([a-zA-Z0-9_-]+)\}`)

// ////////////////////////////////////////////////////////////////////////////////// //

// NewRecipe create new recipe struct
func NewRecipe(file string) *Recipe {
	return &Recipe{
		File:        file,
		LockWorkdir: true,
	}
}

// NewCommand create new command struct
func NewCommand(args []string) *Command {
	command := &Command{}

	switch len(args) {
	case 2:
		command.Cmdline = args[0]
		command.Description = args[1]
	case 1:
		command.Cmdline = args[0]
	}

	return command
}

// NewAction create new action struct
func NewAction(name string, args []string, isNegative bool) *Action {
	return &Action{name, args, isNegative, nil}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// AddCommand appends command to command slice
func (r *Recipe) AddCommand(cmd *Command) {
	cmd.Recipe = r
	r.Commands = append(r.Commands, cmd)
}

// AddVariable adds new RO variable
func (r *Recipe) AddVariable(name, value string) {
	if r.variables == nil {
		r.variables = make(map[string]*Variable)
	}

	r.variables[name] = &Variable{value, true}
}

// SetVariable sets RW variable
func (r *Recipe) SetVariable(name, value string) error {
	if r.variables == nil {
		r.variables = make(map[string]*Variable)
	}

	varInfo, ok := r.variables[name]

	if !ok {
		r.variables[name] = &Variable{value, false}
		return nil
	}

	if !varInfo.ReadOnly {
		r.variables[name].Value = value
		return nil
	}

	return fmt.Errorf("Can't set read-only variable %s", name)
}

// GetVariable returns variable value as string
func (r *Recipe) GetVariable(name string) string {
	if r.variables == nil {
		return ""
	}

	varInfo, ok := r.variables[name]

	if !ok {
		return ""
	}

	return varInfo.Value
}

// ////////////////////////////////////////////////////////////////////////////////// //

// AddAction appends command to actions slice
func (c *Command) AddAction(action *Action) {
	action.Command = c
	c.Actions = append(c.Actions, action)
}

// Arguments returns command line arguments, including the command as [0]
func (c *Command) Arguments() []string {
	return strutil.Fields(renderVars(c.Cmdline, c.Recipe.variables))
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Has returns true if an argument with given is exist
func (a *Action) Has(index int) bool {
	if len(a.Arguments) <= index {
		return false
	}

	return true
}

// GetS returns argument with given index as string
func (a *Action) GetS(index int) (string, error) {
	if !a.Has(index) {
		return "", fmt.Errorf("Action doesn't have arguments with index %d", index)
	}

	data := a.Arguments[index]

	if isVariable(data) {
		return renderVars(data, a.Command.Recipe.variables), nil
	}

	return data, nil
}

// GetI returns argument with given index as int
func (a *Action) GetI(index int) (int, error) {
	data, err := a.GetS(index)

	if err != nil {
		return 0, err
	}

	valI, err := strconv.ParseInt(data, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("Can't parse integer argument with index %d", index)
	}

	return int(valI), nil
}

// GetF returns argument with given index as float64
func (a *Action) GetF(index int) (float64, error) {
	data, err := a.GetS(index)

	if err != nil {
		return 0.0, err
	}

	valF, err := strconv.ParseFloat(data, 64)

	if err != nil {
		return 0.0, fmt.Errorf("Can't parse float argument with index %d", index)
	}

	return valF, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// isVariable returns true if given data is variable definition
func isVariable(data string) bool {
	return strings.Contains(data, "{") && strings.Contains(data, "}")
}

// renderVars renders variables in given string
func renderVars(data string, vars map[string]*Variable) string {
	if len(vars) == 0 {
		return data
	}

	for _, found := range varRegex.FindAllStringSubmatch(data, -1) {
		varInfo, hasVar := vars[found[1]]

		if !hasVar {
			continue
		}

		data = strings.Replace(data, found[0], varInfo.Value, -1)
	}

	return data
}
