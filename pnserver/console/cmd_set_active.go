// --------------------------------------------------------------------
// cmd_set_active.go -- Sets a project or designer active.
//
// Created 2018-09-27 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/pnsql"
	"fmt"
)

func init() {
	RegistorCmd("set-active", "thing", "Sets a project or designer active.", handle_set_active)
}

func handle_set_active(cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		fmt.Printf("Not enough args.\n")
		return
	}
	thing := args[1]
	if pnsql.IsDesigner(thing) {
		err := pnsql.SetDesignerActive(thing, true)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("Success.\n")
		}
		return
	}

	if pnsql.IsProject(thing) {
		err := pnsql.SetProjectActive(thing, true)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("Success.\n")
		}
		return
	}

	fmt.Printf("The argument %q is neigher a ProjectId or a Designer's name.\n", thing)
}
