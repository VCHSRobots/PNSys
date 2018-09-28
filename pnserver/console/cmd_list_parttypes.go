// --------------------------------------------------------------------
// cmd_list_parttypes.go -- Show part types
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
)

func init() {
	RegistorCmd("list-parttypes", "", "Lists all the part types for epic part numbers.", handle_list_parttypes)
}

func handle_list_parttypes(cmdline string) {
	lst := pnsql.GetPartTypes()
	tbl := util.NewTable("Digit", "Part Type")
	for _, c := range lst {
		tbl.AddRow(c.Digit, c.Description)
	}
	fmt.Printf("\n%s\n", tbl.Text())
}
