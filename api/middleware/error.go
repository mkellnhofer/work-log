package middleware

import (
	"net/http"

	"kellnhofer.com/work-log/api/model"
	"kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	httputil "kellnhofer.com/work-log/util/http"
)

var httpStatusCodeMapping = map[int]int{
	error.AuthUnknown:            http.StatusUnauthorized,
	error.AuthInvalidCredentials: http.StatusUnauthorized,

	error.PermUnknown:             http.StatusForbidden,
	error.PermGetUserData:         http.StatusForbidden,
	error.PermChangeUserData:      http.StatusForbidden,
	error.PermGetUserAccount:      http.StatusForbidden,
	error.PermChangeUserAccount:   http.StatusForbidden,
	error.PermGetEntryCharacts:    http.StatusForbidden,
	error.PermChangeEntryCharacts: http.StatusForbidden,
	error.PermGetAllEntries:       http.StatusForbidden,
	error.PermChangeAllEntries:    http.StatusForbidden,
	error.PermGetOwnEntries:       http.StatusForbidden,
	error.PermChangeOwnEntries:    http.StatusForbidden,

	error.ValUnknown:              http.StatusBadRequest,
	error.ValJsonInvalid:          http.StatusBadRequest,
	error.ValPageNumberInvalid:    http.StatusBadRequest,
	error.ValIdInvalid:            http.StatusBadRequest,
	error.ValFilterInvalid:        http.StatusBadRequest,
	error.ValSortInvalid:          http.StatusBadRequest,
	error.ValOffsetInvalid:        http.StatusBadRequest,
	error.ValLimitInvalid:         http.StatusBadRequest,
	error.ValDateInvalid:          http.StatusBadRequest,
	error.ValStartDateInvalid:     http.StatusBadRequest,
	error.ValEndDateInvalid:       http.StatusBadRequest,
	error.ValStartTimeInvalid:     http.StatusBadRequest,
	error.ValEndTimeInvalid:       http.StatusBadRequest,
	error.ValBreakDurationInvalid: http.StatusBadRequest,
	error.ValDescriptionTooLong:   http.StatusBadRequest,
	error.ValSearchInvalid:        http.StatusBadRequest,
	error.ValSearchQueryInvalid:   http.StatusBadRequest,
	error.ValMonthInvalid:         http.StatusBadRequest,

	error.LogicEntryNotFound:                  http.StatusNotFound,
	error.LogicEntryTypeNotFound:              http.StatusNotFound,
	error.LogicEntryActivityNotFound:          http.StatusNotFound,
	error.LogicEntryActivityDeleteNotAllowed:  http.StatusConflict,
	error.LogicEntryTimeIntervalInvalid:       http.StatusBadRequest,
	error.LogicEntryBreakDurationTooLong:      http.StatusBadRequest,
	error.LogicEntrySearchDateIntervalInvalid: http.StatusBadRequest,
	error.LogicRoleNotFound:                   http.StatusNotFound,
	error.LogicUserNotFound:                   http.StatusNotFound,
	error.LogicUserAlreadyExists:              http.StatusConflict,
}

func getHttpStatusCode(errorCode int) int {
	sc, ok := httpStatusCodeMapping[errorCode]
	if !ok {
		log.Debugf("Unexpected error code %d. Using status code 500.", errorCode)
		return http.StatusInternalServerError
	}
	return sc
}

// ErrorMiddleware catches errors and returns a error response.
type ErrorMiddleware struct {
}

// NewErrorMiddleware create a new error middleware.
func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

// ServeHTTP processes requests.
func (m *ErrorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Verb("Before API error check.")

	defer func() {
		if err := recover(); err != nil {
			log.Verb("Catching error.")
			m.handleError(w, r, err)
		}
	}()

	next(w, r)

	log.Verb("After API error check.")
}

func (m *ErrorMiddleware) handleError(w http.ResponseWriter, r *http.Request, err interface{}) {
	switch e := err.(type) {
	case *error.Error:
		// We can retrieve the status here and write out a specific status code.
		sc := getHttpStatusCode(e.Code)
		er := model.NewError(e.Code, e.Message)
		httputil.WriteHttpResponse(w, sc, er)
	default:
		// Any error types we don't specifically look out for default to serving a HTTP 500
		log.Errorf("%s", e)
		sc := http.StatusInternalServerError
		er := model.NewError(error.SysUnknown, "Internal Server Error")
		httputil.WriteHttpResponse(w, sc, er)
	}
}
