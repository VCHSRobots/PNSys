// --------------------------------------------------------------------
// cmd_delete.go -- Deletes anything.
//
// Created 2018-09-27 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/pnsql"
	"fmt"
)

var gTopic_delete_thing string = `
The delete command will delete either a part, a subsystem, a project, 
or a designer.  In some cases (such as designers and parts) no safety
checks are done, and the object is deleted with mercy.  In other cases,
such as for projects and subsystems, safety checks are performed.

The format of the command is:

    delete thing

where thing can be a designer's name (in quotes), a part number (in the
form "ppp-ss-000", a subsystem Id (in the form "ppp-ss") or a project id
(in the form "ppp"). 

`

func init() {
	RegistorCmd("delete", "thing", "Delete anything -- be careful! (see topic).", handle_delete_thing)
	RegistorTopic("delete", gTopic_delete_thing)
}

func handle_delete_thing(cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		fmt.Printf("Not enough args.\n")
		return
	}
	thing := args[1]
	prtresult := func(err error) {
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("Success.\n")
		}
	}

	err = pnsql.CheckDesignerNameText(thing)
	if err == nil {
		prtresult(pnsql.DeleteDesigner(thing))
		return
	}

	pt, _ := pnsql.ClassifyPN(thing)

	if pt == pnsql.PNType_Supplier {
		part, err := pnsql.GetSupplierPart(thing)
		if err != nil || part == nil {
			fmt.Printf("Part %s not found.\n", thing)
			return
		}
		prtresult(pnsql.DeleteSupplierPart(part))
		return
	}

	if pt == pnsql.PNType_Epic {
		pn, err := pnsql.StrToEpicPN(thing)
		if err == nil {
			part, err := pnsql.GetEpicPart(pn.PNString())
			if err != nil || part == nil {
				fmt.Printf("Part %s not found.\n", pn.PNString())
				return
			}
			prtresult(pnsql.DeleteEpicPart(part))
			return
		}
	}

	projectid, subsystemid, err := pnsql.SplitProjectId(thing)
	if err == nil {
		if subsystemid == "" {
			prtresult(pnsql.DeleteProject(projectid))
		} else {
			prtresult(pnsql.DeleteSubsystem(projectid, subsystemid))
		}
		return
	}
	fmt.Printf("%q cannot be indentified.\n", thing)
}
