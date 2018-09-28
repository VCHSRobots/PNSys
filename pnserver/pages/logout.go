// --------------------------------------------------------------------
// logout.go -- Logout Page
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

import (
	"github.com/gin-gonic/gin"
)

func init() {
	RegisterPage("Logout", Invoke_GET, handle_logout)
}

func handle_logout(c *gin.Context) {
	// Dummy data for now.
	data := &HeaderData{}
	data.PageTabTitle = "Epic PN"
	data.IsLoggedIn = false
	data.IsAdmin = false

	html, err := MakePage(data, "header", "logout", "footer")
	if err != nil {
		// Log has already been writen to...
		c.AbortWithError(400, err)
		return
	}

	c.Data(200, "text/html", html)
}
