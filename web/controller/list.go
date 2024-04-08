package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	"kellnhofer.com/work-log/web/pages"
)

// ListController handles requests for list endpoints.
type ListController struct {
	baseController

	mapper *mapper.ListMapper
}

// NewListController creates a new list controller.
func NewListController(uServ *service.UserService, eServ *service.EntryService) *ListController {
	return &ListController{
		baseController: baseController{
			uServ: uServ,
			eServ: eServ,
		},
		mapper: mapper.NewListMapper(),
	}
}

// --- Endpoints ---

// GetListHandler returns a handler for "GET /list".
func (c *ListController) GetListHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /list.")
		return c.handleShowList(eCtx)
	}
}

// --- Handler functions ---

func (c *ListController) handleShowList(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID and user contract
	userId, userContract, err := c.getUserIdAndUserContract(ctx)
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
		workSummary, err = c.eServ.GetTotalWorkSummaryByUserId(ctx, userId)
		if err != nil {
			return err
		}
	}

	// Get entries
	entries, cnt, err := c.eServ.GetDateEntriesByUserId(ctx, userId, offset, limit)
	if err != nil {
		return err
	}
	// Get entry master data
	entryTypesMap, entryActivitiesMap, err := c.getEntryMasterDataMap(ctx)
	if err != nil {
		return err
	}

	// Create view model
	model := c.mapper.CreateListViewModel(userContract, workSummary, pageNum, pageSize, cnt, entries,
		entryTypesMap, entryActivitiesMap)

	// Save current URL to be able to used later for back navigation
	saveCurrentUrl(eCtx)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.ListEntriesPage(model))
}
