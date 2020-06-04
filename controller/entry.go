package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/loc"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/service"
	"kellnhofer.com/work-log/util"
	"kellnhofer.com/work-log/view"
	vm "kellnhofer.com/work-log/view/model"
)

const pageSize = 7

const dateFormat = "2006-01-02"
const timeFormat = "15:04"
const dateTimeFormat = "2006-01-02 15:04"

const searchDateTimeFormat = "200601021504"

type entryFormInput struct {
	typeId        string
	date          string
	startTime     string
	endTime       string
	breakDuration string
	activityId    string
	description   string
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
func (c *EntryController) GetListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /list.")
		c.handleShowList(w, r)
	}
}

// GetCreateHandler returns a handler for "GET /create".
func (c *EntryController) GetCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /create.")
		c.handleShowCreate(w, r)
	}
}

// PostCreateHandler returns a handler for "POST /create".
func (c *EntryController) PostCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /create.")
		c.handleExecuteCreate(w, r)
	}
}

// GetEditHandler returns a handler for "GET /edit/{id}".
func (c *EntryController) GetEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /edit/{id}.")
		c.handleShowEdit(w, r)
	}
}

// PostEditHandler returns a handler for "POST /edit/{id}".
func (c *EntryController) PostEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /edit/{id}.")
		c.handleExecuteEdit(w, r)
	}
}

// GetCopyHandler returns a handler for "GET /copy/{id}".
func (c *EntryController) GetCopyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /copy/{id}.")
		c.handleShowCopy(w, r)
	}
}

// PostCopyHandler returns a handler for "POST /copy/{id}".
func (c *EntryController) PostCopyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /copy/{id}.")
		c.handleExecuteCopy(w, r)
	}
}

// PostDeleteHandler returns a handler for "POST /delete/{id}".
func (c *EntryController) PostDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /delete/{id}.")
		c.handleExecuteDelete(w, r)
	}
}

// GetSearchHandler returns a handler for "GET /search".
func (c *EntryController) GetSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /search.")
		// Get search string
		sq := getSearchQueryParam(r)
		if sq == nil {
			c.handleShowSearch(w, r)
		} else {
			c.handleShowListSearch(w, r, *sq)
		}
	}
}

// PostSearchHandler returns a handler for "POST /search".
func (c *EntryController) PostSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /search.")
		c.handleExecuteSearch(w, r)
	}
}

// GetOverviewHandler returns a handler for "GET /overview".
func (c *EntryController) GetOverviewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /overview.")
		c.handleShowOverview(w, r)
	}
}

// PostOverviewHandler returns a handler for "POST /overview".
func (c *EntryController) PostOverviewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /overview.")
		c.handleExecuteOverviewChange(w, r)
	}
}

// GetOverviewExportHandler returns a handler for "GET /overview/export".
func (c *EntryController) GetOverviewExportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /overview/export.")
		c.handleExportOverview(w, r)
	}
}

// --- List handler functions ---

func (c *EntryController) handleShowList(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	userId := getCurrentUserId(r)
	// Get user contract
	userContract := c.getUserContract(userId)

	// Get page number, offset and limit
	pageNum, offset, limit := c.getListPagingParams(r)

	// Get work summary (only for first page)
	var workSummary *model.WorkSummary
	if pageNum == 1 {
		var gwsErr *e.Error
		workSummary, gwsErr = c.eServ.GetTotalWorkSummary(userId)
		if gwsErr != nil {
			panic(gwsErr)
		}
	}

	// Get work entries
	entries, cnt, gesErr := c.eServ.GetDateEntries(userId, offset, limit)
	if gesErr != nil {
		panic(gesErr)
	}
	// Get work entry types
	entryTypesMap := c.getEntryTypesMap()
	// Get work entry activities
	entryActivitiesMap := c.getEntryActivitiesMap()

	// Create view model
	model := c.createListViewModel(userContract, workSummary, pageNum, cnt, entries, entryTypesMap,
		entryActivitiesMap)

	// Save current URL to be able to used later for back navigation
	saveCurrentUrl(r)

	// Render
	view.RenderListEntriesTemplate(w, model)
}

// --- Create handler functions ---

func (c *EntryController) handleShowCreate(w http.ResponseWriter, r *http.Request) {
	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	prevUrl := getPreviousUrl(r)
	entryTypeId := 0
	if len(entryTypes) > 0 {
		entryTypeId = entryTypes[0].Id
	}
	model := c.createCreateViewModel(prevUrl, "", entryTypeId, getDateString(time.Now()), "00:00",
		"00:00", "0", 0, "", entryTypes, entryActivities)

	// Render
	view.RenderCreateEntryTemplate(w, model)
}

func (c *EntryController) handleExecuteCreate(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	userId := getCurrentUserId(r)

	// Get form inputs
	input := c.getEntryFormInput(r)

	// Create model
	entry, cmErr := c.createEntryModel(0, userId, input)
	if cmErr != nil {
		c.handleCreateError(w, r, cmErr, input)
	}

	// Create work entry
	if ceErr := c.eServ.CreateEntry(entry); ceErr != nil {
		c.handleCreateError(w, r, ceErr, input)
	}

	c.handleCreateSuccess(w, r)
}

func (c *EntryController) handleCreateSuccess(w http.ResponseWriter, r *http.Request) {
	prevUrl := getPreviousUrl(r)
	http.Redirect(w, r, prevUrl, http.StatusFound)
}

func (c *EntryController) handleCreateError(w http.ResponseWriter, r *http.Request, err *e.Error,
	input *entryFormInput) {
	// Get error message
	em := loc.GetErrorMessageString(err.Code)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	prevUrl := getPreviousUrl(r)
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.createCreateViewModel(prevUrl, em, entryTypeId, input.date, input.startTime,
		input.endTime, input.breakDuration, entryActivityId, input.description, entryTypes,
		entryActivities)

	// Render
	view.RenderCreateEntryTemplate(w, model)
}

// --- Edit handler functions ---

func (c *EntryController) handleShowEdit(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	userId := getCurrentUserId(r)

	// Get ID
	entryId := getIdPathVar(r)

	// Get work entry
	entry := c.getEntry(entryId, userId)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	prevUrl := getPreviousUrl(r)
	model := c.createEditViewModel(prevUrl, "", entry.Id, entry.TypeId, getDateString(entry.StartTime),
		getTimeString(entry.StartTime), getTimeString(entry.EndTime), getMinutesString(
			entry.BreakDuration), entry.ActivityId, entry.Description, entryTypes, entryActivities)

	// Render
	view.RenderEditEntryTemplate(w, model)
}

func (c *EntryController) handleExecuteEdit(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	userId := getCurrentUserId(r)

	// Get ID
	entryId := getIdPathVar(r)

	// Get form inputs
	input := c.getEntryFormInput(r)

	// Create model
	entry, cmErr := c.createEntryModel(entryId, userId, input)
	if cmErr != nil {
		c.handleEditError(w, r, cmErr, entryId, input)
	}

	// Update work entry
	if ueErr := c.eServ.UpdateEntry(entry, userId); ueErr != nil {
		c.handleEditError(w, r, ueErr, entryId, input)
	}

	c.handleEditSuccess(w, r)
}

func (c *EntryController) handleEditSuccess(w http.ResponseWriter, r *http.Request) {
	prevUrl := getPreviousUrl(r)
	http.Redirect(w, r, prevUrl, http.StatusFound)
}

func (c *EntryController) handleEditError(w http.ResponseWriter, r *http.Request, err *e.Error,
	id int, input *entryFormInput) {
	// Get error message
	em := loc.GetErrorMessageString(err.Code)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	prevUrl := getPreviousUrl(r)
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.createEditViewModel(prevUrl, em, id, entryTypeId, input.date, input.startTime,
		input.endTime, input.breakDuration, entryActivityId, input.description, entryTypes,
		entryActivities)

	// Render
	view.RenderEditEntryTemplate(w, model)
}

// --- Copy handler functions ---

func (c *EntryController) handleShowCopy(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	userId := getCurrentUserId(r)

	// Get ID
	entryId := getIdPathVar(r)

	// Get work entry
	entry := c.getEntry(entryId, userId)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	prevUrl := getPreviousUrl(r)
	model := c.createCopyViewModel(prevUrl, "", entry.Id, entry.TypeId, getDateString(entry.StartTime),
		getTimeString(entry.StartTime), getTimeString(entry.EndTime),
		getMinutesString(entry.BreakDuration), entry.ActivityId, entry.Description, entryTypes,
		entryActivities)

	// Render
	view.RenderCopyEntryTemplate(w, model)
}

func (c *EntryController) handleExecuteCopy(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	userId := getCurrentUserId(r)

	// Get ID
	entryId := getIdPathVar(r)

	// Get form inputs
	input := c.getEntryFormInput(r)

	// Create model
	entry, cmErr := c.createEntryModel(0, userId, input)
	if cmErr != nil {
		c.handleCopyError(w, r, cmErr, entryId, input)
	}

	// Create work entry
	if ceErr := c.eServ.CreateEntry(entry); ceErr != nil {
		c.handleCopyError(w, r, ceErr, entryId, input)
	}

	c.handleCopySuccess(w, r)
}

func (c *EntryController) handleCopySuccess(w http.ResponseWriter, r *http.Request) {
	prevUrl := getPreviousUrl(r)
	http.Redirect(w, r, prevUrl, http.StatusFound)
}

func (c *EntryController) handleCopyError(w http.ResponseWriter, r *http.Request, err *e.Error,
	id int, input *entryFormInput) {
	// Get error message
	em := loc.GetErrorMessageString(err.Code)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	prevUrl := getPreviousUrl(r)
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.createCopyViewModel(prevUrl, em, id, entryTypeId, input.date, input.startTime,
		input.endTime, input.breakDuration, entryActivityId, input.description, entryTypes,
		entryActivities)

	// Render
	view.RenderCopyEntryTemplate(w, model)
}

// --- Delete handler functions ---

func (c *EntryController) handleExecuteDelete(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	userId := getCurrentUserId(r)

	// Get ID
	entryId := getIdPathVar(r)

	// Delete work entry
	if deErr := c.eServ.DeleteEntryById(entryId, userId); deErr != nil {
		panic(deErr)
	}

	c.handleDeleteSuccess(w, r)
}

func (c *EntryController) handleDeleteSuccess(w http.ResponseWriter, r *http.Request) {
	prevUrl := getPreviousUrl(r)
	http.Redirect(w, r, prevUrl, http.StatusFound)
}

// --- Search handler functions ---

func (c *EntryController) handleShowSearch(w http.ResponseWriter, r *http.Request) {
	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	prevUrl := getPreviousUrl(r)
	entryTypeId := 0
	if len(entryTypes) > 0 {
		entryTypeId = entryTypes[0].Id
	}
	model := c.createSearchViewModel(prevUrl, "", false, entryTypeId, false, getDateString(time.Now()),
		getDateString(time.Now()), false, 0, false, "", entryTypes, entryActivities)

	// Render
	view.RenderSearchEntriesTemplate(w, model)
}

func (c *EntryController) handleShowListSearch(w http.ResponseWriter, r *http.Request, query string) {
	// Get current user ID from session
	userId := getCurrentUserId(r)

	// Get page number, offset and limit
	pageNum, offset, limit := c.getListPagingParams(r)

	// Create search params model from query string
	params := c.parseSearchQueryString(query)

	// Get work entries
	entries, cnt, gesErr := c.eServ.SearchDateEntries(userId, params, offset, limit)
	if gesErr != nil {
		panic(gesErr)
	}
	// Get work entry types
	entryTypesMap := c.getEntryTypesMap()
	// Get work entry activities
	entryActivitiesMap := c.getEntryActivitiesMap()

	// Create view model
	model := c.createListSearchViewModel(constant.PathDefault, query, pageNum, cnt, entries,
		entryTypesMap, entryActivitiesMap)

	// Save current URL to be able to used later for back navigation
	saveCurrentUrl(r)

	// Render
	view.RenderListSearchEntriesTemplate(w, model)
}

func (c *EntryController) handleExecuteSearch(w http.ResponseWriter, r *http.Request) {
	// Get form inputs
	input := c.getSearchEntriesFormInput(r)

	// Create search params model from inputs
	params, cmErr := c.createSearchEntriesParamsModel(input)
	if cmErr != nil {
		c.handleSearchError(w, r, cmErr, input)
	}

	c.handleSearchSuccess(w, r, params)
}

func (c *EntryController) handleSearchSuccess(w http.ResponseWriter, r *http.Request,
	params *model.SearchEntriesParams) {
	http.Redirect(w, r, "/search?query="+c.buildSearchQueryString(params), http.StatusFound)
}

func (c *EntryController) handleSearchError(w http.ResponseWriter, r *http.Request, err *e.Error,
	input *searchEntriesFormInput) {
	// Get error message
	em := loc.GetErrorMessageString(err.Code)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	prevUrl := getPreviousUrl(r)
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
	view.RenderSearchEntriesTemplate(w, model)
}

// --- Overview handler functions ---

func (c *EntryController) handleShowOverview(w http.ResponseWriter, r *http.Request) {
	// Get view data
	model := c.getOverviewViewData(r)

	// Render
	view.RenderListOverviewEntriesTemplate(w, model)
}

func (c *EntryController) handleExportOverview(w http.ResponseWriter, r *http.Request) {
	// Get view data
	model := c.getOverviewViewData(r)

	// Create file
	fileName := fmt.Sprintf("work-log-export-%s.xlsx", model.CurrMonth)
	file := exportOverviewEntries(model)

	// Write file
	writeFile(w, fileName, file)
}

func (c *EntryController) getOverviewViewData(r *http.Request) *vm.ListOverviewEntries {
	// Get current user ID from session
	userId := getCurrentUserId(r)
	// Get user contract
	userContract := c.getUserContract(userId)

	// Get user setting
	showDetails, gusErr := c.uServ.GetSettingShowOverviewDetails(userId)
	if gusErr != nil {
		panic(gusErr)
	}

	// Get year and month
	year, month := c.getOverviewParams(r)

	// Get work entries
	entries, gesErr := c.eServ.GetMonthEntries(userId, year, month)
	if gesErr != nil {
		panic(gesErr)
	}
	// Get work entry types
	entryTypesMap := c.getEntryTypesMap()
	// Get work entry activities
	entryActivitiesMap := c.getEntryActivitiesMap()

	// Create view model
	prevUrl := getPreviousUrl(r)
	model := c.createListOverviewViewModel(prevUrl, year, month, userContract, entries,
		entryTypesMap, entryActivitiesMap, showDetails)

	return model
}

func (c *EntryController) getOverviewParams(r *http.Request) (int, int) {
	// Get year and month
	y, m := getMonthQueryParam(r)

	// Was a year and month provided?
	if y != nil && m != nil {
		// Use these
		return *y, *m
	} else {
		// Get current year/month
		t := time.Now()
		return t.Year(), int(t.Month())
	}
}

func (c *EntryController) handleExecuteOverviewChange(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	userId := getCurrentUserId(r)

	// Get form inputs
	input := c.getOverviewFormInput(r)

	// Validate month param
	parseMonthParam(input.month)

	// Update user setting
	showDetails := input.showDetails == "on"
	c.uServ.UpdateSettingShowOverviewDetails(userId, showDetails)

	// Redirect
	http.Redirect(w, r, "/overview?month="+input.month, http.StatusFound)
}

// --- Viem model converter functions ---

func (c *EntryController) createListViewModel(userContract *model.UserContract,
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

	// Create work entries
	lesvm.Days = c.createEntriesViewModel(userContract, entries, entryTypesMap, entryActivitiesMap,
		true)

	return lesvm
}

func (c *EntryController) createListSummaryViewModel(userContract *model.UserContract,
	workSummary *model.WorkSummary) *vm.ListEntriesSummary {
	// If no user contract or work summary was provided: Skip calculation
	if userContract == nil || workSummary == nil {
		return nil
	}

	// Calulate durations
	overtime := c.calculateOvertimeDuration(userContract, workSummary)
	remainingVacation := c.calculateRemainingVacationDuration(userContract, workSummary)

	// Create summary
	lessvm := vm.NewListEntriesSummary()
	lessvm.OvertimeHours = getHoursString(overtime)
	lessvm.RemainingVacationDays = getDaysString(remainingVacation, userContract.DailyWorkingDuration)
	return lessvm
}

func (c *EntryController) calculateOvertimeDuration(userContract *model.UserContract,
	workSummary *model.WorkSummary) time.Duration {
	// Calculate work days since first work day
	start := userContract.FirstWorkDay
	end := time.Now()
	workDays := util.CalculateWorkingDays(start, end)
	log.Verbf("Work days: %s - %s: %d", getDateString(start), getDateString(end), workDays)

	// Calculate target duration
	targetDuration := time.Duration(workDays) * userContract.DailyWorkingDuration
	log.Verbf("Target work duration: %.0f min", targetDuration.Minutes())

	// Calculate actual duration
	var actualDuration time.Duration
	for _, workDuration := range workSummary.WorkDurations {
		actualDuration = actualDuration + workDuration.WorkDuration - workDuration.BreakDuration
	}
	log.Verbf("Actual work duration: %.0f min", actualDuration.Minutes())

	// Calculate overtime
	return userContract.InitOvertimeDuration + actualDuration - targetDuration
}

func (c *EntryController) calculateRemainingVacationDuration(userContract *model.UserContract,
	workSummary *model.WorkSummary) time.Duration {
	// Calculate years since first work day
	years := time.Now().Year() - userContract.FirstWorkDay.Year() + 1
	log.Verbf("Years: %d", years)

	// Calculate total vacation
	totalVacationDays := float32(years)*userContract.AnnualVacationDays + userContract.InitVacationDays
	totalVacation := time.Duration(totalVacationDays) * userContract.DailyWorkingDuration
	log.Verbf("Total vacation hours: %.0f", totalVacation.Hours())

	// Calculate taken vacation
	var takenVacation time.Duration
	for _, workDuration := range workSummary.WorkDurations {
		if workDuration.TypeId == constant.EntryTypeVacation {
			takenVacation = takenVacation + workDuration.WorkDuration - workDuration.BreakDuration
		}
	}
	log.Verbf("Taken vacation hours: %.0f", takenVacation.Hours())

	// Calculate remaining vacation
	return totalVacation - takenVacation
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
	date string, startTime string, endTime string, breakDuration string, activityId int,
	description string, types []*model.EntryType, activities []*model.EntryActivity) *vm.CreateEntry {
	cevm := vm.NewCreateEntry()
	cevm.PreviousUrl = prevUrl
	cevm.ErrorMessage = errorMessage
	cevm.Entry = c.createEntryViewModel(0, typeId, date, startTime, endTime, breakDuration,
		activityId, description)
	cevm.EntryTypes = c.createEntryTypesViewModel(types)
	cevm.EntryActivities = c.createEntryActivitiesViewModel(activities)
	return cevm
}

func (c *EntryController) createEditViewModel(prevUrl string, errorMessage string, id int,
	typeId int, date string, startTime string, endTime string, breakDuration string, activityId int,
	description string, types []*model.EntryType, activities []*model.EntryActivity) *vm.EditEntry {
	eevm := vm.NewEditEntry()
	eevm.PreviousUrl = prevUrl
	eevm.ErrorMessage = errorMessage
	eevm.Entry = c.createEntryViewModel(id, typeId, date, startTime, endTime, breakDuration,
		activityId, description)
	eevm.EntryTypes = c.createEntryTypesViewModel(types)
	eevm.EntryActivities = c.createEntryActivitiesViewModel(activities)
	return eevm
}

func (c *EntryController) createCopyViewModel(prevUrl string, errorMessage string, id int,
	typeId int, date string, startTime string, endTime string, breakDuration string, activityId int,
	description string, types []*model.EntryType, activities []*model.EntryActivity) *vm.CopyEntry {
	cevm := vm.NewCopyEntry()
	cevm.PreviousUrl = prevUrl
	cevm.ErrorMessage = errorMessage
	cevm.Entry = c.createEntryViewModel(id, typeId, date, startTime, endTime, breakDuration,
		activityId, description)
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

	// Create work entries
	lesvm.Days = c.createEntriesViewModel(nil, entries, entryTypesMap, entryActivitiesMap, false)

	return lesvm
}

func (c *EntryController) createEntriesViewModel(userContract *model.UserContract,
	entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity,
	checkMissingOrOverlapping bool) []*vm.ListEntriesDay {
	ldsvm := make([]*vm.ListEntriesDay, 0, pageSize)

	var calcTargetWorkDurationReached bool
	var targetWorkDuration time.Duration

	// If no user contract was provided: Skip target calculation
	if userContract != nil {
		calcTargetWorkDurationReached = true
		targetWorkDuration = userContract.DailyWorkingDuration
	}

	var ldvm *vm.ListEntriesDay
	prevDate := ""
	var prevStartTime *time.Time
	var totalNetWorkDuration time.Duration
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
			totalNetWorkDuration = 0
			totalBreakDuration = 0
			wasTargetWorkDurationReached = ""

			// Create and add new day
			ldvm = vm.NewListEntriesDay()
			ldvm.Date = view.FormatDate(entry.StartTime)
			ldvm.Weekday = view.GetWeekdayName(entry.StartTime)
			ldvm.Entries = make([]*vm.ListEntry, 0, 10)
			ldsvm = append(ldsvm, ldvm)
		}

		// Calculate work duration
		workDuration := entry.EndTime.Sub(entry.StartTime)
		netWorkDuration := workDuration - entry.BreakDuration
		totalNetWorkDuration = totalNetWorkDuration + netWorkDuration
		totalBreakDuration = totalBreakDuration + entry.BreakDuration
		if calcTargetWorkDurationReached {
			reached := (totalNetWorkDuration - targetWorkDuration) >= 0
			wasTargetWorkDurationReached = strconv.FormatBool(reached)
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
		levm.BreakDuration = view.FormatHours(entry.BreakDuration)
		levm.WorkDuration = view.FormatHours(netWorkDuration)
		levm.EntryActivity = c.getEntryActivityDescription(entryActivitiesMap, entry.ActivityId)
		levm.Description = entry.Description
		ldvm.Entries = append(ldvm.Entries, levm)
		ldvm.WorkDuration = view.FormatHours(totalNetWorkDuration)
		ldvm.BreakDuration = view.FormatHours(totalBreakDuration)
		ldvm.WasTargetWorkDurationReached = wasTargetWorkDurationReached
	}

	return ldsvm
}

func (c *EntryController) createEntryViewModel(id int, typeId int, date string, startTime string,
	endTime string, breakDuration string, activityId int, description string) *vm.Entry {
	evm := vm.NewEntry()
	evm.Id = id
	evm.TypeId = typeId
	evm.Date = date
	evm.StartTime = startTime
	evm.EndTime = endTime
	evm.BreakDuration = breakDuration
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

func (c *EntryController) createEntryActivitiesViewModel(entryActivities []*model.EntryActivity) []*vm.
	EntryActivity {
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
	userContract *model.UserContract, entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
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

	// Create work entries
	lesvm.ShowDetails = showDetails
	lesvm.Days = c.createOverviewEntriesViewModel(year, month, entries, entryTypesMap,
		entryActivitiesMap, showDetails)

	return lesvm
}

func (c *EntryController) createOverviewSummaryViewModel(year int, month int,
	userContract *model.UserContract, entries []*model.Entry) *vm.ListOverviewEntriesSummary {

	// Calculate type durations
	var actWork, actTrav, actVaca, actHoli, actIlln time.Duration
	for _, entry := range entries {
		workDuration := entry.EndTime.Sub(entry.StartTime)
		netWorkDuration := workDuration - entry.BreakDuration

		switch entry.TypeId {
		case constant.EntryTypeWork:
			actWork = actWork + netWorkDuration
		case constant.EntryTypeTravel:
			actTrav = actTrav + netWorkDuration
		case constant.EntryTypeVacation:
			actVaca = actVaca + netWorkDuration
		case constant.EntryTypeHoliday:
			actHoli = actHoli + netWorkDuration
		case constant.EntryTypeIllness:
			actIlln = actIlln + netWorkDuration
		}
	}

	// Calculate days
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	workDays := util.CalculateWorkingDays(start, end)

	// Calculate target, actual and balance durations
	var tar time.Duration = time.Duration(workDays) * userContract.DailyWorkingDuration
	var act time.Duration = actWork + actTrav + actVaca + actHoli + actIlln
	var bal time.Duration = act - tar

	// Create summary
	lessvm := vm.NewListOverviewEntriesSummary()
	lessvm.ActualWorkHours = getHoursString(actWork)
	lessvm.ActualTravelHours = getHoursString(actTrav)
	lessvm.ActualVacationHours = getHoursString(actVaca)
	lessvm.ActualHolidayHours = getHoursString(actHoli)
	lessvm.ActualIllnessHours = getHoursString(actIlln)
	lessvm.TargetHours = getHoursString(tar)
	lessvm.ActualHours = getHoursString(act)
	lessvm.BalanceHours = getHoursString(bal)
	return lessvm
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
		var colBreakDuration, colNetWorkDuration time.Duration
		var dailyBreakDuration, dailyNetWorkDuration time.Duration
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
				colBreakDuration = 0
				colNetWorkDuration = 0
				preEntryTypeId = 0
				break
			}

			// Reset collected break and net work duration
			if entry.TypeId != preEntryTypeId {
				colBreakDuration = 0
				colNetWorkDuration = 0
			}

			// Calculate work duration
			workDuration := entry.EndTime.Sub(entry.StartTime)
			netWorkDuration := workDuration - entry.BreakDuration
			colBreakDuration = colBreakDuration + entry.BreakDuration
			colNetWorkDuration = colNetWorkDuration + netWorkDuration
			dailyBreakDuration = dailyBreakDuration + entry.BreakDuration
			dailyNetWorkDuration = dailyNetWorkDuration + netWorkDuration

			// Create and add new entry
			if showDetails {
				levm = vm.NewListOverviewEntry()
				levm.Id = entry.Id
				levm.EntryType = c.getEntryTypeDescription(entryTypesMap, entry.TypeId)
				levm.StartTime = view.FormatTime(entry.StartTime)
				levm.EndTime = view.FormatTime(entry.EndTime)
				levm.BreakDuration = view.FormatHours(entry.BreakDuration)
				levm.WorkDuration = view.FormatHours(netWorkDuration)
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
				levm.BreakDuration = view.FormatHours(colBreakDuration)
				levm.WorkDuration = view.FormatHours(colNetWorkDuration)
			}

			// Update previous entry type ID
			preEntryTypeId = entry.TypeId

			// Update entry index
			entryIndex++
		}
		ldvm.BreakDuration = view.FormatHours(dailyBreakDuration)
		ldvm.WorkDuration = view.FormatHours(dailyNetWorkDuration)

		// If next month is reached: Abort
		curDate = curDate.Add(24 * time.Hour)
		if curDate.Month() != time.Month(month) {
			break
		}
	}

	return ldsvm
}

// --- Paging functions ---

func (c *EntryController) getListPagingParams(r *http.Request) (int, int, int) {
	// Get page number
	pnqp := getPageNumberQueryParam(r)
	pageNum := 1
	if pnqp != nil {
		pageNum = *pnqp
	}

	// Calculate offset and limit
	offset := (pageNum - 1) * pageSize
	limit := pageSize

	return pageNum, offset, limit
}

// --- Form input retrieval functions ---

func (c *EntryController) getEntryFormInput(r *http.Request) *entryFormInput {
	i := entryFormInput{}
	i.typeId = r.FormValue("type")
	i.date = r.FormValue("date")
	i.startTime = r.FormValue("start-time")
	i.endTime = r.FormValue("end-time")
	i.breakDuration = r.FormValue("break-duration")
	i.activityId = r.FormValue("activity")
	i.description = r.FormValue("description")
	return &i
}

func (c *EntryController) getSearchEntriesFormInput(r *http.Request) *searchEntriesFormInput {
	i := searchEntriesFormInput{}
	i.byType = r.FormValue("by-type")
	i.typeId = r.FormValue("type")
	i.byDate = r.FormValue("by-date")
	i.startDate = r.FormValue("start-date")
	i.endDate = r.FormValue("end-date")
	i.byActivity = r.FormValue("by-activity")
	i.activityId = r.FormValue("activity")
	i.byDescription = r.FormValue("by-description")
	i.description = r.FormValue("description")
	return &i
}

func (c *EntryController) getOverviewFormInput(r *http.Request) *overviewFormInput {
	i := overviewFormInput{}
	i.month = r.FormValue("month")
	i.showDetails = r.FormValue("show-details")
	return &i
}

// --- Model converter functions ---

func (c *EntryController) createEntryModel(id int, userId int, input *entryFormInput) (
	*model.Entry, *e.Error) {
	entry := model.NewEntry()
	entry.Id = id
	entry.UserId = userId

	var err *e.Error

	// Convert type ID
	entry.TypeId = c.convertId(input.typeId)

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

	// Convert break duration
	entry.BreakDuration, err = c.convertDuration(input.breakDuration, e.ValBreakDurationInvalid)
	if err != nil {
		return nil, err
	}

	// Convert activity ID
	entry.ActivityId = c.convertId(input.activityId)

	// Validate description
	if err = c.validateString(input.description, 200, e.ValDescriptionTooLong); err != nil {
		return nil, err
	}
	entry.Description = input.description

	return entry, nil
}

func (c *EntryController) createSearchEntriesParamsModel(input *searchEntriesFormInput) (
	*model.SearchEntriesParams, *e.Error) {
	params := model.NewSearchEntriesParams()

	var err *e.Error

	// Convert type ID
	params.ByType = input.byType == "on"
	params.TypeId = c.convertId(input.typeId)

	// Convert start/end time
	params.ByTime = input.byDate == "on"
	params.StartTime, err = c.convertDateTime(input.startDate, "00:00", e.ValStartDateInvalid)
	if err != nil {
		return nil, err
	}
	params.EndTime, err = c.convertDateTime(input.endDate, "23:59", e.ValEndDateInvalid)
	if err != nil {
		return nil, err
	}
	if params.EndTime.Before(params.StartTime) {
		err := e.NewError(e.LogicEntrySearchDateIntervalInvalid, fmt.Sprintf("End date %s before "+
			"start time %s.", input.endDate, input.startDate))
		log.Debug(err.StackTrace())
		return nil, err
	}

	// Convert activity ID
	params.ByActivity = input.byActivity == "on"
	params.ActivityId = c.convertId(input.activityId)

	// Validate description
	params.ByDescription = input.byDescription == "on"
	if err = c.validateString(input.description, 200, e.ValDescriptionTooLong); err != nil {
		return nil, err
	}
	params.Description = input.description

	// Check if search query is empty
	if !params.ByType && !params.ByTime && !params.ByActivity && !params.ByDescription {
		err = e.NewError(e.ValSearchInvalid, "Search query is empty.")
		log.Debug(err.StackTrace())
		return nil, err
	}

	return params, nil
}

func (c *EntryController) convertId(in string) int {
	out, cErr := strconv.Atoi(in)
	if cErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID. (ID must be numeric.)", cErr)
		log.Debug(err.StackTrace())
		panic(err)
	}
	return out
}

func (c *EntryController) convertDateTime(inDate string, inTime string, code int) (time.Time,
	*e.Error) {
	dt := inDate + " " + inTime
	out, pErr := time.Parse(dateTimeFormat, dt)
	if pErr != nil {
		err := e.WrapError(code, fmt.Sprintf("Could not parse time %s.", inTime), pErr)
		log.Debug(err.StackTrace())
		return time.Now(), err
	}
	return out, nil
}

func (c *EntryController) convertDuration(in string, code int) (time.Duration, *e.Error) {
	m, cErr := strconv.Atoi(in)
	if cErr != nil {
		err := e.WrapError(code, fmt.Sprintf("Could not parse duration %s.", in), cErr)
		log.Debug(err.StackTrace())
		return 0, err
	}
	out, pErr := time.ParseDuration(fmt.Sprintf("%dm", m))
	if pErr != nil {
		err := e.WrapError(code, fmt.Sprintf("Could not parse duration %s.", in), pErr)
		log.Debug(err.StackTrace())
		return 0, err
	}
	return out, nil
}

func (c *EntryController) validateString(in string, length int, code int) *e.Error {
	if len(in) >= length {
		err := e.NewError(code, fmt.Sprintf("String too long. (Must be "+
			"< %d characters.)", length))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

// --- Search query functions ---

func (c *EntryController) buildSearchQueryString(params *model.SearchEntriesParams) string {
	var qps []string
	// Add parameter/value for entry type
	if params.ByType {
		qps = append(qps, fmt.Sprintf("typ:%d", params.TypeId))
	}
	// Add parameter/value for entry start/end time
	if params.ByTime {
		qps = append(qps, fmt.Sprintf("tim:%s-%s", formatSearchDate(params.StartTime),
			formatSearchDate(params.EndTime)))
	}
	// Add parameter/value for entry activity
	if params.ByActivity {
		qps = append(qps, fmt.Sprintf("act:%d", params.ActivityId))
	}
	// Add parameter/value for entry description
	if params.ByDescription {
		qps = append(qps, fmt.Sprintf("des:%s", util.EncodeBase64(params.Description)))
	}
	return strings.Join(qps[:], "|")
}

func (c *EntryController) parseSearchQueryString(query string) *model.SearchEntriesParams {
	params := model.NewSearchEntriesParams()

	qps := strings.Split(query, "|")

	// Check if query is empty
	if len(qps) < 1 {
		err := e.NewError(e.ValSearchQueryInvalid, "Search query is empty.")
		log.Debug(err.StackTrace())
		panic(err)
	}

	for _, qp := range qps {
		pv := strings.Split(qp, ":")
		// Check if query part is invalid
		if len(pv) < 2 {
			err := e.NewError(e.ValSearchQueryInvalid, "Search query part is invalid.")
			log.Debug(err.StackTrace())
			panic(err)
		}

		p := pv[0]
		v := pv[1]
		var cErr error

		// Handle specific conversion
		switch p {
		// Convert value for entry type
		case "typ":
			params.ByType = true
			params.TypeId, cErr = strconv.Atoi(v)
		// Convert values for entry start/end time
		case "tim":
			params.ByTime = true
			se := strings.Split(v, "-")
			if len(se) < 2 {
				cErr = errors.New("invalid range")
				break
			}
			params.StartTime, cErr = parseSearchDate(se[0])
			if cErr != nil {
				break
			}
			params.EndTime, cErr = parseSearchDate(se[1])
			if cErr != nil {
				break
			}
		// Convert value for entry activity
		case "act":
			params.ByActivity = true
			params.ActivityId, cErr = strconv.Atoi(v)
		// Convert value for entry description
		case "des":
			params.ByDescription = true
			params.Description, cErr = util.DecodeBase64(v)
		// Unknown parameter
		default:
			err := e.NewError(e.ValSearchQueryInvalid, fmt.Sprintf("Search query parameter '%s' "+
				"is unknown.", p))
			log.Debug(err.StackTrace())
			panic(err)
		}

		// Check if a error occurred
		if cErr != nil {
			err := e.WrapError(e.ValSearchQueryInvalid, fmt.Sprintf("Search query parameter '%s' "+
				"has invalid value.", p), cErr)
			log.Debug(err.StackTrace())
			panic(err)
		}
	}
	return params
}

func formatSearchDate(d time.Time) string {
	return d.Format(searchDateTimeFormat)
}

func parseSearchDate(d string) (time.Time, error) {
	return time.Parse(searchDateTimeFormat, d)
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
	f.SetCellValue(sheet, "E13", loc.CreateString("tableColBreak"))
	f.SetCellValue(sheet, "F13", loc.CreateString("tableColNet"))
	if overviewEntries.ShowDetails {
		f.SetCellValue(sheet, "G13", loc.CreateString("tableColActivity"))
		f.SetCellValue(sheet, "H13", loc.CreateString("tableColDescription"))
	}
	f.SetCellStyle(sheet, "A13", "F13", styleTableHeader)
	if overviewEntries.ShowDetails {
		f.SetCellStyle(sheet, "G13", "H13", styleTableHeader)
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
			f.SetCellValue(sheet, getCellName("F", row), "-")
			row++
		} else {
			for _, entry := range day.Entries {
				f.SetCellValue(sheet, getCellName("B", row), entry.EntryType)
				f.SetCellValue(sheet, getCellName("C", row), entry.StartTime)
				f.SetCellValue(sheet, getCellName("D", row), entry.EndTime)
				f.SetCellValue(sheet, getCellName("E", row), entry.BreakDuration)
				f.SetCellValue(sheet, getCellName("F", row), entry.WorkDuration)
				f.SetCellValue(sheet, getCellName("G", row), entry.EntryActivity)
				f.SetCellValue(sheet, getCellName("H", row), entry.Description)
				row++
			}
		}
		if len(day.Entries) > 1 {
			f.SetCellValue(sheet, getCellName("E", row), day.BreakDuration)
			f.SetCellValue(sheet, getCellName("F", row), day.WorkDuration)
			row++
		}
	}
	f.SetCellStyle(sheet, "A14", getCellName("F", row-1), styleTableBody)
	if overviewEntries.ShowDetails {
		f.SetCellStyle(sheet, "G14", getCellName("H", row-1), styleTableBody)
	}

	return f
}

func getCellName(col string, row int) string {
	return col + strconv.Itoa(row)
}

// --- Helper functions ---

func (c *EntryController) getUserContract(userId int) *model.UserContract {
	userContract, gucErr := c.uServ.GetUserContractByUserId(userId)
	if gucErr != nil {
		panic(gucErr)
	}
	return userContract
}

func (c *EntryController) getEntry(entryId int, userId int) *model.Entry {
	entry, geErr := c.eServ.GetEntryById(entryId, userId)
	if geErr != nil {
		panic(geErr)
	}
	if entry == nil {
		err := e.NewError(e.LogicEntryNotFound, fmt.Sprintf("Could not find work entry %d.", entryId))
		log.Debug(err.StackTrace())
		panic(err)
	}
	return entry
}

func (c *EntryController) getEntryTypes() []*model.EntryType {
	entryTypes, getsErr := c.eServ.GetEntryTypes()
	if getsErr != nil {
		panic(getsErr)
	}
	return entryTypes
}

func (c *EntryController) getEntryTypesMap() map[int]*model.EntryType {
	entryTypesMap, getsErr := c.eServ.GetEntryTypesMap()
	if getsErr != nil {
		panic(getsErr)
	}
	return entryTypesMap
}

func (c *EntryController) getEntryActivities() []*model.EntryActivity {
	entryActivities, geasErr := c.eServ.GetEntryActivities()
	if geasErr != nil {
		panic(geasErr)
	}
	return entryActivities
}

func (c *EntryController) getEntryActivitiesMap() map[int]*model.EntryActivity {
	entryActivitiesMap, geasErr := c.eServ.GetEntryActivitiesMap()
	if geasErr != nil {
		panic(geasErr)
	}
	return entryActivitiesMap
}

func getDateString(t time.Time) string {
	return t.Format(dateFormat)
}

func getTimeString(t time.Time) string {
	return t.Format(timeFormat)
}

func getDaysString(d time.Duration, wd time.Duration) string {
	rd := d.Round(time.Hour)
	h := int(rd.Hours())
	wh := int(wd.Hours())
	days := float32(h) / float32(wh)
	return loc.CreateString("daysValue", days)
}

func getHoursString(d time.Duration) string {
	rd := d.Round(time.Minute)
	return loc.CreateString("hoursValue", rd.Hours())
}

func getMinutesString(d time.Duration) string {
	rd := d.Round(time.Minute)
	return fmt.Sprintf("%d", int(rd.Minutes()))
}
