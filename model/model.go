package model

// Sorting defines a sorting direction.
type Sorting int

// Sorting direction constants.
const (
	DescSorting Sorting = -1
	NoSorting   Sorting = 0
	AscSorting  Sorting = 1
)
