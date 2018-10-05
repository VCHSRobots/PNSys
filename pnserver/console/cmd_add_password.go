// --------------------------------------------------------------------
// cmd_add_password.go -- Adds a new password.
//
// Created 2018-10-04 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/pwhash"
	"epic/lib/util"
	"epic/pnserver/pnsql"
	pv "epic/pnserver/privilege"
	"strings"
)

var gTopic_add_password string = `
The add-password command adds a new password to the database. The format of the
command is:

  add-password user=name privilege=pppp pw="something"

where all three parameters must be provided.  The name is a designer's name, 
privilege is one of "admin", "user", or "guest".  And "something" is the
clear text password. A hash will be generated for the password, and it will
never be seen again. 

For universal passwords, use "universal" for the name.

`

func init() {
	RegistorCmd("add-password", "", "Adds a new password (see topic).", handle_add_password)
	RegistorTopic("add-password", gTopic_add_password)
}

func handle_add_password(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	name, ok1 := util.MapAlias(params, "user", "name", "Name", "User")
	spriv, ok2 := util.MapAlias(params, "privilege", "Privilege", "priv")
	pw, ok3 := util.MapAlias(params, "pw", "password", "Password")
	if !ok1 || !ok2 || !ok3 {
		c.Printf("One of the required parameters not provided. See topic.\n")
		return
	}
	if util.Blank(name) {
		c.Printf("Name cannot be blank.\n")
		return
	}
	if strings.ToLower(name) == "universal" {
		name = ""
	}
	priv, err := pv.StrToPrivilege(spriv)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if util.Blank(pw) {
		c.Printf("The password cannot be blank.\n")
		return
	}

	hash, err := pwhash.HashPassword(pw)
	if err != nil {
		c.Printf("Unable to hash the password. Err=%v\n", err)
		return
	}
	err = pnsql.AddPassword(name, priv, hash)
	if err != nil {
		c.Printf("Unable to save password. Err=%v\n", err)
		return
	}
	c.Printf("Success.\n")
}
