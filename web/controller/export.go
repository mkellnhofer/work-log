package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/constant"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/export"
	"kellnhofer.com/work-log/web/mapper"
	vm "kellnhofer.com/work-log/web/model"
)

// ExportController handles requests for export endpoints.
type ExportController struct {
	handlerHelper
	baseEntryController
	entryFilterHelper

	mapper   *mapper.EntryMapper
	exporter *export.EntriesExporter
}

// NewExportController creates a new export controller.
func NewExportController(eServ *service.EntryService) *ExportController {
	return &ExportController{
		baseEntryController: *newBaseEntryController(eServ),
		mapper:              mapper.NewEntryMapper(),
		exporter:            export.NewEntriesExporter(),
	}
}

// GetExportHandler returns a handler for "GET /export".
func (c *ExportController) GetExportHandler() echo.HandlerFunc {
	return c.handler(func(eCtx echo.Context, ctx context.Context) error {
		isAdvanced, query := c.getGetExportParams(eCtx)

		exportFilter, err := c.parseQueryString(getCurrentUserId(ctx), isAdvanced, query)
		if err != nil {
			return err
		}

		exportFilterDetails, err := c.getFilterDetailsViewData(ctx, exportFilter)
		if err != nil {
			return err
		}

		exportEntries, err := c.getExportEntriesViewData(ctx, exportFilter)
		if err != nil {
			return err
		}

		timestamp := time.Now().Format(constant.ExportTimestampFormat)
		fileName := fmt.Sprintf(constant.ExportFileNameTemplate, timestamp, "xlsx")
		file := c.exporter.ExportEntries(exportFilterDetails, exportEntries)

		return web.WriteFile(eCtx, fileName, file)
	})
}

func (c *ExportController) getExportEntriesViewData(ctx context.Context,
	exportFilter model.EntryFilter) (*vm.ListEntries, error) {
	if c.isFilterEmpty(exportFilter) {
		return &vm.ListEntries{}, nil
	}

	// Create entries searchSort
	exportSort := model.NewEntrySort()
	exportSort.ByTime = model.AscSorting

	// Get all entries (no pagination for export)
	entries, _, err := c.eServ.GetDateEntries(ctx, exportFilter, exportSort, 0, 0)
	if err != nil {
		return nil, err
	}
	// Get entry master data
	entryTypesMap, entryActivitiesMap, err := c.getEntryMasterDataMap(ctx)
	if err != nil {
		return nil, err
	}

	// Create view model
	return c.mapper.CreateListEntriesViewModel(entries, entryTypesMap, entryActivitiesMap), nil
}

// --- Helper functions ---

func (c *ExportController) getGetExportParams(eCtx echo.Context) (bool, string) {
	isAdvanced := getAdvancedQueryParam(eCtx)
	query := getQueryQueryParam(eCtx)
	return isAdvanced, query
}
