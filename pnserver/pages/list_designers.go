// --------------------------------------------------------------------
// list_designers.go -- List Designers Page
//
// Created 2018-09-22 DLB
// --------------------------------------------------------------------

package pages

import (
	"epic/pnserver/pnsql"
	"github.com/gin-gonic/gin"
)

func init() {
	RegisterPage("ListDesigners", Invoke_GET, handle_list_designers)
}

func handle_list_designers(c *gin.Context) {
	rows := pnsql.GetDesigners()
	data := TablePageData{}
	data.Head = []string{"Name", "Year0"}
	data.Rows = make([]TColumn, 0, len(rows))
	for _, r := range rows {
		data.Rows = append(data.Rows, TColumn{Cols: []string{r.Name, r.Year0}})
	}

	// Other dummy data for now.
	data.PageTabTitle = "Epic PN"
	data.PageTitle = "Designers"
	data.IsLoggedIn = true
	data.UserFormattedName = "D. Brandon"
	data.IsAdmin = true

	html, err := MakePage(data, "header", "nav", "tablepage", "footer")
	if err != nil {
		// Log has already been writen to...
		c.AbortWithError(400, err)
		return
	}

	c.Data(200, "text/html", html)
}
