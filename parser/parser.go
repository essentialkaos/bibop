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
	"strings"

	"pkg.re/essentialkaos/ek.v10/fsutil"
	"pkg.re/essentialkaos/ek.v10/strutil"

	"github.com/essentialkaos/bibop/recipe"
)

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
	if !fsutil.IsExist(file) {
		return fmt.Errorf("File %s doesn't exist", file)
	}

	if !fsutil.IsReadable(file) {
		return fmt.Errorf("File %s is not readable", file)
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
	var (
		err        error
		lineNum    int
		token      recipe.TokenInfo
		args       []string
		isNegative bool
		line       string
	)

	result := recipe.NewRecipe(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		lineNum++
		line = scanner.Text()

		// Skip empty lines and comments
		if isUselessRecipeLine(line) {
			continue
		}

		token, args, isNegative, err = parseToken(line)

		if err != nil {
			return nil, fmt.Errorf("Parsing error in line %d: %v", lineNum, err)
		}

		if !token.Global && len(result.Commands) == 0 {
			return nil, fmt.Errorf("Parsing error in line %d: keyword \"%s\" is not allowed there", lineNum, token.Keyword)
		}

		err = appendData(result, token, args, isNegative)

		if err != nil {
			return nil, fmt.Errorf("Parsing error in line %d: %v", lineNum, err)
		}
	}

	if result.Dir == "" {
		result.Dir, err = os.Getwd()

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// parseToken parse line from recipe
func parseToken(line string) (recipe.TokenInfo, []string, bool, error) {
	var isGlobal bool

	if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
		// Remove prefix
		line = strings.TrimLeft(line, " \t")
	} else {
		isGlobal = true
	}

	cmd := strutil.Fields(line)

	if len(cmd) == 0 {
		return recipe.TokenInfo{}, nil, false, fmt.Errorf("Can't parse token data")
	}

	t := getTokenInfo(cmd[0])

	if t.Keyword == "" || t.Global != isGlobal {
		switch isGlobal {
		case true:
			return recipe.TokenInfo{}, nil, false, fmt.Errorf("Global keyword \"%s\" is not supported", cmd[0])
		case false:
			return recipe.TokenInfo{}, nil, false, fmt.Errorf("Keyword \"%s\" is not supported", cmd[0])
		}
	}

	isNegative := strings.HasPrefix(cmd[0], "!")

	if isNegative && !t.AllowNegative {
		return recipe.TokenInfo{}, nil, false, fmt.Errorf("Action \"%s\" does not support negative results", cmd[0])
	}

	argsNum := len(cmd) - 1

	switch {
	case argsNum > t.MaxArgs:
		return recipe.TokenInfo{}, nil, false, fmt.Errorf("Action \"%s\" has too many arguments (maximum is %d)", t.Keyword, t.MaxArgs)
	case argsNum < t.MinArgs:
		return recipe.TokenInfo{}, nil, false, fmt.Errorf("Action \"%s\" has too few arguments (minimum is %d)", t.Keyword, t.MinArgs)
	}

	return t, cmd[1:], isNegative, nil
}

// appendData append data to recipe struct
func appendData(r *recipe.Recipe, t recipe.TokenInfo, args []string, isNegative bool) error {
	if t.Global {
		return applyGlobalOptions(r, t.Keyword, args)
	}

	lastCommand := r.Commands[len(r.Commands)-1]
	lastCommand.AddAction(recipe.NewAction(t.Keyword, args, isNegative))

	return nil
}

// applyGlobalOptions applies global options to recipe
func applyGlobalOptions(r *recipe.Recipe, keyword string, args []string) error {
	var err error

	switch keyword {
	case "dir":
		r.Dir = args[0]

	case "var":
		r.AddVariable(args[0], args[1])

	case "command":
		r.AddCommand(recipe.NewCommand(args))

	case "unsafe-actions":
		r.UnsafeActions, err = getOptionBoolValue(keyword, args[0])
		if err != nil {
			return err
		}

	case "require-root":
		r.RequireRoot, err = getOptionBoolValue(keyword, args[0])
		if err != nil {
			return err
		}
	}

	return nil
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
