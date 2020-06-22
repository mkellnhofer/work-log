package model

// EntryList holds a list of entries.
type EntryList struct {
	Offset int64    `json:"offset"`
	Limit  int      `json:"limit"`
	Total  int64    `json:"total"`
	Items  []*Entry `json:"items"`
}

// NewEntryList creates a new duty list.
func NewEntryList(offset int64, limit int, total int64, items []*Entry) *EntryList {
	return &EntryList{offset, limit, total, items}
}
