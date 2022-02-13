package query

import (
	"fmt"
	"github.com/nbs-go/nsql/option"
)

type bindVar struct{}

func (b *bindVar) VariableQuery() string {
	return "?"
}

type betweenBindVar struct{}

func (b *betweenBindVar) VariableQuery() string {
	return "? AND ?"
}

type inBindVar struct {
	argCount int
}

func (b *inBindVar) VariableQuery() string {
	switch b.argCount {
	case 0:
		panic(fmt.Errorf("invalid bindVar for IN query, does not have argument"))
	case 1:
		return "(?)"
	}

	// Write bind var query
	q := "(?"
	for i := 1; i < b.argCount; i++ {
		q += ", ?"
	}
	return q + ")"
}

type namedVar struct {
	column string
}

func (v *namedVar) VariableQuery() string {
	return ":" + v.column
}

func IntVar(i int) option.SetOptionFn {
	return func(o *option.Options) {
		o.KV[option.VariableKey] = &intVar{value: i}
	}
}

type intVar struct {
	value int
}

func (v *intVar) VariableQuery() string {
	return fmt.Sprintf("%d", v.value)
}
