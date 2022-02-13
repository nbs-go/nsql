package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
)

type whereCompareWriter struct {
	nsql.ColumnWriter
	op       op.Operator
	variable nsql.VariableWriter
	as       string
}

func (w *whereCompareWriter) SelectQuery() string {
	q := w.WhereQuery()

	if w.as != "" {
		q += fmt.Sprintf(` AS "%s"`, w.as)
	}

	return q
}

func (w *whereCompareWriter) IsAllColumns() bool {
	return false
}

func (w *whereCompareWriter) GetVariable() nsql.VariableWriter {
	return w.variable
}

func (w *whereCompareWriter) SetVariable(v nsql.VariableWriter) {
	w.variable = v
}

func (w *whereCompareWriter) WhereQuery() string {
	var operator string
	switch w.op {
	case op.Equal:
		operator = "="
	case op.NotEqual:
		operator = "!="
	case op.GreaterThan:
		operator = ">"
	case op.GreaterThanEqual:
		operator = ">="
	case op.LessThan:
		operator = "<"
	case op.LessThanEqual:
		operator = "<="
	case op.Like:
		operator = "LIKE"
	case op.NotLike:
		operator = "NOT LIKE"
	case op.ILike:
		operator = "ILIKE"
	case op.NotILike:
		operator = "NOT ILIKE"
	case op.Between:
		operator = "BETWEEN"
	case op.NotBetween:
		operator = "NOT BETWEEN"
	case op.In:
		operator = "IN"
	case op.NotIn:
		operator = "NOT IN"
	default:
		return ""
	}

	q := fmt.Sprintf(`%s %s %s`, w.ColumnQuery(), operator, w.variable.VariableQuery())

	return q
}
