// --------------------------------------------------------------------
// new_epic_pn.go -- Shows Epic PN input form
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

type EpicPNData struct {
	HeaderData
	ProjectsJson  string
	DesignersJson string
	PartTypes     []string
}

type EpicPNDataPost struct {
	HeaderData
	NewPN          string
	Project        string
	Subsystem      string
	Designer       string
	PartType       string
	SequenceNumber string
	Description    string
	DateCreated    string
	ErrorMsg       string
}

type TEpicSubmitData struct {
	Project     string `form:"Project"`
	Subsystem   string `form:"Subsystem"`
	Designer    string `form:"Designer"`
	PartType    string `form:"PartType"`
	Description string `form:"Description"`
}

func init() {
	RegisterPage("NewEpicPN", Invoke_GET, handle_new_epic_pn)
	RegisterPage("SubmitNewEpicPN", Invoke_POST, handle_new_epic_pn_post)
}

func handle_new_epic_pn(c *gin.Context) {
	// Dummy data for now.
	data := &EpicPNData{}
	data.PageTabTitle = "Epic PN"
	data.IsLoggedIn = true
	data.UserFormattedName = "D. Brandon"
	data.IsAdmin = true
	data.PageTitle = "Create New Epic Part Number"
	data.Instructions = "Fill form out and click Submit."
	data.StyleSheets = []string{"new_epic_pn"}
	data.OnLoadFuncJS = "startUp"
	data.PartTypes = pnsql.GetPartTypeSelStrings()

	var err error
	prjs := pnsql.GetProjects()
	prj_bytes, err := json.MarshalIndent(prjs, "", "  ")
	if err != nil {
		log.Errorf("Unable to convert to json. Err=%v", err)
		c.AbortWithError(400, err)
		return
	}
	data.ProjectsJson = string(prj_bytes)

	des := pnsql.GetDesigners()
	des_bytes, err := json.MarshalIndent(des, "", "  ")
	if err != nil {
		log.Errorf("Unable to convert to json. Err=%v", err)
		c.AbortWithError(400, err)
		return
	}
	data.DesignersJson = string(des_bytes)

	html, err := MakePage(data, "header", "nav", "newpn_menubar", "new_epic_pn", "footer")
	if err != nil {
		// Log has already been writen to...
		c.AbortWithError(400, err)
		return
	}

	c.Data(200, "text/html", html)
}

func getword(x string) string {
	wrds := strings.Split(x, " ")
	if len(wrds) <= 0 {
		return x
	}
	return strings.TrimSpace(wrds[0])
}

func handle_new_epic_pn_post(c *gin.Context) {
	var sdata TEpicSubmitData
	err := c.ShouldBind(&sdata)
	if err != nil {
		log.Errorf("Unable to bind data. Err=%v", err)
		c.AbortWithError(400, err)
		return
	}

	// Dummy data for now.
	data := &EpicPNDataPost{}
	data.PageTabTitle = "Epic PN"
	data.IsLoggedIn = true
	data.UserFormattedName = "D. Brandon"
	data.IsAdmin = true
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

	html, err := MakePage(data, "header", "nav", "newpn_menubar", "new_epic_pn_post", "footer")
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.Data(200, "text/html", html)
}
