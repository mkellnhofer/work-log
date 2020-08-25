package controller

import (
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/loc"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/service"
	"kellnhofer.com/work-log/util/security"
	"kellnhofer.com/work-log/view"
	"kellnhofer.com/work-log/view/middleware"
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

// GetPasswordChangeHandler returns a handler for "GET /password_change".
func (c *AuthController) GetPasswordChangeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /password_change.")
		c.handleShowPasswordChange(w, r)
	}
}

// PostPasswordChangeHandler returns a handler for "POST /password_change".
func (c *AuthController) PostPasswordChangeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /password_change.")
		c.handleExecutePasswordChange(w, r)
	}
}

// GetLogoutHandler returns a handler for "GET /logout".
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
	// Create system context
	sysCtx := security.CreateSystemContext(r.Context())

	// Get form inputs
	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Debugf("User %s is trying to authenticate ...", username)

	// Validate inputs
	if err := validateLoginInputs(username, password); err != nil {
		c.handleLoginError(w, r, err)
		return
	}

	// Find user
	user, guErr := c.uServ.GetUserByUsername(sysCtx, username)
	if guErr != nil {
		panic(guErr)
	}
	if user == nil {
		err := e.NewError(e.AuthCredentialsInvalid, fmt.Sprintf("Invalid credentials. (Unknown "+
			"username %s.)", username))
		log.Debug(err.StackTrace())
		c.handleLoginError(w, r, err)
		return
	}

	// Check password
	if cpwErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); cpwErr != nil {
		err := e.WrapError(e.AuthCredentialsInvalid, "Invalid credentials. (Wrong password.)", cpwErr)
		log.Debug(err.StackTrace())
		c.handleLoginError(w, r, err)
		return
	}

	log.Debugf("User %s has successfully authenticated.", username)

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

// --- Password change handler functions ---

func (c *AuthController) handleShowPasswordChange(w http.ResponseWriter, r *http.Request) {
	// Get error code
	ecqp := getErrorCodeQueryParam(r)
	ec := 0
	if ecqp != nil {
		ec = *ecqp
	}

	// Create view model
	model := c.createShowPasswordChangeViewModel(ec)

	// Render
	view.RenderPasswordChangeTemplate(w, model)
}

func (c *AuthController) createShowPasswordChangeViewModel(ec int) *vm.PasswordChange {
	pcvm := vm.NewPasswordChange()
	if ec != 0 {
		pcvm.ErrorMessage = loc.GetErrorMessageString(ec)
	}
	return pcvm
}

func (c *AuthController) handleExecutePasswordChange(w http.ResponseWriter, r *http.Request) {
	// Get context
	ctx := r.Context()

	// Get current user ID
	userId := getCurrentUserId(ctx)

	log.Debugf("User %d is trying to change password ...", userId)

	// Get form inputs
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")

	// Validate inputs
	if err := validatePasswordChangeInputs(password1, password2); err != nil {
		c.handlePasswordChangeError(w, r, err)
		return
	}

	// Update password
	upwErr := c.uServ.UpdateCurrentUserPassword(ctx, password1)
	if upwErr != nil {
		panic(upwErr)
	}

	log.Debugf("User %d has successfully changed password.", userId)

	// Get session holder from context
	sessHolder := r.Context().Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)

	// Get current session
	preSess := sessHolder.Get()

	// Was a previous request stored?
	if preSess != nil && preSess.PreviousUrl != "" {
		// Redirect to previous path
		c.handlePasswordChangeSuccess(w, r, preSess.PreviousUrl)
	} else {
		// Redirect to root path
		c.handlePasswordChangeSuccess(w, r, "/")
	}
}

func (c *AuthController) handlePasswordChangeSuccess(w http.ResponseWriter, r *http.Request,
	url string) {
	http.Redirect(w, r, url, http.StatusFound)
}

func (c *AuthController) handlePasswordChangeError(w http.ResponseWriter, r *http.Request,
	err *e.Error) {
	http.Redirect(w, r, fmt.Sprintf("/password-change?error=%d", err.Code), http.StatusFound)
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

// --- Validator functions ---

func validateLoginInputs(username string, password string) *e.Error {
	if len(username) > model.MaxLengthUserUsername || len(password) > model.MaxLengthUserPassword {
		err := e.NewError(e.AuthCredentialsInvalid, "Invalid credentials. (Invalid username or "+
			"password length.)")
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func validatePasswordChangeInputs(password1 string, password2 string) *e.Error {
	if len(password1) == 0 {
		err := e.NewError(e.ValPasswordEmpty, "Password must not be empty.")
		log.Debug(err.StackTrace())
		return err
	}
	if len(password1) < model.MinLengthUserPassword {
		err := e.NewError(e.ValPasswordTooShort, fmt.Sprintf("Password must be at least %d long.",
			model.MinLengthUserPassword))
		log.Debug(err.StackTrace())
		return err
	}
	if len(password1) > model.MaxLengthUserPassword {
		err := e.NewError(e.ValPasswordTooLong, fmt.Sprintf("Password must not be longer than %d.",
			model.MaxLengthUserPassword))
		log.Debug(err.StackTrace())
		return err
	}
	r := regexp.MustCompile("^[" + model.ValidUserPasswordCharacters + "]+$")
	if !r.MatchString(password1) {
		err := e.NewError(e.ValPasswordInvalid, "Password contains contains illegal character.")
		log.Debug(err.StackTrace())
		return err
	}
	if password1 != password2 {
		err := e.NewError(e.ValPasswordsNotMatching, "Passwords do not match.")
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}
