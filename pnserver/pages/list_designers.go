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
	RegisterPage("ListDesigners", Invoke_GET, authorizer, handle_list_designers)
}

func handle_list_designers(c *gin.Context) {
	rows := pnsql.GetDesigners()
	data := &TablePageData{}
	data.HeaderData = GetHeaderData(c)
	data.Head = []string{"Name", "Year0", "Active"}
	data.Rows = make([]TColumn, 0, len(rows))
	for _, r := range rows {
		sactive := ""
		if r.Active {
			sactive = "Yes"
		}
		data.Rows = append(data.Rows, TColumn{Cols: []string{r.Name, r.Year0, sactive}})
	}
	SendPage(c, data, "header", "nav", "tablepage", "footer")
}
