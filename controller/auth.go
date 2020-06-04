package controller

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/loc"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/middleware"
	"kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/service"
	"kellnhofer.com/work-log/view"
	vm "kellnhofer.com/work-log/view/model"
)

// AuthController handles requests for login/logout endpoints.
type AuthController struct {
	uServ *service.UserService
}

// NewAuthController creates a new auth controller.
func NewAuthController(uServ *service.UserService) *AuthController {
	return &AuthController{uServ}
}

// --- Endpoints ---

// GetLoginHandler returns a handler for "GET /login".
func (c *AuthController) GetLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /login.")
		c.handleShowLogin(w, r)
	}
}

// PostLoginHandler returns a handler for "POST /login".
func (c *AuthController) PostLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /login.")
		c.handleExecuteLogin(w, r)
	}
}

// PostLogoutHandler returns a handler for "GET /logout".
func (c *AuthController) GetLogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /logout.")
		c.handleExecuteLogout(w, r)
	}
}

// --- Login handler functions ---

func (c *AuthController) handleShowLogin(w http.ResponseWriter, r *http.Request) {
	// Get error code
	ecqp := getErrorCodeQueryParam(r)
	ec := 0
	if ecqp != nil {
		ec = *ecqp
	}

	// Create view model
	model := c.createShowLoginViewModel(ec)

	// Render
	view.RenderLoginTemplate(w, model)
}

func (c *AuthController) createShowLoginViewModel(ec int) *vm.Login {
	lvm := vm.NewLogin()
	if ec != 0 {
		lvm.ErrorMessage = loc.GetErrorMessageString(ec)
	}
	return lvm
}

func (c *AuthController) handleExecuteLogin(w http.ResponseWriter, r *http.Request) {
	// Get form inputs
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Find user
	user, guErr := c.uServ.GetUserByUsername(username)
	if guErr != nil {
		panic(guErr)
	}
	if user == nil {
		err := e.NewError(e.AuthInvalidCredentials, fmt.Sprintf("Invalid credentials. (Unknown "+
			"username %s.)", username))
		log.Debug(err.StackTrace())
		c.handleLoginError(w, r, err)
		return
	}

	// Check password
	if cpwErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); cpwErr != nil {
		err := e.WrapError(e.AuthInvalidCredentials, "Invalid credentials. (Wrong password.)", cpwErr)
		log.Debug(err.StackTrace())
		c.handleLoginError(w, r, err)
		return
	}

	// Get session holder from context
	sessHolder := r.Context().Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)

	// Get current session
	preSess := sessHolder.Get()

	// Create new session
	newSess := model.NewSession()
	newSess.UserId = user.Id

	// Set new session
	sessHolder.Set(newSess)

	// Set session cookie
	sessCookie := &http.Cookie{Name: constant.SessionCookieName, Value: newSess.Id, Path: "/",
		HttpOnly: true}
	http.SetCookie(w, sessCookie)

	// Was a previous request stored?
	if preSess != nil && preSess.PreviousUrl != "" {
		// Redirect to previous path
		c.handleLoginSuccess(w, r, preSess.PreviousUrl)
	} else {
		// Redirect to root path
		c.handleLoginSuccess(w, r, "/")
	}
}

func (c *AuthController) handleLoginSuccess(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusFound)
}

func (c *AuthController) handleLoginError(w http.ResponseWriter, r *http.Request, err *e.Error) {
	http.Redirect(w, r, fmt.Sprintf("/login?error=%d", err.Code), http.StatusFound)
}

// --- Logout handler functions ---

func (c *AuthController) handleExecuteLogout(w http.ResponseWriter, r *http.Request) {
	// Get session holder from context
	sessHolder := r.Context().Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)

	// Close session
	sessHolder.Clear()

	// Redirect to login page
	c.handleLogoutSuccess(w, r)
}

func (c *AuthController) handleLogoutSuccess(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
