package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/constant"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/pkg/util/security"
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

// CreateHandler creates a new handler to process requests.
func (m *SessionMiddleware) CreateHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Verb("Before create/restore session.")

		err := m.process(next, c)

		log.Verb("After create/restore session.")

		return err
	}
}

func (m *SessionMiddleware) process(next echo.HandlerFunc, c echo.Context) error {
	// Get request
	req := c.Request()

	// Create system context
	sysCtx := security.CreateSystemContext(req.Context())

	// Get session cookie
	sessCookie, _ := req.Cookie(constant.SessionCookieName)

	// Get session ID
	var rawSessId string
	if sessCookie != nil {
		rawSessId = sessCookie.Value
		if !model.IsValidSessionId(rawSessId) {
			log.Debug("Invalid session ID was provided in cookie.")
			rawSessId = ""
		}
	}

	// Try to load session
	iniSess, err := m.getSession(sysCtx, rawSessId)
	if err != nil {
		return err
	}
	if iniSess != nil {
		log.Debugf("Session '%s' is still valid.", iniSess.GetShortRawId())
	}

	// If no session found: Create new session
	if iniSess == nil {
		iniSess = m.createSession()
		log.Debugf("New session '%s' was created.", iniSess.GetShortRawId())
		sessCookie = &http.Cookie{
			Name:     constant.SessionCookieName,
			Value:    iniSess.RawId,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		c.SetCookie(sessCookie)
	}

	// Create session holder
	sessHolder := &SessionHolder{iniSess}

	// Update context
	ctx := context.WithValue(req.Context(), constant.ContextKeySessionHolder, sessHolder)
	c.SetRequest(req.WithContext(ctx))

	// Forward to next handler
	if err := next(c); err != nil {
		return err
	}

	// Get session from session holder
	altSess := sessHolder.session

	// If session was replaced/deleted: Delete old session
	wasSessionReplaced := altSess != nil && altSess != iniSess
	if wasSessionReplaced {
		if err := m.deleteSession(sysCtx, iniSess.RawId); err != nil {
			return err
		}
		log.Debugf("Session '%s' was replaced by session '%s'.", iniSess.GetShortRawId(),
			altSess.GetShortRawId())
	}
	wasSessionClosed := altSess == nil
	if wasSessionClosed {
		if err := m.deleteSession(sysCtx, iniSess.RawId); err != nil {
			return err
		}
		log.Debugf("Session '%s' was closed.", iniSess.GetShortRawId())
	}

	// Save current session
	if !wasSessionClosed {
		if err := m.saveSession(sysCtx, altSess); err != nil {
			return err
		}
	}

	return nil
}

func (m *SessionMiddleware) createSession() *model.Session {
	return model.NewSession()
}

func (m *SessionMiddleware) getSession(ctx context.Context, rawSessId string) (*model.Session,
	error) {
	if rawSessId == "" {
		return nil, nil
	}

	sess, err := m.sServ.GetSession(ctx, rawSessId)
	if err != nil {
		return nil, err
	}

	if sess == nil || sess.IsExpired() {
		return nil, nil
	}

	return sess, nil
}

func (m *SessionMiddleware) saveSession(ctx context.Context, sess *model.Session) error {
	sess.Renew()
	return m.sServ.SaveSession(ctx, sess)
}

func (m *SessionMiddleware) deleteSession(ctx context.Context, rawSessId string) error {
	return m.sServ.DeleteSession(ctx, rawSessId)
}
