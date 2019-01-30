package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"

	"pkg.re/essentialkaos/ek.v10/system"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionUserExist is action processor for "user-exist"
func actionUserExist(action *recipe.Action) error {
	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	user, _ := system.LookupUser(username)

	switch {
	case !action.Negative && user == nil:
		return fmt.Errorf("User %s doesn't exist on the system", username)
	case action.Negative && user != nil:
		return fmt.Errorf("User %s exists on the system", username)
	}

	return nil
}

// actionUserID action processor for "user-id"
func actionUserID(action *recipe.Action) error {
	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	uid, err := action.GetI(1)

	if err != nil {
		return err
	}

	if uid < 0 {
		return fmt.Errorf("UID can't be less than 0 (%d)", uid)
	}

	user, _ := system.LookupUser(username)

	if user == nil {
		return fmt.Errorf("User %s doesn't exist on the system", username)
	}

	switch {
	case !action.Negative && user.UID != uid:
		return fmt.Errorf("User %s has different UID (%d ≠ %d)", username, user.UID, uid)
	case action.Negative && user.UID == uid:
		return fmt.Errorf("User %s has invalid UID (%d)", username, uid)
	}

	return nil
}

// actionUserGID action processor for "user-gid"
func actionUserGID(action *recipe.Action) error {
	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	gid, err := action.GetI(1)

	if err != nil {
		return err
	}

	if gid < 0 {
		return fmt.Errorf("GID can't be less than 0 (%d)", gid)
	}

	user, _ := system.LookupUser(username)

	if user == nil {
		return fmt.Errorf("User %s doesn't exist on the system", username)
	}

	switch {
	case !action.Negative && user.GID != gid:
		return fmt.Errorf("User %s has different GID (%d ≠ %d)", username, user.GID, gid)
	case action.Negative && user.GID == gid:
		return fmt.Errorf("User %s has invalid GID (%d)", username, gid)
	}

	return nil
}

// actionUserGroup is action processor for "user-group"
func actionUserGroup(action *recipe.Action) error {
	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	groupname, err := action.GetS(1)

	if err != nil {
		return err
	}

	user, _ := system.LookupUser(username)

	if user == nil {
		return fmt.Errorf("User %s doesn't exist on the system", username)
	}

	group, _ := system.LookupGroup(groupname)

	if group == nil {
		return fmt.Errorf("Group %s doesn't exist on the system", groupname)
	}

	var hasGroup bool

	for _, group := range user.Groups {
		if group.Name == groupname {
			hasGroup = true
		}
	}

	switch {
	case !action.Negative && !hasGroup:
		return fmt.Errorf("User %s is not a part of group %s", username, groupname)
	case action.Negative && hasGroup:
		return fmt.Errorf("User %s is a part of group %s", username, groupname)
	}

	return nil
}

// actionUserShell is action processor for "user-shell"
func actionUserShell(action *recipe.Action) error {
	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	shell, err := action.GetS(1)

	if err != nil {
		return err
	}

	user, _ := system.LookupUser(username)

	if user == nil {
		return fmt.Errorf("User %s doesn't exist on the system", username)
	}

	switch {
	case !action.Negative && user.Shell != shell:
		return fmt.Errorf("User %s has different shell (%d ≠ %d)", username, user.Shell, shell)
	case action.Negative && user.Shell == shell:
		return fmt.Errorf("User %s has invalid shell (%d)", username, shell)
	}

	return nil
}

// actionUserHome is action processor for "user-home"
func actionUserHome(action *recipe.Action) error {
	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	homeDir, err := action.GetS(1)

	if err != nil {
		return err
	}

	user, _ := system.LookupUser(username)

	if user == nil {
		return fmt.Errorf("User %s doesn't exist on the system", username)
	}

	switch {
	case !action.Negative && user.HomeDir != homeDir:
		return fmt.Errorf("User %s has different home directory (%d ≠ %d)", username, user.HomeDir, homeDir)
	case action.Negative && user.HomeDir == homeDir:
		return fmt.Errorf("User %s has invalid home directory (%d)", username, homeDir)
	}

	return nil
}

// actionGroupExist is action processor for "group-exist"
func actionGroupExist(action *recipe.Action) error {
	groupname, err := action.GetS(0)

	if err != nil {
		return err
	}

	group, _ := system.LookupGroup(groupname)

	switch {
	case !action.Negative && group == nil:
		return fmt.Errorf("Group %s doesn't exist on the system", groupname)
	case action.Negative && group != nil:
		return fmt.Errorf("Group %s exists on the system", groupname)
	}

	return nil
}

// actionGroupID is action processor for "group-id"
func actionGroupID(action *recipe.Action) error {
	groupname, err := action.GetS(0)

	if err != nil {
		return err
	}

	gid, err := action.GetI(1)

	if err != nil {
		return err
	}

	group, _ := system.LookupGroup(groupname)

	if group == nil {
		return fmt.Errorf("Group %s doesn't exist on the system", groupname)
	}

	switch {
	case !action.Negative && group.GID != gid:
		return fmt.Errorf("Group %s has different GID (%d ≠ %d)", groupname, group.GID, gid)
	case action.Negative && group.GID == gid:
		return fmt.Errorf("Group %s has invalid GID (%d)", groupname, gid)
	}

	return nil
}
