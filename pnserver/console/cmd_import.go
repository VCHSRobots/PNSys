// --------------------------------------------------------------------
// cmd_import.go -- Imports old part numbers.
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

var gTopic_import string = `
The import command imports old part numbers into the database from a csv
file.  The format of the command is:

    import filename type=xxx

where filename is the path to an csv file, and xxx is either 'epic' or
'supplier' depending on the type of parts being imported.  The type 
parameter is required.

Epic csv files must have the following columns, in order, without any
header:  sequence number, part number string, description, designer, and date.

Supplier csv files must have the following columns, in order, without any
header: sequence number, part number string, vendor name, vendor part number,
description, designer, date, weblink.

In both cases the sequence number is simply a line number that is ignored. 
All the other named fields are used.  Columns after the named fields are ignored.

`

func init() {
	RegistorCmd("import", "file", "Imports old part numbers (see topic).", handle_import)
	RegistorTopic("import", gTopic_import)
}

func handle_import(cmdline string) {
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
		import_epic_csv(fn)
		return
	}
	if filetype == "supplier" {
		import_supplier_csv(fn)
		return
	}
	fmt.Printf("Unknown type (%q). Use either 'epic' or 'supplier'.\n", filetype)
}

func import_supplier_csv(fn string) {
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
	nadd := 0
	t0 := time.Now()
	for {
		elp := time.Now().Sub(t0)
		if elp > 20*time.Second {
			t0 = time.Now()
			fmt.Printf("Still working. Lines Processed=%d, Parts Added=%d\n", ilinenum, nadd)
		}
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

		_, err = strconv.Atoi(sseq)
		if err != nil {
			fmt.Printf("Bad sequence number on line %d\n", ilinenum)
			continue
		}
		date, err := util.ParseGenericTime(sdate)
		if err != nil {
			fmt.Printf("Bad date format on line %d. Err=%v\n", ilinenum, err)
			continue
		}
		pn, err := pnsql.StrToSupplierPartPN(spn)
		if err != nil {
			fmt.Printf("%v. On line %d.\n", err, ilinenum)
			continue
		}
		if util.Blank(sdesc) {
			fmt.Printf("Blank description on line %d\n", ilinenum)
			continue
		}

		p := &pnsql.SupplierPart{}
		p.SupplierPartPN = pn
		p.Description = util.CleanStr(sdesc, "|")
		p.Vendor = util.CleanStr(svendor, "|")
		p.VendorPN = util.CleanStr(svendorpn, "|")
		p.WebLink = util.CleanStr(slink, "|")
		p.Designer = sdesigner
		p.DateIssued = date

		err = pnsql.AddSupplierPart(p)
		if err != nil {
			fmt.Printf("Unable to add supplier part (line %d, pn %s). Err=%v\n", ilinenum,
				p.PNString(), err)
			continue
		} else {
			nadd++
		}
	}
	fmt.Printf("Finished!\n")
	fmt.Printf("%d lines read.  %d parts added to database.\n", ilinenum, nadd)
}

func import_epic_csv(fn string) {
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
	nadd := 0
	t0 := time.Now()
	for {
		elp := time.Now().Sub(t0)
		if elp > 20*time.Second {
			t0 = time.Now()
			fmt.Printf("Still working. Lines Processed=%d, Parts Added=%d\n", ilinenum, nadd)
		}
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
		date, err := time.Parse("01/02/06", sdate)
		if err != nil {
			fmt.Printf("Bad date format on line %d. Err=%v\n", ilinenum, err)
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
		p := &pnsql.EpicPart{}
		p.Description = sdesc
		p.Designer = sdesigner
		p.DateIssued = date
		p.EpicPN, err = pnsql.StrToEpicPN(spn)
		if err != nil {
			fmt.Printf("%v. On line %d\n", err, ilinenum)
			continue
		}
		err = pnsql.AddEpicPart(p)
		if err != nil {
			fmt.Printf("Unable to add epic part (line %d, pn %s). Err=%v\n", ilinenum,
				p.PNString(), err)
			continue
		} else {
			nadd++
		}
	}
	fmt.Printf("Finished!\n")
	fmt.Printf("%d lines read.  %d parts added to database.\n", ilinenum, nadd)
}
