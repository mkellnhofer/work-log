package model

const (
	LoginStepEnterCredentials = iota // 0
	LoginStepChangePassword   = iota // 1
)

const PageNavItems = 5

type EntryFilterDetails struct {
	ByType     bool
	Type       string
	ByDate     bool
	Date       string
	ByActivity bool
	Activity   string
	ByLabels   bool
	Labels     []string
	ByText     bool
	Text       string
}
