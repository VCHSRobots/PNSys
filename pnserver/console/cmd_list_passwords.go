// --------------------------------------------------------------------
// cmd_list_passwords.go -- Lists all the passwords
//
// Created 2018-10-04 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
)

var gTopic_list_passords string = `
The list-passwords command will list all the passwords. The
format of the command is:

    list-passwords user

where user is optional.  If provided, only passwords for that
user will be listed.  To list just the universal passwords, use
'universal' as the name of the user.

 `

func init() {
	RegistorCmd("list-passwords", "", "Lists all the passwords (see topic).", handle_list_passwords)
	RegistorTopic("list-passwords", gTopic_list_passords)
}

func handle_list_passwords(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	all := true
	name := ""
	if len(args) > 1 {
		name = args[1]
		all = false
	}

	lst := pnsql.GetPasswords()
	tbl := util.NewTable("Name", "Type", "pw.Hash")
	nc := 0
	for _, pw := range lst {
		if all || name == pw.Name || (name == "universal" && pw.Name == "") {
			name = pw.Name
			if name == "" {
				name = "<universal>"
			}
			hash := util.FixStrLen(pw.Hash, 30, "...")
			tbl.AddRow(name, pw.Privilege.String(), hash)
			nc++
		}
	}
	if nc <= 0 {
		c.Printf("Nothing found.\n")
	}
	c.Printf("\n%s%d Passwords found.\n", tbl.Text(), nc)
}
