// --------------------------------------------------------------------
// login.go -- Handles Login Page
//
// Created 2018-09-22 DLB
// --------------------------------------------------------------------

package pages

import (
	"encoding/json"
	"epic/lib/log"
	"epic/pnserver/pnsql"
	pv "epic/pnserver/privilege"
	"epic/pnserver/sessions"
	"github.com/gin-gonic/gin"
)

type LoginData struct {
	*HeaderData
	DesignersJson string
}

func init() {
	RegisterPage("/Login", Invoke_GET, handle_login)
	RegisterPage("/LoginPost", Invoke_POST, handle_login_post)
}

func make_default_login_data() *LoginData {
	data := &LoginData{}
	data.HeaderData = &HeaderData{}
	data.PageTabTitle = "PnSys Login"
	data.OnLoadFuncJS = "startUp"
	data.HideLoginLink = true
	return data
}

func handle_login(c *gin.Context) {
	kill_session(c)
	data := make_default_login_data()
	show_login_page(c, data)
}

type LoginSubmitData struct {
	Designer string `form:"Designer"`
	Password string `form:"Password"`
}

func handle_login_post(c *gin.Context) {
	kill_session(c)
	data := make_default_login_data()

	var ld LoginSubmitData
	err := c.ShouldBind(&ld)
	if err != nil {
		log.Errorf("Login data failed to bind. Err=%v", err)
		data.ErrorMessage = "Web app failure!  Programming problem?"
		show_login_page(c, data)
		return
	}
	if !pnsql.IsDesigner(ld.Designer) {
		if ld.Designer != "" {
			log.Infof("Attempt to login with unknown designer (%q). Hacking attempt?", ld.Designer)
		} else {
			log.Infof("Attempt to login with blank designer.")
		}
		data.ErrorMessage = "Login Failed."
		show_login_page(c, data)
		return
	}

	priv, ok := sessions.CheckPassword(ld.Designer, ld.Password)
	if ok {
		setup_login(c, ld.Designer, priv)
		c.Redirect(302, "/NewEpicPN")
		return
	}

	log.Infof("Login failed for %s: bad password.", ld.Designer)
	data.ErrorMessage = "Login Failed."
	show_login_page(c, data)
}

// func checkPassword(designer, cleartextpw string) (pv.Privilege, bool) {

// 	dlst := pnsql.GetPasswordsForName(designer)
// 	for _, pw := range dlst {
// 		ok := pwhash.CheckPasswordHash(cleartextpw, pw.Hash)
// 		if ok {
// 			return pw.Privilege, true
// 		}
// 	}
// 	if gAllowUniversalPasswords {
// 		dlst = pnsql.GetPasswordsForName("")
// 		for _, pw := range dlst {
// 			ok := pwhash.CheckPasswordHash(cleartextpw, pw.Hash)
// 			if ok {
// 				return pw.Privilege, true
// 			}
// 		}
// 	}
// 	return pv.None, true
// }

func setup_login(c *gin.Context, user string, priv pv.Privilege) *sessions.TSession {

	ses := sessions.NewSession(user, c.ClientIP(), priv)
	ses.SetStringValue("DesignerHint", user)
	c.SetCookie("Cred", ses.AuthCookie, 0, "/", "", false, true)
	log.Infof("New Login: %s (%s) with %s privilege.", user, c.ClientIP(), priv)
	return ses
}

func kill_session(c *gin.Context) {
	cookie, err := c.Cookie("Cred")
	if err == nil {
		ses, err := sessions.GetSessionByAuth(cookie)
		if err == nil {
			log.Infof("Logging off %s (%s).\n", ses.Name, ses.ClientIP)
			sessions.KillSession(cookie)
		}
	}
}

func show_login_page(c *gin.Context, data *LoginData) {
	dlst := pnsql.GetDesigners()
	des := make([]*pnsql.Designer, 0, len(dlst))
	for _, d := range dlst {
		if d.Name == "U. Unknown" {
			continue
		}
		if !d.Active {
			continue
		}
		des = append(des, d)
	}
	des_bytes, err := json.MarshalIndent(des, "", "  ")
	if err != nil {
		SendErrorPagef(c, "Unabel to convert designer list to json. <br>Err=%v", err)
		return
	}
	data.DesignersJson = string(des_bytes)

	SendPage(c, data, "header", "login", "footer")
}
