// --------------------------------------------------------------------
// cmd_update_project.go -- Updates a project.
//
// Created 2018-09-27 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
)

var gTopic_update_project string = `
The update-project command allows you to change the description or Year0 of a
project or subsystem. It works as follows for a project:

  update-project ppp desc="comments about the part" year0=2006

or for a subsystem:

  update-project ppp-ss desc="xxxx"

where ppp is the Project Id, and ss is the Subsystem Id.

One or both parameters (desc, year0) can be provided.

`

func init() {
	RegistorCmd("update-project", "prj", "Updates a project or subsystem's parameters (see topic).", handle_update_project)
	RegistorTopic("update-project", gTopic_update_project)
}

func handle_update_project(c *Context, cmdline string) {
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
	desc, _ := util.MapAlias(params, "Description", "description", "Desc", "desc")
	year0, _ := util.MapAlias(params, "Year0", "year0", "Year", "year")

	if subsystemid == "" {
		// Updateing a project.
		if !pnsql.IsProject(projectid) {
			c.Printf("Project %s does not exist.\n", projectid)
			return
		}
		nupdates := 0
		if desc != "" {
			err = pnsql.SetProjectDescription(projectid, desc)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			nupdates++
			c.Printf("Description for project %s updated successfully.\n", projectid)
		}
		if year0 != "" {
			err = pnsql.SetProjectYear0(projectid, year0)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			nupdates++
			c.Printf("Year0 for project %s updated successfully.\n", projectid)
		}
		if nupdates == 0 {
			c.Printf("Nothing updated for project %s.\n", projectid)
		}
		return
	} else {
		// Updateing a subsystem.
		if !pnsql.IsSubsystem(projectid, subsystemid) {
			c.Printf("Subsystem %s-%s does not exist.\n", projectid, subsystemid)
			return
		}
		nupdates := 0
		if desc != "" {
			err = pnsql.SetSubsystemDescription(projectid, subsystemid, desc)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			nupdates++
			c.Printf("Description for subsystem %s-%s updated successfully.\n", projectid, subsystemid)
		}
		if nupdates == 0 {
			c.Printf("Nothing updated for subsystem %s-%s.\n", projectid, subsystemid)
		}
		return
	}
}
