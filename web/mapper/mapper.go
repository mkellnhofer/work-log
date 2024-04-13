package mapper

import (
	"fmt"
	"sort"
	"time"

	"golang.org/x/text/message"

	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/model"
	vm "kellnhofer.com/work-log/web/model"
)

type dailyWorkingDuration struct {
	fromDate time.Time
	duration time.Duration
}

type monthlyVacationDays struct {
	fromDate time.Time
	days     float32
}

type mapper struct {
}

func (m *mapper) createEntriesViewModel(userContract *model.Contract, entries []*model.Entry,
	entryTypesMap map[int]*model.EntryType, entryActivitiesMap map[int]*model.EntryActivity,
	checkMissingOrOverlapping bool) []*vm.ListEntriesDay {
	ldsvm := make([]*vm.ListEntriesDay, 0, 10)

	var calcTargetWorkDurationReached bool
	var targetWorkDurations []dailyWorkingDuration
	targetWorkDuration := time.Duration(0)

	// If no user contract was provided: Skip target calculation
	if userContract != nil {
		calcTargetWorkDurationReached = true
		targetWorkDurations = m.convertWorkingHours(userContract.WorkingHours)
	}

	var ldvm *vm.ListEntriesDay
	prevDate := ""
	var prevStartTime *time.Time
	var totalWorkDuration time.Duration
	var totalBreakDuration time.Duration
	var wasTargetWorkDurationReached bool

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
			wasTargetWorkDurationReached = false

			// Get target work duration
			if calcTargetWorkDurationReached {
				targetWorkDuration = m.findWorkingDurationForDate(targetWorkDurations,
					entry.StartTime)
			}

			// Create and add new day
			ldvm = &vm.ListEntriesDay{
				Date:    formatDate(entry.StartTime),
				Weekday: getWeekdayName(entry.StartTime),
				Entries: make([]*vm.ListEntry, 0, 10),
			}
			ldsvm = append(ldsvm, ldvm)
		}

		// Calculate work duration
		duration := entry.EndTime.Sub(entry.StartTime)
		totalWorkDuration = totalWorkDuration + duration

		// Calculate if target work duration was reached
		if calcTargetWorkDurationReached {
			reached := (totalWorkDuration - targetWorkDuration) >= 0
			wasTargetWorkDurationReached = reached
		}

		// Calculate break duration
		if prevStartTime != nil && prevStartTime.After(entry.EndTime) {
			breakDuration := prevStartTime.Sub(entry.EndTime)
			totalBreakDuration = totalBreakDuration + breakDuration
		}

		// Check for missing or overlapping entry
		if checkMissingOrOverlapping {
			if prevStartTime != nil && prevStartTime.After(entry.EndTime) {
				ldvm.Entries = append(ldvm.Entries, &vm.ListEntry{
					IsMissing: true,
				})
			} else if prevStartTime != nil && prevStartTime.Before(entry.EndTime) {
				ldvm.Entries = append(ldvm.Entries, &vm.ListEntry{
					IsOverlapping: true,
				})
			}
		}
		prevStartTime = &entry.StartTime

		// Create and add new entry
		ldvm.Entries = append(ldvm.Entries, &vm.ListEntry{
			Id:            entry.Id,
			EntryType:     m.getEntryTypeDescription(entryTypesMap, entry.TypeId),
			StartTime:     formatTime(entry.StartTime),
			EndTime:       formatTime(entry.EndTime),
			Duration:      formatHours(duration),
			EntryActivity: m.getEntryActivityDescription(entryActivitiesMap, entry.ActivityId),
			Description:   entry.Description,
		})

		// Set work/break durations
		ldvm.WorkDuration = formatHours(totalWorkDuration)
		ldvm.BreakDuration = formatHours(totalBreakDuration)
		ldvm.WasTargetWorkDurationReached = wasTargetWorkDurationReached
	}

	return ldsvm
}

func (m *mapper) createEntryTypesViewModel(entryTypes []*model.EntryType) []*vm.EntryType {
	etsvm := make([]*vm.EntryType, 0, 10)
	for _, entryType := range entryTypes {
		etsvm = append(etsvm, m.createEntryTypeViewModel(entryType.Id, entryType.Description))
	}
	return etsvm
}

func (m *mapper) createEntryTypeViewModel(id int, description string) *vm.EntryType {
	return vm.NewEntryType(id, description)
}

func (m *mapper) createEntryActivitiesViewModel(entryActivities []*model.EntryActivity,
) []*vm.EntryActivity {
	easvm := make([]*vm.EntryActivity, 0, 10)
	easvm = append(easvm, m.createEntryActivityViewModel(0, "-"))
	for _, entryActivity := range entryActivities {
		easvm = append(easvm, m.createEntryActivityViewModel(entryActivity.Id,
			entryActivity.Description))
	}
	return easvm
}

func (m *mapper) createEntryActivityViewModel(id int, description string) *vm.EntryActivity {
	return vm.NewEntryActivity(id, description)
}

func (m *mapper) createEntryViewModel(id int, typeId int, date string, startTime string,
	endTime string, activityId int, description string) *vm.Entry {
	return &vm.Entry{
		Id:          id,
		TypeId:      typeId,
		Date:        date,
		StartTime:   startTime,
		EndTime:     endTime,
		ActivityId:  activityId,
		Description: description,
	}
}

func (m *mapper) getEntryTypeDescription(entryTypesMap map[int]*model.EntryType, id int) string {
	et, ok := entryTypesMap[id]
	if ok {
		return et.Description
	}
	return ""
}

func (m *mapper) getEntryActivityDescription(entryActivitiesMap map[int]*model.EntryActivity,
	id int) string {
	ea, ok := entryActivitiesMap[id]
	if ok {
		return ea.Description
	}
	return ""
}

func (m *mapper) convertVacationDays(vacationDays []model.ContractVacationDays,
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

func (m *mapper) findVacationDaysForDate(monthlyDays []monthlyVacationDays, date time.Time,
) float32 {
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

func (m *mapper) convertWorkingHours(workingHours []model.ContractWorkingHours,
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

func (m *mapper) findWorkingDurationForDate(dailyDurations []dailyWorkingDuration, date time.Time,
) time.Duration {
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

// --- Date and time helpers ---

var weekdayKeys = map[int]string{
	0: "weekdaySun",
	1: "weekdayMon",
	2: "weekdayTue",
	3: "weekdayWed",
	4: "weekdayThu",
	5: "weekdayFri",
	6: "weekdaySat",
}

var monthKeys = map[int]string{
	1:  "monthJan",
	2:  "monthFeb",
	3:  "monthMar",
	4:  "monthApr",
	5:  "monthMay",
	6:  "monthJun",
	7:  "monthJul",
	8:  "monthAug",
	9:  "monthSep",
	10: "monthOct",
	11: "monthNov",
	12: "monthDec",
}

func getMonthName(m int) string {
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf(monthKeys[m])
}

func getWeekdayName(t time.Time) string {
	wd := t.Weekday()
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf(weekdayKeys[int(wd)])
}

func getShortWeekdayName(t time.Time) string {
	d := getWeekdayName(t)
	return fmt.Sprintf("%s.", d[0:2])
}

func getDateString(t time.Time) string {
	return t.Format("2006-01-02")
}

func getTimeString(t time.Time) string {
	return t.Format("15:04")
}

func formatDate(t time.Time) string {
	return t.Format("02.01.2006")
}

func formatShortDate(t time.Time) string {
	return t.Format("02.01.")
}

func formatTime(t time.Time) string {
	return t.Format("15:04")
}

func createDaysString(days float32) string {
	return loc.CreateString("daysValue", days)
}

func createHoursString(hours float32) string {
	return loc.CreateString("hoursValue", hours)
}

func createRoundedHoursString(d time.Duration) string {
	return createHoursString(getRoundedHours(d))
}

func getRoundedHours(d time.Duration) float32 {
	rd := roundDuration(d)
	return float32(rd.Hours())
}

func formatHours(d time.Duration) string {
	h := d.Hours()
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf("%.2f", h)
}

func formatMinutes(d time.Duration) string {
	m := d.Minutes()
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf("%d", int(m))
}

func roundDuration(d time.Duration) time.Duration {
	return d.Round(time.Minute)
}
