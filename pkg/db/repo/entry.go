package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
)

type dbEntry struct {
	id          int
	userId      int
	typeId      int
	startTime   string
	endTime     string
	activityId  sql.NullInt64
	description sql.NullString
}

type dbWorkDuration struct {
	typeId       int
	workDuration int
}

// EntryRepo retrieves and stores entry related entities.
type EntryRepo struct {
	repo
}

// NewEntryRepo creates a new entry repository.
func NewEntryRepo(db *sql.DB) *EntryRepo {
	return &EntryRepo{repo{db}}
}

// --- Entry functions ---

// CountDateEntries counts entries (over date).
func (r *EntryRepo) CountDateEntries(ctx context.Context, filter *model.EntriesFilter) (int,
	*e.Error) {
	q, qa := r.buildCountDateEntriesQuery(filter)

	sr, qErr := r.queryRow(ctx, &scanIntHelper{}, q, qa...)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries (over date) in database.",
			qErr)
		log.Error(err.StackTrace())
		return 0, err
	}

	return sr.(int), nil
}

func (r *EntryRepo) buildCountDateEntriesQuery(filter *model.EntriesFilter) (
	string, []interface{}) {
	qr, qra := r.buildEntriesFilterQueryRestriction(filter)

	q := "SELECT COUNT(DISTINCT(DATE(e.start_time))) FROM entry e " + qr

	return q, qra
}

// GetDateEntries retrieves entries (over date).
func (r *EntryRepo) GetDateEntries(ctx context.Context, filter *model.EntriesFilter,
	sort *model.EntriesSort, offset int, limit int) ([]*model.Entry, *e.Error) {
	qr, qra := r.buildGetDateEntriesRangeQuery(filter, sort, offset, limit)

	start, end, qrErr := r.getDateRange(ctx, qr, qra...)
	if qrErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query range for entries (over date) from "+
			"database.", qrErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	if start == "" || end == "" {
		return make([]*model.Entry, 0), nil
	}

	q, qa := r.buildGetDateEntriesQuery(filter, sort, start, end)

	sr, qErr := r.query(ctx, &scanEntryHelper{}, q, qa...)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query entries (over date) from database.",
			qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	entries := sr.([]*model.Entry)

	return entries, nil
}

func (r *EntryRepo) buildGetDateEntriesRangeQuery(filter *model.EntriesFilter,
	sort *model.EntriesSort, offset int, limit int) (string, []interface{}) {
	qr, qra := r.buildEntriesFilterQueryRestriction(filter)
	var qo string
	if sort != nil && (sort.ByTime == model.NoSorting || sort.ByTime == model.AscSorting) {
		qo = "ORDER BY date ASC"
	} else if sort != nil && sort.ByTime == model.DescSorting {
		qo = "ORDER BY date DESC"
	}

	q := "SELECT DISTINCT(DATE(e.start_time)) AS date " +
		"FROM entry e " +
		qr + " " +
		qo + " " +
		createQueryLimitString(offset, limit)

	return q, qra
}

func (r *EntryRepo) buildGetDateEntriesQuery(filter *model.EntriesFilter, sort *model.EntriesSort,
	start string, end string) (string, []interface{}) {
	qr, qra := r.buildEntriesFilterQueryRestriction(filter)
	qo := r.buildEntriesSortQueryClause(sort)

	q := "SELECT e.id, e.user_id, e.type_id, e.start_time, e.end_time, e.activity_id, " +
		"e.description " +
		"FROM entry e " +
		qr + " " +
		"AND e.start_time BETWEEN ? AND ? " +
		qo

	qa := append(qra, start, end)

	return q, qa
}

// CountDateEntriesByUserId counts all entries (over date) of an user.
func (r *EntryRepo) CountDateEntriesByUserId(ctx context.Context, userId int) (int, *e.Error) {
	q, qa := r.buildCountDateEntriesByUserIdQuery(userId)

	sr, qErr := r.queryRow(ctx, &scanIntHelper{}, q, qa...)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries (over date) in database.",
			qErr)
		log.Error(err.StackTrace())
		return 0, err
	}

	return sr.(int), nil
}

func (r *EntryRepo) buildCountDateEntriesByUserIdQuery(userId int) (string, []interface{}) {
	q := "SELECT COUNT(DISTINCT(DATE(e.start_time))) " +
		"FROM entry e " +
		"WHERE e.user_id = ?"

	qa := []interface{}{userId}

	return q, qa
}

// GetDateEntriesByUserId retrieves all entries (over date) of an user.
func (r *EntryRepo) GetDateEntriesByUserId(ctx context.Context, userId int, offset int, limit int) (
	[]*model.Entry, *e.Error) {
	qr, qra := r.buildGetDateEntriesByUserIdRangeQuery(userId, offset, limit)

	start, end, qrErr := r.getDateRange(ctx, qr, qra...)
	if qrErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query range for entries (over date) "+
			"from database.", qrErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	if start == "" || end == "" {
		return make([]*model.Entry, 0), nil
	}

	q, qa := r.buildGetDateEntriesByUserIdQuery(userId, start, end)

	sr, qErr := r.query(ctx, &scanEntryHelper{}, q, qa...)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query entries (over date) from database.",
			qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	entries := sr.([]*model.Entry)

	return entries, nil
}

func (r *EntryRepo) buildGetDateEntriesByUserIdRangeQuery(userId int, offset int, limit int) (string,
	[]interface{}) {
	q := "SELECT DISTINCT(DATE(e.start_time)) AS date " +
		"FROM entry e " +
		"WHERE e.user_id = ? " +
		"ORDER BY date DESC " +
		createQueryLimitString(offset, limit)

	qa := []interface{}{userId}

	return q, qa
}

func (r *EntryRepo) buildGetDateEntriesByUserIdQuery(userId int, start string, end string) (string,
	[]interface{}) {
	q := "SELECT e.id, e.user_id, e.type_id, e.start_time, e.end_time, e.activity_id, " +
		"e.description " +
		"FROM entry e " +
		"WHERE e.user_id = ? " +
		"AND e.start_time BETWEEN ? AND ? " +
		"ORDER BY e.start_time DESC, e.end_time DESC"

	qa := []interface{}{userId, start, end}

	return q, qa
}

// CountEntries counts all entries.
func (r *EntryRepo) CountEntries(ctx context.Context, filter *model.EntriesFilter) (int, *e.Error) {
	qr, qra := r.buildEntriesFilterQueryRestriction(filter)

	q := "SELECT COUNT(*) FROM entry e " + qr

	var sr interface{}
	var qErr error
	if len(qra) > 0 {
		sr, qErr = r.queryRow(ctx, &scanIntHelper{}, q, qra)
	} else {
		sr, qErr = r.queryRow(ctx, &scanIntHelper{}, q)
	}
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries in database.", qErr)
		log.Error(err.StackTrace())
		return 0, err
	}

	return sr.(int), nil
}

// GetEntries retrieves all entries.
func (r *EntryRepo) GetEntries(ctx context.Context, filter *model.EntriesFilter,
	sort *model.EntriesSort, offset int, limit int) ([]*model.Entry, *e.Error) {
	qr, qra := r.buildEntriesFilterQueryRestriction(filter)
	qo := r.buildEntriesSortQueryClause(sort)

	q := "SELECT e.id, e.user_id, e.type_id, e.start_time, e.end_time, e.activity_id, " +
		"e.description " +
		"FROM entry e " +
		qr + " " +
		qo + " " +
		createQueryLimitString(offset, limit)

	var sr interface{}
	var qErr error
	if len(qra) > 0 {
		sr, qErr = r.query(ctx, &scanEntryHelper{}, q, qra)
	} else {
		sr, qErr = r.query(ctx, &scanEntryHelper{}, q)
	}
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query entries from database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	entries := sr.([]*model.Entry)

	return entries, nil
}

// GetEntryById retrieves an entry.
func (r *EntryRepo) GetEntryById(ctx context.Context, id int) (*model.Entry, *e.Error) {
	q := "SELECT id, user_id, type_id, start_time, end_time, activity_id, description " +
		"FROM entry WHERE id = ?"

	sr, qErr := r.queryRow(ctx, &scanEntryHelper{}, q, id)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read entry %d from database.",
			id), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	if sr == nil {
		return nil, nil
	}
	entry := sr.(*model.Entry)

	return entry, nil
}

// GetEntryByIdAndUserId retrieves an entry of an user.
func (r *EntryRepo) GetEntryByIdAndUserId(ctx context.Context, id int, userId int) (*model.Entry,
	*e.Error) {
	q := "SELECT id, user_id, type_id, start_time, end_time, activity_id, description " +
		"FROM entry WHERE id = ? AND user_id = ?"

	sr, qErr := r.queryRow(ctx, &scanEntryHelper{}, q, id, userId)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read entry %d from database.",
			id), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	if sr == nil {
		return nil, nil
	}
	entry := sr.(*model.Entry)

	return entry, nil
}

// ExistsEntryById checks if a entry exists.
func (r *EntryRepo) ExistsEntryById(ctx context.Context, id int) (bool, *e.Error) {
	cnt, cErr := r.count(ctx, "entry", "id = ?", id)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries from database.", cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// ExistsEntryByIdAndUserId checks if a entry exists for an user.
func (r *EntryRepo) ExistsEntryByIdAndUserId(ctx context.Context, id int, userId int) (bool,
	*e.Error) {
	cnt, cErr := r.count(ctx, "entry", "id = ? AND user_id = ?", id, userId)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries from database.", cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// ExistsEntryByActivityId checks if a entry exists for an activity.
func (r *EntryRepo) ExistsEntryByActivityId(ctx context.Context, activityId int) (bool, *e.Error) {
	cnt, cErr := r.count(ctx, "entry", "activity_id = ?", activityId)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries from database.", cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// CreateEntry creates a new entry.
func (r *EntryRepo) CreateEntry(ctx context.Context, entry *model.Entry) *e.Error {
	etr := toDbEntry(entry)

	q := "INSERT INTO entry (user_id, type_id, start_time, end_time, activity_id, description) " +
		"VALUES (?, ?, ?, ?, ?, ?)"

	id, cErr := r.insert(ctx, q, etr.userId, etr.typeId, etr.startTime, etr.endTime,
		etr.activityId, etr.description)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not create entry in database.", cErr)
		log.Error(err.StackTrace())
		return err
	}

	entry.Id = id

	return nil
}

// UpdateEntry updates a entry.
func (r *EntryRepo) UpdateEntry(ctx context.Context, entry *model.Entry) *e.Error {
	etr := toDbEntry(entry)

	q := "UPDATE entry SET user_id = ?, type_id = ?, start_time = ?, end_time = ?, " +
		"activity_id = ?, description = ? WHERE id = ?"

	uErr := r.exec(ctx, q, etr.userId, etr.typeId, etr.startTime, etr.endTime, etr.activityId,
		etr.description, etr.id)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update entry %d in database.",
			entry.Id), uErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// DeleteEntryById deletes a entry.
func (r *EntryRepo) DeleteEntryById(ctx context.Context, id int) *e.Error {
	q := "DELETE FROM entry WHERE id = ?"

	dErr := r.exec(ctx, q, id)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not delete entry %d from "+
			"database.", id), dErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// GetMonthEntries retrieves all entries of a month.
func (r *EntryRepo) GetMonthEntries(ctx context.Context, userId int, year int, month int) (
	[]*model.Entry, *e.Error) {
	q := "SELECT id, user_id, type_id, start_time, end_time, activity_id, description " +
		"FROM entry " +
		"WHERE user_id = ? " +
		"AND YEAR(start_time) = ? AND MONTH(start_time) = ? " +
		"ORDER BY start_time ASC, end_time ASC"

	sr, qErr := r.query(ctx, &scanEntryHelper{}, q, userId, year, month)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query month entries from database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	entries := sr.([]*model.Entry)

	return entries, nil
}

// --- Entry activity functions ---

// GetEntryActivities retrieves all entry activities.
func (r *EntryRepo) GetEntryActivities(ctx context.Context) ([]*model.EntryActivity, *e.Error) {
	q := "SELECT id, description FROM entry_activity"

	sr, qErr := r.query(ctx, &scanEntryActivityHelper{}, q)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query entry activities from database.",
			qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	return sr.([]*model.EntryActivity), nil
}

// GetEntryActivityByDescription retrieves a entry activity by its description.
func (r *EntryRepo) GetEntryActivityByDescription(ctx context.Context, description string) (
	*model.EntryActivity, *e.Error) {
	q := "SELECT id, description FROM entry_activity WHERE description = ?"

	sr, qErr := r.queryRow(ctx, &scanEntryActivityHelper{}, q, description)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not query entry activity '%s' "+
			"from database.", description), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	if sr == nil {
		return nil, nil
	}
	return sr.(*model.EntryActivity), nil
}

// ExistsEntryActivityById checks if a entry activity exists.
func (r *EntryRepo) ExistsEntryActivityById(ctx context.Context, id int) (bool, *e.Error) {
	cnt, cErr := r.count(ctx, "entry_activity", "id = ?", id)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entry activities in database.", cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// CreateEntryActivity creates a new entry activity.
func (r *EntryRepo) CreateEntryActivity(ctx context.Context,
	entryActivity *model.EntryActivity) *e.Error {
	q := "INSERT INTO entry_activity (description) VALUES (?)"

	id, cErr := r.insert(ctx, q, entryActivity.Description)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not create entry activity in database.", cErr)
		log.Error(err.StackTrace())
		return err
	}

	entryActivity.Id = id

	return nil
}

// UpdateEntryActivity updates a entry activity.
func (r *EntryRepo) UpdateEntryActivity(ctx context.Context,
	entryActivity *model.EntryActivity) *e.Error {
	q := "UPDATE entry_activity SET description = ? WHERE id = ?"

	uErr := r.exec(ctx, q, entryActivity.Description, entryActivity.Id)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update entry activity %d "+
			"in database.", entryActivity.Id), uErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// DeleteEntryActivityById deletes a entry activity.
func (r *EntryRepo) DeleteEntryActivityById(ctx context.Context, id int) *e.Error {
	q := "DELETE FROM entry_activity WHERE id = ?"

	dErr := r.exec(ctx, q, id)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not delete entry activity %d "+
			"from database.", id), dErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// --- Work summary functions ---

// GetWorkSummary gets the work summary for a specific period.
func (r *EntryRepo) GetWorkSummary(ctx context.Context, userId int, start time.Time, end time.Time) (
	*model.WorkSummary,
	*e.Error) {
	q := "SELECT type_id, SUM(TIMESTAMPDIFF(MINUTE, start_time , end_time)) " +
		"FROM entry " +
		"WHERE user_id = ? " +
		"AND start_time >= ? AND end_time <= ? " +
		"GROUP BY type_id"

	sr, qErr := r.query(ctx, &scanWorkDurationHelper{}, q, userId, *formatTimestamp(&start),
		*formatTimestamp(&end))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query work durations from database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	workDurations := sr.([]*model.WorkDuration)

	workSummary := model.NewWorkSummary()
	workSummary.UserId = userId
	workSummary.StartTime = start
	workSummary.EndTime = end
	workSummary.WorkDurations = workDurations

	return workSummary, nil
}

// --- Filter helper functions ---

func (r *EntryRepo) buildEntriesFilterQueryRestriction(filter *model.EntriesFilter) (
	string, []interface{}) {
	var qrs []string
	var qas []interface{}
	if filter == nil {
		return "", qas
	}

	if filter.ByUser {
		qrs = append(qrs, fmt.Sprintf("e.user_id = %d", filter.UserId))
	}
	if filter.ByType {
		qrs = append(qrs, fmt.Sprintf("e.type_id = %d", filter.TypeId))
	}
	if filter.ByTime {
		qrs = append(qrs, fmt.Sprintf("(e.start_time BETWEEN '%s' AND '%s')",
			*formatTimestamp(&filter.StartTime), *formatTimestamp(&filter.EndTime)))
	}
	if filter.ByActivity {
		if filter.ActivityId == 0 {
			qrs = append(qrs, "e.activity_id IS NULL")
		} else {
			qrs = append(qrs, fmt.Sprintf("e.activity_id = %d", filter.ActivityId))
		}
	}
	if filter.ByDescription {
		if filter.Description == "" {
			qrs = append(qrs, "e.description IS NULL")
		} else {
			qrs = append(qrs, "e.description LIKE ?")
			qas = append(qas, "%"+escapeRestrictionString(filter.Description)+"%")
		}
	}

	qr := ""
	if len(qrs) > 0 {
		qr = "WHERE " + strings.Join(qrs[:], " AND ")
	}

	return qr, qas
}

// --- Sort ---

func (r *EntryRepo) buildEntriesSortQueryClause(sort *model.EntriesSort) string {
	if sort == nil {
		return ""
	}

	if sort.ByTime == model.NoSorting || sort.ByTime == model.AscSorting {
		return "ORDER BY e.start_time ASC, e.end_time ASC"
	} else {
		return "ORDER BY e.start_time DESC, e.end_time DESC"
	}
}

// --- Date range helper functions ---

func (r *EntryRepo) getDateRange(ctx context.Context, query string, args ...interface{}) (
	string, string, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	noRows := true
	min := time.Date(9999, time.January, 1, 0, 0, 0, 0, time.Local)
	max := time.Date(1000, time.January, 1, 0, 0, 0, 0, time.Local)
	for rows.Next() {
		noRows = false
		var s string
		if err := rows.Scan(&s); err != nil {
			return "", "", err
		}
		d := *parseDate(&s)
		if d.Before(min) {
			min = d
		}
		if d.After(max) {
			max = d
		}
	}
	if err := rows.Err(); err != nil {
		return "", "", err
	}
	if noRows {
		return "", "", nil
	}

	start := *formatDate(&min)
	end := *formatDate(&max)

	return start + " 00:00:00", end + " 23:59:59", nil
}

// --- Helper functions ---

type scanEntryHelper struct {
}

func (h *scanEntryHelper) makeSlice() interface{} {
	return make([]*model.Entry, 0, 100)
}

func (h *scanEntryHelper) scan(s scanner) (interface{}, error) {
	var dbE dbEntry

	err := s.Scan(&dbE.id, &dbE.userId, &dbE.typeId, &dbE.startTime, &dbE.endTime, &dbE.activityId,
		&dbE.description)
	if err != nil {
		return nil, err
	}

	entry := fromDbEntry(&dbE)

	return entry, nil
}

func (h *scanEntryHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]*model.Entry), item.(*model.Entry))
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

type scanWorkDurationHelper struct {
}

func (h *scanWorkDurationHelper) makeSlice() interface{} {
	return make([]*model.WorkDuration, 0, 10)
}

func (h *scanWorkDurationHelper) scan(s scanner) (interface{}, error) {
	var dbWd dbWorkDuration

	err := s.Scan(&dbWd.typeId, &dbWd.workDuration)
	if err != nil {
		return nil, err
	}

	workDuration := fromDbWorkDuration(&dbWd)

	return workDuration, nil
}

func (h *scanWorkDurationHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]*model.WorkDuration), item.(*model.WorkDuration))
}

func toDbEntry(in *model.Entry) *dbEntry {
	var out dbEntry
	out.id = in.Id
	out.userId = in.UserId
	out.typeId = in.TypeId
	out.startTime = *formatTimestamp(&in.StartTime)
	out.endTime = *formatTimestamp(&in.EndTime)
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

func toDbWorkDuration(in *model.WorkDuration) *dbWorkDuration {
	var out dbWorkDuration
	out.typeId = in.TypeId
	out.workDuration = *formatDuration(&in.WorkDuration)
	return &out
}

func fromDbWorkDuration(in *dbWorkDuration) *model.WorkDuration {
	var out model.WorkDuration
	out.TypeId = in.typeId
	out.WorkDuration = *parseDuration(&in.workDuration)
	return &out
}
