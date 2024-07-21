package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/export"
	"kellnhofer.com/work-log/web/mapper"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/view/hx"
	"kellnhofer.com/work-log/web/view/page"
)

// OverviewController handles requests for overview endpoints.
type OverviewController struct {
	baseController

	mapper   *mapper.OverviewMapper
	exporter *export.OverviewExporter
}

// NewOverviewController creates a new overview controller.
func NewOverviewController(uServ *service.UserService, eServ *service.EntryService,
) *OverviewController {
	overviewMapper := mapper.NewOverviewMapper()
	overviewExporter := export.NewOverviewExporter()
	return &OverviewController{
		baseController: baseController{
			uServ:  uServ,
			eServ:  eServ,
			mapper: &overviewMapper.Mapper,
		},
		mapper:   overviewMapper,
		exporter: overviewExporter,
	}
}

// --- Endpoints ---

// GetOverviewHandler returns a handler for "GET /overview".
func (c *OverviewController) GetOverviewHandler() echo.HandlerFunc {
	return c.handler(func(eCtx echo.Context, ctx context.Context, isHtmxReq bool) error {
		year, month, err := c.getGetOverviewParams(eCtx)
		if err != nil {
			return err
		}

		monthStr := formatMonth(year, month)

		if !isHtmxReq {
			return c.handleShowOverview(eCtx, ctx, monthStr)
		} else {
			return c.handleNavOverview(eCtx, monthStr)
		}
	})
}

// GetOverviewContentHandler returns a handler for "GET /overview/content".
func (c *OverviewController) GetOverviewContentHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		year, month, err := c.getGetOverviewParams(eCtx)
		if err != nil {
			return err
		}
		return c.handleOverviewContent(eCtx, ctx, year, month)
	})
}

// GetOverviewExportHandler returns a handler for "GET /overview/export".
func (c *OverviewController) GetOverviewExportHandler() echo.HandlerFunc {
	return c.resourceHandler(func(eCtx echo.Context, ctx context.Context) error {
		year, month, err := c.getGetOverviewParams(eCtx)
		if err != nil {
			return err
		}
		return c.handleExportOverview(eCtx, ctx, year, month)
	})
}

// --- Handler functions ---

func (c *OverviewController) handleShowOverview(eCtx echo.Context, ctx context.Context, month string,
) error {
	userInfo, err := c.getUserInfoViewData(ctx)
	if err != nil {
		return err
	}

	return web.RenderPage(eCtx, http.StatusOK, page.Overview(userInfo, month))
}

func (c *OverviewController) handleNavOverview(eCtx echo.Context, month string) error {
	return web.RenderHx(eCtx, http.StatusOK, hx.OverviewNav(month))
}

func (c *OverviewController) handleOverviewContent(eCtx echo.Context, ctx context.Context, year int,
	month int) error {
	// Get view data
	overviewEntries, err := c.getOverviewViewData(ctx, year, month)
	if err != nil {
		return err
	}

	// Render
	return web.RenderHx(eCtx, http.StatusOK, hx.OverviewContent(overviewEntries))
}

func (c *OverviewController) handleExportOverview(eCtx echo.Context, ctx context.Context, year int,
	month int) error {
	// Get view data
	data, err := c.getOverviewViewData(ctx, year, month)
	if err != nil {
		return err
	}

	// Create file
	fileName := fmt.Sprintf("work-log-export-%s.xlsx", data.CurrMonth)
	file := c.exporter.ExportOverviewEntries(data)

	// Write file
	return web.WriteFile(eCtx, fileName, file)
}

func (c *OverviewController) getOverviewViewData(ctx context.Context, year int, month int,
) (*vm.OverviewEntries, error) {
	// Get current user information
	userId := getCurrentUserId(ctx)
	userContract, err := c.getUserContract(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Get entries
	entries, err := c.eServ.GetMonthEntriesByUserId(ctx, userId, year, month)
	if err != nil {
		return nil, err
	}
	// Get entry master data
	entryTypesMap, entryActivitiesMap, err := c.getEntryMasterDataMap(ctx)
	if err != nil {
		return nil, err
	}

	// Create view model
	return c.mapper.CreateOverviewEntriesViewModel(userContract, year, month, entries,
		entryTypesMap, entryActivitiesMap), nil
}

// --- Helper functions ---

func (c *OverviewController) getGetOverviewParams(ctx echo.Context) (int, int, error) {
	// Get year and month
	y, m, avail, err := getMonthQueryParam(ctx)
	if err != nil {
		return 0, 0, err
	}

	// Was a year and month provided?
	if !avail {
		// Get current year/month
		t := time.Now()
		return t.Year(), int(t.Month()), nil
	} else {
		// Use these
		return y, m, nil
	}
}
