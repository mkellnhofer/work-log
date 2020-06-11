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

// SecurityMiddleware creates the security context.
type SecurityMiddleware struct {
	uServ *service.UserService
}

// NewSecurityMiddleware create a new SecurityMiddleware.
func NewSecurityMiddleware(uServ *service.UserService) *SecurityMiddleware {
	return &SecurityMiddleware{uServ}
}

// ServeHTTP processes requests.
func (m *SecurityMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Verb("Before security init.")

	// Create system context
	sysCtx := security.CreateSystemContext(r.Context())

	// Get session from context
	sessHolder := r.Context().Value(constant.ContextKeySessionHolder).(*SessionHolder)
	sess := sessHolder.Get()

	// Get current user ID
	userId := sess.UserId

	// Create security context
	var secCtx *model.SecurityContext
	if userId != model.AnonymousUserId {
		userRoles, err := m.uServ.GetUserRoles(sysCtx, userId)
		if err != nil {
			panic(err)
		}
		secCtx = model.NewSecurityContext(userId, userRoles)
	} else {
		secCtx = model.GetAnonymousUserSecurityContext()
	}

	// Update context
	ctx := context.WithValue(r.Context(), constant.ContextKeySecurityContext, secCtx)

	// Forward to next handler
	next(w, r.WithContext(ctx))

	log.Verb("After security init.")
}

// AuthCheckMiddleware ensures that a user was authenticated.
type AuthCheckMiddleware struct {
}

// NewAuthCheckMiddleware create a new AuthCheckMiddleware.
func NewAuthCheckMiddleware() *AuthCheckMiddleware {
	return &AuthCheckMiddleware{}
}

// ServeHTTP processes requests.
func (m *AuthCheckMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Verb("Before auth check.")

	// Get session holder
	sessHolder := r.Context().Value(constant.ContextKeySessionHolder).(*SessionHolder)
	sess := sessHolder.Get()

	// Get security context
	secCtx := r.Context().Value(constant.ContextKeySecurityContext).(*model.SecurityContext)

	// Authenticated?
	if secCtx.UserId != model.AnonymousUserId {
		log.Debugf("User %d is authenticated. (Session: '%s')", secCtx.UserId, sess.Id)

		// Forward to next handler
		next(w, r)
	} else {
		log.Debugf("User must authenticate. (Session: '%s')", sess.Id)

		// Get requested URL
		path := r.URL.EscapedPath()
		query := r.URL.RawQuery
		req := path
		if query != "" {
			req = req + "?" + query
		}

		// Save requested URL in session
		sess.PreviousUrl = req

		// Redirect to login page
		log.Debug("Redirecting to login page ...")
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	log.Verb("After auth check.")
}
