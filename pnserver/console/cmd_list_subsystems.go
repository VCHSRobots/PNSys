// --------------------------------------------------------------------
// cmd_list_subsystems.go -- Lists subsystems
//
// Created 2018-09-26 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
)

var gTopic_list_subsystems string = `
The command list-subsystems and be used with an optional arguement to limit
the list to just those for a particular project, such as:
   
    list-subsystems P17 active=t/f year0=2006

You can also use the active and year0 parameters to filter the output
(based on their project's status).

`

func init() {
	RegistorCmd("list-subsystems", "[project]", "Lists subsystems (see topic)", handle_list_subsystems)
	RegistorTopic("list-subsystems", gTopic_list_subsystems)
}

func handle_list_subsystems(cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	act, doactive, err := ParseActive(params)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	year0, doyear := util.MapAlias(params, "Year0", "year0", "Year", "year")

	prj := ""
	if len(args) > 1 {
		prj = args[1]
	}
	if prj != "" {
		if !pnsql.IsProject(prj) {
			fmt.Printf("%q is not a project.  (Use project Id three letter Id.)\n", prj)
			return
		}
	}
	icount := 0
	tbl := util.NewTable("Project Id", "Subsystem Id", "Subsystem Description", "Active", "Year0")
	for _, p := range pnsql.GetProjects() {
		if doyear && p.Year0 != year0 {
			continue
		}
		if doactive && p.Active != act {
			continue
		}
		sactive := ""
		if p.Active {
			sactive = "Yes"
		}
		if prj == "" || prj == p.ProjectId {
			for _, sub := range p.Subsystems {
				tbl.AddRow(p.ProjectId, sub.SubsystemId, sub.Description, sactive, p.Year0)
				icount++
			}
		}
	}
	fmt.Printf("\n%s\n%d subsystems found.\n", tbl.Text(), icount)
}
