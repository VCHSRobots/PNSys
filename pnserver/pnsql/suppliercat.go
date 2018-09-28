// --------------------------------------------------------------------
// suppliercat.go -- Manage the Supplier Category Table
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package pnsql

import (
	"epic/lib/log"
	"sync"
)

type SupplierCategory struct {
	Category    string
	Description string
}

var gSupplierCategories []*SupplierCategory
var gSupplierCategoriesLock sync.Mutex

func InvalidateSupplierCategoryCache() {
	gSupplierCategoriesLock.Lock()
	defer gSupplierCategoriesLock.Unlock()
	gSupplierCategories = nil
}

func GetSupplierCategories() []*SupplierCategory {
	gSupplierCategoriesLock.Lock()
	defer gSupplierCategoriesLock.Unlock()
	if gSupplierCategories != nil {
		return gSupplierCategories
	}
	lst := make([]*SupplierCategory, 0, 11)
	rows, err := m_db.Query("Select Category, Description from SupplierCategory")
	if err != nil {
		log.Errorf("Err getting SupplierCategory. Returning null slice. Err=%v\n", err)
		return lst
	}
	for rows.Next() {
		var cat string
		var desc string
		err = rows.Scan(&cat, &desc)
		if err != nil {
			log.Errorf("Err during row scan in GetSupplierCategories. Err=%v\n", err)
			continue
		}
		lst = append(lst, &SupplierCategory{cat, desc})
	}
	gSupplierCategories = lst
	return gSupplierCategories
}

func IsSupplierCategory(cat string) bool {
	for _, c := range GetSupplierCategories() {
		if c.Category == cat {
			return true
		}
	}
	return false
}
