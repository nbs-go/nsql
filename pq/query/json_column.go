package query

import (
	"errors"
	"fmt"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

// JsonColumnWriter implements query.SelectWriter for a column that referred to a JSON
type JsonColumnWriter struct {
	name      string
	attrs     []string
	tableName string
	tableAs   string
	as        string
}

func (w *JsonColumnWriter) GetSchemaRef() schema.Reference {
	if w.tableAs != "" {
		return schema.Reference(w.tableAs)
	}
	return schema.Reference(w.tableName)
}

func JsonColumn(column string, args ...interface{}) *JsonColumnWriter {
	// If column does not contain ".", then panic
	tmp := strings.Split(column, ".")
	if len(tmp) < 2 {
		panic(errors.New(`nsql: invalid JsonColumn value, attributes is not defined`))
	}

	if tmp[1] == "" {
		panic(errors.New(`nsql: invalid JsonColumn value, attributes is not defined`))
	}

	// Set columns on the first index, the rest refer to attributes
	col, attrs := tmp[0], tmp[1:]

	// Evaluate options
	opts := option.EvaluateOptions(args)

	// Get tableName
	s := opts.GetSchema()
	var tableName, tableAs string
	if s == nil {
		tableName = fromTableFlag
	} else {
		tableName = s.TableName()
		tableAs = s.As()
	}

	// Get alias
	as, _ := opts.GetString(option.AsKey)

	return &JsonColumnWriter{
		name:      col,
		attrs:     attrs,
		tableName: tableName,
		tableAs:   tableAs,
		as:        as,
	}
}

func (w *JsonColumnWriter) GetColumn() string {
	return w.name
}

func (w *JsonColumnWriter) SetSchema(s *schema.Schema) {
	// Check if column is part of schema
	if !s.IsColumnExist(w.name) {
		// Mark column to be skipped from writer
		w.tableName = skipTableFlag
		return
	}
	w.tableName = s.TableName()
	w.tableAs = s.As()
}

func (w *JsonColumnWriter) ColumnQuery() string {
	return w.query()
}

func (w *JsonColumnWriter) IsAllColumns() bool {
	return false
}

func (w *JsonColumnWriter) SelectQuery() string {
	q := w.query()
	if w.as != "" {
		q = fmt.Sprintf(`(%s) AS "%s"`, q, w.as)
	}
	return q
}

func (w *JsonColumnWriter) GetTableName() string {
	return w.tableName
}

func (w *JsonColumnWriter) SetTableAs(as string) {
	// Ignore if alias is empty
	if as == "" {
		return
	}
	w.tableAs = as
}

func (w *JsonColumnWriter) SetFormat(_ op.ColumnFormat) {
	// Do nothing
}

// query Generate column query
func (w *JsonColumnWriter) query() string {
	// Set table alias
	tableName := w.tableName
	if w.tableAs != "" {
		tableName = w.tableAs
	}

	// Write attribute query
	var aq string
	attrsCount := len(w.attrs)
	if attrsCount == 1 {
		aq = "->>'" + w.attrs[0] + "'"
	} else {
		// Pop leaf from attrs
		leaf, attrs := w.attrs[attrsCount-1], w.attrs[:attrsCount-1]
		// Write attribute query
		aq = "->'" + strings.Join(attrs, "'->'") + "'->>'" + leaf + "'"
	}

	// Write column query
	return fmt.Sprintf(`"%s"."%s"%s`, tableName, w.name, aq)
}
