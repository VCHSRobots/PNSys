// --------------------------------------------------------------------
// context.go -- Context holder for commands
//
// Created 2018-10-03 DLB
// --------------------------------------------------------------------

package console

import (
	"bytes"
	"fmt"
)

type ContextType string

const (
	Context_Internal ContextType = "Internal"
	Context_External ContextType = "External"
)

type Context struct {
	mode    ContextType
	outdata *bytes.Buffer
}

func NewContext(mode ContextType) *Context {
	c := &Context{}
	c.mode = mode
	c.outdata = new(bytes.Buffer)
	return c
}

func (c *Context) Printf(f string, args ...interface{}) {
	fmt.Fprintf(c.outdata, f, args...)
}

func (c *Context) Output() string {
	return string(c.outdata.Bytes())
}

func (c *Context) IsInternal() bool {
	return c.mode == Context_Internal
}

func (c *Context) IsExternal() bool {
	return c.mode == Context_External
}
