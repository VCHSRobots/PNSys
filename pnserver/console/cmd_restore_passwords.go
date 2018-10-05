// --------------------------------------------------------------------
// cmd_restore_passwords.go -- Restores original passwords.
//
// Created 2018-10-04 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/pnsql"
)

var gTopic_restore_passwords string = `
The restore-passwords command will restore the original built-in,
universal, passwords into the system.  It will remove the existing
universal passwords before doing so, but it will not harm other
passwords stored for individual users.

The format of the command is:

    restore-passwords 

Note: these passwords are known to Mr. Brandon and are written
down in his password book.

`

func init() {
	RegistorCmd("restore-passwords", "", "Restores original passwords (for sys startup).", handle_restore_passwords)
	RegistorTopic("restore-passwords", gTopic_restore_passwords)
}

func handle_restore_passwords(c *Context, cmdline string) {

	err := pnsql.RestoreOriginalPasswords()
	if err != nil {
		c.Printf("%v\n", err)
	}
	c.Printf("Success.\n")
}
