package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query/op"
)

type baseWhereComparisonWriter struct {
	tableName string
	column    string
	op        op.Operator
}

func (w *baseWhereComparisonWriter) SetTableAlias(alias string) {
	w.tableName = alias
}

func (w *baseWhereComparisonWriter) GetTableName() string {
	return w.tableName
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
	return fmt.Sprintf(`"%s"."%s" %s ?`, w.tableName, w.column, operator)
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
	return fmt.Sprintf(`"%s"."%s" %s ? AND ?`, w.tableName, w.column, operator)
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
	return fmt.Sprintf(`"%s"."%s" %s (?)`, w.tableName, w.column, operator)
}
