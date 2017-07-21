package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"io"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/hash"
	"pkg.re/essentialkaos/ek.v9/mathutil"
	"pkg.re/essentialkaos/ek.v9/system"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionWait is action processor for "exit"
func actionWait(action *recipe.Action) bool {
	durSec, err := action.GetF(0)

	if err != nil {
		return false
	}

	durSec = mathutil.BetweenF64(durSec, 0.01, 3600.0)

	time.Sleep(secondsToDuration(durSec))

	return true
}

// actionExpect is action processor for "expect"
func actionExpect(action *recipe.Action, output *outputStore) bool {
	var (
		err     error
		start   time.Time
		substr  string
		maxWait float64
	)

	substr, err = action.GetS(0)

	if err != nil {
		return false
	}

	if len(action.Arguments) > 1 {
		maxWait, err = action.GetF(1)

		if err != nil {
			return false
		}
	} else {
		maxWait = 5.0
	}

	maxWait = mathutil.BetweenF64(maxWait, 0.01, 3600.0)
	start = time.Now()

	for {
		if bytes.Contains(output.data.Bytes(), []byte(substr)) {
			return true
		}

		if time.Since(start) >= secondsToDuration(maxWait) {
			return false
		}

		time.Sleep(15 * time.Millisecond)
	}

	return false
}

// actionInput is action processor for "input"
func actionInput(action *recipe.Action, input io.Writer) bool {
	text, err := action.GetS(0)

	if err != nil {
		return false
	}

	if !strings.HasSuffix(text, "\n") {
		text = text + "\n"
	}

	_, err = input.Write([]byte(text))

	return err == nil
}

// actionExit is action processor for "exit"
func actionExit(action *recipe.Action, cmd *exec.Cmd) bool {
	var (
		err      error
		start    time.Time
		exitCode int
		maxWait  float64
	)

	go cmd.Wait()

	exitCode, err = action.GetI(0)

	if err != nil {
		return false
	}

	if len(action.Arguments) > 1 {
		maxWait, err = action.GetF(1)

		if err != nil {
			return false
		}
	} else {
		maxWait = 60.0
	}

	start = time.Now()

	for {
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			break
		}

		if time.Since(start) > secondsToDuration(maxWait) {
			return false
		}
	}

	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)

	if !ok {
		return false
	}

	return status.ExitStatus() == exitCode
}

// actionOutputEqual is action processor for "output-equal"
func actionOutputEqual(action *recipe.Action, output *outputStore) bool {
	data, err := action.GetS(0)

	if err != nil {
		return false
	}

	return output.String() == data
}

// actionOutputContains is action processor for "output-contains"
func actionOutputContains(action *recipe.Action, output *outputStore) bool {
	data, err := action.GetS(0)

	if err != nil {
		return false
	}

	return strings.Contains(output.String(), data)
}

// actionOutputPrefix is action processor for "output-prefix"
func actionOutputPrefix(action *recipe.Action, output *outputStore) bool {
	data, err := action.GetS(0)

	if err != nil {
		return false
	}

	return strings.HasPrefix(output.String(), data)
}

// actionOutputSuffix is action processor for "output-suffix"
func actionOutputSuffix(action *recipe.Action, output *outputStore) bool {
	data, err := action.GetS(0)

	if err != nil {
		return false
	}

	return strings.HasSuffix(output.String(), data)
}

// actionOutputLength is action processor for "output-length"
func actionOutputLength(action *recipe.Action, output *outputStore) bool {
	size, err := action.GetI(0)

	if err != nil {
		return false
	}

	return len(output.String()) == size
}

// actionOutputTrim is action processor for "output-trim"
func actionOutputTrim(action *recipe.Action, output *outputStore) bool {
	output.clear = true
	return true
}

// actionPerms is action processor for "perms"
func actionPerms(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	perms, err := action.GetS(1)

	if err != nil {
		return false
	}

	filePerms := fsutil.GetPerms(file)
	filePermsStr := strconv.FormatUint(uint64(filePerms), 8)

	return perms == filePermsStr
}

// actionOwner is action processor for "owner"
func actionOwner(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	owner, err := action.GetS(1)

	if err != nil {
		return false
	}

	uid, _, err := fsutil.GetOwner(file)

	if err != nil {
		return false
	}

	user, err := system.LookupUser(strconv.Itoa(uid))

	if err != nil {
		return false
	}

	return user.Name == owner
}

// actionExist is action processor for "exist"
func actionExist(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	return fsutil.IsExist(file)
}

// actionNotExist is action processor for "exist"
func actionNotExist(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	return !fsutil.IsExist(file)
}

// actionReadable is action processor for "readable"
func actionReadable(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	return fsutil.IsReadable(file)
}

// actionNotReadable is action processor for "not-readable"
func actionNotReadable(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	return !fsutil.IsReadable(file)
}

// actionWritable is action processor for "writable"
func actionWritable(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	return fsutil.IsWritable(file)
}

// actionNotWritable is action processor for "not-writable"
func actionNotWritable(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	return !fsutil.IsWritable(file)
}

// actionDirectory is action processor for "directory"
func actionDirectory(action *recipe.Action) bool {
	dir, err := action.GetS(0)

	if err != nil {
		return false
	}

	return fsutil.IsDir(dir)
}

// actionNotDirectory is action processor for "not-directory"
func actionNotDirectory(action *recipe.Action) bool {
	dir, err := action.GetS(0)

	if err != nil {
		return false
	}

	return !fsutil.IsDir(dir)
}

// actionEmptyDirectory is action processor for "empty-directory"
func actionEmptyDirectory(action *recipe.Action) bool {
	dir, err := action.GetS(0)

	if err != nil {
		return false
	}

	return fsutil.IsEmptyDir(dir)
}

// actionNotEmptyDirectory is action processor for "not-empty-directory"
func actionNotEmptyDirectory(action *recipe.Action) bool {
	dir, err := action.GetS(0)

	if err != nil {
		return false
	}

	return !fsutil.IsEmptyDir(dir)
}

// actionEmpty is action processor for "empty"
func actionEmpty(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	return !fsutil.IsNonEmpty(file)
}

// actionNotEmpty is action processor for "not-empty"
func actionNotEmpty(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	return fsutil.IsNonEmpty(file)
}

// actionChecksum is action processor for "checksum"
func actionChecksum(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	mustHash, err := action.GetS(0)

	if err != nil {
		return false
	}

	fileHash := hash.FileHash(file)

	return fileHash == mustHash
}

// actionFileContains is action processor for "checksum"
func actionFileContains(action *recipe.Action) bool {
	file, err := action.GetS(0)

	if err != nil {
		return false
	}

	substr, err := action.GetS(0)

	if err != nil {
		return false
	}

	data, err := ioutil.ReadFile(file)

	if err != nil {
		return false
	}

	return bytes.Contains(data, []byte(substr))
}

// actionProcessWorks is action processor for "process-works"
func actionProcessWorks(action *recipe.Action) bool {
	pidFile, err := action.GetS(0)

	if err != nil {
		return false
	}

	pidFileData, err := ioutil.ReadFile(pidFile)

	if err != nil {
		return false
	}

	pid := strings.TrimRight(string(pidFileData), "\n\r")

	return fsutil.IsExist("/proc/" + pid)
}
