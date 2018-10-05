// --------------------------------------------------------------------
// cmd_dev_bypass.go -- Set or unset developer bypass
//
// Created 2018-10-04 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/config"
	"epic/pnserver/pnsql"
	"epic/pnserver/sessions"
)

const gTopic_devbypass = `
The set-dev-bypass command will enable or disable the developer bypass mode.  When
this mode is enabled the login page does NOT check passwords, and therfore allows
anybody easy access.  In fact, the login page is completely bypassed for anybody
who has a connection to the web site -- so be careful.

The format of this command is:

    set-dev-bypass designer

The designer must be in quotes, and be in the database.  The designer name 
given will be the person that is automatically logged in, and he/she will have
Admin privilege.  

To disable this mode, issue the command without a designer name.

Note this mode will only stay active untill the next time the server restarts.
To override this, use the save-config command.

`

func init() {
	RegistorCmd("dev-bypass", "", "Enables or disables developer bypass mode (see topic).", handle_devbypass)
	RegistorTopic("dev-bypass", gTopic_devbypass)
}

func handle_devbypass(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		sessions.SetDeveloperBypass("")
		config.RemoveParam("dev_bypass")
		c.Printf("Success. (Dev bypass disabled, no designer name given.)\n")
		return
	}
	designer := args[1]
	if !pnsql.IsDesigner(designer) {
		c.Printf("%q is not a known desiger.\n", designer)
		return
	}
	sessions.SetDeveloperBypass(designer)
	config.SetParam("dev_bypass", designer)
	c.Printf("Success. (Dev bypass enabled for %s.)\n", designer)
}
