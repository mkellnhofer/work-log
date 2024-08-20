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

type dbReadEntry struct {
	id          int
	userId      int
	typeId      int
	startTime   string
	endTime     string
	activityId  sql.NullInt64
	project     sql.NullString
	description sql.NullString
	labels      sql.NullString
}

type dbWriteEntry struct {
	id          int
	userId      int
	typeId      int
	startTime   string
	endTime     string
	activityId  sql.NullInt64
	projectId   sql.NullInt64
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
func (r *EntryRepo) CountDateEntries(ctx context.Context, filter *model.FieldEntryFilter) (int,
	error) {
	q, qa := r.buildCountDateEntriesQuery(filter)

	sh := newIntScanHelper()
	count, _, qErr := sh.scanRow(r.queryRow(ctx, q, qa...))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries (over date) in database.",
			qErr)
		log.Error(err.StackTrace())
		return 0, err
	}
	return count, nil
}

func (r *EntryRepo) buildCountDateEntriesQuery(filter *model.FieldEntryFilter) (string, []any) {
	qr, qra := r.buildEntryFilterQueryRestriction(filter)

	q := "SELECT COUNT(DISTINCT(DATE(e.start_time))) " +
		"FROM entry e " +
		"LEFT JOIN project p ON e.project_id = p.id " +
		qr

	return q, qra
}

// GetDateEntries retrieves entries (over date).
func (r *EntryRepo) GetDateEntries(ctx context.Context, filter *model.FieldEntryFilter,
	sort *model.EntrySort, offset int, limit int) ([]*model.Entry, error) {
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

	sh := newEntryScanHelper()
	entries, qErr := sh.scanRows(r.query(ctx, q, qa...))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query entries (over date) from database.",
			qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	return entries, nil
}

func (r *EntryRepo) buildGetDateEntriesRangeQuery(filter *model.FieldEntryFilter,
	sort *model.EntrySort, offset int, limit int) (string, []any) {
	qr, qra := r.buildEntryFilterQueryRestriction(filter)
	var qo string
	if sort != nil && (sort.ByTime == model.NoSorting || sort.ByTime == model.AscSorting) {
		qo = "ORDER BY date ASC"
	} else if sort != nil && sort.ByTime == model.DescSorting {
		qo = "ORDER BY date DESC"
	}

	q := "SELECT DISTINCT(DATE(e.start_time)) AS date " +
		"FROM entry e " +
		"LEFT JOIN project p ON e.project_id = p.id " +
		qr + " " +
		qo + " " +
		createQueryLimitString(offset, limit)

	return q, qra
}

func (r *EntryRepo) buildGetDateEntriesQuery(filter *model.FieldEntryFilter, sort *model.EntrySort,
	start string, end string) (string, []any) {
	qr, qra := r.buildEntryFilterQueryRestriction(filter)
	qo := r.buildEntrySortQueryClause(sort)

	q := "SELECT " + r.getEntrySelectColumns() + " " +
		"FROM " + r.getEntrySelectTables() + " " +
		qr + " " +
		"AND e.start_time BETWEEN ? AND ? " +
		"GROUP BY " + r.getEntrySelectGroupByColumns() + " " +
		qo

	qa := append(qra, start, end)

	return q, qa
}

// CountDateEntriesByUserId counts all entries (over date) of an user.
func (r *EntryRepo) CountDateEntriesByUserId(ctx context.Context, userId int) (int, error) {
	q, qa := r.buildCountDateEntriesByUserIdQuery(userId)

	sh := newIntScanHelper()
	count, _, qErr := sh.scanRow(r.queryRow(ctx, q, qa...))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries (over date) in database.",
			qErr)
		log.Error(err.StackTrace())
		return 0, err
	}
	return count, nil
}

func (r *EntryRepo) buildCountDateEntriesByUserIdQuery(userId int) (string, []any) {
	q := "SELECT COUNT(DISTINCT(DATE(e.start_time))) " +
		"FROM entry e " +
		"WHERE e.user_id = ?"

	qa := []any{userId}

	return q, qa
}

// GetDateEntriesByUserId retrieves all entries (over date) of an user.
func (r *EntryRepo) GetDateEntriesByUserId(ctx context.Context, userId int, offset int, limit int) (
	[]*model.Entry, error) {
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

	sh := newEntryScanHelper()
	entries, qErr := sh.scanRows(r.query(ctx, q, qa...))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query entries (over date) from database.",
			qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	return entries, nil
}

func (r *EntryRepo) buildGetDateEntriesByUserIdRangeQuery(userId int, offset int, limit int) (string,
	[]any) {
	q := "SELECT DISTINCT(DATE(e.start_time)) AS date " +
		"FROM entry e " +
		"WHERE e.user_id = ? " +
		"ORDER BY date DESC " +
		createQueryLimitString(offset, limit)

	qa := []any{userId}

	return q, qa
}

func (r *EntryRepo) buildGetDateEntriesByUserIdQuery(userId int, start string, end string) (string,
	[]any) {
	q := "SELECT " + r.getEntrySelectColumns() + " " +
		"FROM " + r.getEntrySelectTables() + " " +
		"WHERE e.user_id = ? " +
		"AND e.start_time BETWEEN ? AND ? " +
		"GROUP BY " + r.getEntrySelectGroupByColumns() + " " +
		"ORDER BY e.start_time DESC, e.end_time DESC"

	qa := []any{userId, start, end}

	return q, qa
}

// GetMonthEntries retrieves all entries of a month.
func (r *EntryRepo) GetMonthEntries(ctx context.Context, userId int, year int, month int) (
	[]*model.Entry, error) {
	q, qa := r.buildGetMonthEntriesQuery(userId, year, month)

	sh := newEntryScanHelper()
	entries, qErr := sh.scanRows(r.query(ctx, q, qa...))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query month entries from database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	return entries, nil
}

func (r *EntryRepo) buildGetMonthEntriesQuery(userId int, year int, month int) (string, []any) {
	q := "SELECT " + r.getEntrySelectColumns() + " " +
		"FROM " + r.getEntrySelectTables() + " " +
		"WHERE e.user_id = ? " +
		"AND YEAR(e.start_time) = ? AND MONTH(e.start_time) = ? " +
		"GROUP BY " + r.getEntrySelectGroupByColumns() + " " +
		"ORDER BY e.start_time ASC, e.end_time ASC"

	qa := []any{userId, year, month}

	return q, qa
}

// CountEntries counts all entries.
func (r *EntryRepo) CountEntries(ctx context.Context, filter *model.FieldEntryFilter) (int, error) {
	qr, qra := r.buildEntryFilterQueryRestriction(filter)

	q := "SELECT COUNT(*) "+
		"FROM entry e " + 
		"LEFT JOIN project p ON e.project_id = p.id " +
		qr

	sh := newIntScanHelper()
	count, _, qErr := sh.scanRow(r.queryRow(ctx, q, qra...))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries in database.", qErr)
		log.Error(err.StackTrace())
		return 0, err
	}
	return count, nil
}

// GetEntries retrieves all entries.
func (r *EntryRepo) GetEntries(ctx context.Context, filter *model.FieldEntryFilter,
	sort *model.EntrySort, offset int, limit int) ([]*model.Entry, error) {
	qr, qra := r.buildEntryFilterQueryRestriction(filter)
	qo := r.buildEntrySortQueryClause(sort)

	q := "SELECT " + r.getEntrySelectColumns() + " " +
		"FROM " + r.getEntrySelectTables() + " " +
		qr + " " +
		"GROUP BY " + r.getEntrySelectGroupByColumns() + " " +
		qo + " " +
		createQueryLimitString(offset, limit)

	sh := newEntryScanHelper()
	entries, qErr := sh.scanRows(r.query(ctx, q, qra...))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query entries from database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	return entries, nil
}

// GetEntryById retrieves an entry.
func (r *EntryRepo) GetEntryById(ctx context.Context, id int) (*model.Entry, error) {
	q := "SELECT " + r.getEntrySelectColumns() + " " +
		"FROM " + r.getEntrySelectTables() + " " +
		"WHERE e.id = ? " +
		"GROUP BY " + r.getEntrySelectGroupByColumns()

	sh := newEntryScanHelper()
	entry, found, qErr := sh.scanRow(r.queryRow(ctx, q, id))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read entry %d from database.",
			id), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return entry, nil
}

// GetEntryByIdAndUserId retrieves an entry of an user.
func (r *EntryRepo) GetEntryByIdAndUserId(ctx context.Context, id int, userId int) (*model.Entry,
	error) {
	q := "SELECT " + r.getEntrySelectColumns() + " " +
		"FROM " + r.getEntrySelectTables() + " " +
		"WHERE e.id = ? AND e.user_id = ? " +
		"GROUP BY " + r.getEntrySelectGroupByColumns()

	sh := newEntryScanHelper()
	entry, found, qErr := sh.scanRow(r.queryRow(ctx, q, id, userId))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read entry %d from database.",
			id), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return entry, nil
}

func (r *EntryRepo) getEntrySelectColumns() string {
	return r.getEntrySelectBaseColumns() + ", " +
		r.getEntrySelectProjectColumn() + ", " +
		r.getEntrySelectLabelsColumn()
}

func (r *EntryRepo) getEntrySelectTables() string {
	return "entry e " +
		"LEFT JOIN project p ON e.project_id = p.id " +
		"LEFT JOIN entry_label el ON e.id = el.entry_id " +
		"LEFT JOIN label l ON el.label_id = l.id"
}

func (r *EntryRepo) getEntrySelectGroupByColumns() string {
	return r.getEntrySelectBaseColumns() + ", " +
		r.getEntrySelectProjectColumn()
}

func (r *EntryRepo) getEntrySelectBaseColumns() string {
	return "e.id, e.user_id, e.type_id, e.start_time, e.end_time, e.activity_id, e.description"
}

func (r *EntryRepo) getEntrySelectProjectColumn() string {
	return "p.name"
}

func (r *EntryRepo) getEntrySelectLabelsColumn() string {
	return "GROUP_CONCAT(l.name ORDER BY l.name SEPARATOR ',') AS labels"
}

// ExistsEntryById checks if a entry exists.
func (r *EntryRepo) ExistsEntryById(ctx context.Context, id int) (bool, error) {
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
	error) {
	cnt, cErr := r.count(ctx, "entry", "id = ? AND user_id = ?", id, userId)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries from database.", cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// ExistsEntryByActivityId checks if a entry exists for an activity.
func (r *EntryRepo) ExistsEntryByActivityId(ctx context.Context, activityId int) (bool, error) {
	cnt, cErr := r.count(ctx, "entry", "activity_id = ?", activityId)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not count entries from database.", cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// CreateEntry creates a new entry.
func (r *EntryRepo) CreateEntry(ctx context.Context, entry *model.Entry) error {
	return r.executeInTransaction(ctx, func(tx *sql.Tx) error {
		projectId, cpErr := r.getOrCreateProject(tx, entry.Project)
		if cpErr != nil {
			return cpErr
		}

		etr := toDbEntry(0, entry.UserId, entry.TypeId, entry.StartTime, entry.EndTime,
			entry.ActivityId, projectId, entry.Description)

		q := "INSERT INTO entry (user_id, type_id, start_time, end_time, activity_id, project_id, " +
			"description) VALUES (?, ?, ?, ?, ?, ?, ?)"

		id, cErr := r.insertWithTx(tx, q, etr.userId, etr.typeId, etr.startTime, etr.endTime,
			etr.activityId, etr.projectId, etr.description)
		if cErr != nil {
			err := e.WrapError(e.SysDbInsertFailed, "Could not create entry in database.", cErr)
			log.Error(err.StackTrace())
			return err
		}

		entry.Id = id

		return r.setEntryLabels(tx, entry.Id, entry.Labels)
	})
}

// UpdateEntry updates a entry.
func (r *EntryRepo) UpdateEntry(ctx context.Context, entry *model.Entry) error {
	return r.executeInTransaction(ctx, func(tx *sql.Tx) error {
		projectId, cpErr := r.getOrCreateProject(tx, entry.Project)
		if cpErr != nil {
			return cpErr
		}

		etr := toDbEntry(entry.Id, entry.UserId, entry.TypeId, entry.StartTime, entry.EndTime,
			entry.ActivityId, projectId, entry.Description)

		q := "UPDATE entry SET user_id = ?, type_id = ?, start_time = ?, end_time = ?, " +
			"activity_id = ?, project_id = ?, description = ? WHERE id = ?"

		uErr := r.execWithTx(tx, q, etr.userId, etr.typeId, etr.startTime, etr.endTime,
			etr.activityId, etr.projectId, etr.description, etr.id)
		if uErr != nil {
			err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update entry %d in "+
				"database.", entry.Id), uErr)
			log.Error(err.StackTrace())
			return err
		}

		if ulErr := r.setEntryLabels(tx, entry.Id, entry.Labels); ulErr != nil {
			return ulErr
		}

		if dpErr := r.deleteOrphanedProjects(tx); dpErr != nil {
			return dpErr
		}

		return r.deleteOrphanedLabels(tx)
	})
}

// DeleteEntryById deletes a entry.
func (r *EntryRepo) DeleteEntryById(ctx context.Context, id int) error {
	return r.executeInTransaction(ctx, func(tx *sql.Tx) error {
		q := "DELETE FROM entry WHERE id = ?"

		dErr := r.execWithTx(tx, q, id)
		if dErr != nil {
			err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not delete entry %d from "+
				"database.", id), dErr)
			log.Error(err.StackTrace())
			return err
		}

		if dpErr := r.deleteOrphanedProjects(tx); dpErr != nil {
			return dpErr
		}

		return r.deleteOrphanedLabels(tx)
	})
}

func (r *EntryRepo) getOrCreateProject(tx *sql.Tx, name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	id, err := r.getProject(tx, name)
	if err != nil {
		return 0, err
	}
	if id != 0 {
		return id, nil
	}
	return r.createProject(tx, name)
}

func (r *EntryRepo) getProject(tx *sql.Tx, name string) (int, error) {
	q := "SELECT id FROM project WHERE name = ?"

	sh := newIdScanHelper()
	id, found, qErr := sh.scanRow(r.queryRowWithTx(tx, q, name))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not get project from database.", qErr)
		log.Error(err.StackTrace())
		return 0, err
	}
	if !found {
		return 0, nil
	}
	return id, nil
}

func (r *EntryRepo) createProject(tx *sql.Tx, name string) (int, error) {
	q := "INSERT INTO project (name) VALUES (?)"

	id, cErr := r.insertWithTx(tx, q, name)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not create project in database.", cErr)
		log.Error(err.StackTrace())
		return 0, err
	}
	return id, nil
}

func (r *EntryRepo) deleteOrphanedProjects(tx *sql.Tx) error {
	q := "DELETE FROM project " +
		"WHERE id NOT IN (SELECT DISTINCT project_id FROM entry WHERE project_id IS NOT NULL)"

	if dErr := r.execWithTx(tx, q); dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, "Could not delete orphaned projects from database.",
			dErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}

func (r *EntryRepo) setEntryLabels(tx *sql.Tx, entryId int, labels []string) error {
	q := "DELETE FROM entry_label WHERE entry_id = ?"
	if dErr := r.execWithTx(tx, q, entryId); dErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, "Could not delete entry labels from database.", dErr)
		log.Error(err.StackTrace())
		return err
	}

	seen := make(map[string]bool)
	for _, labelName := range labels {
		if seen[labelName] {
			continue
		}

		labelId, err := r.getOrCreateLabel(tx, labelName)
		if err != nil {
			return err
		}

		q := "INSERT INTO entry_label (entry_id, label_id) VALUES (?, ?)"
		if cErr := r.execWithTx(tx, q, entryId, labelId); cErr != nil {
			err := e.WrapError(e.SysDbUpdateFailed, "Could not insert entry label into database.",
				cErr)
			log.Error(err.StackTrace())
			return err
		}

		seen[labelName] = true
	}

	return nil
}

func (r *EntryRepo) getOrCreateLabel(tx *sql.Tx, name string) (int, error) {
	id, err := r.getLabel(tx, name)
	if err != nil {
		return 0, err
	}
	if id != 0 {
		return id, nil
	}
	return r.createLabel(tx, name)
}

func (r *EntryRepo) getLabel(tx *sql.Tx, name string) (int, error) {
	q := "SELECT id FROM label WHERE name = ?"

	sh := newIdScanHelper()
	id, found, qErr := sh.scanRow(r.queryRowWithTx(tx, q, name))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not get label from database.", qErr)
		log.Error(err.StackTrace())
		return 0, err
	}
	if !found {
		return 0, nil
	}
	return id, nil
}

func (r *EntryRepo) createLabel(tx *sql.Tx, name string) (int, error) {
	q := "INSERT INTO label (name) VALUES (?)"

	id, cErr := r.insertWithTx(tx, q, name)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not insert label into database.", cErr)
		log.Error(err.StackTrace())
		return 0, err
	}
	return id, nil
}

func (r *EntryRepo) deleteOrphanedLabels(tx *sql.Tx) error {
	q := "DELETE FROM label " +
		"WHERE id NOT IN (SELECT DISTINCT label_id FROM entry_label)"

	if dErr := r.execWithTx(tx, q); dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, "Could not delete orphaned labels from database.",
			dErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}

// --- Entry activity functions ---

// GetEntryActivities retrieves all entry activities.
func (r *EntryRepo) GetEntryActivities(ctx context.Context) ([]*model.EntryActivity, error) {
	q := "SELECT id, description FROM entry_activity"

	sh := newEntryActivityScanHelper()
	activities, qErr := sh.scanRows(r.query(ctx, q))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query entry activities from database.",
			qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	return activities, nil
}

// GetEntryActivityByDescription retrieves a entry activity by its description.
func (r *EntryRepo) GetEntryActivityByDescription(ctx context.Context, description string) (
	*model.EntryActivity, error) {
	q := "SELECT id, description FROM entry_activity WHERE description = ?"

	sh := newEntryActivityScanHelper()
	activity, found, qErr := sh.scanRow(r.queryRow(ctx, q, description))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not query entry activity '%s' "+
			"from database.", description), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return activity, nil
}

// ExistsEntryActivityById checks if a entry activity exists.
func (r *EntryRepo) ExistsEntryActivityById(ctx context.Context, id int) (bool, error) {
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
	entryActivity *model.EntryActivity) error {
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
	entryActivity *model.EntryActivity) error {
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
func (r *EntryRepo) DeleteEntryActivityById(ctx context.Context, id int) error {
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
	error) {
	q := "SELECT type_id, SUM(TIMESTAMPDIFF(MINUTE, start_time , end_time)) " +
		"FROM entry " +
		"WHERE user_id = ? " +
		"AND start_time >= ? AND end_time <= ? " +
		"GROUP BY type_id"

	sh := newWorkDurationScanHelper()
	workDurations, qErr := sh.scanRows(r.query(ctx, q, userId, *formatTimestamp(&start),
		*formatTimestamp(&end)))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query work durations from database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	workSummary := model.NewWorkSummary()
	workSummary.UserId = userId
	workSummary.StartTime = start
	workSummary.EndTime = end
	workSummary.WorkDurations = workDurations

	return workSummary, nil
}

// --- Filter helper functions ---

func (r *EntryRepo) buildEntryFilterQueryRestriction(filter *model.FieldEntryFilter) (string, []any) {
	var qrs []string
	var qas []any
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

	if filter.ByProject {
		if filter.Project == "" {
			qrs = append(qrs, "p.name IS NULL")
		} else {
			qrs = append(qrs, "p.name LIKE ?")
			qas = append(qas, "%"+escapeRestrictionString(filter.Project)+"%")
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

	if filter.ByLabel {
		labels := filter.Labels
		if len(labels) == 0 {
			sq := "SELECT 1 FROM entry_label el WHERE el.entry_id = e.id"
			qrs = append(qrs, "NOT EXISTS ("+sq+")")
		} else {
			sq := "SELECT 1 " +
				"FROM entry_label el " +
				"JOIN label l ON el.label_id = l.id " +
				"WHERE el.entry_id = e.id " +
				"AND l.name IN (" + createPlaceholderString(len(labels)) + ")"
			qrs = append(qrs, "EXISTS ("+sq+")")
			for _, label := range labels {
				qas = append(qas, label)
			}
		}
	}

	qr := ""
	if len(qrs) > 0 {
		qr = "WHERE " + strings.Join(qrs[:], " AND ")
	}

	return qr, qas
}

// --- Sort ---

func (r *EntryRepo) buildEntrySortQueryClause(sort *model.EntrySort) string {
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

func (r *EntryRepo) getDateRange(ctx context.Context, query string, args ...any) (
	string, string, error) {
	rows, err := r.getDbHandle(ctx).Query(query, args...)
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

func newEntryScanHelper() *scanHelper[*model.Entry] {
	return newScanHelper(10, scanEntryFunc)
}

func scanEntryFunc(s scanner) (*model.Entry, error) {
	var dbE dbReadEntry

	err := s.Scan(&dbE.id, &dbE.userId, &dbE.typeId, &dbE.startTime, &dbE.endTime, &dbE.activityId,
		&dbE.description, &dbE.project, &dbE.labels)
	if err != nil {
		return nil, err
	}

	entry := fromDbEntry(&dbE)

	return entry, nil
}

func newEntryActivityScanHelper() *scanHelper[*model.EntryActivity] {
	return newScanHelper(10, scanEntryActivityFunc)
}

func scanEntryActivityFunc(s scanner) (*model.EntryActivity, error) {
	var et model.EntryActivity

	err := s.Scan(&et.Id, &et.Description)
	if err != nil {
		return nil, err
	}

	return &et, nil
}

func newWorkDurationScanHelper() *scanHelper[*model.WorkDuration] {
	return newScanHelper(10, scanWorkDurationFunc)
}

func scanWorkDurationFunc(s scanner) (*model.WorkDuration, error) {
	var dbWd dbWorkDuration

	err := s.Scan(&dbWd.typeId, &dbWd.workDuration)
	if err != nil {
		return nil, err
	}

	workDuration := fromDbWorkDuration(&dbWd)

	return workDuration, nil
}

func toDbEntry(id int, userId int, typeId int, startTime time.Time, endTime time.Time,
	activityId int, projectId int, description string) *dbWriteEntry {
	var out dbWriteEntry
	out.id = id
	out.userId = userId
	out.typeId = typeId
	out.startTime = *formatTimestamp(&startTime)
	out.endTime = *formatTimestamp(&endTime)
	if activityId != 0 {
		out.activityId = sql.NullInt64{Int64: int64(activityId), Valid: true}
	} else {
		out.activityId = sql.NullInt64{Int64: 0, Valid: false}
	}
	if projectId != 0 {
		out.projectId = sql.NullInt64{Int64: int64(projectId), Valid: true}
	} else {
		out.projectId = sql.NullInt64{Int64: 0, Valid: false}
	}
	if description != "" {
		out.description = sql.NullString{String: description, Valid: true}
	} else {
		out.description = sql.NullString{String: "", Valid: false}
	}
	return &out
}

func fromDbEntry(in *dbReadEntry) *model.Entry {
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
	if in.project.Valid {
		out.Project = in.project.String
	} else {
		out.Project = ""
	}
	if in.description.Valid {
		out.Description = in.description.String
	} else {
		out.Description = ""
	}
	if in.labels.Valid && in.labels.String != "" {
		out.Labels = strings.Split(in.labels.String, ",")
	} else {
		out.Labels = []string{}
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
