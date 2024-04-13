package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/constant"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/pkg/util"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	"kellnhofer.com/work-log/web/view/pages"
)

const searchDateTimeFormat = "200601021504"

type searchEntriesFormInput struct {
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
	return &SearchController{
		baseController: baseController{
			uServ: uServ,
			eServ: eServ,
		},
		mapper: mapper.NewSearchMapper(),
	}
}

// --- Endpoints ---

// GetSearchHandler returns a handler for "GET /search".
func (c *SearchController) GetSearchHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /search.")
		// Get search string
		searchQuery, avail := getSearchQueryParam(eCtx)
		if !avail {
			return c.handleShowSearch(eCtx)
		} else {
			return c.handleShowListSearch(eCtx, searchQuery)
		}
	}
}

// PostSearchHandler returns a handler for "POST /search".
func (c *SearchController) PostSearchHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /search.")
		return c.handleExecuteSearch(eCtx)
	}
}

// --- Handler functions ---

func (c *SearchController) handleShowSearch(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	entryTypeId := 0
	if len(entryTypes) > 0 {
		entryTypeId = entryTypes[0].Id
	}
	model := c.mapper.CreateInitialSearchViewModel(prevUrl, entryTypeId, time.Now(), 0, entryTypes,
		entryActivities)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.SearchPage(model))
}

func (c *SearchController) handleShowListSearch(eCtx echo.Context, query string) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user
	user, err := c.getUser(ctx)
	if err != nil {
		return err
	}

	// Get page number, offset and limit
	pageNum, _, err := getPageNumberQueryParam(eCtx)
	if err != nil {
		return err
	}
	offset, limit := calculateOffsetLimitFromPageNumber(pageNum)

	// Create entries filter from query string
	filter, err := c.parseSearchQueryString(user.Id, query)
	if err != nil {
		return err
	}
	// Create entries sort
	sort := model.NewEntriesSort()
	sort.ByTime = model.DescSorting

	// Get entries
	entries, cnt, err := c.eServ.GetDateEntries(ctx, filter, sort, offset, limit)
	if err != nil {
		return err
	}
	// Get entry master data
	entryTypesMap, entryActivitiesMap, err := c.getEntryMasterDataMap(ctx)
	if err != nil {
		return err
	}

	// Create view model
	userModel := c.mapper.CreateUserInfoViewModel(user)
	model := c.mapper.CreateSearchEntriesViewModel(constant.ViewPathDefault, query, pageNum,
		pageSize, cnt, entries, entryTypesMap, entryActivitiesMap)

	// Save current URL to be able to used later for back navigation
	saveCurrentUrl(eCtx)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.SearchEntriesPage(userModel, model))
}

func (c *SearchController) handleExecuteSearch(eCtx echo.Context) error {
	// Get form inputs
	input := c.getSearchEntriesFormInput(eCtx)

	// Create entries filter from inputs
	filter, err := c.createEntriesFilter(input)
	if err != nil {
		return c.handleSearchError(eCtx, err, input)
	}

	return c.handleSearchSuccess(eCtx, filter)
}

func (c *SearchController) handleSearchSuccess(eCtx echo.Context, filter *model.EntriesFilter) error {
	return eCtx.Redirect(http.StatusFound, "/search?query="+c.buildSearchQueryString(filter))
}

func (c *SearchController) handleSearchError(eCtx echo.Context, err error,
	input *searchEntriesFormInput) error {
	// Get context
	ctx := getContext(eCtx)

	// Get error message
	ec := getErrorCode(err)
	em := loc.GetErrorMessageString(ec)

	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	byEntryType, _ := strconv.ParseBool(input.byType)
	entryTypeId, _ := strconv.Atoi(input.typeId)
	byEntryDate, _ := strconv.ParseBool(input.byDate)
	byEntryActivity, _ := strconv.ParseBool(input.byActivity)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	byEntryDescription, _ := strconv.ParseBool(input.byDescription)
	model := c.mapper.CreateSearchViewModel(prevUrl, em, byEntryType, entryTypeId,
		byEntryDate, input.startDate, input.endDate, byEntryActivity, entryActivityId,
		byEntryDescription, input.description, entryTypes, entryActivities)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.SearchPage(model))
}

// --- Form input retrieval functions ---

func (c *SearchController) getSearchEntriesFormInput(eCtx echo.Context) *searchEntriesFormInput {
	i := searchEntriesFormInput{}
	i.byType = eCtx.FormValue("by-type")
	i.typeId = eCtx.FormValue("type")
	i.byDate = eCtx.FormValue("by-date")
	i.startDate = eCtx.FormValue("start-date")
	i.endDate = eCtx.FormValue("end-date")
	i.byActivity = eCtx.FormValue("by-activity")
	i.activityId = eCtx.FormValue("activity")
	i.byDescription = eCtx.FormValue("by-description")
	i.description = eCtx.FormValue("description")
	return &i
}

// --- Model converter functions ---

func (c *SearchController) createEntriesFilter(input *searchEntriesFormInput) (*model.EntriesFilter,
	error) {
	filter := model.NewEntriesFilter()

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

// --- Search query functions ---

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

	qps := strings.Split(query, "|")

	// Check if query is empty
	if len(qps) < 1 {
		err := e.NewError(e.ValSearchQueryInvalid, "Search query is empty.")
		log.Debug(err.StackTrace())
		return nil, err
	}

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
