package controller

import (
	"fmt"
	"strings"

	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
)

type sort struct {
	name     string
	operator string
}

const (
	sortOpAsc  = "asc"
	sortOpDesc = "desc"
)

// --- Entry sort functions ---

const (
	sortNameStartTime = "startTime"
)

func getEntriesSort(str string) (*model.EntriesSort, *e.Error) {
	entrySort := model.NewEntriesSort()

	// Get start time sorting value
	var err *e.Error
	if entrySort.ByTime, err = createSort(str, sortNameStartTime); err != nil {
		return nil, err
	}

	return entrySort, nil
}

// --- Sorting parsing functions ---

func createSort(str string, name string) (model.Sorting, *e.Error) {
	// If sort string is empty: Abort
	if str == "" {
		return model.NoSorting, nil
	}

	// Extract field name and operator from sort string
	sortParts := strings.Split(str, ";")
	if err := checkSortParts(sortParts); err != nil {
		return model.NoSorting, err
	}
	sort := sort{sortParts[0], sortParts[1]}

	// If unsupported field name: Abort
	if sort.name != name {
		err := e.NewError(e.ValSortInvalid, fmt.Sprintf("Invalid sorting. (Unknown/unsupported "+
			" field name '%s'.)", sort.name))
		log.Debug(err.StackTrace())
		return model.NoSorting, err
	}

	// Get sorting
	switch sort.operator {
	case sortOpAsc:
		return model.AscSorting, nil
	case sortOpDesc:
		return model.DescSorting, nil
	default:
		err := e.NewError(e.ValSortInvalid, fmt.Sprintf("Invalid sorting. (Unsupported operator "+
			"'%s'.)", sort.operator))
		log.Debug(err.StackTrace())
		return model.NoSorting, err
	}
}

func checkSortParts(parts []string) *e.Error {
	// Check if structure is invalid
	if len(parts) < 2 || len(parts) > 2 {
		err := e.NewError(e.ValSortInvalid, "Invalid sorting. (A sorting must have following "+
			"structure: [field name];[operator])")
		log.Debug(err.StackTrace())
		return err
	}

	// Check if field name is invalid
	if parts[0] == "" {
		err := e.NewError(e.ValSortInvalid, "Invalid sorting. (Field name cannot be empty.)")
		log.Debug(err.StackTrace())
		return err
	}

	// Check if operator is invalid
	if parts[1] == "" {
		err := e.NewError(e.ValSortInvalid, "Invalid sorting. (Operator cannot be empty.)")
		log.Debug(err.StackTrace())
		return err
	}

	return nil
}
