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

func OpenDatabase(pw string) error {
	var err error
	connection := fmt.Sprintf("root:%s@/PnSysData", pw)
	m_db, err = sql.Open("mysql", connection)
	if err != nil {
		err := fmt.Errorf("Unable to open database. Err=%v\n", err)
		return err
	}
	return nil
}
