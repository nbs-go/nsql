package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	opt "github.com/nbs-go/nsql/query/option"
	"github.com/nbs-go/nsql/schema"
	"strings"
)

type UpdateBuilder struct {
	schema  *schema.Schema
	columns []string
	where   query.WhereWriter
}

func (b *UpdateBuilder) Build(args ...interface{}) string {
	// If no column is defined, then panic
	count := len(b.columns)
	if count == 0 {
		panic(fmt.Errorf(`"no column defined on update table "%s"`, b.schema.TableName()))
	}

	// Get variable format option
	opts := opt.EvaluateOptions(args)
	format, ok := opts.GetVariableFormat()
	if !ok {
		// If var format is not defined, then set default to query.NamedVar
		format = query.NamedVar
	}

	// Set variable format in conditions
	if b.where == nil {
		// Set where to id
		b.where = Equal(b.schema, b.schema.PrimaryKey())
	}

	// Set format in conditions
	setUpdateFormat(b.where, b.schema, format)

	// Write assignments queries
	// TODO: Refactor as AssignmentsWriter query
	assignmentQueries := make([]string, count)
	for i, v := range b.columns {
		var q string
		switch format {
		case query.BindVar:
			q = fmt.Sprintf(`"%s" = ?`, v)
		case query.NamedVar:
			q = fmt.Sprintf(`"%s" = :%s`, v, v)
		}
		assignmentQueries[i] = q
	}
	assignments := strings.Join(assignmentQueries, query.Separator)

	// Write where
	where := b.where.WhereQuery()

	return fmt.Sprintf(`UPDATE "%s" SET %s WHERE %s`, b.schema.TableName(), assignments, where)
}

func (b *UpdateBuilder) Where(w query.WhereWriter) *UpdateBuilder {
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

func setUpdateFormat(ww query.WhereWriter, s *schema.Schema, format query.VariableFormat) {
	switch w := ww.(type) {
	case query.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			setUpdateFormat(cw, s, format)
		}
	case query.WhereCompareWriter:
		// Get column
		cw, ok := w.(query.ColumnWriter)
		if !ok {
			panic(fmt.Errorf("update condition did not implement query.ColumnWriter"))
		}

		// Check if column is not part if schema
		col := cw.GetColumn()
		if !s.IsColumnExist(cw.GetColumn()) {
			panic(fmt.Errorf(`"invalid column "%s" is not defined in Schema "%s"`, cw.GetColumn(), s.TableName()))
		}

		// Set format
		cw.SetFormat(query.ColumnOnly)

		// Set variable format
		switch format {
		case query.BindVar:
			w.SetVariable(new(bindVar))
		case query.NamedVar:
			w.SetVariable(&namedVar{column: col})
		}
	}
}
