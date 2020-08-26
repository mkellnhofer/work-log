package controller

import (
	"context"
	"fmt"
	"net/http"

	"kellnhofer.com/work-log/api/mapper"
	"kellnhofer.com/work-log/api/model"
	"kellnhofer.com/work-log/api/validator"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	m "kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/service"
	httputil "kellnhofer.com/work-log/util/http"
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
	// # Filtering
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
	// # Sorting
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
	//   description: Size of the entries result page. (default=50)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get filter from request
		f, err := getEntriesFilter(getFilterQueryParam(r))
		if err != nil {
			panic(err)
		}

		// Get sort from request
		s, err := getEntriesSort(getSortQueryParam(r))
		if err != nil {
			panic(err)
		}

		// Get offset and limit from request
		o := getOffsetQueryParam(r)
		l := getLimitQueryParam(r)
		if l == 0 {
			l = defaultPageSize
		}

		// Execute action
		entries, cnt, err := c.eServ.GetEntries(r.Context(), f, s, o, l)
		if err != nil {
			panic(err)
		}

		// Convert to API model and write response
		aes := mapper.ToEntries(entries, o, l, cnt)
		httputil.WriteHttpResponse(w, http.StatusOK, aes)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Read API model from request
		var ace model.CreateEntry
		httputil.ReadHttpBody(r, &ace)

		// Validate model
		if err := validator.ValidateCreateEntry(&ace); err != nil {
			panic(err)
		}

		// Convert to logic model
		entry := mapper.FromCreateEntry(&ace)

		// Execute action
		err := c.eServ.CreateEntry(r.Context(), entry)
		if err != nil {
			panic(err)
		}

		// Convert to API model and write response
		ae := mapper.ToEntry(entry)
		httputil.WriteHttpResponse(w, http.StatusOK, ae)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from request
		id := getIdPathVar(r)

		// Execute action
		entry, err := c.eServ.GetEntryById(r.Context(), id)
		if err != nil {
			if err.IsPermissionError() {
				err = c.convertPermissionError(r.Context(), id, err)
			}
			panic(err)
		}

		// Check if a entry was found
		if entry == nil {
			err = e.NewError(e.LogicEntryNotFound, fmt.Sprintf("Could not find entry %d.", id))
			log.Debug(err.StackTrace())
			panic(err)
		}

		// Convert to API model and write response
		ae := mapper.ToEntry(entry)
		httputil.WriteHttpResponse(w, http.StatusOK, ae)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from request
		id := getIdPathVar(r)

		// Read API model from request
		var aue model.UpdateEntry
		httputil.ReadHttpBody(r, &aue)

		// Validate model
		if err := validator.ValidateUpdateEntry(&aue); err != nil {
			panic(err)
		}

		// Convert to logic model
		entry := mapper.FromUpdateEntry(id, &aue)

		// Execute action
		err := c.eServ.UpdateEntry(r.Context(), entry)
		if err != nil {
			if err.IsPermissionError() {
				err = c.convertPermissionError(r.Context(), id, err)
			}
			panic(err)
		}

		// Convert to API model and write response
		ae := mapper.ToEntry(entry)
		httputil.WriteHttpResponse(w, http.StatusOK, ae)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from request
		id := getIdPathVar(r)

		// Execute action
		err := c.eServ.DeleteEntryById(r.Context(), id)
		if err != nil {
			if err.IsPermissionError() {
				err = c.convertPermissionError(r.Context(), id, err)
			}
			panic(err)
		}

		// Write response
		httputil.WriteHttpResponse(w, http.StatusNoContent, nil)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute action
		entryTypes, err := c.eServ.GetEntryTypes(r.Context())
		if err != nil {
			panic(err)
		}

		// Convert to API model and write response
		aets := mapper.ToEntryTypes(entryTypes)
		httputil.WriteHttpResponse(w, http.StatusOK, aets)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute action
		entryActivities, err := c.eServ.GetEntryActivities(r.Context())
		if err != nil {
			panic(err)
		}

		// Convert to API model and write response
		aeas := mapper.ToEntryActivities(entryActivities)
		httputil.WriteHttpResponse(w, http.StatusOK, aeas)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Read API model from request
		var acea model.CreateEntryActivity
		httputil.ReadHttpBody(r, &acea)

		// Validate model
		if err := validator.ValidateCreateEntryActivity(&acea); err != nil {
			panic(err)
		}

		// Convert to logic model
		entryActivity := mapper.FromCreateEntryActivity(&acea)

		// Execute action
		err := c.eServ.CreateEntryActivity(r.Context(), entryActivity)
		if err != nil {
			panic(err)
		}

		// Convert to API model and write response
		aea := mapper.ToEntryActivity(entryActivity)
		httputil.WriteHttpResponse(w, http.StatusOK, aea)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from request
		id := getIdPathVar(r)

		// Read API model from request
		var auea model.UpdateEntryActivity
		httputil.ReadHttpBody(r, &auea)

		// Validate model
		if err := validator.ValidateUpdateEntryActivity(&auea); err != nil {
			panic(err)
		}

		// Convert to logic model
		entryActivity := mapper.FromUpdateEntryActivity(id, &auea)

		// Execute action
		err := c.eServ.UpdateEntryActivity(r.Context(), entryActivity)
		if err != nil {
			panic(err)
		}

		// Convert to API model and write response
		aea := mapper.ToEntryActivity(entryActivity)
		httputil.WriteHttpResponse(w, http.StatusOK, aea)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from request
		id := getIdPathVar(r)

		// Execute action
		err := c.eServ.DeleteEntryActivityById(r.Context(), id)
		if err != nil {
			panic(err)
		}

		// Write response
		httputil.WriteHttpResponse(w, http.StatusNoContent, nil)
	}
}

// --- Permission helper functions ---

func (c *EntryController) convertPermissionError(ctx context.Context, id int, pErr *e.Error) *e.Error {
	if !hasCurrentUserRight(ctx, m.RightGetAllEntries) {
		return e.WrapError(e.LogicEntryNotFound, fmt.Sprintf("Could not find entry %d.", id), pErr)
	}
	return pErr
}
