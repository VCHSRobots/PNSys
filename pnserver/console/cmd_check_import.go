// --------------------------------------------------------------------
// cmd_check_import.go -- Checks import file for errors.
//
// Created 2018-09-25 DLB
// --------------------------------------------------------------------

package console

import (
	"encoding/csv"
	"epic/lib/util"
	"epic/pnserver/pnsql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var gTopic_check_import string = `
The import command checks import file for errors. The format of the command is:

    check_import filename type=xxx

where filename is the path to an csv file, and xxx is either 'epic' or
'supplier' depending on the type of parts being imported.  The type 
parameter is required. 

The propoer contents of the csv files is descripbed under 'help import'.

`

func init() {
	RegistorCmd("check-import", "file", "Checks import file for errors.", handle_check_import)
	RegistorTopic("check-import", gTopic_check_import)
}

func handle_check_import(cmdline string) {
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
	fn := args[1]
	filetype, ok := util.MapAlias(params, "Type", "type")
	if !ok {
		fmt.Printf("The type parameters must be specified. See 'help import'.\n")
		return
	}
	filetype = strings.ToLower(filetype)
	if filetype == "epic" {
		check_epic_import(fn)
		return
	}
	if filetype == "supplier" {
		check_supplier_csv(fn)
		return
	}
	fmt.Printf("Unknown type (%q). Use either 'epic' or 'supplier'.\n", filetype)
}

func check_supplier_csv(fn string) {
	fi, err := os.Open(fn)
	if err != nil {
		fmt.Printf("Unable to open file %q. Err=%v\n", fn, err)
		return
	}
	defer fi.Close()

	rdr := csv.NewReader(fi)
	rdr.FieldsPerRecord = 0     // Must all be the same as first record.
	rdr.LazyQuotes = false      // No illegal quoting allowed
	rdr.TrimLeadingSpace = true // Ignore leading white space in a field.
	rdr.ReuseRecord = false     // Allocate mem for each new record.
	ilinenum := 0
	for {
		ilinenum++
		record, err := rdr.Read()
		if err == io.EOF {
			ilinenum--
			err = nil
			break
		}
		if err != nil {
			fmt.Printf("Read error on line %d. Aborting. Err=%v\n", ilinenum, err)
			return
		}
		if len(record) < 8 {
			fmt.Printf("Line %d has too few fields, skipping.\n", ilinenum)
			continue
		}

		sseq := record[0]      // Sequence number
		spn := record[1]       // Part number
		svendor := record[2]   // Vendor name
		svendorpn := record[3] // Vendor part number
		sdesc := record[4]     // Descriptio
		sdesigner := record[5] // Designer
		sdate := record[6]     // Date
		slink := record[7]     // Web link

		pn, _ := pnsql.GetSupplierPart(spn)
		if pn != nil {
			fmt.Printf("Part %s already exsits, on line %d\n", pn.PNString(), ilinenum)
			continue
		}

		if !pnsql.IsDesigner(sdesigner) {
			fmt.Printf("Designer %q does not exist, on line %d\n", sdesigner, ilinenum)
		}

		_, err = strconv.Atoi(sseq)
		if err != nil {
			fmt.Printf("Bad sequence number on line %d\n", ilinenum)
			continue
		}
		_, err = util.ParseGenericTime(sdate)
		if err != nil {
			fmt.Printf("Bad date format on line %d. Err=%v\n", ilinenum, err)
			continue
		}
		_, err = pnsql.StrToSupplierPartPN(spn)
		if err != nil {
			fmt.Printf("%v. On line %d.\n", err, ilinenum)
			continue
		}
		if util.Blank(sdesc) {
			fmt.Printf("Blank description on line %d\n", ilinenum)
			continue
		}
		if sdesc != util.CleanStr(sdesc, "|") {
			fmt.Printf("Illegal chars in Description on line %d\n", ilinenum)
		}
		if svendor != util.CleanStr(svendor, "|") {
			fmt.Printf("Illegal chars in Vendor on line %d\n", ilinenum)
		}
		if svendorpn != util.CleanStr(svendorpn, "|") {
			fmt.Printf("Illegal chars in VenderPN on line %d\n", ilinenum)
		}
		if slink != util.CleanStr(slink, "|") {
			fmt.Printf("Illegal chars in WebLink on line %d\n", ilinenum)
		}
	}
	fmt.Printf("Finished. Lines Processed = %d.\n", ilinenum)
}

func check_epic_import(fn string) {
	fi, err := os.Open(fn)
	if err != nil {
		fmt.Printf("Unable to open file %q. Err=%v\n", fn, err)
		return
	}
	defer fi.Close()

	rdr := csv.NewReader(fi)
	rdr.FieldsPerRecord = 0     // Must all be the same as first record.
	rdr.LazyQuotes = false      // No illegal quoting allowed
	rdr.TrimLeadingSpace = true // Ignore leading white space in a field.
	rdr.ReuseRecord = false     // Allocate mem for each new record.
	ilinenum := 0
	for {
		ilinenum++
		record, err := rdr.Read()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			fmt.Printf("Read error on line %d. Aborting. Err=%v\n", ilinenum, err)
			return
		}
		if len(record) < 5 {
			fmt.Printf("Line %d has too few fields, skipping.\n", ilinenum)
			continue
		}
		sseq := record[0]      // Sequence number
		spn := record[1]       // Part number
		sdesc := record[2]     // Description
		sdesigner := record[3] // Designer
		sdate := record[4]     // date

		_, err = strconv.Atoi(sseq)
		if err != nil {
			fmt.Printf("Bad sequence number on line %d\n", ilinenum)
			continue
		}

		pn, _ := pnsql.GetEpicPart(spn)
		if pn != nil {
			fmt.Printf("Part %s already exsits, on line %d\n", pn.PNString(), ilinenum)
			continue
		}

		err = checkepicpn(spn)
		if err != nil {
			fmt.Printf("%v. On line %d.\n", err, ilinenum)
			continue
		}
		if util.Blank(sdesc) {
			fmt.Printf("Blank description on line %d\n", ilinenum)
			continue
		}
		if !pnsql.IsDesigner(sdesigner) {
			fmt.Printf("Designer %q is unknown. On line %d\n", sdesigner, ilinenum)
		}
		// Check valid date
		date, err := time.Parse("01/02/06", sdate)
		if err != nil {
			fmt.Printf("Bad date format on line %d. Err=%v\n", ilinenum, err)
			continue
		}
		t0 := time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)
		t1 := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
		if date.Before(t0) || date.After(t1) {
			fmt.Printf("Date is out of range on line %d.", ilinenum)
			continue
		}

	}
	fmt.Printf("%d lines read and examined.\n", ilinenum)
}

func checkepicpn(s string) error {
	pn, err := pnsql.StrToEpicPN(s)
	if err != nil {
		return err
	}
	haveit := false
	for _, t := range pnsql.GetPartTypes() {
		if t.Digit == pn.PartType {
			haveit = true
			break
		}
	}
	if !haveit {
		return fmt.Errorf("Part type (%s) is not known.", pn.PartType)
	}
	if !pnsql.IsProject(pn.ProjectId) {
		return fmt.Errorf("Project Id (%s) is not known.", pn.ProjectId)
	}
	if !pnsql.IsSubsystem(pn.ProjectId, pn.SubsystemId) {
		return fmt.Errorf("Subsystem Id (%s-%s) is not known.", pn.ProjectId, pn.SubsystemId)
	}
	return nil
}
