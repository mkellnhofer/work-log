package service

import (
	"context"
	"fmt"

	"kellnhofer.com/work-log/db/repo"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
)

// EntryService contains work entry related logic.
type EntryService struct {
	eRepo *repo.EntryRepo
}

// NewEntryService create a new work entry service.
func NewEntryService(er *repo.EntryRepo) *EntryService {
	return &EntryService{er}
}

// --- Work entry functions ---

// GetEntries gets all work entries.
func (s *EntryService) GetEntries(userId int, offset int, limit int) ([]*model.Entry, int, int, int,
	*e.Error) {
	// Get work entries
	entries, err := s.eRepo.GetEntries(userId, offset, limit)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// Count all available work entries
	cnt, err := s.eRepo.CountEntries(userId)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return entries, offset, limit, cnt, nil
}

// GetEntryById gets a work entry by its ID.
func (s *EntryService) GetEntryById(id int, userId int) (*model.Entry, *e.Error) {
	return s.eRepo.GetEntryById(id, userId)
}

// CreateEntry creates a new work entry.
func (s *EntryService) CreateEntry(entry *model.Entry) *e.Error {
	// Check if work entry type exists
	if err := s.checkIfEntryTypeExists(entry.TypeId); err != nil {
		return err
	}
	// Check if work entry activity exists
	if err := s.checkIfEntryActivityExists(entry.ActivityId); err != nil {
		return err
	}

	// Check work entry
	if err := s.checkEntry(entry); err != nil {
		return err
	}

	// Create work entry
	return s.eRepo.CreateEntry(entry)
}

// UpdateEntry updates a work entry.
func (s *EntryService) UpdateEntry(entry *model.Entry, userId int) *e.Error {
	// Get existing work entry
	existingEntry, err := s.eRepo.GetEntryById(entry.Id, userId)
	if err != nil {
		return err
	}

	// Check if work entry exists
	if err := s.checkIfEntryExists(entry.Id, existingEntry); err != nil {
		return err
	}

	// Check if work entry type exists
	if err := s.checkIfEntryTypeExists(entry.TypeId); err != nil {
		return err
	}
	// Check if work entry activity exists
	if err := s.checkIfEntryActivityExists(entry.ActivityId); err != nil {
		return err
	}

	// Check work entry
	if err := s.checkEntry(entry); err != nil {
		return err
	}

	// Update work entry
	return s.eRepo.UpdateEntry(entry)
}

// DeleteEntryById deletes a work entry by its ID.
func (s *EntryService) DeleteEntryById(id int, userId int) *e.Error {
	// Get existing work entry
	existingEntry, err := s.eRepo.GetEntryById(id, userId)
	if err != nil {
		return err
	}

	// Check if work entry exists
	if err := s.checkIfEntryExists(id, existingEntry); err != nil {
		return err
	}

	// Delete work entry
	return s.eRepo.DeleteEntryById(id)
}

func (s *EntryService) checkIfEntryExists(id int, entry *model.Entry) *e.Error {
	if entry == nil {
		err := e.NewError(e.LogicEntryNotFound, fmt.Sprintf("Could not find work entry %d.", id))
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

// --- Work entry type functions ---

// GetEntryTypes gets all work entry types.
func (s *EntryService) GetEntryTypes(ctx context.Context) ([]*model.EntryType, *e.Error) {
	return s.eRepo.GetEntryTypes()
}

func (s *EntryService) checkIfEntryTypeExists(id int) *e.Error {
	exist, err := s.eRepo.ExistsEntryTypeById(id)
	if err != nil {
		return err
	}
	if !exist {
		err = e.NewError(e.LogicEntryTypeNotFound, fmt.Sprintf("Could not find work entry type %d.",
			id))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

// --- Work entry activity functions ---

// GetEntryActivities gets all work entry activities.
func (s *EntryService) GetEntryActivities() ([]*model.EntryActivity, *e.Error) {
	return s.eRepo.GetEntryActivities()
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
		err = e.NewError(e.LogicEntryActivityNotFound, fmt.Sprintf("Could not find work entry "+
			"activity %d.", id))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}
