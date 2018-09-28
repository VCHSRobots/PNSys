// --------------------------------------------------------------------
// cmd_delete_part.go -- Deletes a part.
//
// Created 2018-09-26 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/pnsql"
	"fmt"
)

var gTopic_delete_part string = `
The delete-part command will delete a part from the database without
mercy.  The format of the command is:

    delete-part ppp-ss-0000

where ppp-ss-0000 denotes the part number.  There is no recovery once
a part is deleted, other than manually re-entering all the data for the 
part.

`

func init() {
	RegistorCmd("delete-part", "pn", "Deletes a part number.", handle_delete_part)
	RegistorTopic("delete-part", gTopic_delete_part)
}

func handle_delete_part(cmdline string) {
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
	spn := args[1]
	part, err := pnsql.GetEpicPart(spn)
	if err != nil {
		fmt.Printf("Part %q does not exist.\n")
		return
	}

	err = pnsql.DeleteEpicPart(part)
	if err != nil {
		fmt.Printf("Error removing part %s. Err=%v", part.PNString(), err)
		return
	}
	fmt.Printf("Success. Part %s removed.\n", part.PNString())
	return
}
