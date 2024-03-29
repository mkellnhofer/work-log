package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"kellnhofer.com/work-log/pkg/constant"
	"kellnhofer.com/work-log/pkg/db/tx"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
)

const defaultPageSize = 100

type repo struct {
	db *sql.DB
}

type dbHandle interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// --- Limit helper functions ---

func createQueryLimitString(offset int, limit int) string {
	if offset <= 0 && limit <= 0 {
		return ""
	}
	off := strconv.Itoa(offset)
	var lim string
	if limit > 0 {
		lim = strconv.Itoa(limit)
	} else {
		lim = strconv.Itoa(defaultPageSize)
	}
	return "LIMIT " + off + ", " + lim
}

// --- Scan helper functions ---

type scanner interface {
	Scan(dest ...interface{}) error
}

type scanHelper interface {
	makeSlice() interface{}
	scan(s scanner) (interface{}, error)
	appendSlice(interface{}, interface{}) interface{}
}

type scanIntHelper struct {
}

func (h *scanIntHelper) makeSlice() interface{} {
	return make([]int, 0, 10)
}

func (h *scanIntHelper) scan(s scanner) (interface{}, error) {
	var id int

	err := s.Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *scanIntHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]int), item.(int))
}

type scanIdHelper struct {
	scanIntHelper
}

// --- DB functions ---

func (r *repo) begin() (*sql.Tx, error) {
	tx, bErr := r.db.Begin()
	if bErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not begin database transaction.", bErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	return tx, nil
}

func (r *repo) commit(tx *sql.Tx) error {
	cErr := tx.Commit()
	if cErr != nil {
		err := e.WrapError(e.SysDbTransactionFailed, "Could not commit database transaction.", cErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}

func (r *repo) rollback(tx *sql.Tx) error {
	rErr := tx.Rollback()
	if rErr != nil {
		err := e.WrapError(e.SysDbTransactionFailed, "Could not rollback database transaction.", rErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}

func (r *repo) getCurrentTransaction(ctx context.Context) *sql.Tx {
	th := ctx.Value(constant.ContextKeyTransactionHolder).(*tx.TransactionHolder)
	return th.Get()
}

func (r *repo) executeInTransaction(ctx context.Context, txf func(tx *sql.Tx) error) error {
	tx := r.getCurrentTransaction(ctx)
	isExistingTx := tx != nil

	if !isExistingTx {
		var err error
		if tx, err = r.begin(); err != nil {
			return err
		}
	}

	if err := txf(tx); err != nil {
		if !isExistingTx {
			r.rollback(tx)
		}
		return err
	}

	if !isExistingTx {
		if err := r.commit(tx); err != nil {
			return err
		}
	}

	return nil
}

func (r *repo) getDbHandle(ctx context.Context) dbHandle {
	tx := r.getCurrentTransaction(ctx)
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *repo) count(ctx context.Context, table string, restriction string, args ...interface{}) (
	int, error) {
	return countInternal(r.getDbHandle(ctx), table, restriction, args...)
}

func (r *repo) countWithTx(tx *sql.Tx, table string, restriction string, args ...interface{}) (int,
	error) {
	return countInternal(tx, table, restriction, args...)
}

func countInternal(db dbHandle, table string, restriction string, args ...interface{}) (int, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE "+restriction, args...)
	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *repo) query(ctx context.Context, sh scanHelper, query string, args ...interface{}) (
	interface{}, error) {
	return queryInternal(r.getDbHandle(ctx), sh, query, args...)
}

func (r *repo) queryWithTx(tx *sql.Tx, sh scanHelper, query string, args ...interface{}) (interface{},
	error) {
	return queryInternal(tx, sh, query, args...)
}

func queryInternal(db dbHandle, sh scanHelper, query string, args ...interface{}) (interface{},
	error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sr, err := scanRows(rows, sh)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

func scanRows(rows *sql.Rows, sh scanHelper) (interface{}, error) {
	objs := sh.makeSlice()
	for rows.Next() {
		obj, err := sh.scan(rows)
		if err != nil {
			return nil, err
		}
		objs = sh.appendSlice(objs, obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objs, nil
}

func (r *repo) queryRow(ctx context.Context, sh scanHelper, query string, args ...interface{}) (
	interface{}, error) {
	return queryRowInternal(r.getDbHandle(ctx), sh, query, args...)
}

func (r *repo) queryRowWithTx(tx *sql.Tx, sh scanHelper, query string, args ...interface{}) (
	interface{}, error) {
	return queryRowInternal(tx, sh, query, args...)
}

func queryRowInternal(db dbHandle, sh scanHelper, query string, args ...interface{}) (interface{},
	error) {
	row := db.QueryRow(query, args...)

	sr, err := sh.scan(row)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
	}

	return sr, nil
}

func (r *repo) queryValue(ctx context.Context, value interface{}, query string,
	args ...interface{}) error {
	return queryValueInternal(r.getDbHandle(ctx), value, query, args...)
}

func (r *repo) queryValueWithTx(tx *sql.Tx, value interface{}, query string,
	args ...interface{}) error {
	return queryValueInternal(tx, value, query, args...)
}

func queryValueInternal(db dbHandle, value interface{}, query string, args ...interface{}) error {
	return db.QueryRow(query, args...).Scan(value)
}

func (r *repo) insert(ctx context.Context, query string, args ...interface{}) (int, error) {
	return insertInternal(r.getDbHandle(ctx), query, args...)
}

func (r *repo) insertWithTx(tx *sql.Tx, query string, args ...interface{}) (int, error) {
	return insertInternal(tx, query, args...)
}

func insertInternal(db dbHandle, query string, args ...interface{}) (int, error) {
	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repo) exec(ctx context.Context, query string, args ...interface{}) error {
	return execInternal(r.getDbHandle(ctx), query, args...)
}

func (r *repo) execWithTx(tx *sql.Tx, query string, args ...interface{}) error {
	return execInternal(tx, query, args...)
}

func execInternal(db dbHandle, query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	return err
}

// --- Helper functions ---

func parseDate(ts *string) *time.Time {
	if ts == nil {
		return nil
	}

	t, pErr := time.ParseInLocation(constant.DbDateFormat, *ts, time.Local)
	if pErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not parse date.", pErr)
		log.Error(err.StackTrace())
		panic(err)
	}
	return &t
}

func formatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}

	tl := t.Local()
	ts := tl.Format(constant.DbDateFormat)
	return &ts
}

func parseTimestamp(ts *string) *time.Time {
	if ts == nil {
		return nil
	}

	t, pErr := time.ParseInLocation(constant.DbTimestampFormat, *ts, time.Local)
	if pErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not parse timestamp.", pErr)
		log.Error(err.StackTrace())
		panic(err)
	}
	return &t
}

func formatTimestamp(t *time.Time) *string {
	if t == nil {
		return nil
	}

	tl := t.Local()
	ts := tl.Format(constant.DbTimestampFormat)
	return &ts
}

func parseDuration(min *int) *time.Duration {
	if min == nil {
		return nil
	}

	d, pErr := time.ParseDuration(fmt.Sprintf("%dm", *min))
	if pErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not parse duration.", pErr)
		log.Error(err.StackTrace())
		panic(err)
	}
	return &d
}

func formatDuration(d *time.Duration) *int {
	if d == nil {
		return nil
	}

	md := d.Round(time.Minute)
	min := int(md.Minutes())
	return &min
}

func createSelectionString(values []string) string {
	s := quoteStrings(values)
	return strings.Join(s, ",")
}

func escapeRestrictionString(restriction string) string {
	return strings.NewReplacer("%", "\\%", "_", "\\_").Replace(restriction)
}

func quoteStrings(in []string) []string {
	out := make([]string, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = "'" + in[i] + "'"
	}
	return out
}
