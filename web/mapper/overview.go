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
	Mapper
}

// NewOverviewMapper creates a new overview mapper.
func NewOverviewMapper() *OverviewMapper {
	return &OverviewMapper{}
}

// CreateOverviewEntriesViewModel creates a view model for the overview page.
func (m *OverviewMapper) CreateOverviewEntriesViewModel(userContract *model.Contract, year int,
	month int, entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity) *vm.OverviewEntries {
	oesvm := &vm.OverviewEntries{}

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
	oesvm.Summary = m.createSummaryViewModel(userContract, year, month, entries)

	// Create entries
	oesvm.Days = m.createEntriesViewModel(year, month, entries, entryTypesMap,
		entryActivitiesMap)

	return oesvm
}

func (m *OverviewMapper) createSummaryViewModel(userContract *model.Contract, year int, month int,
	entries []*model.Entry) *vm.OverviewEntriesSummary {
	// Calculate monthly actual hours per type
	monthTypeActualHours := m.calculateMonthTypeActualHours(entries)

	// Calculate monthly target, actual and balance
	monthTargetHours := m.calculateMonthTargetHours(userContract, year, month)
	monthActualHours := m.calculateMonthActualHours(monthTypeActualHours)
	monthBalanceHours := monthTargetHours - monthActualHours
	monthTotalHours := m.calculateMonthTotalHours(monthActualHours, monthTargetHours)
	monthRemainingHours := m.calculateMonthRemainingHours(monthActualHours, monthTargetHours)

	// Calculate monthly actual percentages per type
	monthTypeActualPercent := m.calculateMonthTypeActualPercentages(monthTypeActualHours,
		monthTotalHours)
	monthActualPercent := m.calculateMonthActualPercentage(monthTypeActualPercent)
	monthRemainingPercent := 100 - monthActualPercent

	// Create summary
	return &vm.OverviewEntriesSummary{
		MonthTargetHours:    getHoursString(monthTargetHours),
		MonthActualHours:    getHoursString(monthActualHours),
		MonthBalanceHours:   getHoursString(monthBalanceHours),
		TypePercentages:     monthTypeActualPercent,
		RemainingPercentage: monthRemainingPercent,
		TypeHours: map[int]string{
			model.EntryTypeIdWork:     getHoursString(monthTypeActualHours[model.EntryTypeIdWork]),
			model.EntryTypeIdTravel:   getHoursString(monthTypeActualHours[model.EntryTypeIdTravel]),
			model.EntryTypeIdVacation: getHoursString(monthTypeActualHours[model.EntryTypeIdVacation]),
			model.EntryTypeIdHoliday:  getHoursString(monthTypeActualHours[model.EntryTypeIdHoliday]),
			model.EntryTypeIdIllness:  getHoursString(monthTypeActualHours[model.EntryTypeIdIllness]),
		},
		RemainingHours: getHoursString(monthRemainingHours),
	}
}

func (m *OverviewMapper) calculateMonthTypeActualHours(entries []*model.Entry) map[int]float32 {
	// Calculate actual durations
	var workDuration, travDuration, vacaDuration, holiDuration, illnDuration time.Duration
	for _, entry := range entries {
		duration := entry.EndTime.Sub(entry.StartTime)
		switch entry.TypeId {
		case model.EntryTypeIdWork:
			workDuration = workDuration + duration
		case model.EntryTypeIdTravel:
			travDuration = travDuration + duration
		case model.EntryTypeIdVacation:
			vacaDuration = vacaDuration + duration
		case model.EntryTypeIdHoliday:
			holiDuration = holiDuration + duration
		case model.EntryTypeIdIllness:
			illnDuration = illnDuration + duration
		}
	}

	// Return rounded hours
	return map[int]float32{
		model.EntryTypeIdWork:     getRoundedHours(workDuration),
		model.EntryTypeIdTravel:   getRoundedHours(travDuration),
		model.EntryTypeIdVacation: getRoundedHours(vacaDuration),
		model.EntryTypeIdHoliday:  getRoundedHours(holiDuration),
		model.EntryTypeIdIllness:  getRoundedHours(illnDuration),
	}
}

func (m *OverviewMapper) calculateMonthActualHours(actualHours map[int]float32) float32 {
	return actualHours[model.EntryTypeIdWork] +
		actualHours[model.EntryTypeIdTravel] +
		actualHours[model.EntryTypeIdVacation] +
		actualHours[model.EntryTypeIdHoliday] +
		actualHours[model.EntryTypeIdIllness]
}

func (m *OverviewMapper) calculateMonthTargetHours(userContract *model.Contract, year int,
	month int) float32 {
	// Get target working durations
	targetWorkDurations := m.convertWorkingHours(userContract.WorkingHours)
	// Abort if no target working durations were set
	if len(targetWorkDurations) == 0 {
		return 0.0
	}

	// Calculate days
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	workDays := util.CalculateWorkingDays(start, end)

	// Calculate actual and target durations
	monthTargetWorkDuration := time.Duration(workDays) * m.findWorkingDurationForDate(
		targetWorkDurations, start)

	// Return rounded hours
	return getRoundedHours(monthTargetWorkDuration)
}

func (m *OverviewMapper) calculateMonthTotalHours(actualHours float32, targetHours float32) float32 {
	totalHours := targetHours
	if actualHours > totalHours {
		totalHours = actualHours
	}
	return totalHours
}

func (m *OverviewMapper) calculateMonthRemainingHours(actualHours float32, targetHours float32,
) float32 {
	remainingHours := targetHours - actualHours
	if remainingHours < 0 {
		remainingHours = 0
	}
	return remainingHours
}

func (m *OverviewMapper) calculateMonthTypeActualPercentages(actualHours map[int]float32,
	totalHours float32) map[int]int {
	wk, wv := m.calculateMonthTypeActualPercentage(model.EntryTypeIdWork, actualHours, totalHours)
	tk, tv := m.calculateMonthTypeActualPercentage(model.EntryTypeIdTravel, actualHours, totalHours)
	vk, vv := m.calculateMonthTypeActualPercentage(model.EntryTypeIdVacation, actualHours, totalHours)
	hk, hv := m.calculateMonthTypeActualPercentage(model.EntryTypeIdHoliday, actualHours, totalHours)
	ik, iv := m.calculateMonthTypeActualPercentage(model.EntryTypeIdIllness, actualHours, totalHours)
	return map[int]int{
		wk: wv,
		tk: tv,
		vk: vv,
		hk: hv,
		ik: iv,
	}
}

func (m *OverviewMapper) calculateMonthTypeActualPercentage(id int, actualHours map[int]float32,
	totalHours float32) (int, int) {
	return id, m.calculatePercentage(actualHours[id], totalHours)
}

func (m *OverviewMapper) calculateMonthActualPercentage(typeActualPercent map[int]int) int {
	return typeActualPercent[model.EntryTypeIdWork] +
		typeActualPercent[model.EntryTypeIdTravel] +
		typeActualPercent[model.EntryTypeIdVacation] +
		typeActualPercent[model.EntryTypeIdHoliday] +
		typeActualPercent[model.EntryTypeIdIllness]
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
