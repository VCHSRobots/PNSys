// --------------------------------------------------------------------
// cmd_list_designers.go -- Lists all the designers
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
	"strings"
)

var gTopic_list_designers string = `
The list-designers command will list all the designers or a subset. The
format of the command is:

    list-designers active=[true/false] year0=yyyy

 where the parameters active and year0 are optional.  If provided the
 list will be filtered accordingly.

 `

func init() {
	RegistorCmd("list-designers", "", "Lists all the designers (see topic).", handle_list_designers)
	RegistorTopic("list-designers", gTopic_list_designers)
}

func handle_list_designers(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	var act bool
	sactive, doactive := util.MapAlias(params, "active")
	year0, doyear := util.MapAlias(params, "Year0", "year0", "Year", "year")
	if doactive {
		sact := strings.ToLower(sactive)
		if sact == "true" || sact == "yes" || sact == "t" || sact == "y" {
			act = true
		} else if sact == "false" || sact == "no" || sact == "f" || sact == "n" {
			act = false
		} else {
			c.Printf("Invalid value for active (%q)\n", sactive)
			return
		}
	}

	lst := pnsql.GetDesigners()
	icount := 0
	tbl := util.NewTable("Ref#", "Name", "Year0", "Active")
	for i, c := range lst {
		if doactive && act != c.Active {
			continue
		}
		if doyear && year0 != c.Year0 {
			continue
		}
		sactive := ""
		if c.Active {
			sactive = "Yes"
		}
		tbl.AddRow(fmt.Sprintf("%d", i), c.Name, c.Year0, sactive)
		icount++
	}
	c.Printf("\n%s%d desingers found.\n", tbl.Text(), icount)
}
