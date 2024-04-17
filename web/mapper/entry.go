package mapper

import (
	"time"

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
func (m *EntryMapper) CreateInitialCreateViewModel(prevUrl string, typeId int, date time.Time,
	activityId int, types []*model.EntryType, activities []*model.EntryActivity) *vm.CreateEntry {
	return m.CreateCreateViewModel(prevUrl, "", typeId, getDateString(date), "00:00", "00:00",
		activityId, "", types, activities)
}

// CreateCreateViewModel creates a view model for the create entry page.
func (m *EntryMapper) CreateCreateViewModel(prevUrl string, errorMessage string, typeId int,
	date string, startTime string, endTime string, activityId int, description string,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.CreateEntry {
	return &vm.CreateEntry{
		PreviousUrl:  prevUrl,
		ErrorMessage: errorMessage,
		Entry: m.createEntryViewModel(0, typeId, date, startTime, endTime, activityId,
			description),
		EntryTypes:      m.createEntryTypesViewModel(types),
		EntryActivities: m.createEntryActivitiesViewModel(activities),
	}
}

// CreateEntryViewModel creates a view model for the edit entry page.
func (m *EntryMapper) CreateInitialEditViewModel(prevUrl string, id int, typeId int,
	startTime time.Time, endTime time.Time, activityId int, description string,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.EditEntry {
	return m.CreateEditViewModel(prevUrl, "", id, typeId, getDateString(startTime),
		getTimeString(startTime), getTimeString(endTime), activityId, description, types, activities)
}

// CreateEditViewModel creates a view model for the edit entry page.
func (m *EntryMapper) CreateEditViewModel(prevUrl string, errorMessage string, id int, typeId int,
	date string, startTime string, endTime string, activityId int, description string,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.EditEntry {
	return &vm.EditEntry{
		PreviousUrl:  prevUrl,
		ErrorMessage: errorMessage,
		Entry: m.createEntryViewModel(id, typeId, date, startTime, endTime, activityId,
			description),
		EntryTypes:      m.createEntryTypesViewModel(types),
		EntryActivities: m.createEntryActivitiesViewModel(activities),
	}
}

// CreateInitialCopyViewModel creates a view model for the copy entry page.
func (m *EntryMapper) CreateInitialCopyViewModel(prevUrl string, id int, typeId int,
	startTime time.Time, endTime time.Time, activityId int, description string,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.CopyEntry {
	return m.CreateCopyViewModel(prevUrl, "", id, typeId, getDateString(startTime),
		getTimeString(startTime), getTimeString(endTime), activityId, description, types, activities)
}

// CreateCopyViewModel creates a view model for the copy entry page.
func (m *EntryMapper) CreateCopyViewModel(prevUrl string, errorMessage string, id int, typeId int,
	date string, startTime string, endTime string, activityId int, description string,
	types []*model.EntryType, activities []*model.EntryActivity) *vm.CopyEntry {
	return &vm.CopyEntry{
		PreviousUrl:  prevUrl,
		ErrorMessage: errorMessage,
		Entry: m.createEntryViewModel(id, typeId, date, startTime, endTime, activityId,
			description),
		EntryTypes:      m.createEntryTypesViewModel(types),
		EntryActivities: m.createEntryActivitiesViewModel(activities),
	}
}

// CreateDeleteViewModel creates a view model for the delete entry page.
func (m *EntryMapper) CreateDeleteViewModel(prevUrl string, errorMessage string, id int,
) *vm.DeleteEntry {
	return &vm.DeleteEntry{
		PreviousUrl:  prevUrl,
		ErrorMessage: errorMessage,
		Id:           id,
	}
}
