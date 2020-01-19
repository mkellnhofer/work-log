package repo

import (
	"database/sql"
	"fmt"

	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
)

type dbEntry struct {
	id            int
	userId        int
	typeId        int
	startTime     string
	endTime       string
	breakDuration int
	activityId    sql.NullInt64
	description   sql.NullString
}

// EntryRepo retrieves and stores work entry and work entry type records.
type EntryRepo struct {
	repo
}

// NewEntryRepo creates a new work entry repository.
func NewEntryRepo(db *sql.DB) *EntryRepo {
	return &EntryRepo{repo{db}}
}

// --- Work entry functions ---

// CountDateEntries counts all work entries (over date).
func (r *EntryRepo) CountDateEntries(userId int) (int, *e.Error) {
	q := "SELECT COUNT(DISTINCT(DATE(start_time))) FROM entry WHERE user_id = ?"

	sr, qErr := r.queryRow(&scanIntHelper{}, q, userId)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count work entries (over date) in database.",
			qErr)
		log.Error(err.StackTrace())
		return 0, err
	}

	return sr.(int), nil
}

// GetDateEntries retrieves all work entries (over date).
func (r *EntryRepo) GetDateEntries(userId int, offset int, limit int) ([]*model.Entry, *e.Error) {
	ql := createQueryLimitString(offset, limit)

	dq := "SELECT DISTINCT(DATE(start_time)) AS date " +
		"FROM entry " +
		"ORDER BY date DESC" +
		ql

	q := "SELECT e.id, e.user_id, e.type_id, e.start_time, e.end_time, e.break_duration, " +
		"e.activity_id, e.description " +
		"FROM (SELECT id, user_id, type_id, DATE(start_time) AS date, start_time, end_time, " +
		"break_duration, activity_id, description FROM entry) e " +
		"INNER JOIN (" + dq + ") AS d ON e.date = d.date " +
		"WHERE user_id = ? " +
		"ORDER BY e.date DESC, e.start_time DESC, e.end_time DESC"

	sr, qErr := r.query(&scanEntryHelper{}, q, userId)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query work entries (over date) from "+
			"database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	entries := sr.([]*model.Entry)

	return entries, nil
}

// CountEntries counts all work entries.
func (r *EntryRepo) CountEntries(userId int) (int, *e.Error) {
	q := "SELECT COUNT(*) FROM entry WHERE user_id = ?"

	sr, qErr := r.queryRow(&scanIntHelper{}, q, userId)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count work entries in database.", qErr)
		log.Error(err.StackTrace())
		return 0, err
	}

	return sr.(int), nil
}

// GetEntries retrieves all work entries.
func (r *EntryRepo) GetEntries(userId int, offset int, limit int) ([]*model.Entry, *e.Error) {
	q := "SELECT id, user_id, type_id, start_time, end_time, break_duration, activity_id, " +
		"description FROM entry WHERE user_id = ? ORDER BY start_time DESC, end_time DESC"

	ql := createQueryLimitString(offset, limit)

	sr, qErr := r.query(&scanEntryHelper{}, q+ql, userId)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query work entries from database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	entries := sr.([]*model.Entry)

	return entries, nil
}

// GetEntryById retrieves a work entry by its ID.
func (r *EntryRepo) GetEntryById(id int, userId int) (*model.Entry, *e.Error) {
	q := "SELECT id, user_id, type_id, start_time, end_time, break_duration, activity_id, " +
		"description FROM entry WHERE id = ? AND user_id = ?"

	sr, qErr := r.queryRow(&scanEntryHelper{}, q, id, userId)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read work entry %d from "+
			"database.", id), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	if sr == nil {
		return nil, nil
	}
	entry := sr.(*model.Entry)

	return entry, nil
}

// ExistsEntryById checks if a work entry exists.
func (r *EntryRepo) ExistsEntryById(id int, userId int) (bool, *e.Error) {
	cnt, cErr := r.count("entry", "id = ? AND user_id = ?", id, userId)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read work entry %d from "+
			"database.", id), cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// CreateEntry creates a new work entry.
func (r *EntryRepo) CreateEntry(entry *model.Entry) *e.Error {
	etr := toDbEntry(entry)

	q := "INSERT INTO entry (user_id, type_id, start_time, end_time, break_duration, activity_id, " +
		"description) VALUES (?, ?, ?, ?, ?, ?, ?)"

	id, cErr := r.insert(q, etr.userId, etr.typeId, etr.startTime, etr.endTime, etr.breakDuration,
		etr.activityId, etr.description)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not create work entry in database.", cErr)
		log.Error(err.StackTrace())
		return err
	}

	entry.Id = id

	return nil
}

// UpdateEntry updates a work entry.
func (r *EntryRepo) UpdateEntry(entry *model.Entry) *e.Error {
	etr := toDbEntry(entry)

	q := "UPDATE entry SET user_id = ?, type_id = ?, start_time = ?, end_time = ?, " +
		"break_duration = ?, activity_id = ?, description = ? WHERE id = ?"

	uErr := r.exec(q, etr.userId, etr.typeId, etr.startTime, etr.endTime, etr.breakDuration,
		etr.activityId, etr.description, etr.id)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update work entry %d in "+
			"database.", entry.Id), uErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// DeleteEntryById deletes a work entry by by its ID.
func (r *EntryRepo) DeleteEntryById(id int) *e.Error {
	q := "DELETE FROM entry WHERE id = ?"

	dErr := r.exec(q, id)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not delete work entry %d from "+
			"database.", id), dErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// --- Work entry type functions ---

// GetEntryTypes retrieves all work entry types.
func (r *EntryRepo) GetEntryTypes() ([]*model.EntryType, *e.Error) {
	q := "SELECT id, description FROM entry_type"

	sr, qErr := r.query(&scanEntryTypeHelper{}, q)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query work entry types from database.",
			qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	return sr.([]*model.EntryType), nil
}

// GetEntryTypeByDescription retrieves a work entry type by its description.
func (r *EntryRepo) GetEntryTypeByDescription(description string) (*model.EntryType, *e.Error) {
	q := "SELECT id, description FROM entry_type WHERE description = ?"

	sr, qErr := r.queryRow(&scanEntryTypeHelper{}, q, description)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not query work entry type '%s' "+
			"from database.", description), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	if sr == nil {
		return nil, nil
	}
	return sr.(*model.EntryType), nil
}

// ExistsEntryTypeById checks if a work entry type exists.
func (r *EntryRepo) ExistsEntryTypeById(id int) (bool, *e.Error) {
	cnt, cErr := r.count("entry_type", "id = ?", id)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read work entry type %d "+
			"from database.", id), cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// CreateEntryType creates a new work entry type.
func (r *EntryRepo) CreateEntryType(entryType *model.EntryType) *e.Error {
	q := "INSERT INTO entry_type (description) VALUES (?)"

	id, cErr := r.insert(q, entryType.Description)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not create work entry type in database.", cErr)
		log.Error(err.StackTrace())
		return err
	}

	entryType.Id = id

	return nil
}

// UpdateEntryType updates a work entry type.
func (r *EntryRepo) UpdateEntryType(entryType *model.EntryType) *e.Error {
	q := "UPDATE entry_type SET description = ? WHERE id = ?"

	uErr := r.exec(q, entryType.Description, entryType.Id)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update work entry type %d "+
			"in database.", entryType.Id), uErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// DeleteEntryTypeById deletes a work entry type by its ID.
func (r *EntryRepo) DeleteEntryTypeById(id int) *e.Error {
	q := "DELETE FROM entry_type WHERE id = ?"

	dErr := r.exec(q, id)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not delete work entry type %d "+
			"from database.", id), dErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// --- Work entry activity functions ---

// GetEntryActivities retrieves all work entry activities.
func (r *EntryRepo) GetEntryActivities() ([]*model.EntryActivity, *e.Error) {
	q := "SELECT id, description FROM entry_activity"

	sr, qErr := r.query(&scanEntryActivityHelper{}, q)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query work entry activities from database.",
			qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	return sr.([]*model.EntryActivity), nil
}

// GetEntryActivityByDescription retrieves a work entry activity by its description.
func (r *EntryRepo) GetEntryActivityByDescription(description string) (*model.EntryActivity,
	*e.Error) {
	q := "SELECT id, description FROM entry_activity WHERE description = ?"

	sr, qErr := r.queryRow(&scanEntryActivityHelper{}, q, description)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not query work entry activity "+
			"'%s' from database.", description), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	if sr == nil {
		return nil, nil
	}
	return sr.(*model.EntryActivity), nil
}

// ExistsEntryActivityById checks if a work entry activity exists.
func (r *EntryRepo) ExistsEntryActivityById(id int) (bool, *e.Error) {
	cnt, cErr := r.count("entry_activity", "id = ?", id)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read work entry activity %d "+
			"from database.", id), cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// CreateEntryActivity creates a new work entry activity.
func (r *EntryRepo) CreateEntryActivity(entryActivity *model.EntryActivity) *e.Error {
	q := "INSERT INTO entry_activity (description) VALUES (?)"

	id, cErr := r.insert(q, entryActivity.Description)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not create work entry activity in database.",
			cErr)
		log.Error(err.StackTrace())
		return err
	}

	entryActivity.Id = id

	return nil
}

// UpdateEntryActivity updates a work entry activity.
func (r *EntryRepo) UpdateEntryActivity(entryActivity *model.EntryActivity) *e.Error {
	q := "UPDATE entry_activity SET description = ? WHERE id = ?"

	uErr := r.exec(q, entryActivity.Description, entryActivity.Id)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update work entry activity "+
			"%d in database.", entryActivity.Id), uErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// DeleteEntryActivityById deletes a work entry activity by its ID.
func (r *EntryRepo) DeleteEntryActivityById(id int) *e.Error {
	q := "DELETE FROM entry_activity WHERE id = ?"

	dErr := r.exec(q, id)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not delete work entry activity "+
			"%d from database.", id), dErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// --- Helper functions ---

type scanEntryHelper struct {
}

func (h *scanEntryHelper) makeSlice() interface{} {
	return make([]*model.Entry, 0, 100)
}

func (h *scanEntryHelper) scan(s scanner) (interface{}, error) {
	var dbE dbEntry

	err := s.Scan(&dbE.id, &dbE.userId, &dbE.typeId, &dbE.startTime, &dbE.endTime, &dbE.breakDuration,
		&dbE.activityId, &dbE.description)
	if err != nil {
		return nil, err
	}

	entry := fromDbEntry(&dbE)

	return entry, nil
}

func (h *scanEntryHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]*model.Entry), item.(*model.Entry))
}

type scanEntryTypeHelper struct {
}

func (h *scanEntryTypeHelper) makeSlice() interface{} {
	return make([]*model.EntryType, 0, 10)
}

func (h *scanEntryTypeHelper) scan(s scanner) (interface{}, error) {
	var et model.EntryType

	err := s.Scan(&et.Id, &et.Description)
	if err != nil {
		return nil, err
	}

	return &et, nil
}

func (h *scanEntryTypeHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]*model.EntryType), item.(*model.EntryType))
}

type scanEntryActivityHelper struct {
}

func (h *scanEntryActivityHelper) makeSlice() interface{} {
	return make([]*model.EntryActivity, 0, 10)
}

func (h *scanEntryActivityHelper) scan(s scanner) (interface{}, error) {
	var et model.EntryActivity

	err := s.Scan(&et.Id, &et.Description)
	if err != nil {
		return nil, err
	}

	return &et, nil
}

func (h *scanEntryActivityHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]*model.EntryActivity), item.(*model.EntryActivity))
}

func toDbEntry(in *model.Entry) *dbEntry {
	var out dbEntry
	out.id = in.Id
	out.userId = in.UserId
	out.typeId = in.TypeId
	out.startTime = *formatTimestamp(&in.StartTime)
	out.endTime = *formatTimestamp(&in.EndTime)
	out.breakDuration = *formatDuration(&in.BreakDuration)
	if in.ActivityId != 0 {
		out.activityId = sql.NullInt64{Int64: int64(in.ActivityId), Valid: true}
	} else {
		out.activityId = sql.NullInt64{Int64: 0, Valid: false}
	}
	if in.Description != "" {
		out.description = sql.NullString{String: in.Description, Valid: true}
	} else {
		out.description = sql.NullString{String: "", Valid: false}
	}
	return &out
}

func fromDbEntry(in *dbEntry) *model.Entry {
	var out model.Entry
	out.Id = in.id
	out.UserId = in.userId
	out.TypeId = in.typeId
	out.StartTime = *parseTimestamp(&in.startTime)
	out.EndTime = *parseTimestamp(&in.endTime)
	out.BreakDuration = *parseDuration(&in.breakDuration)
	if in.activityId.Valid {
		out.ActivityId = int(in.activityId.Int64)
	} else {
		out.ActivityId = 0
	}
	if in.description.Valid {
		out.Description = in.description.String
	} else {
		out.Description = ""
	}
	return &out
}
