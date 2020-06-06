package service

import (
	"fmt"
	"time"

	"kellnhofer.com/work-log/db/repo"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/loc"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
)

// EntryService contains entry related logic.
type EntryService struct {
	eRepo *repo.EntryRepo
}

// NewEntryService create a new entry service.
func NewEntryService(er *repo.EntryRepo) *EntryService {
	return &EntryService{er}
}

// --- Entry functions ---

// GetDateEntries gets all entries (over date).
func (s *EntryService) GetDateEntries(userId int, offset int, limit int) ([]*model.Entry, int,
	*e.Error) {
	// Get entries
	entries, err := s.eRepo.GetDateEntries(userId, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Count all available entries
	cnt, err := s.eRepo.CountDateEntries(userId)
	if err != nil {
		return nil, 0, err
	}

	return entries, cnt, nil
}

// GetEntries gets all entries.
func (s *EntryService) GetEntries(userId int, offset int, limit int) ([]*model.Entry, int,
	*e.Error) {
	// Get entries
	entries, err := s.eRepo.GetEntries(userId, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Count all available entries
	cnt, err := s.eRepo.CountEntries(userId)
	if err != nil {
		return nil, 0, err
	}

	return entries, cnt, nil
}

// GetEntryById gets a entry by its ID.
func (s *EntryService) GetEntryById(id int, userId int) (*model.Entry, *e.Error) {
	return s.eRepo.GetEntryById(id, userId)
}

// CreateEntry creates a new entry.
func (s *EntryService) CreateEntry(entry *model.Entry) *e.Error {
	// Check if entry type exists
	if err := s.checkIfEntryTypeExists(entry.TypeId); err != nil {
		return err
	}
	// Check if entry activity exists
	if err := s.checkIfEntryActivityExists(entry.ActivityId); err != nil {
		return err
	}

	// Check entry
	if err := s.checkEntry(entry); err != nil {
		return err
	}

	// Create entry
	return s.eRepo.CreateEntry(entry)
}

// UpdateEntry updates a entry.
func (s *EntryService) UpdateEntry(entry *model.Entry, userId int) *e.Error {
	// Get existing entry
	existingEntry, err := s.eRepo.GetEntryById(entry.Id, userId)
	if err != nil {
		return err
	}

	// Check if entry exists
	if err := s.checkIfEntryExists(entry.Id, existingEntry); err != nil {
		return err
	}

	// Check if entry type exists
	if err := s.checkIfEntryTypeExists(entry.TypeId); err != nil {
		return err
	}
	// Check if entry activity exists
	if err := s.checkIfEntryActivityExists(entry.ActivityId); err != nil {
		return err
	}

	// Check entry
	if err := s.checkEntry(entry); err != nil {
		return err
	}

	// Update entry
	return s.eRepo.UpdateEntry(entry)
}

// DeleteEntryById deletes a entry by its ID.
func (s *EntryService) DeleteEntryById(id int, userId int) *e.Error {
	// Get existing entry
	existingEntry, err := s.eRepo.GetEntryById(id, userId)
	if err != nil {
		return err
	}

	// Check if entry exists
	if err := s.checkIfEntryExists(id, existingEntry); err != nil {
		return err
	}

	// Delete entry
	return s.eRepo.DeleteEntryById(id)
}

// SearchEntries searches entries (over date).
func (s *EntryService) SearchDateEntries(userId int, params *model.SearchEntriesParams, offset int,
	limit int) ([]*model.Entry, int, *e.Error) {
	// Get entries
	entries, err := s.eRepo.SearchDateEntries(userId, params, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Count all available entries
	cnt, err := s.eRepo.CountSearchDateEntries(userId, params)
	if err != nil {
		return nil, 0, err
	}

	return entries, cnt, nil
}

// GetMonthEntries gets all entries of a month.
func (s *EntryService) GetMonthEntries(userId int, year int, month int) ([]*model.Entry, *e.Error) {
	return s.eRepo.GetMonthEntries(userId, year, month)
}

func (s *EntryService) checkIfEntryExists(id int, entry *model.Entry) *e.Error {
	if entry == nil {
		err := e.NewError(e.LogicEntryNotFound, fmt.Sprintf("Could not find entry %d.", id))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func (s *EntryService) checkEntry(entry *model.Entry) *e.Error {
	if entry.StartTime.After(entry.EndTime) {
		err := e.NewError(e.LogicEntryTimeIntervalInvalid, fmt.Sprintf("End time %s before "+
			"start time %s.", entry.EndTime, entry.StartTime))
		log.Debug(err.StackTrace())
		return err
	}

	workDuration := entry.EndTime.Sub(entry.StartTime)
	if entry.BreakDuration > workDuration {
		err := e.NewError(e.LogicEntryBreakDurationTooLong, fmt.Sprintf("Break duration %s too long.",
			entry.BreakDuration))
		log.Debug(err.StackTrace())
		return err
	}

	return nil
}

// --- Entry type functions ---

// GetEntryTypes gets all entry types.
func (s *EntryService) GetEntryTypes() []*model.EntryType {
	return []*model.EntryType{
		model.NewEntryType(model.EntryTypeIdWork, loc.CreateString("entryTypeWork")),
		model.NewEntryType(model.EntryTypeIdTravel, loc.CreateString("entryTypeTravel")),
		model.NewEntryType(model.EntryTypeIdVacation, loc.CreateString("entryTypeVacation")),
		model.NewEntryType(model.EntryTypeIdHoliday, loc.CreateString("entryTypeHoliday")),
		model.NewEntryType(model.EntryTypeIdIllness, loc.CreateString("entryTypeIllness")),
	}
}

// GetEntryTypesMap gets a map of all entry types.
func (s *EntryService) GetEntryTypesMap() map[int]*model.EntryType {
	// Get entry types
	entryTypes := s.GetEntryTypes()

	// Convert into map
	m := make(map[int]*model.EntryType)
	for _, entryType := range entryTypes {
		m[entryType.Id] = entryType
	}

	return m
}

func (s *EntryService) checkIfEntryTypeExists(id int) *e.Error {
	exist := id == model.EntryTypeIdWork || id == model.EntryTypeIdTravel ||
		id == model.EntryTypeIdVacation || id == model.EntryTypeIdHoliday ||
		id == model.EntryTypeIdIllness
	if !exist {
		err := e.NewError(e.LogicEntryTypeNotFound, fmt.Sprintf("Could not find entry type %d.", id))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

// --- Entry activity functions ---

// GetEntryActivities gets all entry activities.
func (s *EntryService) GetEntryActivities() ([]*model.EntryActivity, *e.Error) {
	return s.eRepo.GetEntryActivities()
}

// GetEntryActivitiesMap gets a map of all entry activities.
func (s *EntryService) GetEntryActivitiesMap() (map[int]*model.EntryActivity, *e.Error) {
	// Get entry activities
	entryActivities, err := s.eRepo.GetEntryActivities()
	if err != nil {
		return nil, err
	}

	// Convert into map
	m := make(map[int]*model.EntryActivity)
	for _, entryActivity := range entryActivities {
		m[entryActivity.Id] = entryActivity
	}

	return m, nil
}

func (s *EntryService) checkIfEntryActivityExists(id int) *e.Error {
	if id == 0 {
		return nil
	}
	exist, err := s.eRepo.ExistsEntryActivityById(id)
	if err != nil {
		return err
	}
	if !exist {
		err = e.NewError(e.LogicEntryActivityNotFound, fmt.Sprintf("Could not find entry activity "+
			"%d.", id))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

// --- Work summary functions ---

// GetTotalWorkSummary gets the total work summary.
func (s *EntryService) GetTotalWorkSummary(userId int) (*model.WorkSummary, *e.Error) {
	start := time.Time{}
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return s.eRepo.GetWorkSummary(userId, start, end)
}
