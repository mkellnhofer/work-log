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

	sysCtx := security.CreateSystemContext(r.Context())
	m.handle(sysCtx, w, r, next)

	log.Verb("After security init.")
}

func (m *SecurityMiddleware) handle(sysCtx context.Context, w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	// Get session from context
	sessHolder := r.Context().Value(constant.ContextKeySessionHolder).(*SessionHolder)
	sess := sessHolder.Get()

	// Get current user ID
	userId := sess.UserId

	// Create security context
	var secCtx *model.SecurityContext
	if userId == model.AnonymousUserId {
		secCtx = model.GetAnonymousUserSecurityContext()
	} else {
		userRoles := m.getUserRoles(sysCtx, userId)
		secCtx = model.NewSecurityContext(userId, userRoles)
	}

	// Update context
	ctx := context.WithValue(r.Context(), constant.ContextKeySecurityContext, secCtx)

	// Forward to next handler
	next(w, r.WithContext(ctx))
}

func (m *SecurityMiddleware) getUserRoles(sysCtx context.Context, userId int) []model.Role {
	userRoles, err := m.uServ.GetUserRoles(sysCtx, userId)
	if err != nil {
		panic(err)
	}
	return userRoles
}

// AuthCheckMiddleware ensures that a user was authenticated.
type AuthCheckMiddleware struct {
	uServ *service.UserService
}

// NewAuthCheckMiddleware create a new AuthCheckMiddleware.
func NewAuthCheckMiddleware(uServ *service.UserService) *AuthCheckMiddleware {
	return &AuthCheckMiddleware{uServ}
}

// ServeHTTP processes requests.
func (m *AuthCheckMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Verb("Before auth check.")

	sysCtx := security.CreateSystemContext(r.Context())
	m.handle(sysCtx, w, r, next)

	log.Verb("After auth check.")
}

func (m *AuthCheckMiddleware) handle(sysCtx context.Context, w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	// Get user ID
	secCtx := r.Context().Value(constant.ContextKeySecurityContext).(*model.SecurityContext)
	userId := secCtx.UserId

	// Get session holder
	sessHolder := r.Context().Value(constant.ContextKeySessionHolder).(*SessionHolder)
	sess := sessHolder.Get()

	// If user is not authenticated: Redirect to login page
	if userId == model.AnonymousUserId {
		log.Debugf("User must authenticate. (Session: '%s')", sess.Id)
		m.redirectLogin(sess, w, r)
		return
	}

	// Get requested path
	reqPath := getRequestPath(r)

	// Get user
	user := m.getUser(sysCtx, userId)

	// If user must change password: Redirect to password change page
	if reqPath != "/password-change" && user.MustChangePassword {
		log.Debugf("User %d must change password. (Session: '%s')", userId, sess.Id)
		m.redirectPasswordChange(sess, w, r)
		return
	}

	// Forward to next handler
	log.Debugf("User %d is authenticated. (Session: '%s')", userId, sess.Id)
	next(w, r)
}

func (m *AuthCheckMiddleware) getUser(sysCtx context.Context, userId int) *model.User {
	user, err := m.uServ.GetUserById(sysCtx, userId)
	if err != nil {
		panic(err)
	}
	return user
}

func (m *AuthCheckMiddleware) redirectLogin(sess *model.Session, w http.ResponseWriter,
	r *http.Request) {
	log.Debug("Redirecting to login page ...")
	m.redirect(sess, w, r, "/login")
}

func (m *AuthCheckMiddleware) redirectPasswordChange(sess *model.Session, w http.ResponseWriter,
	r *http.Request) {
	log.Debug("Redirecting to password change page ...")
	m.redirect(sess, w, r, "/password-change")
}

func (m *AuthCheckMiddleware) redirect(sess *model.Session, w http.ResponseWriter, r *http.Request,
	url string) {
	// Get requested URL
	reqUrl := getRequestUrl(r)

	// Save requested URL in session
	sess.PreviousUrl = reqUrl

	// Redirect to login page
	http.Redirect(w, r, url, http.StatusFound)
}

func getRequestUrl(r *http.Request) string {
	path := r.URL.EscapedPath()
	query := r.URL.RawQuery
	req := path
	if query != "" {
		req = req + "?" + query
	}
	return req
}

func getRequestPath(r *http.Request) string {
	return r.URL.EscapedPath()
}
