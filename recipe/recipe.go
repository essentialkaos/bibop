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

// TokenInfo contains info about supported token
type TokenInfo struct {
	Keyword       string
	MinArgs       int
	MaxArgs       int
	Global        bool
	AllowNegative bool
}

// Tokens is slice with tokens info
var Tokens = []TokenInfo{
	{"var", 2, 2, true, false},

	{"dir", 1, 1, true, false},
	{"unsafe-actions", 1, 1, true, false},
	{"require-root", 1, 1, true, false},
	{"command", 1, 2, true, false},

	{"exit", 1, 2, false, true},
	{"expect", 1, 2, false, false},
	{"output-match", 1, 1, false, true},
	{"output-prefix", 1, 1, false, true},
	{"output-suffix", 1, 1, false, true},
	{"output-length", 1, 1, false, true},
	{"output-contains", 1, 1, false, true},
	{"output-equal", 1, 1, false, true},
	{"output-trim", 0, 0, false, false},
	{"print", 1, 1, false, false},
	{"wait", 1, 1, false, false},

	{"perms", 2, 2, false, true},
	{"owner", 2, 2, false, true},
	{"exist", 1, 1, false, true},
	{"readable", 2, 2, false, true},
	{"writable", 2, 2, false, true},
	{"executable", 2, 2, false, true},
	{"directory", 1, 1, false, true},
	{"empty", 1, 1, false, true},
	{"empty-directory", 1, 1, false, true},

	{"checksum", 2, 2, false, true},
	{"checksum-read", 2, 2, false, false},
	{"file-contains", 2, 2, false, true},

	{"copy", 2, 2, false, false},
	{"move", 2, 2, false, false},
	{"touch", 1, 1, false, false},
	{"mkdir", 1, 1, false, false},
	{"remove", 1, 1, false, false},
	{"chmod", 2, 2, false, false},

	{"process-works", 1, 1, false, true},
	{"wait-pid", 1, 2, false, true},
	{"connect", 2, 2, false, true},
	{"app", 1, 1, false, true},
	{"env", 2, 2, false, true},

	{"user-exist", 1, 1, false, true},
	{"user-id", 2, 2, false, true},
	{"user-gid", 2, 2, false, true},
	{"user-group", 2, 2, false, true},
	{"user-shell", 2, 2, false, true},
	{"user-home", 2, 2, false, true},
	{"group-exist", 1, 1, false, true},
	{"group-id", 2, 2, false, true},

	{"service-present", 1, 1, false, true},
	{"service-enabled", 1, 1, false, true},
	{"service-works", 1, 1, false, true},

	{"http-status", 3, 3, false, true},
	{"http-header", 4, 4, false, true},
	{"http-contains", 3, 3, false, true},

	{"lib-loaded", 1, 1, false, true},
	{"lib-header", 1, 1, false, true},
}

// ////////////////////////////////////////////////////////////////////////////////// //

// varRegex is regexp for parsing variables
var varRegex = regexp.MustCompile(`\{([a-zA-Z0-9_-]+)\}`)

// ////////////////////////////////////////////////////////////////////////////////// //

// NewRecipe create new recipe struct
func NewRecipe(file string) *Recipe {
	return &Recipe{File: file}
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
