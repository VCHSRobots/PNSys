// --------------------------------------------------------------------
// cmd_universal_pw.go -- Set or unset the universal password mode.
//
// Created 2018-10-04 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/config"
	"epic/pnserver/sessions"
	"strings"
)

const gTopic_universal_pw = `
The universal-pw command enables or disables the universal password mode.  If this
mode is enabled, everybody can use the same set of passwords to log in.  If this 
mode is disabled, each designer must use an individually assigned password.

The format of this command is:

    universal-pw [true/false]

where either 'true' or 'false' must be provided.  

Note this mode will only stay active untill the next time the server restarts.
To override this, use the save-config command.

`

func init() {
	RegistorCmd("universal-pw", "t/f", "Enables or disables the universal password mode (see topic).", handle_universal_pw)
	RegistorTopic("universal-pw", gTopic_universal_pw)
}

func handle_universal_pw(c *Context, cmdline string) {
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
		sessions.SetAllowUniversalPasswords(true)
		config.SetBoolParam("allow_universal_pw", true)
		c.Printf("Success.\n")

	} else if mode == "false" {
		sessions.SetAllowUniversalPasswords(false)
		config.SetBoolParam("allow_universal_pw", false)
		c.Printf("Success.\n")

	} else {
		c.Printf("True or False not found.\n")
		return
	}
}
