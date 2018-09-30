// --------------------------------------------------------------------
// help.go -- Help page
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

import (
	"github.com/gin-gonic/gin"
)

func init() {
	RegisterPage("Help", Invoke_GET, guest_auth, handle_help)
}

func handle_help(c *gin.Context) {
	data := GetHeaderData(c)
	SendPage(c, data, "header", "nav", "help", "footer")
}
