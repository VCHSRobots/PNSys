// --------------------------------------------------------------------
// about.go -- Help page
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

import (
	"github.com/gin-gonic/gin"
)

func init() {
	RegisterPage("About", Invoke_GET, handle_about)
}

func handle_about(c *gin.Context) {
	// Dummy data for now.
	data := &HeaderData{}
	data.PageTabTitle = "Epic PN"
	data.IsLoggedIn = true
	data.UserFormattedName = "D. Brandon"
	data.IsAdmin = true
	data.PageTitle = "About..."
	data.Instructions = ""

	html, err := MakePage(data, "header", "nav", "about", "footer")
	if err != nil {
		// Log has already been writen to...
		c.AbortWithError(400, err)
		return
	}

	c.Data(200, "text/html", html)
}
