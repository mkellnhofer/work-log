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
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
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
	Scan(dest ...any) error
}

type scanFunc[T any] func(s scanner) (T, error)

type scanHelper[T any] struct {
	capacity int
	scanFn   scanFunc[T]
}

func (sh *scanHelper[T]) makeSlice() []T {
	return make([]T, 0, sh.capacity)
}

func (sh *scanHelper[T]) appendSlice(items []T, item T) []T {
	return append(items, item)
}

func (sh *scanHelper[T]) scanRows(rows *sql.Rows, err error) ([]T, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	objs := sh.makeSlice()
	for rows.Next() {
		obj, err := sh.scanFn(rows)
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

func (sh *scanHelper[T]) scanRow(row *sql.Row) (T, bool, error) {
	sr, err := sh.scanFn(row)

	switch {
	case err == sql.ErrNoRows:
		return sr, false, nil
	case err != nil:
		return sr, false, err
	default:
		return sr, true, nil
	}
}

func newScanHelper[T any](capacity int, scanFn scanFunc[T]) *scanHelper[T] {
	return &scanHelper[T]{capacity: capacity, scanFn: scanFn}
}

func newIntScanHelper() *scanHelper[int] {
	return newScanHelper(10, scanIntFunc)
}

func newIdScanHelper() *scanHelper[int] {
	return newScanHelper(10, scanIntFunc)
}

func scanIntFunc(s scanner) (int, error) {
	var val int
	err := s.Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
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

func (r *repo) count(ctx context.Context, table string, restriction string, args ...any) (
	int, error) {
	return countInternal(r.getDbHandle(ctx), table, restriction, args...)
}

func (r *repo) countWithTx(tx *sql.Tx, table string, restriction string, args ...any) (int, error) {
	return countInternal(tx, table, restriction, args...)
}

func countInternal(db dbHandle, table string, restriction string, args ...any) (int, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE "+restriction, args...)
	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *repo) query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return queryInternal(r.getDbHandle(ctx), query, args...)
}

func (r *repo) queryWithTx(tx *sql.Tx, query string, args ...any) (*sql.Rows, error) {
	return queryInternal(tx, query, args...)
}

func queryInternal(db dbHandle, query string, args ...any) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func (r *repo) queryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return queryRowInternal(r.getDbHandle(ctx), query, args...)
}

func (r *repo) queryRowWithTx(tx *sql.Tx, query string, args ...any) *sql.Row {
	return queryRowInternal(tx, query, args...)
}

func queryRowInternal(db dbHandle, query string, args ...any) *sql.Row {
	return db.QueryRow(query, args...)
}

func (r *repo) queryValue(ctx context.Context, value any, query string, args ...any) error {
	return queryValueInternal(r.getDbHandle(ctx), value, query, args...)
}

func (r *repo) queryValueWithTx(tx *sql.Tx, value any, query string, args ...any) error {
	return queryValueInternal(tx, value, query, args...)
}

func queryValueInternal(db dbHandle, value any, query string, args ...any) error {
	return db.QueryRow(query, args...).Scan(value)
}

func (r *repo) insert(ctx context.Context, query string, args ...any) (int, error) {
	return insertInternal(r.getDbHandle(ctx), query, args...)
}

func (r *repo) insertWithTx(tx *sql.Tx, query string, args ...any) (int, error) {
	return insertInternal(tx, query, args...)
}

func insertInternal(db dbHandle, query string, args ...any) (int, error) {
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

func (r *repo) exec(ctx context.Context, query string, args ...any) error {
	return execInternal(r.getDbHandle(ctx), query, args...)
}

func (r *repo) execWithTx(tx *sql.Tx, query string, args ...any) error {
	return execInternal(tx, query, args...)
}

func execInternal(db dbHandle, query string, args ...any) error {
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

func createPlaceholderString(count int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ",")
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
