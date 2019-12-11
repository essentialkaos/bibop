package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"pkg.re/essentialkaos/ek.v11/fsutil"
	"pkg.re/essentialkaos/ek.v11/strutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type entity struct {
	info       recipe.TokenInfo
	args       []string
	tag        string
	isNegative bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// tagRegex is regexp for parsing command tag
var tagRegex = regexp.MustCompile(`^command:([a-zA-Z_0-9]+)`)

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
			return nil, fmt.Errorf("Parsing error in line %d: keyword \"%s\" is not allowed there", lineNum, e.info.Keyword)
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

	cmd := strutil.Fields(line)

	if len(cmd) == 0 {
		return nil, fmt.Errorf("Can't parse token data")
	}

	info := getTokenInfo(cmd[0])
	tag := extractTag(cmd[0])

	if info.Keyword == "" || info.Global != isGlobal {
		switch isGlobal {
		case true:
			return nil, fmt.Errorf("Global keyword \"%s\" is not supported", cmd[0])
		case false:
			return nil, fmt.Errorf("Keyword \"%s\" is not supported", cmd[0])
		}
	}

	isNegative := strings.HasPrefix(cmd[0], "!")

	if isNegative && !info.AllowNegative {
		return nil, fmt.Errorf("Action \"%s\" does not support negative results", cmd[0])
	}

	argsNum := len(cmd) - 1

	switch {
	case argsNum > info.MaxArgs:
		return nil, fmt.Errorf("Action \"%s\" has too many arguments (maximum is %d)", info.Keyword, info.MaxArgs)
	case argsNum < info.MinArgs:
		return nil, fmt.Errorf("Action \"%s\" has too few arguments (minimum is %d)", info.Keyword, info.MinArgs)
	}

	return &entity{info, cmd[1:], tag, isNegative}, nil
}

// appendData append data to recipe struct
func appendData(r *recipe.Recipe, e *entity, line uint16) error {
	if e.info.Global {
		return applyGlobalOptions(r, e, line)
	}

	action := &recipe.Action{
		Name:      e.info.Keyword,
		Arguments: e.args,
		Negative:  e.isNegative,
		Line:      line,
	}

	lastCommand := r.Commands[len(r.Commands)-1]
	lastCommand.AddAction(action)

	return nil
}

// applyGlobalOptions applies global options to recipe
func applyGlobalOptions(r *recipe.Recipe, e *entity, line uint16) error {
	var err error

	switch e.info.Keyword {
	case recipe.KEYWORD_VAR:
		r.AddVariable(e.args[0], e.args[1])

	case recipe.KEYWORD_COMMAND:
		r.AddCommand(recipe.NewCommand(e.args, line), e.tag)

	case recipe.KEYWORD_PACKAGE:
		r.Packages = e.args

	case recipe.OPTION_UNSAFE_ACTIONS:
		r.UnsafeActions, err = getOptionBoolValue(e.info.Keyword, e.args[0])

	case recipe.OPTION_REQUIRE_ROOT:
		r.RequireRoot, err = getOptionBoolValue(e.info.Keyword, e.args[0])

	case recipe.OPTION_FAST_FINISH:
		r.FastFinish, err = getOptionBoolValue(e.info.Keyword, e.args[0])

	case recipe.OPTION_LOCK_WORKDIR:
		r.LockWorkdir, err = getOptionBoolValue(e.info.Keyword, e.args[0])
	}

	return err
}

// getOptionBoolValue parse bool option value
func getOptionBoolValue(keyword, value string) (bool, error) {
	switch strings.ToLower(value) {
	case "false", "no":
		return false, nil
	case "true", "yes":
		return true, nil
	}

	return false, fmt.Errorf("\"%s\" is not allowed as value for %s", value, keyword)
}

// getTokenInfo return token info by keyword
func getTokenInfo(keyword string) recipe.TokenInfo {
	if strings.HasPrefix(keyword, recipe.KEYWORD_COMMAND+":") {
		keyword = recipe.KEYWORD_COMMAND
	}

	for _, token := range recipe.Tokens {
		if token.Keyword == keyword || "!"+token.Keyword == keyword {
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
	if !strings.HasPrefix(data, recipe.KEYWORD_COMMAND+":") {
		return ""
	}

	if !tagRegex.MatchString(data) {
		return ""
	}

	return tagRegex.FindStringSubmatch(data)[1]
}
