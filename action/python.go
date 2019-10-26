package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os/exec"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// PythonModule is action processor for "python-module"
func PythonModule(action *recipe.Action) error {
	module, err := action.GetS(0)

	if err != nil {
		return err
	}

	if checkPythonModuleLoad(2, module) != nil {
		return fmt.Errorf("Python2 module %s cannot be loaded", module)
	}

	return nil
}

// Python3Module is action processor for "python3-module"
func Python3Module(action *recipe.Action) error {
	module, err := action.GetS(0)

	if err != nil {
		return err
	}

	if checkPythonModuleLoad(3, module) != nil {
		return fmt.Errorf("Python3 module %s cannot be loaded", module)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// checkPythonModuleLoad returns true if module can be loaded
func checkPythonModuleLoad(pythonVersion int, moduleName string) error {
	pythonBinary := "python"

	if pythonVersion == 3 {
		pythonBinary = "python3"
	}

	return exec.Command(pythonBinary, "-c", fmt.Sprintf("import %s", moduleName)).Run()
}
