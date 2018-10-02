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
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type EpicPNData struct {
	*HeaderData
	ProjectsJson  string
	DesignersJson string
	PartTypesJson string
}

type EpicPNDataPost struct {
	*HeaderData
	PartNumber     string
	Project        string
	Subsystem      string
	Designer       string
	PartType       string
	SequenceNumber string
	Description    string
	DateCreated    string
	ErrorMsg       string
}

func init() {
	RegisterPage("/NewEpicPN", Invoke_GET, authorizer, handle_new_epic_pn)
	RegisterPage("/SubmitNewEpicPN", Invoke_POST, authorizer, handle_new_epic_pn_post)
}

func handle_new_epic_pn(c *gin.Context) {
	data := &EpicPNData{}
	data.HeaderData = GetHeaderData(c)
	data.PageTitle = "Create New Epic Part Number"
	data.Instructions = ""
	data.StyleSheets = []string{"new_epic_pn"}
	data.OnLoadFuncJS = "startUp"

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

	ptlst := pnsql.GetPartTypes()
	pt_bytes, err := json.MarshalIndent(ptlst, "", "  ")
	if err != nil {
		log.Errorf("Unable to convert to json. Err=%v", err)
		c.AbortWithError(400, err)
		return
	}
	data.PartTypesJson = string(pt_bytes)

	SendPage(c, data, "header", "menubar", "new_epic_pn", "footer")
}

func handle_new_epic_pn_post(c *gin.Context) {
	data := &EpicPNDataPost{}
	data.HeaderData = GetHeaderData(c)

	type submitdata struct {
		Project     string `form:"Project"`
		Subsystem   string `form:"Subsystem"`
		Designer    string `form:"Designer"`
		PartType    string `form:"PartType"`
		Description string `form:"Description"`
	}

	var sdata submitdata
	err := c.ShouldBind(&sdata)
	if err != nil {
		err = fmt.Errorf("Bind error for SubmitNewEpicPN. Err=%v", err)
		log.Errorf("%v", err)
		SendErrorPage(c, err)
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
	data.PartNumber = getword(sdata.Project) + "-" + getword(sdata.Subsystem) +
		"-" + getword(sdata.PartType) + data.SequenceNumber

	SendPage(c, data, "header", "menubar", "new_epic_pn_post", "footer")
}

func getword(x string) string {
	wrds := strings.Split(x, " ")
	if len(wrds) <= 0 {
		return x
	}
	return strings.TrimSpace(wrds[0])
}
