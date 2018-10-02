// --------------------------------------------------------------------
// new_supplier_pn.go -- Shows Supplier PN input form
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

import (
	"encoding/json"
	"epic/lib/log"
	"epic/pnserver/pnsql"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type SupplierPNData struct {
	*HeaderData
	DesignersJson    string
	CategoriesJson   string
	KnownVendorsJson string
	DesignerHint     string
	VendorHint       string
	CategoryHint     string
}

type SupplierPNDataPost struct {
	*HeaderData
	NewPN          string
	Category       string
	Vendor         string
	VendorPN       string
	WebLink        string
	WebLinkAddr    string
	Designer       string
	SequenceNumber string
	Description    string
	DateCreated    string
	ErrorMsg       string
}

func init() {
	RegisterPage("/NewSupplierPN", Invoke_GET, authorizer, handle_new_supplier_pn)
	RegisterPage("/SubmitNewSupplierPN", Invoke_POST, authorizer, handle_new_supplier_pn_post)
}

func handle_new_supplier_pn(c *gin.Context) {
	data := &SupplierPNData{}
	data.HeaderData = GetHeaderData(c)
	data.PageTitle = "Create New Supplier Part Number"
	data.Instructions = ""
	data.StyleSheets = []string{"new_supplier_pn"}
	data.OnLoadFuncJS = "startUp"

	des := pnsql.GetDesigners()
	des_bytes, err := json.MarshalIndent(des, "", "  ")
	if err != nil {
		log.Errorf("Unable to convert to json. Err=%v", err)
		c.AbortWithError(400, err)
		return
	}
	data.DesignersJson = string(des_bytes)

	catlst := pnsql.GetSupplierCategories()
	cat_bytes, err := json.MarshalIndent(catlst, "", "  ")
	if err != nil {
		log.Errorf("Unable to convert to json. Err=%v", err)
		c.AbortWithError(400, err)
		return
	}
	data.CategoriesJson = string(cat_bytes)

	vendorlst := pnsql.GetVendors()
	vbytes, err := json.MarshalIndent(vendorlst, "", "  ")
	if err != nil {
		log.Errorf("Unable to convert to json. Err=%v", err)
		c.AbortWithError(400, err)
		return
	}
	data.KnownVendorsJson = string(vbytes)

	ses := GetSession(c)
	data.DesignerHint = ses.Name

	SendPage(c, data, "header", "menubar", "new_supplier_pn", "footer")
}

func handle_new_supplier_pn_post(c *gin.Context) {
	data := &SupplierPNDataPost{}
	data.HeaderData = GetHeaderData(c)

	type submitdata struct {
		Category    string `form:"Category"`
		Vendor      string `form:"Vendor"`
		VendorPN    string `form:"VendorPN"`
		WebLink     string `form:"WebLink"`
		Designer    string `form:"Designer"`
		Description string `form:"Description"`
	}

	var sdata submitdata
	err := c.ShouldBind(&sdata)
	if err != nil {
		err = fmt.Errorf("Bind error for SubmitNewSupplierPN. Err=%v", err)
		log.Errorf("%v", err)
		SendErrorPage(c, err)
		return
	}

	data.PageTitle = "New Supplier Part Number"
	data.StyleSheets = []string{}
	data.Category = sdata.Category
	data.Vendor = sdata.Vendor
	data.VendorPN = sdata.VendorPN
	data.WebLink = sdata.WebLink
	data.WebLinkAddr = FixWebLinkAddr(sdata.WebLink)
	data.Designer = sdata.Designer
	data.SequenceNumber = "087"
	data.Description = sdata.Description
	data.DateCreated = time.Now().Format("2006-01-02")
	data.NewPN = "SP-" + data.Category + "-" + data.SequenceNumber
	SendPage(c, data, "header", "menubar", "new_supplier_pn_post", "footer")
}
