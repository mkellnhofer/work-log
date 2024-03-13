package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
)

// ReadHttpBody reads the HTTP request
func ReadHttpRequestBody(r *http.Request, data interface{}) error {
	return readRequestBody(r, data)
}

// WriteHttpResponse writes the HTTP response.
func WriteHttpResponse(r *echo.Response, statusCode int, data interface{}) error {
	return writeResponse(r.Writer, statusCode, data)
}

// WriteHttpError writes the HTTP error response.
func WriteHttpError(r *echo.Response, statusCode int, error string) error {
	return writeResponse(r.Writer, statusCode, error)
}

func readRequestBody(r *http.Request, data interface{}) error {
	decoder := json.NewDecoder(r.Body)
	jErr := decoder.Decode(data)
	if jErr != nil {
		var err *e.Error
		switch cErr := jErr.(type) {
		case *json.SyntaxError:
			err = e.WrapError(e.ValJsonInvalid, "Could not decode JSON. (Invalid JSON syntax.)",
				jErr)
		case *json.UnmarshalTypeError:
			err = e.WrapError(e.ValJsonInvalid, fmt.Sprintf("Could not decode JSON. (Invalid type "+
				"for field '%s'.)", cErr.Field), jErr)
		default:
			err = e.WrapError(e.ValJsonInvalid, "Could not decode JSON.", jErr)
		}
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func writeResponse(writer http.ResponseWriter, statusCode int, data interface{}) error {
	var body []byte
	var jErr error
	if data != nil {
		body, jErr = json.Marshal(data)
		if jErr != nil {
			err := e.WrapError(e.SysUnknown, "Could not encode JSON.", jErr)
			log.Error(err.StackTrace())
			return err
		}
	}

	writer.WriteHeader(statusCode)
	if body != nil {
		writer.Header().Set("Content-Type", "application/json")
		_, wErr := writer.Write(body)
		if wErr != nil {
			log.Error("Could not write response body!")
			log.Error(wErr.Error())
		}
	}

	return nil
}
