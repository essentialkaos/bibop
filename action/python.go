package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os/exec"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Python2Package is action processor for "python2-package"
func Python2Package(action *recipe.Action) error {
	pkg, err := action.GetS(0)

	if err != nil {
		return err
	}

	if checkPythonPackageLoad(2, pkg) != nil {
		return fmt.Errorf("Python2 package %s cannot be loaded", pkg)
	}

	return nil
}

// Python3Package is action processor for "python3-package"
func Python3Package(action *recipe.Action) error {
	pkg, err := action.GetS(0)

	if err != nil {
		return err
	}

	if checkPythonPackageLoad(3, pkg) != nil {
		return fmt.Errorf("Python3 package %s cannot be loaded", pkg)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// checkPythonPackageLoad returns true if package can be loaded
func checkPythonPackageLoad(pythonVersion int, moduleName string) error {
	pythonBinary := "python"

	if pythonVersion == 3 {
		pythonBinary = "python3"
	}

	return exec.Command(pythonBinary, "-c", "import "+moduleName).Run()
}
