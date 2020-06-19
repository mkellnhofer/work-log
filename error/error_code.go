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
	ValJsonInvalid          = -301
	ValPageNumberInvalid    = -302
	ValIdInvalid            = -303
	ValDateInvalid          = -304
	ValStartDateInvalid     = -305
	ValEndDateInvalid       = -306
	ValStartTimeInvalid     = -307
	ValEndTimeInvalid       = -308
	ValBreakDurationInvalid = -309
	ValDescriptionTooLong   = -310
	ValSearchInvalid        = -311
	ValSearchQueryInvalid   = -312
	ValMonthInvalid         = -313

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
