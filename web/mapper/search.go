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

// CreateInitialSearchViewModel creates a view model for the search page.
func (m *SearchMapper) CreateInitialSearchViewModel(typeId int, date time.Time, activityId int,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.Search {
	return m.CreateSearchViewModel("", false, typeId, false, getDateString(date),
		getDateString(date), false, activityId, false, "", types, activities)
}

// CreateSearchViewModel creates a view model for the search page.
func (m *SearchMapper) CreateSearchViewModel(errorMessage string, byType bool,
	typeId int, byDate bool, startDate string, endDate string, byActivity bool, activityId int,
	byDescription bool, description string, types []*model.EntryType,
	activities []*model.EntryActivity) *vm.Search {
	return &vm.Search{
		ErrorMessage: errorMessage,
		SearchInput: m.createSearchInputViewModel(byType, typeId, byDate, startDate, endDate,
			byActivity, activityId, byDescription, description),
		EntryTypes:      m.createEntryTypesViewModel(types),
		EntryActivities: m.createEntryActivitiesViewModel(activities),
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
func (m *SearchMapper) CreateSearchEntriesViewModel(query string, pageNum int, pageSize int, cnt int,
	entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity) *vm.SearchEntries {
	sesvm := &vm.SearchEntries{}
	sesvm.Query = query

	// Calculate previous/next page numbers
	sesvm.HasPrevPage = pageNum > 1
	sesvm.HasNextPage = (pageNum * pageSize) < cnt
	sesvm.PrevPageNum = pageNum - 1
	sesvm.NextPageNum = pageNum + 1

	// Create entries
	sesvm.Days = m.createEntriesViewModel(nil, entries, entryTypesMap, entryActivitiesMap, false)

	return sesvm
}
