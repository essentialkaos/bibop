package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// XMLRenderer is XML renderer
type XMLRenderer struct {
	Version string

	start time.Time
	data  strings.Builder
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Start prints info about started test
func (rr *XMLRenderer) Start(r *recipe.Recipe) {
	rr.start = time.Now()

	rr.data.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n")
	rr.data.WriteString(fmt.Sprintf("<!-- bibop %s | recipe report -->\n", rr.Version))
	rr.data.WriteString("<report>\n")

	recipeFile, _ := filepath.Abs(r.File)
	workingDir, _ := filepath.Abs(r.Dir)

	rr.data.WriteString("  <recipe>\n")
	rr.data.WriteString(fmt.Sprintf("    <recipe-file>%s</recipe-file>\n", rr.escapeData(recipeFile)))
	rr.data.WriteString(fmt.Sprintf("    <working-dir>%s</working-dir>\n", rr.escapeData(workingDir)))
	rr.data.WriteString(fmt.Sprintf("    <unsafe-actions>%t</unsafe-actions>\n", r.UnsafeActions))
	rr.data.WriteString(fmt.Sprintf("    <require-root>%t</require-root>\n", r.RequireRoot))
	rr.data.WriteString(fmt.Sprintf("    <fast-finish>%t</fast-finish>\n", r.FastFinish))
	rr.data.WriteString(fmt.Sprintf("    <lock-workdir>%t</lock-workdir>\n", r.LockWorkdir))
	rr.data.WriteString(fmt.Sprintf("    <unbuffer>%t</unbuffer>\n", r.Unbuffer))
	rr.data.WriteString("  </recipe>\n")
	rr.data.WriteString("  <commands>\n")
}

// CommandStarted prints info about started command
func (rr *XMLRenderer) CommandStarted(c *recipe.Command) {
	rr.data.WriteString("    <command")

	if c.User != "" {
		rr.data.WriteString(fmt.Sprintf(" user=%q", c.User))
	}

	if c.Tag != "" {
		rr.data.WriteString(fmt.Sprintf(" tag=%q", c.Tag))
	}

	rr.data.WriteString(">\n")

	rr.data.WriteString(fmt.Sprintf("      <cmdline>%s</cmdline>\n", rr.escapeData(c.GetCmdline())))
	rr.data.WriteString(fmt.Sprintf("      <description>%s</description>\n", rr.escapeData(c.Description)))

	if len(c.Env) != 0 {
		rr.data.WriteString(fmt.Sprint("      <environment>\n"))

		for _, variable := range c.Env {
			rr.data.WriteString(fmt.Sprintf(
				"        <variable name=%q value=%q />\n",
				rr.escapeData(strutil.ReadField(variable, 0, false, '=')),
				rr.escapeData(strutil.ReadField(variable, 1, false, '=')),
			))
		}

		rr.data.WriteString(fmt.Sprint("      </environment>\n"))
	}

	rr.data.WriteString("      <actions>\n")
}

// CommandSkipped prints info about skipped command
func (rr *XMLRenderer) CommandSkipped(c *recipe.Command, isLast bool) {
	return
}

// CommandFailed prints info about failed command
func (rr *XMLRenderer) CommandFailed(c *recipe.Command, err error) {
	rr.data.WriteString("      </actions>\n")
	rr.data.WriteString(fmt.Sprintf("      <status failed=\"true\">%v</status>\n", err))
	rr.data.WriteString("    </command>\n")
}

// CommandFailed prints info about executed command
func (rr *XMLRenderer) CommandDone(c *recipe.Command, isLast bool) {
	rr.data.WriteString("      </actions>\n")
	rr.data.WriteString("      <status failed=\"false\"></status>\n")
	rr.data.WriteString("    </command>\n")
}

// ActionStarted prints info about action in progress
func (rr *XMLRenderer) ActionStarted(a *recipe.Action) {
	rr.data.WriteString("        <action>\n")
	rr.data.WriteString(fmt.Sprintf("          <name>%s</name>\n", rr.formatActionName(a)))
	rr.data.WriteString("          <arguments>\n")

	for index := range a.Arguments {
		arg, _ := a.GetS(index)
		rr.data.WriteString(fmt.Sprintf(
			"            <argument>%s</argument>\n",
			rr.escapeData(arg),
		))
	}

	rr.data.WriteString("          </arguments>\n")
}

// ActionFailed prints info about failed action
func (rr *XMLRenderer) ActionFailed(a *recipe.Action, err error) {
	rr.data.WriteString(fmt.Sprintf("          <status failed=\"true\">%v</status>\n", err))
	rr.data.WriteString("        </action>\n")
}

// ActionDone prints info about successfully finished action
func (rr *XMLRenderer) ActionDone(a *recipe.Action, isLast bool) {
	rr.data.WriteString("          <status failed=\"false\"></status>\n")
	rr.data.WriteString("        </action>\n")
}

// Result prints info about test results
func (rr *XMLRenderer) Result(passes, fails, skips int) {
	rr.data.WriteString("  </commands>\n")
	rr.data.WriteString(fmt.Sprintf(
		"  <result passed=\"%d\" failed=\"%d\" skipped=\"%d\" duration=\"%g\" />\n",
		passes, fails, skips, time.Since(rr.start).Seconds(),
	))
	rr.data.WriteString("</report>")

	fmt.Println(rr.data.String())
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (rr *XMLRenderer) escapeData(data string) string {
	data = strings.ReplaceAll(data, "<", "&lt;")
	data = strings.ReplaceAll(data, ">", "&gt;")
	data = strings.ReplaceAll(data, "&", "&amp;")

	return data
}

// formatActionName format action name
func (rr *XMLRenderer) formatActionName(a *recipe.Action) string {
	if a.Negative {
		return rr.escapeData("!" + a.Name)
	}

	return rr.escapeData(a.Name)
}
