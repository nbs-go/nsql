package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
)

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
