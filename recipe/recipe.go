package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strconv"

	"pkg.re/essentialkaos/ek.v9/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Recipe is basic bibop recipe
type Recipe struct {
	File     string     // Path to recipe
	Dir      string     // Working dir
	Commands []*Command // Commands
}

// Command contains command with all actions
type Command struct {
	Cmdline     string
	Description string
	Actions     []*Action
}

// Action contains action name and slice with arguments
type Action struct {
	Name      string
	Arguments []string
}

// ////////////////////////////////////////////////////////////////////////////////// //

// TokenInfo contains info about supported token
type TokenInfo struct {
	Keyword string
	Global  bool
	MinArgs int
	MaxArgs int
}

var Tokens = []TokenInfo{
	{"dir", true, 1, 1},
	{"command", true, 1, 2},

	{"exit", false, 1, 2},
	{"expect", false, 1, 2},
	{"output-match", false, 1, 1},
	{"output-prefix", false, 1, 1},
	{"output-suffix", false, 1, 1},
	{"output-length", false, 1, 1},
	{"output-contains", false, 1, 1},
	{"output-equal", false, 1, 1},
	{"output-trim", false, 0, 0},
	{"print", false, 1, 1}, {"input", false, 1, 1},
	{"wait", false, 1, 1}, {"sleep", false, 1, 1},

	{"perms", false, 2, 2},
	{"owner", false, 2, 2},
	{"exist", false, 1, 1},
	{"readable", false, 1, 1},
	{"writable", false, 1, 1},
	{"directory", false, 1, 1},
	{"empty", false, 1, 1},
	{"empty-directory", false, 1, 1},
	{"not-exist", false, 1, 1},
	{"not-readable", false, 1, 1},
	{"not-writable", false, 1, 1},
	{"not-directory", false, 1, 1},
	{"not-empty", false, 1, 1},
	{"not-empty-directory", false, 1, 1},

	{"checksum", false, 2, 2},
	{"file-contains", false, 2, 2},

	{"copy", false, 2, 2},
	{"move", false, 2, 2},
	{"touch", false, 1, 1},
	{"mkdir", false, 1, 1},
	{"remove", false, 1, 1},
	{"chmod", false, 2, 2},

	{"process-works", false, 1, 1},
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
func NewAction(name string, args []string) *Action {
	return &Action{name, args}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// GetFullCommand return full command
func (c *Command) GetFullCommand() []string {
	return strutil.Fields(c.Cmdline)
}

// GetS return argument with given index as string
func (a *Action) GetS(index int) (string, error) {
	if len(a.Arguments) < index {
		return "", fmt.Errorf("Action %s doesn't have arguments with index %d", a.Name, index)
	}

	return a.Arguments[index], nil
}

// GetI return argument with given index as int
func (a *Action) GetI(index int) (int, error) {
	if len(a.Arguments) < index {
		return 0, fmt.Errorf("Action %s doesn't have arguments with index %d", a.Name, index)
	}

	valI, err := strconv.ParseInt(a.Arguments[index], 10, 64)

	if err != nil {
		return -1, fmt.Errorf("Can't parse integer argument with index %d in action %s", index, a.Name)
	}

	return int(valI), nil
}

// GetF return argument with given index as float64
func (a *Action) GetF(index int) (float64, error) {
	if len(a.Arguments) < index {
		return 0.0, fmt.Errorf("Action %s doesn't have arguments with index %d", a.Name, index)
	}

	valF, err := strconv.ParseFloat(a.Arguments[index], 64)

	if err != nil {
		return -1, fmt.Errorf("Can't parse float argument with index %d in action %s", index, a.Name)
	}

	return valF, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //
