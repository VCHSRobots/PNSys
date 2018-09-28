// --------------------------------------------------------------------
// help.go -- Help page
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

import (
	"github.com/gin-gonic/gin"
)

type HelpData struct {
	HeaderData
	HelpInfo string
}

func init() {
	RegisterPage("Help", Invoke_GET, handle_help)
}

func handle_help(c *gin.Context) {
	// Dummy data for now.
	data := &HelpData{}
	data.PageTabTitle = "Epic PN"
	data.IsLoggedIn = true
	data.UserFormattedName = "D. Brandon"
	data.IsAdmin = true
	data.PageTitle = "Help Page"
	data.Instructions = ""
	data.HelpInfo = "Sorry, not much help yet."

	html, err := MakePage(data, "header", "nav", "help", "footer")
	if err != nil {
		// Log has already been writen to...
		c.AbortWithError(400, err)
		return
	}

	c.Data(200, "text/html", html)
}
