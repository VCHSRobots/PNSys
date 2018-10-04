// --------------------------------------------------------------------
// cmd_add_designer.go -- Adds a new designer.
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"time"
)

var gTopic_add_designer string = `
The add-designer command adds a new designer to the database. The format of the
command is:

  add-designer name year0="yyyy"

 where the name is in quotes, and is in the formant "X. YYYY" where X is the
 first initial, and YYYY is the last name.  The year0 parameter is optional, 
 and if not provided, the current year will be used as the default.
 
 `

func init() {
	RegistorCmd("add-designer", "name", "Adds a new designer (see topic).", handle_add_designer)
	RegistorArg("name", "Desinger's name, in 'X. LastName' format.")
	RegistorTopic("add-designer", gTopic_add_designer)
}

func handle_add_designer(c *Context, cmdline string) {
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
	name := args[1]
	year0, haveyear := util.MapAlias(params, "Year0", "year0", "Year", "year")
	if !haveyear {
		year0 = time.Now().Format("2006")
	}
	err = pnsql.CheckYear0Text(year0)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	if pnsql.IsDesigner(name) {
		c.Printf("Designer %s already exists.\n")
		return
	}

	err = pnsql.AddDesigner(name, year0)
	if err != nil {
		c.Printf("Error: %v\n", err)
		return
	}
	c.Printf("Success.\n")
}
