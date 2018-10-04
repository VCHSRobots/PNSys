// --------------------------------------------------------------------
// cmd_delete_project.go -- Deletes a project.
//
// Created 2018-09-27 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/pnsql"
)

var gTopic_delete_project string = `
The delete-designer will delete a project. The format of the command is:

    delete-project ppp

where ppp is the Project Id.  A project must be empty to be deleted.

`

func init() {
	RegistorCmd("delete-project", "ppp", "Deletes a project.", handle_delete_project)
	RegistorTopic("delete-project", gTopic_delete_project)
}

func handle_delete_project(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		c.Printf("Not enought args.\n")
		return
	}
	projectid := args[1]
	err = pnsql.DeleteProject(projectid)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	c.Printf("Success.\n")
}
