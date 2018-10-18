// --------------------------------------------------------------------
// new_epic_pn.go -- Shows Epic PN input form
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

import (
	"encoding/base64"
	"epic/lib/log"
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
	"github.com/gin-gonic/gin"
)

type EpicPageDefaults struct {
	Designer          string
	Project           string
	Subsystem         string
	PartType          string
	Description       string
	DescriptionBase64 string
}

type EpicPNData struct {
	*HeaderData
	*SelectionBoxData
	Defaults *EpicPageDefaults
}

func init() {
	RegisterPage("/NewEpicPN", Invoke_GET, authorizer, handle_new_epic_pn)
	RegisterPage("/SubmitNewEpicPN", Invoke_POST, authorizer, handle_new_epic_pn_post)
}

func handle_new_epic_pn(c *gin.Context) {
	handle_new_epic_pn_with_error(c, "")
}

func handle_new_epic_pn_with_error(c *gin.Context, errmsg string) {
	data := &EpicPNData{}
	data.HeaderData = GetHeaderData(c)
	data.PageTitle = "Create New Epic Part Number"
	data.Instructions = ""
	data.StyleSheets = []string{"new_epic_pn"}
	data.OnLoadFuncJS = "startUp"
	data.ErrorMessage = errmsg

	var err error
	data.SelectionBoxData, err = GetSelectionBoxData()
	if err != nil {
		SendErrorPage(c, err)
		return
	}

	var sd *EpicPageDefaults
	ses := GetSession(c)
	t, ok := ses.Data["EpicPageDefaults"]
	if !ok {
		sd = &EpicPageDefaults{}
	} else {
		sd, ok = t.(*EpicPageDefaults)
		if !ok {
			log.Errorf("Unable to type convert EpicPageDefaults in handle_new_epic_pn.")
			sd = &EpicPageDefaults{}
		}
	}
	if util.Blank(sd.Designer) {
		sd.Designer = data.HeaderData.Designer
	}
	data.Defaults = sd
	txt := data.Defaults.Description
	data.Defaults.DescriptionBase64 = base64.StdEncoding.EncodeToString([]byte(txt))
	SendPage(c, data, "header", "menubar", "new_epic_pn", "footer")
}

func handle_new_epic_pn_post(c *gin.Context) {

	type submitdata struct {
		Project     string `form:"Project"`
		Subsystem   string `form:"Subsystem"`
		Designer    string `form:"Designer"`
		PartType    string `form:"PartType"`
		Description string `form:"Description"`
	}

	var data submitdata
	err := c.ShouldBind(&data)
	if err != nil {
		SendErrorPage(c, err)
		return
	}

	sd := &EpicPageDefaults{}
	sd.Designer = data.Designer
	sd.Project = data.Project
	sd.Subsystem = data.Subsystem
	sd.PartType = data.PartType
	sd.Description = data.Description
	ses := GetSession(c)
	ses.Data["EpicPageDefaults"] = sd

	if util.Blank(data.Subsystem) {
		handle_new_epic_pn_with_error(c, "No Subsystems for the project yet. Have an admin create one.")
		return
	}

	if util.Blank(data.Description) {
		handle_new_epic_pn_with_error(c, "Please provide at least a few words that describes the part.")
		return
	}

	if !HasWritePrivilege(c) {
		handle_new_epic_pn_with_error(c, "Sorry, users logged in as guest cannot add parts.")
		return
	}

	pn, err := pnsql.NewEpicPartInSequence(data.Designer, data.Project, data.Subsystem, data.PartType,
		data.Description)
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
	ses.Data["EpicPageDefaults"] = sd

	url := fmt.Sprintf("/ShowPart?pn=%s", pn.PNString())
	c.Redirect(303, url)

	// msg := fmt.Sprintf("Submitted Data: <br>")
	// msg += fmt.Sprintf("Project     = %s<br>", data.Project)
	// msg += fmt.Sprintf("Subsystem   = %s<br>", data.Subsystem)
	// msg += fmt.Sprintf("Designer    = %s<br>", data.Designer)
	// msg += fmt.Sprintf("PartType    = %s<br>", data.PartType)
	// msg += fmt.Sprintf("Description = %s<br>", data.Description)

	// SendMessagePagef(c, msg)
}
