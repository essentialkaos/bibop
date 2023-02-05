package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"

	"github.com/essentialkaos/ek/v12/initsystem"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// ServicePresent is action processor for "service-present"
func ServicePresent(action *recipe.Action) error {
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

// ServiceEnabled is action processor for "service-enabled"
func ServiceEnabled(action *recipe.Action) error {
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

// ServiceWorks is action processor for "service-works"
func ServiceWorks(action *recipe.Action) error {
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
