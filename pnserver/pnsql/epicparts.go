// --------------------------------------------------------------------
// epicparts.go -- Manage epic parts
//
// Created 2018-09-25 DLB
// --------------------------------------------------------------------

package pnsql

import (
	"epic/lib/log"
	"epic/lib/util"
	"epic/lib/uuid"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// EpicPN struct contains a split version of a part number.
type EpicPN struct {
	ProjectId   string // 3 Chars long
	SubsystemId string // 2 Chars long
	PartType    string // 1 Char long
	SequenceNum int    // 0-999
}

type EpicPart struct {
	PID uuid.UUID
	*EpicPN
	Designer    string
	Description string
	DateIssued  time.Time
}

var gEpicParts []*EpicPart
var gEpicPartsLock sync.Mutex

func InvalidateEpicPartsCache() {
	gEpicPartsLock.Lock()
	defer gEpicPartsLock.Unlock()
	gEpicParts = nil
}

// StrToEpicPN converts a string to an Epic Part Number, if possible.
func StrToEpicPN(s string) (*EpicPN, error) {
	var err error
	parts := strings.Split(s, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Baddly formatted part number (%s) -- no parts.", s)
	}
	pn := &EpicPN{ProjectId: parts[0], SubsystemId: parts[1]}
	if len(pn.ProjectId) != 3 || len(pn.SubsystemId) != 2 {
		return pn, fmt.Errorf("Baddly formatted part number (%s).", s)
	}
	digits := parts[2]
	if len(digits) == 3 && pn.ProjectId == "Y14" {
		// Fix this special case
		digits = "9" + digits
	}
	if len(digits) != 4 {
		return pn, fmt.Errorf("Baddly formatted part number (%s) -- bad ending digits.", s)
	}
	pn.PartType = digits[0:1]
	pn.SequenceNum, err = strconv.Atoi(digits[1:])
	if err != nil {
		return pn, fmt.Errorf("Baddly formatted part number (%q)-- sequence not integer.", s)
	}
	return pn, nil
}

func (p *EpicPart) ProjectDesc() string {
	for _, prj := range GetProjects() {
		if prj.ProjectId == p.ProjectId {
			return prj.Description
		}
	}
	return ""
}

func (p *EpicPart) SubsystemDesc() string {
	for _, prj := range GetProjects() {
		if prj.ProjectId == p.ProjectId {
			for _, sub := range prj.Subsystems {
				if p.SubsystemId == sub.SubsystemId {
					return sub.Description
				}
			}
			return ""
		}
	}
	return ""
}

func (p *EpicPart) PartTypeDesc() string {
	for _, pt := range GetPartTypes() {
		if pt.Digit == p.PartType {
			return pt.Description
		}
	}
	return ""
}

func (p *EpicPart) PNString() string {
	return (p.EpicPN.PNString())
}

func (t *EpicPN) PNString() string {
	s := t.ProjectId + "-" + t.SubsystemId + "-" + t.PartType + fmt.Sprintf("%03d", t.SequenceNum)
	return s
}

func (p *EpicPart) CheckPNFormat() error {
	return p.EpicPN.CheckPNFormat()
}

func (t *EpicPN) CheckPNFormat() error {
	if t.SequenceNum < 0 || t.SequenceNum > 999 {
		return fmt.Errorf("Sequence number (%d) out of range.", t.SequenceNum)
	}
	pns := t.PNString()
	_, err := StrToEpicPN(pns)
	return err
}

func (p *EpicPart) SamePN(p2 *EpicPart) bool {
	return p.EpicPN.SamePN(p2.EpicPN)
}

func (t *EpicPN) SamePN(t2 *EpicPN) bool {
	if t.ProjectId != t2.ProjectId {
		return false
	}
	if t.SubsystemId != t2.SubsystemId {
		return false
	}
	if t.PartType != t2.PartType {
		return false
	}
	if t.SequenceNum != t2.SequenceNum {
		return false
	}
	return true
}

func GetEpicParts() []*EpicPart {
	gEpicPartsLock.Lock()
	defer gEpicPartsLock.Unlock()
	if gEpicParts != nil {
		return gEpicParts
	}
	rows, err := m_db.Query("Select PID, ProjectId, SubsystemId, PartType, SequenceNum, Designer, DateIssued, Description from EpicParts")
	if err != nil {
		log.Errorf("Err getting epic parts. Err=%v\n", err)
		return make([]*EpicPart, 0, 0)
	}
	parts := make([]*EpicPart, 0, 10000)
	for rows.Next() {
		var spid, projectid, subsystemid, parttype, designer, sdate, desc string
		var seqnum int
		err = rows.Scan(&spid, &projectid, &subsystemid, &parttype, &seqnum, &designer, &sdate, &desc)
		if err != nil {
			log.Errorf("Err during scan of EpicParts. Err=%v", err)
			continue
		}
		id, err := uuid.FromString(spid)
		if err != nil {
			log.Errorf("Bad uuid for EpidPart in database. (%s) Err=%v", spid, err)
			continue
		}
		dateissued, err := time.Parse("2006-01-02", sdate)
		if err != nil {
			log.Errorf("Bad date for EpicPart in database. (%s) Err=%v", spid, err)
			continue
		}
		p := &EpicPart{PID: id, Designer: designer, DateIssued: dateissued, Description: desc}
		p.EpicPN = &EpicPN{projectid, subsystemid, parttype, seqnum}
		parts = append(parts, p)
	}
	gEpicParts = parts
	return gEpicParts
}

// Retrieve an epic part given its part number string (XX-YYY-ZZZZ)
func GetEpicPart(pns string) (*EpicPart, error) {
	pn, err := StrToEpicPN(pns)
	if err != nil {
		return nil, err
	}
	for _, p := range GetEpicParts() {
		if p.EpicPN.SamePN(pn) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("Part %s not found.", pn.PNString())
}

func DeleteEpicPart(p *EpicPart) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Epic part not identified correctly. PID=0.")
	}
	stmt, err := m_db.Prepare("delete from EpicParts where PID=?")
	if err != nil {
		log.Errorf("Err deleting from EpicParts. Err=%v", err)
		return err
	}
	r, err := stmt.Exec(p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err deleting from EpicParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateEpicPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Err deleting from EpicParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func SetEpicPartDesigner(p *EpicPart, designer string) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Epic part not identified correctly. PID=0.")
	}
	if !IsDesigner(designer) {
		return fmt.Errorf("%q is not a known designer.")
	}
	stmt, err := m_db.Prepare("update EpicParts set Designer=? where PID=?")
	if err != nil {
		err := fmt.Errorf("Err updating EpicParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(designer, p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err updating EpicParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateEpicPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Err updating EpicParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func SetEpicPartDescription(p *EpicPart, description string) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Epic part not identified correctly. PID=0.")
	}
	stmt, err := m_db.Prepare("update EpicParts set Description=? where PID=?")
	if err != nil {
		err := fmt.Errorf("Err updating EpicParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(description, p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err updating EpicParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateEpicPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err := fmt.Errorf("Err updating EpicParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func SetEpicPartDateIssued(p *EpicPart, dateissued time.Time) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Epic part not identified correctly. PID=0.")
	}
	stmt, err := m_db.Prepare("update EpicParts set DateIssued=? where PID=?")
	if err != nil {
		err := fmt.Errorf("Err updating EpicParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(dateissued.Format("2006-01-02"), p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err updating EpicParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateEpicPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Err updating EpicParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func AddEpicPart(p *EpicPart) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		p.PID = uuid.New()
	}
	if p.DateIssued.IsZero() {
		p.DateIssued = time.Now()
	}
	ep, _ := GetEpicPart(p.PNString())
	if ep != nil {
		return fmt.Errorf("Part already exists.")
	}

	t0 := time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2040, time.January, 1, 0, 0, 0, 0, time.UTC)
	if p.DateIssued.Before(t0) || p.DateIssued.After(t1) {
		return fmt.Errorf("Date Issued (%s) is out of range.", p.DateIssued.Format("2006-01-02"))
	}

	if !IsDesigner(p.Designer) {
		return fmt.Errorf("Designer %s unknown.", p.Designer)
	}

	if !IsProject(p.ProjectId) {
		return fmt.Errorf("Project Id %s unknown.", p.ProjectId)
	}

	if !IsSubsystem(p.ProjectId, p.SubsystemId) {
		return fmt.Errorf("Subsystem %s-%s unknown.", p.ProjectId, p.SubsystemId)
	}

	stmt, err := m_db.Prepare("insert EpicParts set" +
		" PID=?," +
		" ProjectId=?," +
		" SubsystemId=?," +
		" PartType=?," +
		" SequenceNum=?," +
		" Designer=?," +
		" DateIssued=?," +
		" Description=?")
	if err != nil {
		err := fmt.Errorf("Err inserting into EpicPart. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	_, err = stmt.Exec(p.PID.String(), p.ProjectId, p.SubsystemId, p.PartType, p.SequenceNum,
		p.Designer, p.DateIssued.Format("2006-01-02"), p.Description)
	if err != nil {
		err := fmt.Errorf("Err inserting into EpicPart. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	// Now put the part into the cache.  No need to refresh directly from disk?
	gEpicPartsLock.Lock()
	defer gEpicPartsLock.Unlock()
	gEpicParts = append(gEpicParts, p)
	// Think about refreshing directly from disk?

	return nil
}

// FilterEpicParts returns a list of parts that have been filtered by the
// parameters in params. The parameters are: ProjectId, SubsystemId,
// PartType, DateBefore, DateAfter, Designer, and Description.  There
// are various aliases for each of these.  In general, if a parameter
// is present only parts that have an exact match are included in
// the return.  Dates can be in one of four formats: yyyy-mm-dd,
// yy-mm-dd, mm/dd/yy, or mm/dd/yyyy.  A match for Describtion is
// case insensitive, and any substring.
func FilterEpicParts(params map[string]string) []*EpicPart {
	mainlst := GetEpicParts()

	projectid, ok := util.MapAlias(params, "ProjectId", "projectid", "project", "proj", "prj")
	if ok {
		newlst := make([]*EpicPart, 0, len(mainlst))
		for _, p := range mainlst {
			if p.ProjectId == projectid {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	subsystemid, ok := util.MapAlias(params, "SubsystemId", "subsystemid", "subsystem", "subsys", "sub")
	if ok {
		newlst := make([]*EpicPart, 0, len(mainlst))
		for _, p := range mainlst {
			if p.SubsystemId == subsystemid {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	parttype, ok := util.MapAlias(params, "PartType", "parttype", "part")
	if ok {
		if len(parttype) > 1 {
			for _, x := range GetPartTypes() {
				if strings.ToLower(x.Description) == strings.ToLower(parttype) {
					parttype = x.Digit
				}
			}
		}
		newlst := make([]*EpicPart, 0, len(mainlst))
		for _, p := range mainlst {
			if p.PartType == parttype {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	designer, ok := util.MapAlias(params, "Designer", "designer")
	if ok {
		newlst := make([]*EpicPart, 0, len(mainlst))
		for _, p := range mainlst {
			if p.Designer == designer {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	datebefore, ok := util.MapAlias(params, "DateBefore", "datebefore", "date0", "before")
	if ok {
		newlst := make([]*EpicPart, 0, len(mainlst))
		t, err := util.ParseGenericTime(datebefore)
		if err == nil {
			for _, p := range mainlst {
				if p.DateIssued.Before(t) {
					newlst = append(newlst, p)
				}
			}
		}
		mainlst = newlst
	}

	dateafter, ok := util.MapAlias(params, "DateBefore", "datebefore", "date0", "after")
	if ok {
		newlst := make([]*EpicPart, 0, len(mainlst))
		t, err := util.ParseGenericTime(dateafter)
		if err == nil {
			for _, p := range mainlst {
				if p.DateIssued.After(t) {
					newlst = append(newlst, p)
				}
			}
		}
		mainlst = newlst
	}

	desc, ok := util.MapAlias(params, "Description", "description", "desc")
	if ok {
		desc = strings.ToLower(desc)
		newlst := make([]*EpicPart, 0, len(mainlst))
		for _, p := range mainlst {
			s := strings.ToLower(p.Description)
			i := strings.Index(s, desc)
			if i >= 0 {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	return mainlst
}
