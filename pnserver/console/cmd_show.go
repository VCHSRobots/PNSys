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

func handle_show_thing(cmdline string) {
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
	err = pnsql.CheckDesignerNameText(thing)
	if err == nil {
		print_designer_info(thing)
		return
	}
	pt, _ := pnsql.ClassifyPN(thing)
	if pt == pnsql.PNType_Supplier {
		print_supplier_info(thing)
		return
	}

	projectid, subsystemid, err := pnsql.SplitProjectId(thing)
	if err == nil {
		if subsystemid == "" {
			print_project_info(projectid)
		} else {
			print_subsystem_info(projectid, subsystemid)
		}
		return
	}
	pn, err := pnsql.StrToEpicPN(thing)
	if err == nil {
		print_partnumber_info(pn)
		return
	}
	fmt.Printf("%q cannot be indentified.\n", thing)
}

func print_supplier_info(pns string) {
	part, err := pnsql.GetSupplierPart(pns)
	if err != nil {
		fmt.Printf("Database error while searching for %s. Err=%v\n", pns, err)
		return
	}
	if part == nil {
		fmt.Printf("Part %s does not exist.\n", pns)
		return
	}
	fmt.Printf("\n")
	fmt.Printf("Found Supplier Part\n")
	fmt.Printf("-------------------\n")
	fmt.Printf("Part Number = %s\n", part.PNString())
	fmt.Printf("Description = %s\n", part.Description)
	fmt.Printf("Designer    = %s\n", part.Designer)
	fmt.Printf("Date Issued = %s\n", part.DateIssued.Format("2006-01-02"))
	fmt.Printf("Vendor      = %s\n", part.Vendor)
	fmt.Printf("VendorPN    = %s\n", part.VendorPN)
	fmt.Printf("Weblink     = %s\n", part.WebLink)
	fmt.Printf("PID         = %s\n", part.PID)
}

func print_partnumber_info(pn *pnsql.EpicPN) {
	part, err := pnsql.GetEpicPart(pn.PNString())
	if err != nil {
		fmt.Printf("Database error while searching for part %s. Err=%v\n", pn.PNString(), err)
		return
	}
	if part == nil {
		fmt.Printf("Part %s not found.", pn.PNString())
		return
	}
	fmt.Printf("\n")
	fmt.Printf("Found Epic Part\n")
	fmt.Printf("---------------\n")
	fmt.Printf("Part Number = %s\n", part.PNString())
	fmt.Printf("Description = %s\n", part.Description)
	fmt.Printf("Designer    = %s\n", part.Designer)
	fmt.Printf("Date Issued = %s\n", part.DateIssued.Format("2006-01-02"))
	fmt.Printf("PID         = %s\n", part.PID)
}

func print_designer_info(name string) {
	designer, err := pnsql.GetDesigner(name)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	params := make(map[string]string, 5)
	params["Designer"] = designer.Name
	prts := pnsql.FilterEpicParts(params)
	nepic := len(prts)
	spprts := pnsql.FilterSupplierParts(params)
	nsp := len(spprts)
	fmt.Printf("\n")
	fmt.Printf("Found Designer\n")
	fmt.Printf("--------------\n")
	fmt.Printf("Name    = %s\n", designer.Name)
	fmt.Printf("Year0   = %s\n", designer.Year0)
	fmt.Printf("Active  = %t\n", designer.Active)
	fmt.Printf("# Epic Parts = %d\n", nepic)
	fmt.Printf("# SP Parts   = %d\n", nsp)
}

func print_project_info(projectid string) {
	prj, err := pnsql.GetProject(projectid)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	params := make(map[string]string, 5)
	params["ProjectId"] = projectid
	prts := pnsql.FilterEpicParts(params)
	nparts := len(prts)
	fmt.Printf("\n")
	fmt.Printf("Found Project\n")
	fmt.Printf("-------------\n")
	fmt.Printf("ProjectId   = %s\n", prj.ProjectId)
	fmt.Printf("Description = %s\n", prj.Description)
	fmt.Printf("Year0       = %s\n", prj.Year0)
	fmt.Printf("Active      = %t\n", prj.Active)
	fmt.Printf("# Parts     = %d\n", nparts)
	if len(prj.Subsystems) <= 0 {
		fmt.Printf("No Subsystems.")
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
		fmt.Printf("%s", tbl.Text())
	}
}

func print_subsystem_info(projectid, subsystemid string) {
	ss, err := pnsql.GetSubsystem(projectid, subsystemid)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	params := make(map[string]string, 5)
	fmt.Printf("\n")
	fmt.Printf("Found Subsystem\n")
	fmt.Printf("---------------\n")
	params["ProjectId"] = ss.ProjectId
	params["SubsystemId"] = ss.SubsystemId
	prts := pnsql.FilterEpicParts(params)
	nparts := len(prts)
	fmt.Printf("ProjectId    = %s\n", ss.ProjectId)
	fmt.Printf("SubsystemId  = %s\n", ss.SubsystemId)
	fmt.Printf("Description  = %s\n", ss.Description)
	fmt.Printf("# Parts      = %d\n", nparts)
}
