package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"kellnhofer.com/work-log/pkg/constant"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/pkg/util/security"
)

// SecurityMiddleware creates the security context.
type SecurityMiddleware struct {
	uServ *service.UserService
}

// NewSecurityMiddleware create a new SecurityMiddleware.
func NewSecurityMiddleware(us *service.UserService) *SecurityMiddleware {
	return &SecurityMiddleware{us}
}

// CreateHandler creates a new handler to process requests.
func (m *SecurityMiddleware) CreateHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Verb("Before API auth check.")

		err := m.process(next, c)

		log.Verb("After API auth check.")

		return err
	}
}

func (m *SecurityMiddleware) process(next echo.HandlerFunc, c echo.Context) error {
	// Get request
	req := c.Request()

	// Create system context
	sysCtx := security.CreateSystemContext(req.Context())

	userId := model.AnonymousUserId

	// Was authentication data provided?
	if m.hasAuthenticationData(req) {
		// Get user credentials
		username, password, err := m.getUserCredentials(req)
		if err != nil {
			return err
		}

		log.Debugf("Authenticating user '%s' ...", username)

		// Try to authenticate user
		user, err := m.authenticateUser(sysCtx, username, password)
		if err != nil {
			return err
		}

		// Check user is activated
		if err := m.checkUserActivated(req, user); err != nil {
			return err
		}

		userId = user.Id
	}

	// Create security context
	secCtx, err := m.createSecurityContext(sysCtx, userId)
	if err != nil {
		return err
	}

	// Update context
	ctx := context.WithValue(req.Context(), constant.ContextKeySecurityContext, secCtx)
	c.SetRequest(req.WithContext(ctx))

	// Forward to next handler
	return next(c)
}

func (m *SecurityMiddleware) hasAuthenticationData(r *http.Request) bool {
	return r.Header.Get("Authorization") != ""
}

func (m *SecurityMiddleware) getUserCredentials(r *http.Request) (string, string, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		err := e.NewError(e.AuthDataInvalid, "Invalid authentication data.")
		log.Debug(err.StackTrace())
		return "", "", err
	}
	return username, password, nil
}

func (m *SecurityMiddleware) authenticateUser(ctx context.Context, username string,
	password string) (*model.User, error) {
	user, guErr := m.uServ.GetUserByUsername(ctx, username)
	if guErr != nil {
		return nil, guErr
	}
	if user == nil {
		err := e.NewError(e.AuthCredentialsInvalid, "Invalid credentials.")
		log.Debug(err.StackTrace())
		return nil, err
	}

	if cpwErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); cpwErr != nil {
		err := e.WrapError(e.AuthCredentialsInvalid, "Invalid credentials.", cpwErr)
		log.Debug(err.StackTrace())
		return nil, err
	}

	return user, nil
}

func (m *SecurityMiddleware) checkUserActivated(r *http.Request, user *model.User) error {
	path := r.URL.EscapedPath()
	path = strings.TrimPrefix(path, constant.ApiPath)

	if path != "/user" && !strings.HasPrefix(path, "/user/") && user.MustChangePassword {
		err := e.NewError(e.AuthUserNotActivated, "User must change password.")
		log.Debug(err.StackTrace())
		return err
	}

	return nil
}

func (m *SecurityMiddleware) createSecurityContext(ctx context.Context, userId int) (
	*model.SecurityContext, error) {
	if userId == model.AnonymousUserId {
		return model.GetAnonymousUserSecurityContext(), nil
	}

	userRoles, err := m.uServ.GetUserRoles(ctx, userId)
	if err != nil {
		return nil, err
	}

	return model.NewSecurityContext(userId, userRoles), nil
}

// AuthCheckMiddleware ensures that a user was authenticated.
type AuthCheckMiddleware struct {
}

// NewAuthCheckMiddleware create a new AuthCheckMiddleware.
func NewAuthCheckMiddleware() *AuthCheckMiddleware {
	return &AuthCheckMiddleware{}
}

// CreateHandler creates a new handler to process requests.
func (m *AuthCheckMiddleware) CreateHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Verb("Before API auth check.")

		err := m.process(next, c)

		log.Verb("After API auth check.")

		return err
	}
}

func (m *AuthCheckMiddleware) process(next echo.HandlerFunc, c echo.Context) error {
	// Get request
	req := c.Request()

	// Get security context
	secCtx := req.Context().Value(constant.ContextKeySecurityContext).(*model.SecurityContext)

	// Wad anonymous user?
	if secCtx.UserId == model.AnonymousUserId {
		log.Debug("User must authenticate.")

		// Request authentication
		log.Debug("Requesting authentication ...")
		res := c.Response().Writer
		res.Header().Set("WWW-Authenticate", "Basic")
		return c.NoContent(http.StatusUnauthorized)
	}

	log.Debugf("User %d is authenticated.", secCtx.UserId)

	// Forward to next handler
	return next(c)
}
