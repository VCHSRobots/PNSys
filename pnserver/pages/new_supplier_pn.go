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
    "github.com/gin-gonic/gin"
    "strings"
    "time"
)

type SupplierPNData struct {
    *HeaderData
    DesignersJson string
    Categories     []string
    KnownVendorsJaon string
    DesignerHint string
    VendorHint string
    CategoryHint string
}

type SupplierPNDataPost struct {
    *HeaderData
    NewPN          string
    Category       string
    Designer       string
    SequenceNumber string
    Description    string
    DateCreated    string
    ErrorMsg       string
}

type TSupplierSubmitData struct {
    Category    string `form:"Category"`
    Vendor      string `form:"Vendor"`
    VendorPN    string `form:"VendorPN"`
    WebLink     string `form:"WebLink"`
    Designer    string `form:"Designer"`
    Description string `form:"Description"`
}

func init() {
    RegisterPage("/NewSupplierPN", Invoke_GET, authorizer, handle_new_supplier_pn)
    RegisterPage("/SubmitNewSupplierPN", Invoke_POST, authorizer, handle_new_supplier_pn_post)
}

func handle_new_supplier_pn(c *gin.Context) {
    data := &EpicPNData{}
    data.HeaderData = GetHeaderData(c)
    data.PageTitle = "Create New Supplier Part Number"
    data.Instructions = ""
    data.StyleSheets = []string{"new_supplier_pn"}
    data.OnLoadFuncJS = "startUp"

    DesignersJson string
    Categories     []string
    KnownVendors []string
    DesignerHint string
    VendorHint string
    CategoryHint string


    data.PartTypes = pnsql.GetPartTypeSelStrings()

    des := pnsql.GetDesigners()
    des_bytes, err := json.MarshalIndent(des, "", "  ")
    if err != nil {
        log.Errorf("Unable to convert to json. Err=%v", err)
        c.AbortWithError(400, err)
        return
    }
    data.DesignersJson = string(des_bytes)

    SendPage(c, data, "header", "nav", "newpn_menubar", "new_epic_pn", "footer")
}

func handle_new_supplier_pn_post(c *gin.Context) {
    data := &EpicPNDataPost{}
    data.HeaderData = GetHeaderData(c)

    var sdata TEpicSubmitData
    err := c.ShouldBind(&sdata)
    if err != nil {
        log.Errorf("Unable to bind data. Err=%v", err)
        c.AbortWithError(400, err)
        return
    }

    data.PageTitle = "New Epic Part Number"
    data.StyleSheets = []string{}
    data.Project = sdata.Project
    data.Subsystem = sdata.Subsystem
    data.Designer = sdata.Designer
    data.PartType = sdata.PartType
    data.SequenceNumber = "087"
    data.Description = sdata.Description
    data.DateCreated = time.Now().Format("2006-01-02")
    data.NewPN = getword(sdata.Project) + "-" + getword(sdata.Subsystem) +
        "-" + getword(sdata.PartType) + data.SequenceNumber

    SendPage(c, data, "header", "nav", "newpn_menubar", "new_epic_pn_post", "footer")
}

func getword(x string) string {
    wrds := strings.Split(x, " ")
    if len(wrds) <= 0 {
        return x
    }
    return strings.TrimSpace(wrds[0])
}
