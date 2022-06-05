package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	am "kellnhofer.com/work-log/api/model"
	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	m "kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/util/security"
)

const defaultPageSize = 50

// Information about the error.
// swagger:response ErrorResponse
type Error struct {
	// in: body
	Body am.Error
}

func getCurrentUserId(ctx context.Context) int {
	return security.GetCurrentUserId(ctx)
}

func hasCurrentUserRight(ctx context.Context, right m.Right) bool {
	return security.HasCurrentUserRight(ctx, right)
}

func getIdPathVar(r *http.Request) int {
	v, ok := getPathVar(r, "id")
	if !ok {
		err := e.NewError(e.ValIdInvalid, "Invalid ID variable. (Varible missing.)")
		log.Debug(err.StackTrace())
		panic(err)
	}

	id, pErr := strconv.Atoi(v)
	if pErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID variable. (Varible must be numeric.)", pErr)
		log.Debug(err.StackTrace())
		panic(err)
	}

	if id <= 0 {
		err := e.NewError(e.ValIdInvalid, "Invalid ID variable. (Varible must be positive.)")
		log.Debug(err.StackTrace())
		panic(err)
	}

	return id
}

func getPathVar(r *http.Request, n string) (string, bool) {
	vars := mux.Vars(r)
	elem, ok := vars[n]
	return elem, ok
}

func getFilterQueryParam(r *http.Request) string {
	return getStringQueryParam(r, "filter")
}

func getSortQueryParam(r *http.Request) string {
	return getStringQueryParam(r, "sort")
}

func getOffsetQueryParam(r *http.Request) int {
	i, err := getIntQueryParam(r, "offset")
	if err != nil {
		err := e.NewError(e.ValOffsetInvalid, "Invalid offset. (Offset must be numeric (int32).)")
		log.Debug(err.StackTrace())
		panic(err)
	}
	if i < 0 {
		err := e.NewError(e.ValOffsetInvalid, "Invalid offset. (Offset must be positive.)")
		log.Debug(err.StackTrace())
		panic(err)
	}
	return i
}

func getLimitQueryParam(r *http.Request) int {
	i, err := getIntQueryParam(r, "limit")
	if err != nil {
		err := e.NewError(e.ValLimitInvalid, "Invalid limit. (Limit must be numeric (int32).)")
		log.Debug(err.StackTrace())
		panic(err)
	}
	if i < 0 {
		err := e.NewError(e.ValLimitInvalid, "Invalid limit. (Limit must be positive.)")
		log.Debug(err.StackTrace())
		panic(err)
	}
	return i
}

func getStringQueryParam(r *http.Request, n string) string {
	qvs := r.URL.Query()
	return qvs.Get(n)
}

func getIntQueryParam(r *http.Request, n string) (int, error) {
	qvs := r.URL.Query()
	qv := qvs.Get(n)
	if qv == "" {
		return 0, nil
	}
	return strconv.Atoi(qv)
}

func parseDate(d string) (time.Time, error) {
	return time.ParseInLocation(constant.ApiDateFormat, d, time.Local)
}

func parseTimestamp(ts string) (time.Time, error) {
	return time.ParseInLocation(constant.ApiTimestampFormat, ts, time.Local)
}
