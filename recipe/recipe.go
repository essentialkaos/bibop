package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strconv"

	"pkg.re/essentialkaos/ek.v10/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Recipe contains recipe data
type Recipe struct {
	File        string     // Path to recipe
	Dir         string     // Working dir
	UnsafePaths bool       // Allow insafe paths
	RequireRoot bool       // Require root privileges
	Commands    []*Command // Commands
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
	{"dir", 1, 1, true, false},
	{"unsafe-paths", 1, 1, true, false},
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
	{"readable", 1, 1, false, true},
	{"writable", 1, 1, false, true},
	{"directory", 1, 1, false, true},
	{"empty", 1, 1, false, true},
	{"empty-directory", 1, 1, false, true},

	{"checksum", 2, 2, false, true},
	{"file-contains", 2, 2, false, true},

	{"copy", 2, 2, false, false},
	{"move", 2, 2, false, false},
	{"touch", 1, 1, false, false},
	{"mkdir", 1, 1, false, false},
	{"remove", 1, 1, false, false},
	{"chmod", 2, 2, false, false},

	{"process-works", 1, 1, false, true},

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
	{"http-header", 3, 3, false, true},
	{"http-contains", 3, 3, false, true},

	{"lib-loaded", 1, 1, false, true},
	{"lib-header", 1, 1, false, true},
}

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

// AddCommand append command to command slice
func (r *Recipe) AddCommand(cmd *Command) {
	cmd.Recipe = r
	r.Commands = append(r.Commands, cmd)
}

// AddAction append command to actions slice
func (c *Command) AddAction(action *Action) {
	action.Command = c
	c.Actions = append(c.Actions, action)
}

// GetFullCommand return full command
func (c *Command) GetFullCommand() []string {
	return strutil.Fields(c.Cmdline)
}

// GetS return argument with given index as string
func (a *Action) GetS(index int) (string, error) {
	if len(a.Arguments) < index {
		return "", fmt.Errorf("Action doesn't have arguments with index %d", index)
	}

	return a.Arguments[index], nil
}

// GetI return argument with given index as int
func (a *Action) GetI(index int) (int, error) {
	if len(a.Arguments) < index {
		return 0, fmt.Errorf("Action doesn't have arguments with index %d", index)
	}

	valI, err := strconv.ParseInt(a.Arguments[index], 10, 64)

	if err != nil {
		return -1, fmt.Errorf("Can't parse integer argument with index %d", index)
	}

	return int(valI), nil
}

// GetF return argument with given index as float64
func (a *Action) GetF(index int) (float64, error) {
	if len(a.Arguments) < index {
		return 0.0, fmt.Errorf("Action doesn't have arguments with index %d", index)
	}

	valF, err := strconv.ParseFloat(a.Arguments[index], 64)

	if err != nil {
		return -1, fmt.Errorf("Can't parse float argument with index %d", index)
	}

	return valF, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //
