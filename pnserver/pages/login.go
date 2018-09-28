// --------------------------------------------------------------------
// login.go -- Handles Login Page
//
// Created 2018-09-22 DLB
// --------------------------------------------------------------------

package pages

import (
	"bytes"
	"epic/lib/log"
	"epic/pnserver/pnsql"
	"epic/pnserver/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
)

type LoginData struct {
	PageTitle   string
	PageHeader  string
	Designers   []string
	BadPassword bool
}

func Login(c *gin.Context) {
	show_login_page(c, false)
}

func show_login_page(c *gin.Context, badpassword bool) {
	tmpl_name := "./static/templates/login.tmpl"
	tmpl_bytes, err := ioutil.ReadFile(tmpl_name)
	if err != nil {
		log.Errorf("Missing template %s. Err=%v", tmpl_name, err)
		c.AbortWithError(400, err)
		return
	}
	tmpl, err := template.New("Login").Parse(string(tmpl_bytes))
	if err != nil {
		log.Errorf("Invalid template %s. Err=%v", tmpl_name, err)
		c.AbortWithError(400, err)
		return
	}

	data := LoginData{PageTitle: "PNSys", PageHeader: "Login"}
	data.BadPassword = badpassword
	dlst := pnsql.GetDesigners()
	data.Designers = make([]string, 0, len(dlst))
	for _, d := range dlst {
		data.Designers = append(data.Designers, d.Name)
	}
	html := new(bytes.Buffer)
	tmpl.Execute(html, data)
	c.Data(200, "text/html", html.Bytes())
}

type TLoginData struct {
	Name     string `form:"name"`
	Password string `form:"password"`
}

func LoginPost(c *gin.Context) {
	var ld TLoginData
	err := c.ShouldBind(&ld)
	if err != nil {
		log.Errorf("Login data failed to bind. Err=%v", err)
		c.AbortWithError(400, err)
		return
	}
	if ld.Name == "admin" {
		if ld.Password != "loveepic" {
			log.Infof("Attempt to login to admin fails.")
			show_login_page(c, true)
			return
		}
		ses := sessions.NewSession("admin", c.ClientIP())
		c.SetCookie("Cred", ses.AuthCookie.String(), 0, "/", "", false, true)
		show_main_page(c)
		return
	}
	if ld.Name == "guest" {
		ses := sessions.NewSession("guest", c.ClientIP())
		c.SetCookie("Cred", ses.AuthCookie.String(), 0, "/", "", false, true)
		show_main_page(c)
		return
	}
	if ld.Password != "epic4fun" {
		show_login_page(c, true)
		return
	}

	haveit := false
	designers := pnsql.GetDesigners()
	for _, d := range designers {
		if d.Name == ld.Name {
			haveit = true
			break
		}
	}
	if !haveit {
		show_login_page(c, true)
		return
	}
	ses := sessions.NewSession(ld.Name, c.ClientIP())
	c.SetCookie("Cred", ses.AuthCookie.String(), 0, "/", "", false, true)
	show_main_page(c)
}
