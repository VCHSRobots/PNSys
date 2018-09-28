// --------------------------------------------------------------------
// pnsys.go -- Access to PnSysData database.
//
// Created 2018-09-20 DLB
// --------------------------------------------------------------------

package pnsql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var m_db *sql.DB

func init() {
	var err error
	m_db, err = sql.Open("mysql", "root:loveepic@/PnSysData")
	if err != nil {
		fmt.Printf("Unable to open database. Err=%v\n", err)
		panic(err)
	}
}
