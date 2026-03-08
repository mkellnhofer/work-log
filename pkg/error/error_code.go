package error

const (
	// Authentication errors
	AuthUnknown            = -100
	AuthDataInvalid        = -101
	AuthCredentialsInvalid = -102
	AuthUserNotActivated   = -103
	AuthTokenInvalid       = -104
	AuthTokenNotAllowed    = -105

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

	// General validation erros
	ValUnknown              = -300
	ValJsonInvalid          = -301
	ValPageNumberInvalid    = -302
	ValIdInvalid            = -303
	ValFilterInvalid        = -304
	ValSortInvalid          = -305
	ValOffsetInvalid        = -306
	ValLimitInvalid         = -307
	ValFieldNil             = -308
	ValNumberNegative       = -309
	ValNumberNegativeOrZero = -310
	ValStringEmpty          = -311
	ValStringTooLong        = -312
	ValDateInvalid          = -313
	ValTimestampInvalid     = -314
	ValArrayEmpty           = -315
	ValRoleInvalid          = -316
	ValUsernameInvalid      = -317
	ValPasswordInvalid      = -318
	ValLabelInvalid         = -319
	// View validation errors
	ValStartDateInvalid     = -350
	ValEndDateInvalid       = -351
	ValStartTimeInvalid     = -352
	ValEndTimeInvalid       = -353
	ValProjectNameTooLong   = -354
	ValDescriptionTooLong   = -355
	ValLabelTooShort        = -356
	ValLabelTooLong         = -357
	ValQueryInvalid         = -358
	ValQueryEmpty           = -359
	ValMonthInvalid         = -360
	ValPasswordEmpty        = -361
	ValPasswordTooShort     = -362
	ValPasswordTooLong      = -363
	ValPasswordsNotMatching = -364

	// Logic errors
	LogicUnknown                       = -400
	LogicEntryNotFound                 = -401
	LogicEntryTypeNotFound             = -402
	LogicEntryActivityNotFound         = -403
	LogicEntryActivityDeleteNotAllowed = -404
	LogicEntryTimeIntervalInvalid      = -405
	LogicEntryDateIntervalInvalid      = -406
	LogicRoleNotFound                  = -407
	LogicUserNotFound                  = -408
	LogicUserAlreadyExists             = -409
	LogicContractWorkingHoursInvalid   = -410
	LogicContractVacationDaysInvalid   = -411
	LogicEntryActivityNotAllowed       = -412
	LogicTokenNotFound                 = -413

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
