package mapper

import (
	"kellnhofer.com/work-log/pkg/model"
	vm "kellnhofer.com/work-log/web/model"
)

// EntryMapper creates view models for the create/edit/copy/delete entry modal.
type EntryMapper struct {
	mapper
}

// NewEntryMapper creates a new entry mapper.
func NewEntryMapper() *EntryMapper {
	return &EntryMapper{}
}

// CreateEntryDataViewModel creates a view model for the entry modal.
func (m *EntryMapper) CreateEntryDataViewModel(entry *model.Entry, types []*model.EntryType,
	activities []*model.EntryActivity) *vm.EntryData {
	return &vm.EntryData{
		Entry:           m.CreateEntryViewModel(entry),
		EntryTypes:      m.CreateEntryTypesViewModel(types),
		EntryActivities: m.CreateEntryActivitiesViewModel(activities),
	}
}

// CreateListEntriesViewModel creates a view model for the entries list.
func (m *EntryMapper) CreateListEntriesViewModel(entries []*model.Entry,
	entryTypesMap map[int]*model.EntryType,
	entryActivitiesMap map[int]*model.EntryActivity) *vm.ListEntries {
	lesvm := &vm.ListEntries{}

	// Create entries
	lesvm.Days = m.createEntriesViewModel(nil, entries, entryTypesMap, entryActivitiesMap, false)

	return lesvm
}
