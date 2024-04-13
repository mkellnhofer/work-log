package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	"kellnhofer.com/work-log/web/view/pages"
)

// LogController handles requests for log endpoints.
type LogController struct {
	baseController

	mapper *mapper.LogMapper
}

// NewLogController creates a new log controller.
func NewLogController(uServ *service.UserService, eServ *service.EntryService) *LogController {
	return &LogController{
		baseController: baseController{
			uServ: uServ,
			eServ: eServ,
		},
		mapper: mapper.NewLogMapper(),
	}
}

// --- Endpoints ---

// GetLogHandler returns a handler for "GET /log".
func (c *LogController) GetLogHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /log.")
		return c.handleShowLog(eCtx)
	}
}

// --- Handler functions ---

func (c *LogController) handleShowLog(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user and user contract
	user, userContract, err := c.getUserAndUserContract(ctx)
	if err != nil {
		return err
	}

	// Get page number, offset and limit
	pageNum, _, err := getPageNumberQueryParam(eCtx)
	if err != nil {
		return err
	}
	offset, limit := calculateOffsetLimitFromPageNumber(pageNum)

	// Get work summary (only for first page)
	var workSummary *model.WorkSummary
	if pageNum == 1 {
		workSummary, err = c.eServ.GetTotalWorkSummaryByUserId(ctx, user.Id)
		if err != nil {
			return err
		}
	}

	// Get entries
	entries, cnt, err := c.eServ.GetDateEntriesByUserId(ctx, user.Id, offset, limit)
	if err != nil {
		return err
	}
	// Get entry master data
	entryTypesMap, entryActivitiesMap, err := c.getEntryMasterDataMap(ctx)
	if err != nil {
		return err
	}

	// Create view model
	userModel := c.mapper.CreateUserInfoViewModel(user)
	model := c.mapper.CreateLogViewModel(userContract, workSummary, pageNum, pageSize, cnt, entries,
		entryTypesMap, entryActivitiesMap)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.LogEntriesPage(userModel, model))
}
