package query

import (
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

// columnSchemaWriter implements query.SelectWriter for columnSchemaWriter in a Schema
type columnSchemaWriter struct {
	schema    *schema.Schema
	columns   []string
	tableName string
	format    op.ColumnFormat
}

func (w *columnSchemaWriter) SetSchema(s *schema.Schema) {
	w.schema = s
	w.tableName = s.TableName()
}

func (w *columnSchemaWriter) GetTableName() string {
	if w.schema != nil {
		return w.schema.TableName()
	}
	return w.tableName
}

func (w *columnSchemaWriter) SelectQuery() string {
	// Create query
	var queries []string
	for _, col := range w.columns {
		// Skip if columnWriter is not set
		if !w.schema.IsColumnExist(col) {
			continue
		}

		// Create columnWriter
		q := writeColumn(w.tableName, col, w.format)
		queries = append(queries, q)
	}
	// Generate query
	return strings.Join(queries, nsql.Separator)
}

func (w *columnSchemaWriter) SetFormat(format op.ColumnFormat) {
	w.format = format
}

func (w *columnSchemaWriter) SetTableAs(as string) {
	if as == "" {
		return
	}
	w.tableName = as
}

func (w *columnSchemaWriter) IsAllColumns() bool {
	return false
}
