// --------------------------------------------------------------------
// parttypes_test.go -- Test parttypes.
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package pnsql

import (
	"fmt"
	"testing"
)

func Test_PartTypes(t *testing.T) {
	lst := GetPartTypes()
	for _, x := range lst {
		fmt.Printf("%s -- %s\n", x.Digit, x.Description)
	}
}
