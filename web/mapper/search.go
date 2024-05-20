package mapper

import (
	"time"

	"kellnhofer.com/work-log/pkg/model"
	vm "kellnhofer.com/work-log/web/model"
)

// SearchMapper creates view models for the search page.
type SearchMapper struct {
	Mapper
}

// NewSearchMapper creates a new search mapper.
func NewSearchMapper() *SearchMapper {
	return &SearchMapper{}
}

// CreateSearchViewModel creates a view model for the search page.
func (m *SearchMapper) CreateSearchViewModel(byType bool, typeId int, byDate bool,
	startDate time.Time, endDate time.Time, byActivity bool, activityId int, byDescription bool,
	description string, types []*model.EntryType, activities []*model.EntryActivity) *vm.Search {
	return &vm.Search{
		SearchInput: m.createSearchInputViewModel(byType, typeId, byDate, getDateString(startDate),
			getDateString(endDate), byActivity, activityId, byDescription, description),
		EntryTypes:      m.CreateEntryTypesViewModel(types),
		EntryActivities: m.CreateEntryActivitiesViewModel(activities),
	}
}

func (m *SearchMapper) createSearchInputViewModel(byType bool, typeId int, byDate bool,
	startDate string, endDate string, byActivity bool, activityId int, byDescription bool,
	description string) *vm.SearchInput {
	return &vm.SearchInput{
		ByType:        byType,
		TypeId:        typeId,
		ByDate:        byDate,
		StartDate:     startDate,
		EndDate:       endDate,
		ByActivity:    byActivity,
		ActivityId:    activityId,
		ByDescription: byDescription,
		Description:   description,
	}
}

// CreateSearchEntriesViewModel creates a view model for the search result page.
func (m *SearchMapper) CreateSearchEntriesViewModel(query string, curPageNum int, totPageNum int,
	entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity) *vm.SearchEntries {
	sesvm := &vm.SearchEntries{}
	sesvm.Query = query

	// Calculate paging nav numbers
	sesvm.CurrentPageNum = curPageNum
	sesvm.FirstPageNum, sesvm.LastPageNum = m.calcPageNavFirstLastPageNums(curPageNum, totPageNum,
		vm.PageNavItems)

	// Create entries
	sesvm.Days = m.createEntriesViewModel(nil, entries, entryTypesMap, entryActivitiesMap, false)

	return sesvm
}
