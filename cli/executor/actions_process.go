package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"strings"

	"pkg.re/essentialkaos/ek.v10/fsutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionProcessWorks is action processor for "process-works"
func actionProcessWorks(action *recipe.Action) error {
	pidFile, err := action.GetS(0)

	if err != nil {
		return err
	}

	pidFileData, err := ioutil.ReadFile(pidFile)

	if err != nil {
		return err
	}

	pid := strings.TrimRight(string(pidFileData), "\n\r")

	if !fsutil.IsExist("/proc/" + pid) {
		return fmt.Errorf(
			"Process with PID %s from PID file %s doesn't exist",
			pid, pidFile,
		)
	}

	return err
}
