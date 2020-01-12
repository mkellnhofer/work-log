package middleware

import (
	"net/http"

	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/service"
)

// AuthMiddleware creates/restores the session.
type AuthMiddleware struct {
	sServ *service.SessionService
}

// NewAuthMiddleware create a new session middleware.
func NewAuthMiddleware(sServ *service.SessionService) *AuthMiddleware {
	return &AuthMiddleware{sServ}
}

// ServeHTTP processes requests.
func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Verb("Before auth check.")

	// Get session from context
	sess := r.Context().Value(constant.ContextKeySession).(*model.Session)

	// Authenticated?
	if sess.UserId != 0 {
		log.Debugf("User %d is authenticated. (Session: '%s')", sess.UserId, sess.Id)

		// Forward to next handler
		next(w, r)
	} else {
		log.Debugf("User must authenticate. (Session: '%s')", sess.Id)

		// Get current request
		path := r.URL.EscapedPath()
		query := r.URL.RawQuery
		req := path
		if query != "" {
			req = req + "?" + query
		}

		// Save current request in session
		sess.PreviousUrl = req
		ssErr := m.sServ.SaveSession(sess)
		if ssErr != nil {
			panic(ssErr)
		}

		// Redirect to login page
		log.Debug("Redirecting to login page ...")
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	log.Verb("After auth check.")
}
