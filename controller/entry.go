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
const timeFormat = "15:04"
const dateTimeFormat = "2006-01-02 15:04"

type formInput struct {
	typeId        string
	date          string
	startTime     string
	endTime       string
	breakDuration string
	activityId    string
	description   string
}

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

// GetEditHandler returns a handler for "GET /edit/{id}".
func (c *EntryController) GetEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /edit.")
		c.handleShowEdit(w, r)
	}
}

// PostEditHandler returns a handler for "POST /edit/{id}".
func (c *EntryController) PostEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /edit.")
		c.handleExecuteEdit(w, r)
	}
}

// GetCopyHandler returns a handler for "GET /copy/{id}".
func (c *EntryController) GetCopyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle GET /copy.")
		c.handleShowCopy(w, r)
	}
}

// PostCopyHandler returns a handler for "POST /copy/{id}".
func (c *EntryController) PostCopyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /copy.")
		c.handleExecuteCopy(w, r)
	}
}

// PostDeleteHandler returns a handler for "POST /delete/{id}".
func (c *EntryController) PostDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Verb("Handle POST /delete.")
		c.handleExecuteDelete(w, r)
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

// --- Create handler functions ---

func (c *EntryController) handleShowCreate(w http.ResponseWriter, r *http.Request) {
	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	entryTypeId := 0
	if len(entryTypes) > 0 {
		entryTypeId = entryTypes[0].Id
	}
	model := c.createCreateViewModel("", entryTypeId, getDateString(time.Now()), "00:00", "00:00",
		"0", 0, "", entryTypes, entryActivities)

	// Render
	view.RenderCreateEntryTemplate(w, model)
}

func (c *EntryController) handleExecuteCreate(w http.ResponseWriter, r *http.Request) {
	// Get current session from context
	sess := r.Context().Value(constant.ContextKeySession).(*model.Session)
	// Get current user ID
	userId := sess.UserId

	// Get form inputs
	input := c.getFormInput(r)

	// Create model
	entry, cemErr := c.createEntryModel(0, userId, input)
	if cemErr != nil {
		c.handleCreateError(w, r, cemErr, input)
	}

	// Create work entry
	if ceErr := c.eServ.CreateEntry(entry); ceErr != nil {
		c.handleCreateError(w, r, ceErr, input)
	}

	c.handleCreateSuccess(w, r)
}

func (c *EntryController) handleCreateSuccess(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, constant.PathListFirstPage, http.StatusFound)
}

func (c *EntryController) handleCreateError(w http.ResponseWriter, r *http.Request, err *e.Error,
	input *formInput) {
	// Get error message
	em := getErrorMessage(err.Code)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.createCreateViewModel(em, entryTypeId, input.date, input.startTime, input.endTime,
		input.breakDuration, entryActivityId, input.description, entryTypes, entryActivities)

	// Render
	view.RenderCreateEntryTemplate(w, model)
}

// --- Edit handler functions ---

func (c *EntryController) handleShowEdit(w http.ResponseWriter, r *http.Request) {
	// Get current session from context
	sess := r.Context().Value(constant.ContextKeySession).(*model.Session)
	// Get current user ID
	userId := sess.UserId

	// Get ID
	entryId := getIdPathVar(r)

	// Get work entry
	entry := c.getEntry(entryId, userId)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	model := c.createEditViewModel("", entry.Id, entry.TypeId, getDateString(entry.StartTime),
		getTimeString(entry.StartTime), getTimeString(entry.EndTime), getDurationString(
			entry.BreakDuration), entry.ActivityId, entry.Description, entryTypes, entryActivities)

	// Render
	view.RenderEditEntryTemplate(w, model)
}

func (c *EntryController) handleExecuteEdit(w http.ResponseWriter, r *http.Request) {
	// Get current session from context
	sess := r.Context().Value(constant.ContextKeySession).(*model.Session)
	// Get current user ID
	userId := sess.UserId

	// Get ID
	entryId := getIdPathVar(r)

	// Get form inputs
	input := c.getFormInput(r)

	// Create model
	entry, cemErr := c.createEntryModel(entryId, userId, input)
	if cemErr != nil {
		c.handleEditError(w, r, cemErr, entryId, input)
	}

	// Update work entry
	if ueErr := c.eServ.UpdateEntry(entry, userId); ueErr != nil {
		c.handleEditError(w, r, ueErr, entryId, input)
	}

	c.handleEditSuccess(w, r)
}

func (c *EntryController) handleEditSuccess(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, constant.PathListFirstPage, http.StatusFound)
}

func (c *EntryController) handleEditError(w http.ResponseWriter, r *http.Request, err *e.Error,
	id int, input *formInput) {
	// Get error message
	em := getErrorMessage(err.Code)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.createEditViewModel(em, id, entryTypeId, input.date, input.startTime, input.endTime,
		input.breakDuration, entryActivityId, input.description, entryTypes, entryActivities)

	// Render
	view.RenderEditEntryTemplate(w, model)
}

// --- Copy handler functions ---

func (c *EntryController) handleShowCopy(w http.ResponseWriter, r *http.Request) {
	// Get current session from context
	sess := r.Context().Value(constant.ContextKeySession).(*model.Session)
	// Get current user ID
	userId := sess.UserId

	// Get ID
	entryId := getIdPathVar(r)

	// Get work entry
	entry := c.getEntry(entryId, userId)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	model := c.createCopyViewModel("", entry.Id, entry.TypeId, getDateString(entry.StartTime),
		getTimeString(entry.StartTime), getTimeString(entry.EndTime),
		getDurationString(entry.BreakDuration), entry.ActivityId, entry.Description, entryTypes,
		entryActivities)

	// Render
	view.RenderCopyEntryTemplate(w, model)
}

func (c *EntryController) handleExecuteCopy(w http.ResponseWriter, r *http.Request) {
	// Get current session from context
	sess := r.Context().Value(constant.ContextKeySession).(*model.Session)
	// Get current user ID
	userId := sess.UserId

	// Get ID
	entryId := getIdPathVar(r)

	// Get form inputs
	input := c.getFormInput(r)

	// Create model
	entry, cemErr := c.createEntryModel(0, userId, input)
	if cemErr != nil {
		c.handleCopyError(w, r, cemErr, entryId, input)
	}

	// Create work entry
	if ceErr := c.eServ.CreateEntry(entry); ceErr != nil {
		c.handleCopyError(w, r, ceErr, entryId, input)
	}

	c.handleCopySuccess(w, r)
}

func (c *EntryController) handleCopySuccess(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, constant.PathListFirstPage, http.StatusFound)
}

func (c *EntryController) handleCopyError(w http.ResponseWriter, r *http.Request, err *e.Error,
	id int, input *formInput) {
	// Get error message
	em := getErrorMessage(err.Code)

	// Get work entry types
	entryTypes := c.getEntryTypes()
	// Get work entry activities
	entryActivities := c.getEntryActivities()

	// Create view model
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.createCopyViewModel(em, id, entryTypeId, input.date, input.startTime, input.endTime,
		input.breakDuration, entryActivityId, input.description, entryTypes, entryActivities)

	// Render
	view.RenderCopyEntryTemplate(w, model)
}

// --- Delete handler functions ---

func (c *EntryController) handleExecuteDelete(w http.ResponseWriter, r *http.Request) {
	// Get current session from context
	sess := r.Context().Value(constant.ContextKeySession).(*model.Session)
	// Get current user ID
	userId := sess.UserId

	// Get ID
	entryId := getIdPathVar(r)

	// Delete work entry
	if deErr := c.eServ.DeleteEntryById(entryId, userId); deErr != nil {
		panic(deErr)
	}

	c.handleDeleteSuccess(w, r)
}

func (c *EntryController) handleDeleteSuccess(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, constant.PathListFirstPage, http.StatusFound)
}

// --- Viem model converter functions ---

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
	ldsvm := make([]*vm.ListDay, 0, pageSize)
	var ldvm *vm.ListDay
	prevDate := ""
	var prevStartTime *time.Time
	var totalNetWorkDuration time.Duration
	var totalBreakDuration time.Duration
	for _, entry := range entries {
		currDate := getDateString(entry.StartTime)

		// If new day: Create and add new work day
		if prevDate != currDate {
			prevDate = currDate
			prevStartTime = nil

			// Reset total work and break duration
			totalNetWorkDuration = 0
			totalBreakDuration = 0

			// Create and add new work day
			ldvm = vm.NewListDay()
			ldvm.Date = view.FormatDate(entry.StartTime)
			ldvm.Weekday = view.FormatWeekday(entry.StartTime)
			ldvm.ListEntries = make([]*vm.ListEntry, 0, 10)
			ldsvm = append(ldsvm, ldvm)
		}

		// Calculate work duration
		workDuration := entry.EndTime.Sub(entry.StartTime)
		netWorkDuration := workDuration - entry.BreakDuration
		totalNetWorkDuration = totalNetWorkDuration + netWorkDuration
		totalBreakDuration = totalBreakDuration + entry.BreakDuration

		// Check for missing or overlapping work entry
		if prevStartTime != nil && prevStartTime.After(entry.EndTime) {
			levm := vm.NewListEntry()
			levm.IsMissing = true
			ldvm.ListEntries = append(ldvm.ListEntries, levm)
		} else if prevStartTime != nil && prevStartTime.Before(entry.EndTime) {
			levm := vm.NewListEntry()
			levm.IsOverlapping = true
			ldvm.ListEntries = append(ldvm.ListEntries, levm)
		}
		prevStartTime = &entry.StartTime

		// Create and add new work entry
		levm := vm.NewListEntry()
		levm.Id = entry.Id
		levm.EntryType = c.getEntryTypeDescription(entryTypesMap, entry.TypeId)
		levm.StartTime = view.FormatTime(entry.StartTime)
		levm.EndTime = view.FormatTime(entry.EndTime)
		levm.BreakDuration = view.FormatHours(entry.BreakDuration)
		levm.WorkDuration = view.FormatHours(netWorkDuration)
		levm.EntryActivity = c.getEntryActivityDescription(entryActivitiesMap, entry.ActivityId)
		levm.Description = entry.Description
		ldvm.ListEntries = append(ldvm.ListEntries, levm)
		ldvm.WorkDuration = view.FormatHours(totalNetWorkDuration)
		ldvm.BreakDuration = view.FormatHours(totalBreakDuration)
	}
	lesvm.ListDays = ldsvm

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

func (c *EntryController) createCreateViewModel(errorMessage string, entryTypeId int, date string,
	startTime string, endTime string, breakDuration string, entryActivityId int, description string,
	entryTypes []*model.EntryType, entryActivities []*model.EntryActivity) *vm.CreateEntry {
	cevm := vm.NewCreateEntry()
	cevm.PreviousUrl = constant.PathListFirstPage
	cevm.ErrorMessage = errorMessage
	cevm.Entry = c.createEntryViewModel(0, entryTypeId, date, startTime, endTime, breakDuration,
		entryActivityId, description)
	cevm.EntryTypes = c.createEntryTypesViewModel(entryTypes)
	cevm.EntryActivities = c.createEntryActivitiesViewModel(entryActivities)
	return cevm
}

func (c *EntryController) createEditViewModel(errorMessage string, entryId int, entryTypeId int,
	date string, startTime string, endTime string, breakDuration string, entryActivityId int,
	description string, entryTypes []*model.EntryType, entryActivities []*model.EntryActivity) *vm.
	EditEntry {
	eevm := vm.NewEditEntry()
	eevm.PreviousUrl = constant.PathListFirstPage
	eevm.ErrorMessage = errorMessage
	eevm.Entry = c.createEntryViewModel(entryId, entryTypeId, date, startTime, endTime,
		breakDuration, entryActivityId, description)
	eevm.EntryTypes = c.createEntryTypesViewModel(entryTypes)
	eevm.EntryActivities = c.createEntryActivitiesViewModel(entryActivities)
	return eevm
}

func (c *EntryController) createCopyViewModel(errorMessage string, entryId int, entryTypeId int,
	date string, startTime string, endTime string, breakDuration string, entryActivityId int,
	description string, entryTypes []*model.EntryType, entryActivities []*model.EntryActivity) *vm.
	CopyEntry {
	cevm := vm.NewCopyEntry()
	cevm.PreviousUrl = constant.PathListFirstPage
	cevm.ErrorMessage = errorMessage
	cevm.Entry = c.createEntryViewModel(entryId, entryTypeId, date, startTime, endTime,
		breakDuration, entryActivityId, description)
	cevm.EntryTypes = c.createEntryTypesViewModel(entryTypes)
	cevm.EntryActivities = c.createEntryActivitiesViewModel(entryActivities)
	return cevm
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

// --- Form input retrieval functions ---

func (c *EntryController) getFormInput(r *http.Request) *formInput {
	i := formInput{}
	i.typeId = r.FormValue("type")
	i.date = r.FormValue("date")
	i.startTime = r.FormValue("start-time")
	i.endTime = r.FormValue("end-time")
	i.breakDuration = r.FormValue("break-duration")
	i.activityId = r.FormValue("activity")
	i.description = r.FormValue("description")
	return &i
}

// --- Model converter functions ---

func (c *EntryController) createEntryModel(id int, userId int, input *formInput) (
	*model.Entry, *e.Error) {
	entry := model.NewEntry()
	entry.Id = id
	entry.UserId = userId

	// Validate type ID
	typeId, ptidErr := strconv.Atoi(input.typeId)
	if ptidErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID. (ID must be numeric.)", ptidErr)
		log.Debug(err.StackTrace())
		panic(err)
	}
	entry.TypeId = typeId

	// Validate date
	_, pdErr := parseDateTimeString(input.date, "00:00")
	if pdErr != nil {
		err := e.WrapError(e.ValEntryDateInvalid, fmt.Sprintf("Could not parse date %s.", input.date),
			pdErr)
		log.Debug(err.StackTrace())
		return nil, err
	}

	// Validate start time
	st, pstErr := parseDateTimeString(input.date, input.startTime)
	if pstErr != nil {
		err := e.WrapError(e.ValEntryStartTimeInvalid, fmt.Sprintf("Could not parse start time %s.",
			input.startTime), pstErr)
		log.Debug(err.StackTrace())
		return nil, err
	}
	entry.StartTime = st

	// Validate end time
	et, petErr := parseDateTimeString(input.date, input.endTime)
	if petErr != nil {
		err := e.WrapError(e.ValEntryEndTimeInvalid, fmt.Sprintf("Could not parse end time %s.",
			input.endTime), petErr)
		log.Debug(err.StackTrace())
		return nil, err
	}
	entry.EndTime = et

	// Validate break duration
	bd, pbdErr := parseDurationString(input.breakDuration)
	if pbdErr != nil {
		err := e.WrapError(e.ValEntryBreakDurationInvalid, fmt.Sprintf("Could not parse break "+
			"duration %s.", input.breakDuration), pbdErr)
		log.Debug(err.StackTrace())
		return nil, err
	}
	entry.BreakDuration = bd

	// Validate activity ID
	activityId, paidErr := strconv.Atoi(input.activityId)
	if paidErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID. (ID must be numeric.)", paidErr)
		log.Debug(err.StackTrace())
		panic(err)
	}
	entry.ActivityId = activityId

	// Validate description
	if len(input.description) >= 200 {
		err := e.NewError(e.ValEntryDescriptionTooLong, "Description too long. (Must be < 200 "+
			"characters.)")
		log.Debug(err.StackTrace())
		return nil, err
	}
	entry.Description = input.description

	return entry, nil
}

// --- Helper functions ---

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

func (c *EntryController) getEntryActivities() []*model.EntryActivity {
	entryActivities, geasErr := c.eServ.GetEntryActivities()
	if geasErr != nil {
		panic(geasErr)
	}
	return entryActivities
}

func getDateString(t time.Time) string {
	return t.Format(dateFormat)
}

func getTimeString(t time.Time) string {
	return t.Format(timeFormat)
}

func getDurationString(d time.Duration) string {
	md := d.Round(time.Minute)
	return fmt.Sprintf("%d", int(md.Minutes()))
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
