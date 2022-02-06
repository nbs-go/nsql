package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
	"github.com/nbs-go/nsql/schema"
)

// whereComparisonWriter

func newWhereComparisonWriter(s *schema.Schema, col string, op op.Operator) *whereComparisonWriter {
	if !s.IsColumnExist(col) {
		panic(fmt.Errorf(`"columnWriter "%s" is not available in table "%s"`, col, s.TableName))
	}

	return &whereComparisonWriter{
		baseWhereComparisonWriter{
			ColumnWriter: &columnWriter{
				name:      col,
				tableName: s.TableName,
			},
			op: op,
		}}
}

func Equal(s *schema.Schema, col string) *whereComparisonWriter {
	return newWhereComparisonWriter(s, col, op.Equal)
}

func NotEqual(s *schema.Schema, col string) *whereComparisonWriter {
	return newWhereComparisonWriter(s, col, op.NotEqual)
}

func GreaterThan(s *schema.Schema, col string) *whereComparisonWriter {
	return newWhereComparisonWriter(s, col, op.GreaterThan)
}

func GreaterThanEqual(s *schema.Schema, col string) *whereComparisonWriter {
	return newWhereComparisonWriter(s, col, op.GreaterThanEqual)
}

func LessThan(s *schema.Schema, col string) *whereComparisonWriter {
	return newWhereComparisonWriter(s, col, op.LessThan)
}

func LessThanEqual(s *schema.Schema, col string) *whereComparisonWriter {
	return newWhereComparisonWriter(s, col, op.LessThanEqual)
}

func Like(s *schema.Schema, col string) *whereComparisonWriter {
	return newWhereComparisonWriter(s, col, op.Like)
}

func NotLike(s *schema.Schema, col string) *whereComparisonWriter {
	return newWhereComparisonWriter(s, col, op.NotLike)
}

func ILike(s *schema.Schema, col string) *whereComparisonWriter {
	return newWhereComparisonWriter(s, col, op.ILike)
}

func NotILike(s *schema.Schema, col string) *whereComparisonWriter {
	return newWhereComparisonWriter(s, col, op.NotILike)
}

// whereBetweenWriter

func newWhereBetweenWriter(s *schema.Schema, col string, op op.Operator) *whereBetweenWriter {
	if !s.IsColumnExist(col) {
		panic(fmt.Errorf(`"columnWriter "%s" is not available in table "%s"`, col, s.TableName))
	}

	return &whereBetweenWriter{
		baseWhereComparisonWriter{
			ColumnWriter: &columnWriter{
				name:      col,
				tableName: s.TableName,
			},
			op: op,
		}}
}

func Between(s *schema.Schema, col string) *whereBetweenWriter {
	return newWhereBetweenWriter(s, col, op.Between)
}

func NotBetween(s *schema.Schema, col string) *whereBetweenWriter {
	return newWhereBetweenWriter(s, col, op.NotBetween)
}

// whereInWriter

func newWhereInWriter(s *schema.Schema, col string, op op.Operator) *whereInWriter {
	if !s.IsColumnExist(col) {
		panic(fmt.Errorf(`"columnWriter "%s" is not available in table "%s"`, col, s.TableName))
	}

	return &whereInWriter{
		baseWhereComparisonWriter{
			ColumnWriter: &columnWriter{
				name:      col,
				tableName: s.TableName,
			},
			op: op,
		}}
}

func In(s *schema.Schema, col string) *whereInWriter {
	return newWhereInWriter(s, col, op.In)
}

func NotIn(s *schema.Schema, col string) *whereInWriter {
	return newWhereInWriter(s, col, op.NotIn)
}

// whereLogicalWriter

func And(c1 query.WhereWriter, cn ...query.WhereWriter) *whereLogicalWriter {
	conditions := []query.WhereWriter{c1}
	if len(cn) > 0 {
		conditions = append(conditions, cn...)
	}

	return &whereLogicalWriter{
		op:         op.And,
		conditions: conditions,
	}
}

func Or(c1 query.WhereWriter, cn ...query.WhereWriter) *whereLogicalWriter {
	conditions := []query.WhereWriter{c1}
	if len(cn) > 0 {
		conditions = append(conditions, cn...)
	}

	return &whereLogicalWriter{
		op:         op.Or,
		conditions: conditions,
	}
}

func resolveFromOfWhereWriters(ww query.WhereWriter, from *schema.Schema) {
	// Switch type
	switch w := ww.(type) {
	case query.LogicalWhereWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			resolveFromOfWhereWriters(cw, from)
		}
	case query.ComparisonWhereWriter:
		// Get alias
		if w.GetTableName() == fromTableFlag {
			w.SetSchema(from)
		}
	}
}

func filterWhereWriters(ww query.WhereWriter, tables map[string]selectTable) query.WhereWriter {
	// Switch type
	switch w := ww.(type) {
	case query.LogicalWhereWriter:
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
	case query.ComparisonWhereWriter:
		// Check if condition is registered in table
		table, ok := tables[w.GetTableName()]
		if !ok {
			return nil
		}

		// Set alias
		w.SetTableAs(table.as)
	}
	return ww
}
