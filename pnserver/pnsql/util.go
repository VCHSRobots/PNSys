// --------------------------------------------------------------------
// util.go -- Utilities for PartNumber sql package.
//
// Created 2018-09-27 DLB
// --------------------------------------------------------------------

package pnsql

import (
	"fmt"
	"regexp"
	"strconv"
)

type PNType string

const (
	PNType_Epic     = "Epic"
	PNType_Supplier = "Supplier"
	PNType_Unknown  = "Unknown"
)

func CheckYear0Text(s string) error {
	iyear, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("Illegal value for Year0.")
	}
	if iyear < 2000 || iyear > 2040 {
		return fmt.Errorf("Year0 is not an appropriate value (%d).", iyear)
	}
	return nil
}

func ClassifyPN(s string) (PNType, error) {
	// _, err1 := StrToSupplierPartPN(s)
	// _, err2 := StrToEpicPN(s)
	// if err1 == nil && err2 == nil {
	// 	return PNType_Unknown, fmt.Errorf("Appears to be both an epic part and suppiler part!")
	// }
	// if err1 != nil && err2 == nil {
	// 	return PNType_Epic, nil
	// }
	// if err1 == nil && err2 != nil {
	// 	return PNType_Supplier, nil
	// }
	// It looks bad both ways.  Do the best we can.

	// Just do pattern analysis for now.
	rx_supplier := regexp.MustCompile(`\ASP\-\d\d\-\d\d\d$`)   // Matches SP-dd-dddd
	rx_epic := regexp.MustCompile(`\A\S\S\S\-\S\S\-\d\d\d\d$`) // Matches aaa-aa-dddd
	if rx_supplier.MatchString(s) {
		return PNType_Supplier, nil
	}
	if rx_epic.MatchString(s) {
		return PNType_Epic, nil
	}
	return PNType_Unknown, fmt.Errorf("Not a known part number format.")
}
