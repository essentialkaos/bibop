package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"

	"pkg.re/essentialkaos/ek.v10/fsutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionBackup is action processor for "backup"
func actionBackup(action *recipe.Action) error {
	path, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, path)

	if err != nil {
		return err
	}

	switch {
	case !isSafePath:
		return fmt.Errorf("Path is unsafe (%s)", path)
	case !fsutil.IsExist(path):
		return fmt.Errorf("File %s does not exist", path)
	case !fsutil.IsRegular(path):
		return fmt.Errorf("Object %s is not a file", path)
	}

	err = fsutil.CopyFile(path, path+".bak")

	if err != nil {
		return fmt.Errorf("Can't create backup file: %v", err)
	}

	return nil
}

// actionBackupRestore is action processor for "backup-restore"
func actionBackupRestore(action *recipe.Action) error {
	path, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, path)

	if err != nil {
		return err
	}

	backupFile := path + ".bak"

	switch {
	case !isSafePath:
		return fmt.Errorf("Path is unsafe (%s)", path)
	case !fsutil.IsExist(backupFile):
		return fmt.Errorf("Backup file %s does not exist", backupFile)
	case !fsutil.IsRegular(backupFile):
		return fmt.Errorf("Object %s is not a file", backupFile)
	}

	err = os.Remove(path)

	if err != nil {
		return fmt.Errorf("Can't remove original file: %v", err)
	}

	err = fsutil.MoveFile(backupFile, path)

	if err != nil {
		return fmt.Errorf("Can't move backup file: %v", err)
	}

	return nil
}
