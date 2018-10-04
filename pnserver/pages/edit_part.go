// --------------------------------------------------------------------
// edit_part.go -- Edits a Part Number
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

type EditPartData struct {
	*HeaderData
	*SelectionBoxData

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
	RegisterPage("/EditPart", Invoke_GET, authorizer, handle_edit_part)
	RegisterPage("/EditEpicPNPost", Invoke_POST, authorizer, handle_edit_epic_post)
	RegisterPage("/EditSupplierPNPost", Invoke_POST, authorizer, handle_edit_supplier_post)
}

func handle_edit_part(c *gin.Context) {
	data := &EditPartData{}
	data.HeaderData = GetHeaderData(c)
	data.PageTitle = "EDIT Part Number"

	data.OnLoadFuncJS = "startUp"
	pn := c.Query("pn")
	if util.Blank(pn) {
		log.Warnf("EditPart page called without a part number!")
		c.Redirect(300, "/NewEpicPN")
		return
	}

	var err error
	data.SelectionBoxData, err = GetSelectionBoxData()
	if err != nil {
		SendErrorPage(c, err)
		return
	}

	ty, _ := pnsql.ClassifyPN(pn)
	if ty == pnsql.PNType_Epic {
		data.IsEpic = true
		data.StyleSheets = []string{"edit_epic_part"}
		part, err := pnsql.GetEpicPart(pn)
		if err != nil {
			err = fmt.Errorf("Database error while searching for part %s. Err=%v", pn, err)
			SendErrorPage(c, err)
			return
		}
		if part == nil {
			err = fmt.Errorf("EditPart called with a bogus part number! (%s)", pn)
			SendErrorPage(c, err)
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
		SendPage(c, data, "header", "menubar", "edit_epic_part", "footer")
		return

	} else if ty == pnsql.PNType_Supplier {
		data.IsSupplier = true
		data.StyleSheets = []string{"edit_supplier_part"}
		part, err := pnsql.GetSupplierPart(pn)
		if err != nil {
			err = fmt.Errorf("Database error while searching for part %s. Err=%v", pn, err)
			SendErrorPage(c, err)
			return
		}
		if part == nil {
			SendErrorPagef(c, "EditPart called with a bogus part number! (%s)", pn)
			return
		}
		data.HavePart = true
		data.PartNumber = part.PNString()
		data.Description = util.CleanForWeb(part.Description)
		data.Designer = part.Designer
		data.DateIssued = part.DateIssued.Format("2006-01-02")
		data.PID = part.PID.String()
		data.SequenceNumber = fmt.Sprintf("%03d", part.SequenceNum)
		data.Category = util.CleanForWeb(part.CategoryDesc())
		data.Vendor = util.CleanForWeb(part.Vendor)
		data.VendorPN = util.CleanForWeb(part.VendorPN)
		data.WebLink = util.CleanForWeb(part.WebLink)
		data.WebLinkAddr = FixWebLinkAddr(data.WebLink)

		SendPage(c, data, "header", "menubar", "edit_supplier_part", "footer")
		return
	} else {
		SendErrorPagef(c, "EditPart called with a bogus part number! (%s)", pn)
		return
	}
}

func handle_edit_epic_post(c *gin.Context) {
	type submitdata struct {
		PartNumber  string `form:"PartNumber"`
		Designer    string `form:"Designer"`
		Description string `form:"Description"`
	}
	var data submitdata
	err := c.ShouldBind(&data)
	if err != nil {
		err = fmt.Errorf("Bind error for EditEpicPNPost. Err=%v", err)
		log.Errorf("%v", err)
		SendErrorPage(c, err)
		return
	}
	msg := fmt.Sprintf("Submitted Data:<br>")
	msg += fmt.Sprintf("PartNumber  = %s<br>", data.PartNumber)
	msg += fmt.Sprintf("Designer    = %s<br>", data.Designer)
	msg += fmt.Sprintf("Description = %s<br>", data.Description)
	SendMessagePagef(c, msg)
}

func handle_edit_supplier_post(c *gin.Context) {
	type submitdata struct {
		PartNumber  string `form:"PartNumber"`
		Vendor      string `form:"Vendor"`
		VendorPN    string `form:"VendorPN"`
		WebLink     string `form:"WebLink"`
		Designer    string `form:"Designer"`
		Description string `form:"Description"`
	}
	var data submitdata
	err := c.ShouldBind(&data)
	if err != nil {
		err = fmt.Errorf("Bind error for EditSupplierPNPost. Err=%v", err)
		SendErrorPage(c, err)
		return
	}
	msg := fmt.Sprintf("Data:<br>")
	msg += fmt.Sprintf("PartNumber  = %s<br>", data.PartNumber)
	msg += fmt.Sprintf("Vendor      = %s<br>", data.Vendor)
	msg += fmt.Sprintf("VendorPN    = %s<br>", data.VendorPN)
	msg += fmt.Sprintf("WebLink     = %s<br>", data.WebLink)
	msg += fmt.Sprintf("Designer    = %s<br>", data.Designer)
	msg += fmt.Sprintf("Description = %s<br>", data.Description)
	SendMessagePagef(c, msg)
}
