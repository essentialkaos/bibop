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

	if err != nil {
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
