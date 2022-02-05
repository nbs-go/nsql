package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
)

type selectCountWriter struct {
	column    string
	tableName string
	as        string
}

func (s *selectCountWriter) SelectQuery() string {
	var q string

	if s.column == "*" {
		q = "COUNT(*)"
	} else {
		q = fmt.Sprintf(`COUNT("%s"."%s")`, s.tableName, s.column)
	}

	// Set "as" query
	if s.as != "" {
		q += fmt.Sprintf(` AS "%s"`, s.as)
	}
	return q
}

func (s *selectCountWriter) GetTableName() string {
	return s.tableName
}

func (s *selectCountWriter) SetTableAlias(alias string) {
	s.tableName = alias
}

func (s *selectCountWriter) SetMode(_ query.ColumnMode) {}
