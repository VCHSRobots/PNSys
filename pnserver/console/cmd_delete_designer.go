// --------------------------------------------------------------------
// cmd_delete_designer.go -- Deletes a designer.
//
// Created 2018-09-22 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/pnsql"
	"strconv"
	"strings"
)

var gTopic_delete_designer string = `
The delete-designer will delete a designer either by a reference number or by a name.
The format of the command is:

    delete-designer name 

or 
 
    delete-designer ref=nnn

where name is in quote and must match the designer's name exactly, or the reference
number is the number next to the names after a list-designer command.  

Warning: nocheck is made on how many parts are assigned to the designer.  If there are
parts under a designer's name, and his/her name is deleted, you will not be able to 
filter part lists by that person's name.

`

func init() {
	RegistorCmd("delete-designer", "name [ref=#]", "Deletes a designer by name or ref #.", handle_delete_designer)
	RegistorTopic("delete-designer", gTopic_delete_designer)
}

func handle_delete_designer(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	alldesigners := pnsql.GetDesigners()
	var name string
	refs, ok := params["ref"]
	if !ok {
		if len(args) < 2 {
			c.Printf("Not enough args.\n")
			return
		}
		name = args[1]
	} else {
		indx, err := strconv.Atoi(strings.TrimSpace(refs))
		if err != nil {
			c.Printf("Bad ref # (%s). Err=%v\n", refs, err)
			return
		}

		if indx < 0 || indx >= len(alldesigners) {
			c.Printf("Ref number (%d) out of range. Total number of desginers is %d\n",
				indx, len(alldesigners))
			return
		}
		name = alldesigners[indx].Name
	}
	iCount := 0
	for _, d := range alldesigners {
		if d.Name == name {
			iCount += 1
		}
	}
	if iCount <= 0 {
		c.Printf("Designer %q not found.\n", name)
		return
	} else if iCount > 1 {
		c.Printf("%d duplicate records found for %q.\n", iCount, name)
	}
	c.Printf("You are about to delete %s from the list of designers.\n", name)
	if !AreYouSure() {
		return
	}
	err = pnsql.DeleteDesigner(name)
	if err != nil {
		c.Printf("Error %v\n", err)
		return
	}
	c.Printf("Success. Designer %q removed.\n", name)
}
