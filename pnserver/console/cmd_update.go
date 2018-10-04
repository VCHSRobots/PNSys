// --------------------------------------------------------------------
// cmd_update.go -- updates about anything.
//
// Created 2018-09-27 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
	"strings"
)

var gTopic_update_thing string = `
The command update allows changing information about most objects in the 
database, such as designers, projects, subsystems, and parts. The format
of the command is:

    update thing p1=xxx p2=yyy p3=zzz

where thing can be a designer's name (in quotes), a part number (in the
form "ppp-ss-000", a subsystem Id (in the form "ppp-ss") or a project id
(in the form "ppp"). Also, thing can be in the form "SP-cc-000" for 
supplier parts.

The p1, p2, p3 denotes the attributes that can be changed.  These are:

    description, desc -- for parts, subsystems, and projects
    year0             -- for projects and designers
    active            -- for projects and desginers
    dateissued, date  --- for all parts
    designer          -- for parts
    vendor, ven       -- for supplier parts
    vendorpn, vpn     -- for supplier parts
    weblink, link     -- for supplier parts

`

func init() {
	RegistorCmd("update", "thing", "Updates all objects (see topic).", handle_update_thing)
	RegistorTopic("update", gTopic_update_thing)
}

func handle_update_thing(c *Context, cmdline string) {
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
		update_designer(c, thing, params)
		return
	}

	pt, _ := pnsql.ClassifyPN(thing)
	if pt == pnsql.PNType_Epic {
		pn, err := pnsql.StrToEpicPN(thing)
		if err == nil {
			update_epic_part(c, pn, params)
		} else {
			fmt.Printf("%v\n", err)
		}
		return
	}

	if pt == pnsql.PNType_Supplier {
		update_supplier_part(c, thing, params)
		return
	}

	projectid, subsystemid, err := pnsql.SplitProjectId(thing)
	if err == nil {
		if subsystemid == "" {
			update_project(c, projectid, params)
		} else {
			update_subsystem(c, projectid, subsystemid, params)
		}
		return
	}
	c.Printf("%q cannot be indentified.\n", thing)
}

func update_epic_part(c *Context, pn *pnsql.EpicPN, params map[string]string) {
	part, err := pnsql.GetEpicPart(pn.PNString())
	if err != nil {
		c.Printf("Database error while searching for part %s. Err=%v\n", pn.PNString(), err)
		return
	}
	if part == nil {
		c.Printf("Part %s not found.\n", pn.PNString())
		return
	}
	nupdates := 0
	desc, havedesc := util.MapAlias(params, "Description", "description", "Desc", "desc")
	sdate, havedate := util.MapAlias(params, "DateIssued", "dateissued", "Date", "date")
	dname, havedesigner := util.MapAlias(params, "Designer", "designer")
	if havedesc {
		if desc != part.Description {
			err := pnsql.SetEpicPartDescription(part, desc)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Description on part %s.\n", part.PNString())
			nupdates++
		} else {
			c.Printf("Description specified, but equal to current value.\n")
		}
	}
	if havedate {
		date, err := util.ParseGenericTime(sdate)
		if err != nil {
			c.Printf("%v\n", err)
			return
		}
		if date.Format("2006-01-02") != part.DateIssued.Format("2006-01-02") {
			err = pnsql.SetEpicPartDateIssued(part, date)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Date Issued on part %s.\n", part.PNString())
			nupdates++
		} else {
			c.Printf("Date specified, but equal to current value.\n")
		}

	}
	if havedesigner {
		if part.Designer != dname {
			err = pnsql.SetEpicPartDesigner(part, dname)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Designer on part %s.\n", part.PNString())
			nupdates++
		} else {
			c.Printf("Designer specified, but equal to current value.\n")
		}
	}
	if nupdates == 0 {
		c.Printf("Nothing updated on part %s.\n", part.PNString())
	}
}

func update_designer(c *Context, name string, params map[string]string) {
	designer, err := pnsql.GetDesigner(name)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	nupdates := 0
	year0, haveyear := util.MapAlias(params, "Year0", "year0", "Year", "year")
	sact, haveactive := util.MapAlias(params, "Active", "active")
	if haveyear {
		if designer.Year0 != year0 {
			err = pnsql.SetDesignerYear0(designer.Name, year0)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Year0 for Designer %s.\n", designer.Name)
			nupdates++
		} else {
			c.Printf("Year0 specified, but equal to current value.\n")
		}
	}
	if haveactive {
		act := false
		ssact := strings.ToLower(sact)
		if ssact == "true" || ssact == "t" || ssact == "yes" || ssact == "y" {
			act = true
		} else if ssact == "false" || ssact == "f" || ssact == "no" || ssact == "n" {
			act = false
		} else {
			c.Printf("Illegal value (%q) for active.", sact)
			return
		}
		if designer.Active != act {
			err = pnsql.SetDesignerActive(designer.Name, act)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Active for Designer %s.\n", designer.Name)
			nupdates++
		} else {
			c.Printf("Active specified, but equal to current value.\n")
		}
	}
	if nupdates == 0 {
		c.Printf("Nothing updated for designer %s.\n", designer.Name)
	}
}

func update_supplier_part(c *Context, pns string, params map[string]string) {

	part, err := pnsql.GetSupplierPart(pns)
	if err != nil {
		c.Printf("Database error while searching for part %s. Err=%v\n", pns, err)
		return
	}
	if part == nil {
		c.Printf("Part %s not found.\n", pns)
		return
	}

	nupdates := 0
	desc, havedesc := util.MapAlias(params, "Description", "description", "Desc", "desc")
	sdate, havedate := util.MapAlias(params, "DateIssued", "dateissued", "Date", "date")
	dname, havedesigner := util.MapAlias(params, "Designer", "designer")
	sven, havevendor := util.MapAlias(params, "Vendor", "vendor", "ven")
	svpn, havevpn := util.MapAlias(params, "Vendor", "vendor", "ven")
	slink, havelink := util.MapAlias(params, "WebLink", "weblink", "link")

	if havedesc {
		if part.Description != desc {
			err := pnsql.SetSupplierPartDescription(part, desc)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Description on part %s.\n", part.PNString())
			nupdates++
		} else {
			c.Printf("Description specified, but equal to current value.\n")
		}
	}
	if havedate {
		date, err := util.ParseGenericTime(sdate)
		if err != nil {
			c.Printf("%v\n", err)
			return
		}
		if date.Format("2006-01-02") != part.DateIssued.Format("2006-01-02") {
			err = pnsql.SetSupplierPartDateIssued(part, date)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Date Issued on part %s.\n", part.PNString())
			nupdates++
		} else {
			c.Printf("Date specified, but equal to current value.\n")
		}
	}
	if havedesigner {
		if part.Designer != dname {
			err = pnsql.SetSupplierPartDesigner(part, dname)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Designer on part %s.\n", part.PNString())
			nupdates++
		} else {
			c.Printf("Designer specified, but equal to current value.\n")
		}
	}

	if havevendor {
		if part.Vendor != sven {
			err = pnsql.SetSupplierPartVendor(part, sven)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Vendor on part %s.\n", part.PNString())
			nupdates++
		} else {
			c.Printf("Vendor specified, but equal to current value.\n")
		}
	}

	if havevpn {
		if part.VendorPN != svpn {
			err = pnsql.SetSupplierPartVendorPN(part, svpn)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of VendorPN on part %s.\n", part.PNString())
			nupdates++
		} else {
			c.Printf("VendorPN specified, but equal to current value.\n")
		}
	}

	if havelink {
		if part.WebLink != slink {
			err = pnsql.SetSupplierPartWebLink(part, slink)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of WebLink on part %s.\n", part.PNString())
			nupdates++
		} else {
			c.Printf("WebLink specified, but equal to current value.\n")
		}
	}

	if nupdates == 0 {
		c.Printf("Nothing updated on part %s.\n", part.PNString())
	}
}

func update_project(c *Context, projectid string, params map[string]string) {
	prj, err := pnsql.GetProject(projectid)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	nupdates := 0
	desc, havedesc := util.MapAlias(params, "Description", "description", "Desc", "desc")
	year0, haveyear := util.MapAlias(params, "Year0", "year0", "Year", "year")
	sact, haveactive := util.MapAlias(params, "Active", "active")
	if havedesc {
		if prj.Description != desc {
			err = pnsql.SetProjectDescription(prj.ProjectId, desc)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Description on project %s.\n", prj.ProjectId)
			nupdates++
		} else {
			c.Printf("Description specified, but equal to current value.\n")
		}
	}
	if haveyear {
		if prj.Year0 != year0 {
			err = pnsql.SetProjectYear0(prj.ProjectId, year0)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Year0 on project %s.\n", prj.ProjectId)
			nupdates++
		} else {
			c.Printf("Year0 specified, but equal to current value.\n")
		}
	}
	if haveactive {
		act := false
		ssact := strings.ToLower(sact)
		if ssact == "true" || ssact == "t" || ssact == "yes" || ssact == "y" {
			act = true
		} else if ssact == "false" || ssact == "f" || ssact == "no" || ssact == "n" {
			act = false
		} else {
			c.Printf("Illegal value (%q) for active.", sact)
			return
		}
		if prj.Active != act {
			err = pnsql.SetProjectActive(prj.ProjectId, act)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Active on project %s.\n", prj.ProjectId)
			nupdates++
		} else {
			c.Printf("Active specified, but equal to current value.\n")
		}
	}
	if nupdates == 0 {
		c.Printf("Nothing updated for project %s.\n", prj.ProjectId)
	}
}

func update_subsystem(c *Context, projectid, subsystemid string, params map[string]string) {
	ss, err := pnsql.GetSubsystem(projectid, subsystemid)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	nupdates := 0
	desc, havedesc := util.MapAlias(params, "Description", "description", "Desc", "desc")
	if havedesc {
		if ss.Description != desc {
			err = pnsql.SetSubsystemDescription(ss.ProjectId, ss.SubsystemId, desc)
			if err != nil {
				c.Printf("%v\n", err)
				return
			}
			c.Printf("Successful update of Description on subsystem %s-%s.\n", ss.ProjectId, ss.SubsystemId)
			nupdates++
		} else {
			c.Printf("Description specified, but equal to current value.\n")
		}
	}
	if nupdates == 0 {
		c.Printf("Nothing updated for subsystem %s-%s.\n", ss.ProjectId, ss.SubsystemId)
	}
}
