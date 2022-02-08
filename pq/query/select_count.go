package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	opt "github.com/nbs-go/nsql/query/option"
)

func Count(column string, args ...interface{}) *selectCountWriter {
	// Get options
	opts := opt.EvaluateOptions(args)
	s := opts.GetSchema()
	as, _ := opts.GetString(opt.AsKey)

	// If all column, then create
	var allColumn bool
	var cw query.ColumnWriter
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
	query.ColumnWriter
	as        string
	allColumn bool
}

func (s *selectCountWriter) IsAllColumns() bool {
	return s.allColumn
}

func (s *selectCountWriter) SelectQuery() string {
	var q string

	if s.allColumn {
		q = "COUNT(*)"
	} else {
		q = fmt.Sprintf(`COUNT(%s)`, s.ColumnQuery())
	}

	// Set "as" query
	if s.as != "" {
		q += fmt.Sprintf(` AS "%s"`, s.as)
	}
	return q
}

func (s *selectCountWriter) SetFormat(_ query.ColumnFormat) {}
