package query

type bindVar struct{}

func (b *bindVar) VariableQuery() string {
	return "?"
}

type betweenBindVar struct{}

func (b *betweenBindVar) VariableQuery() string {
	return "? AND ?"
}

type inBindVar struct{}

func (b *inBindVar) VariableQuery() string {
	return "(?)"
}

type namedVar struct {
	column string
}

func (v *namedVar) VariableQuery() string {
	return ":" + v.column
}
