package nsql

const (
	Separator = ", "
)

type ColumnFormat int8

const (
	NonAmbiguousColumn ColumnFormat = iota
	SelectJoinColumn
	ColumnOnly
)

type VariableFormat uint8

const (
	BindVar VariableFormat = iota
	NamedVar
)
