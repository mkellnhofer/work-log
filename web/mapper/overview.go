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

// CreateOverviewEntriesViewModel creates a view model for the overview page.
func (m *OverviewMapper) CreateOverviewEntriesViewModel(prevUrl string, year int, month int,
	userContract *model.Contract, entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity) *vm.OverviewEntries {
	oesvm := &vm.OverviewEntries{}
	oesvm.PreviousUrl = prevUrl

	// Get current month name
	oesvm.CurrMonthName = fmt.Sprintf("%s %d", getMonthName(month), year)

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
	oesvm.CurrMonth = fmt.Sprintf("%d%02d", year, month)
	oesvm.PrevMonth = fmt.Sprintf("%d%02d", py, pm)
	oesvm.NextMonth = fmt.Sprintf("%d%02d", ny, nm)

	// Calculate summary
	oesvm.Summary = m.createSummaryViewModel(year, month, userContract, entries)

	// Create entries
	oesvm.Days = m.createEntriesViewModel(year, month, entries, entryTypesMap,
		entryActivitiesMap)

	return oesvm
}

func (m *OverviewMapper) createSummaryViewModel(year int, month int, userContract *model.Contract,
	entries []*model.Entry) *vm.OverviewEntriesSummary {
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
	return &vm.OverviewEntriesSummary{
		ActualWorkHours:     createRoundedHoursString(actWork),
		ActualTravelHours:   createRoundedHoursString(actTrav),
		ActualVacationHours: createRoundedHoursString(actVaca),
		ActualHolidayHours:  createRoundedHoursString(actHoli),
		ActualIllnessHours:  createRoundedHoursString(actIlln),
		TargetHours:         createRoundedHoursString(tar),
		ActualHours:         createRoundedHoursString(act),
		BalanceHours:        createRoundedHoursString(bal),
	}
}

func (m *OverviewMapper) createEntriesViewModel(year int, month int, entries []*model.Entry,
	entryTypesMap map[int]*model.EntryType, entryActivitiesMap map[int]*model.EntryActivity,
) []*vm.OverviewEntriesDay {
	dsvm := make([]*vm.OverviewEntriesDay, 0, 31)

	curDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	// Create days
	entryIndex := 0
	for {
		// Create and add new day
		dvm := &vm.OverviewEntriesDay{
			Date:         formatShortDate(curDate),
			Weekday:      getShortWeekdayName(curDate),
			IsWeekendDay: curDate.Weekday() == time.Saturday || curDate.Weekday() == time.Sunday,
			Entries:      make([]*vm.OverviewEntry, 0, 10),
		}
		dsvm = append(dsvm, dvm)

		// Create entries
		var colWorkDuration time.Duration
		var dailyWorkDuration time.Duration
		preEntryTypeId := 0
		var evm *vm.OverviewEntry
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
			evm = &vm.OverviewEntry{
				Id:        entry.Id,
				EntryType: m.getEntryTypeDescription(entryTypesMap, entry.TypeId),
				StartTime: formatTime(entry.StartTime),
				EndTime:   formatTime(entry.EndTime),
				Duration:  formatHours(duration),
				EntryActivity: m.getEntryActivityDescription(entryActivitiesMap,
					entry.ActivityId),
				Description: entry.Description,
			}
			dvm.Entries = append(dvm.Entries, evm)

			// Update previous entry type ID
			preEntryTypeId = entry.TypeId

			// Update entry index
			entryIndex++
		}
		dvm.WorkDuration = formatHours(dailyWorkDuration)

		// If next month is reached: Abort
		curDate = curDate.Add(24 * time.Hour)
		if curDate.Month() != time.Month(month) {
			break
		}
	}

	return dsvm
}
