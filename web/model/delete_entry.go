package model

// DeleteEntry stores data for the delete entry view.
type DeleteEntry struct {
	PreviousUrl  string
	ErrorMessage string
	Id           int
}

// NewDeleteEntry creates a new DeleteEntry view model.
func NewDeleteEntry() *DeleteEntry {
	return &DeleteEntry{}
}
