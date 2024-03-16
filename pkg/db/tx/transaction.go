package tx

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/constant"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
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

// CreateHandler creates a new handler to process requests.
func (m *TransactionMiddleware) CreateHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return m.process(next, c)
	}
}

func (m *TransactionMiddleware) process(next echo.HandlerFunc, c echo.Context) error {
	// Get request
	req := c.Request()

	// Create session holder
	txHolder := &TransactionHolder{}

	// Update context
	ctx := context.WithValue(req.Context(), constant.ContextKeyTransactionHolder, txHolder)
	c.SetRequest(req.WithContext(ctx))

	// Forward to next handler
	return next(c)
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
func (tm *TransactionManager) Begin(ctx context.Context) error {
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
func (tm *TransactionManager) Commit(ctx context.Context) error {
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
func (tm *TransactionManager) Rollback(ctx context.Context) error {
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

// Executes the provided function in a new database transaction.
func (tm *TransactionManager) ExecuteInNewTransaction(ctx context.Context,
	txf func(ctx context.Context) error) error {
	// Start transaction
	if err := tm.Begin(ctx); err != nil {
		return err
	}

	// Execute wrapped function
	if err := txf(ctx); err != nil {
		tm.Rollback(ctx)
		return err
	}

	// Commit transaction
	return tm.Commit(ctx)
}
