package controller

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/view/hx"
	"kellnhofer.com/work-log/web/view/pages"
)

// LogController handles requests for log endpoints.
type LogController struct {
	baseController

	mapper *mapper.LogMapper
}

// NewLogController creates a new log controller.
func NewLogController(uServ *service.UserService, eServ *service.EntryService) *LogController {
	logMapper := mapper.NewLogMapper()
	return &LogController{
		baseController: baseController{
			uServ:  uServ,
			eServ:  eServ,
			mapper: &logMapper.Mapper,
		},
		mapper: logMapper,
	}
}

// --- Endpoints ---

// GetLogHandler returns a handler for "GET /log".
func (c *LogController) GetLogHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle: GET /log")

		isHtmxReq := web.IsHtmxRequest(eCtx)
		pageNum, isPageReq, err := getPageNumberQueryParam(eCtx)
		if err != nil {
			return err
		}

		ctx := getContext(eCtx)

		if !isHtmxReq {
			return c.handleShowLog(eCtx, ctx, pageNum)
		} else if !isPageReq {
			return c.handleHxShowLog(eCtx, ctx, pageNum)
		} else {
			return c.handleHxGetLogPage(eCtx, ctx, pageNum)
		}
	}
}

// --- Handler functions ---

func (c *LogController) handleShowLog(eCtx echo.Context, ctx context.Context, pageNum int) error {
	// Get view data
	userModel, err := c.getUserInfoViewData(ctx)
	if err != nil {
		return err
	}
	model, err := c.getLogViewData(ctx, pageNum)
	if err != nil {
		return err
	}

	// Render
	return web.Render(eCtx, http.StatusOK, pages.Log(userModel, model))
}

func (c *LogController) handleHxShowLog(eCtx echo.Context, ctx context.Context, pageNum int) error {
	// Get view data
	model, err := c.getLogViewData(ctx, pageNum)
	if err != nil {
		return err
	}

	// Render
	return web.Render(eCtx, http.StatusOK, hx.Log(model))
}

func (c *LogController) handleHxGetLogPage(eCtx echo.Context, ctx context.Context, pageNum int) error {
	// TODO!!!
	return nil
}

func (c *LogController) getLogViewData(ctx context.Context, pageNum int) (*vm.LogEntries, error) {
	// Get current user information
	userId := getCurrentUserId(ctx)
	userContract, err := c.getUserContract(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Get work summary data
	workSummary, err := c.eServ.GetTotalWorkSummaryByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Get entries
	offset, limit := calculateOffsetLimitFromPageNumber(pageNum)
	entries, cnt, err := c.eServ.GetDateEntriesByUserId(ctx, userId, offset, limit)
	if err != nil {
		return nil, err
	}
	// Get entry master data
	entryTypesMap, entryActivitiesMap, err := c.getEntryMasterDataMap(ctx)
	if err != nil {
		return nil, err
	}

	// Create view model
	return c.mapper.CreateLogViewModel(userContract, workSummary, pageNum, pageSize, cnt, entries,
		entryTypesMap, entryActivitiesMap), nil
}
