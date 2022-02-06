package query

import (
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

// columnSchemaWriter implements query.SelectWriter for columnSchemaWriter in a Schema
type columnSchemaWriter struct {
	schema    *schema.Schema
	columns   []string
	tableName string
	format    query.ColumnFormat
}

func (w *columnSchemaWriter) SetSchema(s *schema.Schema) {
	w.schema = s
	w.tableName = s.TableName

}
func (w *columnSchemaWriter) GetTableName() string {
	if w.schema != nil {
		return w.schema.TableName
	}
	return w.tableName
}

func (w *columnSchemaWriter) SelectQuery() string {
	// Use schema table name if not set
	if w.tableName == "" {
		w.tableName = w.schema.TableName
	}

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
	return strings.Join(queries, query.Separator)
}

func (w *columnSchemaWriter) SetFormat(format query.ColumnFormat) {
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
