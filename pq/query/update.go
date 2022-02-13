package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

type UpdateBuilder struct {
	schema  *schema.Schema
	columns []string
	where   nsql.WhereWriter
}

func (b *UpdateBuilder) Build(args ...interface{}) string {
	// If no column is defined, then panic
	count := len(b.columns)
	if count == 0 {
		panic(fmt.Errorf(`"no column defined on update table "%s"`, b.schema.TableName()))
	}

	// Get variable format option
	opts := option.EvaluateOptions(args)
	format, ok := opts.GetVariableFormat()
	if !ok {
		// If var format is not defined, then set default to query.NamedVar
		format = op.NamedVar
	}

	// Set variable format in conditions
	if b.where == nil {
		// Set where to id
		b.where = Equal(Column(b.schema.PrimaryKey(), b.schema))
	}

	// Set format in conditions
	setUpdateFormat(b.where, b.schema, format)

	// Write assignments queries
	// TODO: Refactor as AssignmentsWriter query
	assignmentQueries := make([]string, count)
	for i, v := range b.columns {
		var q string
		switch format {
		case op.BindVar:
			q = fmt.Sprintf(`"%s" = ?`, v)
		case op.NamedVar:
			q = fmt.Sprintf(`"%s" = :%s`, v, v)
		}
		assignmentQueries[i] = q
	}
	assignments := strings.Join(assignmentQueries, nsql.Separator)

	// Write where
	where := b.where.WhereQuery()

	return fmt.Sprintf(`UPDATE "%s" SET %s WHERE %s`, b.schema.TableName(), assignments, where)
}

func (b *UpdateBuilder) Where(w nsql.WhereWriter) *UpdateBuilder {
	b.where = w
	return b
}

func Update(s *schema.Schema, column string, columnN ...string) *UpdateBuilder {
	// Init builder
	b := UpdateBuilder{
		schema: s,
	}

	var columns []string
	if column == AllColumns {
		// Get all columns
		columns = s.UpdateColumns()
	} else {
		inColumns := append([]string{column}, columnN...)
		pk := s.PrimaryKey()
		for _, c := range inColumns {
			if s.IsColumnExist(c) && c != pk {
				columns = append(columns, c)
			}
		}
	}

	// Set columns
	b.columns = columns

	return &b
}

func setUpdateFormat(ww nsql.WhereWriter, s *schema.Schema, format op.VariableFormat) {
	switch w := ww.(type) {
	case nsql.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			setUpdateFormat(cw, s, format)
		}
	case nsql.WhereCompareWriter:
		// Get column
		cw, ok := w.(nsql.ColumnWriter)
		if !ok {
			panic(fmt.Errorf("update condition did not implement query.ColumnWriter"))
		}

		// Check if column is not part if schema
		col := cw.GetColumn()
		if !s.IsColumnExist(cw.GetColumn()) {
			panic(fmt.Errorf(`"invalid column "%s" is not defined in Schema "%s"`, cw.GetColumn(), s.TableName()))
		}

		// Set format
		cw.SetFormat(op.ColumnOnly)

		// Set variable format
		switch format {
		case op.BindVar:
			w.SetVariable(new(bindVar))
		case op.NamedVar:
			w.SetVariable(&namedVar{column: col})
		}
	}
}
