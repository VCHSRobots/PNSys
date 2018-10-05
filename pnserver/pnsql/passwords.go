// --------------------------------------------------------------------
// passwords.go -- Manage the passwords table
//
// Created 2018-10-04 DLB
// --------------------------------------------------------------------

package pnsql

import (
	"epic/lib/log"
	pv "epic/pnserver/privilege"
	"fmt"
	"sync"
)

type Password struct {
	Name      string       // Empty names means it is a univeral password
	Privilege pv.Privilege // To grant admin privilege with this password
	Hash      string       // The password hash
}

var gPasswords []*Password   // The password cache
var gPasswordLock sync.Mutex // For General access to the cache

func InvalidatePasswordCache() {
	gPasswordLock.Lock()
	defer gPasswordLock.Unlock()
	gPasswords = nil
}

// GetPasswords returns all the passwords in the system.
func GetPasswords() []*Password {
	gPasswordLock.Lock()
	defer gPasswordLock.Unlock()
	if gPasswords != nil {
		return gPasswords
	}
	lst := make([]*Password, 0, 100)
	rows, err := m_db.Query("Select Name, Privilege, Hash from Passwords")
	if err != nil {
		fmt.Printf("Err getting Passwords. Returning null slice. Err=%v\n", err)
		return lst
	}
	for rows.Next() {
		var name string
		var priv int
		var hash string
		err = rows.Scan(&name, &priv, &hash)
		if err != nil {
			fmt.Printf("Err during row scan in GetPasswords. Err=%v\n", err)
			continue
		}
		lst = append(lst, &Password{name, pv.Privilege(priv), hash})
	}
	gPasswords = lst
	return gPasswords
}

// Returns a list of passwords that are valid for a given name.  Universal
// passwords are not included.  To get the universal passwords, use "" for
// name.
func GetPasswordsForName(name string) []*Password {
	lst := make([]*Password, 0, 10)
	for _, pw := range GetPasswords() {
		if pw.Name == name {
			lst = append(lst, pw)
		}
	}
	return lst
}

// AddPassword will add a password hash to the table under the given name.  If the
// name is blank, it is considered a univeral password and will work for everybody.
func AddPassword(name string, priv pv.Privilege, hash string) error {
	stmt, err := m_db.Prepare("insert Passwords set Name=?, Privilege=?, Hash=?")
	if err != nil {
		err := fmt.Errorf("Err inserting into Passwords. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	_, err = stmt.Exec(name, int(priv), hash)
	InvalidatePasswordCache()
	if err != nil {
		err := fmt.Errorf("Err inserting into Passwords. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	if name == "" {
		log.Infof("Added a new universal password for %s privilege.", priv.String())
	} else {
		log.Infof("Added a new password for user %s, for %s privilege.", name, priv.String())
	}
	return nil
}

// DeletePasswords will remove ALL passwords associated with a name.  If the name is blank
// then the univeral passwords will be removed.
func DeletePasswords(name string) error {
	stmt, err := m_db.Prepare("delete from Passwords where Name=?")
	if err != nil {
		err := fmt.Errorf("Err deleting from Passwords. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	_, err = stmt.Exec(name)
	InvalidatePasswordCache()
	if err != nil {
		err := fmt.Errorf("Err deleting from Passwords. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	if name == "" {
		log.Infof("Deleted all universal passwords.")
	} else {
		log.Infof("Deleted all passwords for user %s.", name)
	}
	return nil
}

// RestoreOriginalPasswords will restore the developer level passwords into the system.
// They "should" be changed once the system is up and running.
func RestoreOriginalPasswords() error {
	const (
		gPwAdmin = "JDJhJDA2JHNHWGt6SkM0RnVJL0EvZ25aNUlZai5wcC4uQWlxL1dNdlJDN2w5dkJJZHluR0xrZEd1djFT"
		gPwUser1 = "JDJhJDA2JHB0S2xmamVtL3k4OHpwZ1JEQzRYLy5qUWc1NjNqQWZzTUxOd3REUVFUM082UDVpTlZiWFpt"
		gPwUser2 = "JDJhJDA2JGJ3Z1hUZW1GaHVxTFhWYnBQb2FyUnVHeEdGdVZhLk9CZTVlWFZuWnpuNE9pSE5BS3RMYXk2"
		gPwUser3 = "JDJhJDA2JEsxNXovdzBaRHhoWVF0N3U1OXlmZU9rU2hQT3VhM3lJa0hBeFNmY3FDMTZhdDliOE90b0Yu"
	)

	type PwInfo struct {
		IsAdmin bool
		Hash    string
	}

	KnownPws := []PwInfo{{true, gPwAdmin}, {false, gPwUser1}, {false, gPwUser2}, {false, gPwUser3}}

	log.Infof("Attempting to restore original passwords.")
	err := DeletePasswords("")
	if err != nil {
		err = fmt.Errorf("Unable to restore original passwords. Unable to remove. Err=%v", err)
		log.Warnf("%v", err)
		return err
	}
	for _, pw := range KnownPws {
		lvl := pv.User
		if pw.IsAdmin {
			lvl = pv.Admin
		}
		err := AddPassword("", lvl, pw.Hash)
		if err != nil {
			err = fmt.Errorf("Unable to restore all original passwords. Err=%v", err)
			log.Warnf("%v", err)
		}
	}
	return nil
}
