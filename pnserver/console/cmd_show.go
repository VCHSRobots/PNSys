// --------------------------------------------------------------------
// cmd_show.go -- Shows info about anything.
//
// Created 2018-09-27 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
)

func init() {
	RegistorCmd("show", "thing", "Shows info about any object.", handle_show_thing)
	RegistorArg("thing", "A project Id, a subsystem Id, a desinger's name, or a part number.")
}

func handle_show_thing(c *Context, cmdline string) {
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
	err = pnsql.CheckDesignerNameText(thing)
	if err == nil {
		print_designer_info(c, thing)
		return
	}
	pt, _ := pnsql.ClassifyPN(thing)
	if pt == pnsql.PNType_Supplier {
		print_supplier_info(c, thing)
		return
	}

	projectid, subsystemid, err := pnsql.SplitProjectId(thing)
	if err == nil {
		if subsystemid == "" {
			print_project_info(c, projectid)
		} else {
			print_subsystem_info(c, projectid, subsystemid)
		}
		return
	}
	pn, err := pnsql.StrToEpicPN(thing)
	if err == nil {
		print_partnumber_info(c, pn)
		return
	}
	c.Printf("%q cannot be indentified.\n", thing)
}

func print_supplier_info(c *Context, pns string) {
	part, err := pnsql.GetSupplierPart(pns)
	if err != nil {
		c.Printf("Database error while searching for %s. Err=%v\n", pns, err)
		return
	}
	if part == nil {
		c.Printf("Part %s does not exist.\n", pns)
		return
	}
	c.Printf("\n")
	c.Printf("Found Supplier Part\n")
	c.Printf("-------------------\n")
	c.Printf("Part Number = %s\n", part.PNString())
	c.Printf("Description = %s\n", part.Description)
	c.Printf("Designer    = %s\n", part.Designer)
	c.Printf("Date Issued = %s\n", part.DateIssued.Format("2006-01-02"))
	c.Printf("Vendor      = %s\n", part.Vendor)
	c.Printf("VendorPN    = %s\n", part.VendorPN)
	c.Printf("Weblink     = %s\n", part.WebLink)
	c.Printf("PID         = %s\n", part.PID)
}

func print_partnumber_info(c *Context, pn *pnsql.EpicPN) {
	part, err := pnsql.GetEpicPart(pn.PNString())
	if err != nil {
		c.Printf("Database error while searching for part %s. Err=%v\n", pn.PNString(), err)
		return
	}
	if part == nil {
		c.Printf("Part %s not found.", pn.PNString())
		return
	}
	c.Printf("\n")
	c.Printf("Found Epic Part\n")
	c.Printf("---------------\n")
	c.Printf("Part Number = %s\n", part.PNString())
	c.Printf("Description = %s\n", part.Description)
	c.Printf("Designer    = %s\n", part.Designer)
	c.Printf("Date Issued = %s\n", part.DateIssued.Format("2006-01-02"))
	c.Printf("PID         = %s\n", part.PID)
}

func print_designer_info(c *Context, name string) {
	designer, err := pnsql.GetDesigner(name)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	params := make(map[string]string, 5)
	params["Designer"] = designer.Name
	prts := pnsql.FilterEpicParts(params)
	nepic := len(prts)
	spprts := pnsql.FilterSupplierParts(params)
	nsp := len(spprts)
	c.Printf("\n")
	c.Printf("Found Designer\n")
	c.Printf("--------------\n")
	c.Printf("Name    = %s\n", designer.Name)
	c.Printf("Year0   = %s\n", designer.Year0)
	c.Printf("Active  = %t\n", designer.Active)
	c.Printf("# Epic Parts = %d\n", nepic)
	c.Printf("# SP Parts   = %d\n", nsp)
}

func print_project_info(c *Context, projectid string) {
	prj, err := pnsql.GetProject(projectid)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	params := make(map[string]string, 5)
	params["ProjectId"] = projectid
	prts := pnsql.FilterEpicParts(params)
	nparts := len(prts)
	c.Printf("\n")
	c.Printf("Found Project\n")
	c.Printf("-------------\n")
	c.Printf("ProjectId   = %s\n", prj.ProjectId)
	c.Printf("Description = %s\n", prj.Description)
	c.Printf("Year0       = %s\n", prj.Year0)
	c.Printf("Active      = %t\n", prj.Active)
	c.Printf("# Parts     = %d\n", nparts)
	if len(prj.Subsystems) <= 0 {
		c.Printf("No Subsystems.")
	} else {
		tbl := util.NewTable("Subsystem", "Description", "# Parts")
		for _, s := range prj.Subsystems {
			nc := 0
			for _, p := range prts {
				if p.SubsystemId == s.SubsystemId {
					nc++
				}
			}
			tbl.AddRow(s.SubsystemId, s.Description, fmt.Sprintf("%d", nc))
		}
		c.Printf("%s", tbl.Text())
	}
}

func print_subsystem_info(c *Context, projectid, subsystemid string) {
	ss, err := pnsql.GetSubsystem(projectid, subsystemid)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	params := make(map[string]string, 5)
	c.Printf("\n")
	c.Printf("Found Subsystem\n")
	c.Printf("---------------\n")
	params["ProjectId"] = ss.ProjectId
	params["SubsystemId"] = ss.SubsystemId
	prts := pnsql.FilterEpicParts(params)
	nparts := len(prts)
	c.Printf("ProjectId    = %s\n", ss.ProjectId)
	c.Printf("SubsystemId  = %s\n", ss.SubsystemId)
	c.Printf("Description  = %s\n", ss.Description)
	c.Printf("# Parts      = %d\n", nparts)
}
