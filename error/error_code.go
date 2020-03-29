package error

const (
	// Authentication errors
	AuthUnknown            = -100
	AuthInvalidCredentials = -101

	// Validation erros
	ValUnknown              = -200
	ValPageNumberInvalid    = -201
	ValIdInvalid            = -202
	ValDateInvalid          = -203
	ValStartDateInvalid     = -204
	ValEndDateInvalid       = -205
	ValStartTimeInvalid     = -206
	ValEndTimeInvalid       = -207
	ValBreakDurationInvalid = -208
	ValDescriptionTooLong   = -209
	ValSearchInvalid        = -210
	ValSearchQueryInvalid   = -211

	// Logic errors
	LogicUnknown                        = -300
	LogicEntryNotFound                  = -301
	LogicEntryTypeNotFound              = -302
	LogicEntryActivityNotFound          = -303
	LogicEntryTimeIntervalInvalid       = -304
	LogicEntryBreakDurationTooLong      = -305
	LogicEntrySearchDateIntervalInvalid = -306

	// System errors
	SysUnknown             = -400
	SysDbUnknown           = -401
	SysDbConnectionFailed  = -402
	SysDbTransactionFailed = -403
	SysDbQueryFailed       = -404
	SysDbInsertFailed      = -405
	SysDbUpdateFailed      = -406
	SysDbDeleteFailed      = -407
	SysJobFailed           = -408
)
