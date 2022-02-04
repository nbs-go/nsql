package query

type ColumnMode = int8

const (
	NonAmbiguousColumn ColumnMode = iota
	JoinColumn
)
