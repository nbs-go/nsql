package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

// columnWriter implement query.SelectFieldWriter to write query for selected field in a schema
type columnWriter struct {
	tableName string
	columns   []string
	mode      query.ColumnMode
}

func (w *columnWriter) SelectQuery() string {
	// Create query
	fieldQueries := make([]string, len(w.columns))
	for i, c := range w.columns {
		var q string

		switch w.mode {
		case query.JoinColumn:
			q = fmt.Sprintf(`"%s"."%s" AS "%s.%s"`, w.tableName, c, w.tableName, c)
		default:
			q = fmt.Sprintf(`"%s"."%s"`, w.tableName, c)
		}

		fieldQueries[i] = q
	}

	// Generate query
	return strings.Join(fieldQueries, query.Separator)
}

func (w *columnWriter) SetTableAlias(alias string) {
	w.tableName = alias
}

func (w *columnWriter) SetMode(mode query.ColumnMode) {
	w.mode = mode
}

func (w *columnWriter) GetTableName() string {
	return w.tableName
}

func newColumnWriter(s *schema.Schema, columns []string) *columnWriter {
	if isAllField(columns) {
		columns = s.GetColumns()
	} else {
		// Only contains columns in schema
		for i, col := range columns {
			if !s.IsColumnExist(col) {
				// Remove from array
				columns = append(columns[:i], columns[i+1:]...)
			}
		}
	}

	return &columnWriter{
		tableName: s.TableName,
		columns:   columns,
	}
}

func isAllField(f []string) bool {
	return len(f) == 1 && f[0] == "*"
}
