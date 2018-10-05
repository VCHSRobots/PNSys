// --------------------------------------------------------------------
// cmd_delete_password.go -- Deletes passwords for a user.
//
// Created 2018-10-04 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"strings"
)

var gTopic_delete_passwords string = `
The delete-passwords command will delete all passwords for a user.  
The format of the command is:

    delete-passwords user=name

where name is in quotes and must match the designer's name exactly.  

All passwords associated with that name are deleted.  If you want
to delete the universal passwords, use "univeral" for the name.

`

func init() {
	RegistorCmd("delete-passwords", "", "Deletes all passwords for a user (see topic).", handle_delete_passwords)
	RegistorTopic("delete-passwords", gTopic_delete_passwords)
}

func handle_delete_passwords(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	name, ok := util.MapAlias(params, "user", "name", "Name", "User")
	if !ok {
		c.Printf("The name parameter not provided. See topic.\n")
		return
	}

	if util.Blank(name) {
		c.Printf("Name cannot be blank.\n")
		return
	}
	if strings.ToLower(name) == "universal" {
		name = ""
	}

	err = pnsql.DeletePasswords(name)
	if err != nil {
		c.Printf("Error deleting passwords. Err=%v\n", err)
		return
	}
	c.Printf("Success.\n")
}
