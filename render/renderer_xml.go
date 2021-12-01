package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// XMLRenderer is XML renderer
type XMLRenderer struct {
	start time.Time
	data  string
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Start prints info about started test
func (rr *XMLRenderer) Start(r *recipe.Recipe) {
	rr.start = time.Now()

	rr.data += "<report>\n"

	recipeFile, _ := filepath.Abs(r.File)
	workingDir, _ := filepath.Abs(r.Dir)

	rr.data += "  <recipe>\n"
	rr.data += fmt.Sprintf("    <recipe-file>%s</recipe-file>\n", rr.escapeData(recipeFile))
	rr.data += fmt.Sprintf("    <working-dir>%s</working-dir>\n", rr.escapeData(workingDir))
	rr.data += fmt.Sprintf("    <unsafe-actions>%t</unsafe-actions>\n", r.UnsafeActions)
	rr.data += fmt.Sprintf("    <require-root>%t</require-root>\n", r.RequireRoot)
	rr.data += fmt.Sprintf("    <fast-finish>%t</fast-finish>\n", r.FastFinish)
	rr.data += fmt.Sprintf("    <lock-workdir>%t</lock-workdir>\n", r.LockWorkdir)
	rr.data += fmt.Sprintf("    <unbuffer>%t</unbuffer>\n", r.Unbuffer)
	rr.data += "  </recipe>\n"
	rr.data += "  <commands>\n"
}

// CommandStarted prints info about started command
func (rr *XMLRenderer) CommandStarted(c *recipe.Command) {
	rr.data += "    <command"

	if c.User != "" {
		rr.data += fmt.Sprintf(" user=\"%s\"", c.User)
	}

	if c.Tag != "" {
		rr.data += fmt.Sprintf(" tag=\"%s\"", c.Tag)
	}

	rr.data += ">\n"

	rr.data += fmt.Sprintf("      <cmdline>%s</cmdline>\n", rr.escapeData(c.GetCmdline()))
	rr.data += fmt.Sprintf("      <description>%s</description>\n", rr.escapeData(c.Description))
	rr.data += "        <actions>\n"
}

// CommandSkipped prints info about skipped command
func (rr *XMLRenderer) CommandSkipped(c *recipe.Command) {
	return
}

// CommandFailed prints info about failed command
func (rr *XMLRenderer) CommandFailed(c *recipe.Command, err error) {
	rr.data += "        </actions>\n"
	rr.data += fmt.Sprintf("        <status failed=\"true\">%v</status>\n", err)
	rr.data += "    </command>\n"
}

// CommandFailed prints info about executed command
func (rr *XMLRenderer) CommandDone(c *recipe.Command, isLast bool) {
	rr.data += "        </actions>\n"
	rr.data += "        <status failed=\"false\"></status>\n"
	rr.data += "    </command>\n"
}

// ActionStarted prints info about action in progress
func (rr *XMLRenderer) ActionStarted(a *recipe.Action) {
	rr.data += "          <action>\n"
	rr.data += fmt.Sprintf("            <name>%s</name>\n", rr.formatActionName(a))
	rr.data += "            <arguments>\n"

	for index := range a.Arguments {
		arg, _ := a.GetS(index)
		rr.data += fmt.Sprintf(
			"              <argument>%s</argument>\n",
			rr.escapeData(arg),
		)
	}

	rr.data += "            </arguments>\n"
}

// ActionFailed prints info about failed action
func (rr *XMLRenderer) ActionFailed(a *recipe.Action, err error) {
	rr.data += fmt.Sprintf("            <status failed=\"true\">%v</status>\n", err)
	rr.data += "          </action>\n"
}

// ActionDone prints info about successfully finished action
func (rr *XMLRenderer) ActionDone(a *recipe.Action, isLast bool) {
	rr.data += "            <status failed=\"false\"></status>\n"
	rr.data += "          </action>\n"
}

// Result prints info about test results
func (rr *XMLRenderer) Result(passes, fails, skips int) {
	rr.data += "  </commands>\n"
	rr.data += fmt.Sprintf(
		"  <result passed=\"%d\" failed=\"%d\" skipped=\"%d\" duration=\"%g\" />\n",
		passes, fails, skips, time.Since(rr.start).Seconds(),
	)
	rr.data += "</report>"
	fmt.Println(rr.data)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (rr *XMLRenderer) escapeData(data string) string {
	data = strings.Replace(data, "<", "&lt;", -1)
	data = strings.Replace(data, ">", "&gt;", -1)
	data = strings.Replace(data, "&", "&amp;", -1)

	return data
}

// formatActionName format action name
func (rr *XMLRenderer) formatActionName(a *recipe.Action) string {
	if a.Negative {
		return rr.escapeData("!" + a.Name)
	}

	return rr.escapeData(a.Name)
}
