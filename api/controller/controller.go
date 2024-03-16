package controller

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	am "kellnhofer.com/work-log/api/model"
	"kellnhofer.com/work-log/pkg/constant"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	m "kellnhofer.com/work-log/pkg/model"
	httputil "kellnhofer.com/work-log/pkg/util/http"
	"kellnhofer.com/work-log/pkg/util/security"
)

const defaultPageSize = 50

// Information about the error.
// swagger:response ErrorResponse
type Error struct {
	// in: body
	Body am.Error
}

func getContext(eCtx echo.Context) context.Context {
	return eCtx.Request().Context()
}

func hasCurrentUserRight(ctx context.Context, right m.Right) bool {
	return security.HasCurrentUserRight(ctx, right)
}

func getIdPathVar(eCtx echo.Context) (int, error) {
	v := eCtx.Param("id")
	if v == "" {
		err := e.NewError(e.ValIdInvalid, "Invalid ID variable. (Varible missing.)")
		log.Debug(err.StackTrace())
		return 0, err
	}

	id, pErr := strconv.Atoi(v)
	if pErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID variable. (Varible must be numeric.)", pErr)
		log.Debug(err.StackTrace())
		return 0, err
	}

	if id <= 0 {
		err := e.NewError(e.ValIdInvalid, "Invalid ID variable. (Varible must be positive.)")
		log.Debug(err.StackTrace())
		return 0, err
	}

	return id, nil
}

func getFilterQueryParam(eCtx echo.Context) string {
	return eCtx.QueryParam("filter")
}

func getSortQueryParam(eCtx echo.Context) string {
	return eCtx.QueryParam("sort")
}

func getOffsetQueryParam(eCtx echo.Context) (int, error) {
	i, err := getIntQueryParam(eCtx, "offset")
	if err != nil {
		err := e.NewError(e.ValOffsetInvalid, "Invalid offset. (Offset must be numeric (int32).)")
		log.Debug(err.StackTrace())
		return 0, err
	}
	if i < 0 {
		err := e.NewError(e.ValOffsetInvalid, "Invalid offset. (Offset must be positive.)")
		log.Debug(err.StackTrace())
		return 0, err
	}
	return i, nil
}

func getLimitQueryParam(eCtx echo.Context) (int, error) {
	i, err := getIntQueryParam(eCtx, "limit")
	if err != nil {
		err := e.NewError(e.ValLimitInvalid, "Invalid limit. (Limit must be numeric (int32).)")
		log.Debug(err.StackTrace())
		return 0, err
	}
	if i < 0 {
		err := e.NewError(e.ValLimitInvalid, "Invalid limit. (Limit must be positive.)")
		log.Debug(err.StackTrace())
		return 0, err
	}
	return i, nil
}

func getIntQueryParam(eCtx echo.Context, name string) (int, error) {
	qv := eCtx.QueryParam(name)
	if qv == "" {
		return 0, nil
	}
	return strconv.Atoi(qv)
}

func parseTimestamp(ts string) (time.Time, error) {
	return time.ParseInLocation(constant.ApiTimestampFormat, ts, time.Local)
}

func readRequestBody(eCtx echo.Context, data interface{}) error {
	return httputil.ReadHttpRequestBody(eCtx.Request(), data)
}

func writeResponse(eCtx echo.Context, statusCode int, data interface{}) error {
	return httputil.WriteHttpResponse(eCtx.Response(), statusCode, data)
}
