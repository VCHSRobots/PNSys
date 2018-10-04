// --------------------------------------------------------------------
// cmd_update_part.go -- Updates a part.
//
// Created 2018-09-26 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"epic/pnserver/pnsql"
)

var gTopic_update_part string = `
The update-part command allows you to change the date issued, the designer or 
the description of a part.  It works as follows:

  update-part ppp-ss-0000 designer="X. Name" desc="comments about the part" date="yyyy-mm-dd"

One, two or all three parameters can be specified.

`

func init() {
	RegistorCmd("update-part", "pn", "Updates a part number (see topic).", handle_update_part)
	RegistorTopic("update-part", gTopic_update_part)
	RegistorArg("pn", "A partnumber in the ppp-ss-0000 or SP-cc-000 format.")
}

func handle_update_part(c *Context, cmdline string) {
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
	spn := args[1]
	part, err := pnsql.GetEpicPart(spn)
	if err != nil || part == nil {
		c.Printf("Part %q not found.\n", spn)
		return
	}

	nupdate := 0

	designer, ok := params["Designer"]
	if !ok {
		designer, ok = params["designer"]
	}
	if ok {
		if !pnsql.IsDesigner(designer) {
			c.Printf("%s is not a current designer.\n", designer)
			return
		}
		err := pnsql.SetEpicPartDesigner(part, designer)
		pnsql.InvalidateEpicPartsCache()
		if err != nil {
			c.Printf("Error on update of designer. Err=%v\n", err)
			return
		}
		nupdate++
	}

	desc, ok := params["Description"]
	if !ok {
		desc, ok = params["description"]
		if !ok {
			desc, ok = params["desc"]
			if !ok {
				desc, ok = params["Desc"]
			}
		}
	}
	if desc != "" {
		err := pnsql.SetEpicPartDescription(part, desc)
		pnsql.InvalidateEpicPartsCache()
		if err != nil {
			c.Printf("Error on update of description. Err=%v\n", err)
			return
		}
		nupdate++
	}

	sdate, ok := params["Date"]
	if !ok {
		sdate, ok = params["date"]
		if !ok {
			sdate, ok = params["DateIssued"]
			if !ok {
				sdate, ok = params["dateissued"]
			}
		}
	}
	if sdate != "" {
		date, err := util.ParseGenericTime(sdate)
		if err != nil {
			c.Printf("Syntax error in date (%q).\n", sdate)
			return
		}
		err = pnsql.SetEpicPartDateIssued(part, date)
		pnsql.InvalidateEpicPartsCache()
		if err != nil {
			c.Printf("Error on update of date issued. Err=%v\n", err)
			return
		}
		nupdate++
	}

	if nupdate <= 0 {
		c.Printf("Nothing updated (see topic). \n")
		return
	}
	c.Printf("Success.\n")
}
