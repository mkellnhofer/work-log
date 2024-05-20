package mapper

import (
	"time"

	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/util"
	vm "kellnhofer.com/work-log/web/model"
)

// LogMapper creates view models for the log page.
type LogMapper struct {
	Mapper
}

// NewLogMapper creates a new log mapper.
func NewLogMapper() *LogMapper {
	return &LogMapper{}
}

// CreateLogSummaryViewModel creates a summary view model for the log page.
func (m *LogMapper) CreateLogSummaryViewModel(userContract *model.Contract, now time.Time,
	totalWorkSummary *model.WorkSummary, monthWorkSummary *model.WorkSummary) *vm.LogSummary {
	// If no user contract or work summary was provided: Skip calculation
	if userContract == nil || totalWorkSummary == nil || monthWorkSummary == nil {
		return nil
	}
	return m.createSummaryViewModel(userContract, now, totalWorkSummary, monthWorkSummary)
}

// CreateLogEntriesViewModel creates a entries view model for the log page.
func (m *LogMapper) CreateLogEntriesViewModel(userContract *model.Contract, curPageNum int,
	totPageNum int, entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity) *vm.ListEntries {
	lesvm := &vm.ListEntries{}

	// Calculate paging nav numbers
	lesvm.CurrentPageNum = curPageNum
	lesvm.FirstPageNum, lesvm.LastPageNum = m.calcPageNavFirstLastPageNums(curPageNum, totPageNum,
		vm.PageNavItems)

	// Create entries
	lesvm.Days = m.createEntriesViewModel(userContract, entries, entryTypesMap, entryActivitiesMap,
		true)

	return lesvm
}

func (m *LogMapper) createSummaryViewModel(userContract *model.Contract, now time.Time,
	totalWorkSummary *model.WorkSummary, monthWorkSummary *model.WorkSummary) *vm.LogSummary {
	// Calculate monthly actual and target
	monthActualHours := m.calculateMonthActualHours(monthWorkSummary)
	monthTargetHours := m.calculateMonthTargetHours(userContract, now)
	monthTotalHours := m.calculateMonthTotalHours(monthActualHours, monthTargetHours)

	// Calculate progress hours
	curLoggedHours := monthActualHours
	curRemainingHours := m.calculateCurrentRemainingHours(monthActualHours, monthTargetHours)
	curRequiredHours := m.calculateCurrentRequiredHours(userContract, now)
	curOvertimeHours, curUndertimeHours := m.calculateCurrentOvertimeUndertimeHours(monthActualHours,
		curRequiredHours)

	// Calculate progress percentages
	curLoggedPercent := m.calculatePercentage(curLoggedHours, monthTotalHours)
	curOvertimePercent := m.calculatePercentage(curOvertimeHours, monthTotalHours)
	curUndertimePercent := m.calculatePercentage(curUndertimeHours, monthTotalHours)
	curRemainingPercent := 100 - curLoggedPercent - curUndertimePercent

	// Calulate total overtime and remaining vacation
	totalOvertimeHours := m.calculateTotalOvertimeHours(userContract, now, totalWorkSummary)
	totalRemainingVacationDays := m.calculateTotalRemainingVacationDays(userContract, now,
		totalWorkSummary)

	// Create summary
	return &vm.LogSummary{
		MonthActualHours:           getHoursString(monthActualHours),
		MonthTargetHours:           getHoursString(monthTargetHours),
		CurrentLoggedPercent:       curLoggedPercent,
		CurrentRemainingPercent:    curRemainingPercent,
		CurrentOvertimePercent:     curOvertimePercent,
		CurrentUndertimePercent:    curUndertimePercent,
		CurrentLoggedHours:         getHoursString(curLoggedHours),
		CurrentRemainingHours:      getHoursString(curRemainingHours),
		CurrentOvertimeHours:       getHoursString(curOvertimeHours),
		CurrentUndertimeHours:      getHoursString(curUndertimeHours),
		TotalOvertimeHours:         getHoursString(totalOvertimeHours),
		TotalRemainingVacationDays: getDaysString(totalRemainingVacationDays),
	}
}

func (m *LogMapper) calculateMonthActualHours(monthWorkSummary *model.WorkSummary) float32 {
	// Calculate actual duration
	var workDuration time.Duration
	for _, wd := range monthWorkSummary.WorkDurations {
		workDuration = workDuration + wd.WorkDuration
	}

	// Return rounded hours
	return getRoundedHours(workDuration)
}

func (m *LogMapper) calculateMonthTargetHours(userContract *model.Contract, now time.Time) float32 {
	// Get target working durations
	targetWorkDurations := m.convertWorkingHours(userContract.WorkingHours)
	// Abort if no target working durations were set
	if len(targetWorkDurations) == 0 {
		return 0.0
	}

	// Create month interval
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local)

	// Calculate work days
	workDays := util.CalculateWorkingDays(start, end)

	// Find target working duration
	targetWorkDuration := m.findWorkingDurationForDate(targetWorkDurations, start)

	// Calculate target duration
	monthTargetWorkDuration := time.Duration(workDays) * targetWorkDuration

	// Return rounded hours
	return getRoundedHours(monthTargetWorkDuration)
}

func (m *LogMapper) calculateMonthTotalHours(actualHours float32, targetHours float32) float32 {
	totalHours := targetHours
	if actualHours > totalHours {
		totalHours = actualHours
	}
	return totalHours
}

func (m *LogMapper) calculateCurrentRemainingHours(actualHours float32, targetHours float32) float32 {
	remainingHours := targetHours - actualHours
	if remainingHours < 0 {
		remainingHours = 0
	}
	return remainingHours
}

func (m *LogMapper) calculateCurrentRequiredHours(userContract *model.Contract, now time.Time,
) float32 {
	// Get target working durations
	targetWorkDurations := m.convertWorkingHours(userContract.WorkingHours)
	// Abort if no target working durations were set
	if len(targetWorkDurations) == 0 {
		return 0.0
	}

	// Create month interval
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)

	// Calculate work days
	workDays := util.CalculateWorkingDays(start, end)

	// Find target working duration
	targetWorkDuration := m.findWorkingDurationForDate(targetWorkDurations, start)

	// Calculate required duration
	requiredWorkDuration := time.Duration(workDays) * targetWorkDuration

	// Return rounded hours
	return getRoundedHours(requiredWorkDuration)
}

func (m *LogMapper) calculateCurrentOvertimeUndertimeHours(actualHours float32, requiredHours float32,
) (float32, float32) {
	overtimeHours := actualHours - requiredHours
	if overtimeHours < 0 {
		overtimeHours = 0
	}
	undertimeHours := requiredHours - actualHours
	if undertimeHours < 0 {
		undertimeHours = 0
	}
	return overtimeHours, undertimeHours
}

func (m *LogMapper) calculateTotalOvertimeHours(userContract *model.Contract, now time.Time,
	workSummary *model.WorkSummary) float32 {
	// Get target working durations
	targetWorkDurations := m.convertWorkingHours(userContract.WorkingHours)
	// Abort if no target working durations were set
	if len(targetWorkDurations) == 0 {
		return 0.0
	}

	// Calculate initial overtime duration
	initOvertimeDuration := time.Duration(int(userContract.InitOvertimeHours*60.0)) * time.Minute

	// Calculate actual duration
	var actualWorkDuration time.Duration
	for _, workDuration := range workSummary.WorkDurations {
		actualWorkDuration = actualWorkDuration + workDuration.WorkDuration
	}

	// Calculate target duration
	start := userContract.FirstDay
	end := now
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

		// Calculate interval target duration
		intTargetWorkDuration := time.Duration(intWorkDays) * targetWorkDurations[i].duration

		// Update target duration
		targetWorkDuration = targetWorkDuration + intTargetWorkDuration
	}

	// Calculate overtime
	overtimeDuration := initOvertimeDuration + actualWorkDuration - targetWorkDuration

	// Return rounded hours
	return getRoundedHours(overtimeDuration)
}

func (m *LogMapper) calculateTotalRemainingVacationDays(userContract *model.Contract, now time.Time,
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

	// Calculate available vacation hours (month by month)
	start := userContract.FirstDay
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

	// Calculate taken vacation hours
	takenVacationHours := float32(0.0)
	for _, workDuration := range workSummary.WorkDurations {
		if workDuration.TypeId == model.EntryTypeIdVacation {
			takenVacationHours = takenVacationHours + float32(workDuration.WorkDuration.Hours())
		}
	}

	// Calculate remaining vacation hours
	remainingVacationHours := initVacationHours + availableVacationHours - takenVacationHours

	// Convert vacation hours to vacation days
	wd := m.findWorkingDurationForDate(workDurations, now)
	var remainingVacationDays float32
	if remainingVacationHours > 0 {
		remainingVacationDays = remainingVacationHours / float32(wd.Hours())
	}

	return remainingVacationDays
}
