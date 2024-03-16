package controller

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"kellnhofer.com/work-log/pkg/constant"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/pkg/util/security"
	view "kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/middleware"
	vm "kellnhofer.com/work-log/web/model"
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
func (c *AuthController) GetLoginHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /login.")
		return c.handleShowLogin(eCtx)
	}
}

// PostLoginHandler returns a handler for "POST /login".
func (c *AuthController) PostLoginHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /login.")
		return c.handleExecuteLogin(eCtx)
	}
}

// GetPasswordChangeHandler returns a handler for "GET /password_change".
func (c *AuthController) GetPasswordChangeHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /password_change.")
		return c.handleShowPasswordChange(eCtx)
	}
}

// PostPasswordChangeHandler returns a handler for "POST /password_change".
func (c *AuthController) PostPasswordChangeHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /password_change.")
		return c.handleExecutePasswordChange(eCtx)
	}
}

// GetLogoutHandler returns a handler for "GET /logout".
func (c *AuthController) GetLogoutHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /logout.")
		return c.handleExecuteLogout(eCtx)
	}
}

// --- Login handler functions ---

func (c *AuthController) handleShowLogin(eCtx echo.Context) error {
	// Get error code
	ec, err := getErrorCodeQueryParam(eCtx)
	if err != nil {
		return err
	}

	// Create view model
	model := c.createShowLoginViewModel(ec)

	// Render
	return view.RenderLoginTemplate(eCtx.Response(), model)
}

func (c *AuthController) createShowLoginViewModel(ec int) *vm.Login {
	lvm := vm.NewLogin()
	if ec != 0 {
		lvm.ErrorMessage = loc.GetErrorMessageString(ec)
	}
	return lvm
}

func (c *AuthController) handleExecuteLogin(eCtx echo.Context) error {
	// Create system context
	syseCtx := security.CreateSystemContext(getContext(eCtx))

	// Get form inputs
	username := eCtx.FormValue("username")
	password := eCtx.FormValue("password")

	log.Debugf("User %s is trying to authenticate ...", username)

	// Validate inputs
	if err := validateLoginInputs(username, password); err != nil {
		return c.handleLoginError(eCtx, err)
	}

	// Find user
	user, err := c.uServ.GetUserByUsername(syseCtx, username)
	if err != nil {
		return err
	}
	if user == nil {
		err := e.NewError(e.AuthCredentialsInvalid, fmt.Sprintf("Invalid credentials. (Unknown "+
			"username %s.)", username))
		log.Debug(err.StackTrace())
		return c.handleLoginError(eCtx, err)
	}

	// Check password
	if cpwErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); cpwErr != nil {
		err := e.WrapError(e.AuthCredentialsInvalid, "Invalid credentials. (Wrong password.)", cpwErr)
		log.Debug(err.StackTrace())
		return c.handleLoginError(eCtx, err)
	}

	log.Debugf("User %s has successfully authenticated.", username)

	// Get session holder from context
	sessHolder := getContext(eCtx).Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)

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
	eCtx.SetCookie(sessCookie)

	// Was a previous request stored?
	if preSess != nil && preSess.PreviousUrl != "" {
		// Redirect to previous path
		return c.handleLoginSuccess(eCtx, preSess.PreviousUrl)
	} else {
		// Redirect to root path
		return c.handleLoginSuccess(eCtx, "/")
	}
}

func (c *AuthController) handleLoginSuccess(eCtx echo.Context, url string) error {
	return eCtx.Redirect(http.StatusFound, url)
}

func (c *AuthController) handleLoginError(eCtx echo.Context, err error) error {
	code := e.SysUnknown
	if er, ok := err.(*e.Error); ok {
		code = er.Code
	}
	return eCtx.Redirect(http.StatusFound, fmt.Sprintf("/login?error=%d", code))
}

// --- Password change handler functions ---

func (c *AuthController) handleShowPasswordChange(eCtx echo.Context) error {
	// Get error code
	ec, err := getErrorCodeQueryParam(eCtx)
	if err != nil {
		return err
	}

	// Create view model
	model := c.createShowPasswordChangeViewModel(ec)

	// Render
	return view.RenderPasswordChangeTemplate(eCtx.Response(), model)
}

func (c *AuthController) createShowPasswordChangeViewModel(ec int) *vm.PasswordChange {
	pcvm := vm.NewPasswordChange()
	if ec != 0 {
		pcvm.ErrorMessage = loc.GetErrorMessageString(ec)
	}
	return pcvm
}

func (c *AuthController) handleExecutePasswordChange(eCtx echo.Context) error {
	// Get current user ID
	userId := getCurrentUserId(getContext(eCtx))

	log.Debugf("User %d is trying to change password ...", userId)

	// Get form inputs
	password1 := eCtx.FormValue("password1")
	password2 := eCtx.FormValue("password2")

	// Validate inputs
	if err := validatePasswordChangeInputs(password1, password2); err != nil {
		return c.handlePasswordChangeError(eCtx, err)
	}

	// Update password
	upwErr := c.uServ.UpdateCurrentUserPassword(getContext(eCtx), password1)
	if upwErr != nil {
		return upwErr
	}

	log.Debugf("User %d has successfully changed password.", userId)

	// Get session holder from context
	sessHolder := getContext(eCtx).Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)

	// Get current session
	preSess := sessHolder.Get()

	// Was a previous request stored?
	if preSess != nil && preSess.PreviousUrl != "" {
		// Redirect to previous path
		return c.handlePasswordChangeSuccess(eCtx, preSess.PreviousUrl)
	} else {
		// Redirect to root path
		return c.handlePasswordChangeSuccess(eCtx, "/")
	}
}

func (c *AuthController) handlePasswordChangeSuccess(eCtx echo.Context, url string) error {
	return eCtx.Redirect(http.StatusFound, url)
}

func (c *AuthController) handlePasswordChangeError(eCtx echo.Context, err error) error {
	code := e.SysUnknown
	if er, ok := err.(*e.Error); ok {
		code = er.Code
	}
	return eCtx.Redirect(http.StatusFound, fmt.Sprintf("/password-change?error=%d", code))
}

// --- Logout handler functions ---

func (c *AuthController) handleExecuteLogout(eCtx echo.Context) error {
	// Get session holder from context
	sessHolder := getContext(eCtx).Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)

	// Close session
	sessHolder.Clear()

	// Redirect to login page
	return c.handleLogoutSuccess(eCtx)
}

func (c *AuthController) handleLogoutSuccess(eCtx echo.Context) error {
	return eCtx.Redirect(http.StatusFound, "/login")
}

// --- Validator functions ---

func validateLoginInputs(username string, password string) error {
	if len(username) > model.MaxLengthUserUsername || len(password) > model.MaxLengthUserPassword {
		err := e.NewError(e.AuthCredentialsInvalid, "Invalid credentials. (Invalid username or "+
			"password length.)")
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func validatePasswordChangeInputs(password1 string, password2 string) error {
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
