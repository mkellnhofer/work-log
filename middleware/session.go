package middleware

import (
	"context"
	"net/http"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/service"
)

// SessionMiddleware creates/restores the session.
type SessionMiddleware struct {
	sServ *service.SessionService
}

// NewSessionMiddleware create a new session middleware.
func NewSessionMiddleware(sServ *service.SessionService) *SessionMiddleware {
	return &SessionMiddleware{sServ}
}

// ServeHTTP processes requests.
func (m *SessionMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Verb("Before create/restore session.")

	sessCookie, _ := r.Cookie(constant.SessionCookieName)

	var sessId string
	if sessCookie != nil {
		sessId = sessCookie.Value
	}

	// Try to load session
	var sess *model.Session
	if sessId != "" {
		var gsErr *e.Error
		sess, gsErr = m.getSession(sessId)
		if gsErr != nil {
			panic(gsErr)
		}
	}

	// Was session found?
	var newSess bool
	if sess != nil {
		// Refresh session
		rsErr := m.sServ.RefreshSession(sess)
		if rsErr != nil {
			panic(rsErr)
		}
		log.Debugf("Session '%s' is still valid.", sessId)
		newSess = false
	} else {
		// Create a new session
		sessId = m.sServ.NewSessionId()
		sess = model.NewSession(sessId)
		csErr := m.sServ.SaveSession(sess)
		if csErr != nil {
			panic(csErr)
		}
		log.Debugf("New session '%s' was created.", sessId)
		newSess = true
	}

	// Get context
	ctx := r.Context()
	ctx = context.WithValue(ctx, constant.ContextKeySession, sess)

	// If a new session was created: Set session cookie
	if newSess {
		sessCookie = &http.Cookie{Name: constant.SessionCookieName, Value: sessId, Path: "/",
			HttpOnly: true}
		http.SetCookie(w, sessCookie)
	}

	// Forward to next handler
	next(w, r.WithContext(ctx))

	log.Verb("After create/restore session.")
}

func (m *SessionMiddleware) getSession(sessId string) (*model.Session, *e.Error) {
	sess, err := m.sServ.GetSession(sessId)
	if err != nil {
		return nil, err
	}

	if sess == nil || sess.IsExpired() {
		return nil, nil
	}

	return sess, nil
}
