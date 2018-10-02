// --------------------------------------------------------------------
// authorizer.go -- Authorizer: checks login and fills header data.
//
// Created 2018-09-29 DLB
// --------------------------------------------------------------------

package pages

import (
	"epic/lib/log"
	"epic/pnserver/sessions"
	//"epic/pnserver/pnsql"
	"fmt"
	"github.com/gin-gonic/gin"
	//"strings"
	//"time"
)

func authorizer(c *gin.Context) {
	// If not logged in, redirect to login page.
	cookie, err := c.Cookie("Cred")
	ses, err := sessions.GetSessionByAuth(cookie)
	if err != nil {
		// We are not logged in!
		log.Infof("Invalid auth cookie.. redirecting to Login.")
		c.Redirect(302, "/Login")
		c.Abort()
		return
	}

	// Session found. Fill up data.
	c.Set("Session", ses)
	hdrdata := &HeaderData{}
	hdrdata.PageTabTitle = "Epic PN"
	hdrdata.IsLoggedIn = true
	hdrdata.UserFormattedName = ses.Name
	hdrdata.IsAdmin = ses.IsAdmin()
	c.Set("HeaderData", hdrdata)
}

func guest_auth(c *gin.Context) {
	cookie, err := c.Cookie("Cred")
	ses, err := sessions.GetSessionByAuth(cookie)
	if err != nil {
		c.Set("Session", sessions.NewGuestSession(c.ClientIP()))
		// We are not logged in!
		// Provide a guest level version of the header.
		hdrdata := &HeaderData{}
		hdrdata.PageTabTitle = "Epic PN"
		hdrdata.IsLoggedIn = false
		hdrdata.UserFormattedName = ""
		hdrdata.IsAdmin = false
		c.Set("HeaderData", hdrdata)
		return
	}

	// Session found. Fill up data.
	c.Set("Session", ses)
	hdrdata := &HeaderData{}
	hdrdata.PageTabTitle = "Epic PN"
	hdrdata.IsLoggedIn = true
	hdrdata.UserFormattedName = ses.Name
	hdrdata.IsAdmin = ses.IsAdmin()
	c.Set("HeaderData", hdrdata)
}

// Note: GetHeaderDate should NEVER return nil for any
// page that has used authorizer or guest_auth.  Therefore
// don't bother checking for nil.  If nil is returned, allow
// your code to painic.
func GetHeaderData(c *gin.Context) *HeaderData {
	hdrdata, ok := c.Get("HeaderData")
	if !ok {
		err := fmt.Errorf("Header Data not avaliable.")
		log.Errorf("%v\n", err)
		c.AbortWithError(400, err)
		return nil
	}
	hdr, ok := hdrdata.(*HeaderData)
	if !ok {
		err := fmt.Errorf("Header Data assert error! Programming BUG.")
		log.Errorf("%v\n", err)
		c.AbortWithError(400, err)
		return nil
	}
	return hdr
}

// GetSession returns session data for the current session.
// This call should always return a valid session for any
// page that used authorizer or guest_auth.
func GetSession(c *gin.Context) *sessions.TSession {
	x, ok := c.Get("Session")
	if !ok {
		err := fmt.Errorf("Session data not avaliable after authincation!?")
		log.Errorf("%v\n", err)
		c.AbortWithError(400, err)
		return nil
	}
	ses, ok := x.(*sessions.TSession)
	if !ok {
		err := fmt.Errorf("Session Data assert error! Programming BUG.")
		log.Errorf("%v\n", err)
		c.AbortWithError(400, err)
		return nil
	}
	return ses
}
