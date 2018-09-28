// --------------------------------------------------------------------
// sessions.go -- Manages logged in users and sessions...
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package sessions

import (
	"epic/lib/uuid"
	"fmt"
	"sync"
	"time"
)

type TSession struct {
	Name       string
	ClientIP   string
	LastAccess time.Time
	LoginTime  time.Time
	Data       map[string]interface{}
	AuthCookie uuid.UUID
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

func NewSession(name, ClientIP string) *TSession {
	session := new(TSession)
	session.Name = name
	session.ClientIP = ClientIP
	session.LastAccess = time.Now()
	session.Data = make(map[string]interface{})
	session.AuthCookie = uuid.New()
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	gSessions = append(gSessions, session)
	return session
}

func GetSessionByAuth(AuthCookie uuid.UUID) (*TSession, error) {
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
