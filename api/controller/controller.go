package controller

import "kellnhofer.com/work-log/api/model"

// Information about the error.
// swagger:response ErrorResponse
type Error struct {
	// in: body
	Body model.Error
}
