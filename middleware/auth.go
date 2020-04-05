package middleware

import (
	"net/http"

	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/log"
)

// AuthMiddleware creates/restores the session.
type AuthMiddleware struct {
}

// NewAuthMiddleware create a new session middleware.
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

// ServeHTTP processes requests.
func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Verb("Before auth check.")

	// Get session from context
	sessHolder := r.Context().Value(constant.ContextKeySessionHolder).(*SessionHolder)
	sess := sessHolder.Get()

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

		// Redirect to login page
		log.Debug("Redirecting to login page ...")
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	log.Verb("After auth check.")
}
