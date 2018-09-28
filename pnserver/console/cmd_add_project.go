// --------------------------------------------------------------------
// cmd_add_project.go -- Adds a project.
//
// Created 2018-09-26 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
	"time"
)

var gTopic_add_project string = `
The add-project command will add a new project to the database. The format of
the command is:

  add-project ppp desc="comments about the project" year0="yyyy"

where ppp is the Project Id.  The desc parameter must be provided. The year0
parameters is options, and defaults to the current year.

`

func init() {
	RegistorCmd("add-project", "ppp", "Adds a new project (see topic).", handle_add_project)
	RegistorTopic("handle_add_project", gTopic_add_project)
	RegistorArg("ppp", "A project Id.")
}

func handle_add_project(cmdline string) {
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
	projectid := args[1]
	err = pnsql.CheckProjectIdText(projectid)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	desc, ok := util.MapAlias(params, "description", "Description", "Desc", "desc")
	if !ok || util.Blank(desc) {
		fmt.Printf("A description must be provided.\n")
		return
	}
	year0, _ := util.MapAlias(params, "year0", "Year0", "year", "Year")
	if year0 == "" {
		year0 = time.Now().Format("2006")
	}
	err = pnsql.AddProject(projectid, desc, year0)
	if err != nil {
		fmt.Printf("Err adding project. Err=%v\n", err)
		return
	}
	fmt.Printf("Success.\n")
}
