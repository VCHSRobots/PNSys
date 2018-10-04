// --------------------------------------------------------------------
// cmd_set_active.go -- Sets a project or designer active.
//
// Created 2018-09-27 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/pnsql"
)

func init() {
	RegistorCmd("set-inactive", "thing", "Sets a project or designer inactive.", handle_set_inactive)
}

func handle_set_inactive(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		c.Printf("Not enough args.\n")
		return
	}
	thing := args[1]
	if pnsql.IsDesigner(thing) {
		err := pnsql.SetDesignerActive(thing, false)
		if err != nil {
			c.Printf("%v\n", err)
		} else {
			c.Printf("Success.\n")
		}
		return
	}

	if pnsql.IsProject(thing) {
		err := pnsql.SetProjectActive(thing, false)
		if err != nil {
			c.Printf("%v\n", err)
		} else {
			c.Printf("Success.\n")
		}
		return
	}

	c.Printf("The argument %q is neigher a ProjectId or a Designer's name.\n", thing)
}
