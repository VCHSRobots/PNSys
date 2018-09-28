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
}

func handle_showlog(cmdline string) {
	log.UseConsole(true)
}

func handle_hidelog(cmdline string) {
	log.UseConsole(false)
}
