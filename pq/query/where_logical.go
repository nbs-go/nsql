package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
	"strings"
)

type whereLogicalWriter struct {
	op         op.Operator
	conditions []query.WhereWriter
}

func (w *whereLogicalWriter) SetConditions(conditions []query.WhereWriter) {
	w.conditions = conditions
}

func (w *whereLogicalWriter) WhereQuery() string {
	// If no conditions, then return empty string
	if len(w.conditions) == 0 {
		return ""
	}

	var separator string
	switch w.op {
	case op.And:
		separator = " AND "
	case op.Or:
		separator = " OR "
	default:
		return ""
	}

	var conditions []string
	for _, cw := range w.conditions {
		// Create query
		cq := cw.WhereQuery()

		// If where query is empty, then skip
		if cq == "" {
			continue
		}

		// If condition is a logical, then add brackets
		if _, ok := cw.(query.LogicalWhereWriter); ok {
			cq = fmt.Sprintf("(%s)", cq)
		}

		conditions = append(conditions, cq)
	}

	return strings.Join(conditions, separator)
}

func (w *whereLogicalWriter) GetConditions() []query.WhereWriter {
	return w.conditions
}
