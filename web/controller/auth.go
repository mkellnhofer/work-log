package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"kellnhofer.com/work-log/pkg/constant"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/pkg/util/security"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/middleware"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/view/hx"
	"kellnhofer.com/work-log/web/view/pages"
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

		if !web.IsHtmxRequest(eCtx) {
			err := e.NewError(e.ValUnknown, "Not a HTMX request.")
			log.Debug(err.StackTrace())
			return err
		}

		return c.handleExecuteLogin(eCtx)
	}
}

// GetLogoutHandler returns a handler for "GET /logout".
func (c *AuthController) GetLogoutHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /logout.")
		return c.handleExecuteLogout(eCtx)
	}
}

// --- Handler functions ---

func (c *AuthController) handleShowLogin(eCtx echo.Context) error {
	// Get security context
	secCtx := security.GetSecurityContext(getContext(eCtx))

	// If user is not authenticated: Show form to enter credentials
	if secCtx.IsAnonymousUser() {
		return c.showEnterCredentials(eCtx)
	}

	// Get user
	sysCtx := security.CreateSystemContext(getContext(eCtx))
	userId := secCtx.UserId
	user, err := c.uServ.GetUserById(sysCtx, userId)
	if err != nil {
		return err
	}

	// If user was not found: Abort
	if user == nil {
		err := e.NewError(e.AuthUnknown, "User is not authenticated.")
		log.Debug(err.StackTrace())
		return err
	}

	// If user must change password: Show form to change password
	if user.MustChangePassword {
		return c.showChangePassword(eCtx)
	}

	// Redirect to home page
	return c.redirectHome(eCtx)
}

func (c *AuthController) handleExecuteLogin(eCtx echo.Context) error {
	// Get step value
	s := eCtx.FormValue("step")
	step, cErr := strconv.Atoi(s)
	if cErr != nil {
		err := e.WrapError(e.AuthDataInvalid, "Invalid login step.", cErr)
		log.Debug(err.StackTrace())
		return err
	}

	// Handle specific step
	switch step {
	case vm.LoginStepEnterCredentials:
		return c.handleEnterCredentials(eCtx)
	case vm.LoginStepChangePassword:
		return c.handleChangePassword(eCtx)
	default:
		err := e.WrapError(e.AuthDataInvalid, "Invalid login step.", cErr)
		log.Debug(err.StackTrace())
		return err
	}
}

func (c *AuthController) handleEnterCredentials(eCtx echo.Context) error {
	// Get form inputs
	username := eCtx.FormValue("username")
	password := eCtx.FormValue("password")

	log.Debugf("User %s is trying to authenticate ...", username)

	// If inputs are invalid: Show error
	if err := c.validateEnterCredentialsInputs(username, password); err != nil {
		return c.showEnterCredentialsError(eCtx, err)
	}

	// Find user
	sysCtx := security.CreateSystemContext(getContext(eCtx))
	user, err := c.uServ.GetUserByUsername(sysCtx, username)
	if err != nil {
		return err
	}

	// If no user was found: Show error
	if user == nil {
		err := e.NewError(e.AuthCredentialsInvalid, "Invalid credentials. (Unknown username.)")
		log.Debug(err.StackTrace())
		return c.showEnterCredentialsError(eCtx, err)
	}

	// If password does not match: Show error
	if err := c.checkPassword(user.Password, password); err != nil {
		return c.showEnterCredentialsError(eCtx, err)
	}

	log.Debugf("User %s has successfully authenticated.", username)

	// Create new session
	sess := c.createNewSession(eCtx, user.Id)

	// If user must change password: Show form to change password
	if user.MustChangePassword {
		return c.showChangePassword(eCtx)
	}

	// Redirect user to saved URL
	return c.redirectSavedUrl(eCtx, sess)
}

func (c *AuthController) validateEnterCredentialsInputs(username string, password string) error {
	if len(username) > model.MaxLengthUserUsername || len(password) > model.MaxLengthUserPassword {
		err := e.NewError(e.AuthCredentialsInvalid, "Invalid credentials. (Invalid username or "+
			"password length.)")
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func (c *AuthController) checkPassword(storedPassword string, enteredPassword string) error {
	cpwErr := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(enteredPassword))
	if cpwErr != nil {
		err := e.WrapError(e.AuthCredentialsInvalid, "Invalid credentials. (Wrong password.)",
			cpwErr)
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func (c *AuthController) showEnterCredentialsError(eCtx echo.Context, err error) error {
	// Get error message
	ec := getErrorCode(err)
	em := loc.GetErrorMessageString(ec)
	// Render
	return web.Render(eCtx, http.StatusOK, hx.LoginPage(vm.LoginStepEnterCredentials, em))
}

func (c *AuthController) handleChangePassword(eCtx echo.Context) error {
	// Get security context
	secCtx := security.GetSecurityContext(getContext(eCtx))

	// If user is not authenticated: Abort
	if secCtx.IsAnonymousUser() {
		err := e.NewError(e.AuthUnknown, "User is not authenticated.")
		log.Debug(err.StackTrace())
		return err
	}

	log.Debugf("User %d is trying to change password ...", secCtx.UserId)

	// Get form inputs
	password1 := eCtx.FormValue("password1")
	password2 := eCtx.FormValue("password2")

	// If inputs are invalid: Show error
	if err := c.validateChangePasswordInputs(password1, password2); err != nil {
		return c.handleChangePasswordError(eCtx, err)
	}

	// Update password
	upwErr := c.uServ.UpdateCurrentUserPassword(getContext(eCtx), password1)
	if upwErr != nil {
		return upwErr
	}

	log.Debugf("User %d has successfully changed password.", secCtx.UserId)

	// Get current session
	sess := c.getCurrentSession(eCtx)

	// Redirect user to saved URL
	return c.redirectSavedUrl(eCtx, sess)
}

func (c *AuthController) validateChangePasswordInputs(password1 string, password2 string) error {
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

func (c *AuthController) handleChangePasswordError(eCtx echo.Context, err error) error {
	// Get error message
	ec := getErrorCode(err)
	em := loc.GetErrorMessageString(ec)
	// Render
	return web.Render(eCtx, http.StatusOK, hx.LoginPage(vm.LoginStepChangePassword, em))
}

func (c *AuthController) createNewSession(eCtx echo.Context, userId int) *model.Session {
	// Get session holder from context
	sessHolder := getContext(eCtx).Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)

	// Get current session
	preSess := sessHolder.Get()

	// Create new session
	newSess := model.NewSession()
	newSess.UserId = userId
	newSess.PreviousUrl = preSess.PreviousUrl

	// Set new session
	sessHolder.Set(newSess)

	// Set session cookie
	sessCookie := &http.Cookie{Name: constant.SessionCookieName, Value: newSess.Id, Path: "/",
		HttpOnly: true}
	eCtx.SetCookie(sessCookie)

	return newSess
}

func (c *AuthController) getCurrentSession(eCtx echo.Context) *model.Session {
	// Get session holder from context
	sessHolder := getContext(eCtx).Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)

	// Get current session
	return sessHolder.Get()
}

func (c *AuthController) showEnterCredentials(eCtx echo.Context) error {
	if web.IsHtmxRequest(eCtx) {
		return web.Render(eCtx, http.StatusOK, hx.LoginPage(vm.LoginStepEnterCredentials, ""))
	} else {
		return web.Render(eCtx, http.StatusOK, pages.LoginPage(vm.LoginStepEnterCredentials))
	}
}

func (c *AuthController) showChangePassword(eCtx echo.Context) error {
	if web.IsHtmxRequest(eCtx) {
		return web.Render(eCtx, http.StatusOK, hx.LoginPage(vm.LoginStepChangePassword, ""))
	} else {
		return web.Render(eCtx, http.StatusOK, pages.LoginPage(vm.LoginStepChangePassword))
	}
}

func (c *AuthController) redirectHome(eCtx echo.Context) error {
	return c.redirect(eCtx, constant.ViewPathDefault)
}

func (c *AuthController) redirectSavedUrl(eCtx echo.Context, sess *model.Session) error {
	url := constant.ViewPathDefault
	if sess != nil && sess.PreviousUrl != "" {
		url = sess.PreviousUrl
		sess.PreviousUrl = ""
	}
	return c.redirect(eCtx, url)
}

func (c *AuthController) redirect(eCtx echo.Context, url string) error {
	if web.IsHtmxRequest(eCtx) {
		web.HtmxRedirectUrl(eCtx, url)
		return eCtx.NoContent(http.StatusOK)
	} else {
		return eCtx.Redirect(http.StatusFound, url)
	}
}

func (c *AuthController) handleExecuteLogout(eCtx echo.Context) error {
	// Get session holder from context
	sessHolder := getContext(eCtx).Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)

	// Close session
	sessHolder.Clear()

	// Redirect to login page
	return eCtx.Redirect(http.StatusFound, "/login")
}
