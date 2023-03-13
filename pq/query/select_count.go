package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
)

func Count(column string, args ...interface{}) *selectCountWriter {
	// Get options
	opts := option.EvaluateOptions(args)
	s := opts.GetSchema()
	as, _ := opts.GetString(option.AsKey)

	// If all column, then create
	var allColumn bool
	var cw nsql.ColumnWriter
	if column == AllColumns {
		cw = &columnWriter{
			name:      column,
			tableName: forceWriteFlag,
		}
		allColumn = true
	} else if s != nil {
		// If column is invalid or not in schema, then skip
		if !s.IsColumnExist(column) {
			panic(fmt.Errorf(`column "%s" is not declared in schema "%s"`, column, s.TableName()))
		}

		// Set writer
		cw = &columnWriter{
			name:      column,
			tableName: s.TableName(),
			tableAs:   s.As(),
		}
	} else {
		// Set writer with FROM flag
		cw = &columnWriter{
			name:      column,
			tableName: fromTableFlag,
		}
	}

	return &selectCountWriter{ColumnWriter: cw, as: as, allColumn: allColumn}
}

type selectCountWriter struct {
	nsql.ColumnWriter
	as        string
	allColumn bool
}

func (s *selectCountWriter) IsAllColumns() bool {
	return s.allColumn
}

func (s *selectCountWriter) ColumnQuery() string {
	if s.allColumn {
		return "COUNT(*)"
	}
	// Print column
	col := s.ColumnWriter.ColumnQuery()
	return fmt.Sprintf(`COUNT(%s)`, col)
}

func (s *selectCountWriter) SelectQuery() string {
	q := s.ColumnQuery()

	// Set "as" query
	if s.as != "" {
		q += fmt.Sprintf(` AS "%s"`, s.as)
	}
	return q
}

func (s *selectCountWriter) SetFormat(_ op.ColumnFormat) {}
