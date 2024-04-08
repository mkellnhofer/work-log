package mapper

import (
	"fmt"
	"time"

	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/util"
	vm "kellnhofer.com/work-log/web/model"
)

// OverviewMapper creates view models for the overview page.
type OverviewMapper struct {
	mapper
}

// NewOverviewMapper creates a new overview mapper.
func NewOverviewMapper() *OverviewMapper {
	return &OverviewMapper{}
}

// CreateOverviewViewModel creates a view model for the overview page.
func (m *OverviewMapper) CreateListOverviewViewModel(prevUrl string, year int, month int,
	userContract *model.Contract, entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity, showDetails bool) *vm.ListOverviewEntries {
	lesvm := vm.NewListOverviewEntries()
	lesvm.PreviousUrl = prevUrl

	// Get current month name
	lesvm.CurrMonthName = fmt.Sprintf("%s %d", getMonthName(month), year)

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
	lesvm.Summary = m.createOverviewSummaryViewModel(year, month, userContract, entries)

	// Create entries
	lesvm.ShowDetails = showDetails
	lesvm.Days = m.createOverviewEntriesViewModel(year, month, entries, entryTypesMap,
		entryActivitiesMap, showDetails)

	return lesvm
}

func (m *OverviewMapper) createOverviewSummaryViewModel(year int, month int,
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
	targetWorkDurations := m.convertWorkingHours(userContract.WorkingHours)

	// Calculate target, actual and balance durations
	var tar time.Duration = time.Duration(workDays) * m.findWorkingDurationForDate(
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

func (m *OverviewMapper) createOverviewEntriesViewModel(year int, month int, entries []*model.Entry,
	entryTypesMap map[int]*model.EntryType, entryActivitiesMap map[int]*model.EntryActivity,
	showDetails bool) []*vm.ListOverviewEntriesDay {
	ldsvm := make([]*vm.ListOverviewEntriesDay, 0, 31)

	curDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	// Create days
	entryIndex := 0
	for {
		// Create and add new day
		ldvm := vm.NewListOverviewEntriesDay()
		ldvm.Date = formatShortDate(curDate)
		ldvm.Weekday = getShortWeekdayName(curDate)
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
				levm.EntryType = m.getEntryTypeDescription(entryTypesMap, entry.TypeId)
				levm.StartTime = formatTime(entry.StartTime)
				levm.EndTime = formatTime(entry.EndTime)
				levm.Duration = formatHours(duration)
				levm.EntryActivity = m.getEntryActivityDescription(entryActivitiesMap, entry.ActivityId)
				levm.Description = entry.Description
				ldvm.Entries = append(ldvm.Entries, levm)
			} else {
				if entry.TypeId != preEntryTypeId {
					levm = vm.NewListOverviewEntry()
					levm.Id = entry.Id
					levm.EntryType = m.getEntryTypeDescription(entryTypesMap, entry.TypeId)
					levm.StartTime = formatTime(entry.StartTime)
					ldvm.Entries = append(ldvm.Entries, levm)
				}
				levm.EndTime = formatTime(entry.EndTime)
				levm.Duration = formatHours(colWorkDuration)
			}

			// Update previous entry type ID
			preEntryTypeId = entry.TypeId

			// Update entry index
			entryIndex++
		}
		ldvm.WorkDuration = formatHours(dailyWorkDuration)

		// If next month is reached: Abort
		curDate = curDate.Add(24 * time.Hour)
		if curDate.Month() != time.Month(month) {
			break
		}
	}

	return ldsvm
}
