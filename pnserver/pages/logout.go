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
	RegisterPage("Logout", Invoke_GET, authorizer, handle_logout)
}

func handle_logout(c *gin.Context) {
	kill_session(c)
	data := &HeaderData{}
	data.PageTabTitle = "Epic PN"
	data.HideLoginLink = false
	data.HideAboutLink = true
	SendPage(c, data, "header", "logout", "footer")
}
