package model

// Sorting defines a sorting direction.
type Sorting int

// Sorting direction constants.
const (
	DescSorting Sorting = -1
	NoSorting   Sorting = 0
	AscSorting  Sorting = 1
)

// Maximum string length constants.
const (
	MaxLengthRoleName                 = 100
	MaxLengthUserName                 = 100
	MaxLengthUserUsername             = 100
	MaxLengthUserPassword             = 100
	MaxLengthEntryTypeDescription     = 50
	MaxLengthEntryActivityDescription = 50
	MaxLengthEntryDescription         = 200
)
