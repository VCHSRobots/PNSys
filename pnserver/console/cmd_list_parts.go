// --------------------------------------------------------------------
// cmd_list_parts.go -- Gets a part.
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
	"sort"
	"strings"
)

var g_topic_list_parts = `
The command list-parts will list all parts in the database, both epic and supplier.
Parameters can be applied to filter (search) the output.  The format of the command 
is:

    list-parts p1=xxx p2=yyy p3=zzz

where you can have as many parameters (p1, p2, ...) as you please to narrow the
results.  The possible parameters are:

  Type, type                               -- for the part type (either epic or supplier)
  ProjectId, projectid, project prj        -- for the project,
  SubsystemId, subsystemid, subsystem, sub -- for the Subsystem,
  PartType, parttype, part                 -- for the Part Type,
  Designer, designer                       -- for the designer,
  DateBefore, datebefore, date0, before    -- for parts before a date,
  DateAfter, dateafter, date1, after       -- for parts after a date,
  Category, category, Cat, cat             -- for parts with a supplier category,
  Vender, vender, ven                      -- for parts of a vender name (lazy match),
  VenderPN, vendernp, vpn                  -- for parts of a vender's PN (lazy match),
  WebLink, Weblink, weblink                -- for the web line (lazy match),
  Description, description, desc           -- for the description (lazy match).

For items that are marked 'lazy' the match is any case insensitive substring.  For items
such as desinger, a perfect match is required.  If a parameter is not used, then all
the parts match that field. 

`

func init() {
	RegistorCmd("list-parts", "field=xxx", "Lists parts with filters (see topic).", handle_list_parts)
	RegistorTopic("list-parts", g_topic_list_parts)
}

func handle_list_parts(cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	bMustBeEpic := false
	bMustBeSupplier := false

	stype, havetype := util.MapAlias(params, "Type", "type")
	if havetype {
		stype = strings.ToLower(stype)
		if stype != "epic" && stype != "supplier" && stype != "sup" {
			fmt.Printf("Illegal parts type specified (%q). Must be either 'epic' or 'supplier'.\n", stype)
			return
		}
		if stype == "sup" {
			stype = "supplier"
		}
		if stype == "supplier" {
			bMustBeSupplier = true
		}
		if stype == "epic" {
			bMustBeEpic = true
		}
	}

	_, haveprj := util.MapAlias(params, "ProjectId", "projectid", "project", "prj")
	_, havesub := util.MapAlias(params, "SubsystemId", "subsystemid", "subsystem", "sub")
	_, havept := util.MapAlias(params, "PartType", "parttype", "part")
	_, havecat := util.MapAlias(params, "Category", "category", "Cat", "cat")
	_, havevendor := util.MapAlias(params, "Vender", "vender", "ven")
	_, havevendorpn := util.MapAlias(params, "VenderPN", "vendernp", "vpn")
	_, haveweblink := util.MapAlias(params, "WebLink", "Weblink", "weblink")

	if haveprj || havesub || havept {
		bMustBeEpic = true
	}
	if havecat || havevendor || havevendorpn || haveweblink {
		bMustBeSupplier = true
	}

	if bMustBeEpic && bMustBeSupplier {
		fmt.Printf("Incompatible parameters -- no parts can be found.\n")
		return
	}

	var epiclst []*pnsql.EpicPart
	var suplst []*pnsql.SupplierPart

	if !bMustBeSupplier {
		epiclst = pnsql.FilterEpicParts(params)
	}
	if !bMustBeEpic {
		suplst = pnsql.FilterSupplierParts(params)
	}

	// Do we have a mixture?
	if len(epiclst) > 0 && len(suplst) > 0 {
		tbl := util.NewTable("Part Number", "Description", "Designer", "Date", "Other Information")
		type bucket struct {
			pn, desc, designer, sdate, other string
		}
		blst := make([]*bucket, 0, len(epiclst)+len(suplst))
		for _, p := range epiclst {
			spn := p.PNString()
			sdesc := util.FixStrLen(p.Description, 40, " ")
			sdate := p.DateIssued.Format("2006-01-02")
			sother := fmt.Sprintf("Project: %s, PartType: %s", p.ProjectDesc(), p.PartTypeDesc())
			blst = append(blst, &bucket{spn, sdesc, p.Designer, sdate, sother})
		}
		for _, p := range suplst {
			spn := p.PNString()
			sdesc := util.FixStrLen(p.Description, 40, " ")
			sdate := p.DateIssued.Format("2006-01-02")
			sother := fmt.Sprintf("Category: %s, Vendor: %s VPN: %s", p.CategoryDesc(), p.Vendor, p.VendorPN)
			blst = append(blst, &bucket{spn, sdesc, p.Designer, sdate, sother})
		}
		sorter := func(i, j int) bool {
			return blst[i].sdate < blst[j].sdate
		}
		sort.Slice(blst, sorter)
		for _, r := range blst {
			tbl.AddRow(r.pn, r.desc, r.designer, r.sdate, r.other)
		}
		fmt.Printf("\n%s%d parts found.\n", tbl.Text(), len(blst))
		return
	} else if len(epiclst) > 0 {
		tbl := util.NewTable("Part Number", "Description", "Designer", "Date", "Project", "Subsystem", "PartType", "PID")
		type bucket struct {
			pn, desc, designer, sdate, prj, subsys, parttype, pid string
		}
		blst := make([]*bucket, 0, len(epiclst))
		for _, p := range epiclst {
			spn := p.PNString()
			sdesc := util.FixStrLen(p.Description, 40, " ")
			sdate := p.DateIssued.Format("2006-01-02")
			blst = append(blst, &bucket{spn, sdesc, p.Designer, sdate, p.ProjectDesc(), p.SubsystemDesc(), p.PartTypeDesc(), p.PID.String()})
		}
		sorter := func(i, j int) bool {
			return blst[i].sdate < blst[j].sdate
		}
		sort.Slice(blst, sorter)
		for _, r := range blst {
			tbl.AddRow(r.pn, r.desc, r.designer, r.sdate, r.prj, r.subsys, r.parttype, r.pid)
		}
		fmt.Printf("\n%s%d parts found.\n", tbl.Text(), len(blst))
		return
	} else if len(suplst) > 0 {
		tbl := util.NewTable("Part Number", "Description", "Designer", "Date", "Category", "Vendor", "VendorPN", "WebLink", "PID")
		type bucket struct {
			pn, desc, designer, sdate, cat, ven, venpn, weblink, pid string
		}
		blst := make([]*bucket, 0, len(suplst))
		for _, p := range suplst {
			spn := p.PNString()
			sdesc := util.FixStrLen(p.Description, 40, " ")
			sdate := p.DateIssued.Format("2006-01-02")
			slink := util.FixStrLen(p.WebLink, 30, "...")
			blst = append(blst, &bucket{spn, sdesc, p.Designer, sdate, p.CategoryDesc(), p.Vendor, p.VendorPN, slink, p.PID.String()})
		}
		sorter := func(i, j int) bool {
			return blst[i].sdate < blst[j].sdate
		}
		sort.Slice(blst, sorter)
		for _, r := range blst {
			tbl.AddRow(r.pn, r.desc, r.designer, r.sdate, r.cat, r.ven, r.venpn, r.weblink, r.pid)
		}
		fmt.Printf("\n%s%d parts found.\n", tbl.Text(), len(blst))
		return
	} else {
		fmt.Printf("No parts found.\n")
		return
	}
}
