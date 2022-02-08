package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	opt "github.com/nbs-go/nsql/query/option"
	"github.com/nbs-go/nsql/schema"
)

func Column(col string, args ...interface{}) *columnWriter {
	// Evaluate options
	opts := opt.EvaluateOptions(args)

	// Get tableName
	s := opts.GetSchema()
	var tableName string
	if s == nil {
		tableName = fromTableFlag
	} else {
		tableName = s.TableName()
	}

	// Get format
	format, ok := opts.GetColumnFormat()
	if !ok {
		format = query.NonAmbiguousColumn
	}

	return &columnWriter{
		name:      col,
		tableName: tableName,
		format:    format,
	}
}

func Columns(column1, column2 string, args ...interface{}) *columnSchemaWriter {
	// Init columns containers
	var inColumns []string

	// Evaluate arguments
	optCopy := opt.NewOptions()
	for _, v := range args {
		switch cv := v.(type) {
		case opt.SetOptionFn:
			cv(optCopy)
		case string:
			inColumns = append(inColumns, cv)
		}
	}

	// Get schema
	s := optCopy.GetSchema()
	var cols []string
	var tableName string
	if s == nil {
		tableName = fromTableFlag
		cols = append([]string{column1, column2}, inColumns...)
	} else {
		tableName = s.TableName()
		inColumns = append([]string{column2}, inColumns...)
		cols = s.Filter(column1, inColumns...)
	}

	// Get format
	format, ok := optCopy.GetColumnFormat()
	if !ok {
		format = query.NonAmbiguousColumn
	}

	// Create schema writer
	return &columnSchemaWriter{
		schema:    s,
		columns:   cols,
		tableName: tableName,
		format:    format,
	}
}

// columnWriter implements query.SelectWriter for a single columnWriter
type columnWriter struct {
	name      string
	tableName string
	format    query.ColumnFormat
}

func (w *columnWriter) VariableQuery() string {
	return w.ColumnQuery()
}

func (w *columnWriter) GetColumn() string {
	return w.name
}

func (w *columnWriter) SetSchema(s *schema.Schema) {
	// Check if column is part of schema
	if w.name != AllColumns && !s.IsColumnExist(w.name) {
		// Mark column to be skipped from writer
		w.tableName = skipTableFlag
		return
	}
	w.tableName = s.TableName()
}

func (w *columnWriter) Expand(args ...interface{}) query.SelectWriter {
	// Get schema
	opts := opt.EvaluateOptions(args)
	s := opts.GetSchema()

	// Expand to columnWriter schema
	return &columnSchemaWriter{
		schema:    s,
		columns:   s.Columns(),
		tableName: s.TableName(),
	}
}

func (w *columnWriter) ColumnQuery() string {
	return writeColumn(w.tableName, w.name, w.format)
}

func (w *columnWriter) IsAllColumns() bool {
	return w.name == AllColumns
}

func (w *columnWriter) SelectQuery() string {
	return w.ColumnQuery()
}

func (w *columnWriter) GetTableName() string {
	return w.tableName
}

func (w *columnWriter) SetTableAs(as string) {
	// Ignore if alias is empty
	if as == "" {
		return
	}
	w.tableName = as
}

func (w *columnWriter) SetFormat(format query.ColumnFormat) {
	w.format = format
}

func writeColumn(tableName string, name string, format query.ColumnFormat) string {
	switch format {
	case query.SelectJoinColumn:
		return fmt.Sprintf(`"%s"."%s" AS "%s.%s"`, tableName, name, tableName, name)
	case query.ColumnOnly:
		return fmt.Sprintf(`"%s"`, name)
	default:
		// If not set, treat as NonAmbiguous column
		return fmt.Sprintf(`"%s"."%s"`, tableName, name)
	}
}
