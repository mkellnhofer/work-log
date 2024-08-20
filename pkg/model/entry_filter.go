package model

import "time"

// EntryFilter is an interface for entry filters.
type EntryFilter interface {
	isEntryFilter()
	IsByUser() bool
	GetUserId() int
	SetUserFilter(userId int)
}

// baseEntryFilter stores common filter parameters shared by all entry filters.
type baseEntryFilter struct {
	ByUser bool // Flag to filter by user
	UserId int  // ID of the user
}

func (*baseEntryFilter) isEntryFilter() {}

func (b *baseEntryFilter) IsByUser() bool {
	return b.ByUser
}

func (b *baseEntryFilter) GetUserId() int {
	return b.UserId
}

func (b *baseEntryFilter) SetUserFilter(userId int) {
	b.ByUser = true
	b.UserId = userId
}

// EmptyEntryFilter is an empty entry filter.
type EmptyEntryFilter struct {
	baseEntryFilter
}

// NewEmptyEntryFilter create a new EmptyEntryFilter model.
func NewEmptyEntryFilter() EntryFilter {
	return &EmptyEntryFilter{}
}

// FieldEntryFilter stores parameters to filter entries by specific fields.
type FieldEntryFilter struct {
	baseEntryFilter
	ByType        bool      // Flag to filter by entry type
	TypeId        int       // ID of the entry type
	ByTime        bool      // Flag to filter by time
	StartTime     time.Time // Start time
	EndTime       time.Time // End time
	ByActivity    bool      // Flag to filter by entry activity
	ActivityId    int       // ID of the entry activity
	ByProject     bool      // Flag to filter by project name
	Project       string    // Project name
	ByDescription bool      // Flag to filter by description
	Description   string    // Description
	ByLabel       bool      // Flag to filter by label
	Labels        []string  // Label names
}

// NewFieldEntryFilter create a new FieldEntryFilter model.
func NewFieldEntryFilter() *FieldEntryFilter {
	return &FieldEntryFilter{}
}

// TextEntryFilter stores parameters to filter entries by specific texts.
type TextEntryFilter struct {
	baseEntryFilter
	Text string // Text to filter by
}

// NewTextEntryFilter create a new TextEntryFilter model.
func NewTextEntryFilter() *TextEntryFilter {
	return &TextEntryFilter{}
}
