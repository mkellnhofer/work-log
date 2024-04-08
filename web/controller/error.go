package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/web"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/pages"
)

// ErrorController handles requests for error endpoints.
type ErrorController struct {
}

// NewErrorController creates a new error controller.
func NewErrorController() *ErrorController {
	return &ErrorController{}
}

// --- Endpoints ---

// GetErrorHandler returns a handler for "GET /error".
func (c *ErrorController) GetErrorHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Verb("Handle GET /error.")
		return c.handleShowError(ctx)
	}
}

// --- Handler functions ---

func (c *ErrorController) handleShowError(eCtx echo.Context) error {
	// Get error code
	ec, err := getErrorCodeQueryParam(eCtx)
	if err != nil {
		ec = e.SysUnknown
	}

	// Create view model
	em := loc.GetErrorMessageString(ec)
	model := vm.NewError(em)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.ErrorPage(model))
}
