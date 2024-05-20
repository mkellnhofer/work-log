package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/util/security"
)

const pageSize = 7

const dateTimeFormat = "2006-01-02 15:04"

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

func getIdPathVar(ctx echo.Context) (int, error) {
	v := ctx.Param("id")
	if v == "" {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (Variable missing.)")
		log.Debug(err.StackTrace())
		return 0, err
	}

	return parseId(v, false)
}

func getTypeIdQueryParam(ctx echo.Context) (int, error) {
	v := ctx.QueryParam("type")
	if v == "" {
		return 0, nil
	}

	return parseId(v, false)
}

func parseId(in string, allowZero bool) (int, error) {
	id, cErr := strconv.Atoi(in)
	if cErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID. (ID must be numeric.)", cErr)
		log.Debug(err.StackTrace())
		return 0, err
	}
	if !allowZero && id <= 0 {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (ID must be positive.)")
		log.Debug(err.StackTrace())
		return 0, err
	}
	if id < 0 {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (ID must be zero or positive.)")
		log.Debug(err.StackTrace())
		return 0, err
	}
	return id, nil
}

func getPageNumberQueryParam(ctx echo.Context) (int, bool, error) {
	v := ctx.QueryParam("page")
	if v == "" {
		return 0, false, nil
	}

	pageNum, pErr := strconv.Atoi(v)
	if pErr != nil {
		err := e.WrapError(e.ValPageNumberInvalid, "Invalid page number. (Variable must be numeric.)",
			pErr)
		log.Debug(err.StackTrace())
		return 0, false, err
	}

	if pageNum <= 0 {
		err := e.NewError(e.ValPageNumberInvalid, "Invalid page number. (Variable must be positive.)")
		log.Debug(err.StackTrace())
		return 0, false, err
	}

	return pageNum, true, nil
}

func calculateOffsetLimitFromPageNumber(pageNum int) (int, int) {
	page := pageNum
	if page == 0 {
		page = 1
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	return offset, limit
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

	return parseMonth(v)
}

func parseMonth(v string) (int, int, bool, error) {
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

func parseDateTime(inDate string, inTime string, code int) (time.Time, error) {
	dt := inDate + " " + inTime
	out, pErr := time.ParseInLocation(dateTimeFormat, dt, time.Local)
	if pErr != nil {
		err := e.WrapError(code, fmt.Sprintf("Could not parse time %s.", inTime), pErr)
		log.Debug(err.StackTrace())
		return time.Now(), err
	}
	return out, nil
}

func validateStringLength(in string, length int, code int) error {
	if len(in) > length {
		err := e.NewError(code, fmt.Sprintf("String too long. (Must be <= %d characters.)", length))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func calculateNumberOfTotalPages(items int, pageSize int) int {
	page := items / pageSize
	remaining := items % pageSize
	if remaining > 0 {
		return page + 1
	}
	return page
}
