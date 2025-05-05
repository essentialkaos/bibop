package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"hash/crc32"
	"os"

	"github.com/essentialkaos/ek/v13/fsutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Backup is action processor for "backup"
func Backup(action *recipe.Action, tmpDir string) error {
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

	pathCRC32 := calcCRC32Q(path)

	err = fsutil.CopyFile(path, tmpDir+"/"+pathCRC32)

	if err != nil {
		return fmt.Errorf("Can't backup file: %v", err)
	}

	err = fsutil.CopyAttr(path, tmpDir+"/"+pathCRC32)

	if err != nil {
		return fmt.Errorf("Can't copy attributes: %v", err)
	}

	return nil
}

// BackupRestore is action processor for "backup-restore"
func BackupRestore(action *recipe.Action, tmpDir string) error {
	path, err := action.GetS(0)

	if err != nil {
		return err
	}

	isSafePath, err := checkPathSafety(action.Command.Recipe, path)

	if err != nil {
		return err
	}

	pathCRC32 := calcCRC32Q(path)
	backupFile := tmpDir + "/" + pathCRC32

	switch {
	case !isSafePath:
		return fmt.Errorf("Path is unsafe (%s)", path)
	case !fsutil.IsExist(backupFile):
		return fmt.Errorf("Backup file for %s does not exist", path)
	}

	err = os.Remove(path)

	if err != nil {
		return fmt.Errorf("Can't remove original file: %v", err)
	}

	err = fsutil.CopyFile(backupFile, path)

	if err != nil {
		return fmt.Errorf("Can't copy backup file: %v", err)
	}

	err = fsutil.CopyAttr(backupFile, path)

	if err != nil {
		return fmt.Errorf("Can't copy attributes: %v", err)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// calcCRC32Q calculate CRC32 checksum
func calcCRC32Q(data string) string {
	table := crc32.MakeTable(0xD5828281)
	return fmt.Sprintf("%08x", crc32.Checksum([]byte(data), table))
}
