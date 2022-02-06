package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
)

type baseWhereComparisonWriter struct {
	query.ColumnWriter
	op op.Operator
}

type whereComparisonWriter struct {
	baseWhereComparisonWriter
}

func (w *whereComparisonWriter) WhereQuery() string {
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
	default:
		return ""
	}

	q := fmt.Sprintf(`%s %s ?`, w.ColumnQuery(), operator)

	return q
}

type whereBetweenWriter struct {
	baseWhereComparisonWriter
}

func (w *whereBetweenWriter) WhereQuery() string {
	var operator string
	switch w.op {
	case op.Between:
		operator = "BETWEEN"
	case op.NotBetween:
		operator = "NOT BETWEEN"
	default:
		return ""
	}
	return fmt.Sprintf(`%s %s ? AND ?`, w.ColumnQuery(), operator)
}

type whereInWriter struct {
	baseWhereComparisonWriter
}

func (w *whereInWriter) WhereQuery() string {
	var operator string
	switch w.op {
	case op.In:
		operator = "IN"
	case op.NotIn:
		operator = "NOT IN"
	default:
		return ""
	}
	return fmt.Sprintf(`%s %s (?)`, w.ColumnQuery(), operator)
}
