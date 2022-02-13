package op

type ColumnFormat int8

const (
	NonAmbiguousColumn ColumnFormat = iota
	SelectJoinColumn
	ColumnOnly
)
