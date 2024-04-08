package mapper

import (
	"time"

	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/util"
	vm "kellnhofer.com/work-log/web/model"
)

// ListMapper creates view models for the list page.
type ListMapper struct {
	mapper
}

// NewListMapper creates a new list mapper.
func NewListMapper() *ListMapper {
	return &ListMapper{}
}

// CreateListViewModel creates a view model for the list page.
func (m *ListMapper) CreateListViewModel(userContract *model.Contract, workSummary *model.WorkSummary,
	pageNum int, pageSize int, cnt int, entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity) *vm.ListEntries {
	lesvm := vm.NewListEntries()

	// Calculate summary
	lesvm.Summary = m.createListSummaryViewModel(userContract, workSummary)

	// Calculate previous/next page numbers
	lesvm.HasPrevPage = pageNum > 1
	lesvm.HasNextPage = (pageNum * pageSize) < cnt
	lesvm.PrevPageNum = pageNum - 1
	lesvm.NextPageNum = pageNum + 1

	// Create entries
	lesvm.Days = m.createEntriesViewModel(userContract, entries, entryTypesMap, entryActivitiesMap,
		true)

	return lesvm
}

func (m *ListMapper) createListSummaryViewModel(userContract *model.Contract,
	workSummary *model.WorkSummary) *vm.ListEntriesSummary {
	// If no user contract or work summary was provided: Skip calculation
	if userContract == nil || workSummary == nil {
		return nil
	}

	// Calulate durations
	overtimeHours := m.calculateOvertimeHours(userContract, workSummary)
	remainingVacationDays := m.calculateRemainingVacationDays(userContract, workSummary)

	// Create summary
	lessvm := vm.NewListEntriesSummary()
	lessvm.OvertimeHours = createHoursString(overtimeHours)
	lessvm.RemainingVacationDays = createDaysString(remainingVacationDays)
	return lessvm
}

func (m *ListMapper) calculateOvertimeHours(userContract *model.Contract,
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
	targetWorkDurations := m.convertWorkingHours(userContract.WorkingHours)
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

func (m *ListMapper) calculateRemainingVacationDays(userContract *model.Contract,
	workSummary *model.WorkSummary) float32 {
	// Get monthly vacation days
	vacationDays := m.convertVacationDays(userContract.VacationDays)
	// Get daily working durations
	workDurations := m.convertWorkingHours(userContract.WorkingHours)
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
		vd := m.findVacationDaysForDate(vacationDays, curMonth)
		wd := m.findWorkingDurationForDate(workDurations, curMonth)
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
	wd := m.findWorkingDurationForDate(workDurations, now)
	var remainingVacationDays float32
	if remainingVacationHours > 0 {
		remainingVacationDays = remainingVacationHours / float32(wd.Hours())
	}

	return remainingVacationDays
}
