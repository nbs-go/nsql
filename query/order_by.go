package query

type SortDirection = uint8

const (
	Ascending SortDirection = iota
	Descending
)

type OrderBy struct {
	TableName string
	Column    string
	Direction SortDirection
}
