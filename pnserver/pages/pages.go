// --------------------------------------------------------------------
// pages.go -- Handle all pages
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

import (
	"github.com/gin-gonic/gin"
)

type InvokeType string

const (
	Invoke_GET  = "GET"
	Invoke_POST = "POST"
)

type Page struct {
	Route   string
	Invoke  InvokeType
	Handler func(c *gin.Context)
}

var gPages []*Page

// Registor a web page
func RegisterPage(route string, invoke InvokeType, handler func(c *gin.Context)) {
	if gPages == nil {
		gPages = make([]*Page, 0, 50)
	}
	p := &Page{route, invoke, handler}
	gPages = append(gPages, p)
}

func GetAllPages() []*Page {
	return gPages
}
