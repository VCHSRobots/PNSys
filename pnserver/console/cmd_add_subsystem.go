// --------------------------------------------------------------------
// cmd_add_subsystem.go -- Adds a subsystem.
//
// Created 2018-09-27 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
)

var gTopic_add_subsystem string = `
The add-subsystem command will add a new subsystem to the database. The format of
the command is:

  add-subsystem ppp-ss desc="comments about the subsystem" 

where ppp is the Project Id, ss is the new subsystem Id, and desc is the description
parameter, which is required.

`

func init() {
	RegistorCmd("add-subsystem", "ppp-ss", "Adds a new subsystem (see topic).", handle_add_subsystem)
	RegistorTopic("add-subsystem", gTopic_add_subsystem)
	RegistorArg("ppp-ss", "A subsystem Id, consisting of both the project and subsystem Ids.")
}

func handle_add_subsystem(c *Context, cmdline string) {
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
	projectid, subsystemid, err := pnsql.SplitProjectId(args[1])
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if subsystemid == "" {
		c.Printf("Subsystem not specified.\n")
		return
	}
	desc, ok := util.MapAlias(params, "description", "Description", "Desc", "desc")
	if !ok || util.Blank(desc) {
		c.Printf("A description must be provided.\n")
		return
	}

	err = pnsql.AddSubsystem(projectid, subsystemid, desc)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	c.Printf("Success.\n")
}
