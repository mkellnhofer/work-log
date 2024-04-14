package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/web"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/view/hx"
	"kellnhofer.com/work-log/web/view/pages"
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

		isHtmx := web.IsHtmxRequest(ctx)

		if !isHtmx {
			return c.handleShowError(ctx)
		} else {
			return c.handleHxShowError(ctx)
		}
	}
}

// --- Handler functions ---

func (c *ErrorController) handleShowError(eCtx echo.Context) error {
	// Get view data
	model := c.getErrorViewData(eCtx)
	// Render
	return web.Render(eCtx, http.StatusOK, pages.ErrorPage(model))
}

func (c *ErrorController) handleHxShowError(eCtx echo.Context) error {
	// Get view data
	model := c.getErrorViewData(eCtx)
	// Render
	return web.Render(eCtx, http.StatusOK, hx.ErrorPage(model))
}

func (c *ErrorController) getErrorViewData(eCtx echo.Context) *vm.Error {
	// Get error code
	ec, err := getErrorCodeQueryParam(eCtx)
	if err != nil {
		ec = e.SysUnknown
	}

	// Create view model
	em := loc.GetErrorMessageString(ec)
	return vm.NewError(em)
}
