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
	baseController

	mapper   *mapper.EntryMapper
	exporter *export.EntriesExporter
}

// NewExportController creates a new export controller.
func NewExportController(uServ *service.UserService, eServ *service.EntryService) *ExportController {
	entryMapper := mapper.NewEntryMapper()
	exporter := export.NewEntriesExporter()
	return &ExportController{
		baseController: baseController{
			uServ:  uServ,
			eServ:  eServ,
			mapper: &entryMapper.Mapper,
		},
		mapper:   entryMapper,
		exporter: exporter,
	}
}

// GetExportHandler returns a handler for "GET /export".
func (c *ExportController) GetExportHandler() echo.HandlerFunc {
	return c.handler(func(eCtx echo.Context, ctx context.Context) error {
		query := getQueryQueryParam(eCtx)

		exportFilter, err := c.parseQueryString(getCurrentUserId(ctx), query)
		if err != nil {
			return err
		}

		exportDetails, err := c.getExportDetailsViewData(ctx, exportFilter)
		if err != nil {
			return err
		}

		exportEntries, err := c.getExportEntriesViewData(ctx, exportFilter)
		if err != nil {
			return err
		}

		timestamp := time.Now().Format(constant.ExportTimestampFormat)
		fileName := fmt.Sprintf(constant.ExportFileNameTemplate, timestamp, "xlsx")
		file := c.exporter.ExportEntries(exportDetails, exportEntries)

		return web.WriteFile(eCtx, fileName, file)
	})
}

func (c *ExportController) getExportDetailsViewData(ctx context.Context,
	exportFilter *model.FieldEntryFilter) (*vm.EntryFilterDetails, error) {
	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx, exportFilter.TypeId)
	if err != nil {
		return nil, err
	}

	// Create view model
	return c.mapper.CreateEntryFilterDetailsViewModel(exportFilter, entryTypes, entryActivities),
		nil
}

func (c *ExportController) getExportEntriesViewData(ctx context.Context,
	exportFilter *model.FieldEntryFilter) (*vm.ListEntries, error) {
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
