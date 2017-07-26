package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
func actionWait(action *recipe.Action) error {
	durSec, err := action.GetF(0)

	if err != nil {
		return err
	}

	durSec = mathutil.BetweenF64(durSec, 0.01, 3600.0)

	time.Sleep(secondsToDuration(durSec))

	return nil
}

// actionExpect is action processor for "expect"
func actionExpect(action *recipe.Action, output *outputStore) error {
	var (
		err     error
		start   time.Time
		substr  string
		maxWait float64
	)

	substr, err = action.GetS(0)

	if err != nil {
		return err
	}

	if len(action.Arguments) > 1 {
		maxWait, err = action.GetF(1)

		if err != nil {
			return err
		}
	} else {
		maxWait = 5.0
	}

	maxWait = mathutil.BetweenF64(maxWait, 0.01, 3600.0)
	start = time.Now()

	for {
		if bytes.Contains(output.data.Bytes(), []byte(substr)) {
			return fmt.Errorf("Output doesn't contains given substring")
		}

		if time.Since(start) >= secondsToDuration(maxWait) {
			return fmt.Errorf("Reached max wait time (%d sec)", maxWait)
		}

		time.Sleep(15 * time.Millisecond)
	}

	return nil
}

// actionInput is action processor for "input"
func actionInput(action *recipe.Action, input io.Writer) error {
	text, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.HasSuffix(text, "\n") {
		text = text + "\n"
	}

	_, err = input.Write([]byte(text))

	return err
}

// actionExit is action processor for "exit"
func actionExit(action *recipe.Action, cmd *exec.Cmd) error {
	if cmd == nil {
		return nil
	}

	var (
		err      error
		start    time.Time
		exitCode int
		maxWait  float64
	)

	go cmd.Wait()

	exitCode, err = action.GetI(0)

	if err != nil {
		return err
	}

	if len(action.Arguments) > 1 {
		maxWait, err = action.GetF(1)

		if err != nil {
			return err
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
			return fmt.Errorf("Reached max wait time (%d sec)", maxWait)
		}
	}

	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)

	if !ok {
		return fmt.Errorf("Can't get exit code from process state")
	}

	if status.ExitStatus() != exitCode {
		return fmt.Errorf("Process exit with unexpected exit code (%d ≠ %d)", status.ExitStatus(), exitCode)
	}

	return nil
}

// actionOutputEqual is action processor for "output-equal"
func actionOutputEqual(action *recipe.Action, output *outputStore) error {
	data, err := action.GetS(0)

	if err != nil {
		return err
	}

	if output.String() != data {
		return fmt.Errorf("Output doesn't equals substring \"%s\"", data)
	}

	return nil
}

// actionOutputContains is action processor for "output-contains"
func actionOutputContains(action *recipe.Action, output *outputStore) error {
	data, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.Contains(output.String(), data) {
		return fmt.Errorf("Output doesn't contains substring \"%s\"", data)
	}

	return nil
}

// actionOutputPrefix is action processor for "output-prefix"
func actionOutputPrefix(action *recipe.Action, output *outputStore) error {
	data, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.HasPrefix(output.String(), data) {
		return fmt.Errorf("Output doesn't have prefix \"%s\"", data)
	}

	return nil
}

// actionOutputSuffix is action processor for "output-suffix"
func actionOutputSuffix(action *recipe.Action, output *outputStore) error {
	data, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !strings.HasSuffix(output.String(), data) {
		return fmt.Errorf("Output doesn't have suffix \"%s\"", data)
	}

	return nil
}

// actionOutputLength is action processor for "output-length"
func actionOutputLength(action *recipe.Action, output *outputStore) error {
	size, err := action.GetI(0)

	if err != nil {
		return err
	}

	outputSize := len(output.String())

	if outputSize == size {
		return fmt.Errorf("Output have different length (%d ≠ %d)", outputSize, size)
	}

	return nil
}

// actionOutputTrim is action processor for "output-trim"
func actionOutputTrim(action *recipe.Action, output *outputStore) error {
	output.clear = true
	return nil
}

// actionPerms is action processor for "perms"
func actionPerms(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	perms, err := action.GetS(1)

	if err != nil {
		return err
	}

	filePerms := fsutil.GetPerms(file)
	filePermsStr := strconv.FormatUint(uint64(filePerms), 8)

	if perms != filePermsStr && "0"+perms != filePermsStr {
		return fmt.Errorf(
			"File %s have different permissions (%s ≠ %s)",
			file, filePermsStr, perms,
		)
	}

	return nil
}

// actionOwner is action processor for "owner"
func actionOwner(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	owner, err := action.GetS(1)

	if err != nil {
		return err
	}

	uid, _, err := fsutil.GetOwner(file)

	if err != nil {
		return err
	}

	user, err := system.LookupUser(strconv.Itoa(uid))

	if err != nil {
		return err
	}

	if user.Name != owner {
		return fmt.Errorf(
			"File %s have different owner (%s ≠ %s)",
			file, user.Name, owner,
		)
	}

	return nil
}

// actionExist is action processor for "exist"
func actionExist(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !fsutil.IsExist(file) {
		return fmt.Errorf("File %s doesn't exist", file)
	}

	return nil
}

// actionNotExist is action processor for "exist"
func actionNotExist(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if fsutil.IsExist(file) {
		return fmt.Errorf("File %s still exist", file)
	}

	return nil
}

// actionReadable is action processor for "readable"
func actionReadable(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !fsutil.IsReadable(file) {
		return fmt.Errorf("File %s is not readable", file)
	}

	return nil
}

// actionNotReadable is action processor for "not-readable"
func actionNotReadable(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if fsutil.IsReadable(file) {
		return fmt.Errorf("File %s is readable", file)
	}

	return nil
}

// actionWritable is action processor for "writable"
func actionWritable(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !fsutil.IsWritable(file) {
		return fmt.Errorf("File %s is not writable", file)
	}

	return nil
}

// actionNotWritable is action processor for "not-writable"
func actionNotWritable(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if fsutil.IsWritable(file) {
		return fmt.Errorf("File %s is writable", file)
	}

	return nil
}

// actionDirectory is action processor for "directory"
func actionDirectory(action *recipe.Action) error {
	dir, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !fsutil.IsDir(dir) {
		return fmt.Errorf("%s is not a directory", dir)
	}

	return nil
}

// actionNotDirectory is action processor for "not-directory"
func actionNotDirectory(action *recipe.Action) error {
	dir, err := action.GetS(0)

	if err != nil {
		return err
	}

	if fsutil.IsDir(dir) {
		return fmt.Errorf("%s is a directory", dir)
	}

	return nil
}

// actionEmptyDirectory is action processor for "empty-directory"
func actionEmptyDirectory(action *recipe.Action) error {
	dir, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !fsutil.IsEmptyDir(dir) {
		return fmt.Errorf("Directory %s is not empty", dir)
	}

	return nil
}

// actionNotEmptyDirectory is action processor for "not-empty-directory"
func actionNotEmptyDirectory(action *recipe.Action) error {
	dir, err := action.GetS(0)

	if err != nil {
		return err
	}

	if fsutil.IsEmptyDir(dir) {
		return fmt.Errorf("Directory %s is empty", dir)
	}

	return nil
}

// actionEmpty is action processor for "empty"
func actionEmpty(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !fsutil.IsNonEmpty(file) {
		return fmt.Errorf("File %s is empty", file)
	}

	return nil
}

// actionNotEmpty is action processor for "not-empty"
func actionNotEmpty(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if fsutil.IsNonEmpty(file) {
		return fmt.Errorf("File %s is not empty", file)
	}

	return nil
}

// actionChecksum is action processor for "checksum"
func actionChecksum(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	mustHash, err := action.GetS(0)

	if err != nil {
		return err
	}

	fileHash := hash.FileHash(file)

	if fileHash != mustHash {
		return fmt.Errorf(
			"File %s have different checksum hash (%s ≠ %s)",
			file, fileHash, mustHash,
		)
	}

	return nil
}

// actionFileContains is action processor for "checksum"
func actionFileContains(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !isSafePath(action.Command.Recipe, file) {
		return fmt.Errorf("Path \"%s\" is unsafe", file)
	}

	substr, err := action.GetS(0)

	if err != err {
		return err
	}

	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	if !bytes.Contains(data, []byte(substr)) {
		return fmt.Errorf("File %s doesn't contain substring \"%s\"", file, substr)
	}

	return nil
}

// actionCopy is action processor for "copy"
func actionCopy(action *recipe.Action) error {
	source, err := action.GetS(0)

	if err != nil {
		return err
	}

	dest, err := action.GetS(1)

	if err != nil {
		return err
	}

	if !isSafePath(action.Command.Recipe, source) {
		return fmt.Errorf("Source have unsafe path (%s)", source)
	}

	if !isSafePath(action.Command.Recipe, dest) {
		return fmt.Errorf("Dest have unsafe path (%s)", dest)
	}

	err = fsutil.CopyFile(source, dest)

	if err != nil {
		return err
	}

	return nil
}

// actionMove is action processor for "move"
func actionMove(action *recipe.Action) error {
	source, err := action.GetS(0)

	if err != nil {
		return err
	}

	dest, err := action.GetS(1)

	if err != nil {
		return err
	}

	if !isSafePath(action.Command.Recipe, source) {
		return fmt.Errorf("Source have unsafe path (%s)", source)
	}

	if !isSafePath(action.Command.Recipe, dest) {
		return fmt.Errorf("Dest have unsafe path (%s)", dest)
	}

	err = os.Rename(source, dest)

	if err != nil {
		return err
	}

	return nil
}

// actionTouch is action processor for "touch"
func actionTouch(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !isSafePath(action.Command.Recipe, file) {
		return fmt.Errorf("Path \"%s\" is unsafe", file)
	}

	err = ioutil.WriteFile(file, []byte(""), 0644)

	if err != nil {
		return err
	}

	return nil
}

// actionMkdir is action processor for "mkdir"
func actionMkdir(action *recipe.Action) error {
	dir, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !isSafePath(action.Command.Recipe, dir) {
		return fmt.Errorf("Path \"%s\" is unsafe", dir)
	}

	err = os.MkdirAll(dir, 0755)

	if err != nil {
		return err
	}

	return nil
}

// actionRemove is action processor for "remove"
func actionRemove(action *recipe.Action) error {
	path, err := action.GetS(0)

	if err != nil {
		return err
	}

	if !isSafePath(action.Command.Recipe, path) {
		return fmt.Errorf("Path \"%s\" is unsafe", path)
	}

	err = os.RemoveAll(path)

	if err != nil {
		return err
	}

	return nil
}

// actionChmod is action processor for "chmod"
func actionChmod(action *recipe.Action) error {
	path, err := action.GetS(0)

	if err != nil {
		return err
	}

	modeStr, err := action.GetS(1)

	if err != nil {
		return err
	}

	mode, err := strconv.ParseUint(modeStr, 8, 32)

	if err != nil {
		return err
	}

	if !isSafePath(action.Command.Recipe, path) {
		return fmt.Errorf("Path \"%s\" is unsafe", path)
	}

	err = os.Chmod(path, os.FileMode(mode))

	if err != nil {
		return err
	}

	return nil
}

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
