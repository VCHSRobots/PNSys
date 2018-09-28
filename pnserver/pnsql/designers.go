// --------------------------------------------------------------------
// designers.go -- Manage the Designers table
//
// Created 2018-09-20 DLB
// --------------------------------------------------------------------

package pnsql

import (
	"epic/lib/log"
	"epic/lib/util"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Designer struct {
	Name   string
	Year0  string
	Active bool
}

var gDesignerCache []*Designer
var gDesignerLock sync.Mutex

func InvalidateDesignersCache() {
	gDesignerLock.Lock()
	defer gDesignerLock.Unlock()
	gDesignerCache = nil
}

func GetDesignerNames(filter bool) []string {
	dlst := GetDesigners()
	lst := make([]string, 0, len(dlst))
	for _, d := range dlst {
		if !filter || d.Active {
			lst = append(lst, d.Name)
		}
	}
	return lst
}

func GetDesigners() []*Designer {
	gDesignerLock.Lock()
	defer gDesignerLock.Unlock()
	if gDesignerCache != nil {
		return gDesignerCache
	}
	lst := make([]*Designer, 0, 11)
	rows, err := m_db.Query("Select Name, Year0, Active from Designers")
	if err != nil {
		log.Errorf("Err getting Designers. Returning null slice. Err=%v\n", err)
		return lst
	}
	for rows.Next() {
		var name string
		var year0 string
		var active int
		err = rows.Scan(&name, &year0, &active)
		if err != nil {
			log.Errorf("Err during row scan in GetDesigners. Err=%v\n", err)
			continue
		}
		bactive := active != 0
		lst = append(lst, &Designer{name, year0, bactive})
	}

	parselastname := func(name string) string {
		if len(name) < 4 {
			return name
		}
		if name[1:3] == ". " {
			return strings.TrimSpace(name[3:])
		}
		return name
	}
	sorter := func(i, j int) bool {
		if lst[i].Active == lst[j].Active {
			if lst[i].Year0 == lst[j].Year0 {
				return parselastname(lst[i].Name) < parselastname(lst[j].Name)
			}
			return lst[i].Year0 < lst[j].Year0
		}
		return lst[i].Active
	}
	sort.Slice(lst, sorter)
	gDesignerCache = lst
	return gDesignerCache
}

func IsDesigner(Name string) bool {
	for _, d := range GetDesigners() {
		if d.Name == Name {
			return true
		}
	}
	return false
}

func GetDesigner(Name string) (*Designer, error) {
	for _, d := range GetDesigners() {
		if d.Name == Name {
			return d, nil
		}
	}
	return nil, fmt.Errorf("Designer %q does not exist.", Name)
}

func CheckDesignerNameText(Name string) error {
	if util.Blank(Name) {
		return fmt.Errorf("Designer name cannot be blank.")
	}
	if len(Name) <= 4 {
		return fmt.Errorf("Designer name (%q) must be at least 5 chars long.", Name)
	}
	if Name[1:3] != ". " {
		return fmt.Errorf("Designer name (%q) does not follow prescribed format.", Name)
	}
	if !util.ContainsOnly(Name[0:1], "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return fmt.Errorf("Designer name (%q) does not have capitalized first initial.", Name)
	}
	if !util.ContainsOnly(Name[3:4], "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return fmt.Errorf("Designer name (%q) does not have capitalized last name.", Name)
	}
	lastpart := strings.ToLower(Name[3:])
	if !util.ContainsOnly(lastpart, "abcdefghijklmnopqrstuvwxyz0123456789_") {
		return fmt.Errorf("Designer name (%q) conttains illegal characters.", Name)
	}
	return nil
}

func AddDesigner(Name, Year0 string) error {
	err := CheckDesignerNameText(Name)
	if err != nil {
		return err
	}
	err = CheckYear0Text(Year0)
	if err != nil {
		return err
	}
	stmt, err := m_db.Prepare("insert Designers set Name=?, Year0=?, Active=1")
	if err != nil {
		log.Errorf("Err inserting into Designers. Err=%v", err)
		return err
	}
	r, err := stmt.Exec(Name, Year0)
	InvalidateDesignersCache()
	if err != nil {
		log.Errorf("Err inserting into Designers. Err=%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Wrong number of rows affected (%d) when adding designer %s.",
			n, Name)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func DeleteDesigner(Name string) error {
	InvalidateDesignersCache()
	if !IsDesigner(Name) {
		return fmt.Errorf("The designer %q does not exist.", Name)
	}

	stmt, err := m_db.Prepare("delete from Designers where Name=?")
	if err != nil {
		log.Errorf("Err deleting from Designers. Err=%v", err)
		return err
	}
	r, err := stmt.Exec(Name)
	InvalidateDesignersCache()
	if err != nil {
		log.Errorf("Err deleting from Designers. Err= %v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Wrong number of rows affected (%d) when deleting desinger %s.",
			n, Name)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func SetDesignerActive(Name string, Active bool) error {
	InvalidateDesignersCache()
	designer, err := GetDesigner(Name)
	if err != nil || designer == nil {
		return fmt.Errorf("Designer %q does not exist.", Name)
	}
	if designer.Active == Active {
		return fmt.Errorf("No change to database requested. Active equal.")
	}
	iactive := 0
	if Active {
		iactive = 1
	}
	stmt, err := m_db.Prepare("update Designers set Active=? where Name=?")
	if err != nil {
		err := fmt.Errorf("Err updating Designers. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(iactive, designer.Name)
	InvalidateDesignersCache()
	if err != nil {
		err := fmt.Errorf("Err updating Designers. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err := fmt.Errorf("Err updating Designers. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func SetDesignerYear0(Name string, Year0 string) error {
	InvalidateDesignersCache()
	designer, err := GetDesigner(Name)
	if err != nil || designer == nil {
		return fmt.Errorf("Designer %q does not exist.", Name)
	}
	err = CheckYear0Text(Year0)
	if err != nil {
		return err
	}
	if designer.Year0 == Year0 {
		return fmt.Errorf("No change to database requested. Year0 equal.")
	}
	stmt, err := m_db.Prepare("update Designers set Year0=? where Name=?")
	if err != nil {
		err := fmt.Errorf("Err updating Designers. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(Year0, designer.Name)
	InvalidateDesignersCache()
	if err != nil {
		err := fmt.Errorf("Err updating Designers. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err := fmt.Errorf("Err updating Designers. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	return nil
}
