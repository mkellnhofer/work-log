package controller

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/api/export"
	"kellnhofer.com/work-log/pkg/constant"
	"kellnhofer.com/work-log/pkg/service"
)

// ExportController handles requests for export endpoints.
type ExportController struct {
	eServ *service.EntryService

	exporter *export.EntriesExporter
}

// NewExportController creates a new export controller.
func NewExportController(eServ *service.EntryService) *ExportController {
	return &ExportController{
		eServ:    eServ,
		exporter: export.NewEntriesExporter(),
	}
}

// GetExportHandler returns a handler for "GET /export".
func (c *ExportController) GetExportHandler() echo.HandlerFunc {
	// swagger:operation GET /export export exportEntries
	//
	// Export entries as CSV.
	//
	// Only entries a user can see are exported.
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
	// | project     | i (is), cn (contains) | null, string               |
	// | description | i (is), cn (contains) | null, string               |
	// | labels      | i (is), in (in)       | null, strings              |
	// &#9432; Filters are connected via logical conjunction (AND).
	//
	// __Filter Syntax:__
	// [field name];[operator];[value-1];...;[value-n]
	//
	// __Examples:__
	//
	// Get entries for a specific time interval: startTime;bt;2019-01-01T00:00:00;2019-01-05T00:00:00
	//
	// Get entries with specific labels (OR logic): labels;in;bug;frontend
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
	// - text/csv
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
	//
	// responses:
	//   '200':
	//     description: CSV file containing the exported entries.
	//     schema:
	//       type: string
	//       format: binary
	//   '400':
	//     description: "__Bad Request__\n\n
	//       ⦁ [-304]: Invalid filter\n
	//       ⦁ [-305]: Invalid sort"
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

		// Get all entries (no pagination for export)
		entries, _, err := c.eServ.GetEntries(getContext(eCtx), f, s, 0, 0)
		if err != nil {
			return err
		}

		// Get entry types and entry activities for lookup
		entryTypes, err := c.eServ.GetEntryTypes(getContext(eCtx))
		if err != nil {
			return err
		}
		entryActivities, err := c.eServ.GetEntryActivities(getContext(eCtx))
		if err != nil {
			return err
		}

		// Create file name
		timestamp := time.Now().Format(constant.ExportTimestampFormat)
		fileName := fmt.Sprintf(constant.ExportFileNameTemplate, timestamp, "csv")
		// Create CSV export
		file := c.exporter.ExportEntries(entries, entryTypes, entryActivities)

		// Write file response
		return writeFileResponse(eCtx, fileName, file)
	}
}
