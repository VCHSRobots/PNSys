// --------------------------------------------------------------------
// cmd_list_projects.go -- Lists projects and sub systems
//
// Created 2018-09-26 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
)

var gTopic_list_projects string = `
The list-gTopic_list_projects command will list all the projects or a subset. The
format of the command is:

    list-projects active=[true/false] year0=yyyy

 where the parameters active and year0 are optional.  If provided the
 list will be filtered accordingly.

`

func init() {
	RegistorCmd("list-projects", "", "Lists all the Epic projects.", handle_list_projects)
	RegistorTopic("list-projects", gTopic_list_projects)
}

func handle_list_projects(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	act, doactive, err := ParseActive(params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	year0, doyear := util.MapAlias(params, "Year0", "year0", "Year", "year")

	icount := 0
	tbl := util.NewTable("Project", "Project Description", "Year0", "Active")
	for _, p := range pnsql.GetProjects() {
		if doactive && p.Active != act {
			continue
		}
		if doyear && p.Year0 != year0 {
			continue
		}
		sactive := ""
		if p.Active {
			sactive = "Yes"
		}
		tbl.AddRow(p.ProjectId, p.Description, p.Year0, sactive)
		icount++
	}
	c.Printf("\n%s%d projects found.\n", tbl.Text(), icount)
}
