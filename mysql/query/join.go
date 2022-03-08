package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
)

type joinWriter struct {
	method      op.JoinMethod
	table       *nsql.Table
	onCondition nsql.WhereWriter
	index       int
}

func (j *joinWriter) GetIndex() int {
	return j.index
}

func (j *joinWriter) SetIndex(n int) {
	j.index = n
}

func (j *joinWriter) GetTableName() string {
	return j.table.Schema.TableName()
}

func (j *joinWriter) JoinQuery() string {
	// Get method query
	var method string
	switch j.method {
	case op.InnerJoin:
		method = "INNER JOIN"
	case op.RightJoin:
		method = "RIGHT JOIN"
	case op.FullJoin:
		method = "FULL OUTER JOIN"
	default:
		// Default to left join
		method = "LEFT JOIN"
	}

	// Generate table name
	table := j.table
	tableName := fmt.Sprintf("`%s`", table.Schema.TableName())
	if table.As != "" {
		tableName += fmt.Sprintf(" AS `%s`", table.As)
	}

	// Write condition
	condition := j.onCondition.WhereQuery()

	return fmt.Sprintf(`%s %s ON %s`, method, tableName, condition)
}
