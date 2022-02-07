package query

type VariableFormat uint8

const (
	BindVar VariableFormat = iota
	NamedVar
)
