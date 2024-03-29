package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/api/mapper"
	"kellnhofer.com/work-log/api/model"
	"kellnhofer.com/work-log/api/validator"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	m "kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
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
func (c *EntryController) GetEntriesHandler() echo.HandlerFunc {
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
	//   '400':
	//     description: "__Bad Request__\n\n
	//       ⦁ [-302]: Invalid page number\n
	//       ⦁ [-304]: Invalid filter\n
	//       ⦁ [-305]: Invalid sort\n
	//       ⦁ [-306]: Invalid offset\n
	//       ⦁ [-307]: Invalid limit"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-207]: No right to get entries of other users"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get filter from request
		f, err := getEntriesFilter(getFilterQueryParam(eCtx))
		if err != nil {
			return err
		}

		// Get sort from request
		s, err := getEntriesSort(getSortQueryParam(eCtx))
		if err != nil {
			return err
		}

		// Get offset and limit from request
		var o, l int
		if o, err = getOffsetQueryParam(eCtx); err != nil {
			return err
		}
		if l, err = getLimitQueryParam(eCtx); err != nil {
			return err
		}
		if l == 0 {
			l = defaultPageSize
		}

		// Execute action
		entries, cnt, err := c.eServ.GetEntries(getContext(eCtx), f, s, o, l)
		if err != nil {
			return err
		}

		// Convert to API model and write response
		aes := mapper.ToEntries(entries, o, l, cnt)
		return writeResponse(eCtx, http.StatusOK, aes)
	}
}

// CreateEntryHandler returns a handler for "POST /entries".
func (c *EntryController) CreateEntryHandler() echo.HandlerFunc {
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
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID\n
	//       ⦁ [-309]: Negative number\n
	//       ⦁ [-312]: Too long string\n
	//       ⦁ [-314]: Invalid timestamp\n
	//       ⦁ [-405]: Invalid time interval"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-208]: No right to create entries for other users"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-402]: Entry type not found\n
	//       ⦁ [-403]: Entry activity not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Read API model from request
		var ace model.CreateEntry
		if err := readRequestBody(eCtx, &ace); err != nil {
			return err
		}

		// Validate model
		if err := validator.ValidateCreateEntry(&ace); err != nil {
			return err
		}

		// Convert to logic model
		entry := mapper.FromCreateEntry(&ace)

		// Execute action
		if err := c.eServ.CreateEntry(getContext(eCtx), entry); err != nil {
			return err
		}

		// Convert to API model and write response
		ae := mapper.ToEntry(entry)
		return writeResponse(eCtx, http.StatusOK, ae)
	}
}

// GetEntryHandler returns a handler for "GET /entries/{id}".
func (c *EntryController) GetEntryHandler() echo.HandlerFunc {
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
	//   '400':
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-401]: Entry not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get ID from request
		id, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Execute action
		entry, err := c.eServ.GetEntryById(getContext(eCtx), id)
		if err != nil {
			return c.convertPermissionError(getContext(eCtx), id, err)
		}

		// Check if a entry was found
		if entry == nil {
			err := e.NewError(e.LogicEntryNotFound, fmt.Sprintf("Could not find entry %d.", id))
			log.Debug(err.StackTrace())
			return err
		}

		// Convert to API model and write response
		ae := mapper.ToEntry(entry)
		return writeResponse(eCtx, http.StatusOK, ae)
	}
}

// UpdateEntryHandler returns a handler for "PUT /entries/{id}".
func (c *EntryController) UpdateEntryHandler() echo.HandlerFunc {
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
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID\n
	//       ⦁ [-309]: Negative number\n
	//       ⦁ [-312]: Too long string\n
	//       ⦁ [-314]: Invalid timestamp\n
	//       ⦁ [-405]: Invalid time interval"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-208]: No right to update entries of other users"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-401]: Entry not found\n
	//       ⦁ [-402]: Entry type not found\n
	//       ⦁ [-403]: Entry activity not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get ID from request
		id, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Read API model from request
		var aue model.UpdateEntry
		if err := readRequestBody(eCtx, &aue); err != nil {
			return err
		}

		// Validate model
		if err := validator.ValidateUpdateEntry(&aue); err != nil {
			return err
		}

		// Convert to logic model
		entry := mapper.FromUpdateEntry(id, &aue)

		// Execute action
		if err := c.eServ.UpdateEntry(getContext(eCtx), entry); err != nil {
			return c.convertPermissionError(getContext(eCtx), id, err)
		}

		// Convert to API model and write response
		ae := mapper.ToEntry(entry)
		return writeResponse(eCtx, http.StatusOK, ae)
	}
}

// DeleteEntryHandler returns a handler for "DELETE /entries/{id}".
func (c *EntryController) DeleteEntryHandler() echo.HandlerFunc {
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
	//   '400':
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-208]: No right to delete entries of other users"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-401]: Entry not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get ID from request
		id, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Execute action
		if err := c.eServ.DeleteEntryById(getContext(eCtx), id); err != nil {
			return c.convertPermissionError(getContext(eCtx), id, err)
		}

		// Write response
		return writeResponse(eCtx, http.StatusNoContent, nil)
	}
}

// GetEntryTypesHandler returns a handler for "GET /entry_types".
func (c *EntryController) GetEntryTypesHandler() echo.HandlerFunc {
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
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Execute action
		entryTypes, err := c.eServ.GetEntryTypes(getContext(eCtx))
		if err != nil {
			return err
		}

		// Convert to API model and write response
		aets := mapper.ToEntryTypes(entryTypes)
		return writeResponse(eCtx, http.StatusOK, aets)
	}
}

// GetEntryActivitiesHandler returns a handler for "GET /entry_activites".
func (c *EntryController) GetEntryActivitiesHandler() echo.HandlerFunc {
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
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Execute action
		entryActivities, err := c.eServ.GetEntryActivities(getContext(eCtx))
		if err != nil {
			return err
		}

		// Convert to API model and write response
		aeas := mapper.ToEntryActivities(entryActivities)
		return writeResponse(eCtx, http.StatusOK, aeas)
	}
}

// CreateEntryActivityHandler returns a handler for "POST /entry_activites".
func (c *EntryController) CreateEntryActivityHandler() echo.HandlerFunc {
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
	//     description: "__Bad Request__\n\n
	//       ⦁ [-311]: Empty string\n
	//       ⦁ [-312]: Too long string"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-206]: No right to create entry activities"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Read API model from request
		var acea model.CreateEntryActivity
		if err := readRequestBody(eCtx, &acea); err != nil {
			return err
		}

		// Validate model
		if err := validator.ValidateCreateEntryActivity(&acea); err != nil {
			return err
		}

		// Convert to logic model
		entryActivity := mapper.FromCreateEntryActivity(&acea)

		// Execute action
		if err := c.eServ.CreateEntryActivity(getContext(eCtx), entryActivity); err != nil {
			return err
		}

		// Convert to API model and write response
		aea := mapper.ToEntryActivity(entryActivity)
		return writeResponse(eCtx, http.StatusOK, aea)
	}
}

// UpdateEntryActivityHandler returns a handler for "PUT /entry_activites/{id}".
func (c *EntryController) UpdateEntryActivityHandler() echo.HandlerFunc {
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
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID\n
	//       ⦁ [-311]: Empty string\n
	//       ⦁ [-312]: Too long string"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-206]: No right to update entry activities"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-403]: Entry activity not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get ID from request
		id, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Read API model from request
		var auea model.UpdateEntryActivity
		if err := readRequestBody(eCtx, &auea); err != nil {
			return err
		}

		// Validate model
		if err := validator.ValidateUpdateEntryActivity(&auea); err != nil {
			return err
		}

		// Convert to logic model
		entryActivity := mapper.FromUpdateEntryActivity(id, &auea)

		// Execute action
		if err := c.eServ.UpdateEntryActivity(getContext(eCtx), entryActivity); err != nil {
			return err
		}

		// Convert to API model and write response
		aea := mapper.ToEntryActivity(entryActivity)
		return writeResponse(eCtx, http.StatusOK, aea)
	}
}

// DeleteEntryActivityHandler returns a handler for "DELETE /entry_activites/{id}".
func (c *EntryController) DeleteEntryActivityHandler() echo.HandlerFunc {
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
	//   '400':
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-206]: No right to delete entry activities"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-403]: Entry activity not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '409':
	//     description: "__Conflict__\n\n
	//       ⦁ [-404]: Entry activity still used"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get ID from request
		id, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Execute action
		if err := c.eServ.DeleteEntryActivityById(getContext(eCtx), id); err != nil {
			return err
		}

		// Write response
		return writeResponse(eCtx, http.StatusNoContent, nil)
	}
}

// --- Permission helper functions ---

func (c *EntryController) convertPermissionError(ctx context.Context, id int, err error) error {
	er, ok := err.(*e.Error)
	if ok && er.IsPermissionError() && !hasCurrentUserRight(ctx, m.RightGetAllEntries) {
		return e.WrapError(e.LogicEntryNotFound, fmt.Sprintf("Could not find entry %d.", id), err)
	}
	return err
}
