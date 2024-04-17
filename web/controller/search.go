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
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/pkg/util"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/view/hx"
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

		searchQuery, _ := getSearchQueryParam(eCtx)
		pageNum, isPageReq, err := getPageNumberQueryParam(eCtx)
		if err != nil {
			return err
		}

		ctx := getContext(eCtx)

		if !isHtmxReq {
			return c.handleShowSearch(eCtx, ctx, searchQuery, pageNum)
		} else if !isPageReq {
			return c.handleHxShowSearch(eCtx, ctx, searchQuery, pageNum)
		} else {
			return c.handleHxGetSearchPage(eCtx, ctx, searchQuery, pageNum)
		}
	}
}

// PostSearchHandler returns a handler for "POST /search".
func (c *SearchController) PostSearchHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /search.")

		// TODO!!!
		return nil
	}
}

// --- Handler functions ---

func (c *SearchController) handleShowSearch(eCtx echo.Context, ctx context.Context, query string,
	pageNum int) error {
	// Create view model
	userModel, err := c.getUserInfoViewData(ctx)
	if err != nil {
		return err
	}
	model, err := c.getSearchViewData(ctx, query, pageNum)
	if err != nil {
		return err
	}

	// Render
	return web.Render(eCtx, http.StatusOK, pages.Search(userModel, model))
}

func (c *SearchController) handleHxShowSearch(eCtx echo.Context, ctx context.Context, query string,
	pageNum int) error {
	// Create view model
	model, err := c.getSearchViewData(ctx, query, pageNum)
	if err != nil {
		return err
	}

	// Render
	return web.Render(eCtx, http.StatusOK, hx.Search(model))
}

func (c *SearchController) handleHxGetSearchPage(eCtx echo.Context, ctx context.Context, query string,
	pageNum int) error {
	// TODO!!!
	return nil
}

func (c *SearchController) getSearchViewData(ctx context.Context, query string, pageNum int,
) (*vm.SearchEntries, error) {
	// Get current user information
	userId := getCurrentUserId(ctx)

	// Create entries filter from query string
	filter, err := c.parseSearchQueryString(userId, query)
	if err != nil {
		return nil, err
	}
	// Create entries sort
	sort := model.NewEntriesSort()
	sort.ByTime = model.DescSorting

	// Get entries
	offset, limit := calculateOffsetLimitFromPageNumber(pageNum)
	entries, cnt, err := c.eServ.GetDateEntries(ctx, filter, sort, offset, limit)
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
