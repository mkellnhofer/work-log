package model

const (
	LoginStepEnterCredentials = iota // 0
	LoginStepChangePassword   = iota // 1
)

const PageNavItems = 5

type EntryFilterDetails interface {
	isEntryFilterDetails()
}

type baseEntryFilterDetails struct {}

func (*baseEntryFilterDetails) isEntryFilterDetails() {}

type BasicEntryFilterDetails struct {
	baseEntryFilterDetails
	Text string
}

type AdvancedEntryFilterDetails struct {
	baseEntryFilterDetails
	ByType        bool
	Type          string
	ByDate        bool
	Date          string
	ByActivity    bool
	Activity      string
	ByProject     bool
	Project       string
	ByDescription bool
	Description   string
	ByLabels      bool
	Labels        []string
}
