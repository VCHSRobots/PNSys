// --------------------------------------------------------------------
// cmd_add_part.go -- Adds a part.
//
// Created 2018-09-26 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"time"
)

var gTopic_add_part string = `
The add-part command will add an arbitary part to the database -- so be careful!
The part added must have a valid project and subsystem, and must not conflict
with an existing part.  However, the sequence number can be anything you like!
Use the command as follows for an epic part:

  add-part ppp-ss-0000 desiger="X. Name" desc="comments about the part" date="yyyy-mm-dd"

  or, for a supplier part:

  add-part SP-cc-000 desiger="X. Name" desc="comments about the part" date="yyyy-mm-dd"

For epic parts, you can also use the parameters 'Vendor', 'VendorPN', and 'WebLink'.
  
`

func init() {
	RegistorCmd("add-part", "pn", "Adds a part number (see topic).", handle_add_part)
	RegistorTopic("add-part", gTopic_add_part)
}

func handle_add_part(c *Context, cmdline string) {
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
	spn := args[1]
	tp, _ := pnsql.ClassifyPN(spn)
	if tp == pnsql.PNType_Epic {
		add_epic_part(c, spn, params)
		return
	}
	if tp == pnsql.PNType_Supplier {
		add_supplier_part(c, spn, params)
		return
	}
	c.Printf("Unrecognizable part number (%q)?\n", spn)
	return
}

func add_epic_part(c *Context, spn string, params map[string]string) {
	part, err := pnsql.GetEpicPart(spn)
	if err != nil {
		c.Printf("Database error while searching for part. Err=%v", err)
		return
	}
	if part != nil {
		c.Printf("Part %q current exists.  Cannot overwrite.\n", part.PNString())
		return
	}
	pns, err := pnsql.StrToEpicPN(spn)
	if err != nil {
		c.Printf("Invalid part number %q. Err=%v\n", spn, err)
		return
	}
	if !pnsql.IsProject(pns.ProjectId) {
		c.Printf("%s is not a current project.\n", pns.ProjectId)
	}
	if !pnsql.IsSubsystem(pns.ProjectId, pns.SubsystemId) {
		c.Printf("%s is not a current subsystem in the project %s.\n", pns.SubsystemId, pns.ProjectId)
	}
	designer, ok := params["Designer"]
	if !ok {
		designer, ok = params["designer"]
	}
	if !ok {
		c.Printf("A designer must be specified.\n")
		return
	}
	if !pnsql.IsDesigner(designer) {
		c.Printf("%s is not a current designer.\n", designer)
		return
	}
	desc, ok := util.MapAlias(params, "Description", "description", "desc")
	if !ok || util.Blank(desc) {
		c.Printf("A description must be provided.\n")
		return
	}
	sdate, ok := util.MapAlias(params, "Date", "date", "DateIssued", "dateissued")
	if sdate == "" || !ok {
		sdate = time.Now().Format("2006-01-02")
	}
	date, err := util.ParseGenericTime(sdate)
	if err != nil {
		c.Printf("Syntax error in date (%q).\n", sdate)
		return
	}

	pt := &pnsql.EpicPart{}
	pt.EpicPN = pns
	pt.Designer = designer
	pt.Description = desc
	pt.DateIssued = date
	err = pnsql.AddEpicPart(pt)
	if err != nil {
		c.Printf("Error adding epic part. Err=%v\n", err)
		return
	}
	pnsql.InvalidateEpicPartsCache()
	c.Printf("Success.\n")
}

func add_supplier_part(c *Context, spn string, params map[string]string) {
	part, err := pnsql.GetSupplierPart(spn)
	if err != nil {
		c.Printf("Database error while searching for %s. Err=%v\n", part.PNString(), err)
		return
	}
	if part != nil {
		c.Printf("Part %q current exists.  Cannot overwrite.\n", part.PNString())
		return
	}
	pn, err := pnsql.StrToSupplierPartPN(spn)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	designer, ok := params["Designer"]
	if !ok {
		designer, ok = params["designer"]
	}
	if !ok {
		c.Printf("A designer must be specified.\n")
		return
	}
	if !pnsql.IsDesigner(designer) {
		c.Printf("%s is not a current designer.\n", designer)
		return
	}
	desc, ok := util.MapAlias(params, "Description", "description", "desc")
	if !ok || util.Blank(desc) {
		c.Printf("A description must be provided.\n")
		return
	}
	sdate, ok := util.MapAlias(params, "Date", "date", "DateIssued", "dateissued")
	if sdate == "" || !ok {
		sdate = time.Now().Format("2006-01-02")
	}
	date, err := util.ParseGenericTime(sdate)
	if err != nil {
		c.Printf("Syntax error in date (%q).\n", sdate)
		return
	}
	vendor, _ := util.MapAlias(params, "Vendor", "vendor", "ven")
	vendorpn, _ := util.MapAlias(params, "VendorPN", "vendorpn", "vpn")
	weblink, _ := util.MapAlias(params, "WebLink", "weblink", "web", "link")

	pt := &pnsql.SupplierPart{}
	pt.SupplierPartPN = pn
	pt.Description = desc
	pt.Vendor = vendor
	pt.VendorPN = vendorpn
	pt.WebLink = weblink
	pt.Designer = designer
	pt.DateIssued = date

	err = pnsql.AddSupplierPart(pt)
	if err != nil {
		c.Printf("Error adding epic part. Err=%v\n", err)
		return
	}
	pnsql.InvalidateSupplierPartsCache()
	c.Printf("Success.\n")
}
