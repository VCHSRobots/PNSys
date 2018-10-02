// --------------------------------------------------------------------
// sessions.go -- Manages logged in users and sessions...
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package sessions

import (
	"epic/lib/log"
	"epic/lib/uuid"
	"fmt"
	"sync"
	"time"
)

type SessionPrivilege string

const (
	Privilege_Guest = "Guest"
	Privilege_User  = "User"
	Privilege_Admin = "Admin"
)

type TSession struct {
	Name       string
	ClientIP   string
	Privilege  SessionPrivilege
	LastAccess time.Time
	LoginTime  time.Time
	Data       map[string]interface{}
	AuthCookie string
}

var gSessions []*TSession
var gSessionLock sync.Mutex
var gTimeToLive int = 7200 // In Seconds

func init() {
	gSessions = make([]*TSession, 0, 20)
	go func() {
		for {
			time.Sleep(30)
			clean_sessions()
		}
	}()
}

func clean_sessions() {
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	t := time.Now()
	lst := make([]*TSession, 0, len(gSessions))
	for _, s := range gSessions {
		elapsed := t.Sub(s.LastAccess)
		if elapsed < time.Duration(gTimeToLive)*time.Second {
			lst = append(lst, s)
		} else {
			log.Infof("Logging out %s (%s) due to inactivity.", s.Name, s.ClientIP)
		}
	}
	gSessions = lst
}

func SetTimeToLive(secs int) {
	gTimeToLive = secs
}

func GetTimeToLive() int {
	return gTimeToLive
}

func NewSession(name, ClientIP string, Privilege SessionPrivilege) *TSession {
	session := new(TSession)
	session.Name = name
	session.ClientIP = ClientIP
	session.Privilege = Privilege
	session.LoginTime = time.Now()
	session.LastAccess = time.Now()
	session.Data = make(map[string]interface{}, 30)
	session.AuthCookie = uuid.New().String()
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	gSessions = append(gSessions, session)
	return session
}

func NewGuestSession(ipaddr string) *TSession {
	session := new(TSession)
	session.Name = "Guest"
	session.ClientIP = ipaddr
	session.Privilege = "Privilege_Guest"
	session.LoginTime = time.Now()
	session.LastAccess = time.Now()
	session.Data = make(map[string]interface{}, 30)
	session.AuthCookie = ""
	return session
}

func KillSession(AuthCookie string) {
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	i := -1
	for j := 0; j < len(gSessions); j++ {
		if gSessions[j].AuthCookie == AuthCookie {
			i = j
			break
		}
	}
	if i < 0 {
		return
	}
	copy(gSessions[i:], gSessions[i+1:])
	gSessions[len(gSessions)-1] = nil
	gSessions = gSessions[:len(gSessions)-1]
}

func GetSessionByAuth(AuthCookie string) (*TSession, error) {
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	for _, ss := range gSessions {
		if ss.AuthCookie == AuthCookie {
			return ss, nil
		}
	}
	return nil, fmt.Errorf("Session not found.")
}

func GetAllSessions() []*TSession {
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	lst := make([]*TSession, 0, len(gSessions))
	for _, s := range gSessions {
		lst = append(lst, s)
	}
	return lst
}

func (ses *TSession) GetStringValue(key string) string {
	t, ok := ses.Data[key]
	if !ok {
		return ""
	}
	s, _ := t.(string)
	return s
}

func (ses *TSession) SetStringValue(key, value string) {
	ses.Data[key] = value
}

func (ses *TSession) IsAdmin() bool {
	return ses.Privilege == Privilege_Admin
}
