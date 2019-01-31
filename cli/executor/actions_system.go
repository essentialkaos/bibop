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
	"net"
	"strings"
	"time"

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

	switch {
	case !action.Negative && !fsutil.IsExist("/proc/"+pid):
		return fmt.Errorf("Process with PID %s from PID file %s doesn't exist", pid, pidFile)
	case action.Negative && fsutil.IsExist("/proc/"+pid):
		return fmt.Errorf("Process with PID %s from PID file %s still exists", pid, pidFile)
	}

	return err
}

// actionConnect is action processor for "connect"
func actionConnect(action *recipe.Action) error {
	network, err := action.GetS(0)

	if err != nil {
		return err
	}

	address, err := action.GetS(1)

	if err != nil {
		return err
	}

	conn, err := net.DialTimeout(network, address, time.Second)

	if conn != nil {
		conn.Close()
	}

	switch {
	case !action.Negative && err != nil:
		return fmt.Errorf("Can't connect to %s (%s)", address, network)
	case action.Negative && err == nil:
		return fmt.Errorf("Successfully connected to %s (%s)", address, network)
	}

	return nil
}
