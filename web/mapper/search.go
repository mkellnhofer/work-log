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

// CreateBasicSearchQueryViewModel creates a view model for the search form.
func (m *SearchMapper) CreateBasicSearchQueryViewModel(filter *model.TextEntryFilter,
) *vm.BasicSearchQuery {
	return &vm.BasicSearchQuery{
		Text: filter.Text,
	}
}

// CreateAdvancedSearchQueryViewModel creates a view model for the search form.
func (m *SearchMapper) CreateAdvancedSearchQueryViewModel(filter *model.FieldEntryFilter,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.AdvancedSearchQuery {
	return &vm.AdvancedSearchQuery{
		Input: &vm.AdvancedSearchQueryInput{
			ByType:         filter.ByType,
			TypeId:         filter.TypeId,
			ByDate:         filter.ByTime,
			StartDate:      formatDate(filter.StartTime),
			StartDateValue: getDateString(filter.StartTime),
			EndDate:        formatDate(filter.EndTime),
			EndDateValue:   getDateString(filter.EndTime),
			ByActivity:     filter.ByActivity,
			ActivityId:     filter.ActivityId,
			ByProject:      filter.ByProject,
			Project:        filter.Project,
			ByDescription:  filter.ByDescription,
			Description:    filter.Description,
			ByLabels:       filter.ByLabel,
			Labels:         filter.Labels,
		},
		EntryTypes:      m.CreateEntryTypesViewModel(types),
		EntryActivities: m.CreateEntryActivitiesViewModel(activities),
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
