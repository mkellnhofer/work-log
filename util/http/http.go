package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
)

// ReadHttpBody reads the HTTP request
func ReadHttpBody(r *http.Request, data interface{}) {
	readRequest(r, data)
}

// WriteHttpResponse writes the HTTP response.
func WriteHttpResponse(writer http.ResponseWriter, statusCode int, data interface{}) {
	writeResponse(writer, statusCode, data)
}

// WriteHttpError writes the HTTP error response.
func WriteHttpError(writer http.ResponseWriter, statusCode int, error string) {
	writeResponse(writer, statusCode, error)
}

func readRequest(r *http.Request, data interface{}) {
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
		panic(err)
	}
}

func writeResponse(writer http.ResponseWriter, statusCode int, data interface{}) {
	var body []byte
	var jErr error
	if data != nil {
		body, jErr = json.Marshal(data)
		if jErr != nil {
			err := e.WrapError(e.SysUnknown, "Could not encode JSON.", jErr)
			log.Error(err.StackTrace())
			panic(err)
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
}
