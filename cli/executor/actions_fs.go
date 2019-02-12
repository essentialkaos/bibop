package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"pkg.re/essentialkaos/ek.v10/fsutil"
	"pkg.re/essentialkaos/ek.v10/hash"
	"pkg.re/essentialkaos/ek.v10/system"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionChdir is action processor for "chdir"
func actionChdir(action *recipe.Action) error {
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

	switch {
	case !action.Negative && perms != filePermsStr:
		return fmt.Errorf("File %s has invalid permissions (%s ≠ %s)", file, filePermsStr, perms)
	case action.Negative && perms == filePermsStr:
		return fmt.Errorf("File %s has invalid permissions (%s)", file, filePermsStr)
	}

	return nil
}

// actionOwner is action processor for "owner"
func actionOwner(action *recipe.Action) error {
	target, err := action.GetS(0)

	if err != nil {
		return err
	}

	owner, err := action.GetS(1)

	if err != nil {
		return err
	}

	uid, _, err := fsutil.GetOwner(target)

	if err != nil {
		return err
	}

	user, err := system.LookupUser(strconv.Itoa(uid))

	if err != nil {
		return err
	}

	switch {
	case !action.Negative && user.Name != owner:
		return fmt.Errorf("Object %s has invalid owner (%s ≠ %s)", target, user.Name, owner)
	case action.Negative && user.Name == owner:
		return fmt.Errorf("Object %s has invalid owner (%s)", target, user.Name)
	}

	return nil
}

// actionExist is action processor for "exist"
func actionExist(action *recipe.Action) error {
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

// actionReadable is action processor for "readable"
func actionReadable(action *recipe.Action) error {
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

// actionWritable is action processor for "writable"
func actionWritable(action *recipe.Action) error {
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

// actionExecutable is action processor for "executable"
func actionExecutable(action *recipe.Action) error {
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

// actionDirectory is action processor for "directory"
func actionDirectory(action *recipe.Action) error {
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

// actionEmptyDirectory is action processor for "empty-directory"
func actionEmptyDirectory(action *recipe.Action) error {
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

// actionEmpty is action processor for "empty"
func actionEmpty(action *recipe.Action) error {
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

	switch {
	case !action.Negative && fileHash != mustHash:
		return fmt.Errorf("File %s has invalid checksum hash (%s ≠ %s)", file, fileHash, mustHash)
	case action.Negative && fileHash == mustHash:
		return fmt.Errorf("File %s has invalid checksum hash (%s)", file, fileHash)
	}

	return nil
}

// actionChecksumRead is action processor for "actionChecksumRead"
func actionChecksumRead(action *recipe.Action) error {
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

// actionFileContains is action processor for "checksum"
func actionFileContains(action *recipe.Action) error {
	file, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, file)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Path \"%s\" is unsafe", file)
	}

	substr, err := action.GetS(0)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	switch {
	case !action.Negative && !bytes.Contains(data, []byte(substr)):
		return fmt.Errorf("File %s doesn't contain substring \"%s\"", file, substr)
	case action.Negative && bytes.Contains(data, []byte(substr)):
		return fmt.Errorf("File %s contains substring \"%s\"", file, substr)
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

	isSafePath, err := checkPathSafety(action.Command.Recipe, source)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Source has unsafe path (%s)", source)
	}

	isSafePath, err = checkPathSafety(action.Command.Recipe, dest)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Dest has unsafe path (%s)", dest)
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

	isSafePath, err := checkPathSafety(action.Command.Recipe, source)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Source has unsafe path (%s)", source)
	}

	isSafePath, err = checkPathSafety(action.Command.Recipe, dest)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Dest has unsafe path (%s)", dest)
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

	isSafePath, err := checkPathSafety(action.Command.Recipe, file)

	if err != nil {
		return err
	}

	if !isSafePath {
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

	isSafePath, err := checkPathSafety(action.Command.Recipe, dir)

	if err != nil {
		return err
	}

	if !isSafePath {
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
	target, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, target)

	if err != nil {
		return err
	}

	if !isSafePath {
		return fmt.Errorf("Path \"%s\" is unsafe", target)
	}

	err = os.RemoveAll(target)

	if err != nil {
		return err
	}

	return nil
}

// actionChmod is action processor for "chmod"
func actionChmod(action *recipe.Action) error {
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
		return fmt.Errorf("Path \"%s\" is unsafe", target)
	}

	err = os.Chmod(target, os.FileMode(mode))

	if err != nil {
		return err
	}

	return nil
}
