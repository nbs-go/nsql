package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
	"strings"
)

type whereLogicWriter struct {
	op         op.Operator
	conditions []query.WhereWriter
}

func (w *whereLogicWriter) SetConditions(conditions []query.WhereWriter) {
	w.conditions = conditions
}

func (w *whereLogicWriter) WhereQuery() string {
	var separator string
	if w.op == op.Or {
		separator = " OR "
	} else {
		separator = " AND "
	}

	var conditions []string
	for _, cw := range w.conditions {
		// Create query
		cq := cw.WhereQuery()

		// If condition is a logical, then add brackets
		if _, ok := cw.(query.WhereLogicWriter); ok {
			cq = fmt.Sprintf("(%s)", cq)
		}

		conditions = append(conditions, cq)
	}

	return strings.Join(conditions, separator)
}

func (w *whereLogicWriter) GetConditions() []query.WhereWriter {
	return w.conditions
}
