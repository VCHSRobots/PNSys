// --------------------------------------------------------------------
// cmd_pwhash.go -- Prints a password hash
//
// Created 2018-09-29 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/pwhash"
	"fmt"
)

func init() {
	RegistorCmd("pwhash", "pw", "Computes password hash (for development).", handle_pwhash)
}

func handle_pwhash(cmdline string) {
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
	pw := args[1]
	hash, err := pwhash.HashPassword(pw)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Hash=%q\n", hash)
}
