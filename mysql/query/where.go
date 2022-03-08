package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/schema"
)

// whereCompareWriter

func newWhereComparisonWriter(col nsql.ColumnWriter, operator op.Operator, args []interface{}) *whereCompareWriter {
	opts := option.EvaluateOptions(args)

	// Get variable writer
	v := opts.GetVariable(option.VariableKey)
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
	as, _ := opts.GetString(option.AsKey)

	return &whereCompareWriter{
		ColumnWriter: col,
		op:           operator,
		variable:     v,
		as:           as,
	}
}

func newInWhereComparisonWriter(col nsql.ColumnWriter, argCount int, operator op.Operator, args []interface{}) *whereCompareWriter {
	opts := option.EvaluateOptions(args)

	// Get variable writer
	v := opts.GetVariable(option.VariableKey)
	if v == nil {
		// Set default variable writer, by operator
		switch operator {
		case op.In, op.NotIn:
			v = &inBindVar{argCount: argCount}
		}
	}

	// Get alias
	as, _ := opts.GetString(option.AsKey)

	return &whereCompareWriter{
		ColumnWriter: col,
		op:           operator,
		variable:     v,
		as:           as,
	}
}

func Equal(col nsql.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.Equal, args)
}

func NotEqual(col nsql.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotEqual, args)
}

func GreaterThan(col nsql.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.GreaterThan, args)
}

func GreaterThanEqual(col nsql.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.GreaterThanEqual, args)
}

func LessThan(col nsql.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.LessThan, args)
}

func LessThanEqual(col nsql.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.LessThanEqual, args)
}

func Like(col nsql.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.Like, args)
}

func NotLike(col nsql.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotLike, args)
}

func Between(col nsql.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.Between, args)
}

func NotBetween(col nsql.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotBetween, args)
}

func In(col nsql.ColumnWriter, argCount int, args ...interface{}) *whereCompareWriter {
	return newInWhereComparisonWriter(col, argCount, op.In, args)
}

func NotIn(col nsql.ColumnWriter, argCount int, args ...interface{}) *whereCompareWriter {
	return newInWhereComparisonWriter(col, argCount, op.NotIn, args)
}

// whereLogicWriter

func newWhereLogicalWriter(operator op.Operator, cn []nsql.WhereWriter) *whereLogicWriter {
	return &whereLogicWriter{
		op:         operator,
		conditions: cn,
	}
}

func And(cn ...nsql.WhereWriter) *whereLogicWriter {
	return newWhereLogicalWriter(op.And, cn)
}

func Or(cn ...nsql.WhereWriter) *whereLogicWriter {
	return newWhereLogicalWriter(op.Or, cn)
}

func resolveFromTableFlag(ww nsql.WhereWriter, from *schema.Schema) {
	// Switch type
	switch w := ww.(type) {
	case nsql.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			resolveFromTableFlag(cw, from)
		}
	case nsql.WhereCompareWriter:
		// Get alias
		if w.GetTableName() == fromTableFlag {
			w.SetSchema(from)
		}
	}
}

func filterWhereWriters(ww nsql.WhereWriter, tables map[string]nsql.Table) nsql.WhereWriter {
	// Switch type
	switch w := ww.(type) {
	case nsql.WhereLogicWriter:
		// Get conditions
		var conditions []nsql.WhereWriter
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
	case nsql.WhereCompareWriter:
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

func resolveJoinTableFlag(ww nsql.WhereWriter, joinTable *schema.Schema) {
	// Switch type
	switch w := ww.(type) {
	case nsql.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			resolveJoinTableFlag(cw, joinTable)
		}
	case nsql.WhereCompareWriter:
		// Get variable
		v := w.GetVariable()

		// Cast as column
		cv, ok := v.(nsql.ColumnWriter)
		if !ok {
			return
		}

		// If join table flag reference is set, then set schema
		if cv.GetTableName() == joinTableFlag {
			cv.SetSchema(joinTable)
		}
	}
}

func setJoinTableAs(ww nsql.WhereWriter, joinTable *nsql.Table, tableRefs map[string]nsql.Table) {
	switch w := ww.(type) {
	case nsql.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			setJoinTableAs(cw, joinTable, tableRefs)
		}
	case nsql.WhereCompareWriter:
		// Get column
		if cw, ok := w.(nsql.ColumnWriter); ok {
			setJoinColumnTableAs(cw, joinTable, tableRefs)
		}

		// Get variable and set table alias
		v := w.GetVariable()
		if cv, ok := v.(nsql.ColumnWriter); ok {
			setJoinColumnTableAs(cv, joinTable, tableRefs)
		}
	}
}

func setJoinColumnTableAs(column nsql.ColumnWriter, joinTable *nsql.Table, tableRefs map[string]nsql.Table) {
	// Get arguments
	tableName := column.GetTableName()

	if tableName == skipTableFlag {
		return
	}

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
