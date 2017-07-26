package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/strutil"

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

	result := recipe.NewRecipe(file)
	reader := bufio.NewReader(fd)
	scanner := bufio.NewScanner(reader)

	var (
		lineNum int
		token   recipe.TokenInfo
		args    []string
		line    string
	)

	for scanner.Scan() {
		lineNum++
		line = scanner.Text()

		// Skip empty lines and comments
		if isUselessRecipeLine(line) {
			continue
		}

		token, args, err = parseToken(line)

		if err != nil {
			return nil, fmt.Errorf("Parsing error in line %d: %v", lineNum, err)
		}

		if !token.Global && len(result.Commands) == 0 {
			return nil, fmt.Errorf("Parsing error in line %d: keyword \"%s\" is not allowed there", lineNum, token.Keyword)
		}

		err = appendData(result, token, args)

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
func parseToken(line string) (recipe.TokenInfo, []string, error) {
	var isGlobal bool

	if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
		// Remove prefix
		line = strings.TrimLeft(line, " \t")
	} else {
		isGlobal = true
	}

	cmd := strutil.Fields(line)
	t := getTokenInfo(cmd[0])

	if t.Keyword == "" || t.Global != isGlobal {
		switch isGlobal {
		case true:
			return recipe.TokenInfo{}, nil, fmt.Errorf("Global keyword \"%s\" is not supported", cmd[0])
		case false:
			return recipe.TokenInfo{}, nil, fmt.Errorf("Keyword \"%s\" is not supported", cmd[0])
		}
	}

	argsNum := len(cmd) - 1

	switch {
	case argsNum > t.MaxArgs:
		return recipe.TokenInfo{}, nil, fmt.Errorf("Property \"%s\" have too many arguments (maximum is %d)", t.Keyword, t.MaxArgs)
	case argsNum < t.MinArgs:
		return recipe.TokenInfo{}, nil, fmt.Errorf("Property \"%s\" have too few arguments (minimum is %d)", t.Keyword, t.MinArgs)
	}

	return t, cmd[1:], nil
}

// appendData append data to recipe struct
func appendData(r *recipe.Recipe, t recipe.TokenInfo, args []string) error {
	if t.Global {
		switch t.Keyword {
		case "dir":
			r.Dir = args[0]
		case "unsafe-paths":
			if args[0] == "true" {
				r.UnsafePaths = true
			} else {
				return fmt.Errorf("Unsupported token value \"%s\"", args[0])
			}

			r.UnsafePaths = args[0] == "true"
		case "command":
			r.AddCommand(recipe.NewCommand(args))
		}

		return nil
	}

	lastCommand := r.Commands[len(r.Commands)-1]
	lastCommand.AddAction(recipe.NewAction(t.Keyword, args))

	return nil
}

// getTokenInfo return token info by keyword
func getTokenInfo(keyword string) recipe.TokenInfo {
	for _, token := range recipe.Tokens {
		if token.Keyword == keyword {
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
