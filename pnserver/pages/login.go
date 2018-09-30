// --------------------------------------------------------------------
// login.go -- Handles Login Page
//
// Created 2018-09-22 DLB
// --------------------------------------------------------------------

package pages

import (
	"encoding/json"
	"epic/lib/log"
	"epic/lib/pwhash"
	"epic/pnserver/pnsql"
	"epic/pnserver/sessions"
	"github.com/gin-gonic/gin"
)

const (
	gPwAdmin = "JDJhJDA2JHNHWGt6SkM0RnVJL0EvZ25aNUlZai5wcC4uQWlxL1dNdlJDN2w5dkJJZHluR0xrZEd1djFT"
	gPwUser  = "JDJhJDA2JHB0S2xmamVtL3k4OHpwZ1JEQzRYLy5qUWc1NjNqQWZzTUxOd3REUVFUM082UDVpTlZiWFpt"
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

	IsAdmin := pwhash.CheckPasswordHash(ld.Password, gPwAdmin)
	if IsAdmin {
		ses := sessions.NewSession(ld.Designer, c.ClientIP(), sessions.Privilege_Admin)
		ses.SetStringValue("DesignerHint", ld.Designer)
		c.SetCookie("Cred", ses.AuthCookie, 0, "/", "", false, true)
		log.Infof("New Login: %s (%s) with privilege: Admin.", ld.Designer, c.ClientIP())
		c.Redirect(302, "/NewEpicPN")
		return
	}

	IsUser := pwhash.CheckPasswordHash(ld.Password, gPwUser)
	if IsUser {
		ses := sessions.NewSession(ld.Designer, c.ClientIP(), sessions.Privilege_User)
		ses.SetStringValue("DesignerHint", ld.Designer)
		c.SetCookie("Cred", ses.AuthCookie, 0, "/", "", false, true)
		log.Infof("New Login: %s (%s) with privilege: User.", ld.Designer, c.ClientIP())
		c.Redirect(302, "/NewEpicPN")
		return
	}

	log.Infof("Login failed: bad password.")
	data.ErrorMessage = "Login Failed."
	show_login_page(c, data)
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
		log.Errorf("Unable to convert to json. Err=%v", err)
		c.AbortWithError(400, err)
		return
	}
	data.DesignersJson = string(des_bytes)

	SendPage(c, data, "header", "login", "footer")
}
