package model

// EntryList
//
// A list of entries.
//
// swagger:model EntryList
type EntryList struct {
	// The entries page offset.
	// min: 0
	// example: 0
	Offset int `json:"offset"`

	// The entries page limit.
	// min: 0
	// example: 0
	Limit int `json:"limit"`

	// The total count of entries available.
	// min: 0
	// example: 0
	Total int `json:"total"`

	// The entries.
	Items []*Entry `json:"items"`
}

// NewEntryList creates a new entry list.
func NewEntryList(offset int, limit int, total int, items []*Entry) *EntryList {
	return &EntryList{offset, limit, total, items}
}
