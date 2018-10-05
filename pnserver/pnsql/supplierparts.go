// --------------------------------------------------------------------
// supplierparts.go -- Manage supplier parts
//
// Created 2018-09-27 DLB
// --------------------------------------------------------------------

package pnsql

import (
	"epic/lib/log"
	"epic/lib/util"
	"epic/lib/uuid"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SupplierPartPN struct contains a split version of a supplier part number.
type SupplierPartPN struct {
	Category    string // 2 Chars long
	SequenceNum int    // 0-999
}

type SupplierPart struct {
	PID uuid.UUID
	*SupplierPartPN
	Description string
	Vendor      string
	VendorPN    string
	WebLink     string
	Designer    string
	DateIssued  time.Time
}

var gSupplierParts []*SupplierPart
var gSupplierPartsLock sync.Mutex
var gSupplierNewPartLock sync.Mutex // When generating a new part with a new sequence number

func InvalidateSupplierPartsCache() {
	gSupplierPartsLock.Lock()
	defer gSupplierPartsLock.Unlock()
	gSupplierParts = nil
}

// StrToPns converts a string to a Supplier Part Number, if possible.
func StrToSupplierPartPN(s string) (*SupplierPartPN, error) {
	var err error
	parts := strings.Split(s, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Baddly formatted part number (%s) -- no parts.", s)
	}
	if parts[0] != "SP" {
		return nil, fmt.Errorf("Baddly formatted part number (%s) -- not SP.", s)
	}
	cat := parts[1]
	snum := parts[2]
	if len(cat) != 2 || len(snum) != 3 {
		return nil, fmt.Errorf("Baddly formatted part number (%s) -- wrong lengths.", s)
	}
	if !IsSupplierCategory(cat) {
		return nil, fmt.Errorf("Baddly formatted part number (%s). %q is not a category.", s, cat)
	}
	iseq, err := strconv.Atoi(snum)
	if err != nil {
		return nil, fmt.Errorf("Baddly formatted part number (%s). %q is not a sequence num.", s, snum)
	}
	if iseq < 0 || iseq > 999 {
		return nil, fmt.Errorf("Baddly formatted part number (%s). Sequence number out of range.", s)
	}
	return &SupplierPartPN{Category: cat, SequenceNum: iseq}, nil
}

func (p *SupplierPart) PNString() string {
	return p.SupplierPartPN.PNString()
}

func (pn *SupplierPartPN) PNString() string {
	s := "SP-" + pn.Category + "-" + fmt.Sprintf("%03d", pn.SequenceNum)
	return s
}

func (p *SupplierPart) CheckPNFormat() error {
	return p.SupplierPartPN.CheckPNFormat()
}

func (pn *SupplierPartPN) CheckPNFormat() error {
	if pn.SequenceNum < 0 || pn.SequenceNum > 999 {
		return fmt.Errorf("Sequence number (%d) out of range.", pn.SequenceNum)
	}
	_, err := StrToSupplierPartPN(pn.PNString())
	return err
}

func (p *SupplierPart) SamePN(p2 *SupplierPart) bool {
	return p.SupplierPartPN.SamePN(p2.SupplierPartPN)
}

func (pn *SupplierPartPN) SamePN(pn2 *SupplierPartPN) bool {
	if pn.Category != pn2.Category {
		return false
	}
	if pn.SequenceNum != pn2.SequenceNum {
		return false
	}
	return true
}

func (p *SupplierPart) CategoryDesc() string {
	for _, c := range GetSupplierCategories() {
		if c.Category == p.Category {
			return fmt.Sprintf("%s -- %s", c.Category, c.Description)
		}
	}
	return ""
}

func GetSupplierParts() []*SupplierPart {
	gSupplierPartsLock.Lock()
	defer gSupplierPartsLock.Unlock()
	if gSupplierParts != nil {
		return gSupplierParts
	}
	rows, err := m_db.Query(
		"Select PID, Category, SeqNum, Description, Vendor, " +
			"VendorPN, WebLink, Designer, DateIssued from SupplierParts")
	if err != nil {
		log.Errorf("Err getting supplier parts. Err=%v\n", err)
		return make([]*SupplierPart, 0, 0)
	}
	parts := make([]*SupplierPart, 0, 10000)
	for rows.Next() {
		var spid, cat, desc, vendor, vendorpn, weblink, designer, sdate string
		var seqnum int
		err = rows.Scan(&spid, &cat, &seqnum, &desc, &vendor, &vendorpn, &weblink, &designer, &sdate)
		if err != nil {
			log.Errorf("Err during scan of SupplierParts. Err=%v", err)
			continue
		}
		id, err := uuid.FromString(spid)
		if err != nil {
			log.Errorf("Bad uuid for SupplierParts in database. (%s) Err=%v", spid, err)
			continue
		}
		dateissued, err := time.Parse("2006-01-02", sdate)
		if err != nil {
			log.Errorf("Bad date for SupplierParts in database. (%s) Err=%v", spid, err)
			continue
		}
		p := &SupplierPart{PID: id, Description: desc, Vendor: vendor, VendorPN: vendorpn,
			WebLink: weblink, Designer: designer, DateIssued: dateissued}
		p.SupplierPartPN = &SupplierPartPN{Category: cat, SequenceNum: seqnum}
		parts = append(parts, p)
	}
	gSupplierParts = parts
	return gSupplierParts
}

// Retrieve an supplier part given its part number string (SP-cc-0000)
// Note, it is valid to return nil for the part without an error, which
// means the part doesn't exist!
func GetSupplierPart(pns string) (*SupplierPart, error) {
	pn, err := StrToSupplierPartPN(pns)
	if err != nil {
		return nil, err
	}
	for _, p := range GetSupplierParts() {
		if p.SupplierPartPN.SamePN(pn) {
			return p, nil
		}
	}
	return nil, nil
}

func DeleteSupplierPart(p *SupplierPart) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Supplier part not identified correctly. PID=0.")
	}
	stmt, err := m_db.Prepare("delete from SupplierParts where PID=?")
	if err != nil {
		log.Errorf("Err deleting from SupplierParts. Err=%v", err)
		return err
	}
	r, err := stmt.Exec(p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err deleting from SupplierParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateSupplierPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Err deleting from SupplierParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	log.Infof("Supplier Part %s deleted from database.", p.PNString())
	return nil
}

func SetSupplierPartDesigner(p *SupplierPart, designer string) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Supplier part not identified correctly. PID=0.")
	}
	if !IsDesigner(designer) {
		return fmt.Errorf("%q is not a known designer.")
	}
	olddesigner := p.Designer
	stmt, err := m_db.Prepare("update SupplierParts set Designer=? where PID=?")
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(designer, p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateSupplierPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Err updating SupplierParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	log.Infof("Supplier Part %s: Designer changed from %s to %s.", p.PNString(), olddesigner, designer)
	return nil
}

func SetSupplierPartVendor(p *SupplierPart, vendor string) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Supplier part not identified correctly. PID=0.")
	}
	oldvendor := p.Vendor
	stmt, err := m_db.Prepare("update SupplierParts set Vendor=? where PID=?")
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(vendor, p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateSupplierPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Err updating SupplierParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	log.Infof("Supplier Part %s: Vendor changed from %s to %s.", p.PNString(), oldvendor, vendor)
	return nil
}

func SetSupplierPartVendorPN(p *SupplierPart, vendorpn string) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Supplier part not identified correctly. PID=0.")
	}
	oldvendorpn := p.VendorPN
	stmt, err := m_db.Prepare("update SupplierParts set VendorPN=? where PID=?")
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(vendorpn, p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateSupplierPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Err updating SupplierParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	log.Infof("Supplier Part %s: VendorPN changed from %s to %s.", p.PNString(), oldvendorpn, vendorpn)
	return nil
}

func SetSupplierPartWebLink(p *SupplierPart, weblink string) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Supplier part not identified correctly. PID=0.")
	}
	stmt, err := m_db.Prepare("update SupplierParts set WebLink=? where PID=?")
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(weblink, p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateSupplierPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Err updating SupplierParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	log.Infof("Supplier Part %s: WebLink changed to %q.", p.PNString(), weblink)
	return nil
}

func SetSupplierPartDescription(p *SupplierPart, description string) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Supplier part not identified correctly. PID=0.")
	}
	stmt, err := m_db.Prepare("update SupplierParts set Description=? where PID=?")
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(description, p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateSupplierPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err := fmt.Errorf("Err updating SupplierParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	log.Infof("Supplier Part %s: Description changed to %q.", p.PNString(), description)
	return nil
}

func SetSupplierPartDateIssued(p *SupplierPart, dateissued time.Time) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		return fmt.Errorf("Supplier part not identified correctly. PID=0.")
	}
	olddate := p.DateIssued
	stmt, err := m_db.Prepare("update SupplierParts set DateIssued=? where PID=?")
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(dateissued.Format("2006-01-02"), p.PID.String())
	if err != nil {
		err := fmt.Errorf("Err updating SupplierParts. Err= %v", err)
		log.Errorf("%v", err)
		return err
	}
	InvalidateSupplierPartsCache()
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Err updating SupplierParts. Wrong affected row count (%d)", n)
		log.Errorf("%v", err)
		return err
	}
	log.Infof("Supplier Part %s: DateIssued changed from %s to: %s.", p.PNString(), olddate.Format("2006-01-02"))
	return nil
}

func NewSupplierPartInSequence(Designer, Category, Vendor, VendorPN, WebLink, Description string) (*SupplierPart, error) {
	if !IsDesigner(Designer) {
		return nil, fmt.Errorf("Designer %s unknown.", Designer)
	}
	if !IsSupplierCategory(Category) {
		return nil, fmt.Errorf("Category %s unknown.", Category)
	}
	if util.Blank(Description) {
		return nil, fmt.Errorf("Description cannot be blank.")
	}

	gSupplierNewPartLock.Lock()
	defer gSupplierNewPartLock.Unlock()
	InvalidateSupplierPartsCache()
	lastseq := 0
	for _, part := range GetSupplierParts() {
		if part.Category != Category {
			continue
		}
		if part.SequenceNum > lastseq {
			lastseq = part.SequenceNum
		}
	}
	p := &SupplierPart{}
	p.SupplierPartPN = &SupplierPartPN{}

	p.PID = uuid.New()
	p.Category = Category
	p.SequenceNum = lastseq + 1
	p.Description = Description
	p.Vendor = Vendor
	p.VendorPN = VendorPN
	p.WebLink = WebLink
	p.Designer = Designer
	p.DateIssued = time.Now()
	err := write_supplier_part(p)
	InvalidateSupplierPartsCache()
	if err != nil {
		log.Infof("NEW Supplier Part: %s added in sequence.", p.PNString())
	}
	return p, err
}

func AddSupplierPart(p *SupplierPart) error {
	if err := p.CheckPNFormat(); err != nil {
		return err
	}
	if p.PID.IsZero() {
		p.PID = uuid.New()
	}
	if p.DateIssued.IsZero() {
		p.DateIssued = time.Now()
	}
	ep, _ := GetSupplierPart(p.PNString())
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

	if !IsSupplierCategory(p.Category) {
		return fmt.Errorf("Supplier category %s unknown.", p.Category)
	}

	if util.Blank(p.Description) {
		return fmt.Errorf("Description cannot be blank.")
	}
	err := write_supplier_part(p)
	// Now put the part into the cache.  No need to refresh directly from disk?
	gSupplierPartsLock.Lock()
	defer gSupplierPartsLock.Unlock()
	if err != nil {
		InvalidateSupplierPartsCache()
	} else {
		gSupplierParts = append(gSupplierParts, p)
	}
	if err != nil {
		log.Infof("NEW Supplier Part: %s added out-of-sequence.", p.PNString())
	}
	return nil
}

func write_supplier_part(p *SupplierPart) error {
	stmt, err := m_db.Prepare("insert SupplierParts set" +
		" PID=?," +
		" Category=?," +
		" SeqNum=?," +
		" Description=?," +
		" Vendor=?," +
		" VendorPN=?," +
		" WebLink=?," +
		" Designer=?," +
		" DateIssued=?")
	if err != nil {
		err := fmt.Errorf("Err inserting into SupplierParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	r, err := stmt.Exec(p.PID.String(), p.Category, p.SequenceNum,
		p.Description, p.Vendor, p.VendorPN, p.WebLink,
		p.Designer, p.DateIssued.Format("2006-01-02"))
	if err != nil {
		err := fmt.Errorf("Err inserting into SupplierParts. Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	if n, _ := r.RowsAffected(); n != 1 {
		log.Errorf("Wrong number of affected rows (%d) on insert of supplier part.", n)
	}
	return nil
}

// FilterSupplierParts returns a list of parts that have been filtered by the
// parameters in params. The parameters are: Category, Description, Vendor,
// VendorPN, WebLink, Designer, DateBefore, DateAfter. There
// are various aliases for each of these.  If Category and/or Designer is
// present only parts that have an exact match are included in
// the return.  However, for Description, Vendor, VendorPN, and WebLink,
// a match is declared for any case insensitive substring.  Dates can be in
// a generic format.
func FilterSupplierParts(params map[string]string) []*SupplierPart {
	mainlst := GetSupplierParts()

	category, ok := util.MapAlias(params, "Category", "category", "Cat", "cat")
	if ok {
		if len(category) > 2 {
			for _, x := range GetSupplierCategories() {
				if strings.ToLower(x.Description) == strings.ToLower(category) {
					category = x.Category
					break
				}
			}
		}
		newlst := make([]*SupplierPart, 0, len(mainlst))
		for _, p := range mainlst {
			if p.Category == category {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	desc, ok := util.MapAlias(params, "Description", "description", "desc")
	if ok {
		desc = strings.ToLower(desc)
		newlst := make([]*SupplierPart, 0, len(mainlst))
		for _, p := range mainlst {
			s := strings.ToLower(p.Description)
			i := strings.Index(s, desc)
			if i >= 0 {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	vendor, ok := util.MapAlias(params, "Vendor", "vendor", "ven")
	if ok {
		vendor = strings.ToLower(vendor)
		newlst := make([]*SupplierPart, 0, len(mainlst))
		for _, p := range mainlst {
			s := strings.ToLower(p.Vendor)
			i := strings.Index(s, vendor)
			if i >= 0 {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	vendorpn, ok := util.MapAlias(params, "VendorPN", "vendorpn", "vpn")
	if ok {
		vendorpn = strings.ToLower(vendorpn)
		newlst := make([]*SupplierPart, 0, len(mainlst))
		for _, p := range mainlst {
			s := strings.ToLower(p.VendorPN)
			i := strings.Index(s, vendorpn)
			if i >= 0 {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	weblink, ok := util.MapAlias(params, "WebLink", "Weblink", "weblink")
	if ok {
		weblink = strings.ToLower(weblink)
		newlst := make([]*SupplierPart, 0, len(mainlst))
		for _, p := range mainlst {
			s := strings.ToLower(p.WebLink)
			i := strings.Index(s, weblink)
			if i >= 0 {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	designer, ok := util.MapAlias(params, "Designer", "designer")
	if ok {
		newlst := make([]*SupplierPart, 0, len(mainlst))
		for _, p := range mainlst {
			if p.Designer == designer {
				newlst = append(newlst, p)
			}
		}
		mainlst = newlst
	}

	datebefore, ok := util.MapAlias(params, "DateBefore", "datebefore", "date0", "before")
	if ok {
		newlst := make([]*SupplierPart, 0, len(mainlst))
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
		newlst := make([]*SupplierPart, 0, len(mainlst))
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

	return mainlst
}

// GetVendors returns a list of known vendors, alphbized.
func GetVendors() []string {
	vmap := make(map[string]int, 200)
	for _, p := range GetSupplierParts() {
		v := strings.TrimSpace(p.Vendor)
		if util.Blank(v) {
			continue
		}
		vmap[v]++
	}
	lst := make([]string, 0, len(vmap))
	for k, _ := range vmap {
		lst = append(lst, k)
	}
	sort.Strings(lst)
	return lst
}
