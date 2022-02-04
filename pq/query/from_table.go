package query

import (
	"fmt"
)

type tableWriter struct {
	tableName string
	alias     string
}

func (s *tableWriter) FromQuery() string {
	q := fmt.Sprintf(`"%s"`, s.tableName)

	if s.alias != "" {
		q += fmt.Sprintf(` AS "%s"`, s.alias)
	}

	return q
}

func (s *tableWriter) SetAlias(alias string) {
	s.alias = alias
}
