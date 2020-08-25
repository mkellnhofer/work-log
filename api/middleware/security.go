package middleware

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
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
func NewSecurityMiddleware(us *service.UserService) *SecurityMiddleware {
	return &SecurityMiddleware{us}
}

// ServeHTTP processes requests.
func (m *SecurityMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Verb("Before API auth check.")

	// Create system context
	sysCtx := security.CreateSystemContext(r.Context())

	userId := model.AnonymousUserId

	// Was authentication data provided?
	if r.Header.Get("Authorization") != "" {
		// Get user credentials
		username, password := m.getUserCredentials(r)

		log.Debugf("Authenticating user '%s' ...", username)

		// Try to authenticate user
		user := m.authenticateUser(sysCtx, username, password)

		// Get requested path
		reqPath := getRequestPath(r)

		// Check is user activated
		if reqPath != "/user" && !strings.HasPrefix(reqPath, "/user/") {
			checkUserActivated(user)
		}

		userId = user.Id
	}

	// Create security context
	secCtx := m.createSecurityContext(sysCtx, userId)

	// Update context
	ctx := context.WithValue(r.Context(), constant.ContextKeySecurityContext, secCtx)

	// Forward to next handler
	next(w, r.WithContext(ctx))

	log.Verb("After API auth check.")
}

func (m *SecurityMiddleware) getUserCredentials(r *http.Request) (string, string) {
	username, password, ok := r.BasicAuth()
	if !ok {
		err := e.NewError(e.AuthUnknown, "Invalid authentication request.")
		log.Debug(err.StackTrace())
		panic(err)
	}
	return username, password
}

func (m *SecurityMiddleware) authenticateUser(ctx context.Context, username string,
	password string) *model.User {
	user, guErr := m.uServ.GetUserByUsername(ctx, username)
	if guErr != nil {
		panic(guErr)
	}
	if user == nil {
		err := e.NewError(e.AuthCredentialsInvalid, "Invalid credentials.")
		log.Debug(err.StackTrace())
		panic(err)
	}

	if cpwErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); cpwErr != nil {
		err := e.WrapError(e.AuthCredentialsInvalid, "Invalid credentials.", cpwErr)
		log.Debug(err.StackTrace())
		panic(err)
	}

	return user
}

func checkUserActivated(user *model.User) {
	if user.MustChangePassword {
		err := e.NewError(e.AuthUserNotActivated, "User must change password.")
		log.Debug(err.StackTrace())
		panic(err)
	}
}

func (m *SecurityMiddleware) createSecurityContext(ctx context.Context,
	userId int) *model.SecurityContext {
	if userId == model.AnonymousUserId {
		return model.GetAnonymousUserSecurityContext()
	}

	userRoles, err := m.uServ.GetUserRoles(ctx, userId)
	if err != nil {
		panic(err)
	}
	return model.NewSecurityContext(userId, userRoles)
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
	log.Verb("Before API auth check.")

	// Get security context
	secCtx := r.Context().Value(constant.ContextKeySecurityContext).(*model.SecurityContext)

	// Authenticated?
	if secCtx.UserId != model.AnonymousUserId {
		log.Debugf("User %d is authenticated.", secCtx.UserId)

		// Forward to next handler
		next(w, r)
	} else {
		log.Debug("User must authenticate.")

		// Request authentication
		log.Debug("Requesting authentication ...")
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", "Basic")
	}

	log.Verb("After API auth check.")
}

func getRequestPath(r *http.Request) string {
	p := r.URL.EscapedPath()
	return strings.TrimPrefix(p, constant.ApiPath)
}
