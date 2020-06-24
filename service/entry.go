package service

import (
	"context"
	"fmt"
	"time"

	"kellnhofer.com/work-log/db/repo"
	"kellnhofer.com/work-log/db/tx"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/loc"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
)

// EntryService contains entry related logic.
type EntryService struct {
	service
	eRepo *repo.EntryRepo
}

// NewEntryService create a new entry service.
func NewEntryService(tm *tx.TransactionManager, er *repo.EntryRepo) *EntryService {
	return &EntryService{service{tm}, er}
}

// --- Entry functions ---

// GetDateEntries gets all entries (over date).
func (s *EntryService) GetDateEntries(ctx context.Context, filter *model.EntriesFilter, offset int,
	limit int) ([]*model.Entry, int, *e.Error) {
	// Check permissions
	userId := 0
	if filter != nil {
		userId = filter.UserId
	}
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, 0, err
	}

	// Get entries
	entries, err := s.eRepo.GetDateEntries(ctx, filter, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Count all available entries
	cnt, err := s.eRepo.CountDateEntries(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return entries, cnt, nil
}

// GetDateEntriesByUserId gets all entries (over date) of an user.
func (s *EntryService) GetDateEntriesByUserId(ctx context.Context, userId int, offset int,
	limit int) ([]*model.Entry, int, *e.Error) {
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, 0, err
	}

	// Get entries
	entries, err := s.eRepo.GetDateEntriesByUserId(ctx, userId, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Count all available entries
	cnt, err := s.eRepo.CountDateEntriesByUserId(ctx, userId)
	if err != nil {
		return nil, 0, err
	}

	return entries, cnt, nil
}

// GetEntries gets all entries.
func (s *EntryService) GetEntries(ctx context.Context, filter *model.EntriesFilter, offset int,
	limit int) ([]*model.Entry, int, *e.Error) {
	// Check permissions
	userId := 0
	if filter != nil {
		userId = filter.UserId
	}
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, 0, err
	}

	// Get entries
	entries, err := s.eRepo.GetEntries(ctx, filter, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Count all available entries
	cnt, err := s.eRepo.CountEntries(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return entries, cnt, nil
}

// GetEntryById gets an entry.
func (s *EntryService) GetEntryById(ctx context.Context, id int) (*model.Entry, *e.Error) {
	// Get entry
	entry, err := s.eRepo.GetEntryById(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, entry.UserId); err != nil {
		return nil, err
	}

	return entry, nil
}

// GetEntryByIdAndUserId gets an entry of an user.
func (s *EntryService) GetEntryByIdAndUserId(ctx context.Context, id int, userId int) (*model.Entry,
	*e.Error) {
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, err
	}

	// Get entry
	return s.eRepo.GetEntryByIdAndUserId(ctx, id, userId)
}

// CreateEntry creates a new entry.
func (s *EntryService) CreateEntry(ctx context.Context, entry *model.Entry) *e.Error {
	// Check permissions
	if err := s.checkHasCurrentUserChangeRight(ctx, entry.UserId); err != nil {
		return err
	}

	// Check if entry type exists
	if err := s.checkIfEntryTypeExists(entry.TypeId); err != nil {
		return err
	}
	// Check if entry activity exists
	if err := s.checkIfEntryActivityExists(ctx, entry.ActivityId); err != nil {
		return err
	}

	// Check entry
	if err := s.checkEntry(entry); err != nil {
		return err
	}

	// Create entry
	return s.eRepo.CreateEntry(ctx, entry)
}

// UpdateEntry updates an entry.
func (s *EntryService) UpdateEntry(ctx context.Context, entry *model.Entry) *e.Error {
	// Get existing entry
	existingEntry, err := s.eRepo.GetEntryByIdAndUserId(ctx, entry.Id, entry.UserId)
	if err != nil {
		return err
	}

	// Check if entry exists
	if err := s.checkIfEntryExists(entry.Id, existingEntry); err != nil {
		return err
	}

	// Check permissions
	if err := s.checkHasCurrentUserChangeRight(ctx, existingEntry.UserId); err != nil {
		return err
	}

	// Check if entry type exists
	if err := s.checkIfEntryTypeExists(entry.TypeId); err != nil {
		return err
	}
	// Check if entry activity exists
	if err := s.checkIfEntryActivityExists(ctx, entry.ActivityId); err != nil {
		return err
	}

	// Check entry
	if err := s.checkEntry(entry); err != nil {
		return err
	}

	// Update entry
	return s.eRepo.UpdateEntry(ctx, entry)
}

// DeleteEntryById deletes an entry.
func (s *EntryService) DeleteEntryById(ctx context.Context, id int) *e.Error {
	// Get existing entry
	existingEntry, err := s.eRepo.GetEntryById(ctx, id)
	if err != nil {
		return err
	}

	// Check if entry exists
	if err := s.checkIfEntryExists(id, existingEntry); err != nil {
		return err
	}

	// Check permissions
	if err := s.checkHasCurrentUserChangeRight(ctx, existingEntry.UserId); err != nil {
		return err
	}

	// Delete entry
	return s.eRepo.DeleteEntryById(ctx, id)
}

// DeleteEntryByIdAndUserId deletes an entry of an user.
func (s *EntryService) DeleteEntryByIdAndUserId(ctx context.Context, id int, userId int) *e.Error {
	// Check permissions
	if err := s.checkHasCurrentUserChangeRight(ctx, userId); err != nil {
		return err
	}

	// Get existing entry
	existingEntry, err := s.eRepo.GetEntryByIdAndUserId(ctx, id, userId)
	if err != nil {
		return err
	}

	// Check if entry exists
	if err := s.checkIfEntryExists(id, existingEntry); err != nil {
		return err
	}

	// Delete entry
	return s.eRepo.DeleteEntryById(ctx, id)
}

// GetMonthEntriesByUserId gets all entries of a month of an user.
func (s *EntryService) GetMonthEntriesByUserId(ctx context.Context, userId int, year int,
	month int) ([]*model.Entry, *e.Error) {
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, err
	}

	// Get entries
	return s.eRepo.GetMonthEntries(ctx, userId, year, month)
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
func (s *EntryService) GetEntryTypes(ctx context.Context) ([]*model.EntryType, *e.Error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetEntryCharacts); err != nil {
		return nil, err
	}

	// Get entry types
	return []*model.EntryType{
		model.NewEntryType(model.EntryTypeIdWork, loc.CreateString("entryTypeWork")),
		model.NewEntryType(model.EntryTypeIdTravel, loc.CreateString("entryTypeTravel")),
		model.NewEntryType(model.EntryTypeIdVacation, loc.CreateString("entryTypeVacation")),
		model.NewEntryType(model.EntryTypeIdHoliday, loc.CreateString("entryTypeHoliday")),
		model.NewEntryType(model.EntryTypeIdIllness, loc.CreateString("entryTypeIllness")),
	}, nil
}

// GetEntryTypesMap gets a map of all entry types.
func (s *EntryService) GetEntryTypesMap(ctx context.Context) (map[int]*model.EntryType, *e.Error) {
	// Get entry types
	entryTypes, err := s.GetEntryTypes(ctx)
	if err != nil {
		return nil, err
	}

	// Convert into map
	m := make(map[int]*model.EntryType)
	for _, entryType := range entryTypes {
		m[entryType.Id] = entryType
	}

	return m, nil
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
func (s *EntryService) GetEntryActivities(ctx context.Context) ([]*model.EntryActivity, *e.Error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetEntryCharacts); err != nil {
		return nil, err
	}

	// Get entry activities
	return s.eRepo.GetEntryActivities(ctx)
}

// GetEntryActivitiesMap gets a map of all entry activities.
func (s *EntryService) GetEntryActivitiesMap(ctx context.Context) (map[int]*model.EntryActivity,
	*e.Error) {
	// Get entry activities
	entryActivities, err := s.eRepo.GetEntryActivities(ctx)
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

// CreateEntryActivity creates a new entry activity.
func (s *EntryService) CreateEntryActivity(ctx context.Context,
	entryActivity *model.EntryActivity) *e.Error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeEntryCharacts); err != nil {
		return err
	}

	// Check if entry activity exists
	if err := s.checkIfEntryActivityExists(ctx, entryActivity.Id); err != nil {
		return err
	}

	// Create entry activity
	return s.eRepo.CreateEntryActivity(ctx, entryActivity)
}

// UpdateEntryActivity updates an entry activity.
func (s *EntryService) UpdateEntryActivity(ctx context.Context,
	entryActivity *model.EntryActivity) *e.Error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeEntryCharacts); err != nil {
		return err
	}

	// Check if entry activity exists
	if err := s.checkIfEntryActivityExists(ctx, entryActivity.Id); err != nil {
		return err
	}

	// Update entry activity
	return s.eRepo.UpdateEntryActivity(ctx, entryActivity)
}

// DeleteEntryActivityById deletes an entry activity.
func (s *EntryService) DeleteEntryActivityById(ctx context.Context, id int) *e.Error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeEntryCharacts); err != nil {
		return err
	}

	// Check if entry activity exists
	if err := s.checkIfEntryActivityExists(ctx, id); err != nil {
		return err
	}

	// Check if entries with this activity exist
	if err := s.checkIfEntryActivityIsUsed(ctx, id); err != nil {
		return err
	}

	// Delete entry activity
	return s.eRepo.DeleteEntryActivityById(ctx, id)
}

func (s *EntryService) checkIfEntryActivityExists(ctx context.Context, id int) *e.Error {
	if id == 0 {
		return nil
	}
	exists, err := s.eRepo.ExistsEntryActivityById(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		err = e.NewError(e.LogicEntryActivityNotFound, fmt.Sprintf("Could not find entry activity "+
			"%d.", id))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func (s *EntryService) checkIfEntryActivityIsUsed(ctx context.Context, id int) *e.Error {
	existsEntry, err := s.eRepo.ExistsEntryByActivityId(ctx, id)
	if err != nil {
		return err
	}
	if existsEntry {
		err = e.NewError(e.LogicEntryActivityDeleteNotAllowed, fmt.Sprintf("Could not delete entry "+
			"activity %d. There are still entries for this activity.", id))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

// --- Work summary functions ---

// GetTotalWorkSummaryByUserId gets the total work summary of an user.
func (s *EntryService) GetTotalWorkSummaryByUserId(ctx context.Context, userId int) (
	*model.WorkSummary, *e.Error) {
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, err
	}

	// Get work summary
	start := time.Time{}
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return s.eRepo.GetWorkSummary(ctx, userId, start, end)
}

// --- Permission helper functions ---

func (s *EntryService) checkHasCurrentUserGetRight(ctx context.Context, userId int) *e.Error {
	if userId == getCurrentUserId(ctx) {
		return checkHasCurrentUserRight(ctx, model.RightGetOwnEntries)
	} else {
		return checkHasCurrentUserRight(ctx, model.RightGetAllEntries)
	}
}

func (s *EntryService) checkHasCurrentUserChangeRight(ctx context.Context, userId int) *e.Error {
	if userId == getCurrentUserId(ctx) {
		return checkHasCurrentUserRight(ctx, model.RightChangeOwnEntries)
	} else {
		return checkHasCurrentUserRight(ctx, model.RightChangeAllEntries)
	}
}
