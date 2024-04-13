package web

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
)

// Render renders a template.
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	tErr := t.Render(ctx.Request().Context(), ctx.Response().Writer)
	if tErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not render template.", tErr)
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}
