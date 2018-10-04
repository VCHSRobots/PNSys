// --------------------------------------------------------------------
// cmd_list_spcat.go -- List supplier categories
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
)

func init() {
	RegistorCmd("list-spcat", "", "Lists the supplier categories.", handle_list_supplier_cats)
}

func handle_list_supplier_cats(c *Context, cmdline string) {
	lst := pnsql.GetSupplierCategories()
	tbl := util.NewTable("Category", "Description")
	for _, c := range lst {
		tbl.AddRow(c.Category, c.Description)
	}
	c.Printf("\n%s\n", tbl.Text())
}
