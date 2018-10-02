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
	RegisterPage("About", Invoke_GET, guest_auth, handle_about)
}

func handle_about(c *gin.Context) {
	data := GetHeaderData(c)
	data.HideAboutLink = true
	SendPage(c, data, "header", "menubar", "about", "footer")
}
