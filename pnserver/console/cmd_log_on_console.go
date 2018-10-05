// --------------------------------------------------------------------
// cmd_log_on_console.go -- show or hide log output on the console
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/log"
	"epic/pnserver/config"
	"strings"
)

const gTopic_log_on_console = `
The log-on-console command enables or disables writing log messages to the console.
The format of the command is:

	log-on-console [true/false]

Note: this command can only be issued at a local connection.  It will not work
over http.  And, over a http connection, log messages are never written to the terminal.

`

func init() {
	RegistorCmd("log-on-console", "", "Shows or hides the log output as it occurs.", handle_log_on_console)
	RegistorTopic("log-on-console", gTopic_log_on_console)
}

func handle_log_on_console(c *Context, cmdline string) {
	if c.IsExternal() {
		c.Printf("Cannot use this command from an external connection.\n")
		return
	}
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		c.Printf("Not enough args.\n")
		return
	}
	mode := strings.ToLower(args[1])
	if mode == "true" {
		log.UseConsole(true)
		config.SetBoolParam("log_on_console", true)
		c.Printf("Success.\n")

	} else if mode == "false" {
		log.UseConsole(false)
		config.SetBoolParam("log_on_console", false)
		c.Printf("Success.\n")

	} else {
		c.Printf("True or False not found.\n")
		return
	}
}
