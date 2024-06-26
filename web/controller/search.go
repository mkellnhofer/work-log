package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/pkg/util"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/view/hx"
	"kellnhofer.com/work-log/web/view/page"
)

const searchDateTimeFormat = "200601021504"

type searchInput struct {
	byType        string
	typeId        string
	byDate        string
	startDate     string
	endDate       string
	byActivity    string
	activityId    string
	byDescription string
	description   string
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

// --- Endpoints ---

// GetSearchHandler returns a handler for "GET /search".
func (c *SearchController) GetSearchHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /search.")

		isHtmxReq := web.IsHtmxRequest(eCtx)

		searchQuery, pageNum, isPageReq, err := c.getSearchParams(eCtx)
		if err != nil {
			return err
		}

		ctx := getContext(eCtx)

		if !isHtmxReq {
			return c.handleShowSearch(eCtx, ctx, searchQuery, pageNum)
		} else if !isPageReq {
			return c.handleHxNavSearch(eCtx, ctx)
		} else {
			return c.handleHxGetSearchPage(eCtx, ctx, searchQuery, pageNum)
		}
	}
}

func (c *SearchController) getSearchParams(eCtx echo.Context) (string, int, bool, error) {
	searchQuery, _ := getSearchQueryParam(eCtx)
	pageNum, pageNumAvail, err := getPageNumberQueryParam(eCtx)
	if err != nil {
		return "", 0, false, err
	}
	if !pageNumAvail {
		pageNum = 1
	}
	return searchQuery, pageNum, pageNumAvail, nil
}

// PostSearchHandler returns a handler for "POST /search".
func (c *SearchController) PostSearchHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /search.")

		isHtmxReq := web.IsHtmxRequest(eCtx)
		if !isHtmxReq {
			err := e.NewError(e.ValUnknown, "Not a HTMX request.")
			log.Debug(err.StackTrace())
			return err
		}

		ctx := getContext(eCtx)
		input := c.getSearchInput(eCtx)

		return c.handleHxExecuteSearch(eCtx, ctx, input)
	}
}

func (c *SearchController) getSearchInput(eCtx echo.Context) *searchInput {
	return &searchInput{
		byType:        eCtx.FormValue("by-type"),
		typeId:        eCtx.FormValue("type"),
		byDate:        eCtx.FormValue("by-date"),
		startDate:     eCtx.FormValue("start-date"),
		endDate:       eCtx.FormValue("end-date"),
		byActivity:    eCtx.FormValue("by-activity"),
		activityId:    eCtx.FormValue("activity"),
		byDescription: eCtx.FormValue("by-description"),
		description:   eCtx.FormValue("description"),
	}
}

// --- Handler functions ---

func (c *SearchController) handleShowSearch(eCtx echo.Context, ctx context.Context, query string,
	pageNum int) error {
	// Create search filter
	searchFilter, err := c.parseSearchQueryString(getCurrentUserId(ctx), query)
	searchErrorMessage := ""
	if err != nil {
		searchFilter = model.NewEntriesFilter()
		searchErrorMessage = loc.GetErrorMessageString(getErrorCode(err))
	}

	// Create view model
	userInfo, err := c.getUserInfoViewData(ctx)
	if err != nil {
		return err
	}
	search, err := c.getSearchViewData(ctx, searchFilter)
	if err != nil {
		return err
	}
	searchEntries, err := c.getSearchEntriesViewData(ctx, searchFilter, pageNum)
	if err != nil {
		return err
	}

	// Render
	return web.Render(eCtx, http.StatusOK, page.Search(userInfo, searchErrorMessage, search,
		searchEntries))
}

func (c *SearchController) handleHxNavSearch(eCtx echo.Context, ctx context.Context) error {
	// Create search filter
	searchFilter, err := c.parseSearchQueryString(getCurrentUserId(ctx), "")
	if err != nil {
		return err
	}

	// Create view model
	search, err := c.getSearchViewData(ctx, searchFilter)
	if err != nil {
		return err
	}
	searchEntries := &vm.SearchEntries{}

	// Render
	return web.Render(eCtx, http.StatusOK, hx.SearchNav(search, searchEntries))
}

func (c *SearchController) handleHxExecuteSearch(eCtx echo.Context, ctx context.Context,
	searchInputs *searchInput) error {
	// Create search filter
	searchFilter, err := c.createSearchFilter(getCurrentUserId(ctx), searchInputs)
	if err != nil {
		searchErrorMessage := loc.GetErrorMessageString(getErrorCode(err))
		searchEntries := &vm.SearchEntries{}
		return web.Render(eCtx, http.StatusOK, hx.Search(searchErrorMessage, nil, searchEntries))
	}

	// Create view model
	search, err := c.getSearchViewData(ctx, searchFilter)
	if err != nil {
		return err
	}
	searchEntries, err := c.getSearchEntriesViewData(ctx, searchFilter, 1)
	if err != nil {
		return err
	}

	// Render
	web.HtmxPushUrl(eCtx, "/search?query="+searchEntries.Query)
	return web.Render(eCtx, http.StatusOK, hx.Search("", search, searchEntries))
}

func (c *SearchController) handleHxGetSearchPage(eCtx echo.Context, ctx context.Context, query string,
	pageNum int) error {
	// Create search filter
	searchFilter, err := c.parseSearchQueryString(getCurrentUserId(ctx), query)
	if err != nil {
		return err
	}

	// Create view model
	searchEntries, err := c.getSearchEntriesViewData(ctx, searchFilter, pageNum)
	if err != nil {
		return err
	}

	// Render
	return web.Render(eCtx, http.StatusOK, hx.SearchPage(searchEntries))
}

func (c *SearchController) getSearchViewData(ctx context.Context, searchFilter *model.EntriesFilter,
) (*vm.Search, error) {
	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return nil, err
	}

	// Create default values
	filter := searchFilter
	if filter == nil {
		filter = model.NewEntriesFilter()
	}
	typeId := filter.TypeId
	if !filter.ByType && len(entryTypes) > 0 {
		typeId = entryTypes[0].Id
	}
	sTime := filter.StartTime
	eTime := filter.EndTime
	if !filter.ByTime {
		sTime = time.Now()
		eTime = time.Now()
	}

	// Create view model
	return c.mapper.CreateSearchViewModel(filter.ByType, typeId, filter.ByTime, sTime, eTime,
		filter.ByActivity, filter.ActivityId, filter.ByDescription, filter.Description, entryTypes,
		entryActivities), nil
}

func (c *SearchController) getSearchEntriesViewData(ctx context.Context,
	searchFilter *model.EntriesFilter, pageNum int) (*vm.SearchEntries, error) {
	if c.isSearchFilterEmpty(searchFilter) {
		return &vm.SearchEntries{}, nil
	}

	// Create query string
	query := c.buildSearchQueryString(searchFilter)
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
	return c.mapper.CreateSearchEntriesViewModel(query, pageNum, pageSize, cnt, entries,
		entryTypesMap, entryActivitiesMap), nil
}

// --- Search query functions ---

func (c *SearchController) createSearchFilter(userId int, input *searchInput) (*model.EntriesFilter,
	error) {
	filter := model.NewEntriesFilter()
	filter.ByUser = true
	filter.UserId = userId

	var err error

	// Convert type ID
	filter.ByType = input.byType == "on"
	filter.TypeId, err = parseId(input.typeId, false)
	if err != nil {
		return nil, err
	}

	// Convert start/end time
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
		err := e.NewError(e.LogicEntrySearchDateIntervalInvalid, fmt.Sprintf("End date %s before "+
			"start time %s.", input.endDate, input.startDate))
		log.Debug(err.StackTrace())
		return nil, err
	}

	// Convert activity ID
	filter.ByActivity = input.byActivity == "on"
	filter.ActivityId, err = parseId(input.activityId, true)
	if err != nil {
		return nil, err
	}

	// Validate description
	filter.ByDescription = input.byDescription == "on"
	if err = validateStringLength(input.description, model.MaxLengthEntryDescription,
		e.ValDescriptionTooLong); err != nil {
		return nil, err
	}
	filter.Description = input.description

	// Check if search query is empty
	if !filter.ByType && !filter.ByTime && !filter.ByActivity && !filter.ByDescription {
		err := e.NewError(e.ValSearchInvalid, "Search query is empty.")
		log.Debug(err.StackTrace())
		return nil, err
	}

	return filter, nil
}

func (c *SearchController) isSearchFilterEmpty(filter *model.EntriesFilter) bool {
	return !filter.ByType && !filter.ByTime && !filter.ByActivity && !filter.ByDescription
}

func (c *SearchController) buildSearchQueryString(filter *model.EntriesFilter) string {
	var qps []string
	// Add parameter/value for entry type
	if filter.ByType {
		qps = append(qps, fmt.Sprintf("typ:%d", filter.TypeId))
	}
	// Add parameter/value for entry start/end time
	if filter.ByTime {
		qps = append(qps, fmt.Sprintf("tim:%s", c.formatSearchDateRange(filter.StartTime,
			filter.EndTime)))
	}
	// Add parameter/value for entry activity
	if filter.ByActivity {
		qps = append(qps, fmt.Sprintf("act:%d", filter.ActivityId))
	}
	// Add parameter/value for entry description
	if filter.ByDescription {
		qps = append(qps, fmt.Sprintf("des:%s", util.EncodeBase64(filter.Description)))
	}
	return strings.Join(qps[:], "|")
}

func (c *SearchController) parseSearchQueryString(userId int, query string) (*model.EntriesFilter,
	error) {
	filter := model.NewEntriesFilter()
	filter.ByUser = true
	filter.UserId = userId

	if query == "" {
		return filter, nil
	}

	qps := strings.Split(query, "|")

	for _, qp := range qps {
		pv := strings.Split(qp, ":")
		// Check if query part is invalid
		if len(pv) < 2 {
			err := e.NewError(e.ValSearchQueryInvalid, "Search query part is invalid.")
			log.Debug(err.StackTrace())
			return nil, err
		}

		p := pv[0]
		v := pv[1]
		var cErr error

		// Handle specific conversion
		switch p {
		// Convert value for entry type
		case "typ":
			filter.ByType = true
			filter.TypeId, cErr = strconv.Atoi(v)
		// Convert values for entry start/end time
		case "tim":
			filter.ByTime = true
			filter.StartTime, filter.EndTime, cErr = c.parseSearchDateRange(v)
		// Convert value for entry activity
		case "act":
			filter.ByActivity = true
			filter.ActivityId, cErr = strconv.Atoi(v)
		// Convert value for entry description
		case "des":
			filter.ByDescription = true
			filter.Description, cErr = util.DecodeBase64(v)
		// Unknown parameter
		default:
			err := e.NewError(e.ValSearchQueryInvalid, fmt.Sprintf("Search query parameter '%s' "+
				"is unknown.", p))
			log.Debug(err.StackTrace())
			return nil, err
		}

		// Check if a error occurred
		if cErr != nil {
			err := e.WrapError(e.ValSearchQueryInvalid, fmt.Sprintf("Search query parameter '%s' "+
				"has invalid value.", p), cErr)
			log.Debug(err.StackTrace())
			return nil, err
		}
	}
	return filter, nil
}

func (c *SearchController) formatSearchDateRange(startDate time.Time, endDate time.Time) string {
	return fmt.Sprintf("%s-%s", c.formatSearchDate(startDate), c.formatSearchDate(endDate))
}

func (c *SearchController) formatSearchDate(date time.Time) string {
	return date.Format(searchDateTimeFormat)
}

func (c *SearchController) parseSearchDateRange(dateRange string) (time.Time, time.Time, error) {
	se := strings.Split(dateRange, "-")
	if len(se) < 2 {
		return time.Time{}, time.Time{}, errors.New("invalid range")
	}
	startTime, err := c.parseSearchDate(se[0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endTime, err := c.parseSearchDate(se[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return startTime, endTime, nil
}

func (c *SearchController) parseSearchDate(date string) (time.Time, error) {
	return time.ParseInLocation(searchDateTimeFormat, date, time.Local)
}
