package mapper

import (
	"kellnhofer.com/work-log/pkg/model"
	vm "kellnhofer.com/work-log/web/model"
)

// EntryMapper creates view models for the create/edit/copy/delete entry modal.
type EntryMapper struct {
	Mapper
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
