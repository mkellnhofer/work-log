package controller

import (
	"net/http"
	"time"

	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
	"kellnhofer.com/work-log/service"
	"kellnhofer.com/work-log/view"
	vm "kellnhofer.com/work-log/view/model"
)

const pageSize = 7

const dateFormat = "2006-01-02"

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

// --- Helper functions ---

func getDateString(t time.Time) string {
	return t.Format(dateFormat)
}
