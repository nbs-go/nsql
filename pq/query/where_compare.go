package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
)

type whereCompareWriter struct {
	query.ColumnWriter
	op       op.Operator
	variable query.VariableWriter
}

func (w *whereCompareWriter) GetVariable() query.VariableWriter {
	return w.variable
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
