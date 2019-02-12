package executor

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"

	"pkg.re/essentialkaos/ek.v10/initsystem"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// actionServicePresent is action processor for "service-present"
func actionServicePresent(action *recipe.Action) error {
	service, err := action.GetS(0)

	if err != nil {
		return err
	}

	isServicePresent := initsystem.IsPresent(service)

	switch {
	case !action.Negative && !isServicePresent:
		return fmt.Errorf("Service %s doesn't present on the system", service)
	case action.Negative && isServicePresent:
		return fmt.Errorf("Service %s present on the system", service)
	}

	return nil
}

// actionServiceEnabled is action processor for "service-enabled"
func actionServiceEnabled(action *recipe.Action) error {
	service, err := action.GetS(0)

	if err != nil {
		return err
	}

	isServiceEnabled, err := initsystem.IsEnabled(service)

	if err != nil {
		return fmt.Errorf("Can't check auto start status for service %s: %v", service, err)
	}

	switch {
	case !action.Negative && !isServiceEnabled:
		return fmt.Errorf("Service %s auto start is not enabled", service)
	case action.Negative && isServiceEnabled:
		return fmt.Errorf("Service %s auto start is enabled", service)
	}

	return nil
}

// actionServiceWorks is action processor for "service-works"
func actionServiceWorks(action *recipe.Action) error {
	service, err := action.GetS(0)

	if err != nil {
		return err
	}

	isServiceWorks, err := initsystem.IsWorks(service)

	if err != nil {
		return fmt.Errorf("Can't check status for service %s: %v", service, err)
	}

	switch {
	case !action.Negative && !isServiceWorks:
		return fmt.Errorf("Service %s is not working", service)
	case action.Negative && isServiceWorks:
		return fmt.Errorf("Service %s is working", service)
	}

	return nil
}
