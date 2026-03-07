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
	"kellnhofer.com/work-log/pkg/util"
	"kellnhofer.com/work-log/pkg/util/security"
)

const (
	basicAuthPrefix = "Basic "
	bearerAuthPrefix = "Bearer "
)

type authType int

const (
	authTypeBasic  authType = iota
	authTypeBearer
)

type authResult struct {
	authType authType
	user     *model.User
}

// SecurityMiddleware creates the security context.
type SecurityMiddleware struct {
	uServ *service.UserService
	tServ *service.TokenService
}

// NewSecurityMiddleware create a new SecurityMiddleware.
func NewSecurityMiddleware(us *service.UserService, ts *service.TokenService) *SecurityMiddleware {
	return &SecurityMiddleware{us, ts}
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
		// Authenticate user
		ar, err := m.authenticate(sysCtx, req)
		if err != nil {
			return err
		}

		// Do post authentication checks
		if err := m.checkAuthentication(req, ar); err != nil {
			return err
		}

		userId = ar.user.Id
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
	return m.getAuthenticationData(r) != ""
}

func (m *SecurityMiddleware) authenticate(ctx context.Context, r *http.Request) (*authResult, error) {
	authData := m.getAuthenticationData(r)

	if m.isBasicAuthRequest(authData) {
		user, err := m.authenticateBasicAuth(ctx, r)
		return &authResult{authTypeBasic, user}, err
	} else if m.isBearerAuthRequest(authData) {
		user, err := m.authenticateBearerAuth(ctx, r)
		return &authResult{authTypeBearer, user}, err
	} else {
		err := e.NewError(e.AuthDataInvalid, "Invalid authentication data.")
		log.Debug(err.StackTrace())
		return nil, err
	}
}

func (m *SecurityMiddleware) isBasicAuthRequest(authData string) bool {
	return m.hasAuthPrefix(authData, basicAuthPrefix)
}

func (m *SecurityMiddleware) isBearerAuthRequest(authData string) bool {
	return m.hasAuthPrefix(authData, bearerAuthPrefix)
}

func (m *SecurityMiddleware) authenticateBasicAuth(ctx context.Context, r *http.Request) (*model.User,
	error) {
	// Get user credentials
	username, password, ok := m.getBasicAuthCredentials(r)
	if !ok {
		err := e.NewError(e.AuthDataInvalid, "Invalid authentication data.")
		log.Debug(err.StackTrace())
		return nil, err
	}

	log.Debugf("Authenticating user '%s' ...", username)

	// Try to authenticate user
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

func (m *SecurityMiddleware) getBasicAuthCredentials(r *http.Request) (string, string, bool) {
	return r.BasicAuth()
}

func (m *SecurityMiddleware) authenticateBearerAuth(ctx context.Context, r *http.Request) (*model.User,
	error) {
	// Get token
	tokenValue, ok := m.getBearerAuthToken(r)
	if !ok {
		err := e.NewError(e.AuthDataInvalid, "Invalid authentication data.")
		log.Debug(err.StackTrace())
		return nil, err
	}

	log.Debugf("Authenticating token '%s' ...", m.createTruncatedToken(tokenValue))

	// Try to authenticate token
	token, gtErr := m.tServ.GetTokenByValue(ctx, tokenValue)
	if gtErr != nil {
		return nil, gtErr
	}
	if token == nil {
		err := e.NewError(e.AuthTokenInvalid, "Invalid token.")
		log.Debug(err.StackTrace())
		return nil, err
	}
	user, guErr := m.uServ.GetUserById(ctx, token.UserId)
	if guErr != nil {
		return nil, guErr
	}
	if user == nil {
		err := e.NewError(e.AuthTokenInvalid, "Invalid token.")
		log.Debug(err.StackTrace())
		return nil, err
	}
	return user, nil
}

func (m *SecurityMiddleware) getBearerAuthToken(r *http.Request) (string, bool) {
	authData := m.getAuthenticationData(r)
	token := strings.TrimPrefix(authData, bearerAuthPrefix)
	if token == "" {
		return "", false
	}
	return token, true
}

func (m *SecurityMiddleware) checkAuthentication(r *http.Request, ar *authResult) error {
	if err := m.checkUserActivated(r, ar.user); err != nil {
		return err
	}

	switch ar.authType {
	case authTypeBasic:
		return nil
	case authTypeBearer:
		return m.checkTokenAllowedForEndpoint(r)
	default:
		err := e.NewError(e.AuthDataInvalid, "Invalid authentication data.")
		log.Debug(err.StackTrace())
		return err
	}
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

func (m *SecurityMiddleware) checkTokenAllowedForEndpoint(r *http.Request) error {
	path := r.URL.EscapedPath()
	path = strings.TrimPrefix(path, constant.ApiPath)

	if strings.HasPrefix(path, "/user/password") || strings.HasPrefix(path, "/user/tokens") {
		err := e.NewError(e.AuthTokenNotAllowed, "Bearer auth is not allowed for this endpoint.")
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

func (m *SecurityMiddleware) getAuthenticationData(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func (m *SecurityMiddleware) hasAuthPrefix(authData string, prefix string) bool {
	return len(authData) > len(prefix) && strings.HasPrefix(authData, prefix)
}

func (m *SecurityMiddleware) createTruncatedToken(token string) string {
	return util.CreateTruncatedString(token, 4)
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
