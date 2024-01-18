package controller

import (
	"net/http"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	view "kellnhofer.com/work-log/web"
	vm "kellnhofer.com/work-log/web/model"
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
func (c *ErrorController) GetErrorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /error.")
		c.handleShowError(w, r)
	}
}

// --- Handler functions ---

func (c *ErrorController) handleShowError(w http.ResponseWriter, r *http.Request) {
	// Get error code
	ecqp := getErrorCodeQueryParam(r)
	ec := e.SysUnknown
	if ecqp != nil {
		ec = *ecqp
	}

	// Create view model
	em := loc.GetErrorMessageString(ec)
	model := vm.NewError(em)

	// Render
	view.RenderErrorTemplate(w, model)
}
