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

// GetOverviewHandler returns a handler for "GET /overview".
func (c *OverviewController) GetOverviewHandler() echo.HandlerFunc {
	return c.handler(func(eCtx echo.Context, ctx context.Context) error {
		userInfo, err := c.getUserInfoViewData(ctx)
		if err != nil {
			return err
		}

		year, month, err := c.getGetOverviewParams(eCtx)
		if err != nil {
			return err
		}
		monthStr := formatMonth(year, month)

		return web.RenderPage(eCtx, http.StatusOK, page.Overview(userInfo, monthStr))
	})
}

// GetHxNavHandler returns a handler for "GET /hx/overview".
func (c *OverviewController) GetHxNavHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		year, month, err := c.getGetOverviewParams(eCtx)
		if err != nil {
			return err
		}

		web.HtmxPushUrl(eCtx, c.buildOverviewUrl(year, month))
		return web.RenderHx(eCtx, http.StatusOK, hx.OverviewNav())
	})
}

// GetHxContentHandler returns a handler for "GET /hx/overview/content".
func (c *OverviewController) GetHxContentHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		year, month, err := c.getGetOverviewParams(eCtx)
		if err != nil {
			return err
		}

		overviewEntries, err := c.getOverviewViewData(ctx, year, month)
		if err != nil {
			return err
		}

		web.HtmxPushUrl(eCtx, c.buildOverviewUrl(year, month))
		return web.RenderHx(eCtx, http.StatusOK, hx.OverviewContent(overviewEntries))
	})
}

// GetOverviewExportHandler returns a handler for "GET /overview/export".
func (c *OverviewController) GetOverviewExportHandler() echo.HandlerFunc {
	return c.resourceHandler(func(eCtx echo.Context, ctx context.Context) error {
		year, month, err := c.getGetOverviewParams(eCtx)
		if err != nil {
			return err
		}

		overviewEntries, err := c.getOverviewViewData(ctx, year, month)
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("work-log-export-%s.xlsx", overviewEntries.CurrMonth)
		file := c.exporter.ExportOverviewEntries(overviewEntries)

		return web.WriteFile(eCtx, fileName, file)
	})
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

func (c *OverviewController) buildOverviewUrl(year int, month int) string {
	if year != 0 && month != 0 {
		return "/overview?" + buildMonthQueryParam(year, month)
	}
	return "/overview"
}

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
