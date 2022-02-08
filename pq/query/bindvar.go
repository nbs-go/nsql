package query

import (
	"fmt"
	opt "github.com/nbs-go/nsql/query/option"
)

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

func IntVar(i int) opt.SetOptionFn {
	return func(o *opt.Options) {
		o.KV[opt.VariableKey] = &intVar{value: i}
	}
}

type intVar struct {
	value int
}

func (v *intVar) VariableQuery() string {
	return fmt.Sprintf("%d", v.value)
}
