package middleware

import (
	"context"
	"fmt"
	"net/http"

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
		userId = m.authenticateUser(sysCtx, username, password)
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
	password string) int {
	user, guErr := m.uServ.GetUserByUsername(ctx, username)
	if guErr != nil {
		panic(guErr)
	}
	if user == nil {
		err := e.NewError(e.AuthInvalidCredentials, fmt.Sprintf("Invalid credentials. (Unknown "+
			"username '%s'.)", username))
		log.Debug(err.StackTrace())
		panic(err)
	}

	if cpwErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); cpwErr != nil {
		err := e.WrapError(e.AuthInvalidCredentials, "Invalid credentials. (Wrong password.)", cpwErr)
		log.Debug(err.StackTrace())
		panic(err)
	}

	return user.Id
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