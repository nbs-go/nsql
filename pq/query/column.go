package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/schema"
)

func Column(col string, args ...interface{}) *columnWriter {
	// Evaluate options
	opts := option.EvaluateOptions(args)

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
		format = op.NonAmbiguousColumn
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
	optCopy := option.NewOptions()
	for _, v := range args {
		switch cv := v.(type) {
		case option.SetOptionFn:
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
		format = op.NonAmbiguousColumn
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
	tableAs   string
	format    op.ColumnFormat
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

func (w *columnWriter) Expand(args ...interface{}) nsql.SelectWriter {
	// Get schema
	opts := option.EvaluateOptions(args)
	s := opts.GetSchema()

	// Expand to columnWriter schema
	return &columnSchemaWriter{
		schema:    s,
		columns:   s.Columns(),
		tableName: s.TableName(),
	}
}

func (w *columnWriter) ColumnQuery() string {
	return writeColumn(w.tableName, w.tableAs, w.name, w.format)
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
	w.tableAs = as
}

func (w *columnWriter) SetFormat(format op.ColumnFormat) {
	w.format = format
}

func writeColumn(tableName string, tableAs string, name string, format op.ColumnFormat) string {
	// Set table alias
	if tableAs != "" {
		tableName = tableAs
	}

	switch format {
	case op.SelectJoinColumn:
		return fmt.Sprintf(`"%s"."%s" AS "%s.%s"`, tableName, name, tableName, name)
	case op.ColumnOnly:
		return fmt.Sprintf(`"%s"`, name)
	default:
		// If not set, treat as NonAmbiguous column
		return fmt.Sprintf(`"%s"."%s"`, tableName, name)
	}
}
