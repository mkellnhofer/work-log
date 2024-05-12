package service

import (
	"context"
	"fmt"
	"time"

	"kellnhofer.com/work-log/pkg/db/repo"
	"kellnhofer.com/work-log/pkg/db/tx"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
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
func (s *EntryService) GetDateEntries(ctx context.Context, filter *model.EntriesFilter,
	sort *model.EntriesSort, offset int, limit int) ([]*model.Entry, int, error) {
	// If user does not have right to get any entry: Add default user ID filter
	if !hasCurrentUserRight(ctx, model.RightGetAllEntries) && !filter.ByUser {
		filter.ByUser = true
		filter.UserId = getCurrentUserId(ctx)
	}

	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, filter.UserId); err != nil {
		return nil, 0, err
	}

	// Get entries
	entries, err := s.eRepo.GetDateEntries(ctx, filter, sort, offset, limit)
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
	limit int) ([]*model.Entry, int, error) {
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
func (s *EntryService) GetEntries(ctx context.Context, filter *model.EntriesFilter,
	sort *model.EntriesSort, offset int, limit int) ([]*model.Entry, int, error) {
	// If user does not have right to get any entry: Add default user ID filter
	if !hasCurrentUserRight(ctx, model.RightGetAllEntries) && !filter.ByUser {
		filter.ByUser = true
		filter.UserId = getCurrentUserId(ctx)
	}

	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, filter.UserId); err != nil {
		return nil, 0, err
	}

	// Get entries
	entries, err := s.eRepo.GetEntries(ctx, filter, sort, offset, limit)
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
func (s *EntryService) GetEntryById(ctx context.Context, id int) (*model.Entry, error) {
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
	error) {
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, err
	}

	// Get entry
	return s.eRepo.GetEntryByIdAndUserId(ctx, id, userId)
}

// CreateEntry creates a new entry.
func (s *EntryService) CreateEntry(ctx context.Context, entry *model.Entry) error {
	// Check permissions
	if err := s.checkHasCurrentUserChangeRight(ctx, entry.UserId); err != nil {
		return err
	}

	// Check if entry type exists
	if err := s.checkEntryTypeExists(entry.TypeId); err != nil {
		return err
	}
	// Check if entry activity exists
	if err := s.checkEntryActivityExistsAllowed(ctx, entry.TypeId, entry.ActivityId); err != nil {
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
func (s *EntryService) UpdateEntry(ctx context.Context, entry *model.Entry) error {
	// Get existing entry
	existingEntry, err := s.eRepo.GetEntryByIdAndUserId(ctx, entry.Id, entry.UserId)
	if err != nil {
		return err
	}

	// Check if entry exists
	if err := s.checkEntryExists(entry.Id, existingEntry); err != nil {
		return err
	}

	// Check permissions
	if err := s.checkHasCurrentUserChangeRight(ctx, existingEntry.UserId); err != nil {
		return err
	}

	// Check if entry type exists
	if err := s.checkEntryTypeExists(entry.TypeId); err != nil {
		return err
	}
	// Check if entry activity exists
	if err := s.checkEntryActivityExistsAllowed(ctx, entry.TypeId, entry.ActivityId); err != nil {
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
func (s *EntryService) DeleteEntryById(ctx context.Context, id int) error {
	// Get existing entry
	existingEntry, err := s.eRepo.GetEntryById(ctx, id)
	if err != nil {
		return err
	}

	// Check if entry exists
	if err := s.checkEntryExists(id, existingEntry); err != nil {
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
func (s *EntryService) DeleteEntryByIdAndUserId(ctx context.Context, id int, userId int) error {
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
	if err := s.checkEntryExists(id, existingEntry); err != nil {
		return err
	}

	// Delete entry
	return s.eRepo.DeleteEntryById(ctx, id)
}

// GetMonthEntriesByUserId gets all entries of a month of an user.
func (s *EntryService) GetMonthEntriesByUserId(ctx context.Context, userId int, year int,
	month int) ([]*model.Entry, error) {
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, err
	}

	// Get entries
	return s.eRepo.GetMonthEntries(ctx, userId, year, month)
}

func (s *EntryService) checkEntryExists(id int, entry *model.Entry) error {
	if entry == nil {
		err := e.NewError(e.LogicEntryNotFound, fmt.Sprintf("Could not find entry %d.", id))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func (s *EntryService) checkEntry(entry *model.Entry) error {
	if entry.StartTime.After(entry.EndTime) {
		err := e.NewError(e.LogicEntryTimeIntervalInvalid, fmt.Sprintf("End time %s before "+
			"start time %s.", entry.EndTime, entry.StartTime))
		log.Debug(err.StackTrace())
		return err
	}

	return nil
}

// --- Entry type functions ---

// GetEntryTypes gets all entry types.
func (s *EntryService) GetEntryTypes(ctx context.Context) ([]*model.EntryType, error) {
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
func (s *EntryService) GetEntryTypesMap(ctx context.Context) (map[int]*model.EntryType, error) {
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

func (s *EntryService) checkEntryTypeExists(typeId int) error {
	exist := typeId == model.EntryTypeIdWork || typeId == model.EntryTypeIdTravel ||
		typeId == model.EntryTypeIdVacation || typeId == model.EntryTypeIdHoliday ||
		typeId == model.EntryTypeIdIllness
	if !exist {
		err := e.NewError(e.LogicEntryTypeNotFound, fmt.Sprintf("Could not find entry type %d.",
			typeId))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

// --- Entry activity functions ---

// GetEntryActivities gets all entry activities.
func (s *EntryService) GetEntryActivities(ctx context.Context) ([]*model.EntryActivity, error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetEntryCharacts); err != nil {
		return nil, err
	}

	// Get entry activities
	return s.eRepo.GetEntryActivities(ctx)
}

// GetEntryActivitiesMap gets a map of all entry activities.
func (s *EntryService) GetEntryActivitiesMap(ctx context.Context) (map[int]*model.EntryActivity,
	error) {
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
func (s *EntryService) CreateEntryActivity(ctx context.Context, entryActivity *model.EntryActivity,
) error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeEntryCharacts); err != nil {
		return err
	}

	// Check if entry activity exists
	if err := s.checkEntryActivityExists(ctx, entryActivity.Id); err != nil {
		return err
	}

	// Create entry activity
	return s.eRepo.CreateEntryActivity(ctx, entryActivity)
}

// UpdateEntryActivity updates an entry activity.
func (s *EntryService) UpdateEntryActivity(ctx context.Context, entryActivity *model.EntryActivity,
) error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeEntryCharacts); err != nil {
		return err
	}

	// Check if entry activity exists
	if err := s.checkEntryActivityExists(ctx, entryActivity.Id); err != nil {
		return err
	}

	// Update entry activity
	return s.eRepo.UpdateEntryActivity(ctx, entryActivity)
}

// DeleteEntryActivityById deletes an entry activity.
func (s *EntryService) DeleteEntryActivityById(ctx context.Context, actId int) error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeEntryCharacts); err != nil {
		return err
	}

	// Check if entry activity exists
	if err := s.checkEntryActivityExists(ctx, actId); err != nil {
		return err
	}

	// Check if entries with this activity exist
	if err := s.checkEntryActivityIsUsed(ctx, actId); err != nil {
		return err
	}

	// Delete entry activity
	return s.eRepo.DeleteEntryActivityById(ctx, actId)
}

func (s *EntryService) checkEntryActivityExistsAllowed(ctx context.Context, typeId int, actId int,
) error {
	if typeId != model.EntryTypeIdWork && actId != 0 {
		err := e.NewError(e.LogicEntryActivityNotAllowed, fmt.Sprintf("The entry activity %d is "+
			"not allowed for the entry type %d.", actId, typeId))
		log.Debug(err.StackTrace())
		return err
	}
	return s.checkEntryActivityExists(ctx, actId)
}

func (s *EntryService) checkEntryActivityExists(ctx context.Context, actId int) error {
	if actId == 0 {
		return nil
	}
	exists, err := s.eRepo.ExistsEntryActivityById(ctx, actId)
	if err != nil {
		return err
	}
	if !exists {
		err := e.NewError(e.LogicEntryActivityNotFound, fmt.Sprintf("Could not find entry activity "+
			"%d.", actId))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func (s *EntryService) checkEntryActivityIsUsed(ctx context.Context, actId int) error {
	existsEntry, err := s.eRepo.ExistsEntryByActivityId(ctx, actId)
	if err != nil {
		return err
	}
	if existsEntry {
		err := e.NewError(e.LogicEntryActivityDeleteNotAllowed, fmt.Sprintf("Could not delete entry "+
			"activity %d. There are still entries for this activity.", actId))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

// --- Work summary functions ---

// GetTotalWorkSummaryByUserId gets the total work summary of an user.
func (s *EntryService) GetTotalWorkSummaryByUserId(ctx context.Context, userId int) (
	*model.WorkSummary, error) {
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, err
	}

	// Get work summary
	start := time.Time{}
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
	return s.eRepo.GetWorkSummary(ctx, userId, start, end)
}

// GetWorkSummaryByUserId gets the month work summary of an user.
func (s *EntryService) GetMonthWorkSummaryByUserId(ctx context.Context, userId int, year int,
	month time.Month) (*model.WorkSummary, error) {
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, err
	}

	// Get work summary
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
	return s.eRepo.GetWorkSummary(ctx, userId, start, end)
}

// --- Permission helper functions ---

func (s *EntryService) checkHasCurrentUserGetRight(ctx context.Context, userId int) error {
	if userId == getCurrentUserId(ctx) {
		return checkHasCurrentUserRight(ctx, model.RightGetOwnEntries)
	} else {
		return checkHasCurrentUserRight(ctx, model.RightGetAllEntries)
	}
}

func (s *EntryService) checkHasCurrentUserChangeRight(ctx context.Context, userId int) error {
	if userId == getCurrentUserId(ctx) {
		return checkHasCurrentUserRight(ctx, model.RightChangeOwnEntries)
	} else {
		return checkHasCurrentUserRight(ctx, model.RightChangeAllEntries)
	}
}
