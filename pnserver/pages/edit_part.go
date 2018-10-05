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

	pn := c.Query("pn")
	if util.Blank(pn) {
		SendErrorPagef(c, "EditPart page called without a part number!")
		return
	}
	handle_edit_part_with_error(c, pn, "")
}

func handle_edit_part_with_error(c *gin.Context, pn, errmsg string) {
	data := &EditPartData{}
	data.HeaderData = GetHeaderData(c)
	data.PageTitle = "EDIT Part Number"
	data.ErrorMessage = errmsg
	data.OnLoadFuncJS = "startUp"

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

	part, err := pnsql.GetEpicPart(data.PartNumber)
	if err != nil {
		SendErrorPage(c, err)
		return
	}
	if part == nil {
		err = fmt.Errorf("Part %s not found in database on Edit Part, but should be here!", data.PartNumber)
		SendErrorPage(c, err)
		return
	}
	nchanges := 0
	if part.Designer != data.Designer {
		err = pnsql.SetEpicPartDesigner(part, data.Designer)
		if err != nil {
			SendErrorPage(c, err)
			return
		}
		log.Infof("Part %s edited by %s. Designer changed to %s.", part.PNString(), GetDesigner(c), data.Designer)
		nchanges++
	}
	if part.Description != data.Description {
		if util.Blank(data.Description) {
			handle_edit_part_with_error(c, part.PNString(), "Description cannot be blank.")
			return
		}
		err = pnsql.SetEpicPartDescription(part, data.Description)
		if err != nil {
			SendErrorPage(c, err)
			return
		}
		log.Infof("Part %s edited by %s. Description changed to %q.", part.PNString(), GetDesigner(c), data.Description)
		nchanges++
	}
	if nchanges <= 0 {
		handle_edit_part_with_error(c, part.PNString(), "No changes made.. ?")
		return
	}
	url := fmt.Sprintf("/ShowPart?pn=%s", part.PNString())
	c.Redirect(303, url)

	// msg := fmt.Sprintf("Submitted Data:<br>")
	// msg += fmt.Sprintf("PartNumber  = %s<br>", data.PartNumber)
	// msg += fmt.Sprintf("Designer    = %s<br>", data.Designer)
	// msg += fmt.Sprintf("Description = %s<br>", data.Description)
	// SendMessagePagef(c, msg)
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

	part, err := pnsql.GetSupplierPart(data.PartNumber)
	if err != nil {
		SendErrorPage(c, err)
		return
	}
	if part == nil {
		err = fmt.Errorf("Part %s not found in database on Edit Part, but should be here!", data.PartNumber)
		SendErrorPage(c, err)
		return
	}
	nchanges := 0
	if part.Designer != data.Designer {
		err = pnsql.SetSupplierPartDesigner(part, data.Designer)
		if err != nil {
			SendErrorPage(c, err)
			return
		}
		log.Infof("Part %s edited by %s. Designer changed to %s.", part.PNString(), GetDesigner(c), data.Designer)
		nchanges++
	}
	if part.Vendor != data.Vendor {
		err = pnsql.SetSupplierPartVendor(part, data.Vendor)
		if err != nil {
			SendErrorPage(c, err)
			return
		}
		log.Infof("Part %s edited by %s. Vendor changed to %s.", part.PNString(), GetDesigner(c), data.Vendor)
		nchanges++
	}
	if part.VendorPN != data.VendorPN {
		err = pnsql.SetSupplierPartVendorPN(part, data.VendorPN)
		if err != nil {
			SendErrorPage(c, err)
			return
		}
		log.Infof("Part %s edited by %s. Vendor PN changed to %s.", part.PNString(), GetDesigner(c), data.VendorPN)
		nchanges++
	}
	if part.WebLink != data.WebLink {
		err = pnsql.SetSupplierPartWebLink(part, data.WebLink)
		if err != nil {
			SendErrorPage(c, err)
			return
		}
		log.Infof("Part %s edited by %s. WebLink changed to %q.", part.PNString(), GetDesigner(c), data.WebLink)
		nchanges++
	}
	if part.Description != data.Description {
		if util.Blank(data.Description) {
			handle_edit_part_with_error(c, part.PNString(), "Description cannot be blank.")
			return
		}
		err = pnsql.SetSupplierPartDescription(part, data.Description)
		if err != nil {
			SendErrorPage(c, err)
			return
		}
		log.Infof("Part %s edited by %s. Description changed to %q.", part.PNString(), GetDesigner(c), data.Description)
		nchanges++
	}
	if nchanges <= 0 {
		handle_edit_part_with_error(c, part.PNString(), "No changes made.. ?")
		return
	}
	url := fmt.Sprintf("/ShowPart?pn=%s", part.PNString())
	c.Redirect(303, url)

	// msg := fmt.Sprintf("Data:<br>")
	// msg += fmt.Sprintf("PartNumber  = %s<br>", data.PartNumber)
	// msg += fmt.Sprintf("Vendor      = %s<br>", data.Vendor)
	// msg += fmt.Sprintf("VendorPN    = %s<br>", data.VendorPN)
	// msg += fmt.Sprintf("WebLink     = %s<br>", data.WebLink)
	// msg += fmt.Sprintf("Designer    = %s<br>", data.Designer)
	// msg += fmt.Sprintf("Description = %s<br>", data.Description)
	// SendMessagePagef(c, msg)
}
