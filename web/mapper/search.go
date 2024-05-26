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

// CreateAdvancedSearchViewModel creates a view model for the search page.
func (m *SearchMapper) CreateSearchQueryViewModel(isAdvanced bool, byType bool, typeId int,
	byDate bool, startDate time.Time, endDate time.Time, byActivity bool, activityId int, text string,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.SearchQuery {
	return &vm.SearchQuery{
		IsAdvanced: isAdvanced,
		Input: &vm.SearchQueryInput{
			ByType:     byType,
			TypeId:     typeId,
			ByDate:     byDate,
			StartDate:  getDateString(startDate),
			EndDate:    getDateString(endDate),
			ByActivity: byActivity,
			ActivityId: activityId,
			Text:       text,
		},
		EntryTypes:      m.CreateEntryTypesViewModel(types),
		EntryActivities: m.CreateEntryActivitiesViewModel(activities),
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
