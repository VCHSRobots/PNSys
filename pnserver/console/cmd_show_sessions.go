// --------------------------------------------------------------------
// cmd_show_sessions.go -- Show all the sessions
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/sessions"
	"fmt"
)

func init() {
	RegistorCmd("show-sessions", "", "Shows current sessions.", handle_showsessions)
}

func handle_showsessions(cmdline string) {
	lst := sessions.GetAllSessions()
	tbl := util.NewTable("Ref#", "Name", "ClientIP", "Login Time", "Last Access", "Auth Cookie")
	for i, c := range lst {
		tbl.AddRow(fmt.Sprintf("%d", i), c.Name, c.ClientIP,
			c.LoginTime.Format("2006-01-02 15:04:05"),
			c.LastAccess.Format("2006-01-02 15:04:05"),
			c.AuthCookie.String())
	}
	fmt.Printf("\n%s\n", tbl.Text())
}
