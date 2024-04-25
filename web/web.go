package web

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
)

// IsHtmxRequest returns true if request is a HTMX request.
func IsHtmxRequest(ctx echo.Context) bool {
	return ctx.Request().Header.Get("HX-Request") == "true"
}

// HtmxPushUrl sets the response header "HX-Push-Url" which instructs HTMX to push the supplied URL
// into the browser's page history.
func HtmxPushUrl(ctx echo.Context, url string) {
	ctx.Response().Header().Set("HX-Push-Url", url)
}

// HtmxRedirectUrl sets the response headers "HX-Redirect" and "HX-Push-Url" which instructs HTMX to
// do a client-side redirect to the supplied URL.
func HtmxRedirectUrl(ctx echo.Context, url string) {
	ctx.Response().Header().Add("HX-Redirect", url)
	ctx.Response().Header().Add("HX-Push-Url", url)
}

// HtmxRetarget sets the response headers "HX-Retarget" which instructs HTMX to load the response
// content into the element with the supplied CSS target selector.
func HtmxRetarget(ctx echo.Context, target string) {
	ctx.Response().Header().Set("HX-Retarget", target)
}

// HtmxRetarget sets the response headers "HX-Trigger" which instructs HTMX to trigger a client-side
// event.
func HtmxTrigger(ctx echo.Context, trigger string) {
	ctx.Response().Header().Add("HX-Trigger", trigger)
}

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
