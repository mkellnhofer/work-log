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
	"kellnhofer.com/work-log/web/view/hx"
	"kellnhofer.com/work-log/web/view/page"
)

type searchInput struct {
	byType     string
	typeId     string
	byDate     string
	startDate  string
	endDate    string
	byActivity string
	activityId string
	text       string
}

// SearchController handles requests for search endpoints.
type SearchController struct {
	baseController

	mapper *mapper.SearchMapper
}

// NewSearchController creates a new search controller.
func NewSearchController(uServ *service.UserService, eServ *service.EntryService) *SearchController {
	searchMapper := mapper.NewSearchMapper()
	return &SearchController{
		baseController: baseController{
			uServ:  uServ,
			eServ:  eServ,
			mapper: &searchMapper.Mapper,
		},
		mapper: searchMapper,
	}
}

// GetSearchHandler returns a handler for "GET /search".
func (c *SearchController) GetSearchHandler() echo.HandlerFunc {
	return c.handler(func(eCtx echo.Context, ctx context.Context) error {
		userInfo, err := c.getUserInfoViewData(ctx)
		if err != nil {
			return err
		}

		isAdvanced, query, pageNum, _, err := c.getGetSearchParams(eCtx)
		if err != nil {
			return err
		}

		searchFilter, err := c.parseQueryString(getCurrentUserId(ctx), query)
		if err != nil {
			return err
		}

		searchQueryString := c.buildQueryString(searchFilter)
		searchDetails, err := c.getSearchDetailsViewData(ctx, searchFilter)
		if err != nil {
			return err
		}

		return web.RenderPage(eCtx, http.StatusOK, page.Search(userInfo, isAdvanced,
			searchQueryString, pageNum, searchDetails))
	})
}

// GetHxContentHandler returns a handler for "GET /hx/search/content".
func (c *SearchController) GetHxContentHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		isAdvanced, query, pageNum, _, err := c.getGetSearchParams(eCtx)
		if err != nil {
			return err
		}

		searchFilter, err := c.parseQueryString(getCurrentUserId(ctx), query)
		if err != nil {
			return err
		}

		searchQueryString := c.buildQueryString(searchFilter)

		searchEntries, err := c.getSearchEntriesViewData(ctx, searchFilter, pageNum)
		if err != nil {
			return err
		}

		web.HtmxPushUrl(eCtx, c.buildSearchUrl(isAdvanced, searchFilter, pageNum))
		return web.RenderHx(eCtx, http.StatusOK, hx.SearchContent(isAdvanced, searchQueryString,
			searchEntries))
	})
}

// GetHxModalHandler returns a handler for "GET /hx/search-modal".
func (c *SearchController) GetHxModalHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		isAdvanced, query, _, _, err := c.getGetSearchParams(eCtx)
		if err != nil {
			return err
		}

		searchFilter, err := c.parseQueryString(getCurrentUserId(ctx), query)
		if err != nil {
			return err
		}

		searchQuery, err := c.getSearchQueryViewData(ctx, searchFilter)
		if err != nil {
			return err
		}

		return web.RenderHx(eCtx, http.StatusOK, hx.SearchModal(isAdvanced, searchQuery))
	})
}

// PostHxModalHandler returns a handler for "POST /hx/search-modal".
func (c *SearchController) PostHxModalHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		isAdvanced := getAdvancedQueryParam(eCtx)

		userId := getCurrentUserId(ctx)
		searchInput := c.getPostSearchInput(eCtx)
		searchFilter, err := c.createSearchFilter(userId, searchInput)
		if err != nil {
			searchErrorMessage := loc.GetErrorMessageString(getErrorCode(err))
			web.HtmxRetarget(eCtx, "#wl-modal-error-container")
			return web.RenderHx(eCtx, http.StatusOK, hx.ModalError(searchErrorMessage))
		}

		searchQueryString := c.buildQueryString(searchFilter)
		searchDetails, err := c.getSearchDetailsViewData(ctx, searchFilter)
		if err != nil {
			return err
		}

		web.HtmxPushUrl(eCtx, c.buildSearchUrl(isAdvanced, searchFilter, 1))
		return web.RenderHx(eCtx, http.StatusOK, hx.Search(isAdvanced, searchQueryString, 1,
			searchDetails))
	})
}

// GetHxModalActivitiesHandler returns a handler for "GET /hx/search-modal/activities".
func (c *SearchController) GetHxModalActivitiesHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		entryTypeId, err := getTypeIdQueryParam(eCtx)
		if err != nil {
			return err
		}

		entryActivities, err := c.getEntryActivities(ctx, entryTypeId)
		if err != nil {
			return err
		}

		viewData := c.mapper.CreateEntryActivitiesViewModel(entryActivities)

		return web.RenderHx(eCtx, http.StatusOK, hx.SearchFormActivityOptions(viewData))
	})
}

// PostHxModalCancelHandler returns a handler for "POST /hx/search-modal/cancel".
func (c *SearchController) PostHxModalCancelHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		return eCtx.NoContent(http.StatusOK)
	})
}

func (c *SearchController) getSearchQueryViewData(ctx context.Context,
	searchFilter *model.EntriesFilter) (*vm.SearchQuery, error) {
	// Create default values
	filter := searchFilter
	if filter == nil {
		filter = model.NewEntriesFilter()
	}
	if !filter.ByType {
		filter.TypeId = model.EntryTypeIdWork
	}
	if !filter.ByTime {
		filter.StartTime = time.Now()
		filter.EndTime = time.Now()
	}

	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx, filter.TypeId)
	if err != nil {
		return nil, err
	}

	// Create view model
	return c.mapper.CreateSearchQueryViewModel(filter, entryTypes, entryActivities), nil
}

func (c *SearchController) getSearchDetailsViewData(ctx context.Context,
	searchFilter *model.EntriesFilter) (*vm.SearchDetails, error) {
	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx, searchFilter.TypeId)
	if err != nil {
		return nil, err
	}

	// Create view model
	return c.mapper.CreateSearchDetailsViewModel(searchFilter, entryTypes, entryActivities), nil
}

func (c *SearchController) getSearchEntriesViewData(ctx context.Context,
	searchFilter *model.EntriesFilter, pageNum int) (*vm.ListEntries, error) {
	if c.isFilterEmpty(searchFilter) {
		return &vm.ListEntries{}, nil
	}

	// Create entries searchSort
	searchSort := model.NewEntriesSort()
	searchSort.ByTime = model.DescSorting

	// Get entries
	offset, limit := calculateOffsetLimitFromPageNumber(pageNum)
	entries, cnt, err := c.eServ.GetDateEntries(ctx, searchFilter, searchSort, offset, limit)
	if err != nil {
		return nil, err
	}
	// Get entry master data
	entryTypesMap, entryActivitiesMap, err := c.getEntryMasterDataMap(ctx)
	if err != nil {
		return nil, err
	}

	// Create view model
	totPageNum := calculateNumberOfTotalPages(cnt, pageSize)
	return c.mapper.CreateSearchEntriesViewModel(pageNum, totPageNum, entries, entryTypesMap,
		entryActivitiesMap), nil
}

// --- Search query functions ---

func (c *SearchController) getPostSearchInput(eCtx echo.Context) *searchInput {
	return &searchInput{
		byType:     eCtx.FormValue("by-type"),
		typeId:     eCtx.FormValue("type"),
		byDate:     eCtx.FormValue("by-date"),
		startDate:  eCtx.FormValue("start-date"),
		endDate:    eCtx.FormValue("end-date"),
		byActivity: eCtx.FormValue("by-activity"),
		activityId: eCtx.FormValue("activity"),
		text:       eCtx.FormValue("text"),
	}
}

func (c *SearchController) createSearchFilter(userId int, input *searchInput) (*model.EntriesFilter,
	error) {
	filter := model.NewEntriesFilter()
	filter.ByUser = true
	filter.UserId = userId

	var err error

	// Create type ID filter
	filter.ByType = input.byType == "on"
	filter.TypeId, err = parseId(input.typeId, false)
	if err != nil {
		return nil, err
	}

	// Create start/end time filter
	filter.ByTime = input.byDate == "on"
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

	// Create activity ID filter
	filter.ByActivity = input.byActivity == "on"
	filter.ActivityId, err = parseId(input.activityId, true)
	if err != nil {
		return nil, err
	}

	// Create description filter
	if err = validateMaxStringLength(input.text, model.MaxLengthEntryDescription,
		e.ValDescriptionTooLong); err != nil {
		return nil, err
	}
	filter.ByDescription = input.text != ""
	filter.Description = input.text

	// If search query is empty: Create empty description filter
	if !filter.ByType && !filter.ByTime && !filter.ByActivity && !filter.ByDescription {
		filter.ByDescription = true
		filter.Description = ""
	}

	return filter, nil
}

// --- Helper functions ---

func (c *SearchController) buildSearchUrl(isAdvanced bool, searchFilter *model.EntriesFilter,
	pageNum int) string {
	url := "/search?"
	if isAdvanced {
		url = url + "adv=1&"
	}
	query := c.buildQueryString(searchFilter)
	if query != "" {
		url = url + "query=" + query + "&"
	}
	if pageNum != 0 {
		url = url + buildPageNumberQueryParam(pageNum)
	}
	return url
}

func (c *SearchController) getGetSearchParams(eCtx echo.Context) (bool, string, int, bool, error) {
	isAdvanced := getAdvancedQueryParam(eCtx)
	searchQuery := getQueryQueryParam(eCtx)
	pageNum, pageNumAvail, err := getPageNumberQueryParam(eCtx)
	if err != nil {
		return false, "", 0, false, err
	}
	if !pageNumAvail {
		pageNum = 1
	}
	return isAdvanced, searchQuery, pageNum, pageNumAvail, nil
}
