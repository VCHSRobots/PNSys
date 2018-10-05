// --------------------------------------------------------------------
// new_supplier_pn.go -- Shows Supplier PN input form
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

import (
	"epic/lib/log"
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SupplierPageDefaults struct {
	Designer    string
	Category    string
	Vendor      string
	VendorPN    string
	WebLink     string
	Description string
}

type SupplierPNData struct {
	*HeaderData
	*SelectionBoxData
	Defaults *SupplierPageDefaults
}

func init() {
	RegisterPage("/NewSupplierPN", Invoke_GET, authorizer, handle_new_supplier_pn)
	RegisterPage("/SubmitNewSupplierPN", Invoke_POST, authorizer, handle_new_supplier_pn_post)
}

func handle_new_supplier_pn(c *gin.Context) {
	handle_new_supplier_pn_with_error(c, "")
}

func handle_new_supplier_pn_with_error(c *gin.Context, errmsg string) {

	data := &SupplierPNData{}
	data.HeaderData = GetHeaderData(c)
	data.PageTitle = "Create New Supplier Part Number"
	data.Instructions = ""
	data.StyleSheets = []string{"new_supplier_pn"}
	data.OnLoadFuncJS = "startUp"
	data.ErrorMessage = errmsg

	var err error
	data.SelectionBoxData, err = GetSelectionBoxData()
	if err != nil {
		SendErrorPage(c, err)
		return
	}

	var sd *SupplierPageDefaults
	ses := GetSession(c)
	t, ok := ses.Data["SupplierPageDefaults"]
	if !ok {
		sd = &SupplierPageDefaults{}
	} else {
		sd, ok = t.(*SupplierPageDefaults)
		if !ok {
			log.Errorf("Unable to type convert SupplierPageDefaults in handle_new_supplier_pn.")
			sd = &SupplierPageDefaults{}
		}
	}
	if util.Blank(sd.Designer) {
		sd.Designer = data.HeaderData.Designer
	}
	data.Defaults = sd
	SendPage(c, data, "header", "menubar", "new_supplier_pn", "footer")
}

func handle_new_supplier_pn_post(c *gin.Context) {

	type submitdata struct {
		Category    string `form:"Category"`
		Vendor      string `form:"Vendor"`
		VendorPN    string `form:"VendorPN"`
		WebLink     string `form:"WebLink"`
		Designer    string `form:"Designer"`
		Description string `form:"Description"`
	}

	var data submitdata
	err := c.ShouldBind(&data)
	if err != nil {
		err = fmt.Errorf("Bind error for SubmitNewSupplierPN. Err=%v", err)
		SendErrorPage(c, err)
		return
	}

	// Remember the defaults...
	sd := &SupplierPageDefaults{}
	sd.Designer = data.Designer
	sd.Category = data.Category
	sd.Vendor = data.Vendor
	sd.VendorPN = data.VendorPN
	sd.WebLink = data.WebLink
	sd.Description = data.Description
	ses := GetSession(c)
	ses.Data["SupplierPageDefaults"] = sd

	if util.Blank(data.Description) {
		handle_new_supplier_pn_with_error(c, "Please provide at least a few words that describes the part.")
		return
	}

	pn, err := pnsql.NewSupplierPartInSequence(data.Designer, data.Category, data.Vendor,
		data.VendorPN, data.WebLink, data.Description)
	if err != nil {
		SendErrorPage(c, err)
		return
	}

	user := GetDesigner(c)
	if user == pn.Designer {
		log.Infof("New part %s designed and and added to database by %s.", pn.PNString(), pn.Designer)
	} else {
		log.Infof("New part %s designed by %s and added to database by %s.", pn.PNString(), pn.Designer, user)
	}

	// Defaults for next time...
	sd.Description = ""
	sd.Designer = ses.Name
	ses.Data["SupplierPageDefaults"] = sd

	url := fmt.Sprintf("/ShowPart?pn=%s", pn.PNString())
	c.Redirect(303, url)

	// msg := fmt.Sprintf("Submitted Data: <br>")
	// msg += fmt.Sprintf("Category     = %s<br>", data.Category)
	// msg += fmt.Sprintf("Vendor       = %s<br>", data.Vendor)
	// msg += fmt.Sprintf("VendorPN     = %s<br>", data.VendorPN)
	// msg += fmt.Sprintf("WebLink      = %s<br>", data.WebLink)
	// msg += fmt.Sprintf("Designer     = %s<br>", data.Designer)
	// msg += fmt.Sprintf("Description  = %s<br>", data.Description)
	// SendMessagePagef(c, msg)
}
