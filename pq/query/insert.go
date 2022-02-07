package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

type InsertBuilder struct {
	tableName string
	columns   []string
	format    query.ColumnFormat
}

func (b *InsertBuilder) Build() string {
	// Write columns
	count := len(b.columns)
	columnQueries := make([]string, count)
	valueQueries := make([]string, count)
	for i, v := range b.columns {
		columnQueries[i] = fmt.Sprintf(`"%s"`, v)
		valueQueries[i] = fmt.Sprintf(`:%s`, v)
	}

	// Write values
	columns := strings.Join(columnQueries, query.Separator)
	values := strings.Join(valueQueries, query.Separator)

	return fmt.Sprintf(`INSERT INTO "%s"(%s) VALUES (%s)`, b.tableName, columns, values)
}

func Insert(s *schema.Schema, column string, columnN ...string) *InsertBuilder {
	// Init builder
	b := InsertBuilder{
		tableName: s.TableName(),
	}

	var columns []string
	if column == AllColumns {
		// Get all columns
		columns = s.InsertColumns()
	} else {
		inColumns := append([]string{column}, columnN...)
		for _, c := range inColumns {
			if s.IsColumnExist(c) {
				columns = append(columns, c)
			}
		}
	}

	// Set columns
	b.columns = columns

	return &b
}
