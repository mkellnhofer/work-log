package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/constant"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/pkg/util"
	view "kellnhofer.com/work-log/web"
	vm "kellnhofer.com/work-log/web/model"
)

const pageSize = 7

const dateFormat = "2006-01-02"
const timeFormat = "15:04"
const dateTimeFormat = "2006-01-02 15:04"

const searchDateTimeFormat = "200601021504"

type dailyWorkingDuration struct {
	fromDate time.Time
	duration time.Duration
}

type monthlyVacationDays struct {
	fromDate time.Time
	days     float32
}

type entryFormInput struct {
	typeId      string
	date        string
	startTime   string
	endTime     string
	activityId  string
	description string
}

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

type overviewFormInput struct {
	month       string
	showDetails string
}

// EntryController handles requests for entry endpoints.
type EntryController struct {
	uServ *service.UserService
	eServ *service.EntryService
}

// NewEntryController creates a new entry controller.
func NewEntryController(uServ *service.UserService, eServ *service.EntryService) *EntryController {
	return &EntryController{uServ, eServ}
}

// --- Endpoints ---

// GetListHandler returns a handler for "GET /list".
func (c *EntryController) GetListHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /list.")
		return c.handleShowList(eCtx)
	}
}

// GetCreateHandler returns a handler for "GET /create".
func (c *EntryController) GetCreateHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /create.")
		return c.handleShowCreate(eCtx)
	}
}

// PostCreateHandler returns a handler for "POST /create".
func (c *EntryController) PostCreateHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /create.")
		return c.handleExecuteCreate(eCtx)
	}
}

// GetEditHandler returns a handler for "GET /edit/{id}".
func (c *EntryController) GetEditHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /edit/{id}.")
		return c.handleShowEdit(eCtx)
	}
}

// PostEditHandler returns a handler for "POST /edit/{id}".
func (c *EntryController) PostEditHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /edit/{id}.")
		return c.handleExecuteEdit(eCtx)
	}
}

// GetCopyHandler returns a handler for "GET /copy/{id}".
func (c *EntryController) GetCopyHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /copy/{id}.")
		return c.handleShowCopy(eCtx)
	}
}

// PostCopyHandler returns a handler for "POST /copy/{id}".
func (c *EntryController) PostCopyHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /copy/{id}.")
		return c.handleExecuteCopy(eCtx)
	}
}

// PostDeleteHandler returns a handler for "POST /delete/{id}".
func (c *EntryController) PostDeleteHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /delete/{id}.")
		return c.handleExecuteDelete(eCtx)
	}
}

// GetSearchHandler returns a handler for "GET /search".
func (c *EntryController) GetSearchHandler() echo.HandlerFunc {
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
func (c *EntryController) PostSearchHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /search.")
		return c.handleExecuteSearch(eCtx)
	}
}

// GetOverviewHandler returns a handler for "GET /overview".
func (c *EntryController) GetOverviewHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /overview.")
		return c.handleShowOverview(eCtx)
	}
}

// PostOverviewHandler returns a handler for "POST /overview".
func (c *EntryController) PostOverviewHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /overview.")
		return c.handleExecuteOverviewChange(eCtx)
	}
}

// GetOverviewExportHandler returns a handler for "GET /overview/export".
func (c *EntryController) GetOverviewExportHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /overview/export.")
		return c.handleExportOverview(eCtx)
	}
}

// --- List handler functions ---

func (c *EntryController) handleShowList(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID and user contract
	userId, userContract, err := c.getUserIdAndUserContract(ctx)
	if err != nil {
		return err
	}

	// Get page number, offset and limit
	pageNum, offset, limit, err := c.getListPagingParams(eCtx)
	if err != nil {
		return err
	}

	// Get work summary (only for first page)
	var workSummary *model.WorkSummary
	if pageNum == 1 {
		workSummary, err = c.eServ.GetTotalWorkSummaryByUserId(ctx, userId)
		if err != nil {
			return err
		}
	}

	// Get entries
	entries, cnt, err := c.eServ.GetDateEntriesByUserId(ctx, userId, offset, limit)
	if err != nil {
		return err
	}
	// Get entry master data
	entryTypesMap, entryActivitiesMap, err := c.getEntryMasterDataMap(ctx)
	if err != nil {
		return err
	}

	// Create view model
	model := c.createListViewModel(userContract, workSummary, pageNum, cnt, entries, entryTypesMap,
		entryActivitiesMap)

	// Save current URL to be able to used later for back navigation
	saveCurrentUrl(eCtx)

	// Render
	return view.RenderListEntriesTemplate(eCtx.Response(), model)
}

// --- Create handler functions ---

func (c *EntryController) handleShowCreate(eCtx echo.Context) error {
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
	model := c.createCreateViewModel(prevUrl, "", entryTypeId, getDateString(time.Now()), "00:00",
		"00:00", 0, "", entryTypes, entryActivities)

	// Render
	return view.RenderCreateEntryTemplate(eCtx.Response(), model)
}

func (c *EntryController) handleExecuteCreate(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get form inputs
	input := c.getEntryFormInput(eCtx)

	// Create model
	entry, err := c.createEntryModel(0, userId, input)
	if err != nil {
		return c.handleCreateError(eCtx, err, input)
	}

	// Create entry
	if err := c.eServ.CreateEntry(ctx, entry); err != nil {
		return c.handleCreateError(eCtx, err, input)
	}

	return c.handleCreateSuccess(eCtx)
}

func (c *EntryController) handleCreateSuccess(eCtx echo.Context) error {
	prevUrl := getPreviousUrl(eCtx)
	return eCtx.Redirect(http.StatusFound, prevUrl)
}

func (c *EntryController) handleCreateError(eCtx echo.Context, err error,
	input *entryFormInput) error {
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
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.createCreateViewModel(prevUrl, em, entryTypeId, input.date, input.startTime,
		input.endTime, entryActivityId, input.description, entryTypes, entryActivities)

	// Render
	return view.RenderCreateEntryTemplate(eCtx.Response(), model)
}

// --- Edit handler functions ---

func (c *EntryController) handleShowEdit(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Get entry
	entry, err := c.getEntry(ctx, entryId, userId)
	if err != nil {
		return err
	}
	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	model := c.createEditViewModel(prevUrl, "", entry.Id, entry.TypeId, getDateString(
		entry.StartTime), getTimeString(entry.StartTime), getTimeString(entry.EndTime),
		entry.ActivityId, entry.Description, entryTypes, entryActivities)

	// Render
	return view.RenderEditEntryTemplate(eCtx.Response(), model)
}

func (c *EntryController) handleExecuteEdit(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Get form inputs
	input := c.getEntryFormInput(eCtx)

	// Create model
	entry, err := c.createEntryModel(entryId, userId, input)
	if err != nil {
		return c.handleEditError(eCtx, err, entryId, input)
	}

	// Update entry
	if err := c.eServ.UpdateEntry(ctx, entry); err != nil {
		return c.handleEditError(eCtx, err, entryId, input)
	}

	return c.handleEditSuccess(eCtx)
}

func (c *EntryController) handleEditSuccess(eCtx echo.Context) error {
	prevUrl := getPreviousUrl(eCtx)
	return eCtx.Redirect(http.StatusFound, prevUrl)
}

func (c *EntryController) handleEditError(eCtx echo.Context, err error, id int,
	input *entryFormInput) error {
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
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.createEditViewModel(prevUrl, em, id, entryTypeId, input.date, input.startTime,
		input.endTime, entryActivityId, input.description, entryTypes, entryActivities)

	// Render
	return view.RenderEditEntryTemplate(eCtx.Response(), model)
}

// --- Copy handler functions ---

func (c *EntryController) handleShowCopy(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Get entry
	entry, err := c.getEntry(ctx, entryId, userId)
	if err != nil {
		return err
	}
	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	model := c.createCopyViewModel(prevUrl, "", entry.Id, entry.TypeId, getDateString(
		entry.StartTime), getTimeString(entry.StartTime), getTimeString(entry.EndTime),
		entry.ActivityId, entry.Description, entryTypes, entryActivities)

	// Render
	return view.RenderCopyEntryTemplate(eCtx.Response(), model)
}

func (c *EntryController) handleExecuteCopy(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Get form inputs
	input := c.getEntryFormInput(eCtx)

	// Create model
	entry, err := c.createEntryModel(0, userId, input)
	if err != nil {
		return c.handleCopyError(eCtx, err, entryId, input)
	}

	// Create entry
	if err := c.eServ.CreateEntry(ctx, entry); err != nil {
		return c.handleCopyError(eCtx, err, entryId, input)
	}

	return c.handleCopySuccess(eCtx)
}

func (c *EntryController) handleCopySuccess(eCtx echo.Context) error {
	prevUrl := getPreviousUrl(eCtx)
	return eCtx.Redirect(http.StatusFound, prevUrl)
}

func (c *EntryController) handleCopyError(eCtx echo.Context, err error, id int,
	input *entryFormInput) error {
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
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.createCopyViewModel(prevUrl, em, id, entryTypeId, input.date, input.startTime,
		input.endTime, entryActivityId, input.description, entryTypes, entryActivities)

	// Render
	return view.RenderCopyEntryTemplate(eCtx.Response(), model)
}

// --- Delete handler functions ---

func (c *EntryController) handleExecuteDelete(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Delete entry
	if err := c.eServ.DeleteEntryByIdAndUserId(ctx, entryId, userId); err != nil {
		return err
	}

	return c.handleDeleteSuccess(eCtx)
}

func (c *EntryController) handleDeleteSuccess(eCtx echo.Context) error {
	prevUrl := getPreviousUrl(eCtx)
	return eCtx.Redirect(http.StatusFound, prevUrl)
}

// --- Search handler functions ---

func (c *EntryController) handleShowSearch(eCtx echo.Context) error {
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
	model := c.createSearchViewModel(prevUrl, "", false, entryTypeId, false, getDateString(time.Now()),
		getDateString(time.Now()), false, 0, false, "", entryTypes, entryActivities)

	// Render
	return view.RenderSearchEntriesTemplate(eCtx.Response(), model)
}

func (c *EntryController) handleShowListSearch(eCtx echo.Context, query string) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get page number, offset and limit
	pageNum, offset, limit, err := c.getListPagingParams(eCtx)
	if err != nil {
		return err
	}

	// Create entries filter from query string
	filter, err := c.parseSearchQueryString(userId, query)
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
	model := c.createListSearchViewModel(constant.ViewPathDefault, query, pageNum, cnt, entries,
		entryTypesMap, entryActivitiesMap)

	// Save current URL to be able to used later for back navigation
	saveCurrentUrl(eCtx)

	// Render
	return view.RenderListSearchEntriesTemplate(eCtx.Response(), model)
}

func (c *EntryController) handleExecuteSearch(eCtx echo.Context) error {
	// Get form inputs
	input := c.getSearchEntriesFormInput(eCtx)

	// Create entries filter from inputs
	filter, err := c.createEntriesFilter(input)
	if err != nil {
		return c.handleSearchError(eCtx, err, input)
	}

	return c.handleSearchSuccess(eCtx, filter)
}

func (c *EntryController) handleSearchSuccess(eCtx echo.Context, filter *model.EntriesFilter) error {
	return eCtx.Redirect(http.StatusFound, "/search?query="+c.buildSearchQueryString(filter))
}

func (c *EntryController) handleSearchError(eCtx echo.Context, err error,
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
	model := c.createSearchViewModel(prevUrl, em, byEntryType, entryTypeId, byEntryDate,
		input.startDate, input.endDate, byEntryActivity, entryActivityId, byEntryDescription,
		input.description, entryTypes, entryActivities)

	// Render
	return view.RenderSearchEntriesTemplate(eCtx.Response(), model)
}

// --- Overview handler functions ---

func (c *EntryController) handleShowOverview(eCtx echo.Context) error {
	// Get view data
	model, err := c.getOverviewViewData(eCtx)
	if err != nil {
		return err
	}

	// Render
	return view.RenderListOverviewEntriesTemplate(eCtx.Response(), model)
}

func (c *EntryController) handleExportOverview(eCtx echo.Context) error {
	// Get view data
	model, err := c.getOverviewViewData(eCtx)
	if err != nil {
		return err
	}

	// Create file
	fileName := fmt.Sprintf("work-log-export-%s.xlsx", model.CurrMonth)
	file := exportOverviewEntries(model)

	// Write file
	return writeFile(eCtx.Response(), fileName, file)
}

func (c *EntryController) getOverviewViewData(eCtx echo.Context) (*vm.ListOverviewEntries, error) {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID and user contract
	userId, userContract, err := c.getUserIdAndUserContract(ctx)
	if err != nil {
		return nil, err
	}

	// Get user setting
	showDetails, err := c.uServ.GetSettingShowOverviewDetails(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Get year and month
	year, month, err := c.getOverviewParams(eCtx)
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
	prevUrl := getPreviousUrl(eCtx)
	model := c.createListOverviewViewModel(prevUrl, year, month, userContract, entries,
		entryTypesMap, entryActivitiesMap, showDetails)

	return model, nil
}

func (c *EntryController) getOverviewParams(eCtx echo.Context) (int, int, error) {
	// Get year and month
	y, m, avail, err := getMonthQueryParam(eCtx)
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

func (c *EntryController) handleExecuteOverviewChange(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get form inputs
	input := c.getOverviewFormInput(eCtx)

	// Validate month param
	_, _, _, err := parseMonthParam(input.month)
	if err != nil {
		return err
	}

	// Update user setting
	showDetails := input.showDetails == "on"
	c.uServ.UpdateSettingShowOverviewDetails(ctx, userId, showDetails)

	// Redirect
	return eCtx.Redirect(http.StatusFound, "/overview?month="+input.month)
}

// --- Viem model converter functions ---

func (c *EntryController) createListViewModel(userContract *model.Contract,
	workSummary *model.WorkSummary, pageNum int, cnt int, entries []*model.Entry,
	entryTypesMap map[int]*model.EntryType, entryActivitiesMap map[int]*model.EntryActivity) *vm.
	ListEntries {
	lesvm := vm.NewListEntries()

	// Calculate summary
	lesvm.Summary = c.createListSummaryViewModel(userContract, workSummary)

	// Calculate previous/next page numbers
	lesvm.HasPrevPage = pageNum > 1
	lesvm.HasNextPage = (pageNum * pageSize) < cnt
	lesvm.PrevPageNum = pageNum - 1
	lesvm.NextPageNum = pageNum + 1

	// Create entries
	lesvm.Days = c.createEntriesViewModel(userContract, entries, entryTypesMap, entryActivitiesMap,
		true)

	return lesvm
}

func (c *EntryController) createListSummaryViewModel(userContract *model.Contract,
	workSummary *model.WorkSummary) *vm.ListEntriesSummary {
	// If no user contract or work summary was provided: Skip calculation
	if userContract == nil || workSummary == nil {
		return nil
	}

	// Calulate durations
	overtimeHours := c.calculateOvertimeHours(userContract, workSummary)
	remainingVacationDays := c.calculateRemainingVacationDays(userContract, workSummary)

	// Create summary
	lessvm := vm.NewListEntriesSummary()
	lessvm.OvertimeHours = view.CreateHoursString(overtimeHours)
	lessvm.RemainingVacationDays = view.CreateDaysString(remainingVacationDays)
	return lessvm
}

func (c *EntryController) calculateOvertimeHours(userContract *model.Contract,
	workSummary *model.WorkSummary) float32 {
	// Calculate initial overtime duration
	initOvertimeDuration := time.Duration(int(userContract.InitOvertimeHours*60.0)) * time.Minute

	// Calculate actual duration
	var actualWorkDuration time.Duration
	for _, workDuration := range workSummary.WorkDurations {
		actualWorkDuration = actualWorkDuration + workDuration.WorkDuration
	}
	log.Verbf("Actual work duration: %.0f min", actualWorkDuration.Minutes())

	// Get target working durations
	targetWorkDurations := c.convertWorkingHours(userContract.WorkingHours)
	// Abort if no target working durations were set
	if len(targetWorkDurations) == 0 {
		return 0.0
	}

	// Calculate target duration
	start := userContract.FirstDay
	end := time.Now()
	targetWorkDuration := time.Duration(0)
	for i := 0; i < len(targetWorkDurations); i++ {
		// Calculate interval start/end
		intStart := start
		if i > 0 {
			intStart = targetWorkDurations[i].fromDate
		}
		intEnd := end
		if i+1 < len(targetWorkDurations) {
			intEnd = targetWorkDurations[i+1].fromDate.AddDate(0, 0, -1)
		}

		// Calculate interval work days
		intWorkDays := util.CalculateWorkingDays(intStart, intEnd)
		log.Verbf("Interval work days: %s - %s: %d", getDateString(intStart), getDateString(intEnd),
			intWorkDays)

		// Calculate interval target duration
		intTargetWorkDuration := time.Duration(intWorkDays) * targetWorkDurations[i].duration
		log.Verbf("Interval daily work duration: %s - %s: %.0f min", getDateString(intStart),
			getDateString(intEnd), targetWorkDurations[i].duration.Minutes())
		log.Verbf("Interval target work duration: %s - %s: %.0f min", getDateString(intStart),
			getDateString(intEnd), intTargetWorkDuration.Minutes())

		// Update target duration
		targetWorkDuration = targetWorkDuration + intTargetWorkDuration
	}

	// Calculate overtime
	overtimeDuration := initOvertimeDuration + actualWorkDuration - targetWorkDuration

	// Return rounded hours
	return getRoundedHours(overtimeDuration)
}

func (c *EntryController) calculateRemainingVacationDays(userContract *model.Contract,
	workSummary *model.WorkSummary) float32 {
	// Get monthly vacation days
	vacationDays := c.convertVacationDays(userContract.VacationDays)
	// Get daily working durations
	workDurations := c.convertWorkingHours(userContract.WorkingHours)
	// Abort if no vacation days or working durations were set
	if len(vacationDays) == 0 || len(workDurations) == 0 {
		return 0.0
	}

	// Calculate initial vacation hours
	initVacationHours := userContract.InitVacationDays * float32(workDurations[0].duration.Hours())
	log.Verbf("Initial vacation: %.0f hours", initVacationHours)

	// Calculate available vacation hours (month by month)
	start := userContract.FirstDay
	now := time.Now()
	curMonth := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.Local)
	endMonth := time.Date(now.Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
	availableVacationHours := float32(0.0)
	for curMonth.Before(endMonth) {
		// Get vacation days and working duration for current month
		vd := c.findVacationDaysForDate(vacationDays, curMonth)
		wd := c.findWorkingDurationForDate(workDurations, curMonth)
		// Calculate vacation hours for current month
		availableVacationHours = availableVacationHours + vd*float32(wd.Hours())
		// Calculate next month
		curMonth = curMonth.AddDate(0, 1, 0)
	}
	log.Verbf("Available vacation: %.0f hours", availableVacationHours)

	// Calculate taken vacation hours
	takenVacationHours := float32(0.0)
	for _, workDuration := range workSummary.WorkDurations {
		if workDuration.TypeId == model.EntryTypeIdVacation {
			takenVacationHours = takenVacationHours + float32(workDuration.WorkDuration.Hours())
		}
	}
	log.Verbf("Taken vacation: %.0f hours", takenVacationHours)

	// Calculate remaining vacation hours
	remainingVacationHours := initVacationHours + availableVacationHours - takenVacationHours
	log.Verbf("Remaining vacation: %.0f hours", remainingVacationHours)

	// Convert vacation hours to vacation days
	wd := c.findWorkingDurationForDate(workDurations, now)
	var remainingVacationDays float32
	if remainingVacationHours > 0 {
		remainingVacationDays = remainingVacationHours / float32(wd.Hours())
	}

	return remainingVacationDays
}

func (c *EntryController) getEntryTypeDescription(entryTypesMap map[int]*model.EntryType,
	id int) string {
	et, ok := entryTypesMap[id]
	if ok {
		return et.Description
	}
	return ""
}

func (c *EntryController) getEntryActivityDescription(entryActivitiesMap map[int]*model.EntryActivity,
	id int) string {
	ea, ok := entryActivitiesMap[id]
	if ok {
		return ea.Description
	}
	return ""
}

func (c *EntryController) createCreateViewModel(prevUrl string, errorMessage string, typeId int,
	date string, startTime string, endTime string, activityId int, description string,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.CreateEntry {
	cevm := vm.NewCreateEntry()
	cevm.PreviousUrl = prevUrl
	cevm.ErrorMessage = errorMessage
	cevm.Entry = c.createEntryViewModel(0, typeId, date, startTime, endTime, activityId,
		description)
	cevm.EntryTypes = c.createEntryTypesViewModel(types)
	cevm.EntryActivities = c.createEntryActivitiesViewModel(activities)
	return cevm
}

func (c *EntryController) createEditViewModel(prevUrl string, errorMessage string, id int,
	typeId int, date string, startTime string, endTime string, activityId int, description string,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.EditEntry {
	eevm := vm.NewEditEntry()
	eevm.PreviousUrl = prevUrl
	eevm.ErrorMessage = errorMessage
	eevm.Entry = c.createEntryViewModel(id, typeId, date, startTime, endTime, activityId,
		description)
	eevm.EntryTypes = c.createEntryTypesViewModel(types)
	eevm.EntryActivities = c.createEntryActivitiesViewModel(activities)
	return eevm
}

func (c *EntryController) createCopyViewModel(prevUrl string, errorMessage string, id int,
	typeId int, date string, startTime string, endTime string, activityId int, description string,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.CopyEntry {
	cevm := vm.NewCopyEntry()
	cevm.PreviousUrl = prevUrl
	cevm.ErrorMessage = errorMessage
	cevm.Entry = c.createEntryViewModel(id, typeId, date, startTime, endTime, activityId,
		description)
	cevm.EntryTypes = c.createEntryTypesViewModel(types)
	cevm.EntryActivities = c.createEntryActivitiesViewModel(activities)
	return cevm
}

func (c *EntryController) createSearchViewModel(prevUrl string, errorMessage string, byType bool,
	typeId int, byDate bool, startDate string, endDate string, byActivity bool, activityId int,
	byDescription bool, description string, types []*model.EntryType,
	activities []*model.EntryActivity) *vm.SearchEntries {
	sevm := vm.NewSearchEntries()
	sevm.PreviousUrl = prevUrl
	sevm.ErrorMessage = errorMessage
	sevm.ByType = byType
	sevm.TypeId = typeId
	sevm.ByDate = byDate
	sevm.StartDate = startDate
	sevm.EndDate = endDate
	sevm.ByActivity = byActivity
	sevm.ActivityId = activityId
	sevm.ByDescription = byDescription
	sevm.Description = description
	sevm.EntryTypes = c.createEntryTypesViewModel(types)
	sevm.EntryActivities = c.createEntryActivitiesViewModel(activities)
	return sevm
}

func (c *EntryController) createListSearchViewModel(prevUrl string, query string, pageNum int,
	cnt int, entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity) *vm.ListSearchEntries {
	lesvm := vm.NewListSearchEntries()
	lesvm.PreviousUrl = prevUrl
	lesvm.Query = query

	// Calculate previous/next page numbers
	lesvm.HasPrevPage = pageNum > 1
	lesvm.HasNextPage = (pageNum * pageSize) < cnt
	lesvm.PrevPageNum = pageNum - 1
	lesvm.NextPageNum = pageNum + 1

	// Create entries
	lesvm.Days = c.createEntriesViewModel(nil, entries, entryTypesMap, entryActivitiesMap, false)

	return lesvm
}

func (c *EntryController) createEntriesViewModel(userContract *model.Contract,
	entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity,
	checkMissingOrOverlapping bool) []*vm.ListEntriesDay {
	ldsvm := make([]*vm.ListEntriesDay, 0, pageSize)

	var calcTargetWorkDurationReached bool
	var targetWorkDurations []dailyWorkingDuration
	targetWorkDuration := time.Duration(0)

	// If no user contract was provided: Skip target calculation
	if userContract != nil {
		calcTargetWorkDurationReached = true
		targetWorkDurations = c.convertWorkingHours(userContract.WorkingHours)
	}

	var ldvm *vm.ListEntriesDay
	prevDate := ""
	var prevStartTime *time.Time
	var totalWorkDuration time.Duration
	var totalBreakDuration time.Duration
	var wasTargetWorkDurationReached string

	// Create entries
	for _, entry := range entries {
		currDate := getDateString(entry.StartTime)

		// If new day: Create and add new day
		if prevDate != currDate {
			prevDate = currDate
			prevStartTime = nil

			// Reset total work and break duration
			totalWorkDuration = 0
			totalBreakDuration = 0
			wasTargetWorkDurationReached = ""

			// Get target work duration
			if calcTargetWorkDurationReached {
				targetWorkDuration = c.findWorkingDurationForDate(targetWorkDurations, entry.StartTime)
			}

			// Create and add new day
			ldvm = vm.NewListEntriesDay()
			ldvm.Date = view.FormatDate(entry.StartTime)
			ldvm.Weekday = view.GetWeekdayName(entry.StartTime)
			ldvm.Entries = make([]*vm.ListEntry, 0, 10)
			ldsvm = append(ldsvm, ldvm)
		}

		// Calculate work duration
		duration := entry.EndTime.Sub(entry.StartTime)
		totalWorkDuration = totalWorkDuration + duration

		// Calculate if target work duration was reached
		if calcTargetWorkDurationReached {
			reached := (totalWorkDuration - targetWorkDuration) >= 0
			wasTargetWorkDurationReached = strconv.FormatBool(reached)
		}

		// Calculate break duration
		if prevStartTime != nil && prevStartTime.After(entry.EndTime) {
			breakDuration := prevStartTime.Sub(entry.EndTime)
			totalBreakDuration = totalBreakDuration + breakDuration
		}

		// Check for missing or overlapping entry
		if checkMissingOrOverlapping {
			if prevStartTime != nil && prevStartTime.After(entry.EndTime) {
				levm := vm.NewListEntry()
				levm.IsMissing = true
				ldvm.Entries = append(ldvm.Entries, levm)
			} else if prevStartTime != nil && prevStartTime.Before(entry.EndTime) {
				levm := vm.NewListEntry()
				levm.IsOverlapping = true
				ldvm.Entries = append(ldvm.Entries, levm)
			}
		}
		prevStartTime = &entry.StartTime

		// Create and add new entry
		levm := vm.NewListEntry()
		levm.Id = entry.Id
		levm.EntryType = c.getEntryTypeDescription(entryTypesMap, entry.TypeId)
		levm.StartTime = view.FormatTime(entry.StartTime)
		levm.EndTime = view.FormatTime(entry.EndTime)
		levm.Duration = view.FormatHours(duration)
		levm.EntryActivity = c.getEntryActivityDescription(entryActivitiesMap, entry.ActivityId)
		levm.Description = entry.Description
		ldvm.Entries = append(ldvm.Entries, levm)
		ldvm.WorkDuration = view.FormatHours(totalWorkDuration)
		ldvm.BreakDuration = view.FormatHours(totalBreakDuration)
		ldvm.WasTargetWorkDurationReached = wasTargetWorkDurationReached
	}

	return ldsvm
}

func (c *EntryController) createEntryViewModel(id int, typeId int, date string, startTime string,
	endTime string, activityId int, description string) *vm.Entry {
	evm := vm.NewEntry()
	evm.Id = id
	evm.TypeId = typeId
	evm.Date = date
	evm.StartTime = startTime
	evm.EndTime = endTime
	evm.ActivityId = activityId
	evm.Description = description
	return evm
}

func (c *EntryController) createEntryTypesViewModel(entryTypes []*model.EntryType) []*vm.EntryType {
	etsvm := make([]*vm.EntryType, 0, 10)
	for _, entryType := range entryTypes {
		etsvm = append(etsvm, c.createEntryTypeViewModel(entryType.Id, entryType.Description))
	}
	return etsvm
}

func (c *EntryController) createEntryTypeViewModel(id int, description string) *vm.EntryType {
	return vm.NewEntryType(id, description)
}

func (c *EntryController) createEntryActivitiesViewModel(entryActivities []*model.EntryActivity,
) []*vm.EntryActivity {
	easvm := make([]*vm.EntryActivity, 0, 10)
	easvm = append(easvm, c.createEntryActivityViewModel(0, "-"))
	for _, entryActivity := range entryActivities {
		easvm = append(easvm, c.createEntryActivityViewModel(entryActivity.Id,
			entryActivity.Description))
	}
	return easvm
}

func (c *EntryController) createEntryActivityViewModel(id int, description string) *vm.EntryActivity {
	return vm.NewEntryActivity(id, description)
}

func (c *EntryController) createListOverviewViewModel(prevUrl string, year int, month int,
	userContract *model.Contract, entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity, showDetails bool) *vm.ListOverviewEntries {
	lesvm := vm.NewListOverviewEntries()
	lesvm.PreviousUrl = prevUrl

	// Get current month name
	lesvm.CurrMonthName = fmt.Sprintf("%s %d", view.GetMonthName(month), year)

	// Calculate previous/next month
	var py, pm, ny, nm int
	if month == 1 {
		py = year - 1
		pm = 12
		ny = year
		nm = month + 1
	} else if month == 12 {
		py = year
		pm = month - 1
		ny = year + 1
		nm = 1
	} else {
		py = year
		pm = month - 1
		ny = year
		nm = month + 1
	}
	lesvm.CurrMonth = fmt.Sprintf("%d%02d", year, month)
	lesvm.PrevMonth = fmt.Sprintf("%d%02d", py, pm)
	lesvm.NextMonth = fmt.Sprintf("%d%02d", ny, nm)

	// Calculate summary
	lesvm.Summary = c.createOverviewSummaryViewModel(year, month, userContract, entries)

	// Create entries
	lesvm.ShowDetails = showDetails
	lesvm.Days = c.createOverviewEntriesViewModel(year, month, entries, entryTypesMap,
		entryActivitiesMap, showDetails)

	return lesvm
}

func (c *EntryController) createOverviewSummaryViewModel(year int, month int,
	userContract *model.Contract, entries []*model.Entry) *vm.ListOverviewEntriesSummary {

	// Calculate type durations
	var actWork, actTrav, actVaca, actHoli, actIlln time.Duration
	for _, entry := range entries {
		duration := entry.EndTime.Sub(entry.StartTime)

		switch entry.TypeId {
		case model.EntryTypeIdWork:
			actWork = actWork + duration
		case model.EntryTypeIdTravel:
			actTrav = actTrav + duration
		case model.EntryTypeIdVacation:
			actVaca = actVaca + duration
		case model.EntryTypeIdHoliday:
			actHoli = actHoli + duration
		case model.EntryTypeIdIllness:
			actIlln = actIlln + duration
		}
	}

	// Calculate days
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	workDays := util.CalculateWorkingDays(start, end)

	// Get target working durations
	targetWorkDurations := c.convertWorkingHours(userContract.WorkingHours)

	// Calculate target, actual and balance durations
	var tar time.Duration = time.Duration(workDays) * c.findWorkingDurationForDate(
		targetWorkDurations, start)
	var act time.Duration = actWork + actTrav + actVaca + actHoli + actIlln
	var bal time.Duration = act - tar

	// Create summary
	lessvm := vm.NewListOverviewEntriesSummary()
	lessvm.ActualWorkHours = createRoundedHoursString(actWork)
	lessvm.ActualTravelHours = createRoundedHoursString(actTrav)
	lessvm.ActualVacationHours = createRoundedHoursString(actVaca)
	lessvm.ActualHolidayHours = createRoundedHoursString(actHoli)
	lessvm.ActualIllnessHours = createRoundedHoursString(actIlln)
	lessvm.TargetHours = createRoundedHoursString(tar)
	lessvm.ActualHours = createRoundedHoursString(act)
	lessvm.BalanceHours = createRoundedHoursString(bal)
	return lessvm
}

func createRoundedHoursString(d time.Duration) string {
	return view.CreateHoursString(getRoundedHours(d))
}

func (c *EntryController) createOverviewEntriesViewModel(year int, month int, entries []*model.Entry,
	entryTypesMap map[int]*model.EntryType, entryActivitiesMap map[int]*model.EntryActivity,
	showDetails bool) []*vm.ListOverviewEntriesDay {
	ldsvm := make([]*vm.ListOverviewEntriesDay, 0, 31)

	curDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	// Create days
	entryIndex := 0
	for {
		// Create and add new day
		ldvm := vm.NewListOverviewEntriesDay()
		ldvm.Date = view.FormatShortDate(curDate)
		ldvm.Weekday = view.GetShortWeekdayName(curDate)
		ldvm.IsWeekendDay = curDate.Weekday() == time.Saturday || curDate.Weekday() == time.Sunday
		ldvm.Entries = make([]*vm.ListOverviewEntry, 0, 10)
		ldsvm = append(ldsvm, ldvm)

		// Create entries
		var colWorkDuration time.Duration
		var dailyWorkDuration time.Duration
		preEntryTypeId := 0
		var levm *vm.ListOverviewEntry
		for {
			// If there are no entries: Abort (No entries exist for this day)
			if len(entries) == 0 || len(entries) == entryIndex {
				break
			}
			// Get entry
			entry := entries[entryIndex]
			entryDate := entry.StartTime
			// If entry date does not match: Abort (All enties have been added for this day)
			_, _, cd := curDate.Date()
			_, _, ed := entryDate.Date()
			if cd != ed {
				colWorkDuration = 0
				preEntryTypeId = 0
				break
			}

			// Reset collected work duration
			if entry.TypeId != preEntryTypeId {
				colWorkDuration = 0
			}

			// Calculate work duration
			duration := entry.EndTime.Sub(entry.StartTime)
			colWorkDuration = colWorkDuration + duration
			dailyWorkDuration = dailyWorkDuration + duration

			// Create and add new entry
			if showDetails {
				levm = vm.NewListOverviewEntry()
				levm.Id = entry.Id
				levm.EntryType = c.getEntryTypeDescription(entryTypesMap, entry.TypeId)
				levm.StartTime = view.FormatTime(entry.StartTime)
				levm.EndTime = view.FormatTime(entry.EndTime)
				levm.Duration = view.FormatHours(duration)
				levm.EntryActivity = c.getEntryActivityDescription(entryActivitiesMap, entry.ActivityId)
				levm.Description = entry.Description
				ldvm.Entries = append(ldvm.Entries, levm)
			} else {
				if entry.TypeId != preEntryTypeId {
					levm = vm.NewListOverviewEntry()
					levm.Id = entry.Id
					levm.EntryType = c.getEntryTypeDescription(entryTypesMap, entry.TypeId)
					levm.StartTime = view.FormatTime(entry.StartTime)
					ldvm.Entries = append(ldvm.Entries, levm)
				}
				levm.EndTime = view.FormatTime(entry.EndTime)
				levm.Duration = view.FormatHours(colWorkDuration)
			}

			// Update previous entry type ID
			preEntryTypeId = entry.TypeId

			// Update entry index
			entryIndex++
		}
		ldvm.WorkDuration = view.FormatHours(dailyWorkDuration)

		// If next month is reached: Abort
		curDate = curDate.Add(24 * time.Hour)
		if curDate.Month() != time.Month(month) {
			break
		}
	}

	return ldsvm
}

// --- Paging functions ---

func (c *EntryController) getListPagingParams(eCtx echo.Context) (int, int, int, error) {
	// Get page number
	pageNum, avail, err := getPageNumberQueryParam(eCtx)
	if err != nil {
		return 0, 0, 0, err
	}
	if !avail {
		pageNum = 1
	}

	// Calculate offset and limit
	offset := (pageNum - 1) * pageSize
	limit := pageSize

	return pageNum, offset, limit, nil
}

// --- Form input retrieval functions ---

func (c *EntryController) getEntryFormInput(eCtx echo.Context) *entryFormInput {
	i := entryFormInput{}
	i.typeId = eCtx.FormValue("type")
	i.date = eCtx.FormValue("date")
	i.startTime = eCtx.FormValue("start-time")
	i.endTime = eCtx.FormValue("end-time")
	i.activityId = eCtx.FormValue("activity")
	i.description = eCtx.FormValue("description")
	return &i
}

func (c *EntryController) getSearchEntriesFormInput(eCtx echo.Context) *searchEntriesFormInput {
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

func (c *EntryController) getOverviewFormInput(eCtx echo.Context) *overviewFormInput {
	i := overviewFormInput{}
	i.month = eCtx.FormValue("month")
	i.showDetails = eCtx.FormValue("show-details")
	return &i
}

// --- Model converter functions ---

func (c *EntryController) createEntryModel(id int, userId int, input *entryFormInput) (
	*model.Entry, error) {
	entry := model.NewEntry()
	entry.Id = id
	entry.UserId = userId

	var err error

	// Convert type ID
	entry.TypeId, err = c.convertId(input.typeId, false)
	if err != nil {
		return nil, err
	}

	// Convert start/end time
	if _, err := c.convertDateTime(input.date, "00:00", e.ValDateInvalid); err != nil {
		return nil, err
	}
	entry.StartTime, err = c.convertDateTime(input.date, input.startTime, e.ValStartTimeInvalid)
	if err != nil {
		return nil, err
	}
	entry.EndTime, err = c.convertDateTime(input.date, input.endTime, e.ValEndTimeInvalid)
	if err != nil {
		return nil, err
	}

	// Convert activity ID
	entry.ActivityId, err = c.convertId(input.activityId, true)
	if err != nil {
		return nil, err
	}

	// Validate description
	if err = c.validateString(input.description, model.MaxLengthEntryDescription,
		e.ValDescriptionTooLong); err != nil {
		return nil, err
	}
	entry.Description = input.description

	return entry, nil
}

func (c *EntryController) createEntriesFilter(input *searchEntriesFormInput) (*model.EntriesFilter,
	error) {
	filter := model.NewEntriesFilter()

	var err error

	// Convert type ID
	filter.ByType = input.byType == "on"
	filter.TypeId, err = c.convertId(input.typeId, false)
	if err != nil {
		return nil, err
	}

	// Convert start/end time
	filter.ByTime = input.byDate == "on"
	filter.StartTime, err = c.convertDateTime(input.startDate, "00:00", e.ValStartDateInvalid)
	if err != nil {
		return nil, err
	}
	filter.EndTime, err = c.convertDateTime(input.endDate, "23:59", e.ValEndDateInvalid)
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
	filter.ActivityId, err = c.convertId(input.activityId, true)
	if err != nil {
		return nil, err
	}

	// Validate description
	filter.ByDescription = input.byDescription == "on"
	if err = c.validateString(input.description, model.MaxLengthEntryDescription,
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

func (c *EntryController) convertId(in string, allowZero bool) (int, error) {
	out, cErr := strconv.Atoi(in)
	if cErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID. (ID must be numeric.)", cErr)
		log.Debug(err.StackTrace())
		return 0, err
	}
	if !allowZero && out <= 0 {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (ID must be positive.)")
		log.Debug(err.StackTrace())
		return 0, err
	}
	if out < 0 {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (ID must be zero or positive.)")
		log.Debug(err.StackTrace())
		return 0, err
	}
	return out, nil
}

func (c *EntryController) convertDateTime(inDate string, inTime string, code int) (time.Time,
	error) {
	dt := inDate + " " + inTime
	out, pErr := time.ParseInLocation(dateTimeFormat, dt, time.Local)
	if pErr != nil {
		err := e.WrapError(code, fmt.Sprintf("Could not parse time %s.", inTime), pErr)
		log.Debug(err.StackTrace())
		return time.Now(), err
	}
	return out, nil
}

func (c *EntryController) validateString(in string, length int, code int) error {
	if len(in) > length {
		err := e.NewError(code, fmt.Sprintf("String too long. (Must be "+
			"<= %d characters.)", length))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

// --- Search query functions ---

func (c *EntryController) buildSearchQueryString(filter *model.EntriesFilter) string {
	var qps []string
	// Add parameter/value for entry type
	if filter.ByType {
		qps = append(qps, fmt.Sprintf("typ:%d", filter.TypeId))
	}
	// Add parameter/value for entry start/end time
	if filter.ByTime {
		qps = append(qps, fmt.Sprintf("tim:%s", formatSearchDateRange(filter.StartTime,
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

func (c *EntryController) parseSearchQueryString(userId int, query string) (*model.EntriesFilter,
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
			filter.StartTime, filter.EndTime, cErr = parseSearchDateRange(v)
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

func formatSearchDateRange(startDate time.Time, endDate time.Time) string {
	return fmt.Sprintf("%s-%s", formatSearchDate(startDate), formatSearchDate(endDate))
}

func formatSearchDate(date time.Time) string {
	return date.Format(searchDateTimeFormat)
}

func parseSearchDateRange(dateRange string) (time.Time, time.Time, error) {
	se := strings.Split(dateRange, "-")
	if len(se) < 2 {
		return time.Time{}, time.Time{}, errors.New("invalid range")
	}
	startTime, err := parseSearchDate(se[0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endTime, err := parseSearchDate(se[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return startTime, endTime, nil
}

func parseSearchDate(date string) (time.Time, error) {
	return time.ParseInLocation(searchDateTimeFormat, date, time.Local)
}

// --- Export functions ---

func exportOverviewEntries(overviewEntries *vm.ListOverviewEntries) *excelize.File {
	f := excelize.NewFile()

	// Configure work book
	now := time.Now()
	f.SetDocProps(&excelize.DocProperties{
		Created:        now.Format(time.RFC3339),
		Creator:        loc.CreateString("appName"),
		Modified:       now.Format(time.RFC3339),
		LastModifiedBy: loc.CreateString("appName"),
		Title: loc.CreateString("exportPropTitle", loc.CreateString("appName"),
			overviewEntries.CurrMonthName),
		Description: loc.CreateString("exportPropDescription", loc.CreateString("appName")),
		Language:    loc.LngTag.String(),
	})

	sheet := "Sheet1"

	// Create default style
	styleDefault, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "left", WrapText: true},
		Font:      &excelize.Font{Size: 10},
	})

	// Create text styles
	styleTitle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 14, Bold: true}})
	styleTextBold, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 10, Bold: true}})

	// Creat tables styles
	borderLeft := excelize.Border{Type: "left", Style: 1, Color: "000000"}
	borderRight := excelize.Border{Type: "right", Style: 1, Color: "000000"}
	borderTop := excelize.Border{Type: "top", Style: 1, Color: "000000"}
	borderBottom := excelize.Border{Type: "bottom", Style: 1, Color: "000000"}
	fill := excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"EFEFEF"}}
	styleTableHeader, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "left", WrapText: true},
		Font:      &excelize.Font{Size: 10, Bold: true},
		Border:    []excelize.Border{borderLeft, borderRight, borderTop, borderBottom},
		Fill:      fill,
	})
	styleTableBody, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "left", WrapText: true},
		Font:      &excelize.Font{Size: 10},
		Border:    []excelize.Border{borderLeft, borderRight, borderTop, borderBottom},
	})
	styleTableBodyAlignmentRight, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "right", WrapText: true},
		Font:      &excelize.Font{Size: 10},
		Border:    []excelize.Border{borderLeft, borderRight, borderTop, borderBottom},
	})

	// Configure work sheet
	f.SetColWidth(sheet, "A", "A", 12)
	f.SetColWidth(sheet, "B", "B", 10.5)
	f.SetColWidth(sheet, "C", "F", 7.5)
	f.SetColWidth(sheet, "G", "G", 16.5)
	f.SetColWidth(sheet, "H", "H", 42)
	f.SetColStyle(sheet, "A:H", styleDefault)

	// Write title
	f.MergeCell(sheet, "A1", "H1")
	f.MergeCell(sheet, "A2", "H2")
	f.MergeCell(sheet, "A3", "H3")
	f.SetCellValue(sheet, "A1", loc.CreateString("exportTitle", loc.CreateString("appName")))
	f.SetCellValue(sheet, "A2", overviewEntries.CurrMonthName)
	f.SetCellStyle(sheet, "A1", "A1", styleTitle)
	f.SetCellStyle(sheet, "A2", "A2", styleTextBold)

	// Write summary
	f.MergeCell(sheet, "A4", "H4")
	f.MergeCell(sheet, "B5", "C5")
	f.MergeCell(sheet, "E5", "F5")
	f.MergeCell(sheet, "B6", "C6")
	f.MergeCell(sheet, "E6", "F6")
	f.MergeCell(sheet, "B7", "C7")
	f.MergeCell(sheet, "E7", "F7")
	f.MergeCell(sheet, "B8", "C8")
	f.MergeCell(sheet, "E8", "F8")
	f.MergeCell(sheet, "B9", "C9")
	f.MergeCell(sheet, "E9", "F9")
	f.MergeCell(sheet, "B10", "C10")
	f.MergeCell(sheet, "E10", "F10")
	f.MergeCell(sheet, "A11", "H11")
	f.MergeCell(sheet, "E11", "F11")
	// Create heading
	f.SetCellValue(sheet, "A4", loc.CreateString("overviewHeadingSummary"))
	f.SetCellStyle(sheet, "A4", "A4", styleTextBold)
	// Create target/actual table
	f.SetCellValue(sheet, "A5", loc.CreateString("overviewSummaryLabelTargetHours"))
	f.SetCellValue(sheet, "A6", loc.CreateString("overviewSummaryLabelActualHours"))
	f.SetCellValue(sheet, "A7", loc.CreateString("overviewSummaryLabelBalanceHours"))
	f.SetCellValue(sheet, "B5", overviewEntries.Summary.TargetHours)
	f.SetCellValue(sheet, "B6", overviewEntries.Summary.ActualHours)
	f.SetCellValue(sheet, "B7", overviewEntries.Summary.BalanceHours)
	f.SetCellStyle(sheet, "A5", "A10", styleTableHeader)
	f.SetCellStyle(sheet, "B5", "C10", styleTableBodyAlignmentRight)
	// Create types table
	f.SetCellValue(sheet, "E5", loc.CreateString("entryTypeWork"))
	f.SetCellValue(sheet, "E6", loc.CreateString("entryTypeTravel"))
	f.SetCellValue(sheet, "E7", loc.CreateString("entryTypeVacation"))
	f.SetCellValue(sheet, "E8", loc.CreateString("entryTypeHoliday"))
	f.SetCellValue(sheet, "E9", loc.CreateString("entryTypeIllness"))
	f.SetCellValue(sheet, "G5", overviewEntries.Summary.ActualWorkHours)
	f.SetCellValue(sheet, "G6", overviewEntries.Summary.ActualTravelHours)
	f.SetCellValue(sheet, "G7", overviewEntries.Summary.ActualVacationHours)
	f.SetCellValue(sheet, "G8", overviewEntries.Summary.ActualHolidayHours)
	f.SetCellValue(sheet, "G9", overviewEntries.Summary.ActualIllnessHours)
	f.SetCellValue(sheet, "G10", overviewEntries.Summary.ActualHours)
	f.SetCellStyle(sheet, "E5", "E10", styleTableHeader)
	f.SetCellStyle(sheet, "G5", "G10", styleTableBodyAlignmentRight)

	// Write entries
	// Create heading
	f.MergeCell(sheet, "A12", "H12")
	f.SetCellValue(sheet, "A12", loc.CreateString("overviewHeadingEntries"))
	f.SetCellStyle(sheet, "A12", "A12", styleTextBold)
	// Create table header
	f.SetCellValue(sheet, "A13", loc.CreateString("tableColDate"))
	f.SetCellValue(sheet, "B13", loc.CreateString("tableColType"))
	f.SetCellValue(sheet, "C13", loc.CreateString("tableColStart"))
	f.SetCellValue(sheet, "D13", loc.CreateString("tableColEnd"))
	f.SetCellValue(sheet, "E13", loc.CreateString("tableColNet"))
	if overviewEntries.ShowDetails {
		f.SetCellValue(sheet, "F13", loc.CreateString("tableColActivity"))
		f.SetCellValue(sheet, "G13", loc.CreateString("tableColDescription"))
	}
	f.SetCellStyle(sheet, "A13", "E13", styleTableHeader)
	if overviewEntries.ShowDetails {
		f.SetCellStyle(sheet, "F13", "G13", styleTableHeader)
	}
	// Create table body
	row := 14
	for _, day := range overviewEntries.Days {
		f.SetCellValue(sheet, getCellName("A", row), day.Weekday+" "+day.Date)
		if len(day.Entries) == 0 {
			f.SetCellValue(sheet, getCellName("B", row), "-")
			f.SetCellValue(sheet, getCellName("C", row), "-")
			f.SetCellValue(sheet, getCellName("D", row), "-")
			f.SetCellValue(sheet, getCellName("E", row), "-")
			row++
		} else {
			for _, entry := range day.Entries {
				f.SetCellValue(sheet, getCellName("B", row), entry.EntryType)
				f.SetCellValue(sheet, getCellName("C", row), entry.StartTime)
				f.SetCellValue(sheet, getCellName("D", row), entry.EndTime)
				f.SetCellValue(sheet, getCellName("E", row), entry.Duration)
				f.SetCellValue(sheet, getCellName("F", row), entry.EntryActivity)
				f.SetCellValue(sheet, getCellName("G", row), entry.Description)
				row++
			}
		}
		if len(day.Entries) > 1 {
			f.SetCellValue(sheet, getCellName("E", row), day.WorkDuration)
			row++
		}
	}
	f.SetCellStyle(sheet, "A14", getCellName("E", row-1), styleTableBody)
	if overviewEntries.ShowDetails {
		f.SetCellStyle(sheet, "F14", getCellName("G", row-1), styleTableBody)
	}

	return f
}

func getCellName(col string, row int) string {
	return col + strconv.Itoa(row)
}

// --- Helper functions ---

func (c *EntryController) getUserIdAndUserContract(ctx context.Context) (int, *model.Contract,
	error) {
	// Get current user ID
	userId := getCurrentUserId(ctx)
	// Get user contract
	userContract, err := c.uServ.GetUserContractByUserId(ctx, userId)
	if err != nil {
		return 0, nil, err
	}
	return userId, userContract, nil
}

func (c *EntryController) getEntry(ctx context.Context, entryId int, userId int) (*model.Entry,
	error) {
	entry, err := c.eServ.GetEntryByIdAndUserId(ctx, entryId, userId)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		err := e.NewError(e.LogicEntryNotFound, fmt.Sprintf("Could not find entry %d.", entryId))
		log.Debug(err.StackTrace())
		return nil, err
	}
	return entry, nil
}

func (c *EntryController) getEntryMasterData(ctx context.Context) ([]*model.EntryType,
	[]*model.EntryActivity, error) {
	// Get entry types
	entryTypes, err := c.getEntryTypes(ctx)
	if err != nil {
		return nil, nil, err
	}
	// Get entry activities
	entryActivities, err := c.getEntryActivities(ctx)
	if err != nil {
		return nil, nil, err
	}
	return entryTypes, entryActivities, nil
}

func (c *EntryController) getEntryMasterDataMap(ctx context.Context) (map[int]*model.EntryType,
	map[int]*model.EntryActivity, error) {
	// Get entry types
	entryTypesMap, err := c.getEntryTypesMap(ctx)
	if err != nil {
		return nil, nil, err
	}
	// Get entry activities
	entryActivitiesMap, err := c.getEntryActivitiesMap(ctx)
	if err != nil {
		return nil, nil, err
	}
	return entryTypesMap, entryActivitiesMap, nil
}

func (c *EntryController) getEntryTypes(ctx context.Context) ([]*model.EntryType, error) {
	return c.eServ.GetEntryTypes(ctx)
}

func (c *EntryController) getEntryTypesMap(ctx context.Context) (map[int]*model.EntryType,
	error) {
	return c.eServ.GetEntryTypesMap(ctx)
}

func (c *EntryController) getEntryActivities(ctx context.Context) ([]*model.EntryActivity,
	error) {
	return c.eServ.GetEntryActivities(ctx)
}

func (c *EntryController) getEntryActivitiesMap(ctx context.Context) (map[int]*model.EntryActivity,
	error) {
	return c.eServ.GetEntryActivitiesMap(ctx)
}

func (c *EntryController) convertWorkingHours(workingHours []model.ContractWorkingHours,
) []dailyWorkingDuration {
	dds := make([]dailyWorkingDuration, 0, 10)

	// Create daily durations
	for _, whs := range workingHours {
		m := int(whs.Hours * 60.0)
		d := time.Duration(m) * time.Minute
		dds = append(dds, dailyWorkingDuration{whs.FirstDay, d})
	}

	// Sort daily durations
	sort.SliceStable(dds, func(i, j int) bool {
		return dds[i].fromDate.Before(dds[j].fromDate)
	})

	return dds
}

func (c *EntryController) findWorkingDurationForDate(dailyDurations []dailyWorkingDuration,
	date time.Time) time.Duration {
	d := time.Duration(0)

	// Find daily duration for supplied date
	for _, dd := range dailyDurations {
		if dd.fromDate.After(date) {
			break
		}
		d = dd.duration
	}

	return d
}

func (c *EntryController) convertVacationDays(vacationDays []model.ContractVacationDays,
) []monthlyVacationDays {
	mds := make([]monthlyVacationDays, 0, 10)

	// Create monthly days
	for _, vds := range vacationDays {
		mds = append(mds, monthlyVacationDays{vds.FirstDay, vds.Days})
	}

	// Sort monthly days
	sort.SliceStable(mds, func(i, j int) bool {
		return mds[i].fromDate.Before(mds[j].fromDate)
	})

	return mds
}

func (c *EntryController) findVacationDaysForDate(monthlyDays []monthlyVacationDays,
	date time.Time) float32 {
	d := float32(0.0)

	// Find monthly days for supplied date
	for _, md := range monthlyDays {
		if md.fromDate.After(date) {
			break
		}
		d = md.days
	}

	return d
}

func getDateString(t time.Time) string {
	return t.Format(dateFormat)
}

func getTimeString(t time.Time) string {
	return t.Format(timeFormat)
}

func getRoundedHours(d time.Duration) float32 {
	rd := roundDuration(d)
	return float32(rd.Hours())
}

func roundDuration(d time.Duration) time.Duration {
	return d.Round(time.Minute)
}
