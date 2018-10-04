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
	"time"
)

func init() {
	RegistorCmd("show-sessions", "", "Shows current sessions.", handle_showsessions)
}

func handle_showsessions(c *Context, cmdline string) {
	lst := sessions.GetAllSessions()
	tbl := util.NewTable("Ref#", "Name", "ClientIP", "Login Time", "Last Access", "Elp (mins)", "Auth Cookie")
	for i, c := range lst {
		elp := time.Now().Sub(c.LastAccess)
		tbl.AddRow(fmt.Sprintf("%d", i), c.Name, c.ClientIP,
			c.LoginTime.Format("2006-01-02 15:04:05"),
			c.LastAccess.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%6.1f", elp.Minutes()),
			c.AuthCookie)
	}
	c.Printf("\n%s\n", tbl.Text())
}
