package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/view/hx"
	"kellnhofer.com/work-log/web/view/page"
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
	return c.handler(func(eCtx echo.Context, ctx context.Context, isHtmxReq bool) error {
		pageNum, err := c.getGetLogParams(eCtx)
		if err != nil {
			return err
		}

		if !isHtmxReq {
			return c.handleShowLog(eCtx, ctx, pageNum)
		} else {
			return c.handleNavLog(eCtx, pageNum)
		}
	})
}

// GetLogContentHandler returns a handler for "GET /log/content".
func (c *LogController) GetLogContentHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		pageNum, err := c.getGetLogParams(eCtx)
		if err != nil {
			return err
		}
		return c.handleGetLogContent(eCtx, ctx, pageNum)
	})
}

// --- Handler functions ---

func (c *LogController) handleShowLog(eCtx echo.Context, ctx context.Context, pageNum int) error {
	userInfo, err := c.getUserInfoViewData(ctx)
	if err != nil {
		return err
	}

	return web.RenderPage(eCtx, http.StatusOK, page.Log(userInfo, pageNum))
}

func (c *LogController) handleNavLog(eCtx echo.Context, pageNum int) error {
	return web.RenderHx(eCtx, http.StatusOK, hx.LogNav(pageNum))
}

func (c *LogController) handleGetLogContent(eCtx echo.Context, ctx context.Context, pageNum int,
) error {
	// Get view data
	summary, entries, err := c.getLogViewData(ctx, pageNum)
	if err != nil {
		return err
	}

	// Render
	return web.RenderHx(eCtx, http.StatusOK, hx.LogContent(summary, entries))
}

func (c *LogController) getLogViewData(ctx context.Context, pageNum int) (*vm.LogSummary,
	*vm.ListEntries, error) {
	// Get current user information
	userId := getCurrentUserId(ctx)
	userContract, err := c.getUserContract(ctx, userId)
	if err != nil {
		return nil, nil, err
	}

	// Get view data
	if pageNum > 1 {
		entries, err := c.getLogEntriesViewData(ctx, userId, userContract, pageNum)
		return nil, entries, err
	} else {
		summary, err := c.getLogSummaryViewData(ctx, userId, userContract)
		if err != nil {
			return nil, nil, err
		}
		entries, err := c.getLogEntriesViewData(ctx, userId, userContract, pageNum)
		if err != nil {
			return nil, nil, err
		}
		return summary, entries, nil
	}
}

func (c *LogController) getLogSummaryViewData(ctx context.Context, userId int,
	userContract *model.Contract) (*vm.LogSummary, error) {
	// Get total work summary data
	totalWorkSummary, err := c.eServ.GetTotalWorkSummaryByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Get month work summary data
	now := time.Now()
	year, month := now.Year(), now.Month()
	monthWorkSummary, err := c.eServ.GetMonthWorkSummaryByUserId(ctx, userId, year, month)
	if err != nil {
		return nil, err
	}

	// Create view model
	return c.mapper.CreateLogSummaryViewModel(userContract, now, totalWorkSummary, monthWorkSummary),
		nil
}

func (c *LogController) getLogEntriesViewData(ctx context.Context, userId int,
	userContract *model.Contract, pageNum int) (*vm.ListEntries, error) {
	// Get entries
	cnt, entries, entryTypesMap, entryActivitiesMap, err := c.getEntryData(ctx, userId, pageNum)
	if err != nil {
		return nil, err
	}

	// Create view model
	totPageNum := calculateNumberOfTotalPages(cnt, pageSize)
	return c.mapper.CreateLogEntriesViewModel(userContract, pageNum, totPageNum, entries,
		entryTypesMap, entryActivitiesMap), nil
}

func (c *LogController) getEntryData(ctx context.Context, userId, pageNum int) (int, []*model.Entry,
	map[int]*model.EntryType, map[int]*model.EntryActivity, error) {
	// Get entries
	offset, limit := calculateOffsetLimitFromPageNumber(pageNum)
	entries, cnt, err := c.eServ.GetDateEntriesByUserId(ctx, userId, offset, limit)
	if err != nil {
		return 0, nil, nil, nil, err
	}

	// Get entry master data
	entryTypesMap, entryActivitiesMap, err := c.getEntryMasterDataMap(ctx)
	if err != nil {
		return 0, nil, nil, nil, err
	}

	return cnt, entries, entryTypesMap, entryActivitiesMap, nil
}

// --- Helper functions ---

func (c *LogController) getGetLogParams(ctx echo.Context) (int, error) {
	pageNum, avail, err := getPageNumberQueryParam(ctx)
	if err != nil {
		return 0, err
	}
	if !avail {
		pageNum = 1
	}
	return pageNum, nil
}
