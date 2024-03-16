package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
)

type filter struct {
	name     string
	operator string
	values   []string
}

const (
	filterOpIs       = "i"
	filterOpEqual    = "eq"
	filterOpContains = "cn"
	filterOpBetween  = "bt"
)

// --- Entry filter functions ---

const (
	filterNameUserId      = "userId"
	filterNameTypeId      = "typeId"
	filterNameStartTime   = "startTime"
	filterNameActivityId  = "activityId"
	filterNameDescription = "description"
)

func getEntriesFilter(str string) (*model.EntriesFilter, error) {
	// Get filter values
	filters, err := getFilters(str)
	if err != nil {
		return nil, err
	}

	// Check for unknown/unsupported filters
	filterNames := []string{filterNameUserId, filterNameTypeId, filterNameStartTime,
		filterNameActivityId, filterNameDescription}
	if err = checkFiltersSupported(filters, filterNames); err != nil {
		return nil, err
	}

	entryFilter := model.NewEntriesFilter()

	// Get user ID filter
	if filter, ok := filters[filterNameUserId]; ok {
		entryFilter.ByUser = true
		entryFilter.UserId, err = getIdFilterValue(filter, false)
		if err != nil {
			return nil, err
		}
	}
	// Get type ID filter
	if filter, ok := filters[filterNameTypeId]; ok {
		entryFilter.ByType = true
		entryFilter.TypeId, err = getIdFilterValue(filter, false)
		if err != nil {
			return nil, err
		}
	}
	// Get start time filter
	if filter, ok := filters[filterNameStartTime]; ok {
		entryFilter.ByTime = true
		entryFilter.StartTime, entryFilter.EndTime, err = getTimeIntervalFilterValue(filter)
		if err != nil {
			return nil, err
		}
	}
	// Get activity ID filter
	if filter, ok := filters[filterNameActivityId]; ok {
		entryFilter.ByActivity = true
		entryFilter.ActivityId, err = getIdFilterValue(filter, true)
		if err != nil {
			return nil, err
		}
	}
	// Get description filter
	if filter, ok := filters[filterNameDescription]; ok {
		entryFilter.ByDescription = true
		entryFilter.Description, err = getStringFilterValue(filter, true)
		if err != nil {
			return nil, err
		}
	}

	return entryFilter, nil
}

// --- Filter parsing functions ---

func getFilters(str string) (map[string]*filter, error) {
	// If filter string is empty: Abort
	if str == "" {
		return make(map[string]*filter), nil
	}

	// Split fiter string
	filtersParts := strings.Split(str, "|")
	fs := make([]*filter, len(filtersParts))
	for i, filtersPart := range filtersParts {
		filterParts := strings.Split(filtersPart, ";")
		if err := checkFilterParts(filterParts); err != nil {
			return nil, err
		}
		fs[i] = &filter{filterParts[0], filterParts[1], filterParts[2:]}
	}

	// Create filters map
	fm := make(map[string]*filter)
	for _, f := range fs {
		fm[f.name] = f
	}

	return fm, nil
}

func checkFilterParts(parts []string) error {
	// Check if structure is invalid
	if len(parts) < 3 {
		err := e.NewError(e.ValFilterInvalid, "Invalid filter. (A filter must have following "+
			"structure: [field name];[operator];[value-1]...[value-n])")
		log.Debug(err.StackTrace())
		return err
	}

	// Check if field name is invalid
	if parts[0] == "" {
		err := e.NewError(e.ValFilterInvalid, "Invalid filter. (Field name cannot be empty.)")
		log.Debug(err.StackTrace())
		return err
	}

	// Check if operator is invalid
	if parts[1] == "" {
		err := e.NewError(e.ValFilterInvalid, "Invalid filter. (Operator cannot be empty.)")
		log.Debug(err.StackTrace())
		return err
	}

	return nil
}

func checkFiltersSupported(filters map[string]*filter, names []string) error {
	for name := range filters {
		if !isValidFilterName(names, name) {
			err := e.NewError(e.ValFilterInvalid, fmt.Sprintf("Invalid filter. (Unknown/unsupported "+
				" field name '%s'.)", name))
			log.Debug(err.StackTrace())
			return err
		}
	}
	return nil
}

func isValidFilterName(names []string, n string) bool {
	for _, name := range names {
		if name == n {
			return true
		}
	}
	return false
}

func getIdFilterValue(f *filter, nullable bool) (int, error) {
	// Get "is null" value
	if nullable && f.operator == filterOpIs {
		if len(f.values) == 1 && f.values[0] == "null" {
			return 0, nil
		}
		return 0, createInvalidFilterValueError(f.name)
	}

	// Get "equals" value
	if f.operator == filterOpEqual {
		if len(f.values) == 1 && f.values[0] != "" {
			if id, err := strconv.Atoi(f.values[0]); err == nil && id > 0 {
				return id, nil
			}
		}
		return 0, createInvalidFilterValueError(f.name)
	}

	return 0, createInvalidFilterOperatorError(f.name, f.operator)
}

func getStringFilterValue(f *filter, nullable bool) (string, error) {
	// Get "is null" value
	if nullable && f.operator == filterOpIs {
		if len(f.values) == 1 && f.values[0] == "null" {
			return "", nil
		}
		return "", createInvalidFilterValueError(f.name)
	}

	// Get "contains" value
	if f.operator == filterOpContains {
		if len(f.values) == 1 && f.values[0] != "" {
			return f.values[0], nil
		}
		return "", createInvalidFilterValueError(f.name)
	}

	return "", createInvalidFilterOperatorError(f.name, f.operator)
}

func getTimeIntervalFilterValue(f *filter) (time.Time, time.Time, error) {
	// Check if wrong operator
	if f.operator != filterOpBetween {
		return time.Now(), time.Now(), createInvalidFilterOperatorError(f.name, f.operator)
	}
	// Check if wrong number of values
	if len(f.values) != 2 {
		return time.Now(), time.Now(), createInvalidFilterValueError(f.name)
	}

	// Parse values
	start, sErr := parseTimestamp(f.values[0])
	end, eErr := parseTimestamp(f.values[1])
	if sErr != nil || eErr != nil {
		return time.Now(), time.Now(), createInvalidFilterValueError(f.name)
	}

	return start, end, nil
}

func createInvalidFilterOperatorError(name string, operator string) error {
	err := e.NewError(e.ValFilterInvalid, fmt.Sprintf("Invalid filter '%s'. (Unsupported "+
		"operator '%s'.)", name, operator))
	log.Debug(err.StackTrace())
	return err
}

func createInvalidFilterValueError(name string) error {
	err := e.NewError(e.ValFilterInvalid, fmt.Sprintf("Invalid filter '%s'. (Invalid value(s).)",
		name))
	log.Debug(err.StackTrace())
	return err
}
