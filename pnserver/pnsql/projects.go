// --------------------------------------------------------------------
// projects.go -- Manage projects and subassemblies
//
// Created 2018-09-24 DLB
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

type Subsystem struct {
	ProjectId   string // Three chars
	SubsystemId string // Two chars
	Description string
}

type Project struct {
	ProjectId   string
	Description string
	Year0       string
	Active      bool
	Subsystems  []*Subsystem
}

var gProjectCache []*Project
var gLockProjects sync.Mutex

func InvalidateProjectsCache() {
	gLockProjects.Lock()
	defer gLockProjects.Unlock()
	gProjectCache = nil
}

func GetProjects() []*Project {
	gLockProjects.Lock()
	defer gLockProjects.Unlock()
	if gProjectCache != nil {
		return gProjectCache
	}
	prjs := make(map[string]*Project, 200)
	rows, err := m_db.Query("Select ProjectId, Description, Year0, Active from Projects")
	if err != nil {
		log.Errorf("Err getting Projects. Returning null slice. Err=%v\n", err)
		return []*Project{}
	}
	for rows.Next() {
		var id string
		var desc string
		var year0 string
		var active int
		err = rows.Scan(&id, &desc, &year0, &active)
		if err != nil {
			log.Errorf("Err during row scan in GetPartTypes. Err=%v\n", err)
			continue
		}
		bactive := active != 0
		prjs[id] = &Project{ProjectId: id, Description: desc, Year0: year0, Active: bactive}
		prjs[id].Subsystems = make([]*Subsystem, 0, 20)
	}

	rows, err = m_db.Query("Select ProjectId, SubsystemId, Description from Subsystems")
	if err != nil {
		log.Errorf("Err getting Subsystems. Returning empty proejcts. Err=%v", err)
		return []*Project{}
	}
	for rows.Next() {
		var prjid string
		var subsystemid string
		var desc string
		err = rows.Scan(&prjid, &subsystemid, &desc)
		if err != nil {
			log.Errorf("Err during row scan in GetPartTypes. Err=%v\n", err)
			continue
		}
		prj, ok := prjs[prjid]
		if !ok {
			log.Warnf("Found orphined project (%s) in SubSystem table.\n", prjid)
			continue
		}
		subsys := &Subsystem{ProjectId: prjid, SubsystemId: subsystemid, Description: desc}
		prj.Subsystems = append(prj.Subsystems, subsys)
		sortfn1 := func(i, j int) bool {
			return prj.Subsystems[i].SubsystemId < prj.Subsystems[j].SubsystemId
		}
		sort.Slice(prj.Subsystems, sortfn1)
	}

	prjlst := make([]*Project, 0, len(prjs))
	for _, v := range prjs {
		prjlst = append(prjlst, v)
	}
	sortfnc := func(i, j int) bool {
		if prjlst[i].Active == prjlst[j].Active {
			if prjlst[i].Year0 == prjlst[j].Year0 {
				return prjlst[i].ProjectId < prjlst[j].ProjectId
			}
			return prjlst[i].Year0 > prjlst[j].Year0
		}
		return prjlst[i].Active
	}
	sort.Slice(prjlst, sortfnc)
	gProjectCache = prjlst
	return gProjectCache
}

func (p *Project) GetSelString() string {
	ss := p.ProjectId + " -- " + p.Description
	ss = util.FixStrLen(ss, 35, " ")
	ss = strings.TrimSpace(ss)
	return ss
}

func (p *Project) GetSubSysSelList() []string {
	lst := make([]string, 0, len(p.Subsystems))
	for _, st := range p.Subsystems {
		ss := st.SubsystemId + " -- " + st.Description
		ss = util.FixStrLen(ss, 35, " ")
		ss = strings.TrimSpace(ss)
		lst = append(lst, ss)
	}
	return lst
}

func IsProject(ProjectId string) bool {
	for _, p := range GetProjects() {
		if p.ProjectId == ProjectId {
			return true
		}
	}
	return false
}

func IsSubsystem(ProjectId, SubsystemId string) bool {
	for _, p := range GetProjects() {
		if p.ProjectId == ProjectId {
			for _, s := range p.Subsystems {
				if s.SubsystemId == SubsystemId {
					return true
				}
			}
			return false
		}
	}
	return false
}

func IsPartType(PartType string) bool {
	for _, t := range GetPartTypes() {
		if t.Digit == PartType {
			return true
		}
	}
	return false
}

func GetProject(ProjectId string) (*Project, error) {
	for _, p := range GetProjects() {
		if p.ProjectId == ProjectId {
			return p, nil
		}
	}
	return nil, fmt.Errorf("No such project (%s).", ProjectId)
}

func GetSubsystem(ProjectId, SubsystemId string) (*Subsystem, error) {
	for _, p := range GetProjects() {
		if p.ProjectId == ProjectId {
			for _, s := range p.Subsystems {
				if s.SubsystemId == SubsystemId {
					return s, nil
				}
			}
			break
		}
	}
	return nil, fmt.Errorf("No such subsystem (%s-%s).", ProjectId, SubsystemId)
}

func CheckProjectIdText(ProjectId string) error {
	s := strings.ToLower(ProjectId)
	if !util.ContainsOnly(s, "abcdefghijklmnopqrstuvwxyz0123456789") {
		err := fmt.Errorf("The Project Id cannot contain punctuation characters or spaces.")
		return err
	}
	if len(s) != 3 {
		err := fmt.Errorf("The Project id (%q) must be exactly three characters long.", ProjectId)
		return err
	}
	return nil
}

func CheckSubsystemIdText(SubsystemId string) error {
	s := strings.ToLower(SubsystemId)
	if !util.ContainsOnly(s, "abcdefghijklmnopqrstuvwxyz0123456789") {
		err := fmt.Errorf("The Subsystem Id cannot contain punctuation characters or spaces.")
		return err
	}
	if len(s) != 2 {
		err := fmt.Errorf("The Subsystem id (%q) must be exactly two characters long.", SubsystemId)
		return err
	}
	return nil
}

func SplitProjectId(text string) (ProjectId, SubsystemId string, err error) {
	wrds := strings.Split(text, "-")
	if len(wrds) == 1 {
		ProjectId = wrds[0]
		SubsystemId = ""
		err := CheckProjectIdText(ProjectId)
		if err != nil {
			return "", "", fmt.Errorf("%q is an illegal Project/Subsystem Id (1: %q).", text, ProjectId)
		}
		return ProjectId, SubsystemId, nil
	}
	if len(wrds) == 2 {
		ProjectId = wrds[0]
		SubsystemId = wrds[1]
		err := CheckProjectIdText(ProjectId)
		if err != nil {
			return "", "", fmt.Errorf("%q is an illegal Project/Subsystem Id (2: %q).", text, ProjectId)
		}
		err = CheckSubsystemIdText(SubsystemId)
		if err != nil {
			return "", "", fmt.Errorf("%q is an illegal Project/Subsystem Id (3: %q).", text, SubsystemId)
		}
		return ProjectId, SubsystemId, nil
	}
	return "", "", fmt.Errorf("%q is an illegal Project/Subsystem Id (4).", text)
}

func AddProject(ProjectId, Description, Year0 string) error {
	InvalidateProjectsCache()
	err := CheckProjectIdText(ProjectId)
	if err != nil {
		return err
	}

	if IsProject(ProjectId) {
		err = fmt.Errorf("Project %s already exists.", ProjectId)
		return err
	}

	if util.Blank(Description) {
		err = fmt.Errorf("Description must be provided and cannot be blank.")
		return err
	}

	err = CheckYear0Text(Year0)
	if err != nil {
		return err
	}

	stmt, err := m_db.Prepare("insert Projects set ProjectId=?, Description=?, Year0=?, Active=1")
	if err != nil {
		err = fmt.Errorf("Err inserting into Projects. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(ProjectId, Description, Year0)
	InvalidateProjectsCache()
	if err != nil {
		err = fmt.Errorf("Err inserting into Projects. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err := fmt.Errorf("Wrong number of rows affected (%d) on add of Project %s.", n, ProjectId)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func DeleteProject(ProjectId string) error {
	InvalidateProjectsCache()
	prj, err := GetProject(ProjectId)
	if err != nil || prj == nil {
		err = fmt.Errorf("Project %s does not exist.", ProjectId)
		return err
	}

	// Here, we should consider locking creation of subsystems and parts.

	if len(prj.Subsystems) > 0 {
		err = fmt.Errorf("Cannot delete project %s, as there are %d subsystems associated with it.",
			prj.ProjectId, len(prj.Subsystems))
		return err
	}

	params := make(map[string]string, 2)
	params["ProjectId"] = prj.ProjectId
	parts := FilterEpicParts(params)
	if len(parts) > 0 {
		err = fmt.Errorf("Cannot delete project %s, as there are %d parts in it.",
			prj.ProjectId, len(parts))
		return err
	}

	// Ok to delete, proceed.
	stmt, err := m_db.Prepare("delete from Projects where ProjectId=?")
	if err != nil {
		err = fmt.Errorf("Err deleting from Projects. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(ProjectId)
	InvalidateProjectsCache()
	if err != nil {
		err = fmt.Errorf("Err deleting from Projects. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Wrong number of rows affected (%d) when deleting project %s.",
			n, ProjectId)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func SetProjectDescription(ProjectId, Description string) error {
	InvalidateProjectsCache()
	prj, err := GetProject(ProjectId)
	if err != nil || prj == nil {
		return fmt.Errorf("Project %s does not exist.", ProjectId)
	}
	if util.Blank(prj.ProjectId) {
		return fmt.Errorf("Epic project not identified correctly. ProjectId is blank.")
	}
	if prj.Description == Description {
		return fmt.Errorf("No change to database requested. Description equal.")
	}
	if util.Blank(Description) {
		return fmt.Errorf("Description cannot be blank.")
	}
	stmt, err := m_db.Prepare("update Projects set Description=? where ProjectId=?")
	if err != nil {
		err := fmt.Errorf("Err updating Projects. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(Description, prj.ProjectId)
	InvalidateProjectsCache()
	if err != nil {
		err := fmt.Errorf("Err updating Projects. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err := fmt.Errorf("Err updating Projects. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func SetProjectYear0(ProjectId, Year0 string) error {
	InvalidateProjectsCache()
	prj, err := GetProject(ProjectId)
	if err != nil || prj == nil {
		return fmt.Errorf("Project %s does not exist.", ProjectId)
	}
	if util.Blank(prj.ProjectId) {
		return fmt.Errorf("Epic project not identified correctly. ProjectId is blank.")
	}
	err = CheckYear0Text(Year0)
	if err != nil {
		return err
	}
	if prj.Year0 == Year0 {
		return fmt.Errorf("No change to database requested. Year0 equal.")
	}
	stmt, err := m_db.Prepare("update Projects set Year0=? where ProjectId=?")
	if err != nil {
		err := fmt.Errorf("Err updating Projects. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(Year0, prj.ProjectId)
	InvalidateProjectsCache()
	if err != nil {
		err := fmt.Errorf("Err updating Projects. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err := fmt.Errorf("Err updating Projects. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func SetProjectActive(ProjectId string, Active bool) error {
	InvalidateProjectsCache()
	prj, err := GetProject(ProjectId)
	if err != nil || prj == nil {
		return fmt.Errorf("Project %s does not exist.", ProjectId)
	}
	if util.Blank(prj.ProjectId) {
		return fmt.Errorf("Epic project not identified correctly. ProjectId is blank.")
	}
	if prj.Active == Active {
		return fmt.Errorf("No change to database requested. Active equal.")
	}
	iactive := 0
	if Active {
		iactive = 1
	}
	stmt, err := m_db.Prepare("update Projects set Active=? where ProjectId=?")
	if err != nil {
		err := fmt.Errorf("Err updating Projects. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(iactive, prj.ProjectId)
	InvalidateProjectsCache()
	if err != nil {
		err := fmt.Errorf("Err updating Projects. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err := fmt.Errorf("Err updating Projects. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func AddSubsystem(ProjectId, SubsystemId, Description string) error {
	InvalidateProjectsCache()
	prj, err := GetProject(ProjectId)
	if err != nil || prj == nil {
		return fmt.Errorf("Project %s does not exist.", ProjectId)
	}
	for _, sub := range prj.Subsystems {
		if sub.SubsystemId == SubsystemId {
			return fmt.Errorf("Subsystem %s-%s already exists.", ProjectId, SubsystemId)
		}
	}
	err = CheckSubsystemIdText(SubsystemId)
	if err != nil {
		return err
	}

	if util.Blank(Description) {
		err = fmt.Errorf("Description must be provided and cannot be blank.")
		return err
	}

	stmt, err := m_db.Prepare("insert Subsystems set ProjectId=?, SubsystemId=?, Description=?")
	if err != nil {
		log.Errorf("Err inserting into Subsystems. Err=%v", err)
		return err
	}
	_, err = stmt.Exec(ProjectId, SubsystemId, Description)
	InvalidateProjectsCache()
	if err != nil {
		log.Errorf("Err inserting into Subsystems. Err=%v", err)
		return err
	}
	return nil
}

func DeleteSubsystem(ProjectId, SubsystemId string) error {
	InvalidateProjectsCache()
	prj, err := GetProject(ProjectId)
	if err != nil || prj == nil {
		return fmt.Errorf("Project %s does not exist.", ProjectId)
	}
	haveit := false
	for _, sub := range prj.Subsystems {
		if sub.SubsystemId == SubsystemId {
			haveit = true
		}
	}
	if !haveit {
		return fmt.Errorf("Subsystem %s-%s does not exist.", ProjectId, SubsystemId)
	}

	params := make(map[string]string, 2)
	params["ProjectId"] = prj.ProjectId
	params["SubsystemId"] = SubsystemId
	parts := FilterEpicParts(params)
	if len(parts) > 0 {
		err = fmt.Errorf("Cannot delete subsystem %s-%s, as there are %d parts in it.",
			prj.ProjectId, SubsystemId, len(parts))
		return err
	}

	// Ok to delete, proceed.
	stmt, err := m_db.Prepare("delete from Subsystems where ProjectId=? and SubsystemId=?")
	if err != nil {
		err = fmt.Errorf("Err deleting from Subsystems. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(ProjectId, SubsystemId)
	InvalidateProjectsCache()
	if err != nil {
		err = fmt.Errorf("Err deleting from Subsystems. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Wrong number of rows affected (%d) when deleting subsystem %s-%s.",
			n, ProjectId, SubsystemId)
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func SetSubsystemDescription(ProjectId, SubsystemId, Description string) error {
	InvalidateProjectsCache()
	prj, err := GetProject(ProjectId)
	if err != nil || prj == nil {
		return fmt.Errorf("Project %s does not exist.", ProjectId)
	}
	var subsys *Subsystem
	for _, sub := range prj.Subsystems {
		if sub.SubsystemId == SubsystemId {
			subsys = sub
		}
	}
	if subsys == nil {
		return fmt.Errorf("Subsystem %s-%s does not exist.", ProjectId, SubsystemId)
	}
	if subsys.Description == Description {
		return fmt.Errorf("No change to database requested. Description equal.")
	}
	if util.Blank(Description) {
		return fmt.Errorf("Description cannot be blank.")
	}
	stmt, err := m_db.Prepare("update Subystems set Description=? where ProjectId=? and SubsystemId=?")
	if err != nil {
		err := fmt.Errorf("Err updating Subsystems. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(Description, prj.ProjectId, subsys.SubsystemId)
	InvalidateProjectsCache()
	if err != nil {
		err := fmt.Errorf("Err updating Subsystems. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err := fmt.Errorf("Err updating Subsystems. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	return nil
}
