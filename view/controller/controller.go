package controller

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/view/middleware"
)

func getStringPathVar(r *http.Request, n string) (string, bool) {
	vs := mux.Vars(r)
	v, ok := vs[n]
	return v, ok
}

func getIdPathVar(r *http.Request) int {
	v, ok := getStringPathVar(r, "id")
	if !ok {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (Variable missing.)")
		log.Debug(err.StackTrace())
		panic(err)
	}

	return parseIdParam(v)
}

func parseIdParam(v string) int {
	id, pErr := strconv.Atoi(v)
	if pErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID. (Variable must be numeric.)", pErr)
		log.Debug(err.StackTrace())
		panic(err)
	}

	if id <= 0 {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (Variable must be positive.)")
		log.Debug(err.StackTrace())
		panic(err)
	}

	return id
}

func getStringQueryParam(r *http.Request, n string) (string, bool) {
	qvs := r.URL.Query()
	qv := qvs.Get(n)
	if qv == "" {
		return qv, false
	}
	return qv, true
}

func getErrorCodeQueryParam(r *http.Request) *int {
	v, ok := getStringQueryParam(r, "error")
	if !ok {
		return nil
	}

	return parseErrorCodeParam(v)
}

func parseErrorCodeParam(v string) *int {
	ec, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	return &ec
}

func getPageNumberQueryParam(r *http.Request) *int {
	v, ok := getStringQueryParam(r, "page")
	if !ok {
		return nil
	}

	return parsePageNumberParam(v)
}

func parsePageNumberParam(v string) *int {
	page, err := strconv.Atoi(v)
	if err != nil {
		err := e.WrapError(e.ValPageNumberInvalid, "Invalid page number. (Variable must be numeric.)",
			err)
		log.Debug(err.StackTrace())
		panic(err)
	}

	if page <= 0 {
		err := e.NewError(e.ValPageNumberInvalid, "Invalid page number. (Variable must be positive.)")
		log.Debug(err.StackTrace())
		panic(err)
	}

	return &page
}

func getSearchQueryParam(r *http.Request) *string {
	v, ok := getStringQueryParam(r, "query")
	if !ok {
		return nil
	}

	return &v
}

func getMonthQueryParam(r *http.Request) (*int, *int) {
	v, ok := getStringQueryParam(r, "month")
	if !ok {
		return nil, nil
	}

	return parseMonthParam(v)
}

func parseMonthParam(v string) (*int, *int) {
	if len(v) != 6 {
		err := e.NewError(e.ValMonthInvalid, "Invalid month. (Variable must have length of 6.)")
		log.Debug(err.StackTrace())
		panic(err)
	}

	ys := v[0:4]
	ms := v[4:6]

	y, err := strconv.Atoi(ys)
	if err != nil {
		err := e.NewError(e.ValMonthInvalid, "Invalid month. (Year part is invalid.)")
		log.Debug(err.StackTrace())
		panic(err)
	}
	m, err := strconv.Atoi(ms)
	if err != nil || m <= 0 || m > 12 {
		err := e.NewError(e.ValMonthInvalid, "Invalid month. (Month part invalid.)")
		log.Debug(err.StackTrace())
		panic(err)
	}

	return &y, &m
}

func getCurrentUserId(ctx context.Context) int {
	sc := ctx.Value(constant.ContextKeySecurityContext).(*model.SecurityContext)
	return sc.UserId
}

func saveCurrentUrl(ctx context.Context, r *http.Request) {
	sh := ctx.Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)
	s := sh.Get()
	path := r.URL.EscapedPath()
	query := r.URL.RawQuery
	req := path
	if query != "" {
		req = req + "?" + query
	}
	s.PreviousUrl = req
}

func getPreviousUrl(ctx context.Context) string {
	sh := ctx.Value(constant.ContextKeySessionHolder).(*middleware.SessionHolder)
	s := sh.Get()
	if s.PreviousUrl != "" {
		return s.PreviousUrl
	} else {
		return constant.ViewPathDefault
	}
}

func writeFile(w http.ResponseWriter, fileName string, wt io.WriterTo) {
	// Write header
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")

	// Write body
	_, wErr := wt.WriteTo(w)
	if wErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not write response.", wErr)
		log.Debug(err.StackTrace())
		panic(err)
	}
}
