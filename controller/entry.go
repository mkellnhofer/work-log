package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/service"
	"kellnhofer.com/work-log/view"
	vm "kellnhofer.com/work-log/view/model"
)

const pageSize = 7

const dateFormat = "2006-01-02"
const dateTimeFormat = "2006-01-02 15:04"

// EntryController handles requests for entry endpoints.
type EntryController struct {
	eServ *service.EntryService
}

// NewEntryController creates a new entry controller.
func NewEntryController(eServ *service.EntryService) *EntryController {
	return &EntryController{eServ}
}

// --- Endpoints ---

// GetListHandler returns a handler for "GET /list/{page}".
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

// --- List handler functions ---

func (c *EntryController) handleShowList(w http.ResponseWriter, r *http.Request) {
	// Get current session from context
	sess := r.Context().Value(constant.ContextKeySession).(*model.Session)
	// Get current user ID
	userId := sess.UserId

	// Get page number
	pageNum := getPageNumberPathVar(r)

	// Calculate offset and limit
	offset := (pageNum - 1) * pageSize
	limit := pageSize

	// Get work entries
	entries, cnt, gesErr := c.eServ.GetDateEntries(userId, offset, limit)
	if gesErr != nil {
		panic(gesErr)
	}
	// Get work entry types
	entryTypesMap, getErr := c.eServ.GetEntryTypesMap()
	if getErr != nil {
		panic(getErr)
	}
	// Get work entry activities
	entryActivitiesMap, geaErr := c.eServ.GetEntryActivitiesMap()
	if geaErr != nil {
		panic(geaErr)
	}

	// Create view model
	model := c.createShowListViewModel(pageNum, cnt, entries, entryTypesMap, entryActivitiesMap)

	// Render
	view.RenderListEntriesTemplate(w, model)
}

func (c *EntryController) createShowListViewModel(pageNum int, cnt int, entries []*model.Entry,
	entryTypesMap map[int]*model.EntryType, entryActivitiesMap map[int]*model.EntryActivity) *vm.
	ListEntries {
	lesvm := vm.NewListEntries()

	// Calculate previous/next page numbers
	lesvm.HasPrevPage = pageNum > 1
	lesvm.HasNextPage = (pageNum * pageSize) < cnt
	lesvm.PrevPageNum = pageNum - 1
	lesvm.NextPageNum = pageNum + 1

	// Create work entries
	dsvm := make([]*vm.Day, 0, pageSize)
	var dvm *vm.Day
	prevDate := ""
	var totalNetWorkDuration time.Duration
	for _, entry := range entries {
		currDate := getDateString(entry.StartTime)

		// If new day: Create and add new work day
		if prevDate != currDate {
			prevDate = currDate

			// Reset work duration
			totalNetWorkDuration = 0

			// Create and add new work day
			dvm = vm.NewDay()
			dvm.Date = view.FormatDate(entry.StartTime)
			dvm.Weekday = view.FormatWeekday(entry.StartTime)
			dvm.Entries = make([]*vm.Entry, 0, 10)
			dsvm = append(dsvm, dvm)
		}

		// Calculate work duration
		workDuration := entry.EndTime.Sub(entry.StartTime)
		netWorkDuration := workDuration - entry.BreakDuration
		totalNetWorkDuration = totalNetWorkDuration + netWorkDuration

		// Create and add new work entry
		evm := vm.NewEntry()
		evm.Id = entry.Id
		evm.EntryType = c.getEntryTypeDescription(entryTypesMap, entry.TypeId)
		evm.StartTime = view.FormatTime(entry.StartTime)
		evm.EndTime = view.FormatTime(entry.EndTime)
		evm.BreakDuration = view.FormatHours(entry.BreakDuration)
		evm.WorkDuration = view.FormatHours(netWorkDuration)
		evm.EntryActivity = c.getEntryActivityDescription(entryActivitiesMap, entry.ActivityId)
		evm.Description = entry.Description
		dvm.Entries = append(dvm.Entries, evm)
		dvm.WorkDuration = view.FormatHours(totalNetWorkDuration)
	}
	lesvm.Days = dsvm

	return lesvm
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

// --- Create handler functions ---

func (c *EntryController) handleShowCreate(w http.ResponseWriter, r *http.Request) {
	// Get work entry types
	entryTypes, getsErr := c.eServ.GetEntryTypes()
	if getsErr != nil {
		panic(getsErr)
	}
	// Get work entry activities
	entryActivities, geasErr := c.eServ.GetEntryActivities()
	if geasErr != nil {
		panic(geasErr)
	}

	// Create view model
	entryTypeId := 0
	if len(entryTypes) > 0 {
		entryTypeId = entryTypes[0].Id
	}
	model := c.createCreateViewModel("", entryTypeId, entryTypes, getDateString(time.Now()),
		"00:00", "00:00", "0", 0, entryActivities, "")

	// Render
	view.RenderCreateEntryTemplate(w, model)
}

func (c *EntryController) handleExecuteCreate(w http.ResponseWriter, r *http.Request) {
	// Get current session from context
	sess := r.Context().Value(constant.ContextKeySession).(*model.Session)
	// Get current user ID
	userId := sess.UserId

	// Get form inputs
	dateVal := r.FormValue("date")
	typeIdVal := r.FormValue("type")
	startTimeVal := r.FormValue("start-time")
	endTimeVal := r.FormValue("end-time")
	breakDurationVal := r.FormValue("break-duration")
	activityIdVal := r.FormValue("activity")
	descriptionVal := r.FormValue("description")

	// Create model
	entry, viErr := c.createEntryModel(0, userId, typeIdVal, dateVal, startTimeVal, endTimeVal,
		breakDurationVal, activityIdVal, descriptionVal)
	if viErr != nil {
		c.handleCreateError(w, r, viErr, typeIdVal, dateVal, startTimeVal, endTimeVal,
			breakDurationVal, activityIdVal, descriptionVal)
	}

	// Create work entry
	if ceErr := c.eServ.CreateEntry(entry); ceErr != nil {
		c.handleCreateError(w, r, ceErr, typeIdVal, dateVal, startTimeVal, endTimeVal,
			breakDurationVal, activityIdVal, descriptionVal)
	}

	c.handleCreateSuccess(w, r)
}

func (c *EntryController) createEntryModel(id int, userId int, typeIdVal string, dateVal string,
	startTimeVal string, endTimeVal string, breakDurationVal string, activityIdVal string,
	descriptionVal string) (*model.Entry, *e.Error) {
	entry := model.NewEntry()
	entry.Id = id
	entry.UserId = userId

	// Validate type ID
	typeId, ptidErr := strconv.Atoi(typeIdVal)
	if ptidErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID. (ID must be numeric.)", ptidErr)
		log.Debug(err.StackTrace())
		panic(err)
	}
	entry.TypeId = typeId

	// Validate date
	_, pdErr := parseDateTimeString(dateVal, "00:00")
	if pdErr != nil {
		err := e.WrapError(e.ValEntryDateInvalid, fmt.Sprintf("Could not parse date %s.", dateVal),
			pdErr)
		log.Debug(err.StackTrace())
		return nil, err
	}

	// Validate start time
	st, pstErr := parseDateTimeString(dateVal, startTimeVal)
	if pstErr != nil {
		err := e.WrapError(e.ValEntryStartTimeInvalid, fmt.Sprintf("Could not parse start time %s.",
			startTimeVal), pstErr)
		log.Debug(err.StackTrace())
		return nil, err
	}
	entry.StartTime = st

	// Validate end time
	et, petErr := parseDateTimeString(dateVal, endTimeVal)
	if petErr != nil {
		err := e.WrapError(e.ValEntryEndTimeInvalid, fmt.Sprintf("Could not parse end time %s.",
			endTimeVal), petErr)
		log.Debug(err.StackTrace())
		return nil, err
	}
	entry.EndTime = et

	// Validate break duration
	bd, pbdErr := parseDurationString(breakDurationVal)
	if pbdErr != nil {
		err := e.WrapError(e.ValEntryBreakDurationInvalid, fmt.Sprintf("Could not parse break "+
			"duration %s.", breakDurationVal), pbdErr)
		log.Debug(err.StackTrace())
		return nil, err
	}
	entry.BreakDuration = bd

	// Validate activity ID
	activityId, paidErr := strconv.Atoi(activityIdVal)
	if paidErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID. (ID must be numeric.)", paidErr)
		log.Debug(err.StackTrace())
		panic(err)
	}
	entry.ActivityId = activityId

	// Validate description
	if len(descriptionVal) >= 200 {
		err := e.NewError(e.ValEntryDescriptionTooLong, "Description too long. (Must be < 200 "+
			"characters.)")
		log.Debug(err.StackTrace())
		return nil, err
	}
	entry.Description = descriptionVal

	return entry, nil
}

func (c *EntryController) handleCreateSuccess(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/list/1", http.StatusFound)
}

func (c *EntryController) handleCreateError(w http.ResponseWriter, r *http.Request, err *e.Error,
	typeIdVal string, dateVal string, startTimeVal string, endTimeVal string, breakDurationVal string,
	activityIdVal string, descriptionVal string) {
	// Get error message
	em := getErrorMessage(err.Code)

	// Get work entry types
	entryTypes, getsErr := c.eServ.GetEntryTypes()
	if getsErr != nil {
		panic(getsErr)
	}
	// Get work entry activities
	entryActivities, geasErr := c.eServ.GetEntryActivities()
	if geasErr != nil {
		panic(geasErr)
	}

	// Create view model
	entryTypeId, _ := strconv.Atoi(typeIdVal)
	entryActivityId, _ := strconv.Atoi(activityIdVal)
	model := c.createCreateViewModel(em, entryTypeId, entryTypes, dateVal, startTimeVal, endTimeVal,
		breakDurationVal, entryActivityId, entryActivities, descriptionVal)

	// Render
	view.RenderCreateEntryTemplate(w, model)
}

func (c *EntryController) createCreateViewModel(errorMessage string, entryTypeId int,
	entryTypes []*model.EntryType, date string, startTime string, endTime string,
	breakDuration string, entryActivityId int, entryActivities []*model.EntryActivity,
	description string) *vm.CreateEntry {
	cevm := vm.NewCreateEntry()
	cevm.ErrorMessage = errorMessage
	cevm.EntryTypeId = entryTypeId
	cevm.EntryTypes = c.createEntryTypesViewModel(entryTypes)
	cevm.Date = date
	cevm.StartTime = startTime
	cevm.EndTime = endTime
	cevm.BreakDuration = breakDuration
	cevm.EntryActivityId = entryActivityId
	cevm.EntryActivities = c.createEntryActivitiesViewModel(entryActivities)
	cevm.Description = description
	return cevm
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

// --- Helper functions ---

func getDateString(t time.Time) string {
	return t.Format(dateFormat)
}

func parseDateTimeString(dat string, tim string) (time.Time, error) {
	dt := dat + " " + tim
	return time.Parse(dateTimeFormat, dt)
}

func parseDurationString(min string) (time.Duration, error) {
	m, err := strconv.Atoi(min)
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(fmt.Sprintf("%dm", m))
}
