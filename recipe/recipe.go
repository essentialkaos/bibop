package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v12/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// MAX_GROUP_ID is maximum group ID
const MAX_GROUP_ID uint8 = 255

// MAX_VAR_NESTING is maximum variables nesting
const MAX_VAR_NESTING int = 32

// MAX_VARIABLE_SIZE is maximum length of variable value
const MAX_VARIABLE_SIZE int = 512

// TEARDOWN_TAG is teardown tag
const TEARDOWN_TAG = "teardown"

// ////////////////////////////////////////////////////////////////////////////////// //

// Recipe contains recipe data
// aligo:ignore
type Recipe struct {
	Packages        []string // Package list
	Commands        Commands // Commands
	File            string   // Path to recipe
	Dir             string   // Working dir
	Delay           float64  // Delay between commands
	UnsafeActions   bool     // Allow unsafe actions
	RequireRoot     bool     // Require root privileges
	FastFinish      bool     // Fast finish flag
	LockWorkdir     bool     // Locking workdir flag
	Unbuffer        bool     // Disabled IO buffering
	HTTPSSkipVerify bool     // Disable certificate verification

	variables map[string]*Variable // Variables
}

// Commands is a slice with commands
type Commands []*Command

// Command contains command with all actions
// aligo:ignore
type Command struct {
	Actions     Actions  // Slice with actions
	User        string   // User name
	Tag         string   // Tag
	Cmdline     string   // Command line
	Description string   // Description
	Env         []string // Environment variables
	Recipe      *Recipe  // Link to recipe
	Line        uint16   // Line in recipe file

	GroupID uint8 // Unique command group ID

	Data *Storage // Data storage
}

// Actions is a slice with actions
type Actions []*Action

// Action contains action name and slice with arguments
type Action struct {
	Arguments []string // Arguments
	Name      string   // Name
	Command   *Command // Link to command
	Line      uint16   // Line in recipe
	Negative  bool     // Negative check flag
}

// Variable contains variable data
type Variable struct {
	Value      string
	IsReadOnly bool
}

// Storage is data storage
type Storage struct {
	data map[string]interface{}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// varRegex is regexp for parsing variables
var varRegex = regexp.MustCompile(`\{([a-zA-Z0-9:_-]+)\}`)

// userRegex is regexp for parsing user in command
var userRegex = regexp.MustCompile(`^([a-zA-Z_0-9\{\}]+):`)

// ////////////////////////////////////////////////////////////////////////////////// //

// NewRecipe create new recipe struct
func NewRecipe(file string) *Recipe {
	return &Recipe{
		File:        file,
		LockWorkdir: true,
	}
}

// NewCommand create new command struct
func NewCommand(args []string, line uint16) *Command {
	return parseCommand(args, line)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// AddCommand appends command to command slice
func (r *Recipe) AddCommand(cmd *Command, tag string, isNested bool) error {
	cmd.Recipe = r
	cmd.Tag = tag

	if cmd.User != "" {
		r.RequireRoot = true

		if isVariable(cmd.User) {
			cmd.User = renderVars(r, cmd.User)
		}
	}

	if isVariable(cmd.Description) {
		cmd.Description = renderVars(r, cmd.Description)
	}

	if len(r.Commands) != 0 {
		if isNested {
			cmd.GroupID = r.Commands.Last().GroupID
		} else {
			cmd.GroupID = r.Commands.Last().GroupID + 1
		}
	}

	r.Commands = append(r.Commands, cmd)

	return nil
}

// AddVariable adds new RO variable
func (r *Recipe) AddVariable(name, value string) error {
	if r.variables == nil {
		r.variables = make(map[string]*Variable)
	}

	if strings.Contains(value, "{"+name+"}") {
		return fmt.Errorf("Can't define variable %q: variable contains itself as a part of value", name)
	}

	r.variables[name] = &Variable{value, true}

	return nil
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

	if !varInfo.IsReadOnly {
		r.variables[name].Value = value
		return nil
	}

	return fmt.Errorf("Can't set read-only variable %s", name)
}

// GetVariable returns variable value as string
func (r *Recipe) GetVariable(name string) string {
	rtv := getRuntimeVariable(name, r)

	if rtv != "" {
		return rtv
	}

	if r.variables == nil {
		return ""
	}

	varInfo, ok := r.variables[name]

	if !ok {
		return ""
	}

	return varInfo.Value
}

// GetPackages flatten packages slice to string
func (r *Recipe) GetPackages() string {
	return strings.Join(r.Packages, " ")
}

// HasTeardown returns true if recipe contains command with teardown tag
func (r *Recipe) HasTeardown() bool {
	for _, c := range r.Commands {
		if c.Tag == TEARDOWN_TAG {
			return true
		}
	}

	return false
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Last returns the last command from slice
func (c Commands) Last() *Command {
	if len(c) == 0 {
		return nil
	}

	return c[len(c)-1]
}

// Has returns true if slice contains command with given index
func (c Commands) Has(index int) bool {
	return c != nil && index < len(c)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// AddAction appends command to actions slice
func (c *Command) AddAction(action *Action) error {
	action.Command = c
	c.Actions = append(c.Actions, action)
	return nil
}

// GetCmdline returns command line with rendered variables
func (c *Command) GetCmdline() string {
	return renderVars(c.Recipe, c.Cmdline)
}

// GetCmdlineArgs returns command line arguments, including the command as [0]
func (c *Command) GetCmdlineArgs() []string {
	return strutil.Fields(c.GetCmdline())
}

// Index returns command index
func (c *Command) Index() int {
	if c.Recipe == nil {
		return -1
	}

	for index, command := range c.Recipe.Commands {
		if command == c {
			return index
		}
	}

	return -1
}

// String returns string representation of command
func (c *Command) String() string {
	info := fmt.Sprintf("%d: ", c.Index())

	if c.Description != "" {
		info += c.Description + " → "
	}

	if c.User != "" {
		info += fmt.Sprintf("(%s) ", c.User)
	}

	if len(c.Env) != 0 {
		info += fmt.Sprintf("[%s] ", strings.Join(c.Env, " "))
	}

	if c.IsHollow() {
		info += "<HOLLOW>"
	} else {
		info += c.Cmdline
	}

	info += fmt.Sprintf(" | Actions: %d", len(c.Actions))

	return fmt.Sprintf("Command{%s}", info)
}

// IsHollow returns true if the current command is "hollow" i.e., this command
// does not execute any of the binaries on the system
func (c *Command) IsHollow() bool {
	return c.Cmdline == ""
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Set adds new object into storage with given key
func (s *Storage) Set(key string, value interface{}) {
	if s.data == nil {
		s.data = make(map[string]interface{})
	}

	s.data[key] = value
}

// Get returns object with given key from storage
func (s *Storage) Get(key string) interface{} {
	if s.data == nil {
		return ""
	}

	return s.data[key]
}

// Has returns true if storage contains object with given key
func (s *Storage) Has(key string) bool {
	if s.data == nil {
		return false
	}

	_, ok := s.data[key]

	return ok
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Last returns the last action from slice
func (a Actions) Last() *Action {
	if len(a) == 0 {
		return nil
	}

	return a[len(a)-1]
}

// Has returns true if slice contains action with given index
func (a Actions) Has(index int) bool {
	return a != nil && index < len(a)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Index returns action index
func (a *Action) Index() int {
	if a.Command == nil {
		return -1
	}

	for index, action := range a.Command.Actions {
		if action == a {
			return index
		}
	}

	return -1
}

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
		return renderVars(a.Command.Recipe, data), nil
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

// parseCommand parse command data
func parseCommand(args []string, line uint16) *Command {
	var cmdline, desc, user string
	var envs []string

	switch len(args) {
	case 2:
		cmdline, desc = args[0], args[1]
	case 1:
		cmdline = args[0]
	}

	if userRegex.MatchString(cmdline) {
		matchData := userRegex.FindStringSubmatch(cmdline)
		cmdline = userRegex.ReplaceAllString(cmdline, "")
		user = matchData[1]
	}

	cmdline = strings.TrimSpace(cmdline)
	cmdline, envs = extractEnvVariables(cmdline)

	return &Command{
		Cmdline:     cmdline,
		Env:         envs,
		Description: desc,
		User:        user,
		Line:        line,

		Data: &Storage{},
	}
}

// extractEnvVariables separates command line from environment variables
func extractEnvVariables(cmdline string) (string, []string) {
	if cmdline == "" || cmdline == "-" {
		return "", nil
	}

	if !strings.Contains(cmdline, "=") {
		return cmdline, nil
	}

	var envs []string

	for {
		variable := strutil.ReadField(cmdline, 0, false, " ")

		if !strings.Contains(variable, "=") {
			break
		}

		envs = append(envs, variable)
		cmdline = strutil.Substr(cmdline, len(variable)+1, 99999)
	}

	return cmdline, envs
}

// isVariable returns true if given data is variable definition
func isVariable(data string) bool {
	return strings.Contains(data, "{") && strings.Contains(data, "}")
}

// renderVars renders variables in given string
func renderVars(r *Recipe, data string) string {
	if r == nil {
		return data
	}

	for i := 0; i < MAX_VAR_NESTING; i++ {
		for _, found := range varRegex.FindAllStringSubmatch(data, -1) {
			varValue := r.GetVariable(found[1])

			if varValue == "" {
				continue
			}

			data = strings.ReplaceAll(data, found[0], varValue)

			if len(data) > MAX_VARIABLE_SIZE {
				return data
			}
		}
	}

	return data
}
