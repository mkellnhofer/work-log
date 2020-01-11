package error

const (
	// Authentication errors
	AuthUnknown = -100

	// Validation erros
	ValUnknown = -200

	// Logic errors
	LogicUnknown                   = -300
	LogicEntryNotFound             = -301
	LogicEntryTypeNotFound         = -302
	LogicEntryActivityNotFound     = -303
	LogicEntryTimeIntervalInvalid  = -304
	LogicEntryBreakDurationTooLong = -305

	// System errors
	SysUnknown             = -400
	SysDbUnknown           = -401
	SysDbConnectionFailed  = -402
	SysDbTransactionFailed = -403
	SysDbQueryFailed       = -404
	SysDbInsertFailed      = -405
	SysDbUpdateFailed      = -406
	SysDbDeleteFailed      = -407
)
