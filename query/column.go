package query

type ColumnFormat int8

const (
	NonAmbiguousColumn ColumnFormat = iota
	SelectJoinColumn
	NamedArgColumn
)
