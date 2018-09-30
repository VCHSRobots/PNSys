// --------------------------------------------------------------------
// main_page.go -- Shows the main page
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package pages

// import (
// 	"bytes"
// 	"epic/lib/log"
// 	//"epic/lib/uuid"
// 	//"epic/pnserver/pnsql"
// 	"github.com/gin-gonic/gin"
// 	"html/template"
// 	"io/ioutil"
// )

// type MainPageData struct {
// 	PageTitle  string
// 	PageHeader string
// }

// func show_main_page(c *gin.Context) {
// 	// if !override {
// 	//     cookie, err := c.Cookie("Cred")
// 	//     if err != nil {
// 	//         err.Errorf("Unable to get Cred. Err=%v", err)
// 	//         show_login_page(c, false)
// 	//         return
// 	//     }
// 	//     id, err := uuid.FromString(cookie)
// 	//     if err != nil {
// 	//         log.Warnf("UUID parse failed.")
// 	//         show_login_page(c, false)
// 	//         return
// 	//     }
// 	//     ss, err = security.GetSessionByAuth(id)
// 	//     if ss == nil || err != nil {
// 	//         log.Warnf("Attempt to access main page while not being logged in.")
// 	//         show_login_page(c, false)
// 	//         return
// 	//     }
// 	// }

// 	tmpl_name := "./static/templates/main.tmpl"
// 	tmpl_bytes, err := ioutil.ReadFile(tmpl_name)
// 	if err != nil {
// 		log.Errorf("Missing template %s. Err=%v", tmpl_name, err)
// 		c.AbortWithError(400, err)
// 		return
// 	}
// 	tmpl, err := template.New("Main").Parse(string(tmpl_bytes))
// 	if err != nil {
// 		log.Errorf("Invalid template %s. Err=%v", tmpl_name, err)
// 		c.AbortWithError(400, err)
// 		return
// 	}

// 	data := LoginData{PageHeader: "Main Page"}
// 	data.PageTitle = "PNSys"
// 	html := new(bytes.Buffer)
// 	tmpl.Execute(html, data)
// 	c.Data(200, "text/html", html.Bytes())
// }
