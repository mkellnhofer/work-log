package mapper

import (
	"time"

	"kellnhofer.com/work-log/pkg/model"
	vm "kellnhofer.com/work-log/web/model"
)

// SearchMapper creates view models for the search page.
type SearchMapper struct {
	mapper
}

// NewSearchMapper creates a new search mapper.
func NewSearchMapper() *SearchMapper {
	return &SearchMapper{}
}

// CreateInitialSearchEntriesViewModel creates a view model for the search page.
func (m *SearchMapper) CreateInitialSearchEntriesViewModel(prevUrl string, typeId int,
	date time.Time, activityId int, types []*model.EntryType, activities []*model.EntryActivity,
) *vm.SearchEntries {
	return m.CreateSearchEntriesViewModel(prevUrl, "", false, typeId, false, getDateString(date),
		getDateString(date), false, activityId, false, "", types, activities)
}

// CreateSearchEntriesViewModel creates a view model for the search page.
func (m *SearchMapper) CreateSearchEntriesViewModel(prevUrl string, errorMessage string, byType bool,
	typeId int, byDate bool, startDate string, endDate string, byActivity bool, activityId int,
	byDescription bool, description string, types []*model.EntryType,
	activities []*model.EntryActivity) *vm.SearchEntries {
	sevm := vm.NewSearchEntries()
	sevm.PreviousUrl = prevUrl
	sevm.ErrorMessage = errorMessage
	sevm.Search = m.createSearchViewModel(byType, typeId, byDate, startDate, endDate, byActivity,
		activityId, byDescription, description)
	sevm.EntryTypes = m.createEntryTypesViewModel(types)
	sevm.EntryActivities = m.createEntryActivitiesViewModel(activities)
	return sevm
}

func (m *SearchMapper) createSearchViewModel(byType bool, typeId int, byDate bool, startDate string,
	endDate string, byActivity bool, activityId int, byDescription bool, description string,
) *vm.Search {
	svm := vm.NewSearch()
	svm.ByType = byType
	svm.TypeId = typeId
	svm.ByDate = byDate
	svm.StartDate = startDate
	svm.EndDate = endDate
	svm.ByActivity = byActivity
	svm.ActivityId = activityId
	svm.ByDescription = byDescription
	svm.Description = description
	return svm
}

// CreateListSearchViewModel creates a view model for the search list page.
func (m *SearchMapper) CreateListSearchViewModel(prevUrl string, query string, pageNum int,
	pageSize int, cnt int, entries []*model.Entry, entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity) *vm.ListSearchEntries {
	lesvm := vm.NewListSearchEntries()
	lesvm.PreviousUrl = prevUrl
	lesvm.Query = query

	// Calculate previous/next page numbers
	lesvm.HasPrevPage = pageNum > 1
	lesvm.HasNextPage = (pageNum * pageSize) < cnt
	lesvm.PrevPageNum = pageNum - 1
	lesvm.NextPageNum = pageNum + 1

	// Create entries
	lesvm.Days = m.createEntriesViewModel(nil, entries, entryTypesMap, entryActivitiesMap, false)

	return lesvm
}
