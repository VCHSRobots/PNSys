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
	// Here, we get data from the session and store it away if possible.
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
	hdrdata := &HeaderData{}
	hdrdata.PageTabTitle = "Epic PN"
	hdrdata.IsLoggedIn = true
	hdrdata.UserFormattedName = ses.Name
	hdrdata.IsAdmin = ses.IsAdmin()
	c.Set("HeaderData", hdrdata)

}

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
