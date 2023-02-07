package query

type nullVar struct{}

func (b *nullVar) VariableQuery() string {
	return "NULL"
}
