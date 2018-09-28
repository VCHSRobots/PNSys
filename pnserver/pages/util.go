// --------------------------------------------------------------------
// util.go -- Utilities for pages and template processing.
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

import (
	"bytes"
	"epic/lib/log"
	"io/ioutil"
	"text/template"
)

func GetTemplate(name string) (*template.Template, error) {
	fn := "./static/templates/" + name + ".tmpl"
	t_bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Errorf("Missing template %q. (file %s). Err=%v",
			name, fn, err)
		return nil, err
	}
	tmpl, err := template.New(name).Parse(string(t_bytes))
	if err != nil {
		log.Errorf("Invalid template %q. Err=%v", name, err)
		return nil, err
	}
	return tmpl, nil
}

func MakePage(data interface{}, template_names ...string) ([]byte, error) {

	tmpls := make([]*template.Template, 0, len(template_names))
	for _, n := range template_names {
		t, err := GetTemplate(n)
		if err != nil {
			return []byte{}, err
		}
		tmpls = append(tmpls, t)
	}
	html := new(bytes.Buffer)
	for i, t := range tmpls {
		err := t.Execute(html, data)
		if err != nil {
			log.Errorf("Error execution template %q. Err=%v", template_names[i], err)
			return html.Bytes(), err
		}
	}
	return html.Bytes(), nil
}
