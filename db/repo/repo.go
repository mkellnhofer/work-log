package repo

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
)

const defaultPageSize = 100

type repo struct {
	db *sql.DB
}

type database interface {
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
	return r.db.Begin()
}

func (r *repo) commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (r *repo) rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (r *repo) count(table string, restriction string, args ...interface{}) (int, error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE "+restriction, args...)
	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *repo) query(sh scanHelper, query string, args ...interface{}) (interface{}, error) {
	return queryInternal(r.db, sh, query, args...)
}

func (r *repo) queryWithTx(tx *sql.Tx, sh scanHelper, query string, args ...interface{}) (interface{},
	error) {
	return queryInternal(tx, sh, query, args...)
}

func queryInternal(db database, sh scanHelper, query string, args ...interface{}) (interface{},
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

func (r *repo) queryRow(sh scanHelper, query string, args ...interface{}) (interface{}, error) {
	return queryRowInternal(r.db, sh, query, args...)
}

func (r *repo) queryRowWithTx(tx *sql.Tx, sh scanHelper, query string, args ...interface{}) (
	interface{}, error) {
	return queryRowInternal(tx, sh, query, args...)
}

func queryRowInternal(db database, sh scanHelper, query string, args ...interface{}) (interface{},
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

func (r *repo) queryValue(value interface{}, query string, args ...interface{}) error {
	return queryValueInternal(r.db, value, query, args...)
}

func (r *repo) queryValueWithTx(tx *sql.Tx, value interface{}, query string, args ...interface{}) error {
	return queryValueInternal(tx, value, query, args...)
}

func queryValueInternal(db database, value interface{}, query string, args ...interface{}) error {
	return db.QueryRow(query, args...).Scan(value)
}

func (r *repo) insert(query string, args ...interface{}) (int, error) {
	return insertInternal(r.db, query, args...)
}

func (r *repo) insertWithTx(tx *sql.Tx, query string, args ...interface{}) (int, error) {
	return insertInternal(tx, query, args...)
}

func insertInternal(db database, query string, args ...interface{}) (int, error) {
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

func (r *repo) exec(query string, args ...interface{}) error {
	return execInternal(r.db, query, args...)
}

func (r *repo) execWithTx(tx *sql.Tx, query string, args ...interface{}) error {
	return execInternal(tx, query, args...)
}

func execInternal(db database, query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	return err
}

// --- Helper functions ---

func parseDate(ts *string) *time.Time {
	if ts == nil {
		return nil
	}

	t, pErr := time.Parse(constant.DbDateFormat, *ts)
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

	ts := t.Format(constant.DbDateFormat)
	return &ts
}

func parseTimestamp(ts *string) *time.Time {
	if ts == nil {
		return nil
	}

	t, pErr := time.Parse(constant.DbTimestampFormat, *ts)
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

	ts := t.Format(constant.DbTimestampFormat)
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

func escapeRestrictionString(s string) string {
	return strings.NewReplacer("%", "\\%", "_", "\\_").Replace(s)
}
