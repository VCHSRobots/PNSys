// --------------------------------------------------------------------
// parttypes.go -- Manage the part-types table
//
// Created 2018-09-20 DLB
// --------------------------------------------------------------------

package pnsql

import (
	"fmt"
	"strings"
)

type PartType struct {
	Digit       string
	Description string
}

var gPartTypeCache []*PartType

func GetPartTypes() []*PartType {
	if gPartTypeCache != nil {
		return gPartTypeCache
	}
	lst := make([]*PartType, 0, 11)
	rows, err := m_db.Query("Select Digit, Description from PartTypes")
	if err != nil {
		fmt.Printf("Err getting PartTypes. Returning null slice. Err=%v\n", err)
		return lst
	}
	for rows.Next() {
		var digit string
		var desc string
		err = rows.Scan(&digit, &desc)
		if err != nil {
			fmt.Printf("Err during row scan in GetPartTypes. Err=%v\n", err)
			continue
		}
		lst = append(lst, &PartType{digit, desc})
	}
	gPartTypeCache = lst
	return gPartTypeCache
}

func GetPartTypeSelStrings() []string {
	plst := GetPartTypes()
	lst := make([]string, 0, len(plst))
	for _, p := range plst {
		s := p.Digit + " -- " + p.Description
		lst = append(lst, s)
	}
	return lst
}

func GetDigitFromSelString(sel string) (Digit string, err error) {
	plst := GetPartTypes()
	wrds := strings.Split(sel, " ")
	if len(wrds) <= 0 {
		return "?", fmt.Errorf("Unknown Part Type.")
	}
	pdigit := strings.TrimSpace(wrds[0])
	for _, p := range plst {
		if pdigit == p.Digit {
			return pdigit, nil
		}
	}
	return "?", fmt.Errorf("Unknown Part Type.")
}
