package controller

import (
	"net/http"

	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/view"
	vm "kellnhofer.com/work-log/view/model"
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
	ec, err := getIntQueryParam(r, "error")
	if err != nil {
		panic(err)
	}

	// Create view model
	em := getErrorMessage(ec)
	model := vm.NewError(em)

	// Render
	view.RenderErrorTemplate(w, model)
}
