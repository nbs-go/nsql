package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
	opt "github.com/nbs-go/nsql/query/option"
	"github.com/nbs-go/nsql/schema"
)

// whereCompareWriter

func newWhereComparisonWriter(s *schema.Schema, col string, operator op.Operator, args []interface{}) *whereCompareWriter {
	if !s.IsColumnExist(col) {
		panic(fmt.Errorf(`"columnWriter "%s" is not available in table "%s"`, col, s.TableName()))
	}

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
		case op.In, op.NotIn:
			v = new(inBindVar)
		}
	}

	return &whereCompareWriter{
		ColumnWriter: &columnWriter{
			name:      col,
			tableName: s.TableName(),
		},
		op:       operator,
		variable: v,
	}
}

//func newWhereNamedCompareWriter(s *schema.Schema, col string, operator op.Operator, args []interface{}) *whereCompareWriter {
//	if !s.IsColumnExist(col) {
//		panic(fmt.Errorf(`"columnWriter "%s" is not available in table "%s"`, col, s.TableName()))
//	}
//
//	// Evaluate options
//	opts := opt.EvaluateOptions(args)
//
//	// Get column format
//	colFmt, ok := opts.GetColumnFormat()
//	if !ok {
//		colFmt = query.ColumnOnly
//	}
//
//	return &whereCompareWriter{
//		ColumnWriter: &columnWriter{
//			name:      col,
//			tableName: s.TableName(),
//			format:    colFmt,
//		},
//		op: operator,
//		variable: &namedVar{
//			column: col,
//		},
//	}
//}

func Equal(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.Equal, args)
}

//func EqualNamed(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
//	return newWhereNamedCompareWriter(s, col, op.Equal, args)
//}

func NotEqual(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.NotEqual, args)
}

func GreaterThan(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.GreaterThan, args)
}

func GreaterThanEqual(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.GreaterThanEqual, args)
}

func LessThan(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.LessThan, args)
}

func LessThanEqual(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.LessThanEqual, args)
}

func Like(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.Like, args)
}

func NotLike(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.NotLike, args)
}

func ILike(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.ILike, args)
}

func NotILike(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.NotILike, args)
}

func Between(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.Between, args)
}

func NotBetween(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.NotBetween, args)
}

func In(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.In, args)
}

func NotIn(s *schema.Schema, col string, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(s, col, op.NotIn, args)
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
