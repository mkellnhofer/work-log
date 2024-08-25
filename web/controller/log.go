package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/view"
	"kellnhofer.com/work-log/web/view/hx"
	"kellnhofer.com/work-log/web/view/page"
)

type exportInput struct {
	startDate string
	endDate   string
}

// LogController handles requests for log endpoints.
type LogController struct {
	handlerHelper
	baseUserController
	baseEntryController
	entryFilterHelper

	mapper *mapper.LogMapper
}

// NewLogController creates a new log controller.
func NewLogController(uServ *service.UserService, eServ *service.EntryService) *LogController {
	return &LogController{
		baseUserController:  *newBaseUserController(uServ),
		baseEntryController: *newBaseEntryController(eServ),
		mapper:              mapper.NewLogMapper(),
	}
}

// GetLogHandler returns a handler for "GET /log".
func (c *LogController) GetLogHandler() echo.HandlerFunc {
	return c.handler(func(eCtx echo.Context, ctx context.Context) error {
		userInfo, err := c.getUserInfoViewData(ctx)
		if err != nil {
			return err
		}

		pageNum, err := c.getGetLogParams(eCtx)
		if err != nil {
			return err
		}

		return web.RenderPage(eCtx, http.StatusOK, page.Log(userInfo, pageNum))
	})
}

// GetHxNavHandler returns a handler for "GET /hx/log".
func (c *LogController) GetHxNavHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		pageNum, err := c.getGetLogParams(eCtx)
		if err != nil {
			return err
		}

		web.HtmxPushUrl(eCtx, c.buildLogUrl(pageNum))
		return web.RenderHx(eCtx, http.StatusOK, hx.Log())
	})
}

// GetHxContentHandler returns a handler for "GET /hx/log/content".
func (c *LogController) GetHxContentHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		pageNum, err := c.getGetLogParams(eCtx)
		if err != nil {
			return err
		}

		logSummary, logEntries, err := c.getLogViewData(ctx, pageNum)
		if err != nil {
			return err
		}

		web.HtmxPushUrl(eCtx, c.buildLogUrl(pageNum))
		return web.RenderHx(eCtx, http.StatusOK, hx.LogContent(logSummary, logEntries))
	})
}

// GetHxExportModalHandler returns a handler for "GET /hx/log-export-modal".
func (c *LogController) GetHxExportModalHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		now := time.Now()
		startDate := now.Format(view.DateStringFormat)
		endDate := now.Format(view.DateStringFormat)
		return web.RenderHx(eCtx, http.StatusOK, hx.LogExportModal(startDate, endDate))
	})
}

// PostHxExportModalHandler returns a handler for "POST /hx/log-export-modal".
func (c *LogController) PostHxExportModalHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		exportInput := c.getPostExportInput(eCtx)
		exportFilter, err := c.createExportFilter(getCurrentUserId(ctx), exportInput)
		if err != nil {
			searchErrorMessage := loc.GetErrorMessageString(getErrorCode(err))
			web.HtmxRetarget(eCtx, "#wl-modal-error-container")
			return web.RenderHx(eCtx, http.StatusOK, hx.ModalError(searchErrorMessage))
		}

		exportQuery := c.buildQueryString(exportFilter)
		exportUrl := c.buildExportUrl(true, exportQuery)

		web.HtmxTriggerAfterSwap(eCtx, fmt.Sprintf("{ \"downloadFile\": \"%s\" }", exportUrl))
		return eCtx.NoContent(http.StatusOK)
	})
}

// PostHxExportModalCancelHandler returns a handler for "POST /hx/log-export-modal/cancel".
func (c *LogController) PostHxExportModalCancelHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		return eCtx.NoContent(http.StatusOK)
	})
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

// --- Export query functions ---

func (c *LogController) getPostExportInput(eCtx echo.Context) *exportInput {
	return &exportInput{
		startDate: eCtx.FormValue("start-date"),
		endDate:   eCtx.FormValue("end-date"),
	}
}

func (c *LogController) createExportFilter(userId int, input *exportInput) (*model.FieldEntryFilter,
	error) {
	filter := model.NewFieldEntryFilter()
	filter.ByUser = true
	filter.UserId = userId

	var err error

	// Create start/end time filter
	filter.ByTime = true
	filter.StartTime, err = parseDateTime(input.startDate, "00:00", e.ValStartDateInvalid)
	if err != nil {
		return nil, err
	}
	filter.EndTime, err = parseDateTime(input.endDate, "23:59", e.ValEndDateInvalid)
	if err != nil {
		return nil, err
	}
	if filter.EndTime.Before(filter.StartTime) {
		err := e.NewError(e.LogicEntryDateIntervalInvalid, fmt.Sprintf("End date %s before "+
			"start time %s.", input.endDate, input.startDate))
		log.Debug(err.StackTrace())
		return nil, err
	}

	return filter, nil
}

// --- Helper functions ---

func (c *LogController) buildLogUrl(pageNum int) string {
	url := "/log"
	if pageNum != 0 {
		url = url + "?" + buildPageNumberQueryParam(pageNum)
	}
	return url
}

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

func (c *LogController) buildExportUrl(isAdvanced bool, query string) string {
	url := "/export?"
	if isAdvanced {
		url = url + "adv=1&"
	}
	if query != "" {
		url = url + "query=" + query
	}
	return url
}
