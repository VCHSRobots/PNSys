// --------------------------------------------------------------------
// find.go -- Shows Epic PN input form
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
	"github.com/gin-gonic/gin"
	"sort"
	"strings"
	"time"
)

type FindPartData struct {
	*HeaderData
	*SelectionBoxData
}

type FindPartPost struct {
	*HeaderData
	*TableData
	MainType    string
	Designer    string
	Project     string
	Subsystem   string
	PartType    string
	Category    string
	Vendor      string
	VendorPN    string
	WebLink     string
	Description string
	DateBefore  string
	DateAfter   string
	ErrorMsg    string
}

type TFindSubmitData struct {
	MainType    string `form:MainType`
	Designer    string `form:Designer`
	Project     string `form:Project`
	Subsystem   string `form:Subsystem`
	PartType    string `form:PartType`
	Category    string `form:Category`
	Vendor      string `form:Vendor`
	VendorPN    string `form:VendorPN`
	WebLink     string `form:WebLink`
	Description string `form:Description`
	DateBefore  string `form:DateBefore`
	DateAfter   string `form:DateAfter`
}

func init() {
	RegisterPage("/Find", Invoke_GET, authorizer, handle_find)
	RegisterPage("/SubmitFind", Invoke_POST, authorizer, handle_find_post)
}

func handle_find(c *gin.Context) {
	handle_find_with_error(c, "")
}

func handle_find_with_error(c *gin.Context, errmsg string) {

	data := &FindPartData{}
	data.HeaderData = GetHeaderData(c)
	data.PageTitle = "Find Epic and Supplier Parts"
	data.Instructions = ""
	data.StyleSheets = []string{"find"}
	data.OnLoadFuncJS = "startUp"
	data.ErrorMessage = errmsg

	var err error
	data.SelectionBoxData, err = GetSelectionBoxData()
	if err != nil {
		SendErrorPage(c, err)
		return
	}

	SendPage(c, data, "header", "menubar", "find", "footer")
}

func handle_find_post(c *gin.Context) {
	var sdata TFindSubmitData
	err := c.ShouldBind(&sdata)
	if err != nil {
		SendErrorPagef(c, "Unable to bind data in Find Post. Err=%v", err)
		return
	}

	data := &FindPartPost{}
	data.HeaderData = GetHeaderData(c)

	data.PageTitle = "Find Results"
	data.MainType = sdata.MainType
	data.Designer = sdata.Designer
	data.Project = sdata.Project
	data.Subsystem = sdata.Subsystem
	data.PartType = sdata.PartType
	data.Category = sdata.Category
	data.Vendor = sdata.Vendor
	data.VendorPN = sdata.VendorPN
	data.WebLink = sdata.WebLink
	data.Description = sdata.Description
	data.DateBefore = sdata.DateBefore
	data.DateAfter = sdata.DateAfter

	//partlst, tbltype, err := search_for_parts(sdata)
	plst, _, err := search_for_parts(sdata)
	if err != nil {
		handle_find_with_error(c, fmt.Sprintf("%v", err))
		return
	}

	if len(plst) <= 0 {
		handle_find_with_error(c, "No parts found.")
		return
	}

	data.TableData = new(TableData)
	data.Head = []string{"PN", "Description", "Designer", "DateIssued"}
	data.Rows = make([]TColumn, 0, 101)
	if len(plst) > 100 {
		data.LimitMsg = fmt.Sprintf("Showing first 100 parts out of %d found.", len(plst))
	}
	for i, p := range plst {
		sdesc := util.FixStrLen(p.Description, 50, "...")
		lnk := fmt.Sprintf(`<a href="/ShowPart?pn=%s">%s</a>`, p.PartNumber, p.PartNumber)
		cols := []string{lnk, sdesc, p.Designer, p.DateIssued.Format("2006-01-02")}
		data.Rows = append(data.Rows, TColumn{Cols: cols})
		if i >= 99 {
			break
		}
	}

	SendPage(c, data, "header", "menubar", "tablepage", "footer")
}

type CPart struct {
	MainType    string
	PartNumber  string
	Description string
	Designer    string
	DateIssued  time.Time
	SequenceNum int
	ProjectId   string
	SubsystemId string
	PartType    string
	Category    string
	Vendor      string
	VendorPN    string
	WebLink     string
	PID         string
}

// Searches for parts.  Returns a slice, and one of: "both", "epic", or "supplier" depending
// on the parts found.
func search_for_parts(spec TFindSubmitData) (parts []*CPart, tabletype string, err error) {
	params := make(map[string]string)
	// Load up parameters
	ldfnc := func(name string, val string) bool {
		if util.Blank(val) {
			return false
		}
		params[name] = val
		return true
	}
	havetype := ldfnc("type", spec.MainType)
	haveprj := ldfnc("project", spec.Project)
	havesub := ldfnc("subsystem", spec.Subsystem)
	havept := ldfnc("parttype", spec.PartType)
	havecat := ldfnc("category", spec.Category)
	havevendor := ldfnc("vendor", spec.Vendor)
	havevendorpn := ldfnc("vendorpn", spec.VendorPN)
	haveweblink := ldfnc("weblink", spec.WebLink)
	ldfnc("after", spec.DateAfter)
	ldfnc("before", spec.DateBefore)
	ldfnc("description", spec.Description)
	ldfnc("designer", spec.Designer)
	bMustBeEpic := false
	bMustBeSupplier := false

	if havetype {
		stype := strings.ToLower(params["type"])
		if stype != "epic" && stype != "supplier" {
			return nil, "", fmt.Errorf("Illegal parts type specified (%q). Must be either 'epic' or 'supplier'.\n", stype)
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
		params["type"] = stype
	}

	if haveprj || havesub || havept {
		bMustBeEpic = true
	}
	if havecat || havevendor || havevendorpn || haveweblink {
		bMustBeSupplier = true
	}

	if bMustBeEpic && bMustBeSupplier {
		return nil, "", fmt.Errorf("Incompatible parameters -- no parts can be found.\n")
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
		blst := make([]*CPart, 0, len(epiclst)+len(suplst))
		for _, p := range epiclst {
			part := &CPart{}
			// Generic Info
			part.PartNumber = p.PNString()
			part.Designer = p.Designer
			part.Description = p.Description
			part.DateIssued = p.DateIssued
			part.SequenceNum = p.SequenceNum
			part.PID = p.PID.String()
			// Epic Info
			part.ProjectId = p.ProjectId
			part.SubsystemId = p.SubsystemId
			part.PartType = p.PartType
			blst = append(blst, part)
		}
		for _, p := range suplst {
			part := &CPart{}
			// Generic Info
			part.PartNumber = p.PNString()
			part.Designer = p.Designer
			part.Description = p.Description
			part.DateIssued = p.DateIssued
			part.SequenceNum = p.SequenceNum
			part.PID = p.PID.String()
			// Supplier Info
			part.Category = p.Category
			part.Vendor = p.Vendor
			part.VendorPN = p.VendorPN
			part.WebLink = p.WebLink
			blst = append(blst, part)
		}
		sorter := func(i, j int) bool {
			return blst[i].DateIssued.Before(blst[j].DateIssued)
		}
		sort.Slice(blst, sorter)
		return blst, "both", nil
	} else if len(epiclst) > 0 {
		blst := make([]*CPart, 0, len(epiclst))
		for _, p := range epiclst {
			part := &CPart{}
			// Generic Info
			part.PartNumber = p.PNString()
			part.Designer = p.Designer
			part.Description = p.Description
			part.DateIssued = p.DateIssued
			part.SequenceNum = p.SequenceNum
			part.PID = p.PID.String()
			// Epic Info
			part.ProjectId = p.ProjectId
			part.SubsystemId = p.SubsystemId
			part.PartType = p.PartType
			blst = append(blst, part)
		}
		sorter := func(i, j int) bool {
			return blst[i].DateIssued.Before(blst[j].DateIssued)
		}
		sort.Slice(blst, sorter)
		return blst, "epic", nil
	} else if len(suplst) > 0 {
		blst := make([]*CPart, 0, len(suplst))
		for _, p := range suplst {
			part := &CPart{}
			// Generic Info
			part.PartNumber = p.PNString()
			part.Designer = p.Designer
			part.Description = p.Description
			part.DateIssued = p.DateIssued
			part.SequenceNum = p.SequenceNum
			part.PID = p.PID.String()
			// Supplier Info
			part.Category = p.Category
			part.Vendor = p.Vendor
			part.VendorPN = p.VendorPN
			part.WebLink = p.WebLink
			blst = append(blst, part)
		}
		sorter := func(i, j int) bool {
			return blst[i].DateIssued.Before(blst[j].DateIssued)
		}
		sort.Slice(blst, sorter)
		return blst, "supplier", nil
	} else {
		return []*CPart{}, "both", nil
	}
}
