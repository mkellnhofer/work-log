package model

const (
	LoginStepEnterCredentials = iota // 0
	LoginStepChangePassword   = iota // 1
)

const PageNavItems = 5

type EntriesFilterDetails struct {
	Type     string
	Date     string
	Activity string
	Text     string
}
