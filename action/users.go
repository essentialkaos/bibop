package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"

	"pkg.re/essentialkaos/ek.v12/system"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// UserExist is action processor for "user-exist"
func UserExist(action *recipe.Action) error {
	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	isUserExist := system.IsUserExist(username)

	switch {
	case !action.Negative && !isUserExist:
		return fmt.Errorf("User %s doesn't exist on the system", username)
	case action.Negative && isUserExist:
		return fmt.Errorf("User %s exists on the system", username)
	}

	return nil
}

// UserID action processor for "user-id"
func UserID(action *recipe.Action) error {
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

// UserGID action processor for "user-gid"
func UserGID(action *recipe.Action) error {
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

// UserGroup is action processor for "user-group"
func UserGroup(action *recipe.Action) error {
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

	if !system.IsGroupExist(groupname) {
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

// UserShell is action processor for "user-shell"
func UserShell(action *recipe.Action) error {
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
		return fmt.Errorf("User %s has different shell (%s ≠ %s)", username, user.Shell, shell)
	case action.Negative && user.Shell == shell:
		return fmt.Errorf("User %s has invalid shell (%s)", username, shell)
	}

	return nil
}

// UserHome is action processor for "user-home"
func UserHome(action *recipe.Action) error {
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
		return fmt.Errorf("User %s has different home directory (%s ≠ %s)", username, user.HomeDir, homeDir)
	case action.Negative && user.HomeDir == homeDir:
		return fmt.Errorf("User %s has invalid home directory (%s)", username, homeDir)
	}

	return nil
}

// GroupExist is action processor for "group-exist"
func GroupExist(action *recipe.Action) error {
	groupname, err := action.GetS(0)

	if err != nil {
		return err
	}

	isGroupExist := system.IsGroupExist(groupname)

	switch {
	case !action.Negative && !isGroupExist:
		return fmt.Errorf("Group %s doesn't exist on the system", groupname)
	case action.Negative && isGroupExist:
		return fmt.Errorf("Group %s exists on the system", groupname)
	}

	return nil
}

// GroupID is action processor for "group-id"
func GroupID(action *recipe.Action) error {
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
