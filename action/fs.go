package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/hash"
	"github.com/essentialkaos/ek/v12/strutil"
	"github.com/essentialkaos/ek/v12/system"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Chdir is action processor for "chdir"
func Chdir(action *recipe.Action) error {
	path, err := action.GetS(0)

	if err != nil {
		return err
	}

	err = os.Chdir(path)

	if err != nil {
		return fmt.Errorf("Can't change current directory to %s: %v", path, err)
	}

	return nil
}

// Mode is action processor for "mode"
func Mode(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	mode, err := action.GetS(1)

	if err != nil {
		return err
	}

	fileMode := fsutil.GetMode(file)
	fileModeStr := strconv.FormatUint(uint64(fileMode), 8)

	switch {
	case !action.Negative && mode != fileModeStr:
		return fmt.Errorf("File %s has invalid mode (%s ≠ %s)", file, fileModeStr, mode)
	case action.Negative && mode == fileModeStr:
		return fmt.Errorf("File %s has invalid mode (%s)", file, fileModeStr)
	}

	return nil
}

// Owner is action processor for "owner"
func Owner(action *recipe.Action) error {
	target, err := action.GetS(0)

	if err != nil {
		return err
	}

	userAndGroup, err := action.GetS(1)

	if err != nil {
		return err
	}

	userName := strutil.ReadField(userAndGroup, 0, false, ":")
	groupName := strutil.ReadField(userAndGroup, 1, false, ":")

	uid, gid, err := fsutil.GetOwner(target)

	if err != nil {
		return err
	}

	user, err := system.LookupUser(strconv.Itoa(uid))

	if err != nil {
		return err
	}

	group, err := system.LookupGroup(strconv.Itoa(gid))

	if err != nil {
		return err
	}

	switch {
	case !action.Negative && user.Name != userName:
		return fmt.Errorf("Object %s has invalid owner (%s ≠ %s)", target, user.Name, userName)
	case action.Negative && user.Name == userName:
		return fmt.Errorf("Object %s has invalid owner (%s)", target, user.Name)
	case groupName != "" && !action.Negative && group.Name != groupName:
		return fmt.Errorf("Object %s has invalid owner group (%s ≠ %s)", target, group.Name, groupName)
	case groupName != "" && action.Negative && group.Name == groupName:
		return fmt.Errorf("Object %s has invalid owner group (%s)", target, group.Name)
	}

	return nil
}

// Exist is action processor for "exist"
func Exist(action *recipe.Action) error {
	target, err := action.GetS(0)

	if err != nil {
		return err
	}

	switch {
	case !action.Negative && !fsutil.IsExist(target):
		return fmt.Errorf("Object %s doesn't exist", target)
	case action.Negative && fsutil.IsExist(target):
		return fmt.Errorf("Object %s exists, but it mustn't", target)
	}

	return nil
}

// Link is action processor for "link"
func Link(action *recipe.Action) error {
	link, err := action.GetS(0)

	if err != nil {
		return err
	}

	target, err := action.GetS(1)

	if err != nil {
		return err
	}

	linkTarget, err := filepath.EvalSymlinks(link)

	if err != nil {
		return err
	}

	if filepath.IsAbs(linkTarget) {
		linkTarget, err = filepath.Abs(linkTarget)

		if err != nil {
			return err
		}
	}

	switch {
	case !action.Negative && target != linkTarget:
		return fmt.Errorf(
			"Link %s points on different object (%s ≠ %s)",
			link, target, linkTarget,
		)
	case action.Negative && target == linkTarget:
		return fmt.Errorf(
			"Link %s points on defined object (%s)",
			link, target,
		)
	}

	return nil
}

// Readable is action processor for "readable"
func Readable(action *recipe.Action) error {
	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	target, err := action.GetS(1)

	if err != nil {
		return err
	}

	if !system.IsUserExist(username) {
		return fmt.Errorf("User %s doesn't exist on the system", username)
	}

	switch {
	case !action.Negative && !fsutil.IsReadableByUser(target, username):
		return fmt.Errorf("Object %s is not readable for user %s", target, username)
	case action.Negative && fsutil.IsReadableByUser(target, username):
		return fmt.Errorf("Object %s is readable for user %s, but it mustn't", target, username)
	}

	return nil
}

// Writable is action processor for "writable"
func Writable(action *recipe.Action) error {
	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	target, err := action.GetS(1)

	if err != nil {
		return err
	}

	if !system.IsUserExist(username) {
		return fmt.Errorf("User %s doesn't exist on the system", username)
	}

	switch {
	case !action.Negative && !fsutil.IsWritableByUser(target, username):
		return fmt.Errorf("Object %s is not writable for user %s", target, username)
	case action.Negative && fsutil.IsWritableByUser(target, username):
		return fmt.Errorf("Object %s is writable for user %s, but it mustn't", target, username)
	}

	return nil
}

// Executable is action processor for "executable"
func Executable(action *recipe.Action) error {
	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	target, err := action.GetS(1)

	if err != nil {
		return err
	}

	if !system.IsUserExist(username) {
		return fmt.Errorf("User %s doesn't exist on the system", username)
	}

	switch {
	case !action.Negative && !fsutil.IsExecutableByUser(target, username):
		return fmt.Errorf("Object %s is not executable for user %s", target, username)
	case action.Negative && fsutil.IsExecutableByUser(target, username):
		return fmt.Errorf("Object %s is executable for user %s, but it mustn't", target, username)
	}

	return nil
}

// Dir is action processor for "dir"
func Dir(action *recipe.Action) error {
	dir, err := action.GetS(0)

	if err != nil {
		return err
	}

	switch {
	case !action.Negative && !fsutil.IsDir(dir):
		return fmt.Errorf("%s is not a directory", dir)
	case action.Negative && fsutil.IsDir(dir):
		return fmt.Errorf("%s is a directory", dir)
	}

	return nil
}

// EmptyDir is action processor for "empty-dir"
func EmptyDir(action *recipe.Action) error {
	dir, err := action.GetS(0)

	if err != nil {
		return err
	}

	switch {
	case !action.Negative && !fsutil.IsEmptyDir(dir):
		return fmt.Errorf("Directory %s is not empty", dir)
	case action.Negative && fsutil.IsEmptyDir(dir):
		return fmt.Errorf("Directory %s is empty", dir)
	}

	return nil
}

// Empty is action processor for "empty"
func Empty(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	switch {
	case !action.Negative && fsutil.IsNonEmpty(file):
		return fmt.Errorf("File %s is not empty", file)
	case action.Negative && !fsutil.IsNonEmpty(file):
		return fmt.Errorf("File %s is empty", file)
	}

	return nil
}

// Checksum is action processor for "checksum"
func Checksum(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	mustHash, err := action.GetS(1)

	if err != nil {
		return err
	}

	fileHash := hash.FileHash(file)

	switch {
	case !action.Negative && fileHash != mustHash:
		return fmt.Errorf("File %s has invalid checksum hash (%s ≠ %s)", file, fileHash, mustHash)
	case action.Negative && fileHash == mustHash:
		return fmt.Errorf("File %s has invalid checksum hash (%s)", file, fileHash)
	}

	return nil
}

// ChecksumRead is action processor for "actionChecksumRead"
func ChecksumRead(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	variable, err := action.GetS(1)

	if err != nil {
		return err
	}

	hash := hash.FileHash(file)

	return action.Command.Recipe.SetVariable(variable, hash)
}

// FileContains is action processor for "checksum"
func FileContains(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, file)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Path %q is unsafe", file)
	}

	substr, err := action.GetS(1)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	switch {
	case !action.Negative && !bytes.Contains(data, []byte(substr)):
		return fmt.Errorf("File %s doesn't contain substring %q", file, substr)
	case action.Negative && bytes.Contains(data, []byte(substr)):
		return fmt.Errorf("File %s contains substring %q", file, substr)
	}

	return nil
}

// Copy is action processor for "copy"
func Copy(action *recipe.Action) error {
	source, err := action.GetS(0)

	if err != nil {
		return err
	}

	dest, err := action.GetS(1)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, source)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Action uses unsafe path (%s)", source)
	}

	isSafePath, err = checkPathSafety(action.Command.Recipe, dest)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Action uses unsafe path (%s)", dest)
	}

	err = fsutil.CopyFile(source, dest)

	if err != nil {
		return err
	}

	return nil
}

// Move is action processor for "move"
func Move(action *recipe.Action) error {
	source, err := action.GetS(0)

	if err != nil {
		return err
	}

	dest, err := action.GetS(1)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, source)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Action uses unsafe path (%s)", source)
	}

	isSafePath, err = checkPathSafety(action.Command.Recipe, dest)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Action uses unsafe path (%s)", dest)
	}

	err = os.Rename(source, dest)

	if err != nil {
		return err
	}

	return nil
}

// Touch is action processor for "touch"
func Touch(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, file)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Path %q is unsafe", file)
	}

	err = ioutil.WriteFile(file, []byte(""), 0644)

	if err != nil {
		return err
	}

	return nil
}

// Mkdir is action processor for "mkdir"
func Mkdir(action *recipe.Action) error {
	dir, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, dir)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Path %q is unsafe", dir)
	}

	err = os.MkdirAll(dir, 0755)

	if err != nil {
		return err
	}

	return nil
}

// Remove is action processor for "remove"
func Remove(action *recipe.Action) error {
	target, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, target)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Path %q is unsafe", target)
	}

	err = os.RemoveAll(target)

	if err != nil {
		return err
	}

	return nil
}

// Cleanup is action processor for "cleanup"
func Cleanup(action *recipe.Action) error {
	target, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, target)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Path %q is unsafe", target)
	}

	if !fsutil.IsDir(target) {
		return fmt.Errorf("Target object %q is not a directory", target)
	}

	if fsutil.IsEmptyDir(target) {
		return nil
	}

	targetObjects := fsutil.List(target, false)
	fsutil.ListToAbsolute(target, targetObjects)

	for _, obj := range targetObjects {
		err = os.RemoveAll(obj)

		if err != nil {
			return fmt.Errorf("Can't remove object %q: %v", err)
		}
	}

	return nil
}

// Chmod is action processor for "chmod"
func Chmod(action *recipe.Action) error {
	target, err := action.GetS(0)

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

	isSafePath, err := checkPathSafety(action.Command.Recipe, target)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Path %q is unsafe", target)
	}

	err = os.Chmod(target, os.FileMode(mode))

	if err != nil {
		return err
	}

	return nil
}

// Chown is action processor for "chown"
func Chown(action *recipe.Action) error {
	var user, group string
	var uid, gid int

	target, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, target)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Path %q is unsafe", target)
	}

	owner, err := action.GetS(1)

	if err != nil {
		return err
	}

	if strings.Contains(owner, ":") {
		user, group, _ = strings.Cut(owner, ":")
	} else {
		user = owner
	}

	if user != "" {
		ownerUser, err := system.LookupUser(user)

		if err != nil {
			return err
		}

		uid = ownerUser.UID
	} else {
		uid = -1
	}

	if group != "" {
		ownerGroup, err := system.LookupGroup(group)

		if err != nil {
			return err
		}

		gid = ownerGroup.GID
	} else {
		gid = -1
	}

	err = os.Chown(target, uid, gid)

	if err != nil {
		return err
	}

	return nil
}

// Truncate is action processor for "truncate"
func Truncate(action *recipe.Action) error {
	target, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, target)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Path %q is unsafe", target)
	}

	return os.Truncate(target, 0)
}
