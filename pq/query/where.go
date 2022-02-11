package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
	opt "github.com/nbs-go/nsql/query/option"
	"github.com/nbs-go/nsql/schema"
)

// whereCompareWriter

func newWhereComparisonWriter(col query.ColumnWriter, operator op.Operator, args []interface{}) *whereCompareWriter {
	opts := opt.EvaluateOptions(args)

	// Get variable writer
	v := opts.GetVariable(opt.VariableKey)
	if v == nil {
		// Set default variable writer, by operator
		switch operator {
		case op.Equal, op.NotEqual, op.GreaterThan, op.GreaterThanEqual, op.LessThan, op.LessThanEqual, op.Like,
			op.NotLike, op.ILike, op.NotILike:
			v = new(bindVar)
		case op.Between, op.NotBetween:
			v = new(betweenBindVar)
		}
	}

	// Get alias
	as, _ := opts.GetString(opt.AsKey)

	return &whereCompareWriter{
		ColumnWriter: col,
		op:           operator,
		variable:     v,
		as:           as,
	}
}

func newInWhereComparisonWriter(col query.ColumnWriter, argCount int, operator op.Operator, args []interface{}) *whereCompareWriter {
	opts := opt.EvaluateOptions(args)

	// Get variable writer
	v := opts.GetVariable(opt.VariableKey)
	if v == nil {
		// Set default variable writer, by operator
		switch operator {
		case op.In, op.NotIn:
			v = &inBindVar{argCount: argCount}
		}
	}

	// Get alias
	as, _ := opts.GetString(opt.AsKey)

	return &whereCompareWriter{
		ColumnWriter: col,
		op:           operator,
		variable:     v,
		as:           as,
	}
}

func Equal(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.Equal, args)
}

func NotEqual(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotEqual, args)
}

func GreaterThan(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.GreaterThan, args)
}

func GreaterThanEqual(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.GreaterThanEqual, args)
}

func LessThan(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.LessThan, args)
}

func LessThanEqual(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.LessThanEqual, args)
}

func Like(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.Like, args)
}

func NotLike(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotLike, args)
}

func ILike(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.ILike, args)
}

func NotILike(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotILike, args)
}

func Between(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.Between, args)
}

func NotBetween(col query.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotBetween, args)
}

func In(col query.ColumnWriter, argCount int, args ...interface{}) *whereCompareWriter {
	return newInWhereComparisonWriter(col, argCount, op.In, args)
}

func NotIn(col query.ColumnWriter, argCount int, args ...interface{}) *whereCompareWriter {
	return newInWhereComparisonWriter(col, argCount, op.NotIn, args)
}

// whereLogicWriter

func newWhereLogicalWriter(operator op.Operator, c1 query.WhereWriter, cn ...query.WhereWriter) *whereLogicWriter {
	conditions := []query.WhereWriter{c1}
	if len(cn) > 0 {
		conditions = append(conditions, cn...)
	}

	return &whereLogicWriter{
		op:         operator,
		conditions: conditions,
	}
}

func And(c1 query.WhereWriter, cn ...query.WhereWriter) *whereLogicWriter {
	return newWhereLogicalWriter(op.And, c1, cn...)
}

func Or(c1 query.WhereWriter, cn ...query.WhereWriter) *whereLogicWriter {
	return newWhereLogicalWriter(op.Or, c1, cn...)
}

func resolveFromTableFlag(ww query.WhereWriter, from *schema.Schema) {
	// Switch type
	switch w := ww.(type) {
	case query.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			resolveFromTableFlag(cw, from)
		}
	case query.WhereCompareWriter:
		// Get alias
		if w.GetTableName() == fromTableFlag {
			w.SetSchema(from)
		}
	}
}

func filterWhereWriters(ww query.WhereWriter, tables map[string]query.Table) query.WhereWriter {
	// Switch type
	switch w := ww.(type) {
	case query.WhereLogicWriter:
		// Get conditions
		var conditions []query.WhereWriter
		for _, cw := range w.GetConditions() {
			c := filterWhereWriters(cw, tables)
			// If no writer is set, then delete from array
			if c == nil {
				continue
			}
			conditions = append(conditions, c)
		}
		// Update conditions
		w.SetConditions(conditions)
	case query.WhereCompareWriter:
		// Check if condition is registered in table
		table, ok := tables[w.GetTableName()]
		if !ok {
			return nil
		}

		// Set alias
		w.SetTableAs(table.As)
	}
	return ww
}

func resolveJoinTableFlag(ww query.WhereWriter, joinTable *schema.Schema) {
	// Switch type
	switch w := ww.(type) {
	case query.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			resolveJoinTableFlag(cw, joinTable)
		}
	case query.WhereCompareWriter:
		// Get variable
		v := w.GetVariable()

		// Cast as column
		cv, ok := v.(query.ColumnWriter)
		if !ok {
			return
		}

		// If join table flag reference is set, then set schema
		if cv.GetTableName() == joinTableFlag {
			cv.SetSchema(joinTable)
		}
	}
}

func setJoinTableAs(ww query.WhereWriter, joinTable *query.Table, tableRefs map[string]query.Table) {
	switch w := ww.(type) {
	case query.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			setJoinTableAs(cw, joinTable, tableRefs)
		}
	case query.WhereCompareWriter:
		// Get column
		if cw, ok := w.(query.ColumnWriter); ok {
			setJoinColumnTableAs(cw, joinTable, tableRefs)
		}

		// Get variable and set table alias
		v := w.GetVariable()
		if cv, ok := v.(query.ColumnWriter); ok {
			setJoinColumnTableAs(cv, joinTable, tableRefs)
		}
	}
}

func setJoinColumnTableAs(column query.ColumnWriter, joinTable *query.Table, tableRefs map[string]query.Table) {
	// Get arguments
	tableName := column.GetTableName()
	col := column.GetColumn()

	// Check in joinSchema
	if tableName == joinTable.Schema.TableName() {
		if !joinTable.Schema.IsColumnExist(col) {
			panic(fmt.Errorf(`column "%s" is not declared in Table "%s"`, col, joinTable.Schema.TableName()))
		}
		// Set alias
		column.SetTableAs(joinTable.As)
		return
	}

	// Check in tableRefs
	tRef, tRefOk := tableRefs[tableName]
	if !tRefOk {
		panic(fmt.Errorf(`table "%s" is not declared in Query Builder`, tableName))
	}
	// Check against column
	if !tRef.Schema.IsColumnExist(col) {
		panic(fmt.Errorf(`column "%s" is not declared in Table "%s"`, tableName, col))
	}
	column.SetTableAs(tRef.As)
}
