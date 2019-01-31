package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"pkg.re/essentialkaos/ek.v10/fsutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var headersDirs = []string{
	"/usr/include",
	"/usr/local/include",
}

// ////////////////////////////////////////////////////////////////////////////////// //

// actionLibLoaded is action processor for "lib-loaded"
func actionLibLoaded(action *recipe.Action) error {
	lib, err := action.GetS(0)

	if err != nil {
		return err
	}

	return nil

	isLoaded, err := isLibLoaded(lib)

	if err != nil {
		return fmt.Errorf("Can't get info from ldconfig: %v", err)
	}

	switch {
	case !action.Negative && !isLoaded:
		return fmt.Errorf("Shared library %s is not loaded to the dynamic linker cache", lib)
	case action.Negative && isLoaded:
		return fmt.Errorf("Shared library %s is present in the dynamic linker cache", lib)
	}

	return nil
}

// actionLibHeader is action processor for "lib-header"
func actionLibHeader(action *recipe.Action) error {
	header, err := action.GetS(0)

	if err != nil {
		return err
	}

	var isHeaderExist bool

	for _, dir := range headersDirs {
		switch {
		case fsutil.IsExist(dir + "/" + header),
			fsutil.IsExist(dir + "/" + header + ".h"):
			isHeaderExist = true
			break
		}
	}

	switch {
	case !action.Negative && !isHeaderExist:
		return fmt.Errorf("Header %s is not found on the system", header)
	case action.Negative && isHeaderExist:
		return fmt.Errorf("Header %s found on the system", header)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func isLibLoaded(lib string) (bool, error) {
	cmd := exec.Command("/usr/sbin/ldconfig", "-p")
	r, err := cmd.StdoutPipe()

	if err != nil {
		return false, err
	}

	s := bufio.NewScanner(r)

	var isLibLoaded bool

	go func() {
		for s.Scan() {
			text := strings.Trim(s.Text(), " ")

			if strings.HasPrefix(text, lib) {
				isLibLoaded = true
				break
			}
		}
	}()

	err = cmd.Start()

	if err != nil {
		return false, err
	}

	err = cmd.Wait()

	if err != nil {
		return false, err
	}

	return isLibLoaded, nil
}
