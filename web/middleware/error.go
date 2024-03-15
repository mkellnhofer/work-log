package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
)

// ErrorMiddleware catches errors and forwards to the error page.
type ErrorMiddleware struct {
}

// NewErrorMiddleware create a new error middleware.
func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

// CreateHandler creates a new handler to process requests.
func (m *ErrorMiddleware) CreateHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Verb("Before error check.")

		err := m.process(next, c)

		log.Verb("After error check.")

		return err
	}
}

func (m *ErrorMiddleware) process(next echo.HandlerFunc, c echo.Context) error {
	err := next(c)
	if err != nil {
		log.Verb("Catching error.")
		m.handleError(c, err)
	}
	return nil
}

func (m *ErrorMiddleware) handleError(c echo.Context, err interface{}) {
	var ec int
	switch tErr := err.(type) {
	case *e.Error:
		ec = tErr.Code
	default:
		log.Errorf("%s", tErr)
		ec = e.SysUnknown
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/error?error=%d", ec))
}
