package mapper

import (
	"kellnhofer.com/work-log/pkg/model"
	vm "kellnhofer.com/work-log/web/model"
)

// EntryMapper creates view models for the create/edit/copy/delete entry page.
type EntryMapper struct {
	Mapper
}

// NewEntryMapper creates a new entry mapper.
func NewEntryMapper() *EntryMapper {
	return &EntryMapper{}
}

// CreateInitialCreateViewModel creates a view model for the create entry page.
func (m *EntryMapper) CreateEntryDataViewModel(entry *model.Entry, types []*model.EntryType,
	activities []*model.EntryActivity) *vm.EntryData {
	return &vm.EntryData{
		Entry:           m.createEntryViewModel(entry),
		EntryTypes:      m.createEntryTypesViewModel(types),
		EntryActivities: m.createEntryActivitiesViewModel(activities),
	}
}
