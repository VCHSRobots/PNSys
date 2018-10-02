// --------------------------------------------------------------------
// show_part.go -- Shows Epic Part Number
//
// Created 2018-09-31 DLB
// --------------------------------------------------------------------

package pages

import (
	"epic/lib/log"
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ShowPNData struct {
	*HeaderData
	IsEpic     bool
	IsSupplier bool
	HavePart   bool
	// Generic Info
	PartNumber     string
	Description    string
	Designer       string
	DateIssued     string
	PID            string
	SequenceNumber string
	// For Epic Parts
	Project   string
	Subsystem string
	PartType  string
	// For Supplier Parts
	Category    string
	Vendor      string
	VendorPN    string
	WebLink     string
	WebLinkAddr string
}

func init() {
	RegisterPage("/ShowPart", Invoke_GET, authorizer, handle_show_part)
}

func handle_show_part(c *gin.Context) {
	data := &ShowPNData{}
	data.HeaderData = GetHeaderData(c)
	data.PageTitle = "Part Number"
	data.StyleSheets = []string{"show_part"}
	pn := c.Query("pn")
	if util.Blank(pn) {
		log.Warnf("ShowPart page called without a part number!")
		c.Redirect(300, "/NewEpicPN")
		return
	}
	ty, _ := pnsql.ClassifyPN(pn)
	if ty == pnsql.PNType_Epic {
		data.IsEpic = true
		part, err := pnsql.GetEpicPart(pn)
		if err != nil {
			err = fmt.Errorf("Database error while searching for part %s. Err=%v", pn, err)
			log.Warnf("%v", err)
			data.ErrorMessage = fmt.Sprintf("%v", err)
			SendPage(c, data, "header", "menubar", "error", "footer")
			return
		}
		if part == nil {
			data.ErrorMessage = fmt.Sprintf("Part %s not in database.", pn)
			SendPage(c, data, "header", "menubar", "error", "footer")
			return
		}
		data.HavePart = true
		data.PartNumber = part.PNString()
		data.Description = part.Description
		data.Designer = part.Designer
		data.DateIssued = part.DateIssued.Format("2006-01-02")
		data.PID = part.PID.String()
		data.SequenceNumber = fmt.Sprintf("%03d", part.SequenceNum)
		data.Project = part.ProjectDesc()
		data.Subsystem = part.SubsystemDesc()
		data.PartType = part.PartTypeDesc()
		SendPage(c, data, "header", "menubar", "show_part", "footer")
		return

	} else if ty == pnsql.PNType_Supplier {
		data.IsSupplier = true
		part, err := pnsql.GetSupplierPart(pn)
		if err != nil {
			err = fmt.Errorf("Database error while searching for part %s. Err=%v", pn, err)
			log.Warnf("%v", err)
			data.ErrorMessage = fmt.Sprintf("%v", err)
			SendPage(c, data, "header", "menubar", "error", "footer")
			return
		}
		if part == nil {
			data.ErrorMessage = fmt.Sprintf("Part %s not in database.", pn)
			SendPage(c, data, "header", "menubar", "error", "footer")
			return
		}
		data.HavePart = true
		data.PartNumber = part.PNString()
		data.Description = part.Description
		data.Designer = part.Designer
		data.DateIssued = part.DateIssued.Format("2006-01-02")
		data.PID = part.PID.String()
		data.SequenceNumber = fmt.Sprintf("%03d", part.SequenceNum)
		data.Category = part.CategoryDesc()
		data.Vendor = part.Vendor
		data.VendorPN = part.VendorPN
		data.WebLink = part.WebLink
		data.WebLinkAddr = FixWebLinkAddr(part.WebLink)
		SendPage(c, data, "header", "menubar", "show_part", "footer")
		return

	} else {
		data.HavePart = false
		data.ErrorMessage = fmt.Sprintf("The part %q is not a known part.", pn)
		SendPage(c, data, "header", "menubar", "error", "footer")
		return
	}
}
