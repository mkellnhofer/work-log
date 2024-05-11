package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/api/model"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	httputil "kellnhofer.com/work-log/pkg/util/http"
)

var httpStatusCodeMapping = map[int]int{
	e.AuthUnknown:            http.StatusUnauthorized,
	e.AuthDataInvalid:        http.StatusUnauthorized,
	e.AuthCredentialsInvalid: http.StatusUnauthorized,
	e.AuthUserNotActivated:   http.StatusPreconditionFailed,

	e.PermUnknown:             http.StatusForbidden,
	e.PermGetUserData:         http.StatusForbidden,
	e.PermChangeUserData:      http.StatusForbidden,
	e.PermGetUserAccount:      http.StatusForbidden,
	e.PermChangeUserAccount:   http.StatusForbidden,
	e.PermGetEntryCharacts:    http.StatusForbidden,
	e.PermChangeEntryCharacts: http.StatusForbidden,
	e.PermGetAllEntries:       http.StatusForbidden,
	e.PermChangeAllEntries:    http.StatusForbidden,
	e.PermGetOwnEntries:       http.StatusForbidden,
	e.PermChangeOwnEntries:    http.StatusForbidden,

	e.ValUnknown:              http.StatusBadRequest,
	e.ValJsonInvalid:          http.StatusBadRequest,
	e.ValPageNumberInvalid:    http.StatusBadRequest,
	e.ValIdInvalid:            http.StatusBadRequest,
	e.ValFilterInvalid:        http.StatusBadRequest,
	e.ValSortInvalid:          http.StatusBadRequest,
	e.ValOffsetInvalid:        http.StatusBadRequest,
	e.ValLimitInvalid:         http.StatusBadRequest,
	e.ValFieldNil:             http.StatusBadRequest,
	e.ValNumberNegative:       http.StatusBadRequest,
	e.ValNumberNegativeOrZero: http.StatusBadRequest,
	e.ValStringEmpty:          http.StatusBadRequest,
	e.ValStringTooLong:        http.StatusBadRequest,
	e.ValDateInvalid:          http.StatusBadRequest,
	e.ValTimestampInvalid:     http.StatusBadRequest,
	e.ValArrayEmpty:           http.StatusBadRequest,
	e.ValRoleInvalid:          http.StatusBadRequest,
	e.ValUsernameInvalid:      http.StatusBadRequest,
	e.ValPasswordInvalid:      http.StatusBadRequest,

	e.LogicEntryNotFound:                  http.StatusNotFound,
	e.LogicEntryTypeNotFound:              http.StatusNotFound,
	e.LogicEntryActivityNotFound:          http.StatusNotFound,
	e.LogicEntryActivityDeleteNotAllowed:  http.StatusConflict,
	e.LogicEntryTimeIntervalInvalid:       http.StatusBadRequest,
	e.LogicEntrySearchDateIntervalInvalid: http.StatusBadRequest,
	e.LogicRoleNotFound:                   http.StatusNotFound,
	e.LogicUserNotFound:                   http.StatusNotFound,
	e.LogicUserAlreadyExists:              http.StatusConflict,
	e.LogicContractWorkingHoursInvalid:    http.StatusBadRequest,
	e.LogicContractVacationDaysInvalid:    http.StatusBadRequest,
	e.LogicEntryActivityNotAllowed:        http.StatusBadRequest,
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

// CreateHandler creates a new handler to process requests.
func (m *ErrorMiddleware) CreateHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Verb("Before API error check.")

		err := m.process(next, c)

		log.Verb("After API error check.")

		return err
	}
}

func (m *ErrorMiddleware) process(next echo.HandlerFunc, c echo.Context) error {
	err := next(c)
	if err != nil {
		log.Verb("Catching error.")
		m.handleError(c.Response(), err)
	}
	return nil
}

func (m *ErrorMiddleware) handleError(r *echo.Response, err error) {
	switch tErr := err.(type) {
	case *e.Error:
		// We can retrieve the status here and write out a specific status code.
		sc := getHttpStatusCode(tErr.Code)
		er := model.NewError(tErr.Code, tErr.Message)
		httputil.WriteHttpResponse(r, sc, er)
	default:
		// Any error types we don't specifically look out for default to serving a HTTP 500
		log.Errorf("%s", tErr)
		sc := http.StatusInternalServerError
		er := model.NewError(e.SysUnknown, "Internal Server Error")
		httputil.WriteHttpResponse(r, sc, er)
	}
}
