package mapper

import (
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

// CreateSearchQueryViewModel creates a view model for the search form.
func (m *SearchMapper) CreateSearchQueryViewModel(filter *model.EntriesFilter,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.SearchQuery {
	return &vm.SearchQuery{
		Input: &vm.SearchQueryInput{
			ByType:         filter.ByType,
			TypeId:         filter.TypeId,
			ByDate:         filter.ByTime,
			StartDate:      formatDate(filter.StartTime),
			StartDateValue: getDateString(filter.StartTime),
			EndDate:        formatDate(filter.EndTime),
			EndDateValue:   getDateString(filter.EndTime),
			ByActivity:     filter.ByActivity,
			ActivityId:     filter.ActivityId,
			Text:           filter.Description,
		},
		EntryTypes:      m.CreateEntryTypesViewModel(types),
		EntryActivities: m.CreateEntryActivitiesViewModel(activities),
	}
}

// CreateSearchDetailsViewModel creates a view model for the search details page.
func (m *SearchMapper) CreateSearchDetailsViewModel(filter *model.EntriesFilter,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.SearchDetails {
	efd := m.CreateEntriesFilterDetailsViewModel(filter, types, activities)
	return &vm.SearchDetails{
		EntriesFilterDetails: *efd,
	}
}

// CreateSearchEntriesViewModel creates a view model for the search result page.
func (m *SearchMapper) CreateSearchEntriesViewModel(curPageNum int, totPageNum int,
	entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity) *vm.ListEntries {
	sesvm := &vm.ListEntries{}

	// Calculate paging nav numbers
	sesvm.CurrentPageNum = curPageNum
	sesvm.FirstPageNum, sesvm.LastPageNum = m.calcPageNavFirstLastPageNums(curPageNum, totPageNum,
		vm.PageNavItems)

	// Create entries
	sesvm.Days = m.createEntriesViewModel(nil, entries, entryTypesMap, entryActivitiesMap, false)

	return sesvm
}
