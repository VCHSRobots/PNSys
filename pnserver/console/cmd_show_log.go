// --------------------------------------------------------------------
// cmd_show_log.go -- show or hide log output
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/log"
)

func init() {
	RegistorCmd("show-log", "", "Shows the log output as it occurs.", handle_showlog)
	RegistorCmd("hide-log", "", "Hides the log output from the console.", handle_hidelog)
	RegistorCmd("show-gin", "", "Shows the gin output on the console as it occurs.", handle_showgin)
	RegistorCmd("hide-gin", "", "Hides the gin output on the console.", handle_hidegin)
}

func handle_showlog(c *Context, cmdline string) {
	if c.IsInternal() {
		c.Printf("Cannot use this command from an external connection.\n")
		return
	}
	log.UseConsole(true)
}

func handle_hidelog(c *Context, cmdline string) {
	if c.IsInternal() {
		c.Printf("Cannot use this command from an external connection.\n")
		return
	}
	log.UseConsole(false)
}

func handle_showgin(c *Context, cmdline string) {
	if c.IsInternal() {
		c.Printf("Cannot use this command from an external connection.\n")
		return
	}
	log.AllowPassOnConsole(true)
}

func handle_hidegin(c *Context, cmdline string) {
	if c.IsInternal() {
		c.Printf("Cannot use this command from an external connection.\n")
		return
	}
	log.AllowPassOnConsole(false)
}
