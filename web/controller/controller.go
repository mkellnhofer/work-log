package controller

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/constant"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/util/security"
	"kellnhofer.com/work-log/web/middleware"
)

func getContext(eCtx echo.Context) context.Context {
	return eCtx.Request().Context()
}

func getCurrentUserId(ctx context.Context) int {
	return security.GetCurrentUserId(ctx)
}

func getErrorCode(err error) int {
	code := e.SysUnknown
	if er, ok := err.(*e.Error); ok {
		code = er.Code
	}
	return code
}

func getIdPathVar(ctx echo.Context) (int, error) {
	v := ctx.Param("id")
	if v == "" {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (Variable missing.)")
		log.Debug(err.StackTrace())
		return 0, err
	}

	id, pErr := strconv.Atoi(v)
	if pErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID. (Variable must be numeric.)", pErr)
		log.Debug(err.StackTrace())
		return 0, err
	}

	if id <= 0 {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (Variable must be positive.)")
		log.Debug(err.StackTrace())
		return 0, err
	}

	return id, nil
}

func getErrorCodeQueryParam(ctx echo.Context) (int, error) {
	v := ctx.QueryParam("error")
	if v == "" {
		return 0, nil
	}

	ec, pErr := strconv.Atoi(v)
	if pErr != nil {
		err := e.WrapError(e.ValUnknown, "Invalid error code. (Variable must be numeric.)", pErr)
		log.Debug(err.StackTrace())
		return 0, err
	}

	return ec, nil
}

func getPageNumberQueryParam(ctx echo.Context) (int, bool, error) {
	v := ctx.QueryParam("page")
	if v == "" {
		return 0, false, nil
	}

	page, pErr := strconv.Atoi(v)
	if pErr != nil {
		err := e.WrapError(e.ValPageNumberInvalid, "Invalid page number. (Variable must be numeric.)",
			pErr)
		log.Debug(err.StackTrace())
		return 0, false, err
	}

	if page <= 0 {
		err := e.NewError(e.ValPageNumberInvalid, "Invalid page number. (Variable must be positive.)")
		log.Debug(err.StackTrace())
		return 0, false, err
	}

	return page, true, nil
}

func getSearchQueryParam(ctx echo.Context) (string, bool) {
	v := ctx.QueryParam("query")
	if v == "" {
		return "", false
	}
	return v, true
}

func getMonthQueryParam(ctx echo.Context) (int, int, bool, error) {
	v := ctx.QueryParam("month")
	if v == "" {
		return 0, 0, false, nil
	}

	if len(v) != 6 {
		err := e.NewError(e.ValMonthInvalid, "Invalid month. (Variable must have length of 6.)")
		log.Debug(err.StackTrace())
		return 0, 0, false, err
	}

	return parseMonthParam(v)
}

func parseMonthParam(v string) (int, int, bool, error) {
	if v == "" {
		return 0, 0, false, nil
	}

	yv := v[0:4]
	mv := v[4:6]

	year, err := strconv.Atoi(yv)
	if err != nil {
		err := e.NewError(e.ValMonthInvalid, "Invalid month. (Year part is invalid.)")
		log.Debug(err.StackTrace())
		return 0, 0, false, err
	}
	month, err := strconv.Atoi(mv)
	if err != nil || month <= 0 || month > 12 {
		err := e.NewError(e.ValMonthInvalid, "Invalid month. (Month part invalid.)")
		log.Debug(err.StackTrace())
		return 0, 0, false, err
	}

	return year, month, true, nil
}

func saveCurrentUrl(ctx echo.Context) {
	sh := getContext(ctx).Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)
	s := sh.Get()
	r := ctx.Request()
	path := r.URL.EscapedPath()
	query := r.URL.RawQuery
	req := path
	if query != "" {
		req = req + "?" + query
	}
	s.PreviousUrl = req
}

func getPreviousUrl(ctx echo.Context) string {
	sh := getContext(ctx).Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)
	s := sh.Get()
	if s.PreviousUrl != "" {
		return s.PreviousUrl
	} else {
		return constant.ViewPathDefault
	}
}

func writeFile(r *echo.Response, fileName string, wt io.WriterTo) error {
	// Write header
	r.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	r.Header().Set("Content-Type", "application/octet-stream")
	r.Header().Set("Content-Transfer-Encoding", "binary")
	r.Header().Set("Expires", "0")

	// Write body
	_, wErr := wt.WriteTo(r.Writer)
	if wErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not write response.", wErr)
		log.Debug(err.StackTrace())
		return err
	}

	return nil
}
