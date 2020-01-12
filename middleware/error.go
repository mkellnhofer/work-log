package middleware

import (
	"fmt"
	"net/http"

	"kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
)

// ErrorMiddleware catches errors and forwards to the error page.
type ErrorMiddleware struct {
}

// NewErrorMiddleware create a new error middleware.
func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

// ServeHTTP processes requests.
func (m *ErrorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Verb("Before error check.")

	defer func() {
		if err := recover(); err != nil {
			log.Verb("Catching error.")
			m.handleError(w, r, err)
		}
	}()

	next(w, r)

	log.Verb("After error check.")
}

func (m *ErrorMiddleware) handleError(w http.ResponseWriter, r *http.Request, err interface{}) {
	var ec int
	switch e := err.(type) {
	case *error.Error:
		ec = e.Code
	default:
		log.Errorf("%s", e)
		ec = error.SysUnknown
	}

	http.Redirect(w, r, fmt.Sprintf("/error?code=%d", ec), http.StatusFound)
}
