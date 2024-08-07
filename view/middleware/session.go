package middleware

import (
	"context"
	"net/http"

	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/service"
	"kellnhofer.com/work-log/util/security"
)

// SessionHolder acts as a wrapper to be able to change the session of the current context.
type SessionHolder struct {
	session *model.Session
}

// Set sets the current session.
func (h *SessionHolder) Set(sess *model.Session) {
	h.session = sess
}

// Get returns the current session.
func (h *SessionHolder) Get() *model.Session {
	return h.session
}

// Clear clears the current session.
func (h *SessionHolder) Clear() {
	h.session = nil
}

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

	// Create system context
	sysCtx := security.CreateSystemContext(r.Context())

	// Get session cookie
	sessCookie, _ := r.Cookie(constant.SessionCookieName)

	// Get session ID
	var sessId string
	if sessCookie != nil {
		sessId = sessCookie.Value
	}

	// Try to load session
	var iniSess *model.Session
	if sessId != "" {
		iniSess = m.getSession(sysCtx, sessId)
	}
	if iniSess != nil {
		log.Debugf("Session '%s' is still valid.", sessId)
	}

	// If no session found: Create new session
	if iniSess == nil {
		iniSess = m.createSession()
		log.Debugf("New session '%s' was created.", iniSess.Id)
		sessCookie = &http.Cookie{Name: constant.SessionCookieName, Value: iniSess.Id, Path: "/",
			HttpOnly: true}
		http.SetCookie(w, sessCookie)
	}

	// Create session holder
	sessHolder := &SessionHolder{iniSess}

	// Update context
	ctx := context.WithValue(r.Context(), constant.ContextKeySessionHolder, sessHolder)

	// Forward to next handler
	next(w, r.WithContext(ctx))

	// Get session from session holder
	altSess := sessHolder.session

	// If session was replaced/deleted: Delete old session
	wasSessionReplaced := altSess != nil && altSess != iniSess
	if wasSessionReplaced {
		m.deleteSession(sysCtx, iniSess.Id)
		log.Debugf("Session '%s' was replaced by session '%s'.", iniSess.Id, altSess.Id)
	}
	wasSessionClosed := altSess == nil
	if wasSessionClosed {
		m.deleteSession(sysCtx, iniSess.Id)
		log.Debugf("Session '%s' was closed.", iniSess.Id)
	}

	// Save current session
	if !wasSessionClosed {
		m.saveSession(sysCtx, altSess)
	}

	log.Verb("After create/restore session.")
}

func (m *SessionMiddleware) createSession() *model.Session {
	return model.NewSession()
}

func (m *SessionMiddleware) getSession(ctx context.Context, sessId string) *model.Session {
	sess, err := m.sServ.GetSession(ctx, sessId)
	if err != nil {
		panic(err)
	}

	if sess == nil || sess.IsExpired() {
		return nil
	}

	return sess
}

func (m *SessionMiddleware) saveSession(ctx context.Context, sess *model.Session) {
	sess.Renew()
	if err := m.sServ.SaveSession(ctx, sess); err != nil {
		panic(err)
	}
}

func (m *SessionMiddleware) deleteSession(ctx context.Context, sessId string) {
	if err := m.sServ.DeleteSession(ctx, sessId); err != nil {
		panic(err)
	}
}
