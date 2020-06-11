package error

const (
	// Authentication errors
	AuthUnknown            = -100
	AuthInvalidCredentials = -101

	// Permission errors
	PermUnknown             = -200
	PermGetUserData         = -201
	PermChangeUserData      = -202
	PermGetUserAccount      = -203
	PermChangeUserAccount   = -204
	PermGetEntryCharacts    = -205
	PermChangeEntryCharacts = -206
	PermGetAllEntries       = -207
	PermChangeAllEntries    = -208
	PermGetOwnEntries       = -209
	PermChangeOwnEntries    = -210

	// Validation erros
	ValUnknown              = -300
	ValPageNumberInvalid    = -301
	ValIdInvalid            = -302
	ValDateInvalid          = -303
	ValStartDateInvalid     = -304
	ValEndDateInvalid       = -305
	ValStartTimeInvalid     = -306
	ValEndTimeInvalid       = -307
	ValBreakDurationInvalid = -308
	ValDescriptionTooLong   = -309
	ValSearchInvalid        = -310
	ValSearchQueryInvalid   = -311
	ValMonthInvalid         = -312

	// Logic errors
	LogicUnknown                        = -400
	LogicEntryNotFound                  = -401
	LogicEntryTypeNotFound              = -402
	LogicEntryActivityNotFound          = -403
	LogicEntryActivityDeleteNotAllowed  = -404
	LogicEntryTimeIntervalInvalid       = -405
	LogicEntryBreakDurationTooLong      = -406
	LogicEntrySearchDateIntervalInvalid = -407
	LogicRoleNotFound                   = -408
	LogicUserNotFound                   = -409
	LogicUserAlreadyExists              = -410

	// System errors
	SysUnknown             = -500
	SysDbUnknown           = -501
	SysDbConnectionFailed  = -502
	SysDbTransactionFailed = -503
	SysDbQueryFailed       = -504
	SysDbInsertFailed      = -505
	SysDbUpdateFailed      = -506
	SysDbDeleteFailed      = -507
	SysJobFailed           = -508
)
