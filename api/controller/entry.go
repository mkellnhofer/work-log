package controller

import (
	"net/http"

	"kellnhofer.com/work-log/api/model"
	"kellnhofer.com/work-log/service"
)

// EntryController handles requests for entry endpoints.
type EntryController struct {
	eServ *service.EntryService
}

// NewEntryController create a new enty controller.
func NewEntryController(es *service.EntryService) *EntryController {
	return &EntryController{es}
}

// --- Parameters ---

// swagger:parameters getEntry
type GetEntryParameters struct {
	// The ID of the entry.
	//
	// in: path
	// required: true
	Id int `json:"id"`
}

// swagger:parameters createEntry
type CreateEntryParameters struct {
	// in: body
	// required: true
	Body model.CreateEntry
}

// swagger:parameters updateEntry
type UpdateEntryParameters struct {
	// The ID of the entry.
	//
	// in: path
	// required: true
	Id int `json:"id"`
	// in: body
	// required: true
	Body model.UpdateEntry
}

// swagger:parameters deleteEntry
type DeleteEntryParameters struct {
	// The ID of the entry.
	//
	// in: path
	// required: true
	Id int `json:"id"`
}

// swagger:parameters createEntryActivity
type CreateEntryActivityParameters struct {
	// in: body
	// required: true
	Body model.CreateEntryActivity
}

// swagger:parameters updateEntryActivity
type UpdateEntryActivityParameters struct {
	// The ID of the entry activity.
	//
	// in: path
	// required: true
	Id int `json:"id"`
	// in: body
	// required: true
	Body model.UpdateEntryActivity
}

// swagger:parameters deleteEntryActivity
type DeleteEntryActivityParameters struct {
	// The ID of the entry activity.
	//
	// in: path
	// required: true
	Id int `json:"id"`
}

// --- Responses ---

// The list of entries.
// swagger:response GetEntriesResponse
type GetEntriesResponse struct {
	// in: body
	Body model.EntryList
}

// The entry.
// swagger:response GetEntryResponse
type GetEntryResponse struct {
	// in: body
	Body model.Entry
}

// The created entry.
// swagger:response CreateEntryResponse
type CreateEntryResponse struct {
	// in: body
	Body model.Entry
}

// The updated entry.
// swagger:response UpdateEntryResponse
type UpdateEntryResponse struct {
	// in: body
	Body model.Entry
}

// The list of entry types.
// swagger:response GetEntryTypesResponse
type GetEntryTypesResponse struct {
	// in: body
	Body []model.EntryType
}

// The list of entry activities.
// swagger:response GetEntryActivitiesResponse
type GetEntryActivitiesResponse struct {
	// in: body
	Body []model.EntryActivity
}

// The created entry activity.
// swagger:response CreateEntryActivityResponse
type CreateEntryActivityResponse struct {
	// in: body
	Body model.EntryActivity
}

// The updated entry activity.
// swagger:response UpdateEntryActivityResponse
type UpdateEntryActivityResponse struct {
	// in: body
	Body model.EntryActivity
}

// --- Endpoints ---

// GetEntriesHandler returns a handler for "GET /entries".
func (c *EntryController) GetEntriesHandler() http.HandlerFunc {
	// swagger:operation GET /entries entries listEntries
	//
	// Lists all entries.
	//
	// Only entries a user can see are returned.
	//
	// # Filtering:
	//
	// The result can be filtered via following fields:
	// | field name  | operators             | data type / allowed values |
	// | ----------- | --------------------- | -------------------------- |
	// | userId      | eq (equal)            | int                        |
	// | typeId      | eq (equal)            | int                        |
	// | startTime   | bt (between)          | datetime strings           |
	// | activityId  | i (is), eq (equal)    | null, int                  |
	// | description | i (is), cn (contains) | null, string               |
	// &#9432; Filters are connected via logical conjunction (AND).
	//
	// __Filter Syntax:__
	// [field name];[operator];[value-1];...;[value-n]
	//
	// __Example:__
	// Get entries for a specific time interval: startTime;bt;2019-01-01T00:00:00;2019-01-05T00:00:00
	//
	// # Sorting:
	//
	// The result can be sorted via following fields:
	// | field name | operators  |
	// | ---------- | ---------- |
	// | startTime  | asc / desc |
	// &#9432; Sorting by multiple fields is not supported.
	//
	// __Sort Syntax:__
	// [field name];[operator]
	//
	// __Example:__
	// Sort entries descending by their date: date;desc
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// parameters:
	// - name: filter
	//   in: query
	//   description: Filtering applied to the entries result.
	//   required: false
	//   type: string
	// - name: sort
	//   in: query
	//   description: Sorting applied to the entries result.
	//   required: false
	//   type: string
	// - name: offset
	//   in: query
	//   description: Start of the entries result page.
	//   required: false
	//   type: integer
	//   format: int32
	// - name: limit
	//   in: query
	//   description: Size of the entries result page. (default=500)
	//   required: false
	//   type: integer
	//   format: int32
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetEntriesResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// CreateEntryHandler returns a handler for "POST /entries".
func (c *EntryController) CreateEntryHandler() http.HandlerFunc {
	// swagger:operation POST /entries entries createEntry
	//
	// Create a entry.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/CreateEntryResponse"
	//   '400':
	//     "$ref": "#/responses/ErrorResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetEntryHandler returns a handler for "GET /entries/{id}".
func (c *EntryController) GetEntryHandler() http.HandlerFunc {
	// swagger:operation GET /entries/{id} entries getEntry
	//
	// Get a entry by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetEntryResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// UpdateEntryHandler returns a handler for "PUT /entries/{id}".
func (c *EntryController) UpdateEntryHandler() http.HandlerFunc {
	// swagger:operation PUT /entries/{id} entries updateEntry
	//
	// Update a entry by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/UpdateEntryResponse"
	//   '400':
	//     "$ref": "#/responses/ErrorResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   '409':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// DeleteEntryHandler returns a handler for "DELETE /entries/{id}".
func (c *EntryController) DeleteEntryHandler() http.HandlerFunc {
	// swagger:operation DELETE /entries/{id} entries deleteEntry
	//
	// Delete a entry by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '204':
	//     description: No content.
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   '409':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetEntryTypesHandler returns a handler for "GET /entry_types".
func (c *EntryController) GetEntryTypesHandler() http.HandlerFunc {
	// swagger:operation GET /entry_types entry_types listEntryTypes
	//
	// Lists all entry types.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetEntryTypesResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetEntryActivitiesHandler returns a handler for "GET /entry_activites".
func (c *EntryController) GetEntryActivitiesHandler() http.HandlerFunc {
	// swagger:operation GET /entry_activities entry_activities listEntryActivities
	//
	// Lists all entry activities.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetEntryActivitiesResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// CreateEntryActivityHandler returns a handler for "POST /entry_activites".
func (c *EntryController) CreateEntryActivityHandler() http.HandlerFunc {
	// swagger:operation POST /entry_activities entry_activities createEntryActivity
	//
	// Create a entry activity.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/CreateEntryActivityResponse"
	//   '400':
	//     "$ref": "#/responses/ErrorResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '409':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// UpdateEntryActivityHandler returns a handler for "PUT /entry_activites/{id}".
func (c *EntryController) UpdateEntryActivityHandler() http.HandlerFunc {
	// swagger:operation PUT /entry_activities/{id} entry_activities updateEntryActivity
	//
	// Update a entry activity by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/UpdateEntryActivityResponse"
	//   '400':
	//     "$ref": "#/responses/ErrorResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   '409':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// DeleteEntryActivityHandler returns a handler for "DELETE /entry_activites/{id}".
func (c *EntryController) DeleteEntryActivityHandler() http.HandlerFunc {
	// swagger:operation DELETE /entry_activities/{id} entry_activities deleteEntryActivity
	//
	// Delete a entry activity by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '204':
	//     description: No content.
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
