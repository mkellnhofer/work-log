package tx

import (
	"context"
	"database/sql"
	"net/http"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
)

// --- Transaction holder ---

// TransactionHolder holds the current database transaction.
type TransactionHolder struct {
	tx *sql.Tx
}

// Set sets the current transaction.
func (h *TransactionHolder) Set(tx *sql.Tx) {
	h.tx = tx
}

// Get returns the current transaction.
func (h *TransactionHolder) Get() *sql.Tx {
	return h.tx
}

// Clear clears the current transaction.
func (h *TransactionHolder) Clear() {
	h.tx = nil
}

// --- Transaction middleware ---

// TransactionMiddleware initializes the transaction holder.
type TransactionMiddleware struct {
}

// NewTransactionMiddleware create a new transaction middleware.
func NewTransactionMiddleware() *TransactionMiddleware {
	return &TransactionMiddleware{}
}

// ServeHTTP processes requests.
func (m *TransactionMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	// Create session holder
	txHolder := &TransactionHolder{}

	// Update context
	ctx := r.Context()
	ctx = context.WithValue(ctx, constant.ContextKeyTransactionHolder, txHolder)

	// Forward to next handler
	next(w, r.WithContext(ctx))
}

// --- Transaction manager ---

// TransactionManager provides methods to begin, commit and rollback database transactions.
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager creates a new transaction manager.
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db}
}

// Begin starts a new database transaction.
func (tm *TransactionManager) Begin(ctx context.Context) *e.Error {
	// Get transaction holder
	th := ctx.Value(constant.ContextKeyTransactionHolder).(*TransactionHolder)

	// Check if a transaction already exists
	tx := th.Get()
	if tx != nil {
		err := e.NewError(e.SysDbTransactionFailed, "There is already a database transaction.")
		log.Error(err.StackTrace())
		return err
	}

	// Begin transaction
	tx, bErr := tm.db.Begin()
	if bErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not begin database transaction.", bErr)
		log.Error(err.StackTrace())
		return err
	}

	// Update transaction holder
	th.Set(tx)

	return nil
}

// Commit commits the current database transaction.
func (tm *TransactionManager) Commit(ctx context.Context) *e.Error {
	// Get transaction holder
	th := ctx.Value(constant.ContextKeyTransactionHolder).(*TransactionHolder)

	// Check if there is no transaction
	tx := th.Get()
	if tx == nil {
		err := e.NewError(e.SysDbTransactionFailed, "There is no database transaction.")
		log.Error(err.StackTrace())
		return err
	}

	// Commit transaction
	cErr := tx.Commit()
	if cErr != nil {
		err := e.WrapError(e.SysDbTransactionFailed, "Could not commit database transaction.", cErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// Rollback rolls back the current database transaction.
func (tm *TransactionManager) Rollback(ctx context.Context) *e.Error {
	// Get transaction holder
	th := ctx.Value(constant.ContextKeyTransactionHolder).(*TransactionHolder)

	// Check if there is no transaction
	tx := th.Get()
	if tx == nil {
		err := e.NewError(e.SysDbTransactionFailed, "There is no database transaction.")
		log.Error(err.StackTrace())
		return err
	}

	// Rollback transaction
	rErr := tx.Rollback()
	if rErr != nil {
		err := e.WrapError(e.SysDbTransactionFailed, "Could not rollback database transaction.", rErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}
