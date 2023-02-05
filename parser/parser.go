package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type entity struct {
	info       recipe.TokenInfo
	args       []string
	tag        string
	isNegative bool
	isGroup    bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// tagRegex is regexp for parsing command tag
var tagRegex = regexp.MustCompile(`^\+?command:([a-zA-Z_0-9]+)`)

// ////////////////////////////////////////////////////////////////////////////////// //

// Parse parse bibop suite
func Parse(file string) (*recipe.Recipe, error) {
	err := checkRecipeFile(file)

	if err != nil {
		return nil, err
	}

	return parseRecipeFile(file)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// checkRecipeFile check recipe file
func checkRecipeFile(file string) error {
	if !fsutil.CheckPerms("FR", file) {
		return fmt.Errorf("File %s doesn't exist or not readable", file)
	}

	if !fsutil.IsNonEmpty(file) {
		return fmt.Errorf("File %s is empty", file)
	}

	return nil
}

// parseRecipeFile parce recipe file
func parseRecipeFile(file string) (*recipe.Recipe, error) {
	fd, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	defer fd.Close()

	reader := bufio.NewReader(fd)

	return parseRecipeData(file, reader)
}

// parseRecipeData parse recipe data
func parseRecipeData(file string, reader io.Reader) (*recipe.Recipe, error) {
	var lineNum uint16

	result := recipe.NewRecipe(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines and comments
		if isUselessRecipeLine(line) {
			continue
		}

		e, err := parseLine(line)

		if err != nil {
			return nil, fmt.Errorf("Parsing error in line %d: %v", lineNum, err)
		}

		if !e.info.Global && len(result.Commands) == 0 {
			return nil, fmt.Errorf("Parsing error in line %d: keyword %q is not allowed there", lineNum, e.info.Keyword)
		}

		err = appendData(result, e, lineNum)

		if err != nil {
			return nil, fmt.Errorf("Parsing error in line %d: %v", lineNum, err)
		}
	}

	result.Dir, _ = os.Getwd()

	return result, nil
}

// parseLine parse line from recipe
func parseLine(line string) (*entity, error) {
	var isGlobal bool

	if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
		line = strings.TrimLeft(line, " \t") // Remove prefix
	} else {
		isGlobal = true
	}

	fields := strutil.Fields(line)

	if len(fields) == 0 {
		return nil, fmt.Errorf("Can't parse token data")
	}

	keyword := fields[0]

	info := getTokenInfo(keyword)
	tag := extractTag(keyword)

	if info.Keyword == "" || info.Global != isGlobal {
		switch isGlobal {
		case true:
			return nil, fmt.Errorf("Global keyword %q is not supported", keyword)
		case false:
			return nil, fmt.Errorf("Keyword %q is not supported", keyword)
		}
	}

	isNegative := strings.HasPrefix(keyword, recipe.SYMBOL_NEGATIVE_ACTION)

	if isNegative && !info.AllowNegative {
		return nil, fmt.Errorf("Action %q does not support negative results", keyword)
	}

	isGroup := strings.HasPrefix(keyword, recipe.SYMBOL_COMMAND_GROUP)

	argsNum := len(fields) - 1

	switch {
	case argsNum > info.MaxArgs:
		return nil, fmt.Errorf("Action %q has too many arguments (maximum is %d)", info.Keyword, info.MaxArgs)
	case argsNum < info.MinArgs:
		return nil, fmt.Errorf("Action %q has too few arguments (minimum is %d)", info.Keyword, info.MinArgs)
	}

	return &entity{info, fields[1:], tag, isNegative, isGroup}, nil
}

// appendData append data to recipe struct
func appendData(r *recipe.Recipe, e *entity, line uint16) error {
	if e.info.Global {
		return processGlobalEntity(r, e, line)
	}

	return r.Commands.Last().AddAction(
		&recipe.Action{
			Name:      e.info.Keyword,
			Arguments: e.args,
			Negative:  e.isNegative,
			Line:      line,
		},
	)
}

// processGlobalEntity creates new global entity (variable/command) or appplies
// global option
func processGlobalEntity(r *recipe.Recipe, e *entity, line uint16) error {
	var err error

	switch e.info.Keyword {
	case recipe.KEYWORD_VAR:
		err = r.AddVariable(e.args[0], e.args[1])

	case recipe.KEYWORD_COMMAND:
		if e.isGroup && len(r.Commands) == 0 {
			return fmt.Errorf("Group command (with prefix +) cannot be defined as first in a recipe")
		}

		err = r.AddCommand(recipe.NewCommand(e.args, line), e.tag, e.isGroup)

	case recipe.KEYWORD_PACKAGE:
		r.Packages = e.args

	default:
		err = applyGlobalOption(r, e, line)
	}

	return err
}

// applyGlobalOption applies global options to the recipe
func applyGlobalOption(r *recipe.Recipe, e *entity, line uint16) error {
	var err error

	switch e.info.Keyword {
	case recipe.OPTION_UNSAFE_ACTIONS:
		r.UnsafeActions, err = getOptionBoolValue(e.info.Keyword, e.args[0])

	case recipe.OPTION_REQUIRE_ROOT:
		r.RequireRoot, err = getOptionBoolValue(e.info.Keyword, e.args[0])

	case recipe.OPTION_FAST_FINISH:
		r.FastFinish, err = getOptionBoolValue(e.info.Keyword, e.args[0])

	case recipe.OPTION_LOCK_WORKDIR:
		r.LockWorkdir, err = getOptionBoolValue(e.info.Keyword, e.args[0])

	case recipe.OPTION_UNBUFFER:
		r.Unbuffer, err = getOptionBoolValue(e.info.Keyword, e.args[0])

	case recipe.OPTION_HTTPS_SKIP_VERIFY:
		r.HTTPSSkipVerify, err = getOptionBoolValue(e.info.Keyword, e.args[0])

	case recipe.OPTION_DELAY:
		r.Delay, err = getOptionFloatValue(e.info.Keyword, e.args[0])
	}

	return err
}

// getOptionBoolValue parses option value as boolean
func getOptionBoolValue(keyword, value string) (bool, error) {
	switch strings.ToLower(value) {
	case "false", "no":
		return false, nil
	case "true", "yes":
		return true, nil
	}

	return false, fmt.Errorf("%q is not allowed as value for %s", value, keyword)
}

// getOptionFloatValue parses option value as float number
func getOptionFloatValue(keyword, value string) (float64, error) {
	v, err := strconv.ParseFloat(value, 64)

	if err != nil {
		return 0, fmt.Errorf("%q is not allowed as value for %s: %v", value, keyword, err)
	}

	return v, nil
}

// getTokenInfo return token info by keyword
func getTokenInfo(keyword string) recipe.TokenInfo {
	switch {
	case strings.HasPrefix(keyword, recipe.KEYWORD_COMMAND+recipe.SYMBOL_SEPARATOR),
		strings.HasPrefix(keyword, recipe.SYMBOL_COMMAND_GROUP+recipe.KEYWORD_COMMAND),
		strings.HasPrefix(keyword, recipe.SYMBOL_COMMAND_GROUP+recipe.KEYWORD_COMMAND+recipe.SYMBOL_SEPARATOR):
		keyword = recipe.KEYWORD_COMMAND
	}

	for _, token := range recipe.Tokens {
		switch {
		case token.Keyword == keyword,
			recipe.SYMBOL_NEGATIVE_ACTION+token.Keyword == keyword:
			return token
		}
	}

	return recipe.TokenInfo{}
}

// isUselessRecipeLine return if line doesn't contains recipe data
func isUselessRecipeLine(line string) bool {
	// Skip empty lines
	if line == "" || strings.Replace(line, " ", "", -1) == "" {
		return true
	}

	// Skip comments
	if strings.HasPrefix(strings.Trim(line, " "), "#") {
		return true
	}

	return false
}

// extractTag extracts tag from command
func extractTag(data string) string {
	if !strings.Contains(data, recipe.KEYWORD_COMMAND+recipe.SYMBOL_SEPARATOR) {
		return ""
	}

	if !tagRegex.MatchString(data) {
		return ""
	}

	return tagRegex.FindStringSubmatch(data)[1]
}
