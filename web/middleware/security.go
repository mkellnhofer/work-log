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
	"kellnhofer.com/work-log/web"
)

// SecurityMiddleware creates the security context.
type SecurityMiddleware struct {
	uServ *service.UserService
}

// NewSecurityMiddleware create a new SecurityMiddleware.
func NewSecurityMiddleware(uServ *service.UserService) *SecurityMiddleware {
	return &SecurityMiddleware{uServ}
}

// CreateHandler creates a new handler to process requests.
func (m *SecurityMiddleware) CreateHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Verb("Before security init.")

		err := m.process(next, c)

		log.Verb("After security init.")

		return err
	}
}

func (m *SecurityMiddleware) process(next echo.HandlerFunc, c echo.Context) error {
	// Get request
	req := c.Request()

	// Create system context
	sysCtx := security.CreateSystemContext(req.Context())

	// Get session from context
	sessHolder := req.Context().Value(constant.ContextKeySessionHolder).(*SessionHolder)
	sess := sessHolder.Get()

	// Get current user ID
	userId := sess.UserId

	// Create security context
	var secCtx *model.SecurityContext
	if userId == model.AnonymousUserId {
		secCtx = model.GetAnonymousUserSecurityContext()
	} else {
		userRoles, err := m.getUserRoles(sysCtx, userId)
		if err != nil {
			return err
		}
		secCtx = model.NewSecurityContext(userId, userRoles)
	}

	// Update context
	ctx := context.WithValue(req.Context(), constant.ContextKeySecurityContext, secCtx)
	c.SetRequest(req.WithContext(ctx))

	// Forward to next handler
	return next(c)
}

func (m *SecurityMiddleware) getUserRoles(sysCtx context.Context, userId int) ([]model.Role, error) {
	userRoles, err := m.uServ.GetUserRoles(sysCtx, userId)
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

// AuthCheckMiddleware ensures that a user was authenticated.
type AuthCheckMiddleware struct {
	uServ *service.UserService
}

// NewAuthCheckMiddleware create a new AuthCheckMiddleware.
func NewAuthCheckMiddleware(uServ *service.UserService) *AuthCheckMiddleware {
	return &AuthCheckMiddleware{uServ}
}

// CreateHandler creates a new handler to process requests.
func (m *AuthCheckMiddleware) CreateHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Verb("Before auth check.")

		err := m.process(next, c)

		log.Verb("After auth check.")

		return err
	}
}

func (m *AuthCheckMiddleware) process(next echo.HandlerFunc, c echo.Context) error {
	// Get request
	req := c.Request()

	// Create system context
	sysCtx := security.CreateSystemContext(req.Context())

	// Get user ID
	secCtx := req.Context().Value(constant.ContextKeySecurityContext).(*model.SecurityContext)
	userId := secCtx.UserId

	// Get session holder
	sessHolder := req.Context().Value(constant.ContextKeySessionHolder).(*SessionHolder)
	sess := sessHolder.Get()

	// If user is not authenticated: Redirect to login page
	if userId == model.AnonymousUserId {
		log.Debugf("User must authenticate. (Session: '%s')", sess.GetShortId())
		return m.redirectLogin(c, sess)
	}

	// Get requested path
	reqPath := getRequestPath(req)

	// Get user
	user, err := m.uServ.GetUserById(sysCtx, userId)
	if err != nil {
		return err
	}

	// If user must change password: Redirect to password change page
	if reqPath != "/password-change" && user.MustChangePassword {
		log.Debugf("User %d must change password. (Session: '%s')", userId, sess.GetShortId())
		return m.redirectPasswordChange(c, sess)
	}

	log.Debugf("User %d is authenticated. (Session: '%s')", userId, sess.GetShortId())

	// Forward to next handler
	return next(c)
}

func (m *AuthCheckMiddleware) redirectLogin(c echo.Context, sess *model.Session) error {
	log.Debug("Redirecting to login page ...")
	return m.redirect(c, sess, "/login")
}

func (m *AuthCheckMiddleware) redirectPasswordChange(c echo.Context, sess *model.Session) error {
	log.Debug("Redirecting to password change page ...")
	return m.redirect(c, sess, "/password-change")
}

func (m *AuthCheckMiddleware) redirect(c echo.Context, sess *model.Session, url string) error {
	// Get requested URL
	reqUrl := getRequestUrl(c.Request())

	// Save requested URL in session
	sess.PreviousUrl = reqUrl

	// Redirect
	if web.IsHtmxRequest(c) {
		web.HtmxRedirectUrl(c, url)
		return c.NoContent(http.StatusOK)
	} else {
		return c.Redirect(http.StatusFound, url)
	}
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
