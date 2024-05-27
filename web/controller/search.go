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

// --- Endpoints ---

// GetSearchHandler returns a handler for "GET /search".
func (c *SearchController) GetSearchHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		isHtmxReq := web.IsHtmxRequest(eCtx)

		isAdvanced, searchQuery, pageNum, isPageReq, err := c.getGetSearchParams(eCtx)
		if err != nil {
			return err
		}

		ctx := getContext(eCtx)

		if !isHtmxReq {
			return c.handleShowSearch(eCtx, ctx, isAdvanced, searchQuery, pageNum)
		} else if !isPageReq {
			return c.handleHxNavSearch(eCtx, ctx)
		} else {
			return c.handleHxGetSearchPage(eCtx, ctx, isAdvanced, searchQuery, pageNum)
		}
	}
}

func (c *SearchController) getGetSearchParams(eCtx echo.Context) (bool, string, int, bool, error) {
	isAdvanced := getSearchAdvancedParam(eCtx)
	searchQuery := getSearchQueryParam(eCtx)
	pageNum, pageNumAvail, err := getPageNumberQueryParam(eCtx)
	if err != nil {
		return false, "", 0, false, err
	}
	if !pageNumAvail {
		pageNum = 1
	}
	return isAdvanced, searchQuery, pageNum, pageNumAvail, nil
}

// GetFormHandler returns a handler for "GET /search/form".
func (c *SearchController) GetFormHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		isHtmxReq := web.IsHtmxRequest(eCtx)
		if !isHtmxReq {
			err := e.NewError(e.ValUnknown, "Not a HTMX request.")
			log.Debug(err.StackTrace())
			return err
		}
		return c.handleHxGetForm(eCtx, getSearchAdvancedParam(eCtx))
	}
}

// GetActivitiesHandler returns a handler for "GET /search/form/activities".
func (c *SearchController) GetFormActivitiesHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		isHtmxReq := web.IsHtmxRequest(eCtx)
		if !isHtmxReq {
			err := e.NewError(e.ValUnknown, "Not a HTMX request.")
			log.Debug(err.StackTrace())
			return err
		}
		return c.handleHxGetFormActivities(eCtx, getContext(eCtx))
	}
}

// PostSearchHandler returns a handler for "POST /search".
func (c *SearchController) PostSearchHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		isHtmxReq := web.IsHtmxRequest(eCtx)
		if !isHtmxReq {
			err := e.NewError(e.ValUnknown, "Not a HTMX request.")
			log.Debug(err.StackTrace())
			return err
		}

		ctx := getContext(eCtx)
		isAdvanced := getSearchAdvancedParam(eCtx)
		input := c.getPostSearchInput(eCtx)

		return c.handleHxExecuteSearch(eCtx, ctx, isAdvanced, input)
	}
}

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

// --- Handler functions ---

func (c *SearchController) handleShowSearch(eCtx echo.Context, ctx context.Context, isAdvanced bool,
	query string, pageNum int) error {
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
	searchQuery, err := c.getSearchQueryViewData(ctx, isAdvanced, searchFilter)
	if err != nil {
		return err
	}
	searchEntries, err := c.getSearchEntriesViewData(ctx, searchFilter, pageNum)
	if err != nil {
		return err
	}

	// Render
	return web.RenderPage(eCtx, http.StatusOK, page.Search(userInfo, searchErrorMessage, searchQuery,
		searchEntries))
}

func (c *SearchController) handleHxNavSearch(eCtx echo.Context, ctx context.Context) error {
	// Create search filter
	searchFilter, err := c.parseSearchQueryString(getCurrentUserId(ctx), "")
	if err != nil {
		return err
	}

	// Create view model
	searchQuery, err := c.getSearchQueryViewData(ctx, false, searchFilter)
	if err != nil {
		return err
	}
	searchEntries := &vm.SearchEntries{}

	// Render
	return web.RenderHx(eCtx, http.StatusOK, hx.SearchNav(searchQuery, searchEntries))
}

func (c *SearchController) handleHxGetForm(eCtx echo.Context, isAdvanced bool) error {
	// Render
	return web.RenderHx(eCtx, http.StatusOK, hx.SearchForm(isAdvanced))
}

func (c *SearchController) handleHxGetFormActivities(eCtx echo.Context, ctx context.Context) error {
	entryTypeId, err := getTypeIdQueryParam(eCtx)
	if err != nil {
		return err
	}

	// Get entry master data
	entryActivities, err := c.getEntryActivities(ctx, entryTypeId)
	if err != nil {
		return err
	}

	// Create view model
	viewData := c.mapper.CreateEntryActivitiesViewModel(entryActivities)

	// Render
	return web.RenderHx(eCtx, http.StatusOK, hx.SearchFormActivityOptions(viewData))
}

func (c *SearchController) handleHxExecuteSearch(eCtx echo.Context, ctx context.Context,
	isAdvanced bool, searchInputs *searchInput) error {
	// Create search filter
	searchFilter, err := c.createSearchFilter(getCurrentUserId(ctx), searchInputs)
	if err != nil {
		searchErrorMessage := loc.GetErrorMessageString(getErrorCode(err))
		return web.RenderHx(eCtx, http.StatusOK, hx.Search(searchErrorMessage, &vm.SearchQuery{},
			&vm.SearchEntries{}))
	}

	// Create view model
	searchQuery, err := c.getSearchQueryViewData(ctx, isAdvanced, searchFilter)
	if err != nil {
		return err
	}
	searchEntries, err := c.getSearchEntriesViewData(ctx, searchFilter, 1)
	if err != nil {
		return err
	}

	// Push search URL into browser history
	url := "/search?"
	if isAdvanced {
		url = url + "adv=1&"
	}
	url = url + "query=" + searchEntries.Query
	web.HtmxPushUrl(eCtx, url)

	// Render
	return web.RenderHx(eCtx, http.StatusOK, hx.Search("", searchQuery, searchEntries))
}

func (c *SearchController) handleHxGetSearchPage(eCtx echo.Context, ctx context.Context,
	isAdvanced bool, query string, pageNum int) error {
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
	return web.RenderHx(eCtx, http.StatusOK, hx.SearchPage(isAdvanced, searchEntries))
}

func (c *SearchController) getSearchQueryViewData(ctx context.Context, isAdvanced bool,
	searchFilter *model.EntriesFilter) (*vm.SearchQuery, error) {
	// Create default values
	filter := searchFilter
	if filter == nil {
		filter = model.NewEntriesFilter()
	}
	typeId := filter.TypeId
	if !filter.ByType {
		typeId = model.EntryTypeIdWork
	}
	sTime := filter.StartTime
	eTime := filter.EndTime
	if !filter.ByTime {
		sTime = time.Now()
		eTime = time.Now()
	}

	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx, typeId)
	if err != nil {
		return nil, err
	}

	// Create view model
	return c.mapper.CreateSearchQueryViewModel(isAdvanced, filter.ByType, typeId, filter.ByTime,
		sTime, eTime, filter.ByActivity, filter.ActivityId, filter.Description, entryTypes,
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
	totPageNum := calculateNumberOfTotalPages(cnt, pageSize)
	return c.mapper.CreateSearchEntriesViewModel(query, pageNum, totPageNum, entries, entryTypesMap,
		entryActivitiesMap), nil
}

// --- Search query functions ---

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
		err := e.NewError(e.LogicEntrySearchDateIntervalInvalid, fmt.Sprintf("End date %s before "+
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
	if err = validateStringLength(input.text, model.MaxLengthEntryDescription,
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
