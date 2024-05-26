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
	return func(eCtx echo.Context) error {
		isHtmxReq := web.IsHtmxRequest(eCtx)

		year, month, isPageReq, err := c.getGetOverviewParams(eCtx)
		if err != nil {
			return err
		}

		ctx := getContext(eCtx)

		if !isHtmxReq {
			return c.handleShowOverview(eCtx, ctx, year, month)
		} else if !isPageReq {
			return c.handleHxNavOverview(eCtx, ctx, year, month)
		} else {
			return c.handleHxGetOverviewPage(eCtx, ctx, year, month)
		}
	}
}

// GetOverviewExportHandler returns a handler for "GET /overview/export".
func (c *OverviewController) GetOverviewExportHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		year, month, _, err := c.getGetOverviewParams(eCtx)
		if err != nil {
			return err
		}
		return c.handleExportOverview(eCtx, getContext(eCtx), year, month)
	}
}

func (c *OverviewController) getGetOverviewParams(eCtx echo.Context) (int, int, bool, error) {
	// Get year and month
	y, m, avail, err := getMonthQueryParam(eCtx)
	if err != nil {
		return 0, 0, false, err
	}

	// Was a year and month provided?
	if !avail {
		// Get current year/month
		t := time.Now()
		return t.Year(), int(t.Month()), false, nil
	} else {
		// Use these
		return y, m, true, nil
	}
}

// --- Handler functions ---

func (c *OverviewController) handleShowOverview(eCtx echo.Context, ctx context.Context, year int,
	month int) error {
	// Get view data
	userInfo, err := c.getUserInfoViewData(ctx)
	if err != nil {
		return err
	}
	overviewEntries, err := c.getOverviewViewData(ctx, year, month)
	if err != nil {
		return err
	}

	// Render
	return web.RenderPage(eCtx, http.StatusOK, page.Overview(userInfo, overviewEntries))
}

func (c *OverviewController) handleHxNavOverview(eCtx echo.Context, ctx context.Context, year int,
	month int) error {
	// Get view data
	overviewEntries, err := c.getOverviewViewData(ctx, year, month)
	if err != nil {
		return err
	}

	// Render
	return web.RenderHx(eCtx, http.StatusOK, hx.OverviewNav(overviewEntries))
}

func (c *OverviewController) handleHxGetOverviewPage(eCtx echo.Context, ctx context.Context, year int,
	month int) error {
	// Get view data
	overviewEntries, err := c.getOverviewViewData(ctx, year, month)
	if err != nil {
		return err
	}

	// Render
	return web.RenderHx(eCtx, http.StatusOK, hx.OverviewPage(overviewEntries))
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
