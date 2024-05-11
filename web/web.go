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

// Render renders a page template.
func RenderPage(ctx echo.Context, statusCode int, t templ.Component) error {
	return render(ctx, statusCode, true, t)
}

// RenderHx renders a HTMX template.
func RenderHx(ctx echo.Context, statusCode int, t templ.Component) error {
	return render(ctx, statusCode, false, t)
}

var cspRules = "default-src 'self';" +
	"script-src 'self';" +
	"img-src 'self' data:;" +
	"style-src 'self' 'unsafe-inline';" +
	"font-src 'self';" +
	"form-action 'self';" +
	"base-uri 'self';"

func render(ctx echo.Context, statusCode int, addCspHeader bool, t templ.Component) error {
	req := ctx.Request()
	res := ctx.Response()

	res.Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	res.Header().Add(echo.HeaderCacheControl, "no-store")

	if addCspHeader {
		res.Header().Add("Content-Security-Policy", cspRules)
		res.Header().Add("X-Content-Security-Policy", cspRules)
	}

	res.Writer.WriteHeader(statusCode)

	tErr := t.Render(req.Context(), res.Writer)
	if tErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not render template.", tErr)
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}
